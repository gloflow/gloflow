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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func http__get_user_address_eth(p_req *http.Request,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (GF_user_address_eth, *gf_core.GF_error) {

	query_args_map := p_req.URL.Query()
	if values_lst, ok := query_args_map["addr_eth"]; ok {
		return GF_user_address_eth(values_lst[0]), nil
	} else {


		gf_err := gf_core.Mongo__handle_error("incoming http request is missing the addr_eth query-string arg",
			"verify__missing_key_error",
			map[string]interface{}{},
			nil, "gf_identity_lib", p_runtime_sys)
		return GF_user_address_eth(""), gf_err
	}
	return GF_user_address_eth(""), nil
}

//---------------------------------------------------
func verify__auth_proof_signature(p_signature_str string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {





	return false, nil
}