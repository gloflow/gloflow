package gf_nft

import (
	"fmt"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_nft/gf_nft_extern_services"
)

//-------------------------------------------------
func indexAddress(pAddressStr string,
	pServiceSourceStr string,
	pConfig           *gf_eth_core.GF_config,
	pCtx              context.Context,
	pRuntimeSys       *gf_core.Runtime_sys) ([]*GFnft, *gf_core.GFerror) {



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
		nftsAlchemyLst, gfErr := gf_nft_extern_services.AlchemyGetAllNFTsForAddress(pAddressStr,
			pConfig.AlchemyAPIkeyStr,
			chainStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		// DB
		gfErr = DBcreateBulkAlchemyNFTs(nftsAlchemyLst,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}



		nftsLst, gfErr := createForAlchemy(nftsAlchemyLst,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		return nftsLst, nil
	}

	return nil, nil
}














