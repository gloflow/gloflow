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
	"context"
	"time"
	"encoding/json"
)

//-------------------------------------------------
// persists a GFerror into the gf_errors SQL table.
// IMPORTANT!! - does not use ErrorCreate() on failure, to avoid infinite recursion.

func DBsqlPutError(pGFerr *GFerror,
	pRuntimeSys *RuntimeSys) error {

	if pRuntimeSys.SQLdb == nil {
		return fmt.Errorf("SQL DB not initialized, cant persist gf_error")
	}

	// DATA_MAP
	dataMapJSONstr := "null"
	if pGFerr.Data_map != nil {
		dataBytesLst, err := json.Marshal(pGFerr.Data_map)
		if err != nil {
			// values in the map may not be JSON-serializable
			dataBytesLst, _ = json.Marshal(map[string]interface{}{
				"_marshal_error": err.Error(),
				"_data_fmt":      fmt.Sprintf("%v", pGFerr.Data_map),
			})
		}
		dataMapJSONstr = string(dataBytesLst)
	}

	// ERROR - the underlying golang error, if any
	var errorStr *string
	if pGFerr.Error != nil {
		s := pGFerr.Error.Error()
		errorStr = &s
	}

	ctx, cancelFun := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFun()

	_, err := pRuntimeSys.SQLdb.ExecContext(ctx, `
		INSERT INTO gf_errors (
			id_str,
			t,
			creation_unix_time_f,
			type_str,
			user_msg_str,
			data_map,
			descr_str,
			error,
			service_name_str,
			subsystem_name_str,
			stack_trace_str,
			func_name_str,
			file_str,
			line_num_int
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9, $10, $11, $12, $13, $14)`,
		pGFerr.Id_str,
		pGFerr.T_str,
		pGFerr.Creation_unix_time_f,
		pGFerr.Type_str,
		pGFerr.User_msg_str,
		dataMapJSONstr,
		pGFerr.Descr_str,
		errorStr,
		pGFerr.ServiceNameStr,
		pGFerr.Subsystem_name_str,
		pGFerr.Stack_trace_str,
		pGFerr.Function_name_str,
		pGFerr.File_str,
		pGFerr.Line_num_int)
	if err != nil {
		return err
	}
	return nil
}
