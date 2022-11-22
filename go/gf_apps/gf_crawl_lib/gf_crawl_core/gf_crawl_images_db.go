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

package gf_crawl_core

import (
	"context"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------

func ImageDBcreate(p_img *GFcrawlerPageImage,
	pRuntime   *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	ctx := context.Background()
	c, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").CountDocuments(ctx, bson.M{
		"t":        "crawler_page_img",
		"hash_str": p_img.Hash_str,
	})

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to count the number of crawler_page_img's in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"img_ref_url_str":             p_img.Url_str,
				"img_ref_origin_page_url_str": p_img.Origin_page_url_str,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return false, gfErr
	}

	// crawler_page_img already exists, from previous crawls, so ignore it
	var exists_bool bool
	if c > 0 {
		pRuntimeSys.LogFun("INFO", yellow(">>>>>>>> - DB PAGE_IMAGE ALREADY EXISTS >-- ")+cyan(p_img.Url_str))
		
		exists_bool = true
		return exists_bool, nil
	} else {
			
		// IMPORTANT!! - only insert the crawler_page_img if it doesnt exist in the DB already
		ctx           := context.Background()
		coll_name_str := "gf_crawl"
		gfErr         := gf_core.MongoInsert(p_img,
			coll_name_str,
			map[string]interface{}{
				"img_ref_url_str":             p_img.Url_str,
				"img_ref_origin_page_url_str": p_img.Origin_page_url_str,
				"caller_err_msg_str":          "failed to insert a crawler_page_img into the DB",
			},
			ctx,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		exists_bool = false
		return exists_bool, nil
	}

	return false, nil
}

//--------------------------------------------------

func ImageDBcreateRef(p_img_ref *GFcrawlerPageImageRef,
	pRuntime    *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {


	ctx := context.Background()
	c, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").CountDocuments(ctx, bson.M{
		"t":        "crawler_page_img_ref",
		"hash_str": p_img_ref.Hash_str,
	})

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to count the number of crawler_page_img_ref's in mongodb",
			"mongodb_find_error",
			map[string]interface{}{
				"img_ref_url_str":             p_img_ref.Url_str,
				"img_ref_origin_page_url_str": p_img_ref.Origin_page_url_str,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}

	// crawler_page_img already exists, from previous crawls, so ignore it
	if c > 0 {
		return nil
	} else {
			
		// IMPORTANT!! - only insert the crawler_page_img if it doesnt exist in the DB already
		ctx           := context.Background()
		coll_name_str := "gf_crawl"
		gf_err        := gf_core.MongoInsert(p_img_ref,
			coll_name_str,
			map[string]interface{}{
				"img_ref_url_str":             p_img_ref.Url_str,
				"img_ref_origin_page_url_str": p_img_ref.Origin_page_url_str,
				"caller_err_msg_str":          "failed to insert a crawler_page_img_ref into the DB",
			},
			ctx,
			pRuntimeSys)
		if gf_err != nil {
			return gf_err
		}
	}

	return nil
}

//--------------------------------------------------

func imageDBget(p_id_str GFcrawlerPageImageID,
	pRuntime     *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) (*GFcrawlerPageImage, *gf_core.GFerror) {

	var img GFcrawlerPageImage
	ctx := context.Background()

	err := pRuntimeSys.Mongo_db.Collection("gf_crawl").FindOne(ctx, bson.M{
			"t":      "crawler_page_img",
			"id_str": p_id_str,
		}).Decode(&img)

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to get crawler_page_img by ID from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"id_str": p_id_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}

	return &img, nil
}

//-------------------------------------------------

