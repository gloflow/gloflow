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
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func Py__run_plugin__get_contract_info(pNewContractAddrStr string,
	pPluginsInfo *gf_eth_core.GF_py_plugins,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {


	py_path_str := fmt.Sprintf("%s/gf_plugin__get_contract_info.py", pPluginsInfo.Base_dir_path_str)
	args_lst := []string{
		fmt.Sprintf("-contract_addr=%s", pNewContractAddrStr),
	}
	stdout_prefix_str := "GF_OUT:"

	// PY_RUN
	outputs_lst, gfErr := gf_core.CLIpyRun(py_path_str,
		args_lst,
		nil,
		stdout_prefix_str,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	for _, o := range outputs_lst {

		fmt.Println(o)

	}

	return nil
}