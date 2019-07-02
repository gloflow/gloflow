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
	"strings"
	"time"
	"net/url"
	"crypto/md5"
	"encoding/hex"
	"github.com/globalsign/mgo/bson"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_utils"
)

//--------------------------------------------------
type Gf_crawler_page_outgoing_link struct {
	Id                    bson.ObjectId `bson:"_id,omitempty"`
	Id_str                string        `bson:"id_str"`
	T_str                 string        `bson:"t"`                    //"crawler_page_outgoing_link"
	Creation_unix_time_f  float64       `bson:"creation_unix_time_f"`
	Crawler_name_str      string        `bson:"crawler_name_str"`     //name of the crawler that discovered this link
	Cycle_run_id_str      string        `bson:"cycle_run_id_str"`
	A_href_str            string        `bson:"a_href_str"`
	Domain_str            string        `bson:"domain_str"`
	Origin_url_str        string        `bson:"origin_url_str"`       //page url from whos html this element was extracted
	Origin_url_domain_str string        `bson:"origin_url_domain_str"`

	//IMPORTANT!! - this is a hash of the . it 
	Hash_str              string        `bson:"hash_str"`


	Valid_for_crawl_bool  bool          `bson:"valid_for_crawl_bool"`  //if the link should be crawled, or if it should be ignored
	Images_processed_bool bool          `bson:"images_processed_bool"` //if all the images in the page have been downloaded/transformed/stored-in-s3
	Fetched_bool          bool          `bson:"fetched_bool"`          //indicator if the link has been fetched (its html downloaded and parsed)
	Fetch_last_id_str     string        `bson:"fetch_last_id_str"`
	Fetch_last_time_f     float64       `bson:"fetch_last_time_f"`
	//-------------------
	//IMPORTANT!! - indicates if this link hasis currently being processed by some 
	//              crawler master/worker in the cluster
	Import__in_progress_bool bool       `bson:"import__in_progress_bool"`
	Import__start_time_f     float64    `bson:"import__start_time_f"` //when has the "in_progress" flag been set. for detecting interrupted/incomplete imports
	//-------------------
	//IMPORTANT!! - last error that occured/interupted processing of this link
	Error_type_str string               `bson:"error_type_str,omitempty"`
	Error_id_str   string               `bson:"error_id_str,omitempty"`
	//-------------------
}

//--------------------------------------------------
func Link__db_index__init(p_runtime_sys *gf_core.Runtime_sys) []*gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_links.Link__db_index__init()")

	indexes_keys_lst := [][]string{
		[]string{"t", "crawler_name_str"}, //all stat queries first match on "t"
		[]string{"t", "id_str"},
		[]string{"t", "hash_str"},
		[]string{"t", "hash_str", "valid_for_crawl_bool", "fetched_bool", "error_type_str", "error_id_str"}, //Link__get_unresolved()
		[]string{"t", "hash_str", "valid_for_crawl_bool"}, //Link__mark_as_resolved()
	}

	gf_errs_lst := gf_core.Mongo__ensure_index(indexes_keys_lst, "gf_crawl", p_runtime_sys)
	return gf_errs_lst
}

