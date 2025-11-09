/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

package main

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
)

//-------------------------------------------------

func debugSerializeStateHistoryToFile(pFilePathStr string,
	pProgramsDebugLst []*gf_lang.GFprogramDebug,
	pRuntimeSys       *gf_core.RuntimeSys) *gf_core.GFerror {
	
	stateHistoriesLst := [][]*gf_lang.GFstate{}
	for _, d := range pProgramsDebugLst {
		stateHistoriesLst = append(stateHistoriesLst, d.StateHistoryLst)
	}


	// JSON
	jsonBytesLst := gf_core.EncodeJSONfromData(stateHistoriesLst)

	// FILE
	gfErr := gf_core.FileCreateWithContent(string(jsonBytesLst),
		pFilePathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil
	}
	return nil
}

//-------------------------------------------------

func debugSerializeOutputToFile(pFilePathStr string,
	pProgramsDebugLst []*gf_lang.GFprogramDebug,
	pRuntimeSys       *gf_core.RuntimeSys) *gf_core.GFerror {
	
	outputsLst := []interface{}{}
	for _, d := range pProgramsDebugLst {
		outputsLst = append(outputsLst, d.OutputLst)
	}

	// JSON
	jsonBytesLst := gf_core.EncodeJSONfromData(outputsLst)

	// FILE
	gfErr := gf_core.FileCreateWithContent(string(jsonBytesLst),
		pFilePathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil
	}
	return nil
}

//-------------------------------------------------

func debugRunPyAnalyzer(pProgramsDebugLst []*gf_lang.GFprogramDebug,
	pPlugins    *GFplugins,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	pyPathStr                         := fmt.Sprintf("%s/gf_debug_analyzer.py", pPlugins.PyBaseDirStr)
	serializedOutputFilePathStr       := "serialized_output.json"
	serializedStateHistoryFilePathStr := "state_history.json"

	stdoutPrefixStr := "GF_OUT:"
	
	argsLst := []string{
		fmt.Sprintf("-serialized_output_file=%s", serializedOutputFilePathStr),
		fmt.Sprintf("-state_history_file=%s", serializedStateHistoryFilePathStr),
	}

	envMap := map[string]string{}

	// PY_RUN
	_, gfErr := gf_core.CLIpyRun(pyPathStr,
		argsLst,
		nil,
		envMap,
		stdoutPrefixStr,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}

	// fmt.Println(outputsLst)

	return nil
}