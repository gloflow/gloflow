/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_identity_lib

import (
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func db__user__create(p_user *GF_user,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str := "gf_users"

	gf_err := gf_core.Mongo__insert(p_user,
		coll_name_str,
		map[string]interface{}{
			"user_id_str":       p_user.Id_str,
			"username_str":      p_user.Username_str,
			"description_str":   p_user.Description_str,
			"addresses_eth_lst": p_user.Addresses_eth_lst, 
			"caller_err_msg_str": "failed to insert GF_user into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}