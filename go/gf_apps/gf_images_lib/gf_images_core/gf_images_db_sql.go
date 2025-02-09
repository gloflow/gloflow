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

func DBsqlGetImagesFlows(pImagesIDsLst []GFimageID,
	pCtx		context.Context,
	pRuntimeSys	*gf_core.RuntimeSys) (map[GFimageID][]string, *gf_core.GFerror) {

	imagesIDsLst := []string{}
	for _, id := range pImagesIDsLst {
		idStr := string(id)
		imagesIDsLst = append(imagesIDsLst, idStr)
	}

	// SQL query to retrieve flows_names grouped by image ID
	sqlStr := `
		SELECT id, flows_names
		FROM gf_images
		WHERE id = ANY($1)
	`

	// Execute the query and collect results
	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr, pq.Array(imagesIDsLst))
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to execute SQL query",
			"sql_query_execute",
			map[string]interface{}{"images_ids_lst": imagesIDsLst},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	flowsByImageMap := make(map[GFimageID][]string)
	for rows.Next() {
		var idStr string
		var flowsNamesLst []string
		if err := rows.Scan(&idStr, pq.Array(&flowsNamesLst)); err != nil {
			gfErr := gf_core.ErrorCreate("failed to scan row for query to get flows_names for images",
				"sql_row_scan",
				map[string]interface{}{"images_ids_lst": imagesIDsLst},
				err, "gf_images_core", pRuntimeSys)
			return nil, gfErr
		}
		id := GFimageID(idStr)
		flowsByImageMap[id] = append(flowsByImageMap[id], flowsNamesLst...)
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("rows iteration error",
			"sql_query_execute",
			map[string]interface{}{"images_ids_lst": imagesIDsLst},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	return flowsByImageMap, nil
}

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
		pImage.IDstr,                   // id
		pImage.UserID,                  // user_id
		pImage.ClientTypeStr,           // client_type
		pImage.TitleStr,                // title
		pq.Array(pImage.FlowsNamesLst), // flows_names
		pImage.Origin_url_str,        // origin_url
		pImage.Origin_page_url_str,   // origin_page_url
		pImage.ThumbnailSmallURLstr,  // thumb_small_url
		pImage.ThumbnailMediumURLstr, // thumb_medium_url
		pImage.ThumbnailLargeURLstr,  // thumb_large_url
		pImage.Format_str,            // format
		pImage.Width_int,             // width
		pImage.Height_int,            // height
		jsonMetaStr,                  // meta_map
		pq.Array(pImage.TagsLst),     // tags_lst
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

func DBsqlGetImage(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	pRuntimeSys.LogNewFun("DEBUG", "retrieving image data...", map[string]interface{}{
		"image_id": pImageIDstr,
	})

	sqlStr := `
		SELECT 
			id,
			EXTRACT(EPOCH FROM creation_time) AS creation_unix_time,
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

			dominant_color_hex,
			palette_colors_hex,
			meta_map,
			tags_lst
		FROM gf_images
		WHERE id = $1 AND deleted = FALSE
		LIMIT 1`
	
	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr, pImageIDstr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get a single image from DB",
			"sql_query_execution",
			map[string]interface{}{
				"image_id_str": pImageIDstr,
			},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	for rows.Next() {
		img, gfErr := LoadImageFromResult(rows, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		// we're getting just a single image, so return right away
		return img, nil
	}

	return nil, nil
}

//---------------------------------------------------

func LoadImageFromResult(pRows *sql.Rows,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	img := &GFimage{}
	
	var clientTypeStr sql.NullString
	var originPageURLstr sql.NullString
	var thumbSmallURLstr, thumbMediumURLstr, thumbLargeURLstr sql.NullString
	var dominantColorHexStr, palleteStr sql.NullString
	var metaMapRaw []byte

	if err := pRows.Scan(
			&img.IDstr,
			&img.Creation_unix_time_f,
			&img.UserID,

			&clientTypeStr,
			&img.TitleStr,
			pq.Array(&img.FlowsNamesLst),

			&img.Origin_url_str,
			&originPageURLstr,

			&thumbSmallURLstr,
			&thumbMediumURLstr,
			&thumbLargeURLstr,

			&img.Format_str,
			&img.Width_int,
			&img.Height_int,

			&dominantColorHexStr,
			&palleteStr,
			&metaMapRaw,
			pq.Array(&img.TagsLst)); err != nil {
				
		gfErr := gf_core.ErrorCreate("failed to scan a row of images",
			"sql_row_scan",
			map[string]interface{}{},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	
	img.ClientTypeStr       = gf_core.DBsqlGetNullStringOrDefault(clientTypeStr, "")
	img.Origin_page_url_str = gf_core.DBsqlGetNullStringOrDefault(originPageURLstr, "")

	img.ThumbnailSmallURLstr  = gf_core.DBsqlGetNullStringOrDefault(thumbSmallURLstr, "")
	img.ThumbnailMediumURLstr = gf_core.DBsqlGetNullStringOrDefault(thumbMediumURLstr, "")
	img.ThumbnailLargeURLstr  = gf_core.DBsqlGetNullStringOrDefault(thumbLargeURLstr, "")

	img.DominantColorHexStr = gf_core.DBsqlGetNullStringOrDefault(dominantColorHexStr, "")
	img.PalleteStr = gf_core.DBsqlGetNullStringOrDefault(palleteStr, "")

	

	// META_MAP
	if err := json.Unmarshal(metaMapRaw, &img.MetaMap); err != nil {
		gfErr := gf_core.ErrorCreate("failed to unmarshal JSON meta_map",
			"json_decode_error",
			map[string]interface{}{},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}


	return img, nil
}

//---------------------------------------------------
// IMAGE_EXISTS_BY_ID

func DBsqlImageExistsByID(pImageID GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	queryStr := "SELECT COUNT(*) FROM gf_images WHERE id = $1 AND deleted = FALSE"

	var countInt int
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, queryStr, pImageID).Scan(&countInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if image exists in the DB",
			"sql_query_execute",
			map[string]interface{}{"image_id_str": pImageID},
			err, "gf_images_core", pRuntimeSys)
		return false, gfErr
	}

	return countInt > 0, nil
}

//---------------------------------------------------
// IMAGES_EXIST_BY_URLS

func DBsqlImagesExistByURLs(pImagesExternURLsLst []string,
	pFlowNameStr   string,
	// pClientTypeStr string,
	pUserID        gf_core.GF_ID,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	
	sqlStr := `
		SELECT
			creation_time,
			id,
			origin_url,
			origin_page_url,
			flows_names,
			tags_lst
		FROM
			gf_images
		WHERE
			(
				(
					$1 = 'all' AND
					(user_id = $3 OR user_id = 'anon') AND
					origin_url = ANY($2)
				)
				OR
				(
					$1 != 'all'
					AND
					(	
						user_id = $3 OR user_id = 'anon'
					)
					AND
					origin_url = ANY($2)
					AND
					$1 = ANY(flows_names)
				)
			);
	  `
  
	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr,
		pFlowNameStr,                   // $1
		pq.Array(pImagesExternURLsLst), // $2
		pUserID)                        // $3
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to execute SQL query to check if images exist",
			"sql_query_execute",
			map[string]interface{}{
				"images_extern_urls_lst": pImagesExternURLsLst,
				"flow_name_str":          pFlowNameStr,
				// "client_type_str":        pClientTypeStr,
			},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	var existingImagesLst []map[string]interface{}
	for rows.Next() {
		var creationTime time.Time
		var idStr, originURLStr, originPageURLStr string
		var flowsNamesLst, tagsLst []string

		if err := rows.Scan(&creationTime, &idStr, &originURLStr, &originPageURLStr, pq.Array(&flowsNamesLst), pq.Array(&tagsLst)); err != nil {
			gfErr := gf_core.ErrorCreate("failed to scan row for images exist check",
				"sql_row_scan",
				map[string]interface{}{
					"images_extern_urls_lst": pImagesExternURLsLst,
					"flow_name_str":          pFlowNameStr,
					// "client_type_str":        pClientTypeStr,
				},
				err, "gf_images_flows", pRuntimeSys)
			return nil, gfErr
		}

		var creationUNIXtimeF = float64(creationTime.Unix()) + float64(creationTime.Nanosecond())/1e9
		existingImagesLst = append(existingImagesLst, map[string]interface{}{
			"creation_unix_time_f": creationUNIXtimeF,
			"id_str":               idStr,
			"origin_url_str":       originURLStr,
			"origin_page_url_str":  originPageURLStr,
			"flows_names_lst":      flowsNamesLst,
			"tags_lst":             tagsLst,
		})
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("error encountered while iterating over query results",
			"sql_row_scan",
			map[string]interface{}{
				"images_extern_urls_lst": pImagesExternURLsLst,
				"flow_name_str":          pFlowNameStr,
				// "client_type_str":        pClientTypeStr,
			},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}

	return existingImagesLst, nil
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

	queryStr := `
		SELECT
			id,
			EXTRACT(EPOCH FROM creation_time) AS creation_unix_time,
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

			dominant_color_hex,
			palette_colors_hex,
			meta_map,
			tags_lst
			
		FROM gf_images 
		WHERE 
				creation_time IS NOT NULL 
			AND 
				$1 = ANY(flows_names)
		LIMIT $2 
		OFFSET $3;`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, queryStr,
		pFlowNameStr,
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
		
		img, gfErr := LoadImageFromResult(rows, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
		
		imgsLst = append(imgsLst, img)
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
// ADD_TAGS_TO_IMAGE

func DBsqlAddTagsToImage(pImageID GFimageID,
	pTagsLst    []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//--------------------
	// PUSH_TAGS
	// Extend `tags_lst` with the new tags
	queryStr := `
		UPDATE gf_images
		SET tags_lst = array_cat(tags_lst, $2::text[])
		WHERE id = $1;`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, queryStr, string(pImageID),
		pq.Array(pTagsLst))
	
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update a gf_image with new tags in DB",
			"sql_query_execute",
			map[string]interface{}{
				"image_id_str": string(pImageID),
				"tags_lst":     pTagsLst,
			},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}

	//--------------------
	return nil
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
		
		-- internal url of the original image file (relative url), unprocessed in its original form.
		-- this is the image that is stored in the system, and is used to generate thumbs;
		-- never served directly to the user.
		original_file_int_url TEXT,

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