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
	"strings"
	"path/filepath"
	"encoding/json"
)

//-------------------------------------------------
func CLI_py__run(p_py_path_str string,
	p_stdout_prefix_str string,
	p_runtime_sys       *Runtime_sys) ([]map[string]interface{}, *Gf_error) {
	


	
	py_abs_path_str, _ := filepath.Abs(p_py_path_str)

	cmd_info := GF_CLI_cmd_info {
		Cmd_lst:          []string{"python3", "-u", py_abs_path_str,},
		Env_vars_map:     map[string]string{},
		Dir_str:          "",
		View_output_bool: true,
	}



	// RUN
	stdout_lst, _, gf_err := CLI__run(&cmd_info, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	// PARSE
	parsed_output_lst, gf_err := cli_py__parse_output(stdout_lst, p_stdout_prefix_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}




	return parsed_output_lst, nil
}


//-------------------------------------------------
func cli_py__parse_output(p_stdout_lst []string,
	p_stdout_prefix_str string,
	p_runtime_sys       *Runtime_sys) ([]map[string]interface{}, *Gf_error) {



	output_lst := []map[string]interface{}{}
	for _, l_str := range p_stdout_lst {




		if strings.HasPrefix(l_str, p_stdout_prefix_str) {



			output_str := strings.Replace(l_str, p_stdout_prefix_str, "", 1)



			var o map[string]interface{}
			err := json.Unmarshal([]byte(output_str), &o)

			if err != nil {
				gf_err := Error__create("failed to parse json output in py program stdout",
					"json_decode_error",
					map[string]interface{}{"stdout_line_str": l_str,},
					err, "gf_core", p_runtime_sys)
				return nil, gf_err
			}
		}
	}

	return output_lst, nil
}