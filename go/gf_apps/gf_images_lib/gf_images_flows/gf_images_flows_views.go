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
	"fmt"
	"context"
	"strconv"
	"bytes"
	"text/template"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------

func renderInitialPage(pFlowNameStr string,
	p_initial_pages_num_int  int, // 6
	p_page_size_int          int, // 5
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	pCtx                     context.Context,
	pRuntimeSys              *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	//---------------------
	// GET_TEMPLATE_DATA

	pagesLst := [][]*gf_images_core.GFimage{}

	for i:=0; i < p_initial_pages_num_int; i++ {

		start_position_int := i*p_page_size_int
		// int end_position_int = start_position_int+p_page_size_int;

		pRuntimeSys.LogFun("INFO", fmt.Sprintf(">>>>>>> start_position_int - %d - %d", start_position_int, p_page_size_int))
		//------------
		// DB GET PAGE

		// initial page might be larger then subsequent pages, that are requested 
		// dynamically by the front-end
		page_lst, gfErr := dbGetPage(pFlowNameStr,
			start_position_int, // p_cursor_start_position_int
			p_page_size_int,    // p_elements_num_int
			pCtx,
			pRuntimeSys)

		if gfErr != nil {
			return "", gfErr
		}

		//------------

		pagesLst = append(pagesLst, page_lst)
	}


	flowPagesNumInt, gfErr := dbGetPagesTotalNum(pFlowNameStr,
		p_page_size_int,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	//---------------------
	templateRenderedStr, gfErr := renderTemplate(pFlowNameStr,
		pagesLst,
		flowPagesNumInt,
		p_tmpl,
		p_subtemplates_names_lst,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return templateRenderedStr, nil
}

//-------------------------------------------------

func renderTemplate(pFlowNameStr string,
	pImagesPagesLst       [][]*gf_images_core.GFimage,
	pFlowPagesNumInt      int64,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// plugin
	metadataFilterDefinedBool := false
	if pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback != nil {
		metadataFilterDefinedBool = true
	}

	sysReleaseInfo := gf_core.GetSysReleseInfo(pRuntimeSys)
	//-------------------------
	imagesPagesLst := [][]map[string]interface{}{}
	for _, imagesPageLst := range pImagesPagesLst {

		pageImagesLst := []map[string]interface{}{}
		for _, image := range imagesPageLst {

			// META
			var filteredMetaJSONstr string
			if metadataFilterDefinedBool {
				filteredMetaMap := pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback(image.MetaMap)
				metaJSONbytesLst, _ := json.Marshal(filteredMetaMap)
				filteredMetaJSONstr = string(metaJSONbytesLst)
			}

			imageInfoMap := map[string]interface{}{
				"creation_unix_time_str":    strconv.FormatFloat(image.Creation_unix_time_f, 'f', 6, 64),
				"id_str":                    image.IDstr,
				"title_str":                 image.TitleStr,
				"meta_json_str":             filteredMetaJSONstr,
				"format_str":                image.Format_str,
				"thumbnail_small_url_str":   image.Thumbnail_small_url_str,
				"thumbnail_medium_url_str":  image.Thumbnail_medium_url_str,
				"thumbnail_large_url_str":   image.Thumbnail_large_url_str,
				"image_origin_page_url_str": image.Origin_page_url_str,
			}

			if len(image.TagsLst) > 0 {
				imageInfoMap["image_has_tags_bool"] = true
				imageInfoMap["tags_lst"]            = image.TagsLst
			} else {
				imageInfoMap["image_has_tags_bool"] = false
			}

			pageImagesLst = append(pageImagesLst, imageInfoMap)
		}
		imagesPagesLst = append(imagesPagesLst, pageImagesLst)
	}

	//-------------------------

	type tmplData struct {
		Flow_name_str      string
		Images_pages_lst   [][]map[string]interface{}
		Flow_pages_num_int int64
		Sys_release_info   gf_core.SysReleaseInfo
		Is_subtmpl_def     func(string) bool //used inside the main_template to check if the subtemplate is defined
	}

	buff := new(bytes.Buffer)
	err := pTemplate.Execute(buff, tmplData{
		Flow_name_str:      pFlowNameStr,
		Images_pages_lst:   imagesPagesLst,
		Flow_pages_num_int: pFlowPagesNumInt,
		Sys_release_info:   sysReleaseInfo,

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
			map[string]interface{}{},
			err, "gf_images_lib", pRuntimeSys)
		return "", gfErr
	}

	templateRenderedStr := buff.String()
	return templateRenderedStr, nil
}