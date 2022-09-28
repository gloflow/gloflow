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
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_nft/gf_nft_extern_services"
)

//-------------------------------------------------
type GFnft struct {
	Vstr               string             `bson:"v_str"` // schema_version
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	IDstr              gf_core.GF_ID      `bson:"id_str"`
	DeletedBool        bool               `bson:"deleted_bool"`
	CreationUNIXtimeF  float64            `bson:"creation_unix_time_f"`

	OwnerAddressStr    string `bson:"owner_address_str"`
	TokenIDstr         string `bson:"token_id_str"`
	ContractAddressStr string `bson:"contract_address_str"`
	ContractNameStr    string `bson:"contract_name_str"`
	ChainStr           string `bson:"chain_str"`

	// URIs
	TokenURIrawStr     string `bson:"token_uri_raw_str"`
	MediaURIrawStr     string `bson:"media_uri_raw_str"`

	// GATEWAY_URIs
	TokenURIgatewayStr string `bson:"token_uri_gateway_str"`
	MediaURIgatewayStr string `bson:"media_uri_gateway_str"`

	GFimageID          gf_images_core.GFimageID `bson:"gf_image_id_str"`
	GFimageThumbURLstr string                   `bson:"gf_image_thumb_url_str"`

	OpenSeaIDstr       gf_core.GF_ID            `bson:"open_sea_id_str"`
	AlchemyIDstr       gf_core.GF_ID            `bson:"alchemy_id_str"`
}

type GFnftExtern struct {
	OwnerAddressStr    string `json:"owner_address_str"`
	TokenIDstr         string `json:"token_id_str"`
	ContractAddressStr string `json:"contract_address_str"`
	ContractNameStr    string `json:"contract_name_str"`
	ChainStr           string `json:"chain_str"`
	
	TokenURIrawStr     string `json:"token_uri_raw_str"`
	MediaURIrawStr     string `json:"media_uri_raw_str"`

	// GATEWAY_URIs
	TokenURIgatewayStr string `json:"token_uri_gateway_str"`
	MediaURIgatewayStr string `json:"media_uri_gateway_str"`



	GFimageID          gf_images_core.GFimageID `json:"gf_image_id_str"`
	GFimageThumbURLstr string                   `json:"gf_image_thumb_url_str"`
}

//-------------------------------------------------
func getNFTexternForm(pNFTsLst []*GFnft) []*GFnftExtern {

	// export NFT data for public usage
	nftsExternLst := []*GFnftExtern{}
	for _, nft := range pNFTsLst {

		nftExtern := &GFnftExtern{
			OwnerAddressStr:    nft.OwnerAddressStr,
			TokenIDstr:         nft.TokenIDstr,
			ContractAddressStr: nft.ContractAddressStr,
			ContractNameStr:    nft.ContractNameStr,
			ChainStr:           nft.ChainStr,

			TokenURIrawStr: nft.TokenURIrawStr,
			MediaURIrawStr: nft.MediaURIrawStr,

			TokenURIgatewayStr: nft.TokenURIgatewayStr,
			MediaURIgatewayStr: nft.MediaURIgatewayStr,

			GFimageID:          nft.GFimageID,
			GFimageThumbURLstr: nft.GFimageThumbURLstr,
		}

		nftsExternLst = append(nftsExternLst, nftExtern)
	}
	return nftsExternLst
}

//-------------------------------------------------
func get(pTokenIDstr string,
	pCollectionNameStr string,
	pCtx               context.Context,
	pRuntimeSys        *gf_core.RuntimeSys) (*GFnft, *gf_core.GFerror) {


	
	return nil, nil
}

//---------------------------------------------------
// CREATE_FOR_ALCHEMY
func createFromAlchemy(pNFTsAlchemyLst []*gf_nft_extern_services.GFnftAlchemy,
	pMetrics    *GFmetrics,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnft, *gf_core.GFerror) {

	NFTsLst := []*GFnft{}
	for _, nftAlchemy := range pNFTsAlchemyLst {

		creationTimeUNIXf := float64(time.Now().UnixNano()) / 1_000_000_000.0

		idStr := gf_nft_extern_services.CreateID([]string{
			nftAlchemy.ContractAddressStr,
			nftAlchemy.TokenIDstr,},
			creationTimeUNIXf)

		nft := &GFnft{
			Vstr:  "0",
			IDstr: idStr,
			CreationUNIXtimeF:  creationTimeUNIXf,

			OwnerAddressStr:    nftAlchemy.OwnerAddressStr,
			TokenIDstr:         nftAlchemy.TokenIDstr,
			ContractAddressStr: nftAlchemy.ContractAddressStr,
			ContractNameStr:    nftAlchemy.TitleStr,
			ChainStr:           nftAlchemy.ChainStr,

			// URIs
			TokenURIrawStr: nftAlchemy.TokenURIrawStr,
			MediaURIrawStr: nftAlchemy.MediaURIrawStr,

			// GATEWAY_URIs
			TokenURIgatewayStr: nftAlchemy.TokenURIgatewayStr,
			MediaURIgatewayStr: nftAlchemy.MediaURIgatewayStr,

			AlchemyIDstr: nftAlchemy.IDstr,
		}

		NFTsLst = append(NFTsLst, nft)
	}

	// DB
	gfErr := DBcreateBulkNFTs(NFTsLst,
		pMetrics,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return NFTsLst, nil
}

//---------------------------------------------------
// CREATE_ID
func createID(pUserIdentifierStr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		pUserIdentifierStr,
	}
	gfIDstr := gf_core.ID__create(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}