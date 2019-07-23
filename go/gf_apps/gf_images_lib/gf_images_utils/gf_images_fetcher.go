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

package gf_images_utils

import (
	"fmt"
	"os"
	"io"
	"math/rand"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type Image_fetch__error struct {
	Id                   bson.ObjectId `json:"-"                    bson:"_id,omitempty"`
	Id_str               string        `json:"id_str"               bson:"id_str"` 
	T_str                string        `json:"-"                    bson:"t"` //img_fetch_error
	Creation_unix_time_f float64       `json:"creation_unix_time_f" bson:"creation_unix_time_f"`
	Image_url_str        string        `json:"image_url_str"        bson:"image_url_str"`
	Status_code_int      int           `json:"status_code_int"      bson:"status_code_int"`
}

//---------------------------------------------------
func Fetcher__get_extern_image(p_image_url_str string,
	p_images_store_local_dir_path_str string,
	p_random_time_delay_bool          bool,
	p_runtime_sys                     *gf_core.Runtime_sys) (string, Gf_image_id, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_fetcher.Fetcher__get_extern_image()")

	if p_random_time_delay_bool {
		//FIX!! - have a sleep time per domain, so that there"s no wait time if the next image
		//        thats processed comes from a different domain
		//        (store these time counters per domain in redis)

		rand.Seed(42)
		max_time_to_sleep_sec_int := 3
		sleep_sec_int             := rand.Intn(max_time_to_sleep_sec_int)
		time.Sleep(time.Second * time.Duration(sleep_sec_int))
	}

	//--------------
	//NEW_IMAGE_LOCAL_FILE_PATH

	//IMPORTANT!! - 0.4 system, image naming, new scheme containing image_id,
	//              instead of the old original_image naming scheme.
	new_image_local_file_path_str, image_id_str, gf_err := Image_ID__create_gf_image_file_path_from_url(p_image_url_str, p_images_store_local_dir_path_str, p_runtime_sys)
	if gf_err != nil {
		return "", "", gf_err
	}	
	//--------------
	//HTTP DOWNLOAD

	gf_err = Download_file(p_image_url_str, new_image_local_file_path_str, p_runtime_sys)
	if gf_err != nil {
		return "", "", gf_err
	}
	//--------------

	//LOG
	analytics__log_image_fetch(p_image_url_str, p_runtime_sys)
	
	//check if local file exists
	if _, err := os.Stat(new_image_local_file_path_str); os.IsNotExist(err) {
		gf_err := gf_core.Error__create("file that was just fetched by the image fetcher doesnt exist in the FS",
			"file_missing_error",
			map[string]interface{}{"new_image_local_file_path_str": new_image_local_file_path_str,},
			err, "gf_images_utils", p_runtime_sys)
		return "", "", gf_err
	}

	return new_image_local_file_path_str, image_id_str, nil
}

//---------------------------------------------------
func analytics__log_image_fetch(p_image_url_str string,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_fetcher.analytics__log_image_fetch()")
}

//---------------------------------------------------
func Download_file(p_image_url_str string,
	p_local_image_file_path_str string,
	p_runtime_sys               *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_fetcher.Download_file()")

	//-----------------------
	gf_http_fetch, gf_err := gf_core.HTTP__fetch_url(p_image_url_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	defer gf_http_fetch.Resp.Body.Close()
	//-----------------------
	//STATUS_CODE CHECK

	//IMPORTANT!! - check if the reponse is as expected
	//				"2" - 2xx - response success
	//              "3" - 3xx - response redirection
	if !(gf_http_fetch.Status_code_int >= 200 && gf_http_fetch.Status_code_int < 400) {

		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := "image_fetch_error__"+fmt.Sprint(creation_unix_time_f)

		fetch_error := &Image_fetch__error{
			Id_str:               id_str,
			T_str:                "img_fetch_error",
			Creation_unix_time_f: creation_unix_time_f,
			Image_url_str:        p_image_url_str,
			Status_code_int:      gf_http_fetch.Status_code_int,
		}

		err := p_runtime_sys.Mongodb_coll.Insert(fetch_error)
		if err != nil {
			gf_err := gf_core.Mongo__handle_error("failed to insert a Image_fetch__error into mongodb",
				"mongodb_insert_error",
				map[string]interface{}{
					"image_url_str":             p_image_url_str,
					"local_image_file_path_str": p_local_image_file_path_str,
				},
				err,"gf_images_utils",p_runtime_sys)
			return gf_err
		}

		gf_err := gf_core.Error__create("image fetching failed with HTTP status error",
			"http_client_req_status_error",
			map[string]interface{}{
				"image_url_str":             p_image_url_str,
				"local_image_file_path_str": p_local_image_file_path_str,
				"status_code_int":           gf_http_fetch.Status_code_int,
			},
			nil, "gf_images_utils", p_runtime_sys)
		return gf_err
	}
	//-----------------------

	final_url_str := gf_http_fetch.Resp.Request.URL.String() //after possible redirects, this is the url
	p_runtime_sys.Log_fun("INFO", "final_url_str - "+final_url_str)

	//--------------
	//WRITE TO FILE
	fmt.Printf("p_local_image_file_path_str - %s\n", p_local_image_file_path_str)

	out, c_err := os.Create(p_local_image_file_path_str)
	defer out.Close()

	if c_err != nil {
		gf_err := gf_core.Error__create("failed to create local file for fetched image",
			"file_create_error",
			map[string]interface{}{"local_image_file_path_str": p_local_image_file_path_str,},
			c_err, "gf_images_utils", p_runtime_sys)
		return gf_err
	}

	_, cp_err := io.Copy(out,gf_http_fetch.Resp.Body)
	if cp_err != nil {
		gf_err := gf_core.Error__create("failed to copy HTTP GET response Body buffer to a image file",
			"file_buffer_copy_error",
			map[string]interface{}{
				"local_image_file_path_str": p_local_image_file_path_str,
				"image_url_str":             p_image_url_str,
			},
			cp_err, "gf_images_utils", p_runtime_sys)
		return gf_err
	}
	//--------------

	return nil
}