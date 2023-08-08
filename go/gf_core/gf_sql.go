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
	"database/sql"
	"fmt"

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
	
	pRuntimeSys.LogNewFun("INFO", "connecting to SQL DB...",
		map[string]interface{}{
			"db_host": pDBhostStr,
		})

	// connect
	db, err := sql.Open("postgres", urlStr)
	if err != nil {
		gfErr := ErrorCreate("failed to connect to a SQL server at target url",
			"sql_failed_to_connect",
			map[string]interface{}{
				"db_host_str": pDBhostStr,
				"db_name_str": pDBnameStr,
			}, err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	// test the connection
	err = db.Ping()
	if err != nil {
		gfErr := ErrorCreate("failed to connect to a SQL server at target url",
			"sql_failed_to_connect",
			map[string]interface{}{
				"db_host_str": pDBhostStr,
				"db_name_str": pDBnameStr,
			}, err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	fmt.Println("connected to SQL DB...")

	return db, nil
}