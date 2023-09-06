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

package gf_tagger_lib

import (
	"fmt"
	"strings"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_address"
)

//---------------------------------------------------
// pTagsStr           - "," separated list of strings
// pObjectExternIDstr - this is an external identifier for an object, not necessarily its internal. 
//                      for posts - their p_object_extern_id_str is their Title, but internally they have
//                      another ID.

func addTagsToObject(pTagsStr string,
	pObjectTypeStr     string,
	pObjectExternIDstr string,
	pMetaMap           map[string]interface{},
	pUserID            gf_core.GF_ID,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	if pObjectTypeStr != "post" &&
		pObjectTypeStr != "image" &&
		pObjectTypeStr != "event" &&
		pObjectTypeStr != "address" {

		gfErr := gf_core.ErrorCreate(fmt.Sprintf("object_type (%s) is not of supported type (post|image|event)",
			pObjectTypeStr),
			"verify__invalid_value_error",
			map[string]interface{}{
				"tags_str":        pTagsStr,
				"object_type_str": pObjectTypeStr,
			},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}
	
	tagsLst, gfErr := parseTags(pTagsStr,
		500, // pMaxTagsBulkSizeInt        int, // 500
		20,  // pMaxTagCharactersNumberInt int, // 20	
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	pRuntimeSys.LogFun("INFO", fmt.Sprintf("tags_lst - %s", tagsLst))

	//---------------
	// POST
	
	switch pObjectTypeStr {
		//---------------
		// POST
		case "post":
			postTitleStr      := pObjectExternIDstr
			existsBool, gfErr := gf_publisher_core.DBmongoCheckPostExists(postTitleStr,
				pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			
			if existsBool {
				pRuntimeSys.LogNewFun("DEBUG", "POST EXISTS", nil)

				gfErr := dbMongoAddTagsToPost(postTitleStr, tagsLst, pRuntimeSys)
				return gfErr

			} else {
				gfErr := gf_core.ErrorCreate(fmt.Sprintf("post with title (%s) doesnt exist, while adding a tags - %s", 
					postTitleStr,
					tagsLst),
					"verify__invalid_value_error",
					map[string]interface{}{
						"post_title_str": postTitleStr,
						"tags_lst":       tagsLst,
					},
					nil, "gf_tagger", pRuntimeSys)
				return gfErr
			}

		//---------------
		// IMAGE
		case "image":
			imageIDstr := pObjectExternIDstr
			imageID    := gf_images_core.GFimageID(imageIDstr)
			exists_bool, gfErr := gf_images_core.DBmongoImageExists(imageID, pCtx, pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			if exists_bool {
				gfErr := dbMongoAddTagsToImage(imageIDstr, tagsLst, pRuntimeSys)
				if gfErr != nil {
					return gfErr
				}
			}

		//---------------
		// WEB3
		case "address":

			chainStr := pMetaMap["chain_str"].(string)

			addressStr := pObjectExternIDstr
			existsBool, gfErr := gf_address.DBmongoExists(addressStr,
				chainStr,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}

			if existsBool {
				gfErr := gf_address.DBmongoAddTag(tagsLst,
					addressStr,
					chainStr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return gfErr
				}
			}

		//---------------
	}
	return nil
}

//---------------------------------------------------

func getObjectsWithTags(pTagsLst []string,
	pObjectTypeStr string,
	pPageIndexInt  int,
	pPageSizeInt   int,
	pRuntimeSys    *gf_core.RuntimeSys) (map[string][]map[string]interface{}, *gf_core.GFerror) {
		
	objectsWithTagsMap := map[string][]map[string]interface{}{}
	for _, tagStr := range pTagsLst {
		objectsWithTagLst, gfErr := getObjectsWithTag(tagStr,
			pObjectTypeStr,
			pPageIndexInt,
			pPageSizeInt,
			pRuntimeSys)

		if gfErr != nil {
			return nil, gfErr
		}
		objectsWithTagsMap[tagStr] = objectsWithTagLst
	}
	return objectsWithTagsMap, nil
}

//---------------------------------------------------

func getObjectsWithTag(pTagStr string,
	pObjectTypeStr string,
	pPageIndexInt  int,
	pPageSizeInt   int,
	pRuntimeSys    *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	// ADD!! - add support for tagging "image" pObjectTypeStr's
	if pObjectTypeStr != "post" {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("trying to get objects with a tag (%s) for objects type thats not supported - %s", pTagStr, pObjectTypeStr),
			"verify__invalid_value_error",
			map[string]interface{}{
				"tag_str":         pTagStr,
				"object_type_str": pObjectTypeStr,
			},
			nil, "gf_tagger", pRuntimeSys)
		return nil, gfErr
	}
	
	postsWithTagLst, gfErr := dbMongoGetPostsWithTag(pTagStr,
		pPageIndexInt,
		pPageSizeInt,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// package up info of each post that was found with tag 
	minPostsInfosLst := []map[string]interface{}{}
	for _, post := range postsWithTagLst {
		postInfoMap := map[string]interface{}{
			"title_str":               post.TitleStr,
			"tags_lst":                post.TagsLst,
			"url_str":                 fmt.Sprintf("/posts/%s", post.TitleStr),
			"object_type_str":         pObjectTypeStr,
			"thumbnail_small_url_str": post.ThumbnailURLstr,
		}
		minPostsInfosLst = append(minPostsInfosLst, postInfoMap)
	}

	objectsInfosLst := minPostsInfosLst
	return objectsInfosLst, nil
}

//---------------------------------------------------

func parseTags(pTagsStr string,
	pMaxTagsBulkSizeInt        int, // 500
	pMaxTagCharactersNumberInt int, // 20
	pRuntimeSys                *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {
	
	tagsLst := strings.Split(pTagsStr," ")
	//---------------------
	if len(tagsLst) > pMaxTagsBulkSizeInt {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("too many tags supplied - max is %d", pMaxTagsBulkSizeInt),
			"verify__value_too_many_error",
			map[string]interface{}{
				"tags_lst":               tagsLst,
				"max_tags_bulk_size_int": pMaxTagsBulkSizeInt,
			},
			nil, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}

	//---------------------
	for _, tagStr := range tagsLst {
		if len(tagStr) > pMaxTagCharactersNumberInt {
			gfErr := gf_core.ErrorCreate(fmt.Sprintf("tag (%s) is too long - max is (%d)", tagStr, pMaxTagCharactersNumberInt),
				"verify__string_too_long_error",
				map[string]interface{}{
					"tag_str":                       tagStr,
					"max_tag_characters_number_int": pMaxTagCharactersNumberInt,
				},
				nil, "gf_tagger_lib", pRuntimeSys)
			return nil, gfErr
		}
	}
	
	//---------------------
	return tagsLst, nil
}