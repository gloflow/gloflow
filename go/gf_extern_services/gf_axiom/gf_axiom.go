/*
GloFlow application and media management/publishing platform
Copyright (C) 2025 Ivan Trajkovic

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

package gf_axiom

import (
	"time"
	"context"
	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/axiom/ingest"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------------------------

func Emit(pEventMap map[string]interface{},
	pDatasetNameStr string,
	pCtx            context.Context,
	pClient         *axiom.Client,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	pEventMap[ingest.TimestampField] = time.Now()
	_, err := pClient.Datasets.IngestEvents(pCtx, pDatasetNameStr, []axiom.Event{
		pEventMap,
	})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to sent event to Axiom",
			"library_error",
			map[string]interface{}{
				"event_map": pEventMap,
				"axiom_dataset_name_str": pDatasetNameStr,
			},
			err, "gf_project", pRuntimeSys)
		return gfErr
	}
	return nil
}