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
	// "fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
)

//---------------------------------------------------

type GFimage struct {
	Id                   primitive.ObjectID `json:"-"               bson:"_id,omitempty"`
	IDstr                GFimageID     `json:"id_str"               bson:"id_str"`
	T_str                string        `json:"-"                    bson:"t"` // "img"
	Creation_unix_time_f float64       `json:"creation_unix_time_f" bson:"creation_unix_time_f"`
	UserID               gf_core.GF_ID `json:"user_id_str"          bson:"user_id_str"`

	//---------------
	ClientTypeStr        string        `json:"-"                    bson:"client_type_str"` // "gchrome_ext"|"gf_crawl_images"|"gf_image_editor"
	TitleStr             string        `json:"title_str"            bson:"title_str"`
	FlowsNamesLst        []string      `json:"flows_names_lst"      bson:"flows_names_lst"` // image can bellong to multiple flows

	//---------------
	/*
	RESOLVED_SOURCE_URL
	IMPORTANT!! - when the image comes from an external url (as oppose to it being 
		created internally, or uploaded directly to the system).
		this is different from Origin_page_url_str in that the page_url is the url 
		of the page in which the image is found, whereas this origin_url is the url
		of the file on some file server from which the image is served
	*/
	Origin_url_str string `json:"origin_url_str" bson:"origin_url_str"`

	// if the image is extracted from a page, this holds the page_url
	Origin_page_url_str string `json:"origin_page_url_str,omitempty" bson:"origin_page_url_str,omitempty"`

	/*
	DEPRECATED!! - is this used? images are stored in S3, and accessible via URL.
		actual path on the OS filesystem, of the fullsized image gotten from origin_url_str durring
		processing (download/transformation/s3_upload).
	*/
	Original_file_internal_uri_str string `json:"original_file_internal_uri_str,omitempty" bson:"original_file_internal_uri_str,omitempty"`

	//---------------
	// relative url"s - "/images/image_name.*"

	ThumbnailSmallURLstr  string `json:"thumbnail_small_url_str"  bson:"thumbnail_small_url_str"`
	ThumbnailMediumURLstr string `json:"thumbnail_medium_url_str" bson:"thumbnail_medium_url_str"`
	ThumbnailLargeURLstr  string `json:"thumbnail_large_url_str"  bson:"thumbnail_large_url_str"`
	
	//---------------
	Format_str string `json:"format_str" bson:"format_str"` // "jpeg" | "png" | "gif"
	Width_int  int    `json:"width_str"  bson:"width_int"`
	Height_int int    `json:"height_str" bson:"height_int"`

	//---------------
	// COLORS
	DominantColorHexStr string `json:"dominant_color_hex_str" bson:"dominant_color_hex_str"`
	PalleteStr          string `json:"pallete_str"            bson:"pallete_str"`

	//---------------
	// META
	MetaMap map[string]interface{} `json:"meta_map" bson:"meta_map"` // metadata external users might assign to an image
	TagsLst []string               `json:"tags_lst" bson:"tags_lst"` // human facing tags assigned to an image

	//---------------
}

type GFimageExport struct {
	Creation_unix_time_f  float64  `json:"creation_unix_time_f"`
	UserNameStr           gf_identity_core.GFuserName `json:"user_name_str"`
	Title_str             string   `json:"title_str"`
	Flows_names_lst       []string `json:"flows_names_lst"`
	Origin_page_url_str   string   `json:"origin_page_url_str"`
	ThumbnailSmallURLstr  string   `json:"thumbnail_small_url_str"`
	ThumbnailMediumURLstr string   `json:"thumbnail_medium_url_str"`
	ThumbnailLargeURLstr  string   `json:"thumbnail_large_url_str"`
	Format_str            string   `json:"format_str"`
	Tags_lst              []string `json:"tags_lst"`
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
	IDstr                          GFimageID
	Title_str                      string
	Flows_names_lst                []string
	Image_client_type_str          string
	Origin_url_str                 string
	Origin_page_url_str            string
	Original_file_internal_uri_str string
	ThumbnailSmallURLstr           string
	ThumbnailMediumURLstr          string
	ThumbnailLargeURLstr           string
	Format_str                     string
	Width_int                      int
	Height_int                     int

	Meta_map map[string]interface{}

	// user that owns this image, that uploaded or added it in some other way
	UserID gf_core.GF_ID
}

//---------------------------------------------------

func ImageCreateNew(pImageInfo *GFimageNewInfo,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	image := &GFimage{
		IDstr:                          pImageInfo.IDstr,
		T_str:                          "img",
		Creation_unix_time_f:           creationUNIXtimeF,
		UserID:                         pImageInfo.UserID,
		ClientTypeStr:                  pImageInfo.Image_client_type_str,
		TitleStr:                       pImageInfo.Title_str,
		FlowsNamesLst:                  pImageInfo.Flows_names_lst,
		Origin_url_str:                 pImageInfo.Origin_url_str,
		Origin_page_url_str:            pImageInfo.Origin_page_url_str,
		Original_file_internal_uri_str: pImageInfo.Original_file_internal_uri_str,
		ThumbnailSmallURLstr:           pImageInfo.ThumbnailSmallURLstr,
		ThumbnailMediumURLstr:          pImageInfo.ThumbnailMediumURLstr,
		ThumbnailLargeURLstr:           pImageInfo.ThumbnailLargeURLstr,
		Format_str:                     pImageInfo.Format_str,
		Width_int:                      pImageInfo.Width_int,
		Height_int:                     pImageInfo.Height_int,

		TagsLst: []string{},
		MetaMap: pImageInfo.Meta_map,
	}

	//----------------------------------
	// DB PERSIST

	/*
	// MONGO
	gfErr := DBmongoPutImage(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	*/
	
	// SQL
	gfErr = DBsqlPutImage(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//----------------------------------

	return image, nil
}

//---------------------------------------------------
// DEPRECATED!! - use Image__create_new and its structured input

/*
func ImageCreate(pImageInfoMap map[string]interface{},
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {
	
	newImageInfoMap, gfErr := VerifyImageInfo(pImageInfoMap, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	title_str       := newImageInfoMap["title_str"].(string)
	flows_names_lst := newImageInfoMap["flows_names_lst"].([]string)
	gf_image_id_str := GFimageID(newImageInfoMap["id_str"].(string))

	image := &GFimage{
		IDstr:                          gf_image_id_str,
		T_str:                          "img",
		Creation_unix_time_f:           float64(time.Now().UnixNano())/1000000000.0,
		ClientTypeStr:                  newImageInfoMap["image_client_type_str"].(string),
		TitleStr:                       title_str,
		FlowsNamesLst:                  flows_names_lst,
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

	gfErr = DBputImage(image, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//----------------------------------

	return image, nil
}
*/