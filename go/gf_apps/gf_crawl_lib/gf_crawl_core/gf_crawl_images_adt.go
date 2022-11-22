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
	"fmt"
	"time"
	"crypto/md5"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_utils"
)

//---------------------------------------------------

type GFcrawlerPageImageID string
type GFcrawlerPageImage struct {
	Id                         primitive.ObjectID       `bson:"_id,omitempty"`
	IDstr                      GFcrawlerPageImageID     `bson:"id_str"`
	T_str                      string                   `bson:"t"`                          // "crawler_page_img"
	Creation_unix_time_f       float64                  `bson:"creation_unix_time_f"`
	Crawler_name_str           string                   `bson:"crawler_name_str"`           // name of the crawler that discovered this image
	Cycle_run_id_str           string                   `bson:"cycle_run_id_str"`
	Img_ext_str                string                   `bson:"img_ext_str"`                // jpg|gif|png
	Url_str                    string                   `bson:"url_str"`
	Domain_str                 string                   `bson:"domain_str"`                 // domain of the url_str
	Origin_page_url_str        string                   `bson:"origin_page_url_str"`        // page url from whos html this element was extracted
	Origin_page_url_domain_str string                   `bson:"origin_page_url_domain_str"` // domain of the origin_page_url_str // NEW_FIELD!! a lot of records dont have this field

	// IMPORTANT!! - this is unique for the image src encountered. this way the same data links are not entered in duplicates, 
	//               and using the hash the DB can qucikly be checked for existence of record
	Hash_str                   string        `bson:"hash_str"`

	// IMPORTANT!! - indicates if the image was fetched from the remote server,
	//               and has been stored on S3 and ready for usage by other services. 
	Downloaded_bool            bool          `bson:"downloaded_bool"`

	// IMPORTANT!! - the usage was determined to be useful for internal applications,
	//               they're not page elements, or other small unimportant parts.
	//               if it is valid for usage then a gf_image for this image should be 
	//               found in the db
	Valid_for_usage_bool       bool                     `bson:"valid_for_usage_bool"`
	S3_stored_bool             bool                     `bson:"s3_stored_bool"` // if persisting to s3 succeeded
	Nsfv_bool                  bool                     `bson:"nsfv_bool"`      // NSFV (not safe for viewing/nudity) flag for the image 
	GFimageIDstr               gf_images_core.GFimageID `bson:"image_id_str"`   // id of the gf_image for this corresponding crawler_page_img //FIX!! - should be "gf_image_id_str"
}
	
// IMPORTANT!! - reference to an image, on a particular page. 
//               the same image, with the same Url_str can appear on multiple pages, and this 
//               struct tracks that, one record per reference
type GFcrawlerPageImageRef struct {
	Id                         primitive.ObjectID `bson:"_id,omitempty"`
	Id_str                     string        `bson:"id_str"`
	T_str                      string        `bson:"t"`                          //"crawler_page_img_ref"
	Creation_unix_time_f       float64       `bson:"creation_unix_time_f"`
	Crawler_name_str           string        `bson:"crawler_name_str"`           //name of the crawler that discovered this image
	Cycle_run_id_str           string        `bson:"cycle_run_id_str" json:"cycle_run_id_str"`
	Url_str                    string        `bson:"url_str"`
	Domain_str                 string        `bson:"domain_str"`
	Origin_page_url_str        string        `bson:"origin_page_url_str"`        //page url from whos html this element was extracted
	Origin_page_url_domain_str string        `bson:"origin_page_url_domain_str"` //NEW_FIELD!! a lot of records dont have this field

	// IMPORTANT!! - this is unique for the image src encountered. this way the same data links are not entered in duplicates, 
	//               and using the hash the DB can qucikly be checked for existence of record
	Hash_str                   string        `bson:"hash_str"`
}

type GFcrawlerRecentImages struct {
	Domain_str               string    `bson:"_id"                      json:"domain_str"`
	Imgs_count_int           int       `bson:"imgs_count_int"           json:"imgs_count_int"`
	Crawler_page_img_ids_lst []string  `bson:"crawler_page_img_ids_lst" json:"crawler_page_img_ids_lst"`
	Creation_times_lst       []float64 `bson:"creation_times_lst"       json:"creation_times_lst"`
	Urls_lst                 []string  `bson:"urls_lst"                 json:"urls_lst"`
	Nsfv_lst                 []bool    `bson:"nsfv_lst"                 json:"nsfv_lst"`
	Origin_page_urls_lst     []string  `bson:"origin_page_urls_lst"     json:"origin_page_urls_lst"`
}

//---------------------------------------------------

