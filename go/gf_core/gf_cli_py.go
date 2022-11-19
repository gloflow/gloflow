/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_core

import (
	// "fmt"
	"strings"
	"path/filepath"
	"encoding/json"
)

//-------------------------------------------------

func CLIpyRun(pPyPathStr string,
	pArgsLst         []string,
	pInputStdinStr   *string,
	pStdoutPrefixStr string,
	pRuntimeSys      *RuntimeSys) ([]map[string]interface{}, *GFerror) {
	

	pyAbsPathStr, _ := filepath.Abs(pPyPathStr)

	cmdLst := []string{"python3", "-u", pyAbsPathStr,}
	cmdLst = append(cmdLst, pArgsLst...)

	cmdInfo := GF_CLI_cmd_info {
		Cmd_lst:          cmdLst,
		Stdin_data_str:   pInputStdinStr,
		Env_vars_map:     map[string]string{},
		Dir_str:          "",
		View_output_bool: true,
	}


	// RUN
	stdoutLst, _, gfErr := CLIrun(&cmdInfo, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}


	// PARSE
	parsedOutputLst, gfErr := cliPyParseOutput(stdoutLst, pStdoutPrefixStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return parsedOutputLst, nil
}

//-------------------------------------------------

func cliPyParseOutput(pStdoutLst []string,
	pStdoutPrefixStr string,
	pRuntimeSys      *RuntimeSys) ([]map[string]interface{}, *GFerror) {

	outputsLst := []map[string]interface{}{}
	for _, l_str := range pStdoutLst {
		
		if strings.HasPrefix(l_str, pStdoutPrefixStr) {

			// remove the stdout prefix of the Py program stdout
			outputStr := strings.Replace(l_str, pStdoutPrefixStr, "", 1)

			// JSON_DECODE
			var outputMap map[string]interface{}
			err := json.Unmarshal([]byte(outputStr), &outputMap)

			if err != nil {
				gfErr := ErrorCreate("failed to parse json output in py program stdout",
					"json_decode_error",
					map[string]interface{}{"stdout_line_str": l_str,},
					err, "gf_core", pRuntimeSys)
				return nil, gfErr
			}

			outputsLst = append(outputsLst, outputMap)
		}
	}

	return outputsLst, nil
}