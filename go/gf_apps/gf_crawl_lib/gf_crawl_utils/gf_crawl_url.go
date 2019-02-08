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

package gf_crawl_utils

import (
	"path"
	"strings"
	"net/url"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
func Complete_url(p_url_str string,
	p_domain_str  string,
	p_runtime_sys *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_crawler_url.complete_url()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//full url, no work to be done
	if strings.HasPrefix(p_url_str, "http://") || strings.HasPrefix(p_url_str, "https://") {
		return p_url_str, nil
	} else {

		//-----------------
		//RELATIVE_URL
		u,err := url.Parse("http://"+p_domain_str)
		if err != nil {
			gf_err := gf_core.Error__create("failed to parse a domain to complete a url",
				"url_parse_error",
				&map[string]interface{}{
					"url_str":   p_url_str,
					"domain_str":p_domain_str,
				},
				err, "gf_crawl_utils", p_runtime_sys)
			return "", gf_err
		}

		//IMPORTANT!! - path.Join() handles cases where p_url_str might start with "/" or not
		u.Path        = path.Join(u.Path,p_url_str)
		full_url_str := u.String()

		p_runtime_sys.Log_fun("INFO", cyan("COMPLETED_URL")+" - "+yellow(full_url_str))
		//-----------------

		return full_url_str, nil
	}	
	return "", nil
}
//--------------------------------------------------
func Get_domain(p_link_url_str string,
	p_origin_url_str string,
	p_runtime_sys    *gf_core.Runtime_sys) (string, string, *gf_core.Gf_error) {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_crawler_url.get_domain()")

	origin_url,err := url.Parse(p_origin_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse p_origin_url_str to get its domain",
			"url_parse_error",
			&map[string]interface{}{
				"link_url_str":  p_link_url_str,
				"origin_url_str":p_origin_url_str,
			},
			err, "gf_crawl_utils", p_runtime_sys)
		return "", "", gf_err
	}


	var domain_str string
	origin_domain_str := strings.TrimPrefix(origin_url.Host, "www.")


	//IMPORTANT!! - "//" - is for "scheme relative" or "protocol relative" URI's, which are 
	//                     correctly parsed by the "url" library 
	if strings.HasPrefix(p_link_url_str, "//") {
		url,err := url.Parse(p_link_url_str)
		if err != nil {
			gf_err := gf_core.Error__create("failed to parse p_link_url_str starting with '//' to get its domain",
				"url_parse_error",
				&map[string]interface{}{
					"link_url_str":  p_link_url_str,
					"origin_url_str":p_origin_url_str,
				},
				err, "gf_crawl_utils", p_runtime_sys)
			return "", "", gf_err
		}

		domain_str = strings.TrimPrefix(url.Host, "www.")

	} else if strings.HasPrefix(p_link_url_str, "/") {
		//IMPORTANT!! - if p_link_url_str starts with "/" it is a relative link, and therefore
		//              shares the domain with the origin_url_str, url of the page from which the link
		//              was extracted.
		domain_str = origin_domain_str //since this is a relative url, url_domain and origin_domain are the same
	} else {
		url,err := url.Parse(p_link_url_str)
		if err != nil {
			gf_err := gf_core.Error__create("failed to parse p_link_url_str with no prefix '//' or '/' to get its domain",
				"url_parse_error",
				&map[string]interface{}{
					"link_url_str":  p_link_url_str,
					"origin_url_str":p_origin_url_str,
				},
				err, "gf_crawl_utils", p_runtime_sys)
			return "", "", gf_err
		}

		domain_str = strings.TrimPrefix(url.Host, "www.")
	}

	return domain_str, origin_domain_str, nil
}