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
	"math"
	"time"

	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//---------------------------------------------------

type GFmigrationStats struct {
	TotalInt       int  // docs in the mongo collection
	ReadInt        int  // docs read from mongo
	InsertedInt    int  // new rows inserted into SQL (or that would be, in a dry run)
	OverwrittenInt int  // existing SQL rows overwritten with the mongo version
	IdenticalInt   int  // already in SQL, with identical values
	SkippedInt     int  // already in SQL with different values, left untouched
	InvalidInt     int  // docs that could not be migrated (missing id, etc.)
	AbortedBool    bool // migration was stopped early by the user
}

type GFmigrateDupAction int

const (
	GFmigrateDupSkip      GFmigrateDupAction = iota // leave the existing SQL record as is
	GFmigrateDupOverwrite                           // overwrite the SQL record with the mongo version
	GFmigrateDupQuit                                // stop the migration
)

// Callbacks that let the caller (CLI) handle user interaction and output,
// so that no terminal IO is done in the lib.
type GFmigrateCallbacks struct {

	// called when a mongo bookmark has the same id as a record already in SQL,
	// and the two differ. Not called in dry runs (those are treated as skipped).
	// If nil, duplicates are skipped.
	OnDuplicate func(pMongo *GFbookmark, pSQL *GFbookmark) GFmigrateDupAction

	// called for every record processed, with what happened to it.
	// pStatusStr is one of: "inserted", "overwritten", "identical", "skipped", "invalid"
	// (in a dry run, "inserted" means "would be inserted").
	OnRecord func(pIndexInt int, pTotalInt int, pStatusStr string, pBookmark *GFbookmark)
}

//---------------------------------------------------
// MIGRATE_MONGO_TO_SQL
// copies all docs from the mongo "gf_bookmarks" collection into the SQL gf_bookmarks table,
// oldest first.
//
// - all docs are copied, including ones marked deleted (the deleted flag is preserved).
// - a doc whose id already exists in SQL is a duplicate: if the values are identical its skipped,
//   otherwise pCallbacks.OnDuplicate decides if the SQL record is overwritten or left alone.
//   SQL records are never changed without that explicit decision.
// - mongo is not modified.
// - in a dry run nothing is written to SQL (and the SQL table is not created).

