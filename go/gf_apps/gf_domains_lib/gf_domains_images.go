/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_domains_lib

import (
	// "fmt"
	"strings"
	"net/url"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
//IMPORTANT!! - this statistic used by the gf_domains GF app, directly by the end-user
//              (not only by the admin user)

type GFdomainImages struct {
	Name_str            string         `bson:"_id"`
	Count_int           int            `bson:"count_int"`           // total count of all subpages counts
	Subpages_Counts_map map[string]int `bson:"subpages_counts_map"` // counts of individual sub-page urls that images come from
}

func GetDomainsImagesDB(pRuntimeSys *gf_core.RuntimeSys) ([]GFdomainImages, *gf_core.GFerror) {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	pRuntimeSys.LogFun("INFO",cyan("AGGREGATE IMAGES DOMAINS ")+yellow(">>>>>>>>>>>>>>>"))



	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{
				"t": "img",
				"origin_page_url_str": bson.M{"$exists": true,},
			}},
		},
		{
			{"$project", bson.M{
				"_id":                 false, // suppression of the "_id" field
				"origin_page_url_str": "$origin_page_url_str",
			}},
		},
		{
			{"$group", bson.M{
				"_id":       "$origin_page_url_str",
				"count_int": bson.M{"$sum": 1},
			}},
		},
		{
			{"$sort", bson.M{"count_int": -1}},
		},
	}


	/*pipe := pRuntimeSys.Mongo_coll.Pipe([]bson.M{
		//-------------------
		bson.M{"$match":bson.M{
				"t":                   "img",
				"origin_page_url_str": bson.M{"$exists": true,},
			},
		},

		//-------------------
		bson.M{"$project":bson.M{
				"_id":                 false, // suppression of the "_id" field
				"origin_page_url_str": "$origin_page_url_str",
			},
		},

		//-------------------
		// IMPORTANT!! - images dont store which domain they are from, instead they hold the URL of the page
		//               from which they originated.
		//               those page url's are then grouped by domain in the application layer
		//               (although idealy that join would be happening as a part of the aggregation pipeline)
		bson.M{"$group":bson.M{
				"_id":       "$origin_page_url_str",
				"count_int": bson.M{"$sum": 1},
			},
		},

		//-------------------
		bson.M{"$sort": bson.M{"count_int": -1},},
	})*/
	
	cursor, err := pRuntimeSys.Mongo_coll.Aggregate(ctx, pipeline)
	if err != nil {

		gfErr := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains images",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_domains_lib", pRuntimeSys)
		return nil, gfErr
	}
	defer cursor.Close(ctx)

	type ImagesOriginPage struct {
		Origin_page_url_str string `bson:"_id"`
		Count_int           int    `bson:"count_int"`
	}

	/*results_lst := []Images_Origin_Page{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains images",
			"mongodb_aggregation_error",
			nil, err, "gf_domains_lib", pRuntimeSys)
		return nil, gfErr
	}*/

	resultsLst := []ImagesOriginPage{}
	for cursor.Next(ctx) {

		var r ImagesOriginPage
		err := cursor.Decode(&r)
		if err != nil {
			gfErr := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains images",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_domains_lib", pRuntimeSys)
			return nil, gfErr
		}
	
		resultsLst = append(resultsLst, r)
	}

	//----------------------
	// FIX!!       - doesnt scale to large numbers of origin_page_url_str's.
	//               this should all be done in the DB
	// IMPORTANT!! - application-layer JOIN. starts with all unique origin_page_url_str's, 
	//               and then indexes their info by the domain to which they belong.

	domainsImagesMap := map[string]GFdomainImages{}
	for _, imagesOriginPage := range resultsLst {

		originPageURLstr := imagesOriginPage.Origin_page_url_str

		u, err := url.Parse(originPageURLstr)
		if err != nil {
			continue
		}

		domain_str := u.Host
		
		//--------------------
		// IMPORTANT!! - mongodb doesnt allow "." in the document keys. origin_page_url is a regular
		//               url with ".". This is used as a key in the Domain_Images "Subpages_Counts_map"
		//               member, and when stored in the mongodb they raise an error if not encoded.
		origin_page_url_no_dots_str := strings.Replace(originPageURLstr, ".", "+_=_+", -1)

		//--------------------

		if domain_images, ok := domainsImagesMap[domain_str]; ok {
			domain_images.Count_int                                        = domain_images.Count_int + imagesOriginPage.Count_int
			domain_images.Subpages_Counts_map[origin_page_url_no_dots_str] = imagesOriginPage.Count_int
		} else {

			//--------------------
			// domain_image - CREATE

			newDomainImages := GFdomainImages{
				Name_str:            domain_str,
				Count_int:           imagesOriginPage.Count_int,
				Subpages_Counts_map: map[string]int{
					origin_page_url_no_dots_str: imagesOriginPage.Count_int,
				},
			}

			domainsImagesMap[domain_str] = newDomainImages

			//--------------------
		}
	}

	// serialize map 
	domainImagesLst := []GFdomainImages{}
	for _, v := range domainsImagesMap {
		domainImagesLst = append(domainImagesLst, v)
	}

	//----------------------

	return domainImagesLst, nil
}