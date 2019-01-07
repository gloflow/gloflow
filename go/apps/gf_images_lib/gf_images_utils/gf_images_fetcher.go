package gf_images_utils

import (
	"fmt"
	"os"
	"io"
	"math/rand"
	"time"
	"github.com/globalsign/mgo/bson"
	"gf_core"
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
//-------------------------------------------------
func Fetch_image(p_image_url_str string,
		p_images_store_local_dir_path_str string,
		p_runtime_sys                     *gf_core.Runtime_sys) (string,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_fetcher.Fetch_image()")

	//----------------------
	local_image_file_path_str,gf_err := Fetcher__get_extern_image(p_image_url_str,
															p_images_store_local_dir_path_str,
															true,
															p_runtime_sys)
	if gf_err != nil {
		return "",gf_err
	}

	//check if local file exists
	if _, err := os.Stat(local_image_file_path_str); os.IsNotExist(err) {
		gf_err := gf_core.Error__create("file that was just fetched by the image fetcher doesnt exist in the FS",
			"file_missing_error",
			&map[string]interface{}{"local_image_file_path_str":local_image_file_path_str,},
			err,"gf_images_utils",p_runtime_sys)
		return "",gf_err
	}
	//----------------------

	return local_image_file_path_str,nil
}
//---------------------------------------------------
func Fetcher__get_extern_image(p_image_url_str string,
				p_images_store_local_dir_path_str string,
				p_random_time_delay_bool          bool,
				p_runtime_sys                     *gf_core.Runtime_sys) (string,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_fetcher.Fetcher__get_extern_image()")

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
	//FS FILE

	/*parsed_url,err := url.Parse(*p_image_url_str)
	if err != nil {
		return nil,errors.New("supplied image_url is not valid - "+*p_image_url_str)
	}

	url_path_str := parsed_url.Path

	//ATTENTION!! - using the full image_path_str to calc file name, instead of the exact image_file_name_str
	//              because on lots of sites images may have the same name, but are differentiated by the admins/authors
	//              by putting them in different folders (paths)... so to avoid this problem Im acounting on the 
	//              full image path being unique across the whole internal image DB. 
	//              also, by using the image_path_str as the file name, I can get uniqueness and also be able to 
	//              extract the actual file name in case the Image DB is lost and all I have is the file names
	image_file_name_str       := strings.TrimLeft(url_path_str,"/") //path.Base(url_path_str)*/

	//IMPORTANT!! - 0.4 system, image naming, new scheme containing image_id,
	//              instead of the old original_image naming scheme
	image_id_str,_ := Image__create_id_from_url(p_image_url_str, p_runtime_sys)
	ext_str,gf_err := Get_image_ext_from_url(p_image_url_str, p_runtime_sys)
	if gf_err != nil {
		return "",gf_err
	}

	local_image_file_path_str := fmt.Sprintf("%s/%s.%s",
										p_images_store_local_dir_path_str,
										image_id_str,
										ext_str)

	p_runtime_sys.Log_fun("INFO","local_image_file_path_str - "+local_image_file_path_str)
	//--------------
	//HTTP DOWNLOAD

	gf_err = Download_file(p_image_url_str,
					local_image_file_path_str,
					p_runtime_sys)
	if gf_err != nil {
		return "",gf_err
	}
	//--------------
	//LOG
	analytics__log_image_fetch(p_image_url_str, p_runtime_sys)
	
	return local_image_file_path_str,nil
}
//---------------------------------------------------
func analytics__log_image_fetch(p_image_url_str string,
					p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_fetcher.analytics__log_image_fetch()")
}
//---------------------------------------------------
func Download_file(p_image_url_str string,
				p_local_image_file_path_str string,
				p_runtime_sys               *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_fetcher.Download_file()")

	//-----------------------
	gf_http_fetch,gf_err := gf_core.HTTP__fetch_url(p_image_url_str,p_runtime_sys)
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
			Id_str              :id_str,
			T_str               :"img_fetch_error",
			Creation_unix_time_f:creation_unix_time_f,
			Image_url_str       :p_image_url_str,
			Status_code_int     :gf_http_fetch.Status_code_int,
		}

		err := p_runtime_sys.Mongodb_coll.Insert(fetch_error)
		if err != nil {
			gf_err := gf_core.Error__create("failed to insert a Image_fetch__error into mongodb",
				"mongodb_insert_error",
				&map[string]interface{}{
					"image_url_str":            p_image_url_str,
					"local_image_file_path_str":p_local_image_file_path_str,
				},
				err,"gf_images_utils",p_runtime_sys)
			return gf_err
		}

		gf_err := gf_core.Error__create("image fetching failed with HTTP status error",
			"http_client_req_status_error",
			&map[string]interface{}{
				"image_url_str":            p_image_url_str,
				"local_image_file_path_str":p_local_image_file_path_str,
				"status_code_int":          gf_http_fetch.Status_code_int,
			},
			nil,"gf_images_utils",p_runtime_sys)
		return gf_err
	}
	//-----------------------

	final_url_str := gf_http_fetch.Resp.Request.URL.String() //after possible redirects, this is the url
	p_runtime_sys.Log_fun("INFO","final_url_str - "+final_url_str)

	//--------------
	//WRITE TO FILE

	out,c_err := os.Create(p_local_image_file_path_str)
	defer out.Close()

	if c_err != nil {
		gf_err := gf_core.Error__create("failed to create local file for fetched image",
			"file_create_error",
			&map[string]interface{}{"local_image_file_path_str":p_local_image_file_path_str,},
			c_err,"gf_images_utils",p_runtime_sys)
		return gf_err
	}

	_,cp_err := io.Copy(out,gf_http_fetch.Resp.Body)
	if cp_err != nil {
		gf_err := gf_core.Error__create("failed to copy HTTP GET response Body buffer to a image file",
			"file_buffer_copy_error",
			&map[string]interface{}{
				"local_image_file_path_str":p_local_image_file_path_str,
				"image_url_str":            p_image_url_str,
			},
			cp_err,"gf_images_utils",p_runtime_sys)
		return gf_err
	}
	//--------------

	return nil
}