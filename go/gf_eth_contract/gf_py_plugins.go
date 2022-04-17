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

package gf_eth_contract

import (
	"fmt"
	// "encoding/json"
	// log "github.com/sirupsen/logrus"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Py__run_plugin__get_contract_info(p_new_contract_addr_str string,
	p_plugins_info *gf_eth_core.GF_py_plugins,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.GF_error {


	py_path_str := fmt.Sprintf("%s/gf_plugin__get_contract_info.py", p_plugins_info.Base_dir_path_str)
	args_lst := []string{
		fmt.Sprintf("-contract_addr=%s", p_new_contract_addr_str),
	}
	stdout_prefix_str := "GF_OUT:"

	// PY_RUN
	outputs_lst, gf_err := gf_core.CLI_py__run(py_path_str,
		args_lst,
		nil,
		stdout_prefix_str,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}


	for _, o := range outputs_lst {

		fmt.Println(o)

	}

	return nil
}