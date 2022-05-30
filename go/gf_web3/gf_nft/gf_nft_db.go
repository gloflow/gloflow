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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_nft/gf_nft_extern_services"
)

//-------------------------------------------------
func DBgetByOwner(pAddressStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnft, *gf_core.GFerror) {


	collNameStr := "gf_web3_nfts"
	

	findOpts := options.Find()
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"owner_address_str": pAddressStr,
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
		gfErr := gf_core.Mongo__handle_error("failed to get all NFTs for an owner address from cursor",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_nft", pRuntimeSys)
		return nil, gfErr
	}


	return nil, nil
}

//-------------------------------------------------
func DBcreateBulkNFTs(pNFTsLst []*GFnft,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_web3_nfts"

	IDsLst     := []string{}
	recordsLst := []interface{}{}
	contractAddressesLst := []string{}
	for _, nft := range pNFTsLst {
		IDsLst     = append(IDsLst, string(nft.IDstr))
		recordsLst = append(recordsLst, interface{}(nft))
		contractAddressesLst = append(contractAddressesLst, nft.ContractAddressStr)
	}


	// DB_INSERT_BULK
	gfErr := gf_core.MongoInsertBulk(IDsLst, recordsLst,
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


	return nil
}

//-------------------------------------------------
func DBcreateBulkAlchemyNFTs(pNFTsLst []*gf_nft_extern_services.GFnftAlchemy,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_web3_nfts_alchemy"

	IDsLst     := []string{}
	recordsLst := []interface{}{}
	contractAddressesLst := []string{}
	for _, nft := range pNFTsLst {
		IDsLst     = append(IDsLst, string(nft.IDstr))
		recordsLst = append(recordsLst, interface{}(nft))
		contractAddressesLst = append(contractAddressesLst, nft.ContractAddressStr)
	}

	// DB_INSERT_BULK
	gfErr := gf_core.MongoInsertBulk(IDsLst, recordsLst,
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

	return nil
}