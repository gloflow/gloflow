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
type Gf_domain struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Id_str        string             `bson:"id_str"`
	T_str         string             `bson:"t"` // "domain"
	Name_str      string             `bson:"name_str"`
	Count_int     int                `bson:"count_int"`
	Domain_posts  Gf_domain_posts    `bson:"posts_domain"`
	Domain_images Gf_domain_images   `bson:"images_domain"`
}

//--------------------------------------------------
func Init_domains_aggregation(p_runtime_sys *gf_core.RuntimeSys) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_domains.Init_domains_aggregation()")

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
			
			gf_err := Discover_domains_in_db(p_runtime_sys)
			if gf_err != nil {
				continue
			}
		}
	}()
}

//--------------------------------------------------
func Discover_domains_in_db(p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_domains.Discover_domains_in_db()")

	// ADD!! - issue the posts/images queries in parallel via their own go-routines
	//---------------
	// POSTS
	posts_domains_lst, gf_err := Get_domains_posts__mongo(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------
	// IMAGES
	images_domains_lst, gf_err := Get_domains_images__mongo(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//---------------
	// APP_LEVEL_JOIN
	domains_map := accumulate_domains(posts_domains_lst, images_domains_lst, p_runtime_sys)
	// DB PERSIST
	gf_err = db__persist_domains(domains_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//--------------------

	return nil
}

//--------------------------------------------------
func accumulate_domains(p_posts_domains_lst []Gf_domain_posts,
	p_images_domains_lst []Gf_domain_images,
	p_runtime_sys        *gf_core.RuntimeSys) map[string]Gf_domain {
	// p_runtime_sys.LogFun("FUN_ENTER", "gf_domains.accumulate_domains()")

	domains_map := map[string]Gf_domain{}

	//--------------------------------------------------
	// POSTS DOMAINS
	// IMPORTANT!! - these run first so they just create a Domain struct without checks

	for _,domain_posts := range p_posts_domains_lst {

		domain_name_str := domain_posts.Name_str

		// IMPORTANT!! - no existing domain with this domain_str has been found
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := fmt.Sprintf("domain:%f", creation_unix_time_f)
		new_domain           := Gf_domain{
			Id_str:       id_str,
			T_str:        "domain",
			Name_str:     domain_name_str,
			Count_int:    domain_posts.Count_int,
			Domain_posts: domain_posts,
		}
		domains_map[domain_name_str] = new_domain
	}

	//--------------------------------------------------
	// IMAGES DOMAINS
	for _,images_domain := range p_images_domains_lst {

		domain_name_str := images_domain.Name_str

		if domain,ok := domains_map[domain_name_str]; ok {
			domain.Domain_images = images_domain
			domain.Count_int     = domain.Count_int + images_domain.Count_int
		} else {
			// IMPORTANT!! - no existing domain with this domain_str has been found
			creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
			id_str               := fmt.Sprintf("domain:%f", creation_unix_time_f)
			new_domain := Gf_domain{
				Id_str:        id_str,
				T_str:         "domain",
				Name_str:      domain_name_str,
				Count_int:     images_domain.Count_int,
				Domain_images: images_domain,
			}
			domains_map[domain_name_str] = new_domain
		}
	}

	//--------------------------------------------------

	return domains_map
}

//--------------------------------------------------
func db__persist_domains(p_domains_map map[string]Gf_domain,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {
	// p_runtime_sys.LogFun("FUN_ENTER","gf_domains.db__persist_domains()")

	// cyan   := color.New(color.FgCyan).SprintFunc()
	// yellow := color.New(color.FgYellow).SprintFunc()
	// white  := color.New(color.FgWhite).SprintFunc()

	ctx := context.Background()
	
	i := 0
	for _, d := range p_domains_map {

		// p_runtime_sys.LogFun("INFO",yellow("persisting ")+white("domain")+yellow(" "+fmt.Sprint(i)+" >---------------- ")+cyan(d.Name_str))

		// IMPORTANT!! -  finds a single document matching the provided selector document 
		//                and modifies it according to the update document. If no document 
		//                matching the selector is found, the update document is applied 
		//                to the selector document and the result is inserted in the collection

		// UPSERT
		query := bson.M{
			"t":        "domain",
			"name_str": d.Name_str,
		}
		gf_err := gf_core.MongoUpsert(query,
			d,
			map[string]interface{}{
				"domain_name_str":    d.Name_str,
				"caller_err_msg_str": "failed to persist a domain in mongodb",},
			p_runtime_sys.Mongo_coll,
			ctx, p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		i+=1
	}
	return nil
}

//--------------------------------------------------
func db__get_domains(p_runtime_sys *gf_core.RuntimeSys) ([]Gf_domain, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_domains.db__get_domains()")

	ctx := context.Background()

	q := bson.M{
		"t":         "domain",
		"count_int": bson.M{"$exists": true}, // "count_int" is a new required field, and we want those records, not the old ones
	}

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"count_int": -1}) // descending - true - sort the highest count first

	cursor, gf_err := gf_core.Mongo__find(q,
		find_opts,
		map[string]interface{}{
			"caller_err_msg_str": "failed to DB fetch all domains",
		},
		p_runtime_sys.Mongo_coll,
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	defer cursor.Close(ctx)





	/*coll_name_str := p_runtime_sys.Mongo_coll.Name()
	cursor, err      := p_runtime_sys.Mongo_db.Collection(coll_name_str).Find(ctx, q)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to get all domains",
			"mongodb_find_error",
			nil,
			err, "gf_domains_lib", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)*/


	results_lst := []Gf_domain{}
	for cursor.Next(ctx) {
		var domain Gf_domain
		err := cursor.Decode(&domain)
		if err != nil {
			gf_err := gf_core.MongoHandleError("failed to decode mongodb result of query to get domains",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_domains_lib", p_runtime_sys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, domain)
	}


	/*var results_lst []Gf_domain
	err := p_runtime_sys.Mongo_coll.Find(bson.M{
			"t":         "domain",
			"count_int": bson.M{"$exists":true}, //"count_int" is a new required field, and we want those records, not the old ones
		}).
		Sort("-count_int").
		All(&results_lst)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to get all domains",
			"mongodb_find_error",
			nil, err, "gf_domains_lib", p_runtime_sys)
		return nil, gf_err
	}*/

	return results_lst, nil
}