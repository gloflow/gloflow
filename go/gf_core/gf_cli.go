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
	"io"
)

//-------------------------------------------------
type GF_CLI_cmd_info struct {
	Cmd_lst          []string
	Stdin_data_str   *string // data to be passed via stdin
	Env_vars_map     map[string]string
	Dir_str          string
	View_output_bool bool
}

//-------------------------------------------------
// RUN_STANDARD
func CLI__run_standard(p_cmd_lst []string,
	p_env_vars_map map[string]string,
	p_runtime_sys  *Runtime_sys) ([]string, []string, *Gf_error) {

	cli_info := &GF_CLI_cmd_info{
		Cmd_lst:          p_cmd_lst,
		Env_vars_map:     p_env_vars_map,
		Dir_str:          "",
		View_output_bool: true,
	}
	stdout_lst, stderr_lst, gf_err := CLI__run(cli_info, p_runtime_sys)
	return stdout_lst, stderr_lst, gf_err
}

//-------------------------------------------------
// RUN
func CLI__run(p_cmd_info *GF_CLI_cmd_info,
	p_runtime_sys *Runtime_sys) ([]string, []string, *Gf_error) {

	stdout_ch, stderr_ch, done_ch, gf_err := CLI__run_core(p_cmd_info,
		true,
		p_runtime_sys)
	if gf_err != nil {
		return nil, nil, gf_err
	}

	//-------------------------------------------------
	consume_fn := func() ([]string, []string) {
		stdout_lst := []string{}
		stderr_lst := []string{}
		for {
			select {
			case stdout_l_str := <- stdout_ch:
				stdout_lst = append(stdout_lst, stdout_l_str)
			case stderr_l_str := <- stderr_ch:
				stderr_lst = append(stdout_lst, stderr_l_str)
			case _ = <- done_ch:
				return stdout_lst, stderr_lst
			}
		}
		return nil, nil
	}

	//-------------------------------------------------
	stdout_lst, stderr_lst := consume_fn()

	return stdout_lst, stderr_lst, nil
}

//-------------------------------------------------
func CLI__run_core(p_cmd_info *GF_CLI_cmd_info,
	p_wait_for_completion_bool bool,
	p_runtime_sys              *Runtime_sys) (chan string, chan string, chan bool, *Gf_error) {

	cmd_str := strings.Join(p_cmd_info.Cmd_lst, " ")
	fmt.Printf("%s\n", cmd_str)



	cmd_name_str := p_cmd_info.Cmd_lst[0]
	cmd_args_lst := p_cmd_info.Cmd_lst[1:]
	p := exec.Command(cmd_name_str, cmd_args_lst...)




	// CMD_DIR
	if p_cmd_info.Dir_str != "" {
		p.Dir = p_cmd_info.Dir_str
	}

	// STDIN
	stdin, err := p.StdinPipe()
	if err != nil {
		gf_err := Error__create("failed to get STDIN pipe of a CLI process",
			"cli_run_error",
			map[string]interface{}{"cmd": cmd_str,},
			err, "gf_core", p_runtime_sys)
		return nil, nil, nil, gf_err 
	}
	// defer stdin.Close()

	// STDOUT
	cmd_stdout__reader, _ := p.StdoutPipe()
	cmd_stdout__buffer    := bufio.NewReader(cmd_stdout__reader)
	
	// STDERR
	cmd_stderr__reader, _ := p.StderrPipe()
	cmd_stderr__buffer    := bufio.NewReader(cmd_stderr__reader)

	done_ch := make(chan bool)
	
	//----------------------
	// START
	err = p.Start()
	if err != nil {
		gf_err := Error__create("failed to Start a CLI command",
			"cli_run_error",
			map[string]interface{}{"cmd": cmd_str,},
			err, "gf_core", p_runtime_sys)
		return nil, nil, nil, gf_err	
	}

	//----------------------
	// STDIN - input is written to stdin after the process is started
	if p_cmd_info.Stdin_data_str != nil {
		io.WriteString(stdin, fmt.Sprintf("%s\n", *p_cmd_info.Stdin_data_str))
	}

	//----------------------
	// STDOUT

	stdout_ch := make(chan string, 100)
	go func() {
		for {
			l, err := cmd_stdout__buffer.ReadString('\n')
			if fmt.Sprint(err) == "EOF" {
				done_ch <- true
				return
			}
			if err != nil {
				continue
			}
			if p_cmd_info.View_output_bool {
				fmt.Printf("%s\n", l)
			}
			l_str := strings.TrimSuffix(l, "\n")
			stdout_ch <- l_str
			// stdout_lst = append(stdout_lst, )
		}
	}()

	//----------------------
	// STDERR

	stderr_ch := make(chan string, 100)
	go func() {
		for {
			l, err := cmd_stderr__buffer.ReadString('\n')
			if fmt.Sprint(err) == "EOF" {
				done_ch <- true
				return
			}
			if err != nil {
				continue
			}
			if p_cmd_info.View_output_bool {
				fmt.Printf("%s\n", l)
			}
			l_str := strings.TrimSuffix(l, "\n")
			stderr_ch <- l_str
			// stderr_lst = append(stderr_lst, strings.TrimSuffix(l, "\n"))
		}
	}()

	//----------------------
	// WAIT

	if p_wait_for_completion_bool {
		err = p.Wait()
		if err != nil {
			gf_err := Error__create("failed to Wait for a CLI command to complete",
				"cli_run_error",
				map[string]interface{}{"cmd": cmd_str,},
				err, "gf_core", p_runtime_sys)
			return nil, nil, nil, gf_err
		}
	}

	//----------------------

	return stdout_ch, stderr_ch, done_ch, nil
}

//-------------------------------------------------
func CLI__prompt() {





	
}