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

	// GET_DATA
	pagesLst, pagesUserNamesLst, flowPagesNumInt, gfErr := getTemplateData(pFlowNameStr,
		pInitialPagesNumInt,
		pPageSizeInt,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	// RENDER
	templateRenderedStr, gfErr := renderTemplate(pFlowNameStr,
		pagesLst,
		pagesUserNamesLst,
		flowPagesNumInt,
		pTmpl,
		pSubtemplatesNamesLst,
		pUserID,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return templateRenderedStr, nil
}

//-------------------------------------------------

func getTemplateData(pFlowNameStr string,
	pInitialPagesNumInt int,
	pPageSizeInt        int,
	pCtx                context.Context,
	pRuntimeSys         *gf_core.RuntimeSys) ([][]*gf_images_core.GFimage, [][]gf_identity_core.GFuserName, int64, *gf_core.GFerror) {

	//---------------------
	// GET_IMAGES

	imagesPagesLst := [][]*gf_images_core.GFimage{}

	for i:=0; i < pInitialPagesNumInt; i++ {

		startPositionInt := i*pPageSizeInt

		pRuntimeSys.LogNewFun("DEBUG", ">>>>>>> get image page from DB",
			map[string]interface{}{
				"start_position_int": startPositionInt,
				"page_size_int":      pPageSizeInt,
			})

		//------------
		// DB GET PAGE
		
		// initial page might be larger then subsequent pages, that are requested 
		// dynamically by the front-end
		pageLst, gfErr := dbGetPage(pFlowNameStr,
			startPositionInt, // p_cursor_start_position_int
			pPageSizeInt,     // p_elements_num_int
			pCtx,
			pRuntimeSys)

		if gfErr != nil {
			return nil, nil, 0, gfErr
		}
		
		/*
		// initial page might be larger then subsequent pages, that are requested 
		// dynamically by the front-end
		pageLst, gfErr := dbMongoGetPage(pFlowNameStr,
			startPositionInt, // p_cursor_start_position_int
			pPageSizeInt,     // p_elements_num_int
			pCtx,
			pRuntimeSys)

		if gfErr != nil {
			return nil, nil, 0, gfErr
		}
		*/

		//------------

		imagesPagesLst = append(imagesPagesLst, pageLst)
	}

	//---------------------
	// RESOLVE_USER_IDS_TO_USERNAMES

	pagesUserNamesLst := resolveUserIDsToUserNames(imagesPagesLst, pCtx, pRuntimeSys)

	//---------------------
	// TOTAL_PAGES_NUM
	flowPagesNumInt, gfErr := dbMongoGetPagesTotalNum(pFlowNameStr,
		pPageSizeInt,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, 0, gfErr
	}

	//---------------------
	return imagesPagesLst, pagesUserNamesLst, flowPagesNumInt, nil
}

//-------------------------------------------------

func renderTemplate(pFlowNameStr string,
	pImagesPagesLst          [][]*gf_images_core.GFimage,
	pImagesPagesUserNamesLst [][]gf_identity_core.GFuserName,
	pFlowPagesNumInt         int64,
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
	//-------------------------
	imagesPagesLst := [][]map[string]interface{}{}
	for i, imagesPageLst := range pImagesPagesLst {

		pageImagesLst := []map[string]interface{}{}
		for j, image := range imagesPageLst {

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
				"flows_names_str":           strings.Join(image.FlowsNamesLst, ","),
				"format_str":                image.Format_str,
				"thumbnail_small_url_str":   image.ThumbnailSmallURLstr,
				"thumbnail_medium_url_str":  image.ThumbnailMediumURLstr,
				"thumbnail_large_url_str":   image.ThumbnailLargeURLstr,
				"image_origin_page_url_str": image.Origin_page_url_str,
				"owner_user_name_str":       string(pImagesPagesUserNamesLst[i][j]),

				// "owner_user_id_str": image.UserID,
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