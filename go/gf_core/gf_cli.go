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
func CLIrunStandard(pCmdLst []string,
	p_env_vars_map map[string]string,
	pRuntimeSys    *RuntimeSys) ([]string, []string, *GFerror) {

	cliInfo := &GF_CLI_cmd_info{
		Cmd_lst:          pCmdLst,
		Env_vars_map:     p_env_vars_map,
		Dir_str:          "",
		View_output_bool: true,
	}
	stdoutLst, stderrLst, gfErr := CLIrun(cliInfo, pRuntimeSys)
	return stdoutLst, stderrLst, gfErr
}

//-------------------------------------------------
// RUN
func CLIrun(pCmdInfo *GF_CLI_cmd_info,
	pRuntimeSys *RuntimeSys) ([]string, []string, *GFerror) {

	stdoutCh, stderrCh, gfErr := CLIrunCore(pCmdInfo,
		true,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	stdoutLst, stderrLst := consumeOutputs(stdoutCh, stderrCh)

	return stdoutLst, stderrLst, nil
}

//-------------------------------------------------
func CLIrunCore(pCmdInfo *GF_CLI_cmd_info,
	pWaitForCompletionBool bool,
	pRuntimeSys            *RuntimeSys) (chan string, chan string, *GFerror) {

	cmd_str := strings.Join(pCmdInfo.Cmd_lst, " ")
	fmt.Printf("%s\n", cmd_str)



	cmd_name_str := pCmdInfo.Cmd_lst[0]
	cmd_args_lst := pCmdInfo.Cmd_lst[1:]
	p := exec.Command(cmd_name_str, cmd_args_lst...)




	// CMD_DIR
	if pCmdInfo.Dir_str != "" {
		p.Dir = pCmdInfo.Dir_str
	}

	// STDIN
	stdin, err := p.StdinPipe()
	if err != nil {
		gfErr := Error__create("failed to get STDIN pipe of a CLI process",
			"cli_run_error",
			map[string]interface{}{"cmd": cmd_str,},
			err, "gf_core", pRuntimeSys)
		return nil, nil, gfErr 
	}
	// defer stdin.Close()

	// STDOUT
	cmd_stdout__reader, _ := p.StdoutPipe()
	cmd_stdout__buffer    := bufio.NewReader(cmd_stdout__reader)
	
	// STDERR
	cmd_stderr__reader, _ := p.StderrPipe()
	cmd_stderr__buffer    := bufio.NewReader(cmd_stderr__reader)

	// done_ch := make(chan bool)
	
	//----------------------
	// START
	err = p.Start()
	if err != nil {
		gfErr := Error__create("failed to Start a CLI command",
			"cli_run_error",
			map[string]interface{}{"cmd": cmd_str,},
			err, "gf_core", pRuntimeSys)
		return nil, nil, gfErr	
	}

	//----------------------
	// STDIN - input is written to stdin after the process is started
	if pCmdInfo.Stdin_data_str != nil {
		io.WriteString(stdin, fmt.Sprintf("%s\n", *pCmdInfo.Stdin_data_str))
	}

	//----------------------
	// STDOUT

	stdout_ch := make(chan string, 100)
	go func(p_stdout_ch chan string) {
		for {
			l, err := cmd_stdout__buffer.ReadString('\n')

			if fmt.Sprint(err) == "EOF" {
				p_stdout_ch <- "EOF"
				// done_ch <- true
				return
			}
			if err != nil {
				continue
			}
			if pCmdInfo.View_output_bool {
				fmt.Printf("%s", l)
			}
			l_str := strings.TrimSuffix(l, "\n")

			p_stdout_ch <- l_str
		}
	}(stdout_ch)

	//----------------------
	// STDERR

	stderr_ch := make(chan string, 100)
	go func() {
		for {
			l, err := cmd_stderr__buffer.ReadString('\n')
			// if fmt.Sprint(err) == "EOF" {
			// 	done_ch <- true
			// 	return
			// }
			if err != nil {
				continue
			}
			if pCmdInfo.View_output_bool {
				fmt.Printf("%s\n", l)
			}
			l_str := strings.TrimSuffix(l, "\n")
			stderr_ch <- l_str
		}
	}()

	//----------------------
	// WAIT

	if pWaitForCompletionBool {

		err = p.Wait()
		if err != nil {
			stdoutLst, stderrLst := consumeOutputs(stdout_ch, stderr_ch)
			gfErr := Error__create("failed to Wait for a CLI command to complete",
				"cli_run_error",
				map[string]interface{}{
					"cmd":        cmd_str,
					"stdout_lst": stdoutLst,
					"stderr_lst": stderrLst,
				},
				err, "gf_core", pRuntimeSys)
			return nil, nil, gfErr
		}
	}

	//----------------------

	return stdout_ch, stderr_ch, nil // done_ch, nil
}

//-------------------------------------------------
func consumeOutputs(pStdoutCh chan string, pStderrCh chan string) ([]string, []string) {
	stdoutLst := []string{}
	stderrLst := []string{}
	for {
		select {
		case stdoutLineStr := <- pStdoutCh:

			if stdoutLineStr == "EOF" {
				return stdoutLst, stderrLst
			} else {
				stdoutLst = append(stdoutLst, stdoutLineStr)
			}
			
		case stderr_l_str := <- pStderrCh:
			stderrLst = append(stdoutLst, stderr_l_str)

		// case _ = <- done_ch:
		// 	return stdoutLst, stderrLst
		
		}
	}
	return nil, nil
}

//-------------------------------------------------
func CLIprompt() {





	
}