//--------------------------------------------------
func Link__get_unresolved(p_crawler_name_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_crawler_page_outgoing_link, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_links.Link__get_unresolved()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	black  := color.New(color.FgBlack).Add(color.BgWhite).SprintFunc()

	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println("INFO",black("GET__UNRESOLVED_LINK")+" - for crawler - "+yellow(p_crawler_name_str))
	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	query := p_runtime_sys.Mongodb_db.C("gf_crawl").Find(bson.M{
			"t":                    "crawler_page_outgoing_link",
			"crawler_name_str":     p_crawler_name_str, //get links that were discovered by this crawler
			"valid_for_crawl_bool": true,
			"fetched_bool":         false,

			//IMPORTANT!! - get all unresolved links that also dont have any errors associated
			//              with them. this way repeated processing of unresolved links that always cause 
			//              an error is avoided (wasted resources)
			"error_type_str": bson.M{"$exists":false,},
			"error_id_str":   bson.M{"$exists":false,},

			/*//-------------------
			//IMPORTANT!! - this gets all unresolved links that come from the domain 
			//              that the crawler is assigned to
			//"origin_domain_str"   :p_crawler_domain_str,
			"$or":domains_query_lst,
			//-------------------*/
		})


	var unresolved_link Gf_crawler_page_outgoing_link
	err := query.One(&unresolved_link)

	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Mongo__handle_error("unresolved links for gf_crawler were not found in mongodb",
			"mongodb_not_found_error",
			map[string]interface{}{"crawler_name_str": p_crawler_name_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get unresolved_link from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"crawler_name_str": p_crawler_name_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}

	//-------------------
	//IMPORTANT!! - some unresolved links in the DB might possibly be urlescaped,
	//              so for proper usage it is unescaped here and stored back in the unresolved_link struct.
	unescaped_unresolved_link_url_str, err := url.QueryUnescape(unresolved_link.A_href_str)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get unresolved_link from mongodb", "url_unescape_error",
			map[string]interface{}{
				"crawler_name_str":        p_crawler_name_str,
				"unresolved_link_url_str": unresolved_link.A_href_str,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}
	unresolved_link.A_href_str = unescaped_unresolved_link_url_str
	//-------------------

	fmt.Printf("unresolved_link URL - %s\n", unresolved_link.A_href_str)
	return &unresolved_link, nil
}

//--------------------------------------------------
func Links__get_outgoing_in_page(p_url_fetch *Gf_crawler_url_fetch,
	p_cycle_run_id_str string,
	p_crawler_name_str string,
	p_runtime          *Gf_crawler_runtime,
	p_runtime_sys      *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_links.Links__get_outgoing_in_page()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue   := color.New(color.FgBlue).SprintFunc()

	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println("INFO","GET__PAGE_LINKS - "+blue(p_url_fetch.Url_str))
	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	crawled_links_lst := []*Gf_crawler_page_outgoing_link{}

	p_url_fetch.goquery_doc.Find("a").Each(func(p_i int, p_elem *goquery.Selection) {

		origin_url_str := p_url_fetch.Url_str
		a_href_str,_   := p_elem.Attr("href")

		fmt.Println(">> "+cyan("a")+" --- crawler_page_outgoing_link FOUND - domain -"+p_url_fetch.Domain_str+" -- "+yellow(fmt.Sprint(a_href_str)))

		//-------------
		if a_href_str == "" {
			return
		}
		//-------------
		//IMPORTANT!! - links on some pages only contain the protocol specifier
		if a_href_str == "http://" {
			return
		}
		//-------------
		//"#" in html <a> tags is an anchor for a section of the page itself, scrolling the user to it
		//so it doesnt represent a new page itself and should not be persisted/used
		if strings.HasPrefix(a_href_str,"#") {
			return
		}
		//-------------
		//IMPORTANT!! - some sites have this javascript string as the a href value, 
		//              and it indicates to do nothing, but still look like a link
		if strings.Contains(a_href_str,"javascript:void(0)") {
			return
		}
		//-------------
		//CREATE_LINK

		link,gf_err := link__create(a_href_str,
			origin_url_str,
			p_cycle_run_id_str,
			p_crawler_name_str,
			p_runtime_sys)
		if gf_err != nil {
			t := "link__complete_url__failed"
			m := "failed completing the url of a_href_str - "+a_href_str
			Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str": p_url_fetch.Url_str,}, a_href_str, p_crawler_name_str,
				gf_err, p_runtime, p_runtime_sys)
			return
		}
		//-------------

		crawled_links_lst = append(crawled_links_lst, link)
	})
	//--------------
	//STAGE - PERSIST ALL LINKS
	for _, link := range crawled_links_lst {

		gf_err := link__db_create(link, p_runtime_sys)
		if gf_err != nil {
			t:="link__db_create__failed"
			m:="failed creating link in the DB - "+link.A_href_str
			Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str":p_url_fetch.Url_str,}, link.A_href_str, p_crawler_name_str,
				gf_err, p_runtime, p_runtime_sys)
			return
		}
	}
	//--------------
}

