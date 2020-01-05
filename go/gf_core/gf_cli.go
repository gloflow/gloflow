/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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
	"fmt"
	"strings"
	"os/exec"
	"bufio"
)

//-------------------------------------------------
type Gf_CLI_cmd_info struct {
	Cmd_lst          []string
	Env_vars_map     map[string]string
	Dir_str          string
	View_output_bool bool
}

//-------------------------------------------------
// RUN_STANDARD
func CLI__run_standard(p_cmd_lst []string,
	p_env_vars_map map[string]string,
	p_runtime      *Runtime_sys) ([]string, []string, *Gf_error) {

	cli_info := &Gf_CLI_cmd_info{
		Cmd_lst:          p_cmd_lst,
		Env_vars_map:     p_env_vars_map,
		Dir_str:          "",
		View_output_bool: true,
	}
	stdout_lst, stderr_lst, gf_err := CLI__run(cli_info, p_runtime)
	return stdout_lst, stderr_lst, gf_err
}

//-------------------------------------------------
// RUN
func CLI__run(p_cmd_info *Gf_CLI_cmd_info,
	p_runtime *Runtime_sys) ([]string, []string, *Gf_error) {



	cmd_str := strings.Join(p_cmd_info.Cmd_lst, " ")
	fmt.Printf("%s\n", cmd_str)



	cmd_name_str := p_cmd_info.Cmd_lst[0]
	cmd_args_lst := p_cmd_info.Cmd_lst[1:]
	cmd := exec.Command(cmd_name_str, cmd_args_lst...)




	// CMD_DIR
	if p_cmd_info.Dir_str != "" {
		cmd.Dir = p_cmd_info.Dir_str
	}


	// STDOUT
	cmd_stdout__reader, _ := cmd.StdoutPipe()
	cmd_stdout__buffer    := bufio.NewReader(cmd_stdout__reader)
	// STDERR
	cmd_stderr__reader, _ := cmd.StderrPipe()
	cmd_stderr__buffer    := bufio.NewReader(cmd_stderr__reader)







	err := cmd.Start()
	if err != nil {
		gf_err := Error__create("failed to Start a CLI command",
			"cli_run_error",
			map[string]interface{}{"cmd": cmd_str,},
			err, "gf_core", p_runtime)
		return nil, nil, gf_err	
	}



	stdout_lst := []string{}
	stderr_lst := []string{}

	// STDOUT
	go func() {
		for {
			l, err := cmd_stdout__buffer.ReadString('\n')
			if fmt.Sprint(err) == "EOF" {
				return
			}
			if err != nil {
				continue
			}
			if p_cmd_info.View_output_bool {
				fmt.Printf("%s\n", l)
			}

			stdout_lst = append(stdout_lst, strings.TrimSuffix(l, "\n"))
		}
	}()

	// STDERR
	go func() {
		for {
			l, err := cmd_stderr__buffer.ReadString('\n')
			if fmt.Sprint(err) == "EOF" {
				return
			}
			if err != nil {
				continue
			}
			if p_cmd_info.View_output_bool {
				fmt.Printf("%s\n", l)
			}
			
			stderr_lst = append(stderr_lst, strings.TrimSuffix(l, "\n"))
		}
	}()





	err = cmd.Wait()
	if err != nil {
		gf_err := Error__create("failed to Wait for a CLI command",
			"cli_run_error",
			map[string]interface{}{"cmd": cmd_str,},
			err, "gf_core", p_runtime)
		return nil, nil, gf_err
	}

	return stdout_lst, stderr_lst, nil
}

//-------------------------------------------------
func CLI__prompt() {





	
}