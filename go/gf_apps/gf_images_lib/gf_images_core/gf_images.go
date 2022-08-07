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
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GFimage  = Gf_image
type GF_image = Gf_image
type Gf_image struct {

	Id                   primitive.ObjectID `json:"-"               bson:"_id,omitempty"`
	Id_str               GFimageID     `json:"id_str"               bson:"id_str"`
	T_str                string        `json:"-"                    bson:"t"` // "img"
	Creation_unix_time_f float64       `json:"creation_unix_time_f" bson:"creation_unix_time_f"`
	
	//---------------
	Client_type_str      string        `json:"-"                    bson:"client_type_str"` // "gchrome_ext"|"gf_crawl_images"|"gf_image_editor"
	Title_str            string        `json:"title_str"            bson:"title_str"`
	Flows_names_lst      []string      `json:"flows_names_lst"      bson:"flows_names_lst"` // image can bellong to multiple flows

	//---------------
	// RESOLVED_SOURCE_URL
	// IMPORTANT!! - when the image comes from an external url (as oppose to it being 
	//               created internally, or uploaded directly to the system).
	//               this is different from Origin_page_url_str in that the page_url is the url 
	//               of the page in which the image is found, whereas this origin_url is the url
	//               of the file on some file server from which the image is served
	Origin_url_str string `json:"origin_url_str" bson:"origin_url_str"`

	// if the image is extracted from a page, this holds the page_url
	Origin_page_url_str string `json:"origin_page_url_str,omitempty" bson:"origin_page_url_str,omitempty"`

	// DEPRECATED!! - is this used? images are stored in S3, and accessible via URL.
	// actual path on the OS filesystem, of the fullsized image gotten from origin_url_str durring
	// processing (download/transformation/s3_upload).
	Original_file_internal_uri_str string `json:"original_file_internal_uri_str,omitempty" bson:"original_file_internal_uri_str,omitempty"`

	//---------------
	// relative url"s - "/images/image_name.*"
	Thumbnail_small_url_str  string `json:"thumbnail_small_url_str"  bson:"thumbnail_small_url_str"`
	Thumbnail_medium_url_str string `json:"thumbnail_medium_url_str" bson:"thumbnail_medium_url_str"`
	Thumbnail_large_url_str  string `json:"thumbnail_large_url_str"  bson:"thumbnail_large_url_str"`
	
	//---------------
	Format_str string `json:"format_str" bson:"format_str"` // "jpeg"|"png"|"gif"
	Width_int  int    `json:"width_str"  bson:"width_int"`
	Height_int int    `json:"height_str" bson:"height_int"`

	//---------------
	// COLORS
	Dominant_color_hex_str string `json:"dominant_color_hex_str"`
	Pallete_str            string `json:"pallete_str"`

	//---------------
	// META
	Meta_map map[string]interface{} `json:"meta_map" bson:"meta_map"` // metadata external users might assign to an image
	Tags_lst []string               `json:"tags_lst" bson:"tags_lst"` // human facing tags assigned to an image

	//---------------

	// DEPRECATED!! - all images have the flows_names_lst member now, so flow_name_str can be removed both here from the 
	//                struct and from DB records
	// Flow_name_str   string   `json:"flow_name_str"   bson:"flow_name_str"`
}

type GF_image_export struct {
	Creation_unix_time_f     float64  `json:"creation_unix_time_f"`
	Title_str                string   `json:"title_str"`
	Flows_names_lst          []string `json:"flows_names_lst"`
	Origin_page_url_str      string   `json:"origin_page_url_str"`
	Thumbnail_small_url_str  string   `json:"thumbnail_small_url_str"`
	Thumbnail_medium_url_str string   `json:"thumbnail_medium_url_str"`
	Thumbnail_large_url_str  string   `json:"thumbnail_large_url_str"`
	Format_str               string   `json:"format_str"`
	Tags_lst                 []string `json:"tags_lst"`
}

