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
	"context"
	"strings"
	"time"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFnftOpenSea struct {
	Vstr               string             `bson:"v_str"` // schema_version
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	IDstr              gf_core.GF_ID      `bson:"id_str"`
	DeletedBool        bool               `bson:"deleted_bool"`
	CreationUNIXtimeF  float64            `bson:"creation_unix_time_f"`
	
	OpenSeaIDstr    string                    `bson:"token_id_str"       mapstructure:"id"`
	NameStr         string                    `bson:"name_str"           mapstructure:"name"`
	DescriptionStr  string                    `bson:"description_str"    mapstructure:"description"`
	ExternalLinkStr string                    `bson:"external_link_str"  mapstructure:"external_link"`
	AssetContract   GFnftOpenSeaAssetContract `bson:"asset_contract_map" mapstructure:"asser_contract"`
	OwnerStr        GFnftOpenSeaOwner         `bson:"owner_str"          mapstructure:"owner"`
	PermaLinkStr    string                    `bson:"permalink_str"      mapstructure:"permalink"`

	TokenIDstr     string `bson:"token_id_str"     mapstructure:"token_id"`
	SalesNumberInt int    `bson:"sales_number_int" mapstructure:"num_sales"`

	// IMAGE
	ImageURLstr          string `bson:"image_url_str"           mapstructure:"image_url"`
	ImagePreviewURLstr   string `bson:"image_preview_url_str"   mapstructure:"image_preview_url"`
	ImageThumbnailURLstr string `bson:"image_thumbnail_url_str" mapstructure:"image_thumbnail_url"`
	ImageOriginalURLstr  string `bson:"image_original_url_str"  mapstructure:"image_original_url"`
	AnimationURLstr      string `bson:"animation_url_str"       mapstructure:"animation_url"`

	Collection GFnftOpenSeaCollection `bson:"collection_map"`

}

type GFnftOpenSeaAssetContract struct {
	AddressStr     string `mapstructure:"address"`
	CreatedDateStr string `mapstructure:"created_date"`
	NameStr        string `mapstructure:"name"`
	OwnerInt       int64  `mapstructure:"owner"`
	SymbolStr      string `mapstructure:"symbol"`
	DescriptionStr string `mapstructure:"description"`
}

type GFnftOpenSeaOwner struct {

}

type GFnftOpenSeaCollection struct {
	NameStr        string
	DescriptionStr string
	ImageURLstr    string
	ExternalURLstr string
	CreatedDateStr string `mapstructure:"created_date"`
}

//-------------------------------------------------

func OpenSeaGetAllNFTsForAddress(pAddressStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]GFnftOpenSea, *gf_core.GFerror) {

	//---------------------
	// GET_ALL_PAGES
	nftsOpenSeaParsedLst := []GFnftOpenSea{}

	offsetInt := 0
	limitInt := 50
	for ;; {
		nftsOpenSeaParsedPageLst, gfErr := OpenSeaQueryByAddress(pAddressStr,
			offsetInt,
			limitInt,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		nftsOpenSeaParsedLst = append(nftsOpenSeaParsedLst, nftsOpenSeaParsedPageLst...)

		offsetInt += limitInt
		break
	}

	//---------------------

	return nftsOpenSeaParsedLst, nil
}

//-------------------------------------------------

func OpenSeaQueryByAddress(pAddressStr string,
	pOffsetInt  int,
	pLimitInt   int,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]GFnftOpenSea, *gf_core.GFerror) {

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
	

	qsMap := map[string]string{
		"owner":           pAddressStr,
		"order_direction": "desc",
		"offset": fmt.Sprintf("%d", pOffsetInt),
		"limit":  fmt.Sprintf("%d", pLimitInt),
	}
	qsLst := []string{}
	for k, v := range qsMap {
		qsLst = append(qsLst, fmt.Sprintf("%s=%s", k, v))
	}
	qsStr  := strings.Join(qsLst, "&")
	urlStr := fmt.Sprintf("https://api.opensea.io/api/v1/assets?%s", qsStr)
	_, bodyStr, errs := HTTPagent.Get(urlStr).End()

	fmt.Printf("url - %s\n", urlStr)
	spew.Dump(bodyStr)

	if (len(errs) > 0) {
		gfErr := gf_core.ErrorCreate("failed to query opensea http API for all assets for address",
			"http_client_req_error",
			map[string]interface{}{
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
		gfErr := gf_core.ErrorCreate("failed to decode JSON response from opensea http API for all assets for address",
			"json_decode_error",
			map[string]interface{}{
				"address_str": pAddressStr,
				"offset_int":  pOffsetInt,
				"limit_int":   pLimitInt,
				"url_str":     urlStr,
			},
			err, "gf_nft_extern_services", pRuntimeSys)
		return nil, gfErr
	}

	spew.Dump(bodyMap)



	nftsLst := bodyMap["assets"].([]interface{})

	nftsOpenSeaParsedPageLst := []GFnftOpenSea{}
	for _, nftMap := range nftsLst {




		// load returned json map into a struct
		var nftOpenSea GFnftOpenSea
		err := mapstructure.Decode(nftMap, &nftOpenSea)
		if err != nil {

			gfErr := gf_core.ErrorCreate("failed to decode JSON nft asset for address into struct GFnftOpenSea",
				"mapstruct_decode",
				map[string]interface{}{
					"address_str": pAddressStr,
				},
				err, "gf_nft_extern_services", pRuntimeSys)
			return nil, gfErr
		}

		


		// standard props
		creationTimeUNIXf := float64(time.Now().UnixNano()) / 1_000_000_000.0
		idStr   := CreateID([]string{string(nftOpenSea.OpenSeaIDstr),}, creationTimeUNIXf)
		nftOpenSea.Vstr              = "0" 
		nftOpenSea.IDstr             = idStr
		nftOpenSea.DeletedBool       = false 
		nftOpenSea.CreationUNIXtimeF = creationTimeUNIXf
	
	
		nftsOpenSeaParsedPageLst = append(nftsOpenSeaParsedPageLst, nftOpenSea)
	}






	return nftsOpenSeaParsedPageLst, nil
}