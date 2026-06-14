/*
GloFlow application and media management/publishing platform
Copyright (C) 2026 Ivan Trajkovic

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

package gf_bookmarks

import (
	"context"
	"database/sql"
	"time"

	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/lib/pq"
)

//---------------------------------------------------
// HYBRID READ - reads from both SQL and Mongo, deduplicates results
// Eventually everything will be migrated to SQL

func dbGetAllBookmarks(pUserIDstr gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFbookmark, *gf_core.GFerror) {

	// Read from SQL
	sqlBookmarks, gfErr := dbSQLgetAllBookmarks(pUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		// If SQL fails, fall back to Mongo only
		pRuntimeSys.LogFun("WARN", "failed to read bookmarks from SQL, using Mongo only")
		return dbMongoGetAllBookmarks(pUserIDstr, pCtx, pRuntimeSys)
	}

	// Read from Mongo
	mongoBookmarks, gfErr := dbMongoGetAllBookmarks(pUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		// If Mongo fails, just return SQL results
		pRuntimeSys.LogFun("WARN", "failed to read bookmarks from Mongo, using SQL only")
		return sqlBookmarks, nil
	}

	// Deduplicate by ID (SQL takes precedence)
	bookmarksMap := make(map[gf_core.GF_ID]*GFbookmark)

	// Add Mongo bookmarks first
	for _, bookmark := range mongoBookmarks {
		bookmarksMap[bookmark.Id_str] = bookmark
	}

	// Add/override with SQL bookmarks (SQL is source of truth for new data)
	for _, bookmark := range sqlBookmarks {
		bookmarksMap[bookmark.Id_str] = bookmark
	}

	// Convert map back to slice, sorted by creation time (descending)
	var bookmarksLst []*GFbookmark
	for _, bookmark := range bookmarksMap {
		bookmarksLst = append(bookmarksLst, bookmark)
	}

	// Sort by creation time descending
	for i := 0; i < len(bookmarksLst)-1; i++ {
		for j := i + 1; j < len(bookmarksLst); j++ {
			if bookmarksLst[i].Creation_unix_time_f < bookmarksLst[j].Creation_unix_time_f {
				bookmarksLst[i], bookmarksLst[j] = bookmarksLst[j], bookmarksLst[i]
			}
		}
	}

	return bookmarksLst, nil
}

//---------------------------------------------------
// SQL FUNCTIONS
//---------------------------------------------------
// CREATE

func dbSQLcreateBookmark(pBookmark *GFbookmark,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_bookmarks (
			v,
			id,
			deleted,
			creation_time,
			user_id,
			url,
			description,
			tags,
			screenshot_image_id,
			screenshot_image_thumbnail_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	creationTime := time.Unix(int64(pBookmark.Creation_unix_time_f), 0)

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		pBookmark.V_str,
		pBookmark.Id_str,
		pBookmark.Deleted_bool,
		creationTime,
		pBookmark.User_id_str,
		pBookmark.Url_str,
		pBookmark.Description_str,
		pq.Array(pBookmark.Tags_lst),
		pBookmark.Screenshot_image_id_str,
		pBookmark.Screenshot_image_thumbnail_url_str)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert bookmark into the DB",
			"sql_query_execute",
			map[string]interface{}{
				"bookmark_id_str": pBookmark.Id_str,
				"user_id_str":     pBookmark.User_id_str,
				"url_str":         pBookmark.Url_str,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// UPDATE

func DBsqlUpdateBookmark(pBookmarkIDstr gf_core.GF_ID,
	pUserIDstr gf_core.GF_ID,
	pDescriptionStr string,
	pTagsLst []string,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		UPDATE gf_bookmarks
		SET description = $1,
		    tags = $2
		WHERE id = $3 AND user_id = $4 AND deleted = FALSE`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		pDescriptionStr,
		pq.Array(pTagsLst),
		pBookmarkIDstr,
		pUserIDstr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update bookmark in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"bookmark_id_str": pBookmarkIDstr,
				"user_id_str":     pUserIDstr,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// UPDATE_SCREENSHOT

func dbSQLupdateBookmarkScreenshot(pBookmarkIDstr gf_core.GF_ID,
	pScreenshotImageIDstr gf_images_core.GFimageID,
	pScreenshotImageThumbnailURLstr string,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		UPDATE gf_bookmarks
		SET screenshot_image_id = $1,
		    screenshot_image_thumbnail_url = $2
		WHERE id = $3 AND deleted = FALSE`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		pScreenshotImageIDstr,
		pScreenshotImageThumbnailURLstr,
		pBookmarkIDstr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update bookmark screenshot in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"bookmark_id_str":                    pBookmarkIDstr,
				"screenshot_image_id_str":            pScreenshotImageIDstr,
				"screenshot_image_thumbnail_url_str": pScreenshotImageThumbnailURLstr,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// GET_ALL

func dbSQLgetAllBookmarks(pUserIDstr gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFbookmark, *gf_core.GFerror) {

	sqlStr := `
		SELECT 
			v,
			id,
			deleted,
			creation_time,
			user_id,
			url,
			description,
			tags,
			screenshot_image_id,
			screenshot_image_thumbnail_url
		FROM gf_bookmarks
		WHERE user_id = $1 AND deleted = FALSE
		ORDER BY creation_time DESC`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr, pUserIDstr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to query bookmarks from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str": pUserIDstr,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	var bookmarksLst []*GFbookmark
	for rows.Next() {
		bookmark := &GFbookmark{}
		var creationTime time.Time
		var tagsArray pq.StringArray
		var screenshotImageIDstr sql.NullString
		var screenshotImageThumbnailURLstr sql.NullString

		err := rows.Scan(
			&bookmark.V_str,
			&bookmark.Id_str,
			&bookmark.Deleted_bool,
			&creationTime,
			&bookmark.User_id_str,
			&bookmark.Url_str,
			&bookmark.Description_str,
			&tagsArray,
			&screenshotImageIDstr,
			&screenshotImageThumbnailURLstr,
		)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to scan bookmark row from DB",
				"sql_row_scan",
				map[string]interface{}{
					"user_id_str": pUserIDstr,
				},
				err, "gf_tagger_lib", pRuntimeSys)
			return nil, gfErr
		}

		bookmark.Creation_unix_time_f = float64(creationTime.Unix())
		bookmark.Tags_lst = tagsArray

		if screenshotImageIDstr.Valid {
			bookmark.Screenshot_image_id_str = gf_images_core.GFimageID(screenshotImageIDstr.String)
		}
		if screenshotImageThumbnailURLstr.Valid {
			bookmark.Screenshot_image_thumbnail_url_str = screenshotImageThumbnailURLstr.String
		}

		bookmarksLst = append(bookmarksLst, bookmark)
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("rows iteration error for bookmarks query",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str": pUserIDstr,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}

	return bookmarksLst, nil
}

//---------------------------------------------------
// GET_BY_ID

func dbSQLgetBookmarkByID(pBookmarkIDstr gf_core.GF_ID,
	pUserIDstr gf_core.GF_ID,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFbookmark, *gf_core.GFerror) {

	sqlStr := `
		SELECT 
			v,
			id,
			deleted,
			creation_time,
			user_id,
			url,
			description,
			tags,
			screenshot_image_id,
			screenshot_image_thumbnail_url
		FROM gf_bookmarks
		WHERE id = $1 AND user_id = $2 AND deleted = FALSE`

	bookmark := &GFbookmark{}
	var creationTime time.Time
	var tagsArray pq.StringArray
	var screenshotImageIDstr sql.NullString
	var screenshotImageThumbnailURLstr sql.NullString

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pBookmarkIDstr, pUserIDstr).Scan(
		&bookmark.V_str,
		&bookmark.Id_str,
		&bookmark.Deleted_bool,
		&creationTime,
		&bookmark.User_id_str,
		&bookmark.Url_str,
		&bookmark.Description_str,
		&tagsArray,
		&screenshotImageIDstr,
		&screenshotImageThumbnailURLstr,
	)

	if err == sql.ErrNoRows {
		gfErr := gf_core.ErrorCreate("bookmark not found",
			"bookmark_not_found",
			map[string]interface{}{
				"bookmark_id_str": pBookmarkIDstr,
				"user_id_str":     pUserIDstr,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to query bookmark by ID from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"bookmark_id_str": pBookmarkIDstr,
				"user_id_str":     pUserIDstr,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}

	bookmark.Creation_unix_time_f = float64(creationTime.Unix())
	bookmark.Tags_lst = tagsArray

	if screenshotImageIDstr.Valid {
		bookmark.Screenshot_image_id_str = gf_images_core.GFimageID(screenshotImageIDstr.String)
	}
	if screenshotImageThumbnailURLstr.Valid {
		bookmark.Screenshot_image_thumbnail_url_str = screenshotImageThumbnailURLstr.String
	}

	return bookmark, nil
}
