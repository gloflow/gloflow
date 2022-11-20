/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_images_plugins

import (
	"fmt"
	"time"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func RunPyImagePlugins(pImageLocalFilePathStr string,
	pPluginsPyDirPathStr string,
	pMetrics             *GFmetrics,
	pCtx                 context.Context,
	pRuntimeSys          *gf_core.RuntimeSys) {

	

	// Py plugins are lost running, so dont block this function and its callers,
	// instead run the Py command in a new goroutine.
	//
	// ADD!! - dont run the Py plugin from scratch each time, instead start it once 
	//         (so that it can import its dependencies once) and then pass in
	//         requests for new image processing via STDIN
	go func() {

		pyPathStr := fmt.Sprintf("%s/gf_images_plugins_main.py", pPluginsPyDirPathStr)
		argsLst := []string{
			fmt.Sprintf("-image_local_file_path=%s", pImageLocalFilePathStr),
		}
		stdoutPrefixStr := "GF_OUT:"
		inputStdinStr   := ""



		runStartUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0


		// PY_RUN
		outputsLst, gfErr := gf_core.CLIpyRun(pyPathStr,
			argsLst,
			&inputStdinStr,
			stdoutPrefixStr,
			pRuntimeSys)

		

		if gfErr != nil {
			return
		}

		runEndUNIXtimeF   := float64(time.Now().UnixNano())/1000000000.0
		runDurrationSecsF := runEndUNIXtimeF - runStartUNIXtimeF
		
		if pMetrics != nil {
			pMetrics.PyPluginsExecDurationGauge.Set(runDurrationSecsF)
		}

		fmt.Println(outputsLst)

	}()
}

//-------------------------------------------------
/*func py__run_plugin__color_palette(p_input_images_local_file_paths_lst []string,
	p_output_dir_path_str string,
	p_plugins_info        *GF_py_plugins,
	p_runtime_sys         *gf_core.RuntimeSys) *gf_core.GFerror {



	median_cut_levels_num_int := 4

	py_path_str       := fmt.Sprintf("%s/gf_color_palette.py", p_plugins_info.Base_dir_path_str)
	stdout_prefix_str := "GF_OUT:"
	args_lst := []string{
		fmt.Sprintf("-input_images_local_file_paths=%s", strings.Join(p_input_images_local_file_paths_lst, ",")),
		fmt.Sprintf("-output_dir_path=%s", p_output_dir_path_str),
		fmt.Sprintf("-median_cut_levels_num=%d", median_cut_levels_num_int),
	}

	// PY_RUN
	outputs_lst, gf_err := gf_core.CLI_py__run(py_path_str,
		args_lst,
		nil, // input_stdin_str,
		stdout_prefix_str,
		p_runtime_sys)

	if gf_err != nil {
		return gf_err
	}



	fmt.Println(outputs_lst)

	return nil
}*/