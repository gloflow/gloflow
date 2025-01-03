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

package gf_images_flows

import (
	"context"
	"encoding/json"
	"database/sql"
	"github.com/lib/pq"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------

func dbSQLgetPage(pFlowNameStr string,
	pCursorStartPositionInt int, // 0
	pElementsNumInt int, // 50
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*gf_images_core.GFimage, *gf_core.GFerror) {

	query := `
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
		WHERE NOT deleted
		AND (
			$1 = ANY(flows_names)
		)
		ORDER BY creation_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, query, pFlowNameStr, pElementsNumInt, pCursorStartPositionInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get a page of images from a flow",
			"sql_query_execution",
			map[string]interface{}{
				"flow_name_str":             pFlowNameStr,
				"cursor_start_position_int": pCursorStartPositionInt,
				"elements_num_int":          pElementsNumInt,
			},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	imagesLst := []*gf_images_core.GFimage{}
	for rows.Next() {
		img := &gf_images_core.GFimage{}

		var DominantColorHexStr sql.NullString
		var PalleteStr sql.NullString
		var metaMapRaw []byte

		if err := rows.Scan(
				&img.IDstr,
				&img.Creation_unix_time_f,
				&img.UserID,

				&img.ClientTypeStr,
				&img.TitleStr,
				pq.Array(&img.FlowsNamesLst),

				&img.Origin_url_str,
				&img.Origin_page_url_str,

				&img.ThumbnailSmallURLstr,
				&img.ThumbnailMediumURLstr,
				&img.ThumbnailLargeURLstr,

				&img.Format_str,
				&img.Width_int,
				&img.Height_int,

				&DominantColorHexStr,
				&PalleteStr,
				&metaMapRaw,
				pq.Array(&img.TagsLst)); err != nil {
					
			gfErr := gf_core.ErrorCreate("failed to scan a row of images",
				"sql_row_scan",
				map[string]interface{}{},
				err, "gf_images_flows", pRuntimeSys)
			return nil, gfErr
		}

		// DOMINANT_COLOR_HEX
		if DominantColorHexStr.Valid {
			img.DominantColorHexStr = DominantColorHexStr.String
		} else {
			img.DominantColorHexStr = "" // Default value for NULL
		}

		// PALLETE
		if PalleteStr.Valid {
			img.PalleteStr = PalleteStr.String
		} else {
			img.PalleteStr = "" // Default value for NULL
		}

		// META_MAP
		if err := json.Unmarshal(metaMapRaw, &img.MetaMap); err != nil {
			gfErr := gf_core.ErrorCreate("failed to unmarshal JSON meta_map",
				"json_decode_error",
				map[string]interface{}{},
				err, "gf_images_flows", pRuntimeSys)
			return nil, gfErr
		}

		imagesLst = append(imagesLst, img)
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("failed to iterate rows",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}

	return imagesLst, nil
}

//---------------------------------------------------

func DBgetFlowByName(pFlowNameStr string,
	pCtx		context.Context,
	pRuntimeSys	*gf_core.RuntimeSys) (*GFflow, *gf_core.GFerror) {

	sqlStr := `
		SELECT
			v,
			id,
			EXTRACT(EPOCH FROM creation_time) AS creation_unix_time,
			name,
			creator_user_id,
			public,
			description
		FROM gf_images_flows
		WHERE name = $1 AND deleted = FALSE
		LIMIT 1;
	`
	var flow GFflow
	var v sql.NullString
	var public sql.NullBool
	var description sql.NullString

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pFlowNameStr).Scan(
		&v,
		&flow.IDstr,
		&flow.CreationUNIXtimeF,
		&flow.NameStr,
		&flow.OwnerUserID,
		&public,
		&description)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get flow from SQL DB...",
			"sql_query_execute",
			map[string]interface{}{
				"flow_name_str": pFlowNameStr,
			},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}

	if v.Valid {
		flow.Vstr = v.String
	} else {
		flow.Vstr = "" // Default value for NULL
	}

	if public.Valid {
		flow.PublicBool = public.Bool
	} else {
		flow.PublicBool = false // Default value for NULL
	}

	if description.Valid {
		flow.DescriptionStr = description.String
	} else {
		flow.DescriptionStr = "" // Default value for NULL
	}

	return &flow, nil
}

//---------------------------------------------------

func DBsqlGetAll(pCtx context.Context, pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	sqlStr := `
		WITH UnwoundFlows AS (
			SELECT
				UNNEST(flows_names) AS flow_name
			FROM
				gf_images
		),
		FlowCounts AS (
			SELECT
				flow_name AS _id,
				COUNT(*) AS count_int
			FROM
				UnwoundFlows
			GROUP BY
				flow_name
		)
		SELECT
			_id AS flow_name,
			count_int
		FROM
			FlowCounts
		ORDER BY
			count_int DESC;
	`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to execute SQL query to get all flow names",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	var resultsLst []map[string]interface{}
	for rows.Next() {
		
		var flowNameStr string
		var countInt int

		if err := rows.Scan(&flowNameStr, &countInt); err != nil {
			gfErr := gf_core.ErrorCreate("failed to scan row for flow names and counts",
				"sql_row_scan",
				map[string]interface{}{},
				err, "gf_images_flows", pRuntimeSys)
			return nil, gfErr
		}
		resultsLst = append(resultsLst, map[string]interface{}{
			"name_str":  flowNameStr,
			"count_int": countInt,
		})
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("error encountered while iterating over query results",
			"sql_row_scan",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}

	return resultsLst, nil
}

//---------------------------------------------------
// GET_FLOWS_IDS

func DBsqlGetFlowsIDs(pFlowsNamesLst []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]gf_core.GF_ID, *gf_core.GFerror) {

	flowsIDsLst := []gf_core.GF_ID{}
	for _, flowNameStr := range pFlowsNamesLst {
		flowIDstr, gfErr := DBsqlGetID(flowNameStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
		flowsIDsLst = append(flowsIDsLst, flowIDstr)
	}
	return flowsIDsLst, nil
}

//---------------------------------------------------
// GET_ID

func DBsqlGetID(pFlowNameStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	const sqlStr = `SELECT id FROM gf_images_flows WHERE name = $1 LIMIT 1`

	var flowIDstr string
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pFlowNameStr).Scan(&flowIDstr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if flow exists in SQL DB, might not exist...",
			"sql_query_execute",
			map[string]interface{}{
				"flow_name_str": pFlowNameStr,
			},
			err, "gf_images_flows", pRuntimeSys)
		return "", gfErr
	}

	return gf_core.GF_ID(flowIDstr), nil
}

//---------------------------------------------------
// CREATE_FLOW

func DBsqlCreateFlow(pFlowID gf_core.GF_ID,
	pFlowNameStr string,
	pOwnerUserID gf_core.GF_ID,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	db := pRuntimeSys.SQLdb

	tx, err := db.Begin()
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to begin the SQL transaction to create a flow",
			"sql_transaction_begin",
			map[string]interface{}{
				"flow_id_str":       string(pFlowID),
				"flow_name_str":     pFlowNameStr,
				"owner_user_id_str": pOwnerUserID,
			},
			err, "gf_images_flows", pRuntimeSys)
		return gfErr
	}

	// The rollback will be ignored commit is successful
	defer tx.Rollback()

	row := tx.QueryRow(`
		INSERT INTO gf_images_flows (
			id,
			name,
			creator_user_id
		)
		VALUES ($1, $2, $3) RETURNING id
		`,
		string(pFlowID),
		pFlowNameStr,
		pOwnerUserID)

	var id string
	err = row.Scan(&id)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create a new images flow in the DB",
			"sql_row_insert",
			map[string]interface{}{
				"flow_id_str":   string(pFlowID),
				"flow_name_str": pFlowNameStr,
				"user_id_str":   pOwnerUserID,
			},
			err, "gf_images_flows", pRuntimeSys)
		return gfErr
	}
	
	// TX_COMMIT
	err = tx.Commit()
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to commit the SQL transaction to create a flow",
			"sql_transaction_commit",
			map[string]interface{}{
				"flow_id_str":       string(pFlowID),
				"flow_name_str":     pFlowNameStr,
				"owner_user_id_str": pOwnerUserID,
			},
			err, "gf_images_flows", pRuntimeSys)
		return gfErr
	}
	
	return nil

	/*
	// EDITORS
	for _, editorID := range pFlow.EditorUserIDs {
		_, err := tx.Exec(
			"INSERT INTO gf_images_flows_editors (flow_id, user_id) VALUES ($1, $2)",
			id,
			editorID,
		)

		if err != nil {
			gfErr := gf_core.ErrorCreate("failed create a new images flow in the DB",
				"sql_row_insert",
				map[string]interface{}{
					"flow_name_str": pFlowNameStr,
					"user_id_str":   pOwnerUserID,
				},
				err, "gf_images_flows", pRuntimeSys)
			return gfErr
		}
	}
	*/
}

//---------------------------------------------------
// CREATE_TABLES

func DBsqlCreateTables(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_images_flows (
		v               VARCHAR(255),
		id              TEXT,
		deleted         BOOLEAN DEFAULT FALSE,
		creation_time   TIMESTAMP DEFAULT NOW(),
		name            TEXT NOT NULL,
		creator_user_id TEXT NOT NULL,
		public          BOOLEAN,
		description     TEXT,

		PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS gf_images_flows_editors (
		v       VARCHAR(255),
		flow_id TEXT REFERENCES gf_images_flows(id),
		user_id TEXT NOT NULL,
		
		PRIMARY KEY(flow_id, user_id)
	);
	`

	_, err := pRuntimeSys.SQLdb.Exec(sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create flow related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// CHECK_FLOW_EXISTS

func DBsqlCheckFlowExists(pFlowNameStr string,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	db := pRuntimeSys.SQLdb

	var existsBool bool
	sqlStr := `SELECT exists(SELECT 1 FROM gf_images_flows WHERE name=$1)`
	err := db.QueryRow(sqlStr, pFlowNameStr).Scan(&existsBool)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No record found, flow does not exist
		}
		
		gfErr := gf_core.ErrorCreate("failed to check if a flow exists in the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return false, gfErr
	}
	return existsBool, nil
}