//--------------------------------------------------
func Link__get_db(p_link_id_str string, p_runtime_sys *gf_core.Runtime_sys) (*Gf_crawler_page_outgoing_link, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_links.Link__get_db()")

	var unresolved_link Gf_crawler_page_outgoing_link
	err := p_runtime_sys.Mongodb_db.C("gf_crawl").Find(bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link_id_str,
		}).One(&unresolved_link)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get crawler_page_outgoing_link by ID from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"link_id_str":p_link_id_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}

	return &unresolved_link, nil	
}

//--------------------------------------------------
func link__create(p_url_str string,
	p_origin_url_str   string,
	p_cycle_run_id_str string,
	p_crawler_name_str string,
	p_runtime_sys      *gf_core.Runtime_sys) (*Gf_crawler_page_outgoing_link,*gf_core.Gf_error) {

	//-------------
	//DOMAIN
	domain_str, origin_url_domain_str, gf_err := gf_crawl_utils.Get_domain(p_url_str, p_origin_url_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-------------
	//COMPLETE_A_HREF - handle urls that are relative (dont contain the domain component), 
	//                  and complete them to get the full url
	
	complete_a_href_str, gf_err := gf_crawl_utils.Complete_url(p_url_str, domain_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-------------

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("outgoing_link:%f", creation_unix_time_f)

	//-------------
	//HASH
	//IMPORTANT!! - this hash uniquely identifies links going to the same target URL that were discovered on the same origin page URL. 
	//              if a particular origin page has several links in it that point to the same target URL, then all those links
	//              will have the same hash (which can be used for efficient queries or grouping).
	hash        := md5.New()
	hash.Write([]byte(p_origin_url_str))
	hash.Write([]byte(p_url_str))
	hash_str := hex.EncodeToString(hash.Sum(nil))
	//-------------

	link__valid_for_crawl_bool := link__verify_for_crawl(p_url_str, domain_str, p_runtime_sys)
	link := &Gf_crawler_page_outgoing_link{
		Id_str:                id_str,
		T_str:                 "crawler_page_outgoing_link",
		Creation_unix_time_f:  creation_unix_time_f,
		Crawler_name_str:      p_crawler_name_str,
		Cycle_run_id_str:      p_cycle_run_id_str,
		A_href_str:            complete_a_href_str,
		Domain_str:            domain_str,
		Origin_url_str:        p_origin_url_str,
		Origin_url_domain_str: origin_url_domain_str,
		Hash_str:              hash_str,
		Valid_for_crawl_bool:  link__valid_for_crawl_bool,
		Fetched_bool:          false,
		Images_processed_bool: false,
	}

	return link, nil
}

//--------------------------------------------------
func link__db_create(p_link *Gf_crawler_page_outgoing_link, p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_links.link__db_create()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	//-------------
	//IMPORTANT!! - REXAMINE!! - to conserve on storage for potentially large savings (should be checked empirically?), links are persisted
	//                           in the DB only if their hash is unique. Hashes are composed of origin page URL and target URL hashed, so multiple links coming from the 
	//                           same origin page URL, and targeting the same URL, are only stored once.
	//                           this is a potentially loss of information, for pages that have a lot of these duplicate links. having this information 
	//                           on pages could maybe prove useful for some kind of analysis or algo. 
	//                           - so maybe store links even if their hashes are duplicates?
	//                           - add some kind of tracking where these duplicates are counted for pages.
	link_exists_bool, gf_err := link__db_exists(p_link.Hash_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//-------------

	//crawler_page_outgoing_link already exists, from previous crawls, so ignore it
	if link_exists_bool {
		fmt.Println(">> "+yellow(">>>>>>>> - DB PAGE_LINK ALREADY EXISTS >-- ")+cyan(fmt.Sprint(p_link.A_href_str)))
		return nil
	} else {

		err := p_runtime_sys.Mongodb_db.C("gf_crawl").Insert(p_link)
		if err != nil {

			gf_err := gf_core.Mongo__handle_error("failed to insert a crawler_page_outgoing_link in mongodb",
				"mongodb_insert_error",
				map[string]interface{}{
					"link_a_href_str": p_link.A_href_str,
				},
				err, "gf_crawl_core", p_runtime_sys)
			return gf_err
		}
	}

	return nil
}

//--------------------------------------------------
func link__db_exists(p_link_hash_str string, p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.Gf_error) {

	c, err := p_runtime_sys.Mongodb_db.C("gf_crawl").Find(bson.M{
		"t":        "crawler_page_outgoing_link",
		"hash_str": p_link_hash_str,
		}).Count()

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to count crawler_page_outgoing_link by its hash",
			"mongodb_find_error",
			map[string]interface{}{"hash_str": p_link_hash_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return false, gf_err
	}

	//crawler_page_outgoing_link already exists, from previous crawls, so ignore it
	if c > 0 {
		return true, nil
	} else {
		return false, nil
	}
	return false, nil
}

//--------------------------------------------------
func link__verify_for_crawl(p_url_str string,
	p_domain_str  string,
	p_runtime_sys *gf_core.Runtime_sys) bool {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_links.link__verify_for_crawl()")

	blacklisted_domains_map := get_domains_blacklist(p_runtime_sys)

	//dont crawl these mainstream sites
	if val_bool,ok := blacklisted_domains_map[p_domain_str]; ok {
		return val_bool
	}

	//unknown domains are whitelisted for crawling
	return true
}

//--------------------------------------------------
func link__mark_as_failed(p_error *Gf_crawler_error,
	p_link        *Gf_crawler_page_outgoing_link,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_links.link__mark_as_failed()")

	err := p_runtime_sys.Mongodb_db.C("gf_crawl").Update(bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link.Id_str,
		},
		bson.M{"$set": bson.M{
				"error_id_str":   p_error.Id_str,
				"error_type_str": p_error.Type_str,
			},
		})

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a crawler_page_outgoing_link in mongodb as failed",
			"mongodb_update_error",
			map[string]interface{}{
				"link_id_str":    p_link.Id_str,
				"error_id_str":   p_error.Id_str,
				"error_type_str": p_error.Type_str,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}

	return nil
}

//--------------------------------------------------
func Link__mark_import_in_progress(p_status_bool bool,
	p_unix_time_f float64,
	p_link        *Gf_crawler_page_outgoing_link,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_links.Link__mark_import_in_progress()")

	//----------------
	update_map := bson.M{
		"import__in_progress_bool":p_status_bool,
	}
	if p_status_bool {
		update_map["import__start_time_f"] = p_unix_time_f
	} else {
		update_map["import__end_time_f"] = p_unix_time_f
	}
	//----------------

	err := p_runtime_sys.Mongodb_db.C("gf_crawl").Update(bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link.Id_str,
		},
		bson.M{"$set": update_map,})

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a crawler_page_outgoing_link in mongodb as in_progress/complete",
			"mongodb_update_error",
			map[string]interface{}{
				"link_id_str": p_link.Id_str,
				"status_bool": p_status_bool,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}
	return nil
}

//--------------------------------------------------
func Link__mark_as_resolved(p_link *Gf_crawler_page_outgoing_link,
	p_fetch_id_str          string,
	p_fetch_creation_time_f float64,
	p_runtime_sys           *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_links.Link__mark_as_resolved()")

	p_link.Fetched_bool = true
	err := p_runtime_sys.Mongodb_db.C("gf_crawl").Update(bson.M{
				"t":                    "crawler_page_outgoing_link",
				"id_str":               p_link.Id_str,
				"valid_for_crawl_bool": true,
			},
			bson.M{"$set": bson.M{
				"fetched_bool":      true,
				"fetch_last_id_str": p_fetch_id_str,
				"fetch_last_time_f": p_fetch_creation_time_f,
			},
		})
	
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a crawler_page_outgoing_link in mongodb as resolved/fetched",
			"mongodb_update_error",
			map[string]interface{}{
				"link_id_str":  p_link.Id_str,
				"fetch_id_str": p_fetch_id_str,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}

	return nil
}