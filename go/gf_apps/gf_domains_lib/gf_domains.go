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
	"fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/globalsign/mgo/bson"
	// "github.com/fatih/color"
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
	posts_domains_lst, gf_err := GetDomainsPostsDB(pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//---------------
	// IMAGES
	images_domains_lst, gf_err := GetDomainsImagesDB(pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//---------------
	// APP_LEVEL_JOIN
	domains_map := accumulateDomains(posts_domains_lst, images_domains_lst, pRuntimeSys)

	// DB PERSIST
	gf_err = dbPersistDomains(domains_map, pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//--------------------

	return nil
}

//--------------------------------------------------

func accumulateDomains(pPostsDomainsLst []GFdomainPosts,
	pImagesDomainsLst []GFdomainImages,
	pRuntimeSys       *gf_core.RuntimeSys) map[string]GFdomain {

	domainsMap := map[string]GFdomain{}

	//--------------------------------------------------
	// POSTS DOMAINS
	// IMPORTANT!! - these run first so they just create a Domain struct without checks

	for _, domain_posts := range pPostsDomainsLst {

		domain_name_str := domain_posts.Name_str

		// IMPORTANT!! - no existing domain with this domain_str has been found
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := fmt.Sprintf("domain:%f", creation_unix_time_f)
		new_domain           := GFdomain{
			Id_str:       id_str,
			T_str:        "domain",
			Name_str:     domain_name_str,
			Count_int:    domain_posts.Count_int,
			Domain_posts: domain_posts,
		}
		domainsMap[domain_name_str] = new_domain
	}

	//--------------------------------------------------
	// IMAGES DOMAINS
	for _, imagesDomain := range pImagesDomainsLst {

		domainNameStr := imagesDomain.Name_str

		if domain,ok := domainsMap[domainNameStr]; ok {
			domain.Domain_images = imagesDomain
			domain.Count_int     = domain.Count_int + imagesDomain.Count_int
		} else {
			// IMPORTANT!! - no existing domain with this domain_str has been found
			creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
			id_str               := fmt.Sprintf("domain:%f", creation_unix_time_f)
			new_domain := GFdomain{
				Id_str:        id_str,
				T_str:         "domain",
				Name_str:      domainNameStr,
				Count_int:     imagesDomain.Count_int,
				Domain_images: imagesDomain,
			}
			domainsMap[domainNameStr] = new_domain
		}
	}

	//--------------------------------------------------

	return domainsMap
}

//--------------------------------------------------

func dbPersistDomains(pDomainsMap map[string]GFdomain,
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

func dbGetDomains(pRuntimeSys *gf_core.RuntimeSys) ([]GFdomain, *gf_core.GFerror) {

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