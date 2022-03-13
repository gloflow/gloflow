/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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
	// "fmt"
	"context"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// GET
func DBgetPolicy(pTargetResourceIDstr gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) (*GFpolicy, *gf_core.GF_error) {


	return nil, nil
}

// CREATE
func DBcreatePolicy(pPolicy *GFpolicy,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {

	return nil
}