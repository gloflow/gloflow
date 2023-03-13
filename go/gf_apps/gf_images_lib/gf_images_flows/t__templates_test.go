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

package gf_images_flows

import (
	"os"
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	// "github.com/davecgh/go-spew/spew"
)

var logFun func(string,string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//---------------------------------------------------

func TestMain(m *testing.M) {
	logFun, logNewFun  = gf_core.LogsInit()
	cliArgsMap = gf_images_core.CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestTemplates(pTest *testing.T) {

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_images_flows_tests",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
	}

	// TEMPLATES
	templatesPathsMap := map[string]string{
		"gf_images_flows_browser": "./../../../../web/src/gf_apps/gf_images/templates/gf_images_flows_browser/gf_images_flows_browser.html",
	}
	
	gfTemplates, gfErr := tmplLoad(templatesPathsMap, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	imagesPagesLst := [][]*gf_images_core.GFimage{
		{
			&gf_images_core.GFimage{
				IDstr:      "some_test_id",
				TitleStr:   "some_test_img",
				MetaMap:    map[string]interface{}{"t_k": "val"},
				Format_str: "jpg",
				Thumbnail_small_url_str:  "url1",
				Thumbnail_medium_url_str: "url2",
				Thumbnail_large_url_str:  "url3",
				Origin_page_url_str:      "url4",
			},
		},
	}

	flowNameStr     := "test_flow" 
	flowPagesNumInt := int64(6)
	templateRenderedStr, gfErr := renderTemplate(flowNameStr,
		imagesPagesLst,
		flowPagesNumInt,
		gfTemplates.flows_browser__tmpl,
		gfTemplates.flows_browser__subtemplates_names_lst,
		runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	fmt.Println(templateRenderedStr)
}