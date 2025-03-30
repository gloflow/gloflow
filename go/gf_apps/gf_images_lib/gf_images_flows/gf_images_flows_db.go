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
	"sort"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------

func dbGetPage(pFlowNameStr string,
	pCursorStartPositionInt int,
	pElementsNumInt         int,
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) ([]*gf_images_core.GFimage, *gf_core.GFerror) {

	//-------------------
	// SQL
	sqlPageLst, gfErr := dbSQLgetPage(pFlowNameStr,
		pCursorStartPositionInt,
		pElementsNumInt,
		pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-------------------
	// MONGO
	mongoPageLst, gfErr := dbMongoGetPage(pFlowNameStr,
		pCursorStartPositionInt,
		pElementsNumInt,
		pCtx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}

	//-------------------
	// MERGE

	imageMap := make(map[string]*gf_images_core.GFimage)

	// add mongo elements to the map first, so they can be overwritten if conflicting
	for _, img := range mongoPageLst {
		imageMap[string(img.IDstr)] = img
	}

	// add SQL elements second, they take precedence
	for _, img := range sqlPageLst {
		imageMap[string(img.IDstr)] = img
	}

	// serialize map into a list
	mergedLst := make([]*gf_images_core.GFimage, 0, len(imageMap))
	for _, img := range imageMap {
		mergedLst = append(mergedLst, img)
	}

	// sort by creation unix time, so that the newest images are first
	sort.Slice(mergedLst, func(i, j int) bool {
		return mergedLst[i].Creation_unix_time_f > mergedLst[j].Creation_unix_time_f
	})

	// trim the list to the required number of elements, so that the page
	// contains only the latest elements
	if len(mergedLst) > pElementsNumInt {
		mergedLst = mergedLst[:pElementsNumInt]
	}

	//-------------------

	return mergedLst, nil
}

//---------------------------------------------------

func DBgetAll(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	
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
	for _, resultMap := range sqlResultsLst {
		flowNameStr := resultMap["name_str"].(string)
		countInt    := resultMap["count_int"].(int)
		mergedMap[flowNameStr] = countInt
	}

	for _, resultMap := range mongoResultsLst {
		
		flowNameStr := resultMap["_id"].(string) // result["flow_name"].(string)
		countInt := int(resultMap["count_int"].(int32))

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

	flowsCountsLst := make([]map[string]interface{}, 0, len(mergedMap))
	for flowNameStr, countInt := range mergedMap {
		flowsCountsLst = append(flowsCountsLst, map[string]interface{}{
			"name_str":  flowNameStr,
			"count_int": countInt,
		})
	}

	//-------------------
	// sort by count
	// needed because after sql/mongo merge, the order is lost in a map
	sort.Slice(flowsCountsLst, func(i, j int) bool {

		return flowsCountsLst[i]["count_int"].(int) > flowsCountsLst[j]["count_int"].(int)
	})

	//-------------------
	return flowsCountsLst, nil
}
