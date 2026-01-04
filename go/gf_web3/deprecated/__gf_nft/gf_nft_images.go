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

package gf_nft

import (
	// "fmt"
	"context"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//---------------------------------------------------
// CREATE_AS_IMAGES_IN_FLOWS

func createAsImagesInFlows(pNFTsLst []*GFnft,
	pFlowsNamesLst []string,
	pUserID        gf_core.GF_ID,
	pJobsMngrCh    chan gf_images_jobs_core.JobMsg,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {
		
	//---------------------
	// GF_IMAGES_JOB

	clientTypeStr := "gf_web3:gf_nft"
	imagesExternURLsLst      := []string{}
	imagesOriginPagesURLsStr := []string{}

	for _, nft := range pNFTsLst {

		nftMediaURLstr := strings.TrimSpace(nft.MediaURIgatewayStr)
		imagesExternURLsLst = append(imagesExternURLsLst, nftMediaURLstr)

		// imageOriginPageURL == "nft" used because this image is not coming from any origin page (not scraped or
		// pinned from any html page), instead it belongs to an NFT.
		imagesOriginPagesURLsStr = append(imagesOriginPagesURLsStr, "nft")
	}

	_, imagesThumbSmallRelativeURLlst, imagesIDsLst, gfErr := gf_images_flows.AddExternImages(imagesExternURLsLst,
		imagesOriginPagesURLsStr,
		pFlowsNamesLst,
		clientTypeStr,
		pUserID,
		pJobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//---------------------
	// DB
	i := 0
	for _, nft := range pNFTsLst {

		gfImageID          := imagesIDsLst[i]
		gfImageThumbURLstr := imagesThumbSmallRelativeURLlst[i]

		gfErr := DBmongoUpdateGFimageProps(nft.IDstr,
			gfImageID,
			gfImageThumbURLstr,
			pCtx,
			pRuntimeSys)
		
		if gfErr != nil {
			// do nothing for now, let other image_ids be updated
		}

		i++
	}

	//---------------------

	return nil
}