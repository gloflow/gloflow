package gf_nft

import (
	"fmt"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_nft/gf_nft_extern_services"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func indexAddress(pAddressStr string,
	pServiceSourceStr string,
	pConfig           *gf_eth_core.GF_config,
	pJobsMngrCh       chan gf_images_jobs_core.JobMsg,
	pUserID           gf_core.GF_ID,
	pMetrics          *GFmetrics,
	pCtx              context.Context,
	pRuntimeSys       *gf_core.RuntimeSys) ([]*GFnft, *gf_core.GFerror) {



	// OPEN_SEA
	if pServiceSourceStr == "opensea" {
		nftsOpenSeaParsedLst, gfErr := gf_nft_extern_services.OpenSeaGetAllNFTsForAddress(pAddressStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		fmt.Println(nftsOpenSeaParsedLst)
	}

	// ALCHEMY
	if pServiceSourceStr == "alchemy" {

		chainStr := "eth"

		// GET_ALL
		nftsAlchemyLst, gfErr := gf_nft_extern_services.AlchemyGetAllNFTsForAddress(pAddressStr,
			pConfig.AlchemyAPIkeyStr,
			chainStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		// DB - persist alchemy records
		gfErr = DBcreateBulkAlchemyNFTs(nftsAlchemyLst,
			pMetrics,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}


		// create NFT records from Alchemy NFT records
		nftsLst, gfErr := createFromAlchemy(nftsAlchemyLst,
			pMetrics,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		// spew.Dump(nftsLst)

		//---------------------
		// GF_IMAGES
		flowsNamesLst := []string{
			fmt.Sprintf("nft:owner:%s", pAddressStr),
		}
		gfErr = createAsImagesInFlows(nftsLst,
			flowsNamesLst,
			pUserID,
			pJobsMngrCh,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		//---------------------
		
		return nftsLst, nil
	}

	return nil, nil
}