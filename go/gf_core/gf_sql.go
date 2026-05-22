/*
MIT License

Copyright (c) 2023 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_core

import (
	"fmt"
	"net/url"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

//-------------------------------------------------

func validatePostgresURL(pURLstr string) (bool, string) {
	parsedURL, err := url.Parse(pURLstr)
	if err != nil {
		return false, fmt.Sprintf("invalid URL format: %s", err.Error())
	}

	// Check scheme
	if parsedURL.Scheme != "postgres" && parsedURL.Scheme != "postgresql" {
		return false, fmt.Sprintf("invalid scheme: expected 'postgres' or 'postgresql', got '%s'", parsedURL.Scheme)
	}

	// Check username
	if parsedURL.User == nil || parsedURL.User.Username() == "" {
		return false, "missing username"
	}

	// Check password
	if _, hasPassword := parsedURL.User.Password(); !hasPassword {
		return false, "missing password"
	}

	// Check host
	if parsedURL.Host == "" {
		return false, "missing host"
	}

	// Check database name (path without leading slash)
	dbName := parsedURL.Path
	if dbName == "" || dbName == "/" {
		return false, "missing database name"
	}

	// Check sslmode parameter
	sslMode := parsedURL.Query().Get("sslmode")
	if sslMode == "" {
		return false, "missing sslmode parameter"
	}

	validSSLModes := map[string]bool{
		"disable": true, "require": true, "verify-ca": true, "verify-full": true,
	}
	if !validSSLModes[sslMode] {
		return false, fmt.Sprintf("invalid sslmode: %s", sslMode)
	}

	return true, ""
}

//-------------------------------------------------

func DBsqlGetNullStringOrDefault(pNullableStr sql.NullString, pDefaultValStr string) string {
    if pNullableStr.Valid {
        return pNullableStr.String
    }
    return pDefaultValStr
}

//-------------------------------------------------

func DBsqlConnect(pDBnameStr string,
	pUserNameStr string,
	pPassStr     string,
	pDBhostStr   string,
	pSSLmodeStr  string,
	pRuntimeSys  *RuntimeSys) (*sql.DB, string, *GFerror) {

	// SSL mode is now configurable
	if pSSLmodeStr == "" {
		pSSLmodeStr = "disable" // default fallback
	}
	dbDSNuriStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		pUserNameStr,
		pPassStr,
		pDBhostStr,
		pDBnameStr,
		pSSLmodeStr)

	// Validate the DSN URL
	if valid, reasonStr := validatePostgresURL(dbDSNuriStr); !valid {
		gfErr := ErrorCreate(fmt.Sprintf("invalid Postgres DSN: %s", reasonStr),
			"generic_error",
			map[string]interface{}{
				"db_host_str": pDBhostStr,
				"db_name_str": pDBnameStr,
				"reason_str":  reasonStr,
			},
			nil, "gf_core", pRuntimeSys)
		return nil, "", gfErr
	}

	//-----------------------
	// CONNECT

	maxRetriesInt := 5
	retryIntervalSecsInt := 2

	var db *sql.DB
	var err error

	for retriesInt := 0; retriesInt < maxRetriesInt; retriesInt++ {

		pRuntimeSys.LogNewFun("INFO", fmt.Sprintf("attempt - connecting to SQL DB... host: %s, db: %s", pDBhostStr, pDBnameStr), nil)
		db, err = sql.Open("postgres", dbDSNuriStr)
		if err == nil {

			// test the connection
			err = db.Ping()
			if err == nil {
				break // successful connection
			}
		}

		if retriesInt < maxRetriesInt-1 {
			pRuntimeSys.LogNewFun("INFO", fmt.Sprintf("retry %d - SQL DB connect in %ds...", retriesInt, retryIntervalSecsInt), nil)
			time.Sleep(time.Duration(retryIntervalSecsInt) * time.Second)
		}
	}

	// if no connection was established even after retrying, return error
	if db == nil {
		gfErr := ErrorCreate("failed to connect to a SQL server at target url",
			"sql_failed_to_connect",
			map[string]interface{}{
				"db_host_str": pDBhostStr,
				"db_name_str": pDBnameStr,
			},
			nil, "gf_core", pRuntimeSys)
		return nil, "", gfErr
	}

	//-----------------------
	// CONFIGURE CONNECTION POOL
	//
	// Without these settings, the Go SQL driver can open unlimited connections to the database,
	// which can exhaust the available connection slots on the PostgreSQL server (especially on
	// AWS RDS where non-superuser connections are limited, typically to ~80-100 on small instances).
	//
	// These settings work together to prevent connection exhaustion while maintaining performance:

	// SetMaxOpenConns(25):
	// - Hard limit on the total number of open connections (both in-use and idle)
	// - When this limit is reached, new DB operations will BLOCK and wait for a connection to become available
	// - Runtime behavior:
	//   * If 25 connections are already open and a 26th query is attempted, it will wait
	//   * The wait has no timeout by default - it blocks until a connection is freed
	//   * This prevents the "remaining connection slots are reserved" error from PostgreSQL
	// - Tuning considerations:
	//   * Too low: requests will block waiting for connections, degrading performance
	//   * Too high: can exhaust DB server connection limit, especially with multiple app instances
	//   * Rule of thumb: (DB max_connections - reserved) / number_of_app_instances
	//   * For RDS with max_connections=100, reserved=22, and 2 instances: (100-22)/2 = 39 per instance
	//   * Set to 25 as a conservative default that works across different deployment sizes
	db.SetMaxOpenConns(25)

	// SetMaxIdleConns(5):
	// - Maximum number of connections kept open in the idle pool (not actively executing queries)
	// - Runtime behavior:
	//   * When a query finishes, the connection is returned to the idle pool
	//   * If idle pool is at capacity (5), the oldest idle connection is closed
	//   * Idle connections are ready to use immediately (no connection handshake overhead)
	// - Impact on connection accumulation:
	//   * Prevents keeping too many idle connections open unnecessarily
	//   * If all 25 connections are created but only 5 are needed for normal load,
	//     the extra 20 will be closed after they become idle
	//   * This allows the pool to shrink during low-traffic periods
	// - Tuning considerations:
	//   * Too low: frequent connection creation/destruction overhead during traffic spikes
	//   * Too high: wastes DB server resources with idle connections
	//   * Should be >= expected concurrent query load during normal operation
	//   * Set to 5 to handle moderate concurrent load while being resource-efficient
	db.SetMaxIdleConns(5)

	// SetConnMaxLifetime(5 * time.Minute):
	// - Maximum amount of time a connection can be reused before it's closed and recreated
	// - Runtime behavior:
	//   * Connection age is tracked from when it was first created
	//   * After 5 minutes, even if the connection is idle and healthy, it will be closed
	//   * The next query needing a connection will create a fresh one
	//   * This happens lazily - connections aren't proactively closed at exactly 5 minutes
	// - Impact on connection accumulation:
	//   * Prevents indefinite connection accumulation from long-lived processes
	//   * Ensures connections are periodically recycled, avoiding:
	//     - Stale connections that the DB server might have closed
	//     - Connections affected by network issues or DB server restarts
	//     - Connections holding onto resources (memory, prepared statements) indefinitely
	// - Tuning considerations:
	//   * Too low: excessive connection churn, overhead from frequent reconnections
	//   * Too high: stale connections, potential issues with DB server connection resets
	//   * Should be less than DB server's connection timeout settings
	//   * 5 minutes is a good balance - short enough to prevent staleness, long enough to avoid churn
	// - Special consideration for RDS:
	//   * AWS RDS may forcibly close idle connections after a timeout (default 8 hours)
	//   * This setting ensures our connections are refreshed well before RDS times them out
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("connected to SQL DB...")

	return db, dbDSNuriStr, nil
}

//-------------------------------------------------

func DBsqlViewTableStructure(pTableNameStr string,
	pRuntimeSys *RuntimeSys) *GFerror {

	rows, err := pRuntimeSys.SQLdb.Query(`
		SELECT column_name, data_type, udt_name
		FROM information_schema.columns
		WHERE table_name = $1`, pTableNameStr)

	if err != nil {
		gfErr := ErrorCreate("failed to run table structure query against the DB...",
			"sql_query_execute",
			map[string]interface{}{
				"table_name_str": pTableNameStr,
			},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}
	defer rows.Close()

	fmt.Printf("Table structure for '%s':\n", pTableNameStr)
	for rows.Next() {

		var columnName, dataType, udtName string
		if err := rows.Scan(&columnName, &dataType, &udtName); err != nil {

			gfErr := ErrorCreate("failed to run table structure query against the DB...",
				"sql_row_scan",
				map[string]interface{}{
					"table_name_str": pTableNameStr,
				},
				err, "gf_core", pRuntimeSys)
			return gfErr
		}
		if dataType == "ARRAY" {
			fmt.Printf("Column: %s, Type: ARRAY[%s]\n", columnName, udtName)
		} else {
			fmt.Printf("Column: %s, Type: %s\n", columnName, dataType)
		}
	}

	if err := rows.Err(); err != nil {
		gfErr := ErrorCreate("failed to run table structure query against the DB...",
			"sql_query_execute",
			map[string]interface{}{
				"table_name_str": pTableNameStr,
			},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}
	return nil
}
