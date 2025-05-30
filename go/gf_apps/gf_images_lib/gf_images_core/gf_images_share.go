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

package gf_images_core

import (
	"fmt"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//---------------------------------------------------

type GFshareInput struct {
	ImageID         GFimageID
	EmailAddressStr string
	EmailSubjectStr string
	EmailBodyStr    string
}

//---------------------------------------------------

func SharePipeline(pInput *GFshareInput,
	pUserID      gf_core.GF_ID,
	pServiceInfo *GFserviceInfo,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// DB
	
	/*
	// SENDER_EMAIL_ADDRESS
	senderAddressStr, gfErr := gf_identity_core.DBsqlGetUserEmailByID(pUserID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	*/

	// SENDER_USER_NAME
	userNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(pUserID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//------------------------
	// IMAGE

	image, gfErr := DBgetImage(pInput.ImageID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	/*
	image, gfErr := DBmongoGetImage(pInput.ImageID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	*/

	//------------------------

	imageURLstr := fmt.Sprintf("https://%s%s", pServiceInfo.DomainBaseStr, image.ThumbnailMediumURLstr)

	//------------------------

	senderAddressStr := pServiceInfo.EmailSharingSenderAddressStr

	msgBodyHTMLstr := fmt.Sprintf(`
		<div>
			<div style="margin-top: 20px;">

				<a href="https://%s" style="color:gray;">gloflow</a> user <span style="
				background-color: #e1e1e1;
				padding: 6px;
				padding-left: 9px;
				padding-right: 8px;
				font-weight: bold;">%s</span> shared this image with you :) 
			</div>
			<div id='user_body'>
				%s
			</div>
			<div id='image' style="margin-top: 40px;">
				<img src='%s' alt='image' style='width:50%%;'></img>
			</div>
		</div>`,
		pServiceInfo.DomainBaseStr,
		userNameStr,
		pInput.EmailBodyStr,
		imageURLstr)

	//------------------------
	// AWS
	gfErr = gf_aws.SESsendMessage(pInput.EmailAddressStr,
		senderAddressStr,
		pInput.EmailSubjectStr,
		msgBodyHTMLstr,
		"", // msgBodyTextStr,
		pRuntimeSys)
	
	if gfErr != nil {
		return gfErr
	}

	//------------------------
	// PLUGIN - email sending

	emailGFtypeStr := "image_sharing"
	
	gfErr = pRuntimeSys.ExternalPlugins.EmailSendingCallback(pInput.EmailAddressStr,
		senderAddressStr,
		pInput.EmailSubjectStr,
		emailGFtypeStr,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//------------------------
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMetaMap := map[string]interface{}{
			"image_id": pInput.ImageID,
			"user_id":  pUserID,
			"email":    pInput.EmailAddressStr,
		}
		gf_events.EmitApp(GF_EVENT_APP__IMAGE_SHARE,
			eventMetaMap,
			pRuntimeSys.AppNameStr,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------------
	

	return nil
}