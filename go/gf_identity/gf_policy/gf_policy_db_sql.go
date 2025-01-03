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

package gf_policy

import (
	"fmt"
	"time"
	"context"
	"strings"
	"github.com/lib/pq"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// UPDATE
//---------------------------------------------------

func DBsqlUpdatePolicyWithNewTargetFlow(pPolicyID gf_core.GF_ID,
	pFlowsIDsLst []gf_core.GF_ID,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	// array_cat - combines the existing target_resource_ids with the new pFlowsIDsLst
	// unnest    - flatten the combined array into individual elements
	// DISTINCT  - ensures that duplicates are removed after combining the arrays
	// array_agg - reconstructs the filtered elements back into an array
	sqlStr := `
		UPDATE gf_policy
		SET target_resource_ids = (
			SELECT DISTINCT unnest(array_cat(target_resource_ids, $1))
		)
		WHERE id = $2 
			AND deleted = false
			AND target_resource_type = 'flow';
	`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, pq.Array(pFlowsIDsLst), pPolicyID)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update policy in DB with new flow ID",
			"sql_query_execute",
			map[string]interface{}{
				"policy_id_str": pPolicyID,
				"flows_ids_lst": pFlowsIDsLst,
			},
			err, "gf_policy", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func DBsqlUpdatePolicy(pPolicyID gf_core.GF_ID,
	pUpdateOp   *GFpolicyUpdateOp, 
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	updates := []string{}
	args := []interface{}{}

	index := 1 // SQL parameter numbering

	// public_view update
	if pUpdateOp.PublicViewBool != nil {
		updates = append(updates, fmt.Sprintf("public_view = $%d", index))
		args = append(args, *pUpdateOp.PublicViewBool)
		index++
	}

	if len(updates) == 0 {

		// no fields to update
		return nil
	}

	sqlStr := fmt.Sprintf(`
		UPDATE gf_policy
		SET %s
		WHERE
			id=%s AND deleted=false`,
		
			strings.Join(updates, ", "), pPolicyID)

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, args...)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update policy in DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_policy", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// VAR
//---------------------------------------------------

func DBsqlGetFlowPolicyIDforUser(pUserID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	sqlStr := `
		SELECT id
		FROM gf_policy
		WHERE
			owner_user_id = $1 AND
			target_resource_type = 'flow' AND
			deleted = false
		LIMIT 1;
	`
	
	var idStr string
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserID)).Scan(
		&idStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get flow policy ID for user in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id": string(pUserID),
			},
			err, "gf_policy", pRuntimeSys)
		return gf_core.GF_ID(""), gfErr
	}

	return gf_core.GF_ID(idStr), nil
}

//---------------------------------------------------
// CREATE_TABLES

