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

package gf_home_lib

import (
	"os"
	"fmt"
	"time"
	"testing"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func TestMain(m *testing.M) {

	logFun     = gf_core.Init_log_fun()
	cliArgsMap = CLIparseArgs(logFun)

	runtimeSys := Tinit()

	templatesPathsMap := map[string]string{
		"gf_home_main": "./../../../web/src/gf_apps/gf_home/templates/gf_home_main/gf_home_main.html",
	}

	testPortInt := 2000
	go func() {

		HTTPmux := http.NewServeMux()

		serviceInfo := &GFserviceInfo{}
		InitService(templatesPathsMap,
			serviceInfo,
			HTTPmux,
			runtimeSys)
		gf_rpc_lib.Server__init_with_mux(testPortInt, HTTPmux)
	}()
	time.Sleep(2*time.Second) // let server startup

	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func TestHomeViz(pTest *testing.T) {

	runtimeSys := Tinit()
	fmt.Println(runtimeSys)



}

//---------------------------------------------------
func TestTemplates(pTest *testing.T) {

	runtimeSys := &gf_core.Runtime_sys{
		Service_name_str: "gf_home_test",
		Log_fun:          logFun,
	}

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_home_main": "./../../../web/src/gf_apps/gf_home/templates/gf_home_main/gf_home_main.html",
	}
	
	templates, gfErr := templatesLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	templateRenderedStr, gfErr := viewRenderTemplateDashboard(templates.mainTmpl,
		templates.mainSubtemplatesNamesLst,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(templateRenderedStr)
}