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

package gf_nft_extern_services

import (
	"fmt"
	"strings"
	"time"
	"encoding/json"
	"context"
	"github.com/parnurzeal/gorequest"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFnftAlchemy struct {
	Vstr               string             `bson:"v_str"` // schema_version
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	IDstr              gf_core.GF_ID      `bson:"id_str"`
	DeletedBool        bool               `bson:"deleted_bool"`
	CreationUNIXtimeF  float64            `bson:"creation_unix_time_f"`

	OwnerAddressStr    string `bson:"owner_address_str"`
	ContractAddressStr string `bson:"contract_address_str"`
	TokenIDstr         string `bson:"token_id_str"`
	TokenTypeStr       string `bson:"token_type_str"`
	TitleStr           string `bson:"title_str"`
	DescriptionStr     string `bson:"description_str"`
	ChainStr           string `bson:"chain_str"`

	// URIs
	TokenURIrawStr     string `bson:"token_uri_raw_str"`
	MediaURIrawStr     string `bson:"media_uri_raw_str"`

	// GATEWAY_URIs
	TokenURIgatewayStr string `bson:"token_uri_gateway_str"`
	MediaURIgatewayStr string `bson:"media_uri_gateway_str"`

	MetadataNameStr        string                   `bson:"metadata_name_str"`
	MetadataImageStr       string                   `bson:"metadata_image_str"`
	MetadataExternalURLstr string                   `bson:"metadata_external_url_str"`
	MetadataAttributesLst  []map[string]interface{} `bson:"metadata_attributes_lst"`
}

//-------------------------------------------------

func AlchemyGetAllNFTsForAddress(pOwnerAddressStr string,
	pAPIkeyStr  string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnftAlchemy, *gf_core.GFerror) {

	//---------------------
	// GET_ALL_PAGES
	nftsParsedLst := []*GFnftAlchemy{}

	var pageKeyStr string
	for ;; {

		nftsParsedPageLst, newPageKeyStr, gfErr := AlchemyQueryByOwnerAddress(pOwnerAddressStr,
			pageKeyStr,
			pAPIkeyStr,
			pChainStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		nftsParsedLst = append(nftsParsedLst, nftsParsedPageLst...)

		// more pages left
		if newPageKeyStr != "" {
			pageKeyStr = newPageKeyStr
			continue

		} else {

			// last page, exit loop
			break
		}
	}

	//---------------------

	return nftsParsedLst, nil
}

//-------------------------------------------------

func AlchemyQueryByOwnerAddress(pOwnerAddressStr string,
	pPageKeyStr string,
	pAPIkeyStr  string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnftAlchemy, string, *gf_core.GFerror) {

	//---------------------
	// HTTP_AGENT
	var HTTPagent *gorequest.SuperAgent
	
	// sometimes OpenSea blocks access from certain locations
	if pRuntimeSys.HTTPproxyServerURIstr != "" {
		HTTPagent = gorequest.New().Proxy(pRuntimeSys.HTTPproxyServerURIstr)
	} else {
		HTTPagent = gorequest.New()
	}

	//---------------------
	
	var alchemyChainStr string 
	if pChainStr == "eth" {
		alchemyChainStr = "eth-mainnet"
	}

	qsMap := map[string]string{
		"owner": pOwnerAddressStr,
	}

	// user acquires a alchemy page key if there are more pages of data
	// left to fetch. it is returned in the response for that collection of data.
	if pPageKeyStr != "" {
		qsMap["pageKey"] = pPageKeyStr
	}

	qsLst := []string{}
	for k, v := range qsMap {
		qsLst = append(qsLst, fmt.Sprintf("%s=%s", k, v))
	}
	qsStr := strings.Join(qsLst, "&")

	// endpoint supported on Ethereum/Polygon/Flow chains
	// https://docs.alchemy.com/reference/getnfts
	urlStr := fmt.Sprintf("https://%s.g.alchemy.com/v2/%s/getNFTs?%s",
		alchemyChainStr,
		pAPIkeyStr,
		qsStr)
	_, bodyStr, errs := HTTPagent.Get(urlStr).End()

	fmt.Printf("url - %s\n", urlStr)
	spew.Dump(bodyStr)

	if (len(errs) > 0) {
		gfErr := gf_core.ErrorCreate("failed to query alchemy http API for all assets for address",
			"http_client_req_error",
			map[string]interface{}{
				"extern_service_str": "alchemy",
				"owner_address_str":  pOwnerAddressStr,
				"url_str":            urlStr,
			},
			errs[0], "gf_nft_extern_services", pRuntimeSys)
		return nil, "", gfErr
	}

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		gfErr := gf_core.ErrorCreate("failed to decode JSON response from alchemy http API for all assets for address",
			"json_decode_error",
			map[string]interface{}{
				"extern_service_str": "alchemy",
				"owner_address_str":  pOwnerAddressStr,
				"url_str":            urlStr,
			},
			err, "gf_nft_extern_services", pRuntimeSys)
		return nil, "", gfErr
	}

	// spew.Dump(bodyMap)

	nftsLst := bodyMap["ownedNfts"].([]interface{})

	nftsParsedPageLst := []*GFnftAlchemy{}
	for _, nft := range nftsLst {

		nftMap := nft.(map[string]interface{})
		contractMap := nftMap["contract"].(map[string]interface{})
		contractAddressStr := contractMap["address"].(string)

		idMap      := nftMap["id"].(map[string]interface{})
		tokenIDstr := idMap["tokenId"].(string)

		creationTimeUNIXf := float64(time.Now().UnixNano()) / 1_000_000_000.0
		idStr             := CreateID([]string{contractAddressStr, tokenIDstr,}, creationTimeUNIXf)

		gfNFT := &GFnftAlchemy{
			Vstr:              "0",
			IDstr:             idStr,
			CreationUNIXtimeF: creationTimeUNIXf,

			OwnerAddressStr:    pOwnerAddressStr,
			ContractAddressStr: contractAddressStr,
			TokenIDstr:         tokenIDstr,
			TokenTypeStr:       nftMap["id"].(map[string]interface{})["tokenMetadata"].(map[string]interface{})["tokenType"].(string),
			TitleStr:           nftMap["title"].(string),
			DescriptionStr:     nftMap["description"].(string),
			ChainStr:           pChainStr,

			// URIs
			TokenURIrawStr:     nftMap["tokenUri"].(map[string]interface{})["raw"].(string),
			MediaURIrawStr:     nftMap["media"].([]interface{})[0].(map[string]interface{})["raw"].(string),

			// GATEWAY_URIs
			TokenURIgatewayStr: nftMap["tokenUri"].(map[string]interface{})["gateway"].(string),
			MediaURIgatewayStr: nftMap["media"].([]interface{})[0].(map[string]interface{})["gateway"].(string),

			MetadataNameStr:        nftMap["metadata"].(map[string]interface{})["name"].(string),
			MetadataImageStr:       nftMap["metadata"].(map[string]interface{})["image"].(string),
			MetadataExternalURLstr: nftMap["metadata"].(map[string]interface{})["external_url"].(string),
		}

		// attributes
		attributesLst := []map[string]interface{}{}
		for _, attr := range nftMap["metadata"].(map[string]interface{})["attributes"].([]interface{}) {
			attrMap := attr.(map[string]interface{})
			attributesLst = append(attributesLst, attrMap)
		}
		gfNFT.MetadataAttributesLst = attributesLst
	
		fmt.Println("ALCHEMY FETCHED NFT >>>>>>>>")
		spew.Dump(gfNFT)

		nftsParsedPageLst = append(nftsParsedPageLst, gfNFT)
	}

	var pageKeyStr string
	if returnedPageKeyStr, ok := bodyMap["pageKey"]; ok {
		pageKeyStr = returnedPageKeyStr.(string)
	}

	return nftsParsedPageLst, pageKeyStr, nil
}