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

	ContractAddressStr string `bson:"contract_address_str"`
	TokenIDstr         string `bson:"token_id_str"`
	TokenTypeStr       string `bson:"token_type_str"`
	TitleStr           string `bson:"title_str"`
	DescriptionStr     string `bson:"description_str"`

	TokenURIrawStr     string `bson:"token_uri_raw_str"`
	TokenURIgatewayStr string `bson:"token_uri_gateway_str"`

	MediaURIrawStr     string `bson:"media_uri_raw_str"`
	MediaURIgatewayStr string `bson:"media_uri_gateway_str"`

	MetadataNameStr        string                   `bson:"metadata_name_str"`
	MetadataImageStr       string                   `bson:"metadata_image_str"`
	MetadataExternalURLstr string                   `bson:"metadata_external_url_str"`
	MetadataAttributesLst  []map[string]interface{} `bson:"metadata_attributes_lst"`
}

//-------------------------------------------------
func AlchemyGetAllNFTsForAddress(pAddressStr string,
	pAPIkeyStr  string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) ([]*GFnftAlchemy, *gf_core.GFerror) {

	//---------------------
	// GET_ALL_PAGES
	nftsParsedLst := []*GFnftAlchemy{}

	offsetInt := 0
	limitInt := 50
	for ;; {
		nftsParsedPageLst, gfErr := AlchemyQueryByAddress(pAddressStr,
			offsetInt,
			limitInt,
			pAPIkeyStr,
			pChainStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		nftsParsedLst = append(nftsParsedLst, nftsParsedPageLst...)

		offsetInt += limitInt
		break
	}

	//---------------------

	return nftsParsedLst, nil
}

//-------------------------------------------------
func AlchemyQueryByAddress(pAddressStr string,
	pOffsetInt  int,
	pLimitInt   int,
	pAPIkeyStr  string,
	pChainStr   string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.Runtime_sys) ([]*GFnftAlchemy, *gf_core.GFerror) {

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
		"owner": pAddressStr,
	}
	qsLst := []string{}
	for k, v := range qsMap {
		qsLst = append(qsLst, fmt.Sprintf("%s=%s", k, v))
	}
	qsStr  := strings.Join(qsLst, "&")
	urlStr := fmt.Sprintf("https://%s.g.alchemy.com/v2/your-api-key/getNFTs?%s", alchemyChainStr, qsStr)
	_, bodyStr, errs := HTTPagent.Get(urlStr).End()

	fmt.Printf("url - %s\n", urlStr)
	spew.Dump(bodyStr)

	if (len(errs) > 0) {
		gfErr := gf_core.Error__create("failed to query opensea http API for all assets for address",
			"http_client_req_error",
			map[string]interface{}{
				"extern_service_str": "alchemy",
				"address_str": pAddressStr,
				"offset_int":  pOffsetInt,
				"limit_int":   pLimitInt,
				"url_str":     urlStr,
			},
			errs[0], "gf_nft_extern_services", pRuntimeSys)
		return nil, gfErr
	}

	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		gfErr := gf_core.Error__create("failed to decode JSON response from opensea http API for all assets for address",
			"json_decode_error",
			map[string]interface{}{
				"extern_service_str": "alchemy",
				"address_str": pAddressStr,
				"offset_int":  pOffsetInt,
				"limit_int":   pLimitInt,
				"url_str":     urlStr,
			},
			err, "gf_nft_extern_services", pRuntimeSys)
		return nil, gfErr
	}

	spew.Dump(bodyMap)



	nftsLst := bodyMap["ownedNfts"].([]interface{})

	nftsParsedPageLst := []*GFnftAlchemy{}
	for _, nft := range nftsLst {

		nftMap := nft.(map[string]interface{})
		contractMap := nftMap["contract"].(map[string]interface{})
		contractAddressStr := contractMap["address"].(string)

		idMap      := nftMap["id"].(map[string]interface{})
		tokenIDstr := idMap["tokenId"].(string)

		creationTimeUNIXf := float64(time.Now().UnixNano()) / 1_000_000_000.0
		idStr             := createID([]string{contractAddressStr, tokenIDstr,}, creationTimeUNIXf)

		gfNFT := &GFnftAlchemy{
			Vstr:              "0",
			IDstr:             idStr,
			CreationUNIXtimeF: creationTimeUNIXf,


			ContractAddressStr: contractAddressStr,
			TokenIDstr:         tokenIDstr,
			TokenTypeStr:       nftMap["id"].(map[string]interface{})["tokenMetadata"].(map[string]interface{})["tokenType"].(string),
			TitleStr:           nftMap["title"].(string),
			DescriptionStr:     nftMap["description"].(string),

			TokenURIrawStr:     nftMap["tokenUri"].(map[string]interface{})["raw"].(string),
			TokenURIgatewayStr: nftMap["tokenUri"].(map[string]interface{})["gateway"].(string),

			MediaURIrawStr:     nftMap["media"].([]map[string]interface{})[0]["raw"].(string),
			MediaURIgatewayStr: nftMap["media"].([]map[string]interface{})[0]["raw"].(string),

			MetadataNameStr:        nftMap["metadata"].(map[string]interface{})["name"].(string),
			MetadataImageStr:       nftMap["metadata"].(map[string]interface{})["image"].(string),
			MetadataExternalURLstr: nftMap["metadata"].(map[string]interface{})["external_url"].(string),
			MetadataAttributesLst:  nftMap["metadata"].(map[string]interface{})["attributes"].([]map[string]interface{}),
		}

	
	
	
		nftsParsedPageLst = append(nftsParsedPageLst, gfNFT)
	}






	return nftsParsedPageLst, nil
}