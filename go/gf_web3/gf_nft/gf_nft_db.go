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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_nft/gf_nft_extern_services"
)


//-------------------------------------------------

func DBupdateGFimageProps(pNFTid gf_core.GF_ID,
	pGFimageID          gf_images_core.GFimageID,
	pGFimageThumbURLstr *string,
	pCtx                context.Context,
	pRuntimeSys         *gf_core.RuntimeSys) *gf_core.GFerror {




	
	collNameStr := "gf_web3_nfts"


	_, err := pRuntimeSys.Mongo_db.Collection(collNameStr).UpdateMany(pCtx, bson.M{
			"id_str":       pNFTid,
			"deleted_bool": false,
		},
		bson.M{"$set": bson.M{
			"gf_image_id_str":        pGFimageID,
			"gf_image_thumb_url_str": pGFimageThumbURLstr,
		}})

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update gf_nft with new gf_image_id",
			"mongodb_update_error",
			map[string]interface{}{
				"nft_id_str":  pNFTid,
				"gf_image_id": pGFimageID,
			},
			err, "gf_nft", pRuntimeSys)
		return gfErr
	}

	return nil
}

//-------------------------------------------------

func DBgetByOwner(pAddressStr string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnft, *gf_core.GFerror) {


	collNameStr := "gf_web3_nfts"
	
	findOpts := options.Find()
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"owner_address_str": pAddressStr,
			"chain_str":         pChainStr,
			"deleted_bool":      false,
		},
		findOpts,
		map[string]interface{}{
			"owner_address_str":  pAddressStr,
			"caller_err_msg_str": "failed to get all NFTs for an owner address from the DB",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return nil, gfErr
	}

	
	
	var nftsLst []*GFnft
	err := cursor.All(pCtx, &nftsLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get all NFTs for an owner address from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_nft", pRuntimeSys)
		return nil, gfErr
	}


	return nftsLst, nil
}

//-------------------------------------------------

func DBcreateBulkNFTs(pNFTsLst []*GFnft,
	pMetrics    *GFmetrics,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_web3_nfts"

	filterDocsByFieldsLst := []map[string]string{}
	recordsLst            := []interface{}{}
	contractAddressesLst  := []string{}
	for _, nft := range pNFTsLst {
		
		// IMPORTANT!! - upsert NFT docs based on their token_id and contract_address,
		//               (without also including the owner_address as a filter).
		//               this way the given NFT record will always reflect
		//               only the latest owner.
		//               NFT ownership history has to be kept track of in some other way.
		docFilterMap := map[string]string{
			"token_id_str":         nft.TokenIDstr,
			"contract_address_str": nft.ContractAddressStr,
		}
		filterDocsByFieldsLst = append(filterDocsByFieldsLst, docFilterMap)


		recordsLst = append(recordsLst, interface{}(nft))
		contractAddressesLst = append(contractAddressesLst, nft.ContractAddressStr)
	}


	// DB_INSERT_BULK
	insertedNewDocsInt, gfErr := gf_core.MongoUpsertBulk(filterDocsByFieldsLst, recordsLst,
		collNameStr,
		map[string]interface{}{
			"contract_addresses_lst": contractAddressesLst,
			"caller_err_msg_str":     "failed to bulk insert NFTs into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// METRICS
	if pMetrics != nil {
		pMetrics.NftDBinsertsCount.Add(float64(insertedNewDocsInt))
	}

	return nil
}

//-------------------------------------------------

func DBcreateBulkAlchemyNFTs(pNFTsLst []*gf_nft_extern_services.GFnftAlchemy,
	pMetrics    *GFmetrics,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_web3_nfts_alchemy"

	filterDocsByFieldsLst := []map[string]string{}
	recordsLst            := []interface{}{}
	contractAddressesLst  := []string{}
	
	for _, nft := range pNFTsLst {
		
		// IMPORTANT!! - upsert NFT docs based on their token_id and contract_address,
		//               (without also including the owner_address as a filter).
		//               this way the given NFT record will always reflect
		//               only the latest owner.
		//               NFT ownership history has to be kept track of in some other way.
		docFilterMap := map[string]string{
			"token_id_str":         nft.TokenIDstr,
			"contract_address_str": nft.ContractAddressStr,
		}
		filterDocsByFieldsLst = append(filterDocsByFieldsLst, docFilterMap)

		recordsLst = append(recordsLst, interface{}(nft))
		contractAddressesLst = append(contractAddressesLst, nft.ContractAddressStr)
	}

	// DB_INSERT_BULK
	insertedNewDocsInt, gfErr := gf_core.MongoUpsertBulk(filterDocsByFieldsLst,
		recordsLst,
		collNameStr,
		map[string]interface{}{
			"contract_addresses_lst": contractAddressesLst,
			"caller_err_msg_str":     "failed to bulk insert Alchemy NFTs into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// METRICS
	if pMetrics != nil {
		pMetrics.NftAlchemyDBinsertsCount.Add(float64(insertedNewDocsInt))
	}

	return nil
}