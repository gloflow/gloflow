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

package gf_images_core

import (
	"fmt"
	"os"
	"io"
	"math/rand"
	"time"
	"context"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type ImageFetchError struct {
	Id                   primitive.ObjectID `json:"-"                    bson:"_id,omitempty"`
	Id_str               string        `json:"id_str"               bson:"id_str"` 
	T_str                string        `json:"-"                    bson:"t"` // img_fetch_error
	Creation_unix_time_f float64       `json:"creation_unix_time_f" bson:"creation_unix_time_f"`
	Image_url_str        string        `json:"image_url_str"        bson:"image_url_str"`
	Status_code_int      int           `json:"status_code_int"      bson:"status_code_int"`
}

//---------------------------------------------------

func FetcherGetExternImage(pImageURLstr string,
	pImagesStoreLocalDirPathStr string,
	pRandomTimeDelayBool        bool,
	pRuntimeSys                 *gf_core.RuntimeSys) (string, GFimageID, *gf_core.GFerror) {

	if pRandomTimeDelayBool {

		// FIX!! - have a sleep time per domain, so that there"s no wait time if the next image
		//         thats processed comes from a different domain
		//         (store these time counters per domain in redis)

		creation_unix_time := time.Now().UnixNano()
		rand.Seed(creation_unix_time)
		max_time_to_sleep_sec_int := 3
		sleep_sec_int             := rand.Intn(max_time_to_sleep_sec_int)
		time.Sleep(time.Second * time.Duration(sleep_sec_int))
	}

	//--------------
	// NEW_IMAGE_LOCAL_FILE_PATH

	newImageLocalFilePathStr, imageIDstr, gfErr := CreateImageFilePathFromURL("",
		pImageURLstr,
		pImagesStoreLocalDirPathStr,
		pRuntimeSys)
	if gfErr != nil {
		return "", "", gfErr
	}

	//--------------
	// HTTP DOWNLOAD

	gfErr = DownloadFile(pImageURLstr, newImageLocalFilePathStr, pRuntimeSys)
	if gfErr != nil {
		return "", "", gfErr
	}
	
	//--------------

	// LOG
	analytics__log_image_fetch(pImageURLstr, pRuntimeSys)
	
	// check if local file exists
	if _, err := os.Stat(newImageLocalFilePathStr); os.IsNotExist(err) {
		gfErr := gf_core.ErrorCreate("file that was just fetched by the image fetcher doesnt exist in the FS",
			"file_missing_error",
			map[string]interface{}{"new_image_local_file_path_str": newImageLocalFilePathStr,},
			err, "gf_images_core", pRuntimeSys)
		return "", "", gfErr
	}

	return newImageLocalFilePathStr, imageIDstr, nil
}

//---------------------------------------------------

func analytics__log_image_fetch(pImageURLstr string,
	pRuntimeSys *gf_core.RuntimeSys) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_fetcher.analytics__log_image_fetch()")
}

//---------------------------------------------------

func DownloadFile(pImageURLstr string,
	p_local_image_file_path_str string,
	pRuntimeSys               *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_fetcher.DownloadFile()")

	//-----------------------
	headersMap, userAgentStr := GetHTTPreqConfig()
	ctx := context.Background()

	HTTPfetch, gfErr := gf_core.HTTPfetchURL(pImageURLstr, headersMap, userAgentStr, ctx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	defer HTTPfetch.Resp.Body.Close()
	
	//-----------------------
	// STATUS_CODE CHECK

	// IMPORTANT!! - check if the reponse is as expected
	// 	  			 "2" - 2xx - response success
	//               "3" - 3xx - response redirection
	if !(HTTPfetch.Status_code_int >= 200 && HTTPfetch.Status_code_int < 400) {

		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := "image_fetch_error__"+fmt.Sprint(creation_unix_time_f)

		fetch_error := &ImageFetchError{
			Id_str:               id_str,
			T_str:                "img_fetch_error",
			Creation_unix_time_f: creation_unix_time_f,
			Image_url_str:        pImageURLstr,
			Status_code_int:      HTTPfetch.Status_code_int,
		}

		ctx := context.Background()
		collNameStr := pRuntimeSys.Mongo_coll.Name()
		gfErr := gf_core.MongoInsert(fetch_error,
			collNameStr,
			map[string]interface{}{
				"image_url_str":             pImageURLstr,
				"local_image_file_path_str": p_local_image_file_path_str,
				"caller_err_msg_str":        "failed to insert a image fetch error into the DB",
			},
			ctx,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		gfErr = gf_core.ErrorCreate("image fetching failed with HTTP status error",
			"http_client_req_status_error",
			map[string]interface{}{
				"image_url_str":             pImageURLstr,
				"local_image_file_path_str": p_local_image_file_path_str,
				"status_code_int":           HTTPfetch.Status_code_int,
			},
			nil, "gf_images_core", pRuntimeSys)
		return gfErr
	}

	//-----------------------

	finalURLstr := HTTPfetch.Resp.Request.URL.String() // after possible redirects, this is the url
	pRuntimeSys.LogFun("INFO", fmt.Sprintf("final_url_str - %s", finalURLstr))

	//--------------
	// WRITE TO FILE
	fmt.Printf("p_local_image_file_path_str - %s\n", p_local_image_file_path_str)

	out, c_err := os.Create(p_local_image_file_path_str)
	defer out.Close()

	if c_err != nil {
		gfErr := gf_core.ErrorCreate("failed to create local file for fetched image",
			"file_create_error",
			map[string]interface{}{"local_image_file_path_str": p_local_image_file_path_str,},
			c_err, "gf_images_core", pRuntimeSys)
		return gfErr
	}

	_, cp_err := io.Copy(out, HTTPfetch.Resp.Body)
	if cp_err != nil {
		gfErr := gf_core.ErrorCreate("failed to copy HTTP GET response Body buffer to a image file",
			"file_buffer_copy_error",
			map[string]interface{}{
				"local_image_file_path_str": p_local_image_file_path_str,
				"image_url_str":             pImageURLstr,
			},
			cp_err, "gf_images_core", pRuntimeSys)
		return gfErr
	}
	
	//--------------

	return nil
}

//---------------------------------------------------

func GetHTTPreqConfig() (map[string]string, string) {
	headersMap   := map[string]string{}
	userAgentStr := "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1"
	return headersMap, userAgentStr
}