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
	// "fmt"
	"context"
	// "strings"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

var logFun func(p_g string, p_m string)
var logNewFun gf_core.GFlogFun
var cliArgsMap map[string]interface{}

//-------------------------------------------------

func CreateTestPolicy(pTargetResourceID gf_core.GF_ID,
	pOwnerUserID gf_core.GF_ID,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GFpolicy, *gf_core.GFerror) {

	targetResourceTypeStr := "flow"

	policy, gfErr := PipelineCreate(pTargetResourceID,
		targetResourceTypeStr,
		pOwnerUserID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	spew.Dump(policy)

	return policy, nil
}

//-------------------------------------------------

func Tinit(pServiceNameStr string,
	pCliArgsMap map[string]interface{}) *gf_core.RuntimeSys {

	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: pServiceNameStr, // "gf_identity_tests",
		LogFun:         logFun,
		LogNewFun:      logNewFun,
		Validator:      gf_core.ValidateInit(),
	}

	//--------------------
	// SQL

	dbNameStr := "gf_tests"
	dbUserStr := "gf"

	dbHostStr := pCliArgsMap["sql_host_str"].(string)

	sqlDB, gfErr := gf_core.DBsqlConnect(dbNameStr,
		dbUserStr,
		"", // config.SQLpassStr,
		dbHostStr,
		runtimeSys)
	if gfErr != nil {
		panic(-1)
	}

	runtimeSys.SQLdb = sqlDB

	//--------------------
	
	return runtimeSys
}