func imagesADTprepareAndCreate(pCrawlerNameStr string,
	pCycleRunIDstr        string,
	p_img_src_url_str     string,
	p_origin_page_url_str string,
	pRuntime              *GFcrawlerRuntime,
	pRuntimeSys           *gf_core.RuntimeSys) (*GFcrawlerPageImage, *gf_core.GFerror) {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//------------------
	// DOMAINS

	imgSrcDomainStr, origin_page_url_domain_str, gfErr := gf_crawl_utils.GetDomain(p_img_src_url_str, p_origin_page_url_str, pRuntimeSys)
	if gfErr != nil {
		t := "images_in_page__get_domain__failed"
		m := "failed to get domain of image with img_src - "+p_img_src_url_str
		CreateErrorAndEvent(t, m, map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,}, p_img_src_url_str, pCrawlerNameStr,
			gfErr, pRuntime, pRuntimeSys)
		return nil, gfErr
	}

	//-------------
	// COMPLETE_A_HREF
	
	complete_img_src_url_str, gfErr := gf_crawl_utils.CompleteURL(p_img_src_url_str, imgSrcDomainStr, pRuntimeSys)
	if gfErr != nil {
		t:="complete_url__failed"
		m:="failed to complete_url of image with img_src - "+p_img_src_url_str
		CreateErrorAndEvent(t, m, map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,}, p_img_src_url_str, pCrawlerNameStr,
			gfErr, pRuntime, pRuntimeSys)
		return nil, gfErr
	}

	//-------------
	// GET_IMG_EXT_FROM_URL

	img_ext_str, gfErr := gf_images_core.GetImageExtFromURL(p_img_src_url_str, pRuntimeSys)
	if gfErr != nil {
		t:="images_in_page__get_img_extension__failed"
		m:="failed to get file extension of image with img_src - "+p_img_src_url_str
		CreateErrorAndEvent(t, m, map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,}, p_img_src_url_str, pCrawlerNameStr,
			gfErr, pRuntime, pRuntimeSys)
		return nil, gfErr
	}

	//-------------
	pRuntimeSys.LogFun("INFO",">>>>> "+cyan("img")+" -- "+yellow(imgSrcDomainStr)+" ------ "+yellow(fmt.Sprint(complete_img_src_url_str)))

	img := imagesADTcreate(pCrawlerNameStr,
		pCycleRunIDstr,
		complete_img_src_url_str,
		img_ext_str,
		imgSrcDomainStr,
		p_origin_page_url_str,
		origin_page_url_domain_str,
		pRuntimeSys)
	return img, nil
}

//---------------------------------------------------

func imagesADTcreateID() (GFcrawlerPageImageID, float64) {
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("crawler_page_img:%f", creation_unix_time_f)
	return GFcrawlerPageImageID(id_str), creation_unix_time_f
}

//---------------------------------------------------

func imagesADTcreate(pCrawlerNameStr string,
	pCycleRunIDstr               string,
	p_img_src_url_str            string,
	p_img_ext_str                string,
	p_imgSrcDomainStr            string,
	p_origin_page_url_str        string,
	p_origin_page_url_domain_str string,
	pRuntimeSys                  *gf_core.RuntimeSys) *GFcrawlerPageImage {

	
	// HASH
	to_hash_str := p_img_src_url_str // one Crawler_page_img for a given page url, no matter on how many pages it is referenced by
	hash        := md5.New()
	hash.Write([]byte(to_hash_str))
	hash_str := hex.EncodeToString(hash.Sum(nil))

	id_str, creation_unix_time_f := imagesADTcreateID()
	img := &GFcrawlerPageImage{
		IDstr:                      id_str,
		T_str:                      "crawler_page_img",
		Creation_unix_time_f:       creation_unix_time_f,
		Crawler_name_str:           pCrawlerNameStr,
		Cycle_run_id_str:           pCycleRunIDstr,
		Img_ext_str:                p_img_ext_str,
		Url_str:                    p_img_src_url_str,
		Domain_str:                 p_imgSrcDomainStr,
		Origin_page_url_str:        p_origin_page_url_str,
		Origin_page_url_domain_str: p_origin_page_url_domain_str,
		Hash_str:                   hash_str,
		Downloaded_bool:            false,
		Valid_for_usage_bool:       false, // all images are initially set as invalid for usage
		S3_stored_bool:             false, 
	}
	return img
}

//---------------------------------------------------

func imagesADTrefCreate(pCrawlerNameStr string,
	pCycleRunIDstr               string,
	p_image_url_str              string,
	p_image_url_domain_str       string,
	p_origin_page_url_str        string,
	p_origin_page_url_domain_str string,
	pRuntimeSys                  *gf_core.RuntimeSys) *GFcrawlerPageImageRef {

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	ref_id_str           := fmt.Sprintf("img_ref:%f", creation_unix_time_f)

	// HASH
	// IMPORTANT!! - one Crawler_page_img_ref per page img reference, so if the same image is linked on several pages
	//               each of those references will have a different hash_str and will created a new Crawler_page_img_ref
	to_hash_str := p_image_url_str+p_origin_page_url_str
	hash        := md5.New()
	hash.Write([]byte(to_hash_str))
	hash_str := hex.EncodeToString(hash.Sum(nil))

	gf_img_ref := &GFcrawlerPageImageRef{
		Id_str:                     ref_id_str,
		T_str:                      "crawler_page_img_ref",
		Creation_unix_time_f:       creation_unix_time_f,
		Crawler_name_str:           pCrawlerNameStr,
		Cycle_run_id_str:           pCycleRunIDstr,
		Url_str:                    p_image_url_str,        // complete_img_src_str,
		Domain_str:                 p_image_url_domain_str, // imgSrcDomainStr,
		Origin_page_url_str:        p_origin_page_url_str,
		Origin_page_url_domain_str: p_origin_page_url_domain_str,
		Hash_str:                   hash_str,
	}

	return gf_img_ref
}