func DBmigrateMongoToSQL(pDryRunBool bool,
	pCallbacks  *GFmigrateCallbacks,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFmigrationStats, *gf_core.GFerror) {

	if pCallbacks == nil {
		pCallbacks = &GFmigrateCallbacks{}
	}

	collNameStr := "gf_bookmarks"
	coll        := pRuntimeSys.Mongo_db.Collection(collNameStr)

	// make sure the target table exists
	if !pDryRunBool {
		gfErr := DBsqlCreateTables(pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	stats := &GFmigrationStats{}

	totalInt64, err := coll.CountDocuments(pCtx, bson.M{})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to count bookmarks in mongo for migration",
			"mongodb_find_error",
			map[string]interface{}{"coll_name_str": collNameStr},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	stats.TotalInt = int(totalInt64)

	// OLDEST FIRST
	cursor, err := coll.Find(pCtx,
		bson.M{},
		options.Find().SetSort(bson.D{
			{Key: "creation_unix_time_f", Value: 1},
			{Key: "_id", Value: 1},
		}))
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to query all bookmarks from mongo for migration",
			"mongodb_find_error",
			map[string]interface{}{"coll_name_str": collNameStr},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	defer cursor.Close(pCtx)

	reportFun := func(pStatusStr string, pBookmark *GFbookmark) {
		if pCallbacks.OnRecord != nil {
			pCallbacks.OnRecord(stats.ReadInt, stats.TotalInt, pStatusStr, pBookmark)
		}
	}

	for cursor.Next(pCtx) {

		stats.ReadInt++

		var bookmark GFbookmark
		err := cursor.Decode(&bookmark)
		if err != nil {
			stats.InvalidInt++
			pRuntimeSys.LogNewFun("ERROR", "failed to decode bookmark mongo doc, skipping", map[string]interface{}{"err": err.Error()})
			reportFun("invalid", &bookmark)
			continue
		}

		if bookmark.Id_str == "" {
			stats.InvalidInt++
			reportFun("invalid", &bookmark)
			continue
		}

		//------------------
		// DUPLICATE CHECK
		existing, gfErr := dbSQLgetBookmarkByIDanyUser(bookmark.Id_str, pCtx, pRuntimeSys)
		if gfErr != nil {
			if pDryRunBool {
				// in a dry run the SQL table may not exist yet, so nothing can be a duplicate
				existing = nil
			} else {
				return stats, gfErr
			}
		}

		if existing != nil {

			// IDENTICAL - nothing to do
			if len(MigrateDiffFields(&bookmark, existing)) == 0 {
				stats.IdenticalInt++
				reportFun("identical", &bookmark)
				continue
			}

			// DIFFERENT - ask the caller what to do
			action := GFmigrateDupSkip
			if !pDryRunBool && pCallbacks.OnDuplicate != nil {
				action = pCallbacks.OnDuplicate(&bookmark, existing)
			}

			if action == GFmigrateDupQuit {
				stats.ReadInt--
				stats.AbortedBool = true
				return stats, nil
			}
			if action == GFmigrateDupSkip {
				stats.SkippedInt++
				reportFun("skipped", &bookmark)
				continue
			}

			gfErr := dbSQLoverwriteBookmark(&bookmark, pCtx, pRuntimeSys)
			if gfErr != nil {
				return stats, gfErr
			}
			stats.OverwrittenInt++
			reportFun("overwritten", &bookmark)
			continue
		}

		//------------------
		// INSERT
		if !pDryRunBool {
			gfErr := dbSQLinsertMigratedBookmark(&bookmark, pCtx, pRuntimeSys)
			if gfErr != nil {
				return stats, gfErr
			}
		}
		stats.InsertedInt++
		reportFun("inserted", &bookmark)
	}

	if err := cursor.Err(); err != nil {
		gfErr := gf_core.MongoHandleError("mongo cursor error while migrating bookmarks",
			"mongodb_cursor_error",
			map[string]interface{}{"coll_name_str": collNameStr},
			err, "gf_tagger_lib", pRuntimeSys)
		return stats, gfErr
	}

	return stats, nil
}

//---------------------------------------------------
// DIFF - names of the fields that differ between a mongo and a SQL version of a bookmark.

func MigrateDiffFields(pMongo *GFbookmark, pSQL *GFbookmark) []string {

	var diffLst []string

	if pMongo.V_str != pSQL.V_str {
		diffLst = append(diffLst, "v")
	}
	if pMongo.Deleted_bool != pSQL.Deleted_bool {
		diffLst = append(diffLst, "deleted")
	}
	// SQL TIMESTAMP has microsecond precision
	if math.Abs(pMongo.Creation_unix_time_f - pSQL.Creation_unix_time_f) > 0.001 {
		diffLst = append(diffLst, "creation_time")
	}
	if pMongo.User_id_str != pSQL.User_id_str {
		diffLst = append(diffLst, "user_id")
	}
	if pMongo.Url_str != pSQL.Url_str {
		diffLst = append(diffLst, "url")
	}
	if pMongo.Description_str != pSQL.Description_str {
		diffLst = append(diffLst, "description")
	}

	tagsDiffBool := len(pMongo.Tags_lst) != len(pSQL.Tags_lst)
	if !tagsDiffBool {
		for i := range pMongo.Tags_lst {
			if pMongo.Tags_lst[i] != pSQL.Tags_lst[i] {
				tagsDiffBool = true
				break
			}
		}
	}
	if tagsDiffBool {
		diffLst = append(diffLst, "tags")
	}

	if pMongo.Screenshot_image_id_str != pSQL.Screenshot_image_id_str {
		diffLst = append(diffLst, "screenshot_image_id")
	}
	if pMongo.Screenshot_image_thumbnail_url_str != pSQL.Screenshot_image_thumbnail_url_str {
		diffLst = append(diffLst, "screenshot_image_thumbnail_url")
	}
	return diffLst
}

//---------------------------------------------------
// SQL

// mongo stores unix time as float seconds, SQL as TIMESTAMP (no tz).
// stored as UTC wall-clock, which is how the SQL read path interprets it.
func migrateCreationTime(pUnixTimeF float64) time.Time {
	secF, fracF := math.Modf(pUnixTimeF)
	return time.Unix(int64(secF), int64(fracF*1e9)).UTC()
}

// NULLIF - empty screenshot fields are stored as NULL, same as bookmarks that never had a screenshot.
func dbSQLinsertMigratedBookmark(pBookmark *GFbookmark,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// "ON CONFLICT DO NOTHING" - safety net, duplicates are handled before this is called.
	sqlStr := `
		INSERT INTO gf_bookmarks (
			v, id, deleted, creation_time, user_id, url, description, tags,
			screenshot_image_id, screenshot_image_thumbnail_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULLIF($9, ''), NULLIF($10, ''))
		ON CONFLICT (id) DO NOTHING`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		pBookmark.V_str,
		pBookmark.Id_str,
		pBookmark.Deleted_bool,
		migrateCreationTime(pBookmark.Creation_unix_time_f),
		pBookmark.User_id_str,
		pBookmark.Url_str,
		pBookmark.Description_str,
		pq.Array(pBookmark.Tags_lst),
		string(pBookmark.Screenshot_image_id_str),
		pBookmark.Screenshot_image_thumbnail_url_str)
	if err != nil {
		return gf_core.ErrorCreate("failed to insert migrated bookmark into the SQL DB",
			"sql_query_execute",
			map[string]interface{}{
				"bookmark_id_str": pBookmark.Id_str,
				"user_id_str":     pBookmark.User_id_str,
				"url_str":         pBookmark.Url_str,
			},
			err, "gf_tagger_lib", pRuntimeSys)
	}
	return nil
}

func dbSQLoverwriteBookmark(pBookmark *GFbookmark,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		UPDATE gf_bookmarks
		SET v                              = $1,
		    deleted                        = $2,
		    creation_time                  = $3,
		    user_id                        = $4,
		    url                            = $5,
		    description                    = $6,
		    tags                           = $7,
		    screenshot_image_id            = NULLIF($8, ''),
		    screenshot_image_thumbnail_url = NULLIF($9, '')
		WHERE id = $10`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		pBookmark.V_str,
		pBookmark.Deleted_bool,
		migrateCreationTime(pBookmark.Creation_unix_time_f),
		pBookmark.User_id_str,
		pBookmark.Url_str,
		pBookmark.Description_str,
		pq.Array(pBookmark.Tags_lst),
		string(pBookmark.Screenshot_image_id_str),
		pBookmark.Screenshot_image_thumbnail_url_str,
		pBookmark.Id_str)
	if err != nil {
		return gf_core.ErrorCreate("failed to overwrite bookmark in the SQL DB",
			"sql_query_execute",
			map[string]interface{}{
				"bookmark_id_str": pBookmark.Id_str,
			},
			err, "gf_tagger_lib", pRuntimeSys)
	}
	return nil
}

// returns nil (and no error) if there is no bookmark with this id.
// matches regardless of user and of the deleted flag.
func dbSQLgetBookmarkByIDanyUser(pBookmarkIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFbookmark, *gf_core.GFerror) {

	sqlStr := `
		SELECT v, id, deleted, creation_time, user_id, url, description, tags,
			screenshot_image_id, screenshot_image_thumbnail_url
		FROM gf_bookmarks
		WHERE id = $1`

	bookmark := &GFbookmark{}
	var vStr, descriptionStr, screenshotIDstr, screenshotThumbURLstr sql.NullString
	var creationTime time.Time
	var tagsArray pq.StringArray

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pBookmarkIDstr).Scan(
		&vStr,
		&bookmark.Id_str,
		&bookmark.Deleted_bool,
		&creationTime,
		&bookmark.User_id_str,
		&bookmark.Url_str,
		&descriptionStr,
		&tagsArray,
		&screenshotIDstr,
		&screenshotThumbURLstr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, gf_core.ErrorCreate("failed to query bookmark by id from the SQL DB",
			"sql_query_execute",
			map[string]interface{}{"bookmark_id_str": pBookmarkIDstr},
			err, "gf_tagger_lib", pRuntimeSys)
	}

	bookmark.V_str                = vStr.String
	bookmark.Description_str      = descriptionStr.String
	bookmark.Creation_unix_time_f = float64(creationTime.UnixNano()) / 1e9
	bookmark.Tags_lst             = tagsArray
	bookmark.Screenshot_image_id_str            = gf_images_core.GFimageID(screenshotIDstr.String)
	bookmark.Screenshot_image_thumbnail_url_str = screenshotThumbURLstr.String

	return bookmark, nil
}

//---------------------------------------------------
// VERIFY_MONGO_TO_SQL
// READ-ONLY check that every doc in the mongo "gf_bookmarks" collection exists in the SQL
// gf_bookmarks table with identical values. Nothing is written to either DB.
// Also counts the SQL records that are not in mongo (ex. bookmarks created after the switch to SQL),
// which is informational and does not make verification fail.

type GFverifyStats struct {
	TotalInt     int // docs in the mongo collection
	ReadInt      int // docs read from mongo
	IdenticalInt int // in SQL with identical values
	DifferentInt int // in SQL, but with different values
	MissingInt   int // not in SQL
	InvalidInt   int // docs that could not be checked (missing id, etc.)
	SQLonlyInt   int // in SQL, but not in mongo
}

func (p *GFverifyStats) OK() bool {
	return p.ReadInt == p.TotalInt && p.DifferentInt == 0 && p.MissingInt == 0 && p.InvalidInt == 0
}

type GFverifyCallbacks struct {

	// called for every mongo record checked, oldest first.
	// pStatusStr is one of: "identical", "different", "missing", "invalid".
	// pSQL is only set for "different".
	OnRecord func(pIndexInt int, pTotalInt int, pStatusStr string, pMongo *GFbookmark, pSQL *GFbookmark)
}

func DBverifyMongoToSQL(pCallbacks *GFverifyCallbacks,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFverifyStats, *gf_core.GFerror) {

	if pCallbacks == nil {
		pCallbacks = &GFverifyCallbacks{}
	}

	collNameStr := "gf_bookmarks"
	coll        := pRuntimeSys.Mongo_db.Collection(collNameStr)

	stats := &GFverifyStats{}

	totalInt64, err := coll.CountDocuments(pCtx, bson.M{})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to count bookmarks in mongo for verification",
			"mongodb_find_error",
			map[string]interface{}{"coll_name_str": collNameStr},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	stats.TotalInt = int(totalInt64)

	cursor, err := coll.Find(pCtx,
		bson.M{},
		options.Find().SetSort(bson.D{
			{Key: "creation_unix_time_f", Value: 1},
			{Key: "_id", Value: 1},
		}))
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to query all bookmarks from mongo for verification",
			"mongodb_find_error",
			map[string]interface{}{"coll_name_str": collNameStr},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	defer cursor.Close(pCtx)

	reportFun := func(pStatusStr string, pMongo *GFbookmark, pSQL *GFbookmark) {
		if pCallbacks.OnRecord != nil {
			pCallbacks.OnRecord(stats.ReadInt, stats.TotalInt, pStatusStr, pMongo, pSQL)
		}
	}

	mongoIDsMap := map[gf_core.GF_ID]bool{}

	for cursor.Next(pCtx) {

		stats.ReadInt++

		var bookmark GFbookmark
		err := cursor.Decode(&bookmark)
		if err != nil || bookmark.Id_str == "" {
			stats.InvalidInt++
			reportFun("invalid", &bookmark, nil)
			continue
		}
		mongoIDsMap[bookmark.Id_str] = true

		existing, gfErr := dbSQLgetBookmarkByIDanyUser(bookmark.Id_str, pCtx, pRuntimeSys)
		if gfErr != nil {
			return stats, gfErr
		}

		if existing == nil {
			stats.MissingInt++
			reportFun("missing", &bookmark, nil)
			continue
		}

		if len(MigrateDiffFields(&bookmark, existing)) > 0 {
			stats.DifferentInt++
			reportFun("different", &bookmark, existing)
			continue
		}

		stats.IdenticalInt++
		reportFun("identical", &bookmark, nil)
	}

	if err := cursor.Err(); err != nil {
		gfErr := gf_core.MongoHandleError("mongo cursor error while verifying bookmarks",
			"mongodb_cursor_error",
			map[string]interface{}{"coll_name_str": collNameStr},
			err, "gf_tagger_lib", pRuntimeSys)
		return stats, gfErr
	}

	//------------------
	// SQL_ONLY - SQL records that are not in mongo
	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, `SELECT id FROM gf_bookmarks`)
	if err != nil {
		return stats, gf_core.ErrorCreate("failed to query bookmark ids from the SQL DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_tagger_lib", pRuntimeSys)
	}
	defer rows.Close()

	for rows.Next() {
		var idStr gf_core.GF_ID
		err := rows.Scan(&idStr)
		if err != nil {
			return stats, gf_core.ErrorCreate("failed to scan bookmark id from the SQL DB",
				"sql_row_scan",
				map[string]interface{}{},
				err, "gf_tagger_lib", pRuntimeSys)
		}
		if !mongoIDsMap[idStr] {
			stats.SQLonlyInt++
		}
	}
	if err := rows.Err(); err != nil {
		return stats, gf_core.ErrorCreate("rows iteration error for bookmark ids query",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_tagger_lib", pRuntimeSys)
	}

	return stats, nil
}
