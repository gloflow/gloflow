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

	

	emailGFtypeStr := "image_sharing"


	// SENDER_EMAIL_ADDRESS
	senderAddressStr, gfErr := gf_identity_core.DBsqlGetUserEmailByID(pUserID,
		pCtx,
		pRuntimeSys)
	
	userNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(pUserID,
		pCtx,
		pRuntimeSys)

	image, gfErr := DBmongoGetImage(pInput.ImageID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	imageURLstr := image.ThumbnailMediumURLstr

	msgBodyHTMLstr := fmt.Sprintf(`
		<div>
			<div>
				GF user <b>%s</b> shared this image with you :) 
			</div>

			<div id='user_body'>
				%s
			</div>
			<div id='image'>
				<img src='%s' alt='image' style='width:100%%;'>
			</div>
		</div>`,
		userNameStr,
		pInput.EmailBodyStr,
		imageURLstr)

	
	//------------------------
	// PLUGIN
	//------------------------
	// EMAIL_PLUGIN
	gfErr = pRuntimeSys.ExternalPlugins.EmailSendingCallback(pInput.EmailAddressStr,
		senderAddressStr,
		pInput.EmailSubjectStr,
		emailGFtypeStr,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
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
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMeta := map[string]interface{}{
			"image_id": pInput.ImageID,
			"user_id":  pUserID,
			"email":    pInput.EmailAddressStr,
		}
		gf_events.EmitApp(GF_ENVET_APP__IMAGE_SHARE,
			eventMeta,
			pRuntimeSys)
	}

	//------------------------
	

	return nil
}