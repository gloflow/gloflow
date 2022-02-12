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

var log_fun func(string,string)
var cli_args_map map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	log_fun      = gf_core.Init_log_fun()
	cli_args_map = gf_images_core.CLI__parse_args(log_fun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func Test__templates(p_test *testing.T) {

	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_admin_test",
		Log_fun:          log_fun,
	}

	// TEMPLATES
	templates_paths_map := map[string]string{
		"gf_admin_dashboard": "./../../../web/src/gf_apps/gf_admin/templates/gf_admin_dashboard/gf_admin_dashboard.html",
	}
	
	gf_templates, gf_err := tmpl__load(templates_paths_map, runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	template_rendered_str, gf_err := view__render_template_dashboard(gf_templates.dashboard__tmpl,
		gf_templates.dashboard__subtemplates_names_lst,
		runtime_sys)
	if gf_err != nil {
		p_test.Fail()
	}

	fmt.Println(template_rendered_str)
}