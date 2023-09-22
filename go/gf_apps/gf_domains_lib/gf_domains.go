/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

package gf_domains_lib

import (
	"fmt"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
//IMPORTANT!! - this statistic used by the gf_domains GF app, directly by the end-user
//              (not only by the admin user)

type GFdomainImages struct {
	Name_str            string         `bson:"_id"`
	Count_int           int            `bson:"count_int"`           // total count of all subpages counts
	Subpages_Counts_map map[string]int `bson:"subpages_counts_map"` // counts of individual sub-page urls that images come from
}

type GFdomainPosts struct {
	Name_str  string `bson:"name_str"`
	Count_int int    `bson:"count_int"`
}

//--------------------------------------------------

func accumulateDomains(pPostsDomainsLst []GFdomainPosts,
	pImagesDomainsLst []GFdomainImages,
	pRuntimeSys       *gf_core.RuntimeSys) map[string]GFdomain {

	domainsMap := map[string]GFdomain{}

	//--------------------------------------------------
	// IMAGES DOMAINS
	for _, imagesDomain := range pImagesDomainsLst {

		domainNameStr := imagesDomain.Name_str

		if domain,ok := domainsMap[domainNameStr]; ok {
			domain.Domain_images = imagesDomain
			domain.Count_int     = domain.Count_int + imagesDomain.Count_int
		} else {
			// IMPORTANT!! - no existing domain with this domain_str has been found
			creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
			id_str               := fmt.Sprintf("domain:%f", creation_unix_time_f)
			new_domain := GFdomain{
				Id_str:        id_str,
				T_str:         "domain",
				Name_str:      domainNameStr,
				Count_int:     imagesDomain.Count_int,
				Domain_images: imagesDomain,
			}
			domainsMap[domainNameStr] = new_domain
		}
	}
	
	//--------------------------------------------------
	// POSTS DOMAINS
	// IMPORTANT!! - these run first so they just create a Domain struct without checks

	for _, domain_posts := range pPostsDomainsLst {

		domain_name_str := domain_posts.Name_str

		// IMPORTANT!! - no existing domain with this domain_str has been found
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := fmt.Sprintf("domain:%f", creation_unix_time_f)
		new_domain           := GFdomain{
			Id_str:       id_str,
			T_str:        "domain",
			Name_str:     domain_name_str,
			Count_int:    domain_posts.Count_int,
			Domain_posts: domain_posts,
		}
		domainsMap[domain_name_str] = new_domain
	}

	//--------------------------------------------------

	return domainsMap
}