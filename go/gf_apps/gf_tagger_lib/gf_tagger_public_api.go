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
	"context"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_bookmarks"
)

//---------------------------------------------------
// PUBLIC API - Exported functions for use by external packages

// GetAllTags - get all tags from the system
func GetAllTags(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {
	return pipelineGetAllTags(pCtx, pRuntimeSys)
}

//---------------------------------------------------

// AddTagsToObject - add tags to an object (post, image, event, address)
func AddTagsToObject(pTagsStr string,
	pObjectTypeStr string,
	pObjectExternIDstr string,
	pMetaMap map[string]interface{},
	pUserID gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	return addTagsToObject(pTagsStr,
		pObjectTypeStr,
		pObjectExternIDstr,
		pMetaMap,
		pUserID,
		pCtx,
		pRuntimeSys)
}

//---------------------------------------------------
// NOTE: Notes functionality (CreateNote, GetNotes) has been moved to
// gloflow_emergence/gf_apps/gf_tagger_int for emergence-specific implementation

//---------------------------------------------------
// BOOKMARKS - Public API

// CreateBookmark - create a new bookmark
func CreateBookmark(pURLstr string,
	pDescriptionStr string,
	pTagsLst []string,
	pUserID gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	input := &gf_bookmarks.GFbookmarkInputCreate{
		User_id_str:     pUserID,
		Url_str:         pURLstr,
		Description_str: pDescriptionStr,
		Tags_lst:        pTagsLst,
	}

	gfErr := gf_bookmarks.PipelineCreate(input, pUserID, nil, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	// Generate the bookmark ID based on the same logic as in PipelineCreate
	creation_unix_time_f := float64(time.Now().UnixNano()) / 1000000000.0
	unique_vals_for_id_lst := []string{
		pURLstr,
		string(pUserID),
	}
	IDstr := gf_core.IDcreate(unique_vals_for_id_lst, creation_unix_time_f)

	return IDstr, nil
}

//---------------------------------------------------

// GetBookmarks - get all bookmarks for a user
func GetBookmarks(pUserID gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*gf_bookmarks.GFbookmarkExtern, *gf_core.GFerror) {

	input := &gf_bookmarks.GFbookmarkInputGet{
		Response_format_str: "json",
		User_id_str:         pUserID,
	}

	output, gfErr := gf_bookmarks.PipelineGet(input, nil, nil, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return output.Bookmarks_lst, nil
}

//---------------------------------------------------

// GetBookmarkByID - get a specific bookmark by its ID
func GetBookmarkByID(pBookmarkID gf_core.GF_ID,
	pUserID gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_bookmarks.GFbookmarkExtern, *gf_core.GFerror) {

	bookmarks, gfErr := GetBookmarks(pUserID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	for _, bookmark := range bookmarks {
		if bookmark.Id_str == pBookmarkID {
			return bookmark, nil
		}
	}

	gfErr = gf_core.ErrorCreate("bookmark not found",
		"bookmark_not_found",
		map[string]interface{}{
			"bookmark_id_str": pBookmarkID,
			"user_id_str":     pUserID,
		},
		nil, "gf_tagger_lib", pRuntimeSys)
	return nil, gfErr
}

//---------------------------------------------------

// UpdateBookmark - update a bookmark's description and/or tags
func UpdateBookmark(pBookmarkID gf_core.GF_ID,
	pUserID gf_core.GF_ID,
	pDescriptionStr string,
	pTagsLst []string,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	return gf_bookmarks.DBsqlUpdateBookmark(pBookmarkID, pUserID, pDescriptionStr, pTagsLst, pCtx, pRuntimeSys)
}

//---------------------------------------------------
