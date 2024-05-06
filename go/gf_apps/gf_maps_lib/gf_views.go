/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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
	// "fmt"
	"context"
	"strconv"
	"bytes"
	"strings"
	"text/template"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------

func renderInitialPage(pFlowNameStr string,
	pInitialPagesNumInt   int, // 6
	pPageSizeInt          int, // 5
	pTmpl                 *template.Template,
	pSubtemplatesNamesLst []string,
	pUserID               gf_core.GF_ID,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// RENDER
	templateRenderedStr, gfErr := renderTemplate(pTmpl,
		pSubtemplatesNamesLst,
		pUserID,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return templateRenderedStr, nil
}

//-------------------------------------------------

func renderTemplate(
	pTemplate                *template.Template,
	pSubtemplatesNamesLst    []string,
	pUserID                  gf_core.GF_ID,
	pRuntimeSys              *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// plugin
	metadataFilterDefinedBool := false
	if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback != nil {
		metadataFilterDefinedBool = true
	}

	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)

	type tmplData struct {
		Sys_release_info gf_core.SysReleaseInfo
		Is_subtmpl_def   func(string) bool //used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTemplate.Execute(buff, tmplData{
		Sys_release_info: sysReleaseInfo,

		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
		Is_subtmpl_def: func(p_subtemplate_name_str string) bool {
			for _, n := range pSubtemplatesNamesLst {
				if n == p_subtemplate_name_str {
					return true
				}
			}
			return false
		},

		//-------------------------------------------------
	})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to render the images flow template",
			"template_render_error",
			map[string]interface{}{
				"flow_name_str": pFlowNameStr,
				"user_id_str":   pUserID,
			},
			err, "gf_images_lib", pRuntimeSys)
		return "", gfErr
	}

	templateRenderedStr := buff.String()
	return templateRenderedStr, nil
}