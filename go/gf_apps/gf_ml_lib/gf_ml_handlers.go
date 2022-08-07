/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_ml_lib

import (
	"context"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
func initHandlers(pHTTPmux *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GF_error {

	//---------------------
	// DATASETS_CREATE - register a dataset

	gf_rpc_lib.CreateHandlerHTTPwithMux("/ml/datasets/register",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {
			return nil, nil
		},
		pHTTPmux,
		nil,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------

	return nil
}