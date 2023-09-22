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
/*BosaC.Jan30.2020. <3 volim te zauvek*/

package gf_domains_lib

import (
	"time"
	"context"
	"net/url"
	"strings"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
// ADD!! - creation_time

type GFdomain struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Id_str        string             `bson:"id_str"`
	T_str         string             `bson:"t"` // "domain"
	Name_str      string             `bson:"name_str"`
	Count_int     int                `bson:"count_int"`
	Domain_posts  GFdomainPosts      `bson:"posts_domain"`
	Domain_images GFdomainImages     `bson:"images_domain"`
}

//--------------------------------------------------

func InitDomainsAggregation(pRuntimeSys *gf_core.RuntimeSys) {

	go func() {
		for ;; {

			//--------------------
			// IMPORTANT!! - RUN AGGREGATION EVERY Xs (since this is a demanding aggregation)
			//               this is run first, in the loop, so that initialy when this is
			//               initialized it doesnt run, and only later when service active 
			//               for a while it will run for its first iteration.
			time_to_sleep := time.Second*time.Duration(60*5) // 5min
			time.Sleep(time_to_sleep)

			//--------------------
			
			gf_err := DiscoverDomainsInDB(pRuntimeSys)
			if gf_err != nil {
				continue
			}
		}
	}()
}

//--------------------------------------------------

