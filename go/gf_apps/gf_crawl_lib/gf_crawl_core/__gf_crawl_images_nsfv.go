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
	"os"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/fatih/color"
	"github.com/koyachi/go-nude"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
)

//--------------------------------------------------
func images__stage__determine_are_nsfv(p_crawler_name_str string,
	p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	p_origin_page_url_str           string,
	p_runtime                       *Gf_crawler_runtime,
	p_runtime_sys                   *gf_core.RuntimeSys) []*gf_page_img__pipeline_info {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_crawl_images_nsfv.images__stage__determine_are_nsfv")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - determine_are_nsfv")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		//IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		gf_img := page_img__pinfo.page_img

		var is_nsfv_bool bool
		var gf_err       *gf_core.GFerror

		//--------------
		// GIF
		if gf_img.Img_ext_str == "gif" {

			is_nsfv_bool,gf_err = image__is_nsfv__gif(page_img__pinfo.local_file_path_str, gf_img.Url_str, p_runtime_sys)
			if gf_err != nil {
				p_runtime_sys.LogFun("ERROR", "failed to do nudity-detection/filtering in GIF - "+gf_img.Url_str+" - "+fmt.Sprint(gf_err))

				t:="gif_is_nsfv_test__failed"
				m:="failed nsfv testing of GIF with img_url_str - "+gf_img.Url_str
				Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str": p_origin_page_url_str,}, gf_img.Url_str, p_crawler_name_str,
					gf_err, p_runtime, p_runtime_sys)

				page_img__pinfo.gf_error = gf_err
				continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
			}

		//--------------
		// STATIC-IMAGE
		} else {

			//IMPORTANT!! - if image has nudity it is flagged as not valid
			is_nsfv_bool,gf_err = image__is_nsfv(page_img__pinfo.local_file_path_str, p_runtime_sys)
			if gf_err != nil {
				p_runtime_sys.LogFun("ERROR","failed to do nudity-detection/filtering in image - "+gf_img.Url_str+" - "+fmt.Sprint(gf_err))

				t:="image_is_nsfv_test__failed"
				m:="failed nsfv testing of image with img_url_str - "+gf_img.Url_str
				Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str": p_origin_page_url_str,}, gf_img.Url_str, p_crawler_name_str,
					gf_err, p_runtime, p_runtime_sys)

				page_img__pinfo.gf_error = gf_err
				continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
			}
		}

		//--------------
		// FLAG_IMAGE
		
		page_img__pinfo.nsfv_bool = is_nsfv_bool

		//IMPORTANT!! - if image is flagged as NSFV its flag in the DB is updated
		if is_nsfv_bool {
			gf_err = image__flag_as_nsfv(gf_img,p_runtime_sys)
			if gf_err != nil {
				
				t:="image_mark_as_nsfv__failed"
				m:="failed nsfv marking (in DB) of image with img_url_str - "+gf_img.Url_str
				Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str": p_origin_page_url_str,}, gf_img.Url_str, p_crawler_name_str,
					gf_err, p_runtime, p_runtime_sys)

				page_img__pinfo.gf_error = gf_err
				continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
			}
		}

		//--------------
	}
	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
//GIF

func image__is_nsfv__gif(p_img_gif_path_str string,
	p_img_gif_origin_url_str string,
	p_runtime_sys            *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {
	//p_runtime_sys.LogFun("FUN_ENTER","gf_crawl_images_nsfv.image__is_nsfv__gif()")

	cyan  := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgBlack).Add(color.BgGreen).SprintFunc()
	black := color.New(color.FgBlack).Add(color.BgWhite).SprintFunc()

	fmt.Println("INFO", green(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>------------------------------------------------"))
	fmt.Println("INFO", green(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>------------------------------------------------"))
	fmt.Println("INFO", "")
	fmt.Println("INFO", cyan("                          GIF")+" - "+cyan("GET_FRAMES"))
	fmt.Println("INFO", "")
	fmt.Println("INFO", black(p_img_gif_origin_url_str))
	fmt.Println("INFO", green(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>------------------------------------------------"))
	fmt.Println("INFO", green(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>------------------------------------------------"))

	frames_images_dir_path_str := "./"
	new_files_names_lst, gf_err := gf_gif_lib.Gif__frames__save_to_fs(p_img_gif_path_str,
		frames_images_dir_path_str,
		0, //p_frames_num_to_get_int //IMPORTANT!! - if its 0 it gets all frames
		p_runtime_sys)
	if gf_err != nil {
		return false,gf_err
	}

	//-------------------------
	// CLEANUP!! - remove extracted GIF frames in dir frames_images_dir_path_str,
	//             after NSFV analysis is complete
	defer func() {
		for _, f_str := range new_files_names_lst {
			err := os.Remove(f_str)
			if err != nil {

			}
		}
	}()

	//-------------------------
	// IMPORTANT!! - run NSFV detection on each GIF frame, and for the first one that fails the test
	//               use it as a signal to mark the whole GIF as NSFV
	for _,frame_image_file_path_str := range new_files_names_lst {
		is_nsfv_bool,gf_err := image__is_nsfv(frame_image_file_path_str, p_runtime_sys)
		if gf_err != nil {
			return false, gf_err
		}

		//-----------------
		// IMPORTANT!! - first frame that fails the NSFV test indicates the whole GIF is NSFV
		if is_nsfv_bool {
			return is_nsfv_bool, nil
		}

		//-----------------
	}

	//-------------------------

	is_nsfv_bool := false //if all frames pass as non-nsfv then the GIF is not NSFV
	return is_nsfv_bool, nil
}

//--------------------------------------------------
func image__is_nsfv(p_img_path_str string,
	p_runtime_sys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	is_nude_bool,err := nude.IsNude(p_img_path_str)
	if err != nil {
		gf_err := gf_core.ErrorCreate("failed to classify image as NSFV or not, using the 'nude' package",
			"verify__invalid_image_nsfv_error",
			map[string]interface{}{"img_path_str":p_img_path_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return true, gf_err
	}
	p_runtime_sys.LogFun("INFO","image is_nude - "+fmt.Sprint(is_nude_bool))
	return is_nude_bool, nil
}

//--------------------------------------------------
func image__flag_as_nsfv(p_image *Gf_crawler_page_image,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {

	err := p_runtime_sys.Mongodb_db.C("gf_crawl").Update(bson.M{
			"t":      "crawler_page_img",
			"id_str": p_image.Id_str,
		},
		bson.M{
			"$set":bson.M{"nsfv_bool":true},
		})
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update an crawler_page_img NSFV flag by its ID",
			"mongodb_update_error",
			map[string]interface{}{"image_id_str": p_image.Id_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}

	return nil
}