func ImagesDBgetRecent(pRuntimeSys *gf_core.RuntimeSys) ([]GFcrawlerRecentImages, *gf_core.GFerror) {

	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"t": "crawler_page_img"}},
		},
		{
			{"$sort", bson.M{"creation_unix_time_f": -1}},
		},
		{
			bson.E{Key: "$limit", Value: 2000},
		},
		{
			{"$group", bson.M{
				"_id":                      "$origin_page_url_domain_str", // "$domain_str",
				"imgs_count_int":           bson.M{"$sum":  1},
				"crawler_page_img_ids_lst": bson.M{"$push": "$id_str"},
				"creation_times_lst":       bson.M{"$push": "$creation_unix_time_f"},
				"urls_lst":                 bson.M{"$push": "$url_str"},
				"nsfv_ls":                  bson.M{"$push": "$nsfv_bool"},
				"origin_page_urls_lst":     bson.M{"$push": "$origin_page_url_str"},
			}},
		},
	}

	/*pipe := pRuntimeSys.Mongodb_db.C("gf_crawl").Pipe([]bson.M{
		bson.M{"$match": bson.M{
				"t": "crawler_page_img",
			},
		},
		bson.M{"$sort": bson.M{
				"creation_unix_time_f": -1,
			},
		},
		bson.M{"$limit": 2000},
		bson.M{"$group": bson.M{
				"_id":                      "$origin_page_url_domain_str", //"$domain_str",
				"imgs_count_int":           bson.M{"$sum":  1},
				"crawler_page_img_ids_lst": bson.M{"$push": "$id_str"},
				"creation_times_lst":       bson.M{"$push": "$creation_unix_time_f"},
				"urls_lst":                 bson.M{"$push": "$url_str"},
				"nsfv_ls":                  bson.M{"$push": "$nsfv_bool"},
				"origin_page_urls_lst":     bson.M{"$push": "$origin_page_url_str"},
			},
		},
	})*/

	
	
	cursor, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").Aggregate(ctx, pipeline)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get recent_images (crawler_page_img) by domain",
			"mongodb_aggregation_error",
			nil,
			err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)
	
	/*err := pipe.AllowDiskUse().All(&results_lst)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get recent_images (crawler_page_img) by domain",
			"mongodb_aggregation_error",
			nil, err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}*/
	
	results_lst := []GFcrawlerRecentImages{}
	err          = cursor.All(ctx, &results_lst)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get recent_images (crawler_page_img) by domain",
			"mongodb_cursor_decode",
			nil,
			err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}

	return results_lst, nil
}

//--------------------------------------------------

func image__db_mark_as_downloaded(p_image *GFcrawlerPageImage, pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_db.image__db_mark_as_downloaded()")

	ctx := context.Background()

	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"t": "crawler_page_img",

			// IMPORTANT!! - search by "hash_str", not "id_str", because p_image's id_str might not
			//               be the id_str of the p_image (with the same hash_str) that was written to the DB. 
			//               (it might be an old p_image from previous crawler runs. to conserve DB space the crawler
			//               system doesnt write duplicate crawler_page_img's to the DB. 
			"hash_str": p_image.Hash_str,
		},
		bson.M{
			"$set": bson.M{"downloaded_bool": true},
		})
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update an crawler_page_img downloaded flag by its hash",
			"mongodb_update_error",
			map[string]interface{}{"image_hash_str": p_image.Hash_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}

	return nil
}

//--------------------------------------------------

func image__db_set_gf_image_id(pGFimageIDstr gf_images_core.GFimageID,
	pImage      *GFcrawlerPageImage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_db.image__db_set_gf_image_id()")

	ctx := context.Background()

	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"t": "crawler_page_img",

			// IMPORTANT!! - search by "hash_str", not "id_str", because p_image's id_str might not
			//               be the id_str of the p_image (with the same hash_str) that was written to the DB. 
			//               (it might be an old p_image from previous crawler runs. to conserve DB space the crawler
			//               system doesnt write duplicate crawler_page_img's to the DB. 
			"hash_str": pImage.Hash_str,
		},
		bson.M{
			"$set": bson.M{"image_id_str": pGFimageIDstr},
		})
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update an crawler_page_img downloaded flag by its hash",
			"mongodb_update_error",
			map[string]interface{}{"image_hash_str": pImage.Hash_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}

	return nil
}

//--------------------------------------------------

func image__db_update_after_process(pPageImg *GFcrawlerPageImage,
	pGFimageIDstr gf_images_core.GFimageID,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_images_db.image__db_update_after_process()")

	ctx := context.Background()

	pPageImg.Valid_for_usage_bool = true
	pPageImg.GFimageIDstr         = pGFimageIDstr
	
	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"t":      "crawler_page_img",
			"id_str": pPageImg.IDstr,
		},
		bson.M{"$set": bson.M{
				// IMPORTANT!! - gf_image has been created for this page_image, and so the appropriate
				//               image_id_str needs to be set in the page_image DB record
				"image_id_str": pGFimageIDstr,

				// IMPORTANT!! - image has been transformed, and is ready to be used further
				//               by other apps/services, either for display, or further calculation
				"valid_for_usage_bool": true,
			},
		})

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update an crawler_page_img valid_for_usage flag and its image_id (Gf_image) by its ID",
			"mongodb_update_error",
			map[string]interface{}{
				"id_str":          pPageImg.IDstr,
				"gf_image_id_str": pGFimageIDstr,
			}, err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}
	return nil
}