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
	"time"
	"github.com/lib/pq"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// USER_CREATE

func DBsqlUserCreate(pUser *GFuser,
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
		banner_image_url
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
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

	_, err = sqlStatement.Exec(
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
// CREATE_TABLES

func DBsqlCreateTables(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

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
	`

	_, err := pRuntimeSys.SQLdb.Exec(sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create gf_identity related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// GET_BASIC_INFO_BY_ETH_ADDR

func DBsqlGetBasicInfoByETHaddr(pUserAddressETH GFuserAddressETH,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	sqlStr := `
	SELECT id
	FROM gf_users
	WHERE $1 = ANY(addresses_eth) AND deleted = false
	LIMIT 1
	`

	var userIDstr string
	err := pRuntimeSys.SQLdb.QueryRow(sqlStr, pUserAddressETH).
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

	row := pRuntimeSys.SQLdb.QueryRow(sqlStr, pUserID)
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
		&user.BannerImageURLstr,
	)

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
	pRuntimeSys *gf_core.RuntimeSys) (GFuserName, *gf_core.GFerror) {

	sqlStr := `
		SELECT user_name
		FROM gf_users
		WHERE id = $1 AND deleted = false
		LIMIT 1
	`
	row := pRuntimeSys.SQLdb.QueryRow(sqlStr, string(pUserID))

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
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	sqlStr := `
		SELECT COUNT(*)
		FROM gf_users
		WHERE id = $1 AND deleted = false
	`

	var countInt int

	err := pRuntimeSys.SQLdb.QueryRow(sqlStr, string(pUserID)).Scan(&countInt)
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
