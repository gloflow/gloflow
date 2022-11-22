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

package gf_eth_tx

import (
	"fmt"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------
func py__run_plugin__plot_tx_trace(p_tx_id_str string,
	p_tx_trace     *GF_eth__tx_trace,
	p_plugins_info *gf_eth_core.GF_py_plugins,
	p_runtime_sys  *gf_core.RuntimeSys) (string, *gf_core.GFerror) {


	py_path_str := fmt.Sprintf("%s/gf_plugin__plot_tx_trace.py", p_plugins_info.Base_dir_path_str)
	args_lst := []string{
		fmt.Sprintf("-tx_id=%s", p_tx_id_str),

		// write SVG output to stdout instead of to a file
		"-stdout",
	}
	stdout_prefix_str := "GF_OUT:"





	// JSON
	tx_trace_byte_lst, _ := json.Marshal(p_tx_trace)
	tx_trace_byte_str    := string(tx_trace_byte_lst)

	// PY_RUN
	outputs_lst, gfErr := gf_core.CLIpyRun(py_path_str,
		args_lst,
		&tx_trace_byte_str,
		stdout_prefix_str,
		p_runtime_sys)
	if gfErr != nil {
		return "", gfErr
	}

	svg_str := outputs_lst[0]["svg_str"].(string)

	// LOG
	log.WithFields(log.Fields{"tx_id": p_tx_id_str, "py_path": py_path_str}).Info("py_plugin gf_plugin__plot_tx_trace.py complete...")

	return svg_str, nil
}