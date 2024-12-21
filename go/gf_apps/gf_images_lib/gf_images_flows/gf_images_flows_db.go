/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func DBgetAll(pCtx context.Context, pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	
	// MONGO
	mongoResultsLst, gfErr := DBmongoGetAll(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// SQL
	sqlResultsLst, gfErr := DBsqlGetAll(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-------------------
	// MERGE
	mergedMap := make(map[string]int)
	for _, result := range sqlResultsLst {
		flowName := result["flow_name"].(string)
		count := result["count"].(int)
		mergedMap[flowName] = count
	}

	for _, result := range mongoResultsLst {
		
		flowNameStr := result["_id"].(string) // result["flow_name"].(string)
		countInt := result["count"].(int)

		if existingCount, exists := mergedMap[flowNameStr]; exists {
			if countInt > existingCount {
				mergedMap[flowNameStr] = countInt
			}
		} else {
			mergedMap[flowNameStr] = countInt
		}
	}

	//-------------------
	// convert to list

	finalResultsLst := make([]map[string]interface{}, 0, len(mergedMap))
	for flowNameStr, countInt := range mergedMap {
		finalResultsLst = append(finalResultsLst, map[string]interface{}{
			"name_str":  flowNameStr,
			"count_int": countInt,
		})
	}

	//-------------------

	return finalResultsLst, nil
}
