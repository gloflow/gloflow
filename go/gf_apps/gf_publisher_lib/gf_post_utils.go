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

package gf_publisher_lib

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//------------------------------------------
// TAGS
//------------------------------------------

func AddTagsToPostInDB(p_post_title_str string,
	p_tags_lst  []string,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_publisher_core.GFpost, *gf_core.GFerror) {
	
	post, gfErr := gf_publisher_core.DBgetPost(p_post_title_str, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	addTagsToPost(post, p_tags_lst, pRuntimeSys)
	fmt.Println(fmt.Sprintf(" --------- post tags - %s", post.TagsLst))

	gfErr = gf_publisher_core.DBupdatePost(post, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return post, nil
}

//------------------------------------------

func addTagsToPost(p_post *gf_publisher_core.GFpost,
	p_tags_lst    []string,
	pRuntimeSys *gf_core.RuntimeSys) {
	
	if len(p_tags_lst) > 0 {
		p_post.TagsLst = append(p_post.TagsLst, p_tags_lst...)

		//---------------
		// eliminate duplicates from the list, in case 
		// some of the tags just added already exist in the list of all tags

		encountered_map   := map[string]bool{}
		no_dupliactes_lst := []string{}

		for _, t_str := range p_post.TagsLst {
			if encountered_map[t_str] {
				// tuplicate exists
			} else {
				encountered_map[t_str] = true
				no_dupliactes_lst      = append(no_dupliactes_lst, t_str)
 			}
		}

		//---------------
		
		p_post.TagsLst = no_dupliactes_lst
	} else {
		p_post.TagsLst = append(p_post.TagsLst, p_tags_lst...)
	}
}

//---------------------------------------------------

func getPostsSmallThumbnailsURLs(pPostsLst []*gf_publisher_core.GFpost, pRuntimeSys *gf_core.RuntimeSys) map[string][]string {
	
	postsSmallThumbnailsURLsMap := map[string][]string{}
	for _, post := range pPostsLst {

		postSmallThumbnailsURLsLst := []string{}
		for _, postElement := range post.PostElementsLst {

			thumbURLstr               := postElement.ImgThumbnailSmallURLstr
			postSmallThumbnailsURLsLst = append(postSmallThumbnailsURLsLst, thumbURLstr)
		}
		postsSmallThumbnailsURLsMap[post.TitleStr] = postSmallThumbnailsURLsLst
	}
	return postsSmallThumbnailsURLsMap
}