func DBsqlCreateTables(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_policy (
		v             VARCHAR(255),
		id            TEXT,
		deleted       BOOLEAN DEFAULT FALSE,
		creation_time TIMESTAMP DEFAULT NOW(),
		
		target_resource_ids  TEXT[] NOT NULL,
		target_resource_type TEXT NOT NULL,
		owner_user_id        TEXT NOT NULL,
	
		public_view BOOLEAN DEFAULT FALSE,
	
		viewers_user_ids TEXT[] NOT NULL,
		taggers_user_ids TEXT[] NOT NULL,
		editors_user_ids TEXT[] NOT NULL,
		admins_user_ids  TEXT[] NOT NULL,

		PRIMARY KEY(id)
	);
	`

	_, err := pRuntimeSys.SQLdb.Exec(sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create policies related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_policy", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// GET_BY_ID

func DBsqlGetPolicyByID(pPolicyID gf_core.GF_ID,
	pCtx        context.Context, 
	pRuntimeSys *gf_core.RuntimeSys) (*GFpolicy, *gf_core.GFerror) { 

	const sqlStr = `
		SELECT 
			id,
			deleted,
			creation_time,
			target_resource_ids,
			target_resource_type,
			owner_user_id,
			public_view,
			viewers_user_ids,
			taggers_user_ids,
			editors_user_ids,
			admins_user_ids
		FROM 
			gf_policy 
		WHERE 
			id = $1 AND deleted=false;`

	policy := &GFpolicy{}

	var creationTime time.Time
	var idsStr string
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pPolicyID).Scan(
		&policy.ID,
		&policy.DeletedBool,
		&creationTime, // policy.CreationUNIXtimeF,
		&idsStr, // &policy.TargetResourceIDsLst,
		&policy.TargetResourceTypeStr,
		&policy.OwnerUserID,
		&policy.PublicViewBool,
		&policy.ViewersUserIDsLst,
		&policy.TaggersUserIDsLst,
		&policy.EditorsUserIDsLst,
		&policy.AdminsUserIDsLst,)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get policy with ID in the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_policy", pRuntimeSys)
		return nil, gfErr
	}



	creationUNIXtimeF := float64(creationTime.UnixNano()) / 1e9
	policy.CreationUNIXtimeF = creationUNIXtimeF


	return policy, nil
}

//---------------------------------------------------
// GET_POLICIES

func DBsqlGetPolicies(pTargetResourceID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFpolicy, *gf_core.GFerror) {

	const sqlStr = `
		SELECT 
			id,
			deleted,
			creation_time,
			target_resource_ids,
			target_resource_type,
			owner_user_id,
			public_view,
			viewers_user_ids,
			taggers_user_ids,
			editors_user_ids,
			admins_user_ids
		FROM 
			gf_policy 
		WHERE 
			$1 = ANY(target_resource_ids) AND deleted=false;`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr, pTargetResourceID)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to execute query to get policy with target_resource in the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_policy", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	policiesLst := []*GFpolicy{}
	for rows.Next() {

		policy := &GFpolicy{}

		var creationTime time.Time
		// var targetResourceIDsLst []string
		var targetResourceIDsStr string
		var viewersUserIDsStr    *string
		var taggersUserIDsStr    *string
		var editorsUserIDsStr    *string
		var adminsUserIDsStr     *string

		err := rows.Scan(
			&policy.ID,
			&policy.DeletedBool,
			&creationTime,         // policy.CreationUNIXtimeF,
			&targetResourceIDsStr, // &targetResourceIDsLst, // policy.TargetResourceIDsLst,
			&policy.TargetResourceTypeStr,
			&policy.OwnerUserID,
			&policy.PublicViewBool,
			&viewersUserIDsStr,
			&taggersUserIDsStr,
			&editorsUserIDsStr,
			&adminsUserIDsStr)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to scan row to get policy with target_resource in the DB",
				"sql_query_execute",
				map[string]interface{}{},
				err, "gf_policy", pRuntimeSys)
			return nil, gfErr
		}

		creationUNIXtimeF := float64(creationTime.UnixNano()) / 1e9
		policy.CreationUNIXtimeF = creationUNIXtimeF

		//-----------------------
		// PARSE_ARRAYS

		targetResourceIDsLst       := strings.Split(strings.Trim(targetResourceIDsStr, "{}"), ",")
		policy.TargetResourceIDsLst = targetResourceIDsLst

		// check value is not null or empty string
		if viewersUserIDsStr != nil && *viewersUserIDsStr != "" {
			policy.ViewersUserIDsLst = strings.Split(strings.Trim(*viewersUserIDsStr, "{}"), ",")
		}

		if taggersUserIDsStr != nil && *taggersUserIDsStr != "" {
			policy.TaggersUserIDsLst = strings.Split(strings.Trim(*taggersUserIDsStr, "{}"), ",")
		}
		
		if editorsUserIDsStr != nil && *editorsUserIDsStr != "" {
			policy.EditorsUserIDsLst = strings.Split(strings.Trim(*editorsUserIDsStr, "{}"), ",")
		}

		if adminsUserIDsStr != nil && *adminsUserIDsStr != "" {
			policy.AdminsUserIDsLst = strings.Split(strings.Trim(*adminsUserIDsStr, "{}"), ",")
		}

		//-----------------------

		policiesLst = append(policiesLst, policy)
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("failed to scan row to get policy with target_resource in the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_policy", pRuntimeSys)
		return nil, gfErr
	}

	return policiesLst, nil
}

//---------------------------------------------------
// CREATE

func DBsqlCreatePolicy(pPolicy *GFpolicy,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_policy (
			v,
			id,
			target_resource_ids,
			target_resource_type,
			owner_user_id,
			public_view,
			viewers_user_ids,
			taggers_user_ids,
			editors_user_ids,
			admins_user_ids
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	_, err := pRuntimeSys.SQLdb.ExecContext(
		pCtx,
		sqlStr,
		"0", // version
		pPolicy.ID,
		pq.Array(pPolicy.TargetResourceIDsLst),
		pPolicy.TargetResourceTypeStr,
		pPolicy.OwnerUserID,
		pPolicy.PublicViewBool,
		pq.Array(pPolicy.ViewersUserIDsLst),
		pq.Array(pPolicy.TaggersUserIDsLst),
		pq.Array(pPolicy.EditorsUserIDsLst),
		pq.Array(pPolicy.AdminsUserIDsLst),
	)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert policy into the DB",
			"sql_row_insert",
			map[string]interface{}{
				"policy": pPolicy,
			},
			err, "gf_policy", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// EXISTS_BY_USERNAME

func DBsqlExistsByID(pPolicyID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT COUNT(*)
		FROM gf_policies
		WHERE id = $1 AND deleted = false
	`

	var countInt int
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pPolicyID).Scan(&countInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if there is a policy in the DB with a given ID",
			"sql_query_execute",
			map[string]interface{}{
				"policy_id_str": pPolicyID,
			},
			err, "gf_policy", pRuntimeSys)
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}