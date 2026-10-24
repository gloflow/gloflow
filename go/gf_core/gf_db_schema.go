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
	"database/sql"
)

//-------------------------------------------------
// gf_errors table - mirrors the fields of GFerror.

func DBsqlCreateErrorsTable(pDB *sql.DB) error {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_errors (
		id                   BIGSERIAL        PRIMARY KEY,
		id_str               TEXT             NOT NULL,
		t                    TEXT             NOT NULL DEFAULT 'gf_error',
		creation_unix_time_f DOUBLE PRECISION NOT NULL,
		type_str             TEXT             NOT NULL,
		user_msg_str         TEXT             NOT NULL DEFAULT '',
		data_map             JSONB,
		descr_str            TEXT             NOT NULL DEFAULT '',
		error                TEXT,
		service_name_str     TEXT             NOT NULL DEFAULT '',
		subsystem_name_str   TEXT             NOT NULL DEFAULT '',
		stack_trace_str      TEXT             NOT NULL DEFAULT '',
		func_name_str        TEXT             NOT NULL DEFAULT '',
		file_str             TEXT             NOT NULL DEFAULT '',
		line_num_int         INTEGER          NOT NULL DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS gf_errors_creation_time_idx ON gf_errors (creation_unix_time_f DESC);
	CREATE INDEX IF NOT EXISTS gf_errors_type_idx          ON gf_errors (type_str);
	CREATE INDEX IF NOT EXISTS gf_errors_subsystem_idx     ON gf_errors (subsystem_name_str);
	`

	_, err := pDB.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("failed to create gf_errors table: %w", err)
	}
	return nil
}
