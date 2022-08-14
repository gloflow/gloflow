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
	"net/url"
	"strings"
	"time"
	"path/filepath"
	"crypto/md5"
	"encoding/hex"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type Gf_image_id string // DEPRECATED!! - switch to using GF_image_id fully
type GF_image_id = Gf_image_id
type GFimageID   = GF_image_id

//---------------------------------------------------
// CREATES_ID
func CreateIDfromURL(pImageURLstr string,
	pRuntimeSys *gf_core.RuntimeSys) (GF_image_id, *gf_core.GFerror) {
	
	// urlparse() - used so that any possible url query parameters are not used in the 
	//              os.path.basename() result
	url, err := url.Parse(pImageURLstr)
	if err != nil {
		gfErr := gf_core.Error__create("failed to parse image_url to create image ID",
			"url_parse_error",
			map[string]interface{}{"image_url_str": pImageURLstr,},
			err, "gf_images_core", pRuntimeSys)
		return "", gfErr
	}
	
	imageHostStr     := url.Host
	imagePathStr     := url.Path
	imageURIstr      := fmt.Sprintf("%s/%s", imageHostStr, imagePathStr)
	imageExtStr      := filepath.Ext(imagePathStr)
	cleanImageExtStr := strings.Trim(strings.ToLower(imageExtStr),".")
	// imageFileNameStr := path.Base(imagePathStr)
	// imageExtStr      := strings.Split(imageFileNameStr,".")[1]



	fmt.Println("CCCCCCCCCCCCCCCCCCCCCCCCC", cleanImageExtStr)


	_, ok := CheckImageFormat(cleanImageExtStr, pRuntimeSys)
	if !ok {
		usrMsgStr := "invalid image extension found in image url - "+pImageURLstr

		gfErr := gf_core.Error__create(usrMsgStr,
			"verify__invalid_image_extension_error",
			map[string]interface{}{"image_url_str": pImageURLstr,},
			err, "gf_images_core", pRuntimeSys)
		return "", gfErr
	}
	
	//-------------
	imageIDstr := CreateImageID(imageURIstr,
		pRuntimeSys)
	
	//-------------
	return imageIDstr, nil
}

//---------------------------------------------------
// CREATES_ID
// p_image_type_str - :String - "jpeg"|"gif"|"png"

func CreateImageID(pImageURIstr string,
	pRuntimeSys *gf_core.RuntimeSys) GFimageID {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_images_id.Image_ID__create()")
	
	h := md5.New()
		
	current_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	h.Write([]byte(fmt.Sprint(current_unix_time_f)))
	h.Write([]byte(pImageURIstr))

	sum     := h.Sum(nil)
	hex_str := hex.EncodeToString(sum)
	
	gfImageIDstr := GFimageID(hex_str)
	return gfImageIDstr
}