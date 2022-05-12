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
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_nft/gf_nft_extern_services"
)

//-------------------------------------------------
type GFnft struct {
	Vstr               string             `bson:"v_str"` // schema_version
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	IDstr              gf_core.GF_ID      `bson:"id_str"`
	DeletedBool        bool               `bson:"deleted_bool"`
	CreationUNIXtimeF  float64            `bson:"creation_unix_time_f"`



	TokenIDstr         string `bson:"token_id_str"`
	ContractAddressStr string `bson:"contract_address_str"`
	ContractNameStr    string `bson:"contract_name_str"`
	CollectionNameStr  string `bson:"collection_name_str"`

	OpenSeaIDstr       string `bson:"open_sea_nft_id_str"`
}

//-------------------------------------------------
func indexAddress(pAddressStr string,
	pServiceSourceStr string,
	pConfig           gf_eth_core.GF_config,
	pCtx              context.Context,
	pRuntimeSys       *gf_core.Runtime_sys) *gf_core.GFerror {



	// OPEN_SEA
	if pServiceSourceStr == "opensea" {
		nftsOpenSeaParsedLst, gfErr := gf_nft_extern_services.OpenSeaGetAllNFTsForAddress(pAddressStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		fmt.Println(nftsOpenSeaParsedLst)
	}

	// ALCHEMY
	if pServiceSourceStr == "alchemy" {
		chainStr := "eth"
		nftsAlchemyParsedLst, gfErr := gf_nft_extern_services.AlchemyGetAllNFTsForAddress(pAddressStr,
			pConfig.AlchemyAPIkeyStr,
			chainStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		fmt.Println(nftsAlchemyParsedLst)

	}





	
	


	return nil
}

//-------------------------------------------------
func get(pTokenIDstr string,
	pCollectionNameStr string,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.Runtime_sys) (*GFnft, *gf_core.GFerror) {


	
	return nil, nil
}

//---------------------------------------------------
func create() *GFnft {

	nft := &GFnft{

	}
	return nft
}

//---------------------------------------------------
func createID(pUserIdentifierStr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		pUserIdentifierStr,
	}
	gfIDstr := gf_core.ID__create(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}