/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_images_service

import (
	"strings"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------
// INIT_HANDLER_IMAGE_VIEW - registers the "/images/v/" handler, which renders the page of a single image.
//                           its shared by the gf_images service and the gf_images solo service.

func InitHandlerImageView(pAuthSubsystemTypeStr string,
	pTemplatesPathsMap map[string]string,
	pRPChandlerRuntime *gf_rpc_lib.GFrpcHandlerRuntime,
	pServiceInfo       *gf_images_core.GFserviceInfo,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// TEMPLATES
	templates, gfErr := templateLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// VIEW_IMAGE
	// renders a solo image
	
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, false, "/images/v/",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "GET" {

				//-----------------
				// INPUT

				pathStr    := pReq.URL.Path
				imageIDstr := strings.Replace(pathStr, "/images/v/", "", 1)
				imageID    := gf_images_core.GFimageID(imageIDstr)

				// auth is optional here - if user is logged in, we get the user ID, if not, we proceed as anonymous
				var userIDfinalStr gf_core.GF_ID
				userLoggedInBool, userID, gfErr := gf_identity_core.GetUserIDfromReq(pReq,
					pAuthSubsystemTypeStr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				if userLoggedInBool {
					userIDfinalStr = *userID
				}

				//-----------------
				// RENDER_TEMPLATE
				templateRenderedStr, gfErr := renderImageViewPage(imageID,
					templates.imagesViewTmpl,
					templates.imagesViewSubtemplatesNamesLst,
					userIDfinalStr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//-----------------
				pResp.Write([]byte(templateRenderedStr))

				//------------------
				// EVENT
				if pServiceInfo.EnableEventsAppBool {
					eventMeta := map[string]interface{}{
						"user_id":  userID,
						"image_id": imageIDstr,
						"path":     pathStr,
					}
					gf_events.EmitApp(gf_images_core.GF_EVENT_APP__IMAGE_VIEW_PAGE,
						eventMeta,
						pRuntimeSys.AppNameStr,
						userIDfinalStr,
						pCtx,
						pRuntimeSys)
				}

				//------------------
			}

			return nil, nil
		},
		pRPChandlerRuntime,
		pRuntimeSys)

	return nil
}
