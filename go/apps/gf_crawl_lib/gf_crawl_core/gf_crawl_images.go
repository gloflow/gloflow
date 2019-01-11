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

package gf_crawl_core

import (
	"fmt"
	"time"
	"crypto/md5"
	"encoding/hex"
	"github.com/globalsign/mgo/bson"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/apps/gf_crawl_lib/gf_crawl_utils"
)
//--------------------------------------------------
type Crawler_page_img struct {
	Id                         bson.ObjectId `bson:"_id,omitempty"`
	Id_str                     string        `bson:"id_str"`
	T_str                      string        `bson:"t"`                          //"crawler_page_img"
	Creation_unix_time_f       float64       `bson:"creation_unix_time_f"`
	Crawler_name_str           string        `bson:"crawler_name_str"`           //name of the crawler that discovered this image
	Cycle_run_id_str           string        `bson:"cycle_run_id_str"`
	Img_ext_str                string        `bson:"img_ext_str"`                //jpg|gif|png
	Url_str                    string        `bson:"url_str"`
	Domain_str                 string        `bson:"domain_str"`                 //domain of the url_str
	Origin_page_url_str        string        `bson:"origin_page_url_str"`        //page url from whos html this element was extracted
	Origin_page_url_domain_str string        `bson:"origin_page_url_domain_str"` //domain of the origin_page_url_str //NEW_FIELD!! a lot of records dont have this field

	//IMPORTANT!! - this is unique for the image src encountered. this way the same data links are not entered in duplicates, 
	//              and using the hash the DB can qucikly be checked for existence of record
	Hash_str                   string        `bson:"hash_str"`

	//IMPORTANT!! - indicates if the image was fetched from the remote server,
	//              and has been stored on S3 and ready for usage by other services. 
	Downloaded_bool            bool          `bson:"downloaded_bool"`

	//IMPORTANT!! - the usage was determined to be useful for internal applications,
	//              they're not page elements, or other small unimportant parts.
	//              if it is valid for usage then a gf_image for this image should be 
	//              found in the db
	Valid_for_usage_bool       bool          `bson:"valid_for_usage_bool"`
	S3_stored_bool             bool          `bson:"s3_stored_bool"` //if persisting to s3 succeeded
	Nsfv_bool                  bool          `bson:"nsfv_bool"`      //NSFV (not safe for viewing/nudity) flag for the image 
	Image_id_str               string        `bson:"image_id_str"`   //id of the gf_image for this corresponding crawler_page_img //FIX!! - should be "gf_image_id_str"
}
	
//IMPORTANT!! - reference to an image, on a particular page. 
//              the same image, with the same Url_str can appear on multiple pages, and this 
//              struct tracks that, one record per reference
type Crawler_page_img_ref struct {
	Id                         bson.ObjectId `bson:"_id,omitempty"`
	Id_str                     string        `bson:"id_str"`
	T_str                      string        `bson:"t"`                          //"crawler_page_img_ref"
	Creation_unix_time_f       float64       `bson:"creation_unix_time_f"`
	Crawler_name_str           string        `bson:"crawler_name_str"`           //name of the crawler that discovered this image
	Cycle_run_id_str           string        `bson:"cycle_run_id_str" json:"cycle_run_id_str"`
	Url_str                    string        `bson:"url_str"`
	Domain_str                 string        `bson:"domain_str"`
	Origin_page_url_str        string        `bson:"origin_page_url_str"`        //page url from whos html this element was extracted
	Origin_page_url_domain_str string        `bson:"origin_page_url_domain_str"` //NEW_FIELD!! a lot of records dont have this field

	//IMPORTANT!! - this is unique for the image src encountered. this way the same data links are not entered in duplicates, 
	//              and using the hash the DB can qucikly be checked for existence of record
	Hash_str                   string        `bson:"hash_str"`
}

