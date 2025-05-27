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

package gf_tagger_lib

import (
	"context"

	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/lib/pq"
)

//---------------------------------------------------

// dbSQLgetObjectsWithTag retrieves objects (currently only images) with a given tag, paginated.
// It queries gf_tags for object IDs, then fetches images by those IDs.
func dbSQLgetObjectsWithTag(pTagStr string,
	pTargetTypeStr string,
	pOutput        interface{},
	pPageIndexInt  int,
	pPageSizeInt   int,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {

	switch pTargetTypeStr {
	case "img", "image":
		imagesPtr, ok := pOutput.(*[]*gf_images_core.GFimage)
		if !ok {
			return gf_core.ErrorCreate("pOutput is not of type *[]*gf_images_core.GFimage",
				"verify__invalid_input_struct_error",
				map[string]interface{}{},
				nil, "gf_tagger_lib", pRuntimeSys)
		}

		// 1. Query gf_tags for all image IDs with the given tag
		tagIDsQuery := `
			SELECT target_obj_id
			FROM gf_tags
			WHERE name = $1 AND target_obj_type = $2 AND deleted = FALSE
			ORDER BY creation_time ASC
			OFFSET $3 LIMIT $4
		`
		rows, err := pRuntimeSys.SQLdb.QueryContext(
			pCtx,
			tagIDsQuery,
			pTagStr,
			"image",
			pPageIndexInt,
			pPageSizeInt,
		)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to get tagged object IDs from gf_tags",
				"sql_query_execute",
				map[string]interface{}{
					"tag_str":         pTagStr,
					"target_type_str": pTargetTypeStr,
					"page_index_int":  pPageIndexInt,
					"page_size_int":   pPageSizeInt,
				},
				err, "gf_tagger_lib", pRuntimeSys)
			return gfErr
		}
		defer rows.Close()

		var imageIDs []gf_images_core.GFimageID
		for rows.Next() {
			var idStr string
			if err := rows.Scan(&idStr); err != nil {
				gfErr := gf_core.ErrorCreate("failed to scan target_obj_id from gf_tags",
					"sql_row_scan",
					map[string]interface{}{},
					err, "gf_tagger_lib", pRuntimeSys)
				return gfErr
			}
			imageIDs = append(imageIDs, gf_images_core.GFimageID(idStr))
		}
		if err := rows.Err(); err != nil {
			gfErr := gf_core.ErrorCreate("rows iteration error for tagged object IDs",
				"sql_query_execute",
				map[string]interface{}{
					"tag_str":         pTagStr,
					"target_type_str": pTargetTypeStr,
					"page_index_int":  pPageIndexInt,
					"page_size_int":   pPageSizeInt,
				},
				err, "gf_tagger_lib", pRuntimeSys)
			return gfErr
		}

		if len(imageIDs) == 0 {
			*imagesPtr = []*gf_images_core.GFimage{}
			return nil
		}

		// 2. Fetch each image by ID using DBsqlGetImage
		var images []*gf_images_core.GFimage
		for _, imgID := range imageIDs {
			img, gfErr := gf_images_core.DBsqlGetImage(imgID, pCtx, pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}
			if img != nil {
				images = append(images, img)
			}
		}
		*imagesPtr = images
		return nil

	default:
		return gf_core.ErrorCreate("unsupported target type for dbSQLgetObjectsWithTag",
			"verify__invalid_input_struct_error",
			map[string]interface{}{
				"target_type_str": pTargetTypeStr,
			},
			nil, "gf_tagger_lib", pRuntimeSys)
	}
}

//---------------------------------------------------

func pqArrayGFimageID(pIDsLst []gf_images_core.GFimageID) interface{} {
	strs := make([]string, len(pIDsLst))
	for i, id := range pIDsLst {
		strs[i] = string(id)
	}
	return pq.Array(strs)
}

//---------------------------------------------------

func dbSQLcreateTag(pTagID gf_core.GF_ID,
	pTagNameStr       string,
	pCreatorUserID    gf_core.GF_ID,
	pTargetObjID      gf_core.GF_ID,
	pTargetObjTypeStr string,
	pPublicBool       bool,
	pCtx              context.Context,
	pRuntimeSys       *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_tags (
			v,
			id,
			name,
			creator_user_id,
			public,
			target_obj_id,
			target_obj_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		"0",
		pTagID,
		pTagNameStr,
		pCreatorUserID,
		pPublicBool,
		pTargetObjID,
		pTargetObjTypeStr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert tag into the DB",
			"sql_query_execute",
			map[string]interface{}{
				"tag_name_str":        pTagNameStr,
				"creator_user_id_str": pCreatorUserID,
				"target_obj_id_str":   pTargetObjID,
				"target_obj_type_str": pTargetObjTypeStr,
				"public_bool":         pPublicBool,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// CHECK_FLOW_EXISTS

func DBsqlCheckTagExists(pTagStr string,
	pTargetObjID gf_core.GF_ID,
	pRuntimeSys  *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	db := pRuntimeSys.SQLdb

	var existsBool bool
	sqlStr := `SELECT exists(SELECT 1 FROM gf_tags WHERE name=$1 AND target_obj_id=$2)`
	err := db.QueryRow(sqlStr, pTagStr, string(pTargetObjID)).Scan(&existsBool)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if a tag exists in the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_tagger_lib", pRuntimeSys)
		return false, gfErr
	}
	return existsBool, nil
}

//---------------------------------------------------
// CREATE_TABLES

func dbSQLcreateTables(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_tags (
		v               VARCHAR(255),
		id              TEXT,
		deleted         BOOLEAN DEFAULT FALSE,
		creation_time   TIMESTAMP DEFAULT NOW(),
		name            TEXT NOT NULL,
		creator_user_id TEXT NOT NULL,
		public          BOOLEAN,

		-- object that is tagged, its ID and type --
		target_obj_id   TEXT,
		target_obj_type VARCHAR(30),

		PRIMARY KEY(id)
	);
	`

	_, err := pRuntimeSys.SQLdb.Exec(sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create tags related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}
