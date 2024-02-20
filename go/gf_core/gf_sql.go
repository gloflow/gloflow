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
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

//-------------------------------------------------

func DBsqlConnect(pDBnameStr string,
	pUserNameStr string,
	pPassStr     string,
	pDBhostStr   string,
	pRuntimeSys  *RuntimeSys) (*sql.DB, *GFerror) {
	
	// FIX!! - make "sslmode=disable" configurable, dont hardcode it
	urlStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		pUserNameStr,
		pPassStr,
		pDBhostStr,
		pDBnameStr)

	//-----------------------
	// CONNECT

	maxRetriesInt := 5
	retryIntervalSecsInt := 2

	var db *sql.DB
	var err error

	for retriesInt := 0; retriesInt < maxRetriesInt; retriesInt++ {

		pRuntimeSys.LogNewFun("INFO", "attempt - connecting to SQL DB...", nil)
		db, err = sql.Open("postgres", urlStr)
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
		return nil, gfErr
	}

	//-----------------------

	fmt.Println("connected to SQL DB...")

	return db, nil
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
