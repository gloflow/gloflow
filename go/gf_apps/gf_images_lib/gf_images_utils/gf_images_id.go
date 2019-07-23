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
	"net/url"
	"strings"
	"time"
	"path/filepath"
	"crypto/md5"
	"encoding/hex"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type Gf_image_id string

//---------------------------------------------------
//CREATES_ID
func Image_ID__create_from_url(p_image_url_str string, p_runtime_sys *gf_core.Runtime_sys) (Gf_image_id, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_id.Image_ID__create_from_url()")
	
	//urlparse() - used so that any possible url query parameters are not used in the 
	//             os.path.basename() result
	url,err := url.Parse(p_image_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse image_url to create image ID",
			"url_parse_error",
			map[string]interface{}{"image_url_str": p_image_url_str,},
			err, "gf_images_utils", p_runtime_sys)
		return "", gf_err
	}
	
	image_path_str      := url.Path
	image_ext_str       := filepath.Ext(image_path_str)
	clean_image_ext_str := strings.Trim(strings.ToLower(image_ext_str),".")
	//image_file_name_str := path.Base(image_path_str)
	//image_ext_str       := strings.Split(image_file_name_str,".")[1]

	normalized_ext_str,ok := Image__check_image_format(clean_image_ext_str, p_runtime_sys)
	if !ok {
		usr_msg_str := "invalid image extension found in image url - "+p_image_url_str

		gf_err := gf_core.Error__create(usr_msg_str,
			"verify__invalid_image_extension_error",
			map[string]interface{}{"image_url_str":p_image_url_str,},
			err, "gf_images_utils", p_runtime_sys)
		return "", gf_err
	}
	//-------------
	gf_image_id_str := Image_ID__create(image_path_str, normalized_ext_str, p_runtime_sys)
	//-------------
	return gf_image_id_str, nil
}

//---------------------------------------------------
//CREATES_ID
func Image_ID__create_gf_image_file_path_from_url(p_image_url_str string,
	p_images_store_local_dir_path_str string,
	p_runtime_sys                     *gf_core.Runtime_sys) (string, Gf_image_id, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_id.Image_ID__create_gf_image_file_path_from_url()")

	//IMPORTANT!! - 0.4 system, image naming, new scheme containing image_id,
	//              instead of the old original_image naming scheme.
	gf_image_id_str, _ := Image_ID__create_from_url(p_image_url_str, p_runtime_sys)
	ext_str, gf_err := Get_image_ext_from_url(p_image_url_str, p_runtime_sys)
	if gf_err != nil {
		return "", "", gf_err
	}

	local_image_file_name_str := fmt.Sprintf("%s.%s", gf_image_id_str, ext_str)
	local_image_file_path_str := fmt.Sprintf("%s/%s", p_images_store_local_dir_path_str, local_image_file_name_str)

	p_runtime_sys.Log_fun("INFO", fmt.Sprintf("local_image_file_path_str - %s", local_image_file_path_str))
	
	return local_image_file_path_str, gf_image_id_str, nil
}

//---------------------------------------------------
//CREATES_ID
//p_image_type_str - :String - "jpeg"|"gif"|"png"

func Image_ID__create(p_image_path_str string,
	p_image_format_str string,
	p_runtime_sys      *gf_core.Runtime_sys) Gf_image_id {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_id.Image_ID__create()")
	
	h := md5.New()
	
	h.Write([]byte(fmt.Sprint(float64(time.Now().UnixNano())/1000000000.0)))
	h.Write([]byte(p_image_path_str))
	h.Write([]byte(p_image_format_str))

	sum     := h.Sum(nil)
	hex_str := hex.EncodeToString(sum)
	
	gf_image_id_str := Gf_image_id(hex_str)
	return gf_image_id_str
}