type GFimageThumbs struct {
	Small_relative_url_str     string `json:"small_relative_url_str"`
	Medium_relative_url_str    string `json:"medium_relative_url_str"`
	Large_relative_url_str     string `json:"large_relative_url_str"`

	Small_local_file_path_str  string
	Medium_local_file_path_str string
	Large_local_file_path_str  string
}

type GFimageNewInfo struct {
	Id_str                         GFimageID
	Title_str                      string
	Flows_names_lst                []string
	Image_client_type_str          string
	Origin_url_str                 string
	Origin_page_url_str            string
	Original_file_internal_uri_str string
	Thumbnail_small_url_str        string
	Thumbnail_medium_url_str       string
	Thumbnail_large_url_str        string
	Format_str                     string
	Width_int                      int
	Height_int                     int

	Meta_map map[string]interface{}
}

//---------------------------------------------------
func ImageCreateNew(pImageInfo *GFimageNewInfo,
	pCtx       context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GF_image, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	image := &GF_image{
		Id_str:                         pImageInfo.Id_str,
		T_str:                          "img",
		Creation_unix_time_f:           creationUNIXtimeF,
		Client_type_str:                pImageInfo.Image_client_type_str,
		Title_str:                      pImageInfo.Title_str,
		Flows_names_lst:                pImageInfo.Flows_names_lst,
		Origin_url_str:                 pImageInfo.Origin_url_str,
		Origin_page_url_str:            pImageInfo.Origin_page_url_str,
		Original_file_internal_uri_str: pImageInfo.Original_file_internal_uri_str,
		Thumbnail_small_url_str:        pImageInfo.Thumbnail_small_url_str,
		Thumbnail_medium_url_str:       pImageInfo.Thumbnail_medium_url_str,
		Thumbnail_large_url_str:        pImageInfo.Thumbnail_large_url_str,
		Format_str:                     pImageInfo.Format_str,
		Width_int:                      pImageInfo.Width_int,
		Height_int:                     pImageInfo.Height_int,

		Meta_map: pImageInfo.Meta_map,
	}

	//----------------------------------
	// DB PERSIST
	gfErr := DB__put_image(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//----------------------------------

	return image,nil
}

//---------------------------------------------------
// DEPRECATED!! - use Image__create_new and its structured input

func Image__create(pImageInfoMap map[string]interface{},
	pCtx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) (*GF_image, *gf_core.GFerror) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images.Image__create()")
	
	newImageInfoMap, gf_err := Image__verify_image_info(pImageInfoMap, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	title_str       := newImageInfoMap["title_str"].(string)
	flows_names_lst := newImageInfoMap["flows_names_lst"].([]string)
	gf_image_id_str := GFimageID(newImageInfoMap["id_str"].(string))

	image := &GF_image{
		Id_str:                         gf_image_id_str,
		T_str:                          "img",
		Creation_unix_time_f:           float64(time.Now().UnixNano())/1000000000.0,
		Client_type_str:                newImageInfoMap["image_client_type_str"].(string),
		Title_str:                      title_str,
		Flows_names_lst:                flows_names_lst,
		Origin_url_str:                 newImageInfoMap["origin_url_str"].(string),
		Origin_page_url_str:            newImageInfoMap["origin_page_url_str"].(string),
		Original_file_internal_uri_str: newImageInfoMap["original_file_internal_uri_str"].(string),
		Thumbnail_small_url_str:        newImageInfoMap["thumbnail_small_url_str"].(string),
		Thumbnail_medium_url_str:       newImageInfoMap["thumbnail_medium_url_str"].(string),
		Thumbnail_large_url_str:        newImageInfoMap["thumbnail_large_url_str"].(string),
		Format_str:                     newImageInfoMap["format_str"].(string),
	}
	
	//----------------------------------
	// DB PERSIST

	db_gf_err := DB__put_image(image, pCtx, p_runtime_sys)
	if db_gf_err != nil {
		return nil, db_gf_err
	}
	
	//----------------------------------

	return image, nil
}