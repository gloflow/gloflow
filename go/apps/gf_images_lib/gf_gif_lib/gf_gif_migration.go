/*
GloFlow media management/publishing system
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

package gf_gif_lib

import (
	"fmt"
	"time"
	"strings"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
//TEMPORARY - only used for a little while, until all GIF format images are also 
//            created in their GIF record form. new versions of crawler and chrome_ext 
//            logic has proper creation of both image and gif DB records, but old
//            image only DB records need to be processed
func Init_img_to_gif_migration(p_images_store_local_dir_path_str string,
	p_s3_bucket_name_str string,
	p_s3_info            *gf_core.Gf_s3_info,
	p_runtime_sys        *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_gif_migration.Init_img_to_gif_migration()")

	gf_domain_str := "gloflow.com"
	fmt.Println(gf_domain_str)


	/*migrate__fix_gif_urls(p_images_store_local_dir_path_str,
			gf_domain_str,
			p_s3_bucket_name_str,
			p_s3_client,
			p_s3_uploader,
			p_mongodb_coll,
			p_log_fun)*/


	/*migrate__create_gifs_from_images(p_images_store_local_dir_path_str,
			p_s3_bucket_name_str,
			p_s3_client,
			p_s3_uploader,
			p_mongodb_coll,
			p_log_fun)*/

	/*//DELETE_DUPLICATES
	go func() {
		for {
			pipe := p_mongodb_coll.Pipe([]bson.M{
				bson.M{"$match":bson.M{
						"t":"gif",
					},
				},

				bson.M{"$group":bson.M{
						"_id":      "$origin_url_str",
						"count_int":bson.M{"$sum":1},
						"ids_lst":  bson.M{"$push":"$id_str",},
					},
				},
			})

			results_lst := []map[string]interface{}{}
			err         := pipe.All(&results_lst)

			for _,r_map := range results_lst {
				if r_map["count_int"].(int) > 1 {

					for _,id_str := range r_map["ids_lst"].([]string) {

					}
				}
			}
		}
	}*/
	//--------------------------------------------------
	/*//ADD_GIFS_FLOW_NAME_TO_GIF_IMGS
	go func() {
		for {

			//get all images that dont have flows_names_lst
			var imgs_lst []gf_images_utils.GF_Image
			err := p_mongodb_coll.Find(bson.M{
					"t":              "img",
					"format_str":     "gif",
					"flows_names_lst":bson.M{"$nin":[]string{"gifs",}},
				}).All(&imgs_lst)

			//all img's are migrated
			if fmt.Sprint(err) == "not found" {
				return
			}
			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				continue
			}

			//update each one with flows_names_lst field, and add flow_name_str to that list
			for _,img := range imgs_lst {

				err := p_mongodb_coll.Update(bson.M{
						"t":         "img",
						"id_str":    img.Id_str,
						"format_str":"gif",
					},
					bson.M{"$push":bson.M{"flows_names_lst":"gifs",},})
				if err != nil {
					p_log_fun("ERROR",fmt.Sprint(err))
					continue
				}
			}
		}
	}()*/
	
}
//--------------------------------------------------
//FIX_GIF_URLS - try to fetch first frame image of a GIF, and if fails
//               regenerate GIF preview frames.