func DiscoverDomainsInDB(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// ADD!! - issue the posts/images queries in parallel via their own go-routines
	//---------------
	// POSTS
	posts_domains_lst, gf_err := DBmongoGetDomainsPosts(pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//---------------
	// IMAGES
	images_domains_lst, gf_err := DBmongoGetDomainsImages(pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//---------------
	// APP_LEVEL_JOIN
	domains_map := accumulateDomains(posts_domains_lst, images_domains_lst, pRuntimeSys)

	// DB PERSIST
	gf_err = dbMongoPersistDomains(domains_map, pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//--------------------

	return nil
}

//--------------------------------------------------

func dbMongoPersistDomains(pDomainsMap map[string]GFdomain,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// cyan   := color.New(color.FgCyan).SprintFunc()
	// yellow := color.New(color.FgYellow).SprintFunc()
	// white  := color.New(color.FgWhite).SprintFunc()

	ctx := context.Background()
	
	i := 0
	for _, d := range pDomainsMap {

		// pRuntimeSys.LogFun("INFO",yellow("persisting ")+white("domain")+yellow(" "+fmt.Sprint(i)+" >---------------- ")+cyan(d.Name_str))

		// IMPORTANT!! -  finds a single document matching the provided selector document 
		//                and modifies it according to the update document. If no document 
		//                matching the selector is found, the update document is applied 
		//                to the selector document and the result is inserted in the collection

		// UPSERT
		query := bson.M{
			"t":        "domain",
			"name_str": d.Name_str,
		}
		gfErr := gf_core.MongoUpsert(query,
			d,
			map[string]interface{}{
				"domain_name_str":    d.Name_str,
				"caller_err_msg_str": "failed to persist a domain in mongodb",},
			pRuntimeSys.Mongo_coll,
			ctx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		i+=1
	}
	return nil
}

//--------------------------------------------------

func dbMongoGetDomains(pRuntimeSys *gf_core.RuntimeSys) ([]GFdomain, *gf_core.GFerror) {

	ctx := context.Background()

	q := bson.M{
		"t":         "domain",
		"count_int": bson.M{"$exists": true}, // "count_int" is a new required field, and we want those records, not the old ones
	}

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"count_int": -1}) // descending - true - sort the highest count first

	cursor, gfErr := gf_core.MongoFind(q,
		find_opts,
		map[string]interface{}{
			"caller_err_msg_str": "failed to DB fetch all domains",
		},
		pRuntimeSys.Mongo_coll,
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	defer cursor.Close(ctx)

	results_lst := []GFdomain{}
	for cursor.Next(ctx) {
		var domain GFdomain
		err := cursor.Decode(&domain)
		if err != nil {
			gfErr := gf_core.MongoHandleError("failed to decode mongodb result of query to get domains",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_domains_lib", pRuntimeSys)
			return nil, gfErr
		}
	
		results_lst = append(results_lst, domain)
	}

	return results_lst, nil
}

//--------------------------------------------------

func DBmongoGetDomainsImages(pRuntimeSys *gf_core.RuntimeSys) ([]GFdomainImages, *gf_core.GFerror) {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	pRuntimeSys.LogNewFun("DEBUG", cyan("AGGREGATE IMAGES DOMAINS ")+yellow(">>>>>>>>>>>>>>>"), nil)

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


	/*
	pipe := pRuntimeSys.Mongo_coll.Pipe([]bson.M{
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
	})
	*/
	
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

	/*
	results_lst := []Images_Origin_Page{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains images",
			"mongodb_aggregation_error",
			nil, err, "gf_domains_lib", pRuntimeSys)
		return nil, gfErr
	}
	*/

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

//--------------------------------------------------

func DBmongoIndexInit(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	indexesKeysLst := [][]string{
		[]string{"t", }, // all stat queries first match on "t"
		[]string{"t", "origin_page_url_str"},
		[]string{"t", "name_str"},
		[]string{"t", "count_int"},
	}

	indexesNamesLst := []string{
		"by_type",
		"by_type_and_origin_page_url",
		"by_type_and_name",
		"by_type_and_count",
	}
	gfErr := gf_core.MongoEnsureIndex(indexesKeysLst, indexesNamesLst, "data_symphony", pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

/*example post:
{
	"_id" : ObjectId("54bc39f02c42e565d20c432a"),
	"poster_user_name_str" : "Ivan",
	"creation_datetime" : "2012-09-24T04:15:02.315146",
	"colors_lst" : [ ],
	"post_elements_lst" : [
		{
			"height_str" : null,
			"colors_lst" : [ ],
			"description_str" : null,
			"meta_dict" : {
				
			},
			"post_index_3_tpl" : [
				0,
				0,
				0
			],
			"id_str" : "pub_pe:Apollo program - Lunar Module - jetison:(0, 0, 0)",
			"tags_lst" : [ ],
			"type_str" : "link",
			"extern_url_str" : "http://www.hq.nasa.gov/alsj/S66-10987.jpg",
			"width_str" : null
		},
		{
			"height_str" : null,
			"img_thumbnail_large_url_str" : "/images/d/1b8cab084a27a768cc87e662facc2ccb_thumb_large.jpeg",
			"colors_lst" : [ ],
			"img_thumbnail_small_url_str" : "/images/d/1b8cab084a27a768cc87e662facc2ccb_thumb_small.jpeg",
			"meta_dict" : {
				
			},
			"post_index_3_tpl" : [
				0,
				0,
				0
			],
			"description_str" : null,
			"img_thumbnail_medium_url_str" : "/images/d/1b8cab084a27a768cc87e662facc2ccb_thumb_medium.jpeg",
			"id_str" : "pub_pe:Apollo program - Lunar Module - jetison:(1, 0, 0)",
			"tags_lst" : [ ],
			"type_str" : "image",
			"extern_url_str" : "http://www.hq.nasa.gov/alsj/S66-10987.jpg",
			"width_str" : null
		}
	],
	"description_str" : "",
	"title_str" : "Apollo program - Lunar Module - jetison",
	"images_ids_lst" : [
		"1b8cab084a27a768cc87e662facc2ccb"
	],
	"id_str" : "Apollo program - Lunar Module - jetison:2012-09-24 04:15:02.315146",
	"tags_lst" : [
		"nasa",
		"nasa_apollo",
		"moon",
		"space"
	],
	"type_str" : "gchrome_ext",
	"main_image_url_str" : "",
	"main_image_thumbnail_medium_url_str" : null,
	"t" : "post"
}*/

func DBmongoGetDomainsPosts(pRuntimeSys *gf_core.RuntimeSys) ([]GFdomainPosts, *gf_core.GFerror) {

	// cyan   := color.New(color.FgCyan).SprintFunc()
	// yellow := color.New(color.FgYellow).SprintFunc()
	// pRuntimeSys.LogFun("INFO",cyan("AGGREGATE POSTS DOMAINS ")+yellow(">>>>>>>>>>>>>>>"))

	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"t": "post"}},
		},
		{
			// IMPORTANT!! - for each document create a new one that only contains the "post_elements_lst"
			//               field, since this is where the url string is contained, in the post_element
			//               of type 'link'
			{"$project", bson.M{
				"_id":               false, //suppression of the "_id" field
				"post_elements_lst": "$post_elements_lst",
			}},
		},
		{
			// ATTENTION!! - potentially creates a large number of docs. for large amounts of original docs
			//               with large post_elements_lst members, this may overflow memory.
			//               figure out how to horizontally partition this op. 
			// 
			// create as many new docs as there are post_elements in the current doc,
			// where the new docs contain all the fields of the original source doc.
			{"$unwind", "$post_elements_lst"},
		},
		{
			{"$match", bson.M{"post_elements_lst.type_str": "link"}},
		},
		{
			{"$project", bson.M{
				"_id":                 false, // suppression of the "_id" field
				"post_extern_url_str": "$post_elements_lst.extern_url_str",
			}},
		},
	}

	/*pipe := pRuntimeSys.Mongo_coll.Pipe([]bson.M{

		bson.M{"$match": bson.M{"t": "post",},},

		//{"\$limit":10}, //DEBUGGING
		//---------------------
		//IMPORTANT!! - for each document create a new one that only contains the "post_elements_lst"
		//              field, since this is where the url string is contained, in the post_element
		//              of type 'link'
		bson.M{"$project":bson.M{
				"_id":               false, //suppression of the "_id" field
				"post_elements_lst": "$post_elements_lst",
			},
		},
		//---------------------
		//ATTENTION!! - potentially creates a large number of docs. for large amounts of original docs
		//              with large post_elements_lst members, this may overflow memory.
		//              figure out how to horizontally partition this op. 
		// 
		//create as many new docs as there are post_elements in the current doc,
		//where the new docs contain all the fields of the original source doc
		bson.M{"$unwind": "$post_elements_lst",},
		//---------------------
		//"$type_str" - attribute of "post_elements_lst" elements
		bson.M{"$match": bson.M{"post_elements_lst.type_str": "link",},},
		bson.M{"$project": bson.M{
				"_id":                 false, //suppression of the "_id" field
				"post_extern_url_str": "$post_elements_lst.extern_url_str",
			},
		},
	})*/
	
	cursor, err := pRuntimeSys.Mongo_coll.Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains posts",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_domains_lib", pRuntimeSys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	type Posts_extern_urls_result struct {
		Post_extern_url_str string `bson:"post_extern_url_str"`
	}

	/*results_lst := []Posts_extern_urls_result{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains posts",
			"mongodb_aggregation_error",
			nil, err, "gf_domains_lib", pRuntimeSys)
		return nil, gf_err
	}*/
	
	results_lst := []Posts_extern_urls_result{}
	for cursor.Next(ctx) {

		var r Posts_extern_urls_result
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get domains posts",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_domains_lib", pRuntimeSys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, r)
	}

	//---------------
	// IMPORTANT!! - application level join. move this to Db with the apporpriate "domain_str" field

	parsedDomainsMap := map[string]int{}
	for _, r := range results_lst {

		u,err := url.Parse(r.Post_extern_url_str)
		if err != nil {
			continue
		}
		domainStr := u.Host


		if _, ok := parsedDomainsMap[domainStr]; ok {
			parsedDomainsMap[domainStr] += 1
		} else {
			parsedDomainsMap[domainStr] = 1
		}
	}

	//---------------
	domainPostsLst := []GFdomainPosts{}
	for domainStr, countInt := range parsedDomainsMap {

		domainPosts := GFdomainPosts{
			Name_str:  domainStr,
			Count_int: countInt,
		}
		domainPostsLst = append(domainPostsLst, domainPosts)
	}

	// pRuntimeSys.LogFun("INFO", yellow(">>>>>>>> DOMAIN_POSTS FOUND - ")+cyan(fmt.Sprint(len(domainPostsLst))))
	//---------------

	return domainPostsLst, nil

	/*mongo console query - for testing
	db.posts.aggregate(
		{"$match":{"t":"post"}},
		{"$unwind":"$post_elements_lst"},
		{"$match":{"post_elements_lst.type_str":"link"}},
		{"$project":{"_id":0,"title_str":1,"url_domain_str":"$post_elements_lst.extern_url_domain_str"}},
		{"$group":{"_id":"$url_domain_str","titles_lst":{"$addToSet":"$title_str"},"count":{"$sum":1}}},
		{"$sort":{"count":-1}});*/
}