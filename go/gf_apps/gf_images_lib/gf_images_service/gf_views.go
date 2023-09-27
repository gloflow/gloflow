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

package gf_images_service

import (
	"context"
	"bytes"
	"text/template"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------

func renderImageViewPage(pImageID gf_images_core.GFimageID,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pUserID               gf_core.GF_ID,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {


	image, gfErr := gf_images_core.DBmongoGetImage(pImageID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}
	
	// META
	// plugin
	var filteredMetaJSONstr string
	if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback != nil {
		filteredMetaMap := pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback(image.MetaMap)
		metaJSONbytesLst, _ := json.Marshal(filteredMetaMap)
		filteredMetaJSONstr = string(metaJSONbytesLst)
	}

	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)


	type tmplData struct {
		ImageID               gf_images_core.GFimageID
		CreationUNIXtimeF     float64
		OwnerUserID           gf_core.GF_ID
		FlowsNamesLst         []string
		OriginPageURLstr      string
		ThumbnailMediumURLstr string
		TagsLst               []string
		FilteredMetaJSONstr   string
		SysReleaseInfo        gf_core.SysReleaseInfo
		IsSubtmplDef          func(string) bool //used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTemplate.Execute(buff, tmplData{
		ImageID:               pImageID,
		CreationUNIXtimeF:     image.Creation_unix_time_f,
		OwnerUserID:           image.UserID,
		FlowsNamesLst:         image.FlowsNamesLst,
		OriginPageURLstr:      image.Origin_page_url_str,
		ThumbnailMediumURLstr: image.Thumbnail_medium_url_str,
		TagsLst:               image.TagsLst,
		FilteredMetaJSONstr:   filteredMetaJSONstr,
		SysReleaseInfo:        sysReleaseInfo,

		//-------------------------------------------------
		// IS_SUBTEMPLATE_DEFINED
		IsSubtmplDef: func(pSubtemplateNameStr string) bool {
			for _, n := range pSubtemplatesNamesLst {
				if n == pSubtemplateNameStr {
					return true
				}
			}
			return false
		},

		//-------------------------------------------------
	})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to render the image view template",
			"template_render_error",
			map[string]interface{}{
				"image_id_str": pImageID,
				"user_id_str":  pUserID,
			},
			err, "gf_images_lib", pRuntimeSys)
		return "", gfErr
	}



	templateRenderedStr := buff.String()

	return templateRenderedStr, nil
}