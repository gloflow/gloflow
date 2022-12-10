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

package gf_admin_lib

import (
	"os"
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	// "github.com/davecgh/go-spew/spew"
)

var logFun func(string,string)
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, _  = gf_core.LogsInit()
	cliArgsMap = gf_images_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestTemplates(pTest *testing.T) {

	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_admin_test",
		LogFun:           logFun,
	}

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_admin_login":     "./../../../web/src/gf_apps/gf_admin/templates/gf_admin_login/gf_admin_login.html",
		"gf_admin_dashboard": "./../../../web/src/gf_apps/gf_admin/templates/gf_admin_dashboard/gf_admin_dashboard.html",
	}
	
	templates, gfErr := templatesLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	templateRenderedStr, gfErr := viewRenderTemplateDashboard(templates.dashboardTmpl,
		templates.dashboardSubtemplatesNamesLst,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(templateRenderedStr)
}