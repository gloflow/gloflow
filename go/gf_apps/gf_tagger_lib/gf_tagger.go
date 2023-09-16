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
	"time"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_address"
)

//---------------------------------------------------
/*
TEMPORARY - this is mainly needed while tags are held as a property of images
	and discovered there in aggregate to get the total list.
	going forward tags are held in the SQL db and this function
	migrates/creates them in SQL if they dont already exist.

	in the future this function wont be necessary, unless there's some
	need for copying of tags from DB to DB.
*/

func pipelineCreateDiscoveredTags(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	// get all objects tags from the current Mongodb
	allTagsLst, gfErr := dbMongoGetAllObjectsTags(pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//----------------------

	// SQL
	for _, tagInfoMap := range allTagsLst {

		tagNameStr    := tagInfoMap["tag_str"].(string)
		objectID      := gf_core.GF_ID(tagInfoMap["id_str"].(string))
		objectTypeStr := tagInfoMap["t"].(string)
		userID        := gf_core.GF_ID(tagInfoMap["user_id_str"].(string))

		pRuntimeSys.LogNewFun("DEBUG", "creating tag if missing...",
			map[string]interface{}{
				"tag_name":    tagNameStr,
				"object_type": objectTypeStr,
				"user_id":     userID,
			})

		gfErr := CreateIfMissing([]string{tagNameStr},
			objectID,
			objectTypeStr,
			userID,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	return nil
}

//-------------------------------------------------
// CREATE_IF_MISSING

func CreateIfMissing(pTagsLst []string,
	pObjID         gf_core.GF_ID,
	pObjectTypeStr string,
	pUserID        gf_core.GF_ID,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {

	for _, tagStr := range pTagsLst {

		existsBool, gfErr := DBsqlCheckTagExists(tagStr, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		// create tag if it doesnt exist
		if !existsBool {

			//----------------------
			// CREATE_TAG
			gfErr := Create(tagStr,
				pObjID,
				pObjectTypeStr,
				pUserID,
				pCtx,
				pRuntimeSys)

			if gfErr != nil {
				return gfErr
			}

			//----------------------
		}
	}
	return nil
}

//---------------------------------------------------

func Create(pTagStr string,
	pObjID         gf_core.GF_ID,
	pObjectTypeStr string,
	pUserID        gf_core.GF_ID,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {

	tagID := generateTagID()

	// ADD!! - provide a mechanism for users to specify that a tag is private
	publicBool := true

	// DB
	gfErr := dbSQLcreateTag(tagID,
		pTagStr,
		pUserID,
		pObjID,
		pObjectTypeStr,
		publicBool,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	return nil
}

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
			nil, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}
	
	tagsLst, gfErr := parseTags(pTagsStr,
		500, // pMaxTagsBulkSizeInt
		20,  // pMaxTagCharactersNumberInt
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	pRuntimeSys.LogNewFun("DEBUG", "tags to be added to obj...",
		map[string]interface{}{
			"tags_lst":             tagsLst,
			"object_type_str":      pObjectTypeStr,
			"object_extern_id_str": pObjectExternIDstr,
		})

	//---------------
	// DB_SQL


	for _, tagStr := range tagsLst {
		
		objID := gf_core.GF_ID(pObjectExternIDstr)

		gfErr = Create(tagStr,
			objID,
			pObjectTypeStr,
			pUserID,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	//---------------
	switch pObjectTypeStr {
		//---------------
		// POST
		case "post":
			
			//---------------
			// FIX!! - post ID should not be its Title!!!!!
			postTitleStr := pObjectExternIDstr

			//---------------

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
					nil, "gf_tagger_lib", pRuntimeSys)
				return gfErr
			}

		//---------------
		// IMAGE
		case "image":

			imageIDstr := pObjectExternIDstr
			imageID    := gf_images_core.GFimageID(imageIDstr)
			existsBool, gfErr := gf_images_core.DBmongoImageExists(imageID, pCtx, pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			if existsBool {
				gfErr := dbMongoAddTagsToImage(imageIDstr, tagsLst, pCtx, pRuntimeSys)
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

func exportObjectsWithTags(pTagsLst []string,
	pObjectTypeStr string,
	pPageIndexInt  int,
	pPageSizeInt   int,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (map[string][]map[string]interface{}, *gf_core.GFerror) {
		
	objectsWithTagsMap := map[string][]map[string]interface{}{}
	for _, tagStr := range pTagsLst {
		objectsWithTagLst, gfErr := exportObjectsWithTag(tagStr,
			pObjectTypeStr,
			pPageIndexInt,
			pPageSizeInt,
			pCtx,
			pRuntimeSys)

		if gfErr != nil {
			return nil, gfErr
		}
		objectsWithTagsMap[tagStr] = objectsWithTagLst
	}
	return objectsWithTagsMap, nil
}

//---------------------------------------------------

func exportObjectsWithTag(pTagStr string,
	pObjectTypeStr string,
	pPageIndexInt  int,
	pPageSizeInt   int,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	if pObjectTypeStr != "post" && pObjectTypeStr != "image" {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("trying to get objects with a tag (%s) for objects type thats not supported - %s",
			pTagStr,
			pObjectTypeStr),
			"verify__invalid_value_error",
			map[string]interface{}{
				"tag_str":         pTagStr,
				"object_type_str": pObjectTypeStr,
			},
			nil, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	
	var objectsInfosLst []map[string]interface{}

	switch pObjectTypeStr {
	
	//---------------------
	// IMAGES
	case "image":

		var imagesLst []*gf_images_core.GFimage
		gfErr := dbMongoGetObjectsWithTag(pTagStr,
			"img",
			&imagesLst,
			pPageIndexInt,
			pPageSizeInt,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		// plugin
		metadataFilterDefinedBool := false
		if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback != nil {
			metadataFilterDefinedBool = true
		}
		
		exportedImagesInfosLst := []map[string]interface{}{}
		for _, image := range imagesLst {

			// META
			var filteredMetaJSONstr string
			if metadataFilterDefinedBool {
				filteredMetaMap := pRuntimeSys.ExternalPlugins.ImageFilterMetadataCallback(image.MetaMap)
				metaJSONbytesLst, _ := json.Marshal(filteredMetaMap)
				filteredMetaJSONstr = string(metaJSONbytesLst)
			}

			imageInfoMap := map[string]interface{}{
				"id_str": image.IDstr,            
				"creation_unix_time_str":    image.Creation_unix_time_f,
				"owner_user_id_str":         image.UserID,
				"image_origin_page_url_str": image.Origin_page_url_str,
				"thumbnail_medium_url_str":  image.Thumbnail_medium_url_str,
				"thumbnail_large_url_str":   image.Thumbnail_large_url_str,
				"format_str":                image.Format_str,
				"tags_lst":                  image.TagsLst,
				"meta_json_str":             filteredMetaJSONstr,    
			}
			exportedImagesInfosLst = append(exportedImagesInfosLst, imageInfoMap)
		}

		objectsInfosLst = exportedImagesInfosLst
	
	//---------------------
	// POSTS
	case "post":

		var postsLst []*gf_publisher_core.GFpost
		gfErr := dbMongoGetObjectsWithTag(pTagStr,
			"post",
			&postsLst,
			pPageIndexInt,
			pPageSizeInt,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		exportedPostsInfosLst := []map[string]interface{}{}
		for _, post := range postsLst {
			postInfoMap := map[string]interface{}{
				"title_str":               post.TitleStr,
				"tags_lst":                post.TagsLst,
				"url_str":                 fmt.Sprintf("/posts/%s", post.TitleStr),
				"thumbnail_small_url_str": post.ThumbnailURLstr,
			}
			exportedPostsInfosLst = append(exportedPostsInfosLst, postInfoMap)
		}

		objectsInfosLst = exportedPostsInfosLst
	}

	//---------------------

	return objectsInfosLst, nil
}

//---------------------------------------------------

func parseTags(pTagsStr string,
	pMaxTagsBulkSizeInt        int, // 500
	pMaxTagCharactersNumberInt int, // 20
	pRuntimeSys                *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {
	
	tagsLst := strings.Split(pTagsStr, " ")

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

//---------------------------------------------------

func generateTagID() gf_core.GF_ID {

	creationUNIXtimeF  := float64(time.Now().UnixNano())/1000000000.0
	randomStr          := gf_core.StrRandom()
	uniqueValsForIDlst := []string{
		randomStr,
	}
	sessionIDstr := gf_core.IDcreate(uniqueValsForIDlst, creationUNIXtimeF)
	return sessionIDstr
}

//-------------------------------------------------
// RESOLVE_USER_IDS_TO_USER_NAMES

func resolveUserIDStoUserNames(pImagesWithTagLst []map[string]interface{},
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) []gf_identity_core.GFuserName {

	usernamesCacheMap := map[gf_core.GF_ID]gf_identity_core.GFuserName{}
	

	userNamesLst := []gf_identity_core.GFuserName{}
	for _, imageMap := range pImagesWithTagLst {
		
		
		userID := imageMap["owner_user_id_str"].(gf_core.GF_ID)
		var userNameStr gf_identity_core.GFuserName

		// resolve user_id to user_name, or use cached result if its already present.
		if cachedUserNameStr, ok := usernamesCacheMap[userID]; ok {
			userNameStr = cachedUserNameStr
		} else {
			resolvedUserNameStr := gf_identity_core.ResolveUserName(userID, pCtx, pRuntimeSys)
			userNameStr               = resolvedUserNameStr
			usernamesCacheMap[userID] = resolvedUserNameStr
		}
		userNamesLst = append(userNamesLst, userNameStr)
	}
	
	return userNamesLst
}