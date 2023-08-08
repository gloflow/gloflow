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
	"context"
	"strings"
	"github.com/lib/pq"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// CREATE_TABLES

func DBsqlCreateTables(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_policy (
		id            TEXT,
		deleted       BOOLEAN DEFAULT FALSE,
		creation_time TIMESTAMP DEFAULT NOW(),
		
		target_resource_ids  TEXT[] NOT NULL,
		target_resource_type TEXT NOT NULL,
		owner_user_id        TEXT NOT NULL,
	
		public_view BOOLEAN DEFAULT FALSE,
	
		viewers_user_ids TEXT[],
		taggers_user_ids TEXT[],
		editors_user_ids TEXT[],

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
			editors_user_ids
		FROM 
			gf_policy 
		WHERE 
			id = $1 AND deleted=false;`

	policy := &GFpolicy{}
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pPolicyID).Scan(
		&policy.ID,
		&policy.DeletedBool,
		&policy.CreationUNIXtimeF,
		&policy.TargetResourceIDsLst,
		&policy.TargetResourceTypeStr,
		&policy.OwnerUserID,
		&policy.PublicViewBool,
		&policy.ViewersUserIDsLst,
		&policy.TaggersUserIDsLst,
		&policy.EditorsUserIDsLst)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get policy with ID in the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_policy", pRuntimeSys)
		return nil, gfErr
	}

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
			editors_user_ids
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
		err := rows.Scan(
			&policy.ID,
			&policy.DeletedBool,
			&policy.CreationUNIXtimeF,
			&policy.TargetResourceIDsLst,
			&policy.TargetResourceTypeStr,
			&policy.OwnerUserID,
			&policy.PublicViewBool,
			&policy.ViewersUserIDsLst,
			&policy.TaggersUserIDsLst,
			&policy.EditorsUserIDsLst)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to scan row to get policy with target_resource in the DB",
				"sql_query_execute",
				map[string]interface{}{},
				err, "gf_policy", pRuntimeSys)
			return nil, gfErr
		}
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
			id,
			target_resource_ids,
			target_resource_type,
			owner_user_id,
			public_view,
			viewers_user_ids,
			taggers_user_ids,
			editors_user_ids
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err := pRuntimeSys.SQLdb.ExecContext(
		pCtx,
		sqlStr,
		pPolicy.ID,
		pq.Array(pPolicy.TargetResourceIDsLst),
		pPolicy.TargetResourceTypeStr,
		pPolicy.OwnerUserID,
		pPolicy.PublicViewBool,
		pq.Array(pPolicy.ViewersUserIDsLst),
		pq.Array(pPolicy.TaggersUserIDsLst),
		pq.Array(pPolicy.EditorsUserIDsLst),
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
// UPDATE

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

/*
func DBsqlUpdatePolicy(pPolicyIDstr gf_core.GF_ID,
	pUpdateOp   *GFpolicyUpdateOp,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	fieldsTargets := bson.M{}

	if pUpdateOp.PublicViewBool != nil {
		fieldsTargets["public_view_bool"] = *pUpdateOp.PublicViewBool
	}


	_, err := pRuntimeSys.Mongo_db.Collection("gf_policies").UpdateMany(pCtx, bson.M{
		"id_str":       pPolicyIDstr,
		"deleted_bool": false,
	},
	bson.M{"$set": fieldsTargets})
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to to update policy in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"policy_id_str": string(pPolicyIDstr),
			},
			err, "gf_policy", pRuntimeSys)
		return gfErr
	}

	return nil
}
*/

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