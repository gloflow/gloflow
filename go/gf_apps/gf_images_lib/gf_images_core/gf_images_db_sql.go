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

package gf_images_core

import (
	"fmt"
	"context"
	"time"
	"encoding/json"
	"math/rand"
	"database/sql"
	"github.com/lib/pq"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// PUT_IMAGE

func DBsqlPutImage(pImage *GFimage,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	pRuntimeSys.LogNewFun("DEBUG", "upserting image data...", map[string]interface{}{
		"image_id_str": pImage.IDstr,
	})

	//----------------------
	// META_MAP - convert to JSON; meta_map column is of type JSONB
	jsonMetaBytesLst, err := json.Marshal(pImage.MetaMap)
	if err != nil {
		
		gfErr := gf_core.ErrorCreate(
			"failed to json-encode image meta_map, to persist it",
			"json_encode_error",
			map[string]interface{}{
				"image_id_str": pImage.IDstr,
			},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}
	jsonMetaStr := string(jsonMetaBytesLst)


	spew.Dump(pImage)

	//----------------------
	sqlStr := `
		INSERT INTO gf_images (
			id,
			user_id,
			client_type,
			title,
			flows_names, 
			origin_url,
			origin_page_url,
			thumb_small_url,
			thumb_medium_url, 
			thumb_large_url,
			format,
			width,
			height,
			meta_map,
			tags_lst
		)
		VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11, $12, $13, $14, $15
		);
	`

	_, err = pRuntimeSys.SQLdb.ExecContext(
		pCtx,
		sqlStr,
		pImage.IDstr,                 // id
		pImage.UserID,                // user_id
		pImage.ClientTypeStr,         // client_type
		pImage.TitleStr,              // title
		pq.Array(pImage.FlowsNamesLst),         // flows_names
		pImage.Origin_url_str,        // origin_url
		pImage.Origin_page_url_str,   // origin_page_url
		pImage.ThumbnailSmallURLstr,  // thumb_small_url
		pImage.ThumbnailMediumURLstr, // thumb_medium_url
		pImage.ThumbnailLargeURLstr,  // thumb_large_url
		pImage.Format_str,            // format
		pImage.Width_int,             // width
		pImage.Height_int,            // height
		jsonMetaStr,                  // meta_map
		pq.Array(pImage.TagsLst),               // tags_lst
	)

	if err != nil {
		gfErr := gf_core.ErrorCreate(
			"failed to upsert image data in gf_images table",
			"sql_query_execute",
			map[string]interface{}{
				"image_id_str": pImage.IDstr,
			},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// GET_IMAGE

func dbSQLGetImage(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	pRuntimeSys.LogNewFun("DEBUG", "retrieving image data...", map[string]interface{}{
		"image_id": pImageIDstr,
	})

	// SELECT SQL statement
	sqlStr := `
		SELECT 
			id,
			user_id,
			client_type,
			title,
			flows_names, 
			origin_url,
			origin_page_url,
			thumb_small_url,
			thumb_medium_url, 
			thumb_large_url,
			format,
			width,
			height,
			meta_map,
			tags_lst
		FROM gf_images
		WHERE id = $1 AND deleted = FALSE`

	var image GFimage
	row := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pImageIDstr)
	err := row.Scan(
		&image.IDstr,
		&image.UserID,
		&image.ClientTypeStr,
		&image.TitleStr,
		&image.FlowsNamesLst,
		&image.Origin_url_str,
		&image.Origin_page_url_str,
		&image.ThumbnailSmallURLstr,
		&image.ThumbnailMediumURLstr,
		&image.ThumbnailLargeURLstr,
		&image.Format_str,
		&image.Width_int,
		&image.Height_int,
		&image.MetaMap,
		&image.TagsLst,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			gfErr := gf_core.ErrorCreate(
				"image does not exist in gf_images table",
				"sql_query_execute",
				map[string]interface{}{
					"image_id": pImageIDstr,
				},
				err, "gf_images_core", pRuntimeSys)
			return nil, gfErr
		}

		gfErr := gf_core.ErrorCreate(
			"failed to retrieve image data from gf_images table",
			"sql_query_execute",
			map[string]interface{}{
				"image_id": pImageIDstr,
			},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	return &image, nil
}

//---------------------------------------------------
// IMAGE_EXISTS

func DBsqlImageExists(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	query := "SELECT COUNT(*) FROM gf_images WHERE id = ?"

	var count_int int
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, query, pImageIDstr).Scan(&count_int)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if image exists in the DB",
			"sql_query_execute",
			map[string]interface{}{"image_id_str": pImageIDstr},
			err, "gf_images_core", pRuntimeSys)
		return false, gfErr
	}

	return count_int > 0, nil
}

//---------------------------------------------------
// GET_RANDOM_IMAGES_RANGE

func DBsqlGetRandomImagesRange(pImgsNumToGetInt int, // 5
	pMaxRandomCursorPositionInt int, // 2000
	pFlowNameStr                string,
	pUserID                     gf_core.GF_ID,
	pCtx                        context.Context,
	pRuntimeSys                 *gf_core.RuntimeSys) ([]*GFimage, *gf_core.GFerror) {

	// Reseed the random number source
	rand.Seed(time.Now().UnixNano())
	randomCursorPositionInt := rand.Intn(pMaxRandomCursorPositionInt)

	pRuntimeSys.LogNewFun("DEBUG", "imgs_num_to_get_int        - "+fmt.Sprint(pImgsNumToGetInt), nil)
	pRuntimeSys.LogNewFun("DEBUG", "random_cursor_position_int - "+fmt.Sprint(randomCursorPositionInt), nil)

	query := `
		SELECT * FROM gf_images 
		WHERE creation_unix_time_f  IS NOT NULL 
			AND flows_names_lst     LIKE ? 
			AND origin_page_url_str IS NOT NULL 
		LIMIT ? 
		OFFSET ?`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, query,
		"%"+pFlowNameStr+"%",
		pImgsNumToGetInt,
		randomCursorPositionInt)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get random images range from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"imgs_num_to_get_int":            pImgsNumToGetInt,
				"max_random_cursor_position_int": pMaxRandomCursorPositionInt,
				"flow_name_str":                  pFlowNameStr,
			},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	var imgsLst []*GFimage
	for rows.Next() {
		var img GFimage

		if err := rows.Scan(&img.IDstr,
			&img.Creation_unix_time_f,
			&img.FlowsNamesLst,
			&img.Origin_page_url_str); err != nil {
			
			gfErr := gf_core.ErrorCreate("failed to scan row for random images",
				"sql_row_scan",
				map[string]interface{}{
					"imgs_num_to_get_int":            pImgsNumToGetInt,
					"max_random_cursor_position_int": pMaxRandomCursorPositionInt,
					"flow_name_str":                  pFlowNameStr,
				},
				err, "gf_images_core", pRuntimeSys)
			return nil, gfErr
		}
		imgsLst = append(imgsLst, &img)
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("error encountered during rows iteration",
			"sql_row_scan",
			map[string]interface{}{
				"imgs_num_to_get_int":            pImgsNumToGetInt,
				"max_random_cursor_position_int": pMaxRandomCursorPositionInt,
				"flow_name_str":                  pFlowNameStr,
			},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	return imgsLst, nil
}


//---------------------------------------------------
// TABLES
//---------------------------------------------------
// CREATE_TABLES

func DBsqlCreateTables(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_images (
		v                 VARCHAR(255),
		id                TEXT,
		deleted           BOOLEAN DEFAULT FALSE,
		creation_time     TIMESTAMP DEFAULT NOW(),
		user_id           TEXT,

		-- ---------------
		-- "gchrome_ext"|"gf_crawl_images"|"gf_image_editor"
		client_type VARCHAR(255),

		title TEXT,

		-- image can bellong to multiple flows
		flows_names TEXT[],

		-- ---------------
		/*
		RESOLVED_SOURCE_URL
		IMPORTANT!! - when the image comes from an external url (as oppose to it being 
			created internally, or uploaded directly to the system).
			this is different from Origin_page_url_str in that the page_url is the url 
			of the page in which the image is found, whereas this origin_url is the url
			of the file on some file server from which the image is served
		*/
		origin_url TEXT,

		-- if the image is extracted from a page, this holds the page_url
		origin_page_url TEXT,
		
		-- ---------------
		-- THUMBS
		-- relative url"s - "/images/image_name.*"
		thumb_small_url  TEXT,
		thumb_medium_url TEXT,
		thumb_large_url  TEXT,

		-- ---------------
		format TEXT, -- "jpeg" | "png" | "gif"
		width  INT,
		height INT,

		-- ---------------
		-- COLORS
		dominant_color_hex TEXT,
		palette_colors_hex TEXT[],
		
		-- ---------------
		-- META
		meta_map JSONB,  -- metadata external users might assign to an image
		tags_lst TEXT[], -- human facing tags assigned to an image

		-- ---------------

		PRIMARY KEY(id)

		-- for some of the tests to pass right now, we need to have a user_id column that
		-- accepts non-registered user-id''s. for ease of dev and testing.
		-- FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);
	`






	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create gf_identity related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}

	return nil
}