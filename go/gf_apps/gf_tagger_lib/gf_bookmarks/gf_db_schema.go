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
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// CREATE_TABLES

func DBsqlCreateTables(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_bookmarks (
		v                                  VARCHAR(255),
		id                                 TEXT,
		deleted                            BOOLEAN DEFAULT FALSE,
		creation_time                      TIMESTAMP DEFAULT NOW(),
		user_id                            TEXT NOT NULL,

		url                                TEXT NOT NULL,
		description                        TEXT,
		tags                               TEXT[], -- array of tag strings

		-- SCREENSHOT --
		screenshot_image_id                TEXT,
		screenshot_image_thumbnail_url     TEXT,

		PRIMARY KEY(id)
	);

	CREATE INDEX IF NOT EXISTS idx_gf_bookmarks_user_id ON gf_bookmarks(user_id);
	CREATE INDEX IF NOT EXISTS idx_gf_bookmarks_creation_time ON gf_bookmarks(creation_time DESC);
	CREATE INDEX IF NOT EXISTS idx_gf_bookmarks_deleted ON gf_bookmarks(deleted);
	`

	_, err := pRuntimeSys.SQLdb.Exec(sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create bookmarks related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}