func migrate__fix_gif_urls(p_images_store_local_dir_path_str string,
	p_gf_domain_str      string,
	p_s3_bucket_name_str string,
	p_s3_info            *gf_core.Gf_s3_info,
	p_runtime_sys        *gf_core.Runtime_sys) {
	
	go func() {

		fmt.Println("ENTERING LOOP 1")
		for {

			time.Sleep(time.Second * 5)
			fmt.Println("FIX_GIF_URLS ----------- >>>")

			pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
				bson.M{"$match":bson.M{
						"t":  "gif",
						"$or":[]bson.M{
							//IMPORTANT!! - valid_bool is a new field. on most docs it doesnt exist yet,
							//              but on some it does. check that it exists, and if it does check 
							//              that its false (only invalid docs are fixed)
							bson.M{"valid_bool":bson.M{"$exists":false,},},
							bson.M{"valid_bool":false},
						},
					},
				},
				bson.M{"$sample":bson.M{
						"size":1,
					},
				},
			})

			var old_gif Gf_gif
			err := pipe.One(&old_gif)
			if err != nil {
				_ = gf_core.Error__create("failed to run FIX_GIF_URLS migration aggregation_pipeline to get a single GIF",
					"mongodb_aggregation_error",nil,err,"gf_gif_lib",p_runtime_sys)
				continue
			}

			//-----------------------
			//FETCH_GF_URL

			fmt.Println("  > origin_page_url - "+old_gif.Origin_page_url_str)
			fmt.Println("  > origin_url      - "+old_gif.Origin_url_str)
			fmt.Println("  > gf_url          - "+old_gif.Gf_url_str)

			//IMPORTANT!! - old_gif.Gf_url_str in is form - "/images/d/gif"
			gif__full_gf_url_str := fmt.Sprintf("http://%s%s",p_gf_domain_str,old_gif.Gf_url_str)
			gf_http_fetch,gf_err := gf_core.HTTP__fetch_url(gif__full_gf_url_str,p_runtime_sys)

			//IMPORTANT!! - http response body must be closed, regardless of if its used or not (by goland docs)
			defer gf_http_fetch.Resp.Body.Close()
			
			//IMPORTANT!! - common error for malformed url's is:
			//              "parse /images/d/gif/%!!(MISSING)s(*string=0xc4201f4040).gif: invalid URL escape "%!!(MISSING)s""
			//              so for them GIF's are rebuilt as well.
			if gf_err !=nil && strings.HasPrefix(fmt.Sprint(gf_err.Error),"parse") {

				rg__gf_err := migrate__rebuild_gif(&old_gif,
					p_images_store_local_dir_path_str,
					p_s3_bucket_name_str,
					p_s3_info,
					p_runtime_sys)
				if rg__gf_err != nil {
					continue
				}
			}

			if gf_err != nil {
				continue
			}

			

			//FAILED_TO_FETCH_GF_URL
			if !(gf_http_fetch.Status_code_int >= 200 && gf_http_fetch.Status_code_int < 400) {

				//REBUILD_GIF
				rg__gf_err := migrate__rebuild_gif(&old_gif,
					p_images_store_local_dir_path_str,
					p_s3_bucket_name_str,
					p_s3_info,
					p_runtime_sys)
				if rg__gf_err != nil {
					continue
				}
			}
			//-----------------------
			//FETCH_FIRST_PREVIEW_FRAME
			frame_url_str := old_gif.Preview_frames_s3_urls_lst[0]
		

			fpf__gf_http_fetch,gf_err := gf_core.HTTP__fetch_url(frame_url_str,p_runtime_sys)
			
			//IMPORTANT!! - common error for malformed url's is:
			//              "parse /images/d/gif/%!!(MISSING)s(*string=0xc4201f4040).gif: invalid URL escape "%!!(MISSING)s""
			//              so for them GIF's are rebuilt as well.
			if gf_err != nil && strings.HasPrefix(fmt.Sprint(gf_err.Error),"parse") {

				rg__gf_err := migrate__rebuild_gif(&old_gif,
					p_images_store_local_dir_path_str,
					p_s3_bucket_name_str,
					p_s3_info,
					p_runtime_sys)
				if rg__gf_err != nil {
					continue
				}
			}

			if gf_err != nil {
				continue
			}

			//IMPORTANT!! - http response body must be closed, regardless of if its used or not (by goland docs)
			defer fpf__gf_http_fetch.Resp.Body.Close()
			//-----------------------
		}
	}()
}
//--------------------------------------------------
func migrate__create_gifs_from_images(p_images_store_local_dir_path_str string,
	p_s3_bucket_name_str string,
	p_s3_info            *gf_core.Gf_s3_info,
	p_runtime_sys        *gf_core.Runtime_sys) {

	//--------------------------------------------------
	//CREATE_GIFS_FROM_IMAGES - for all 'img' DB objects with format 'gif', process it 
	//                          and create a 'gif' DB object

	go func() {

		fmt.Println("ENTERING LOOP 2")
		for {

			time.Sleep(time.Second * 5)
			fmt.Println("CREATE_GIFS_FROM_IMAGES ----------- >>>")
			
			//---------------------
			//IMPORTANT!! - get a truly random img with GIF format
			pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
				bson.M{"$match":bson.M{
						"t":         "img",
						"format_str":"gif",
					},
				},
				bson.M{"$sample":bson.M{
						"size":1,
					},
				},
			})

			var img gf_images_utils.Gf_image
			err := pipe.One(&img)
			if err != nil {
				_ = gf_core.Error__create("failed to run CREATE_GIFS_FROM_IMAGES migration aggregation_pipeline to get a single GIF",
					"mongodb_aggregation_error",nil,err,"gf_gif_lib",p_runtime_sys)
				continue
			}
			//---------------------

			var gif Gf_gif
			err = p_runtime_sys.Mongodb_coll.Find(bson.M{
					"t":                  "gif",
					"origin_url_str":     img.Origin_url_str,
					"title_str":          bson.M{"$exists":1,}, //IMPORTANT!! - new field added
					"origin_page_url_str":bson.M{"$exists":1,}, //IMPORTANT!! - new field added
					"tags_lst":           bson.M{"$exists":1,}, //IMPORTANT!! - new field added
				}).One(&gif)

			//IMPORTANT!! - a "gif" object was not found in the DB for an "img"
			//              with "gif" format. so a new gif is created
			if fmt.Sprint(err) == "not found" {

				fmt.Println("=================================")
				fmt.Println("")
				fmt.Println("    MIGRATING DB IMG->GIF - "+img.Origin_url_str)
				fmt.Println("")
				fmt.Println("=================================")

				//IMPORTANT!! - emtpy because its not being used here (new GF_Image not created)
				//              p_create_new_db_img_bool is set to 'false'.
				image_client_type_str := ""

				flows_names_lst,gf_err := migrate__get_flows_names(img.Id_str,p_runtime_sys)
				if gf_err != nil {
					continue
				}

				//IMPORTANT!! - image is re-fetched and GIF is processed in full
				_,_,gf_err = Process(img.Origin_url_str, //p_image_source_url_str,
					img.Origin_page_url_str,
					p_images_store_local_dir_path_str,
					image_client_type_str,
					flows_names_lst,
					false, //p_create_new_db_img_bool
					p_s3_bucket_name_str,
					p_s3_info,
					p_runtime_sys)

				if gf_err != nil {
					continue
				}
				continue
			}

			if err != nil {
				continue
			}
		}
	}()
}
//--------------------------------------------------
func migrate__rebuild_gif(p_old_gif *Gf_gif,
	p_images_store_local_dir_path_str string,
	p_s3_bucket_name_str              string,
	p_s3_info                         *gf_core.Gf_s3_info,
	p_runtime_sys                     *gf_core.Runtime_sys) *gf_core.Gf_error {

	//----------------
	//PROCESS_GIF

	//IMPORTANT!! - emtpy because its not being used here (new GF_Image not created)
	//              p_create_new_db_img_bool is set to 'false'.
	image_client_type_str := ""

	flows_names_lst,gf_err := migrate__get_flows_names(p_old_gif.Gf_image_id_str,p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	new_gif,gf_err := Process_and_upload(p_old_gif.Origin_url_str, //p_image_source_url_str
		p_old_gif.Origin_page_url_str, //p_image_origin_page_url_str
		p_images_store_local_dir_path_str,
		image_client_type_str,
		flows_names_lst,
		false, //p_create_new_db_img_bool
		p_s3_bucket_name_str,
		p_s3_info,
		p_runtime_sys)
		
	if gf_err != nil {
		return gf_err
	}
	//----------------
	//UPDATE_GIF_TO_OLD_CREATION_TIME - so that when sorted lists of GIFs from the DB are fetched
	//                                  these newly created GIFs are returned in proper positon.
	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":     "gif",
			"id_str":new_gif.Id_str,
		},
		bson.M{
			"$set":bson.M{"creation_unix_time_f":p_old_gif.Creation_unix_time_f},
		})

	if err != nil {
		gf_err := gf_core.Error__create("failed to update a new migrated GIF with an old creation_unix_time_f (of the old GIF) in mongodb",
			"mongodb_update_error",
			&map[string]interface{}{
				"old_gif_id_str":p_old_gif.Id_str,
				"new_gif_id_str":new_gif.Id_str,
			},
			err,"gf_gif_lib",p_runtime_sys)
		return gf_err
	}
	//----------------
	//DELETE_OLD_GIF - the one that was rebuilt
	gif_db__delete(p_old_gif.Id_str,p_runtime_sys)
	//----------------
	return nil
}
//--------------------------------------------------
func migrate__get_flows_names(p_gif__gf_image_id_str string,
	p_runtime_sys *gf_core.Runtime_sys) ([]string,*gf_core.Gf_error) {
	var flows_names_lst []string

	//IMPORTANT!! - GIF is not linked to a particular GF_Image
	if p_gif__gf_image_id_str != "" {

		var gf_img gf_images_utils.Gf_image
		err := p_runtime_sys.Mongodb_coll.Find(bson.M{"t":"img","id_str":p_gif__gf_image_id_str,}).One(&gf_img)
		if err != nil {
			gf_err := gf_core.Error__create("failed to find images with GIF id_str",
				"mongodb_find_error",
				&map[string]interface{}{"gif__gf_image_id_str":p_gif__gf_image_id_str,},
				err,"gf_gif_lib",p_runtime_sys)
			return nil,gf_err
		}

		//IMPORTANT!! - gf_img.Flows_names_lst is a new field, allowing images to belong to multiple
		//              flows. before there was only the Flow_name_str field.
		//              so in the beginning most GF_Image's will not have "flows_names_lst" set in the DB,
		//              and will contain the default value when loaded by mgo (empty list)
		if len(gf_img.Flows_names_lst) > 0 {

			flows_names_lst = gf_img.Flows_names_lst

			has_gif_flow_bool := false
			for _,s := range flows_names_lst {
				//flows list might contain "gifs" tag
				if s == "gifs" {
					has_gif_flow_bool = true
				}
			}
			//only add "gifs" flow if it doesnt exist
			if !has_gif_flow_bool {
				flows_names_lst = append(flows_names_lst,"gifs")
			}
		} else {
			flows_names_lst = []string{gf_img.Flow_name_str,"gifs",}
		}
	} else {
		flows_names_lst = []string{"gifs",}
	}

	return flows_names_lst,nil
}