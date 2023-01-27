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
	"os"
	"testing"
	log "github.com/sirupsen/logrus"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestMain(m *testing.M) {

	os.Setenv("GF_LOG_LEVEL", "debug")
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestExampleFirstScene(pTest *testing.T) {

	//---------------
	// INIT
	logFun, logNewFun := gf_core.LogsInitNew(true, "DEBUG")
	log.SetOutput(os.Stdout)

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_lang_server",
		EnvStr:         "dev",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}

	externAPI := gf_lang.GetTestExternAPI()

	//---------------

	// program
	localTestProgramStr := "./tests/first_scene.gf"
	programASTlst, gfErr := ParseProgramASTfromFile(localTestProgramStr, runtimeSys)
	if gfErr != nil {
		panic(1)
	}

	// run program
	resultsLst, programsDebugLst, err := gf_lang.Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program in test_basic", map[string]interface{}{"err": err,})
		pTest.Fail()
	}


	spew.Dump(resultsLst)
	// spew.Dump(programsDebugLst[0].StateHistoryLst)

	for _, s := range programsDebugLst[0].StateHistoryLst {
		fmt.Printf("++++ %f     time - %f\n", s.Xf, s.CreationUNIXtimeF)
	}

	debug       := programsDebugLst[0]
	filePathStr := "serialized_output.json"

	gfErr = debugSerializeOutputToFile(filePathStr,
		debug,
		runtimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}
}