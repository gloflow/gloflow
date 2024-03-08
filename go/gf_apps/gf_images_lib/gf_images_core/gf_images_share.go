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
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//---------------------------------------------------

type GFshareInput struct {
	ImageID         GFimageID `mapstructure:"url_str"         validate:"required,min=5,max=400"`
	EmailAddressStr string    `json:"email_address"`
	EmailSubjectStr string    `json:"email_subject"`
	EmailBodyStr    string    `json:"email_body"`
}

//---------------------------------------------------

func SharePipeline(pInput *GFshareInput,
	pUserID     gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	

	emailGFtypeStr := "image_sharing"


	// SENDER_EMAIL_ADDRESS
	senderAddressStr, gfErr := gf_identity_core.DBsqlGetUserEmailByID(pUserID,
		pCtx,
		pRuntimeSys)
	
	userNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(pUserID,
		pCtx,
		pRuntimeSys)

	msgBodyHTMLstr := fmt.Sprintf(`
		<div>
			<div>
				GF user <b>%s</b> shared this image with you :) 
			</div>

			<div id='user_body'>
				%s
			</div>
		</div>`,
		userNameStr,
		pInput.EmailBodyStr)

	
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
	

	return nil
}