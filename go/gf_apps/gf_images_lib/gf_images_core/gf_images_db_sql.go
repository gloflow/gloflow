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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"

)

//---------------------------------------------------

func dbSQLputImage(pImage *GFimage,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	pRuntimeSys.LogNewFun("DEBUG", "upserting image data...", map[string]interface{}{
		"image_id_str": pImage.IDstr,
	})

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

	_, err := pRuntimeSys.SQLdb.ExecContext(
		pCtx,
		sqlStr,
		pImage.IDstr,                 // id
		pImage.UserID,                // user_id
		pImage.ClientTypeStr,         // client_type
		pImage.TitleStr,              // title
		pImage.FlowsNamesLst,         // flows_names
		pImage.Origin_url_str,        // origin_url
		pImage.Origin_page_url_str,   // origin_page_url
		pImage.ThumbnailSmallURLstr,  // thumb_small_url
		pImage.ThumbnailMediumURLstr, // thumb_medium_url
		pImage.ThumbnailLargeURLstr,  // thumb_large_url
		pImage.Format_str,            // format
		pImage.Width_int,             // width
		pImage.Height_int,            // height
		pImage.MetaMap,               // meta_map
		pImage.TagsLst,               // tags_lst
	)

	if err != nil {
		gfErr := gf_core.ErrorCreate(
			"failed to upsert image data in gf_images table",
			"sql_upsert_execute",
			map[string]interface{}{
				"image_id_str": pImage.IDstr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

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
		meta_map JSONB, -- metadata external users might assign to an image
		tags_lst TEXT[], -- human facing tags assigned to an image

		-- ---------------

		PRIMARY KEY(id),
		FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);
	`






	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create gf_identity related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}