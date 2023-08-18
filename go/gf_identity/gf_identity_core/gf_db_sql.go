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

package gf_identity_core

import (
	"fmt"
	"strings"
	"time"
	"context"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// AUTH0
//---------------------------------------------------

func dbSQLAuth0createNewSession(pAuth0session *GFauth0session,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_auth0_session (
			v,
			id,
			deleted,
			login_complete,
			access_token,
			profile
		)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	// serializing the profile map to JSON to store in the database
	profileMapJSON, err := json.Marshal(pAuth0session.ProfileMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to serialize profile_map to JSON for DB",
			"json_encode_error",
			nil,
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	_, err = pRuntimeSys.SQLdb.ExecContext(pCtx,
		sqlStr,
		"0",
		pAuth0session.ID,
		pAuth0session.DeletedBool,
		pAuth0session.LoginCompleteBool,
		pAuth0session.AccessTokenStr,
		profileMapJSON)
	
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert GFauth0session into the DB",
			"sql_query_execute",
			map[string]interface{}{
				"session_id_str": pAuth0session.ID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbSQLauth0getSession(pGFsessionID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFauth0session, *gf_core.GFerror) {

	sqlStr := `
		SELECT
			id,
			deleted,
			creation_time,
			login_complete,
			access_token,
			profile
		FROM gf_auth0_session
		WHERE id = $1 AND deleted = false`

	session := GFauth0session{}
	var creationTime time.Time
	var profileJSON []byte

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pGFsessionID).Scan(
		&session.ID,
		&session.DeletedBool,
		&creationTime,
		&session.LoginCompleteBool,
		&session.AccessTokenStr,
		&profileJSON)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to find Auth0 session by ID in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"session_id_str": pGFsessionID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	session.CreationUNIXtimeF = float64(creationTime.Unix())

	if err := json.Unmarshal(profileJSON, &session.ProfileMap); err != nil {
		gfErr := gf_core.ErrorCreate("failed to unmarshal profile JSON",
			"json_unmarshal_error",
			map[string]interface{}{
				"session_id_str": pGFsessionID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}
	
	return &session, nil
}

//---------------------------------------------------

func dbSQLauth0updateSession(pGFsessionID gf_core.GF_ID,
	pLoginCompleteBool bool,
	pAuth0profileMap   map[string]interface{},
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	profileMapJSONstr, err := json.Marshal(pAuth0profileMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to convert profile map to JSON",
			"json_conversion_error",
			map[string]interface{}{
				"session_id_str": pGFsessionID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	sqlStr := `
		UPDATE gf_auth0_session
		SET login_complete = $1, profile = $2
		WHERE id = $3 AND deleted = false`

	_, err = pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, pLoginCompleteBool, profileMapJSONstr, pGFsessionID)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update Auth0 session in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"session_id_str": pGFsessionID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// USER
//---------------------------------------------------
// USER_CREATE

func DBsqlUserCreate(pUser *GFuser,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// Create a prepared SQL statement to insert a user
	sqlStr := `
	INSERT INTO gf_users (
		v,
		id,
		user_type,
		user_name,
		screen_name,
		description,
		addresses_eth,
		email,
		email_confirmed,
		profile_image_url,
		banner_image_url)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	sqlStatement, err := pRuntimeSys.SQLdb.Prepare(sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to prepare insert statement for GF_user",
			"sql_prepare_statement",
			map[string]interface{}{
				"user_id_str":        pUser.ID,
				"user_name_str":      pUser.UserNameStr,
				"description_str":    pUser.DescriptionStr,
				"addresses_eth_lst":  pUser.AddressesETHlst,
				"caller_err_msg_str": "failed to insert GF_user into the SQL DB",
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}
	defer sqlStatement.Close()

	_, err = sqlStatement.ExecContext(
		pCtx,
		pUser.Vstr,
		pUser.ID,
		pUser.UserTypeStr,
		pUser.UserNameStr,
		pUser.ScreenNameStr,
		pUser.DescriptionStr,
		pq.Array(pUser.AddressesETHlst),
		pUser.EmailStr,
		pUser.EmailConfirmedBool,
		pUser.ProfileImageURLstr,
		pUser.BannerImageURLstr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert user into the sql DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str":        pUser.ID,
				"user_name_str":      pUser.UserNameStr,
				"description_str":    pUser.DescriptionStr,
				"addresses_eth_lst":  pUser.AddressesETHlst,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func DBsqlUserGetAll(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFuser, *gf_core.GFerror) {

	sqlStr := `SELECT
			v,
			id,
			deleted,
			creation_time,
			user_type,
			user_name,
			screen_name,
			description,
			addresses_eth,
			email,
			email_confirmed,
			profile_image_url,
			banner_image_url
		FROM gf_users
		WHERE deleted = FALSE`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get all users records from the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_identity", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	usersLst := []*GFuser{}
	for rows.Next() {

		var user GFuser
		var addressesEth pq.StringArray	
		var creationTime time.Time

		err := rows.Scan(
			&user.Vstr,
			&user.ID,
			&user.DeletedBool,
			&creationTime,
			&user.UserTypeStr,
			&user.UserNameStr,
			&user.ScreenNameStr,
			&user.DescriptionStr,
			&addressesEth,
			&user.EmailStr,
			&user.EmailConfirmedBool,
			&user.ProfileImageURLstr,
			&user.BannerImageURLstr)
		
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to get all users records from rows",
				"sql_row_scan",
				map[string]interface{}{},
				err, "gf_identity_core", pRuntimeSys)
			return nil, gfErr
		}

		for _, addr := range addressesEth {
			user.AddressesETHlst = append(user.AddressesETHlst, GFuserAddressETH(addr))
		}
		user.CreationUNIXtimeF = float64(creationTime.Unix())

		usersLst = append(usersLst, &user)
	}

	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("failed to iterate through rows to get all users from DB",
			"sql_rows_iteration",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	return usersLst, nil
}

//---------------------------------------------------
// GET_BASIC_INFO_BY_ETH_ADDR

func DBsqlGetBasicInfoByETHaddr(pUserAddressETH GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	sqlStr := `
	SELECT id
	FROM gf_users
	WHERE $1 = ANY(addresses_eth) AND deleted = false
	LIMIT 1
	`

	var userIDstr string
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pUserAddressETH).
		Scan(&userIDstr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get user basic_info in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_address_eth_str": pUserAddressETH,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gf_core.GF_ID(""), gfErr
	}

	return gf_core.GF_ID(userIDstr), nil
}

//---------------------------------------------------
// USER_GET_BY_ID

func DBsqlUserGetByID(pUserID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuser, *gf_core.GFerror) {
	
	sqlStr := `
	SELECT 
		v,
		id,
		deleted,
		creation_time,
		user_type,
		user_name,
		screen_name,
		description,
		addresses_eth,
		email,
		email_confirmed,
		profile_image_url,
		banner_image_url
	FROM gf_users
	WHERE id = $1 AND deleted = false
	LIMIT 1
	`

	row := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pUserID)
	user := &GFuser{}

	var addressesEth pq.StringArray	
	var creationTime time.Time

	err := row.Scan(
		&user.Vstr,
		&user.ID,
		&user.DeletedBool,
		&creationTime,
		&user.UserTypeStr,
		&user.UserNameStr,
		&user.ScreenNameStr,
		&user.DescriptionStr,
		&addressesEth,
		&user.EmailStr,
		&user.EmailConfirmedBool,
		&user.ProfileImageURLstr,
		&user.BannerImageURLstr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to find user by ID in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str": pUserID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	for _, addr := range addressesEth {
		user.AddressesETHlst = append(user.AddressesETHlst, GFuserAddressETH(addr))
	}
	user.CreationUNIXtimeF = float64(creationTime.Unix())

	return user, nil
}

//---------------------------------------------------
// GET_USER_NAME_BY_ID

func DBsqlGetUserNameByID(pUserID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (GFuserName, *gf_core.GFerror) {

	sqlStr := `
		SELECT user_name
		FROM gf_users
		WHERE id = $1 AND deleted = false
		LIMIT 1
	`
	row := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserID))

	var userNameStr GFuserName

	err := row.Scan(&userNameStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get user name in the SQL DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str": pUserID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return GFuserName(""), gfErr
	}

	return userNameStr, nil
}

//---------------------------------------------------
// USER_EXISTS_BY_ID

func DBsqlUserExistsByID(pUserID gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT COUNT(*)
		FROM gf_users
		WHERE id = $1 AND deleted = false
	`

	var countInt int

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserID)).Scan(&countInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if there is a user in the DB with a given ID",
			"sql_query_execute",
			map[string]interface{}{
				"id_str": pUserID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------

func DBsqlUserExistsByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT COUNT(*)
		FROM gf_users
		WHERE user_name = $1 AND deleted = false
	`

	var countInt int

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserNameStr)).Scan(&countInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if there is a user in the DB with a given user_name",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str":  pUserNameStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------

func dbSQLuserExistsByETHaddr(pUserAddressETH GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT COUNT(*)
		FROM gf_users
		WHERE $1 = ANY(addresses_eth) AND deleted = false
	`

	var countInt int

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserAddressETH)).Scan(&countInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if there is a user in the DB with a given ETH address",
			"sql_query_execute",
			map[string]interface{}{
				"user_address_eth_str":  pUserAddressETH,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------

func DBsqlGetBasicInfoByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	sqlStr := `
		SELECT id
		FROM gf_users
		WHERE user_name = $1 AND deleted = false
		LIMIT 1
	`

	var userIDstr gf_core.GF_ID

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserNameStr)).Scan(&userIDstr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get user ID by username in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gf_core.GF_ID(""), gfErr
	}

	return userIDstr, nil
}

//---------------------------------------------------

func dbSQLuserGetByETHaddr(pUserAddressETH GFuserAddressETH,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFuser, *gf_core.GFerror) {

	sqlStr := `
		SELECT 
			v,
			id,
			deleted,
			creation_time,
			user_type,
			user_name,
			screen_name,
			description,
			addresses_eth,
			email,
			email_confirmed,
			profile_image_url,
			banner_image_url
		FROM gf_users
		WHERE $1 = ANY(addresses_eth) AND deleted = false`

	user := GFuser{}

	var addressesEth pq.StringArray	
	var creationTime time.Time

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pUserAddressETH).Scan(
		&user.Vstr,
		&user.ID,
		&user.DeletedBool,
		&creationTime,
		&user.UserTypeStr,
		&user.UserNameStr,
		&user.ScreenNameStr,
		&user.DescriptionStr,
		&addressesEth,
		&user.EmailStr,
		&user.EmailConfirmedBool,
		&user.ProfileImageURLstr,
		&user.BannerImageURLstr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to find user by Eth address in the DB",
			"sql_find_error",
			map[string]interface{}{
				"user_address_eth_str": pUserAddressETH,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	for _, addr := range addressesEth {
		user.AddressesETHlst = append(user.AddressesETHlst, GFuserAddressETH(addr))
	}

	user.CreationUNIXtimeF = float64(creationTime.Unix())

	return &user, nil
}

//---------------------------------------------------

func DBsqlUserUpdate(pUserIDstr gf_core.GF_ID,
	pUpdateOp   *GFuserUpdateOp,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	updateFields := []string{}
	paramsLst    := []interface{}{}
	
	i:=1

	if pUpdateOp.DeletedBool != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'deleted' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("deleted = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.DeletedBool)
		i+=1
	}

	if pUpdateOp.UserNameStr != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'user_name' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("user_name = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.UserNameStr)
		i+=1
	}

	if pUpdateOp.ScreenNameStr != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'screen_name' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("screen_name = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.ScreenNameStr)
		i+=1
	}

	if pUpdateOp.DescriptionStr != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'description' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("description = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.DescriptionStr)
		i+=1
	}

	if pUpdateOp.EmailStr != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'email' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("email = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.EmailStr)
		i+=1
	}

	if pUpdateOp.EmailConfirmedBool != nil && *pUpdateOp.EmailConfirmedBool {
		pRuntimeSys.LogNewFun("DEBUG", "user 'email_confirmed' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("email_confirmed = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.EmailConfirmedBool)
		i+=1
	}
	
	if pUpdateOp.MFAconfirmBool != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'mfa_confirm' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("mfa_confirm = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.MFAconfirmBool)
		i+=1
	}

	if pUpdateOp.ProfileImageURLstr != nil {
		pRuntimeSys.LogNewFun("DEBUG", "user 'profile_image_url' column to be updated...", nil)
		updateFields = append(updateFields, fmt.Sprintf("profile_image_url = $%d", i))
		paramsLst = append(paramsLst, *pUpdateOp.ProfileImageURLstr)
		i+=1
	}

	sqlStr := fmt.Sprintf("UPDATE gf_users SET %s WHERE id = $%d AND deleted = false",
		strings.Join(updateFields, ","), i)
	paramsLst = append(paramsLst, pUserIDstr)

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, paramsLst...)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update user info",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str": pUserIDstr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// EMAIL
//---------------------------------------------------

func DBsqlUserEmailIsConfirmed(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT email_confirmed
		FROM gf_users
		WHERE user_name = $1 AND deleted = false
		LIMIT 1
	`

	var emailConfirmedBool bool

	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserNameStr)).Scan(&emailConfirmedBool)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get user email_confirmed from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, gfErr
	}

	return emailConfirmedBool, nil
}

//---------------------------------------------------

func dbSQLuserGetEmailConfirmedByUsername(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT email_confirmed
		FROM gf_users
		WHERE user_name = $1`

	var emailConfirmedBool bool
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pUserNameStr).Scan(&emailConfirmedBool)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get user email_confirm status of a user from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str": string(pUserNameStr),
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, gfErr
	}

	return emailConfirmedBool, nil
}

//---------------------------------------------------
// INVITE_LIST
//---------------------------------------------------

func DBsqlUserGetAllInInviteList(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	sqlStr := `SELECT user_email, creation_time
		FROM gf_users_invite_list
		WHERE deleted = FALSE`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get all records in invite_list from the DB",
			"sql_query_execute",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	inviteListLst := []map[string]interface{}{}
	for rows.Next() {
		
		var userEmail string
		var creationTime time.Time

		if err := rows.Scan(&userEmail, &creationTime); err != nil {
			gfErr := gf_core.ErrorCreate("failed to get all records in invite_list from rows",
				"sql_row_scan",
				map[string]interface{}{},
				err, "gf_identity_core", pRuntimeSys)
			return nil, gfErr
		}
		
		inviteRecord := map[string]interface{}{
			"user_email_str":       userEmail,
			"creation_unix_time_f": float64(creationTime.Unix()),
		}
		inviteListLst = append(inviteListLst, inviteRecord)
	}
	
	if err := rows.Err(); err != nil {
		gfErr := gf_core.ErrorCreate("failed to iterate through rows to get all users in invite list",
			"sql_row_scan",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	return inviteListLst, nil
}

//---------------------------------------------------

func DBsqlUserAddToInviteList(pUserEmailStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_users_invite_list
		(user_email)
		VALUES ($1)
	`
	
	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, pUserEmailStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to add a user to the invite_list in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_email_str": pUserEmailStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------

func DBsqlUserRemoveFromInviteList(pUserEmailStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	sqlStr := `
		UPDATE gf_users_invite_list
		SET deleted = TRUE
		WHERE user_email = $1 AND deleted = FALSE
	`
	
	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, pUserEmailStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to remove user from invite list",
			"sql_query_execute",
			map[string]interface{}{
				"user_email_str": pUserEmailStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}
	return nil
}

//---------------------------------------------------
// CHECK_IN_INVITE_LIST_BY_EMAIL

func dbSQLuserCheckInInvitelistByEmail(pUserEmailStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT COUNT(*) FROM gf_users_invite_list
		WHERE user_email = $1 AND deleted = false`

	var countInt int
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pUserEmailStr).Scan(&countInt)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to check if the user_name is in the invite list",
			"sql_query_execute",
			map[string]interface{}{
				"user_email_str": pUserEmailStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}

//---------------------------------------------------
// USER_CREDS
//---------------------------------------------------

func dbSQLuserCredsCreate(pUserCreds *GFuserCreds,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_users_creds (
			v,
			id,
			user_id,
			user_name,
			pass_salt,
			pass_hash
		)
		VALUES ($1, $2, $3, $4, $5, $6)`
	
	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr,
		pUserCreds.Vstr,
		pUserCreds.ID,
		pUserCreds.UserID,
		pUserCreds.UserNameStr,
		pUserCreds.PassSaltStr,
		pUserCreds.PassHashStr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert gf_user_creds into the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str":   pUserCreds.UserID,
				"user_name_str": pUserCreds.UserNameStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------

func dbSQLuserCredsGetPassHash(pUserNameStr GFuserName,
	pCtx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, string, *gf_core.GFerror) {

	sqlStr := `
		SELECT pass_salt, pass_hash FROM gf_users_creds
		WHERE user_name = $1 AND deleted = false`

	var passSaltStr, passHashStr string
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, string(pUserNameStr)).Scan(
		&passSaltStr,
		&passHashStr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to find user creds by user_name in the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str": pUserNameStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return "", "", gfErr
	}

	return passSaltStr, passHashStr, nil
}

//---------------------------------------------------
// EMAIL
//---------------------------------------------------
// CREATE__EMAIL_CONFIRM

func dbSQLuserEmailConfirmCreate(pUserNameStr GFuserName,
	pUserID         gf_core.GF_ID,
	pConfirmCodeStr string,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
		INSERT INTO gf_users_email_confirm (
			user_name,
			user_id,
			confirm_code
		)
		VALUES ($1, $2, $3)`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, pUserNameStr, pUserID, pConfirmCodeStr)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert user email confirm_code into the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_id_str": pUserID,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}
	
	return nil
}

//---------------------------------------------------
// GET__EMAIL_CONFIRM_CODE

func dbSQLuserEmailConfirmGetCode(pUserNameStr GFuserName,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, float64, *gf_core.GFerror) {

	sqlStr := `
		SELECT confirm_code, creation_time
		FROM gf_users_email_confirm
		WHERE user_name = $1
		ORDER BY creation_time DESC
		LIMIT 1`

	var confirmCodeStr string
	var creationTime time.Time
	err := pRuntimeSys.SQLdb.QueryRowContext(pCtx, sqlStr, pUserNameStr).
		Scan(&confirmCodeStr, &creationTime)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get user email_confirm info from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str": string(pUserNameStr),
			},
			err, "gf_identity_core", pRuntimeSys)
		return "", 0.0, gfErr
	}

	creationUNIXtimeF := float64(creationTime.Unix())

	return confirmCodeStr, creationUNIXtimeF, nil
}

//---------------------------------------------------
// LOGIN_ATTEMPT
//---------------------------------------------------

func dbSQLloginAttemptCreate(pLoginAttempt *GFloginAttempt,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	sqlStr := `
		INSERT INTO gf_login_attempts (
			v,
			id,
			user_type,
			user_name,
			pass_confirmed,
			email_confirmed,
			mfa_confirmed
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr,
		pLoginAttempt.Vstr,
		pLoginAttempt.IDstr,
		pLoginAttempt.UserTypeStr,
		string(pLoginAttempt.UserNameStr),
		pLoginAttempt.PassConfirmedBool,
		pLoginAttempt.EmailConfirmedBool,
		pLoginAttempt.MFAconfirmedBool)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to insert login_attempt into the DB",
			"sql_query_execute",
			map[string]interface{}{
				"login_attempt_id_str": pLoginAttempt.IDstr,
				"user_type_str":        pLoginAttempt.UserTypeStr,
				"user_name_str":        pLoginAttempt.UserNameStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbSQLloginAttemptGetByUsername(pUserNameStr GFuserName,
	pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFloginAttempt, *gf_core.GFerror) {
	
	sqlStr := `
		SELECT
			v,
			id,
			user_type,
			user_name,
			pass_confirmed,
			email_confirmed,
			mfa_confirmed
		FROM gf_login_attempts
		WHERE user_name = $1 AND deleted = FALSE`

	rows, err := pRuntimeSys.SQLdb.QueryContext(pCtx, sqlStr, string(pUserNameStr))
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to get login_attempt by user_name from the DB",
			"sql_query_execute",
			map[string]interface{}{
				"user_name_str": string(pUserNameStr),
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}
	defer rows.Close()

	loginAttemptsLst := []*GFloginAttempt{}
	for rows.Next() {
		var loginAttempt GFloginAttempt
		err := rows.Scan(
			&loginAttempt.Vstr,
			&loginAttempt.IDstr,
			&loginAttempt.UserTypeStr,
			&loginAttempt.UserNameStr,
			&loginAttempt.PassConfirmedBool,
			&loginAttempt.EmailConfirmedBool,
			&loginAttempt.MFAconfirmedBool)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to get login_attempt from cursor",
				"sql_row_scan",
				map[string]interface{}{
					"user_name_str": string(pUserNameStr),
				},
				err, "gf_identity_core", pRuntimeSys)
			return nil, gfErr
		}
		loginAttemptsLst = append(loginAttemptsLst, &loginAttempt)
	}

	if len(loginAttemptsLst) > 0 {
		loginAttempt := loginAttemptsLst[0]
		return loginAttempt, nil
	}
	return nil, nil
}

//---------------------------------------------------

func DBsqlLoginAttemptUpdate(pLoginAttemptID *gf_core.GF_ID,
	pUpdateOp   *GFloginAttemptUpdateOp,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	fieldsTargets := ""
	paramsLst := []interface{}{*pLoginAttemptID}

	if pUpdateOp.PassConfirmedBool != nil {
		fieldsTargets += "pass_confirmed = $2,"
		paramsLst = append(paramsLst, *pUpdateOp.PassConfirmedBool)
	}
	if pUpdateOp.EmailConfirmedBool != nil {
		fieldsTargets += "email_confirmed = $3,"
		paramsLst = append(paramsLst, *pUpdateOp.EmailConfirmedBool)
	}
	if pUpdateOp.MFAconfirmedBool != nil {
		fieldsTargets += "mfa_confirmed = $4,"
		paramsLst = append(paramsLst, *pUpdateOp.MFAconfirmedBool)
	}
	if pUpdateOp.DeletedBool != nil {
		fieldsTargets += "deleted = $5,"
		paramsLst = append(paramsLst, *pUpdateOp.DeletedBool)
	}

	if fieldsTargets == "" {
		return nil // No updates to be made
	}

	fieldsTargets = fieldsTargets[:len(fieldsTargets)-1]

	sqlStr := fmt.Sprintf(`
		UPDATE gf_login_attempts
		SET %s
		WHERE id = $1 AND deleted = FALSE`, fieldsTargets)

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr, paramsLst...)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to to update a login_attempt",
			"sql_query_execute",
			map[string]interface{}{
				"login_attempt_id_str": string(*pLoginAttemptID),
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
	CREATE TABLE IF NOT EXISTS gf_users (
		v                 VARCHAR(255),
		id                TEXT,
		deleted           BOOLEAN DEFAULT FALSE,
		creation_time     TIMESTAMP DEFAULT NOW(),
		user_type         VARCHAR(255), -- "admin" or "standard"
		user_name         TEXT NOT NULL,
		screen_name       TEXT,
		description       TEXT,
		addresses_eth     TEXT[],
		email             TEXT,
		email_confirmed   BOOLEAN DEFAULT FALSE,
		profile_image_url TEXT,
		banner_image_url  TEXT,

		PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS gf_auth0_session (
		v              VARCHAR(255),
		id             TEXT,
		deleted        BOOLEAN DEFAULT FALSE,
		creation_time  TIMESTAMP DEFAULT NOW(),
		login_complete BOOLEAN NOT NULL,
		access_token   TEXT,
		profile        JSON,

		PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS gf_users_creds (
		v             VARCHAR(255),
		id            TEXT,
		deleted       BOOLEAN DEFAULT FALSE,
		creation_time TIMESTAMP DEFAULT NOW(),
		user_id       TEXT,
		user_name     TEXT,
		pass_salt     TEXT,
		pass_hash     TEXT,

		PRIMARY KEY(id),
		FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);

	CREATE TABLE IF NOT EXISTS gf_users_email_confirm (
		user_name          TEXT,
		user_id            TEXT,
		confirm_code       TEXT,
		creation_time      TIMESTAMP DEFAULT NOW(),
	
		PRIMARY KEY(confirm_code),
		FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);

	CREATE TABLE IF NOT EXISTS gf_login_attempts (
		v               VARCHAR(255),
		id              TEXT,
		deleted         BOOLEAN DEFAULT FALSE,
		creation_time   TIMESTAMP DEFAULT NOW(),
		user_type       VARCHAR(255), -- "admin" or "standard"
		user_name       VARCHAR(255),
		pass_confirmed  BOOLEAN,
		email_confirmed BOOLEAN,
		mfa_confirmed   BOOLEAN,
	
		PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS gf_users_invite_list (
		user_email      TEXT UNIQUE NOT NULL,
		creation_time   TIMESTAMP DEFAULT NOW(),
		deleted         BOOLEAN DEFAULT FALSE,
	
		PRIMARY KEY(user_email)
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