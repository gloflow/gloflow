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
	IDstr                 string        `bson:"id_str"`
	T_str                 string        `bson:"t"`                    // "crawler_page_outgoing_link"
	CreationUNIXtimeF     float64       `bson:"creation_unix_time_f"`
	Crawler_name_str      string        `bson:"crawler_name_str"`     // name of the crawler that discovered this link
	Cycle_run_id_str      string        `bson:"cycle_run_id_str"`
	A_href_str            string        `bson:"a_href_str"`
	Domain_str            string        `bson:"domain_str"`
	Origin_url_str        string        `bson:"origin_url_str"`       // page url from whos html this element was extracted
	Origin_url_domain_str string        `bson:"origin_url_domain_str"`

	// IMPORTANT!! - this is a hash of the . it 
	Hash_str string `bson:"hash_str"`


	Valid_for_crawl_bool  bool          `bson:"valid_for_crawl_bool"`  // if the link should be crawled, or if it should be ignored
	Images_processed_bool bool          `bson:"images_processed_bool"` // if all the images in the page have been downloaded/transformed/stored-in-s3
	Fetched_bool          bool          `bson:"fetched_bool"`          // indicator if the link has been fetched (its html downloaded and parsed)
	Fetch_last_id_str     string        `bson:"fetch_last_id_str"`
	Fetch_last_time_f     float64       `bson:"fetch_last_time_f"`

	//-------------------
	// IMPORTANT!! - indicates if this link hasis currently being processed by some 
	//               crawler master/worker in the cluster
	Import__in_progress_bool bool       `bson:"import__in_progress_bool"`
	Import__start_time_f     float64    `bson:"import__start_time_f"` // when has the "in_progress" flag been set. for detecting interrupted/incomplete imports
	
	//-------------------
	// IMPORTANT!! - last error that occured/interupted processing of this link
	Error_type_str string               `bson:"error_type_str,omitempty"`
	Error_id_str   string               `bson:"error_id_str,omitempty"`

	//-------------------
}

//--------------------------------------------------

func linkCreate(pURLstr string,
	pOriginURLstr   string,
	pCycleRunIDstr  string,
	pCrawlerNameStr string,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFcrawlerPageOutgoingLink, *gf_core.GFerror) {

	//-------------
	// DOMAIN
	domainStr, originURLdomainStr, gfErr := gf_crawl_utils.GetDomain(pURLstr, pOriginURLstr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-------------
	// COMPLETE_A_HREF - handle urls that are relative (dont contain the domain component), 
	//                   and complete them to get the full url
	
	completeAhrefStr, gfErr := gf_crawl_utils.CompleteURL(pURLstr, domainStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	//-------------

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	idStr             := fmt.Sprintf("outgoing_link:%f", creationUNIXtimeF)

	//-------------
	// HASH
	// IMPORTANT!! - this hash uniquely identifies links going to the same target URL that were discovered on the same origin page URL. 
	//               if a particular origin page has several links in it that point to the same target URL, then all those links
	//               will have the same hash (which can be used for efficient queries or grouping).
	hash := md5.New()
	hash.Write([]byte(pOriginURLstr))
	hash.Write([]byte(pURLstr))
	hashStr := hex.EncodeToString(hash.Sum(nil))

	//-------------

	link__valid_for_crawl_bool := linkVerifyForCrawl(pURLstr, domainStr, pRuntimeSys)
	link := &GFcrawlerPageOutgoingLink{
		IDstr:                 idStr,
		T_str:                 "crawler_page_outgoing_link",
		CreationUNIXtimeF:     creationUNIXtimeF,
		Crawler_name_str:      pCrawlerNameStr,
		Cycle_run_id_str:      pCycleRunIDstr,
		A_href_str:            completeAhrefStr,
		Domain_str:            domainStr,
		Origin_url_str:        pOriginURLstr,
		Origin_url_domain_str: originURLdomainStr,
		Hash_str:              hashStr,
		Valid_for_crawl_bool:  link__valid_for_crawl_bool,
		Fetched_bool:          false,
		Images_processed_bool: false,
	}

	return link, nil
}

//--------------------------------------------------

func LinksGetOutgoingInPage(pURLfetch *GFcrawlerURLfetch,
	pCycleRunIDstr  string,
	pCrawlerNameStr string,
	pRuntime        *GFcrawlerRuntime,
	pRuntimeSys     *gf_core.RuntimeSys) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links.LinksGetOutgoingInPage()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue   := color.New(color.FgBlue).SprintFunc()

	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println("INFO","GET__PAGE_LINKS - "+blue(pURLfetch.Url_str))
	fmt.Println("INFO",cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	crawledLinksLst := []*GFcrawlerPageOutgoingLink{}

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

		link,gf_err := linkCreate(a_href_str,
			origin_url_str,
			pCycleRunIDstr,
			pCrawlerNameStr,
			pRuntimeSys)
		if gf_err != nil {
			t := "link__complete_url__failed"
			m := "failed completing the url of a_href_str - "+a_href_str
			CreateErrorAndEvent(t, m, map[string]interface{}{"origin_page_url_str": pURLfetch.Url_str,}, a_href_str, pCrawlerNameStr,
				gf_err, pRuntime, pRuntimeSys)
			return
		}

		//-------------

		crawledLinksLst = append(crawledLinksLst, link)
	})

	//--------------
	// STAGE - PERSIST ALL LINKS
	for _, link := range crawledLinksLst {

		gfErr := linkDBcreate(link, pRuntimeSys)
		if gfErr != nil {
			t := "link__db_create__failed"
			m := "failed creating link in the DB - "+link.A_href_str
			CreateErrorAndEvent(t, m, map[string]interface{}{"origin_page_url_str": pURLfetch.Url_str,}, link.A_href_str, pCrawlerNameStr,
				gfErr, pRuntime, pRuntimeSys)
			return
		}
	}

	//--------------
}

//--------------------------------------------------

func linkVerifyForCrawl(pURLstr string,
	pDomainStr  string,
	pRuntimeSys *gf_core.RuntimeSys) bool {

	blacklistedDomainsMap := getDomainsBlacklist(pRuntimeSys)

	// dont crawl these mainstream sites
	if valBool, ok := blacklistedDomainsMap[pDomainStr]; ok {
		return valBool
	}

	// unknown domains are whitelisted for crawling
	return true
}
