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

package main

import (
	"bytes"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_p2p"
)

//-------------------------------------------------
type gfTemplates struct {
	p2pStatus *template.Template
}

//-------------------------------------------------
func viewRenderStatus(pP2pStatus gf_p2p.GFp2pStatus,
	pTmpl       *template.Template,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)

	type tmplData struct {
		SysReleaseInfo gf_core.SysReleaseInfo
		P2pStatus      gf_p2p.GFp2pStatus
	}

	buff := new(bytes.Buffer)
	err := pTmpl.Execute(buff, tmplData{
		SysReleaseInfo: sysReleaseInfo,
		P2pStatus:      pP2pStatus,
	})
	if err != nil {

	}
	templateRenderedStr := buff.String()
	return templateRenderedStr, nil
}

//-------------------------------------------------
func tmplLoad(pTemplatesPathsMap map[string]string, // p_templates_dir_path_str string, 
	pRuntimeSys *gf_core.RuntimeSys) (*gfTemplates, *gf_core.GFerror) {

	
	statusFilepathStr := pTemplatesPathsMap["gf_p2p_status"]

	p2pStatusTmpl, _, gfErr := gf_core.TemplatesLoad(statusFilepathStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	templates := &gfTemplates{
		p2pStatus: p2pStatusTmpl,
	}
	return templates, nil
}