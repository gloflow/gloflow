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
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_utils"
)

//--------------------------------------------------
type GFcrawlerPageOutgoingLink struct {
	Id                    primitive.ObjectID `bson:"_id,omitempty"`
	Id_str                string        `bson:"id_str"`
	T_str                 string        `bson:"t"`                    //"crawler_page_outgoing_link"
	Creation_unix_time_f  float64       `bson:"creation_unix_time_f"`
	Crawler_name_str      string        `bson:"crawler_name_str"`     //name of the crawler that discovered this link
	Cycle_run_id_str      string        `bson:"cycle_run_id_str"`
	A_href_str            string        `bson:"a_href_str"`
	Domain_str            string        `bson:"domain_str"`
	Origin_url_str        string        `bson:"origin_url_str"`       //page url from whos html this element was extracted
	Origin_url_domain_str string        `bson:"origin_url_domain_str"`

	// IMPORTANT!! - this is a hash of the . it 
	Hash_str string `bson:"hash_str"`


	Valid_for_crawl_bool  bool          `bson:"valid_for_crawl_bool"`  //if the link should be crawled, or if it should be ignored
	Images_processed_bool bool          `bson:"images_processed_bool"` //if all the images in the page have been downloaded/transformed/stored-in-s3
	Fetched_bool          bool          `bson:"fetched_bool"`          //indicator if the link has been fetched (its html downloaded and parsed)
	Fetch_last_id_str     string        `bson:"fetch_last_id_str"`
	Fetch_last_time_f     float64       `bson:"fetch_last_time_f"`

	//-------------------
	// IMPORTANT!! - indicates if this link hasis currently being processed by some 
	//               crawler master/worker in the cluster
	Import__in_progress_bool bool       `bson:"import__in_progress_bool"`
	Import__start_time_f     float64    `bson:"import__start_time_f"` //when has the "in_progress" flag been set. for detecting interrupted/incomplete imports
	
	//-------------------
	// IMPORTANT!! - last error that occured/interupted processing of this link
	Error_type_str string               `bson:"error_type_str,omitempty"`
	Error_id_str   string               `bson:"error_id_str,omitempty"`

	//-------------------
}

//--------------------------------------------------
func link__create(pURLstr string,
	pOriginURLstr   string,
	pCycleRunIDstr  string,
	pCrawlerNameStr string,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFcrawlerPageOutgoingLink, *gf_core.GFerror) {

	//-------------
	// DOMAIN
	domain_str, origin_url_domain_str, gf_err := gf_crawl_utils.Get_domain(pURLstr, pOriginURLstr, pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-------------
	// COMPLETE_A_HREF - handle urls that are relative (dont contain the domain component), 
	//                   and complete them to get the full url
	
	complete_a_href_str, gf_err := gf_crawl_utils.Complete_url(pURLstr, domain_str, pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}
	//-------------

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("outgoing_link:%f", creation_unix_time_f)

	//-------------
	// HASH
	// IMPORTANT!! - this hash uniquely identifies links going to the same target URL that were discovered on the same origin page URL. 
	//               if a particular origin page has several links in it that point to the same target URL, then all those links
	//               will have the same hash (which can be used for efficient queries or grouping).
	hash := md5.New()
	hash.Write([]byte(pOriginURLstr))
	hash.Write([]byte(pURLstr))
	hash_str := hex.EncodeToString(hash.Sum(nil))

	//-------------

	link__valid_for_crawl_bool := link__verify_for_crawl(pURLstr, domain_str, pRuntimeSys)
	link := &GFcrawlerPageOutgoingLink{
		Id_str:                id_str,
		T_str:                 "crawler_page_outgoing_link",
		Creation_unix_time_f:  creation_unix_time_f,
		Crawler_name_str:      pCrawlerNameStr,
		Cycle_run_id_str:      pCycleRunIDstr,
		A_href_str:            complete_a_href_str,
		Domain_str:            domain_str,
		Origin_url_str:        pOriginURLstr,
		Origin_url_domain_str: origin_url_domain_str,
		Hash_str:              hash_str,
		Valid_for_crawl_bool:  link__valid_for_crawl_bool,
		Fetched_bool:          false,
		Images_processed_bool: false,
	}

	return link, nil
}

//--------------------------------------------------
func Links__get_outgoing_in_page(pURLfetch *Gf_crawler_url_fetch,
	pCycleRunIDstr  string,
	pCrawlerNameStr string,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_links.Links__get_outgoing_in_page()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue   := color.New(color.FgBlue).SprintFunc()

	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println("INFO","GET__PAGE_LINKS - "+blue(pURLfetch.Url_str))
	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	crawled_links_lst := []*GFcrawlerPageOutgoingLink{}

	pURLfetch.goquery_doc.Find("a").Each(func(p_i int, p_elem *goquery.Selection) {

		origin_url_str := pURLfetch.Url_str
		a_href_str,_   := p_elem.Attr("href")

		fmt.Println(">> "+cyan("<a>")+" --- crawler_page_outgoing_link FOUND - domain - "+pURLfetch.Domain_str+" -- "+yellow(fmt.Sprint(a_href_str)))

		//-------------
		if a_href_str == "" {
			return
		}

		//-------------
		// IMPORTANT!! - links on some pages only contain the protocol specifier
		if a_href_str == "http://" {
			return
		}

		//-------------
		// "#" in html <a> tags is an anchor for a section of the page itself, scrolling the user to it
		// so it doesnt represent a new page itself and should not be persisted/used
		if strings.HasPrefix(a_href_str,"#") {
			return
		}

		//-------------
		// IMPORTANT!! - some sites have this javascript string as the a href value, 
		//               and it indicates to do nothing, but still look like a link
		if strings.Contains(a_href_str,"javascript:void(0)") {
			return
		}

		//-------------
		// CREATE_LINK

		link,gf_err := link__create(a_href_str,
			origin_url_str,
			pCycleRunIDstr,
			pCrawlerNameStr,
			pRuntimeSys)
		if gf_err != nil {
			t := "link__complete_url__failed"
			m := "failed completing the url of a_href_str - "+a_href_str
			Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str": pURLfetch.Url_str,}, a_href_str, pCrawlerNameStr,
				gf_err, pRuntime, pRuntimeSys)
			return
		}

		//-------------

		crawled_links_lst = append(crawled_links_lst, link)
	})

	//--------------
	// STAGE - PERSIST ALL LINKS
	for _, link := range crawled_links_lst {

		gf_err := link__db_create(link, pRuntimeSys)
		if gf_err != nil {
			t := "link__db_create__failed"
			m := "failed creating link in the DB - "+link.A_href_str
			Create_error_and_event(t, m, map[string]interface{}{"origin_page_url_str": pURLfetch.Url_str,}, link.A_href_str, pCrawlerNameStr,
				gf_err, pRuntime, pRuntimeSys)
			return
		}
	}

	//--------------
}

//--------------------------------------------------
func link__verify_for_crawl(pURLstr string,
	p_domain_str  string,
	pRuntimeSys *gf_core.RuntimeSys) bool {
	// pRuntimeSys.Log_fun("FUN_ENTER","gf_crawl_links.link__verify_for_crawl()")

	blacklisted_domains_map := get_domains_blacklist(pRuntimeSys)

	// dont crawl these mainstream sites
	if val_bool,ok := blacklisted_domains_map[p_domain_str]; ok {
		return val_bool
	}

	// unknown domains are whitelisted for crawling
	return true
}