type Crawler__recent_images struct {
	Domain_str               string    `bson:"_id"                      json:"domain_str"`
	Imgs_count_int           int       `bson:"imgs_count_int"           json:"imgs_count_int"`
	Crawler_page_img_ids_lst []string  `bson:"crawler_page_img_ids_lst" json:"crawler_page_img_ids_lst"`
	Creation_times_lst       []float64 `bson:"creation_times_lst"       json:"creation_times_lst"`
	Urls_lst                 []string  `bson:"urls_lst"                 json:"urls_lst"`
	Nsfv_lst                 []bool    `bson:"nsfv_lst"                 json:"nsfv_lst"`
	Origin_page_urls_lst     []string  `bson:"origin_page_urls_lst"     json:"origin_page_urls_lst"`
}
//-------------------------------------------------
func images__prepare_and_create(p_crawler_name_str string,
			p_cycle_run_id_str   string,
			p_img_src_url_str    string,
			p_origin_page_url_str string,
			p_runtime            *Crawler_runtime,
			p_runtime_sys        *gf_core.Runtime_sys) (*Crawler_page_img,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images.images__prepare_and_create()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//------------------
	//DOMAINS

	img_src_domain_str,origin_page_url_domain_str,gf_err := gf_crawl_utils.Get_domain(p_img_src_url_str,p_origin_page_url_str,p_runtime_sys)
	if gf_err != nil {
		t:="images_in_page__get_domain__failed"
		m:="failed to get domain of image with img_src - "+p_img_src_url_str
		Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,},p_img_src_url_str,p_crawler_name_str,
					gf_err,p_runtime,p_runtime_sys)
		return nil,gf_err
	}
	//-------------
	//COMPLETE_A_HREF
	
	complete_img_src_url_str,gf_err := gf_crawl_utils.Complete_url(p_img_src_url_str,img_src_domain_str,p_runtime_sys)
	if gf_err != nil {
		t:="complete_url__failed"
		m:="failed to complete_url of image with img_src - "+p_img_src_url_str
		Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,},p_img_src_url_str,p_crawler_name_str,
					gf_err,p_runtime,p_runtime_sys)
		return nil,gf_err
	}
	//-------------
	//GET_IMG_EXT_FROM_URL

	img_ext_str,gf_err := gf_images_utils.Get_image_ext_from_url(p_img_src_url_str,p_runtime_sys)
	if gf_err != nil {
		t:="images_in_page__get_img_extension__failed"
		m:="failed to get file extension of image with img_src - "+p_img_src_url_str
		Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,},p_img_src_url_str,p_crawler_name_str,
					gf_err,p_runtime,p_runtime_sys)
		return nil,gf_err
	}
	//-------------
	p_runtime_sys.Log_fun("INFO",">>>>> "+cyan("img")+" -- "+yellow(img_src_domain_str)+" ------ "+yellow(fmt.Sprint(complete_img_src_url_str)))

	img := images__create(p_crawler_name_str,
					p_cycle_run_id_str,
					complete_img_src_url_str,
					img_ext_str,
					img_src_domain_str,
					p_origin_page_url_str,
					origin_page_url_domain_str,
					p_runtime_sys)
	return img,nil
}
//-------------------------------------------------
func images__create(p_crawler_name_str string,
			p_cycle_run_id_str           string,
			p_img_src_url_str            string,
			p_img_ext_str                string,
			p_img_src_domain_str         string,
			p_origin_page_url_str        string,
			p_origin_page_url_domain_str string,
			p_runtime_sys                *gf_core.Runtime_sys) *Crawler_page_img {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images.images__create()")

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("crawler_page_img:%f",creation_unix_time_f)

	//HASH
	to_hash_str := p_img_src_url_str //one Crawler_page_img for a given page url, no matter on how many pages it is referenced by
	hash        := md5.New()
	hash.Write([]byte(to_hash_str))
	hash_str := hex.EncodeToString(hash.Sum(nil))

	img := &Crawler_page_img{
		Id_str:                    id_str,
		T_str:                     "crawler_page_img",
		Creation_unix_time_f:      creation_unix_time_f,
		Crawler_name_str:          p_crawler_name_str,
		Cycle_run_id_str:          p_cycle_run_id_str,
		Img_ext_str:               p_img_ext_str,
		Url_str:                   p_img_src_url_str,
		Domain_str:                p_img_src_domain_str,
		Origin_page_url_str:       p_origin_page_url_str,
		Origin_page_url_domain_str:p_origin_page_url_domain_str,
		Hash_str:                  hash_str,
		Downloaded_bool:           false,
		Valid_for_usage_bool:      false, //all images are initially set as invalid for usage
		S3_stored_bool:            false, 
	}

	return img
}
//-------------------------------------------------
func images__ref_create(p_crawler_name_str string,
				p_cycle_run_id_str           string,
				p_image_url_str              string,
				p_image_url_domain_str       string,
				p_origin_page_url_str        string,
				p_origin_page_url_domain_str string,
				p_runtime_sys                *gf_core.Runtime_sys) *Crawler_page_img_ref {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images.images__ref_create()")

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	ref_id_str           := fmt.Sprintf("img_ref:%f",creation_unix_time_f)

	//HASH
	//IMPORTANT!! - one Crawler_page_img_ref per page img reference, so if the same image is linked on several pages
	//              each of those references will have a different hash_str and will created a new Crawler_page_img_ref
	to_hash_str := p_image_url_str+p_origin_page_url_str
	hash        := md5.New()
	hash.Write([]byte(to_hash_str))
	hash_str := hex.EncodeToString(hash.Sum(nil))

	gf_img_ref := &Crawler_page_img_ref{
		Id_str:                    ref_id_str,
		T_str:                     "crawler_page_img_ref",
		Creation_unix_time_f:      creation_unix_time_f,
		Crawler_name_str:          p_crawler_name_str,
		Cycle_run_id_str:          p_cycle_run_id_str,
		Url_str:                   p_image_url_str,        //complete_img_src_str,
		Domain_str:                p_image_url_domain_str, //img_src_domain_str,
		Origin_page_url_str:       p_origin_page_url_str,
		Origin_page_url_domain_str:p_origin_page_url_domain_str,
		Hash_str:                  hash_str,
	}

	return gf_img_ref
}
//-------------------------------------------------
func Images__get_recent(p_runtime_sys *gf_core.Runtime_sys) ([]Crawler__recent_images,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images.Images__get_recent()")

	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t":"crawler_page_img",
			},
		},
		bson.M{"$sort":bson.M{
				"creation_unix_time_f":-1,
			},
		},
		bson.M{"$limit":2000},
		bson.M{"$group":bson.M{
				"_id":                     "$origin_page_url_domain_str", //"$domain_str",
				"imgs_count_int":          bson.M{"$sum" :1},
				"crawler_page_img_ids_lst":bson.M{"$push":"$id_str"},
				"creation_times_lst":      bson.M{"$push":"$creation_unix_time_f"},
				"urls_lst":                bson.M{"$push":"$url_str"},
				"nsfv_ls":                 bson.M{"$push":"$nsfv_bool"},
				"origin_page_urls_lst":    bson.M{"$push":"$origin_page_url_str"},
			},
		},
	})

	results_lst := []Crawler__recent_images{}
	err         := pipe.AllowDiskUse().All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to get recent_images (crawler_page_img) by domain",
			"mongodb_aggregation_error",
			nil,err,"gf_crawl_core",p_runtime_sys)
		return nil,gf_err
	}

	return results_lst,nil
}