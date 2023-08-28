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

package gf_landing_page_lib

import (
	"fmt"
	"context"
	"strconv"
	"net/url"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//------------------------------------------------

type GFfeaturedPost struct {
	TitleStr        string
	ImageURLstr     string
	URLstr          string
	ImagesNumberInt int	
}

type GFfeaturedImage struct {
	TitleStr                   string
	ImageURLstr                string
	ImageThumbnailMediumURLstr string
	ImageOriginPageURLstr      string // for each featured image this is the URL used in links
	ImageOriginPageURLhostStr  string // this is displayed in the user UI for each featured image
	CreationUNIXtimeStr        string
	FlowNameStr                string
	OwnerUserNameStr           gf_identity_core.GFuserName
}

//------------------------------------------
// IMAGES
//------------------------------------------

func getFeaturedImgs(pMaxRandomCursorPositionInt int, // 500
	pElementsNumToGetInt int, // 5
	pFlowNameStr         string,
	pUserID              gf_core.GF_ID,
	pCtx                 context.Context,
	pRuntimeSys          *gf_core.RuntimeSys) ([]*GFfeaturedImage, *gf_core.GFerror) {

	imagesLst, err := gf_images_core.DBmongoGetRandomImagesRange(pElementsNumToGetInt,
		pMaxRandomCursorPositionInt,
		pFlowNameStr,
		pUserID,
		pRuntimeSys)

	if err != nil {
		return nil, err
	}

	featuredImagesLst := []*GFfeaturedImage{}
	usernamesCacheMap := map[gf_core.GF_ID]gf_identity_core.GFuserName{}

	for _, image := range imagesLst {

		// FIX!! - create a proper gfErr
		originPageURL, err := url.Parse(image.Origin_page_url_str)
		if err != nil {
			continue
		}

		//---------------------
		// RESOLVE_USER_ID_TO_USERNAME
		var userNameStr gf_identity_core.GFuserName
		if image.UserID != "" {
			userID := image.UserID

			// check if there is a cached user_name, and use it if present; if not, resolve from DB
			if cachedUserNameStr, ok := usernamesCacheMap[userID]; ok {
				userNameStr = cachedUserNameStr
			} else {
				resolvedUserNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(userID, pCtx, pRuntimeSys)
				if gfErr != nil {
					/*
					failing to resolve username should not fail the rendering
					of the entire flow view.
					*/
					userNameStr = gf_identity_core.GFuserName("?")
					continue
				}
				userNameStr = resolvedUserNameStr
				usernamesCacheMap[userID] = resolvedUserNameStr
			}
		} else {

			// IMPORTANT!! - pre-auth-system images are marked as owned by anonymous users.
			userNameStr = gf_identity_core.GFuserName("anon")
		}

		//---------------------




		featured := &GFfeaturedImage{
			TitleStr:                   image.TitleStr,
			ImageURLstr:                image.Thumbnail_medium_url_str,
			ImageThumbnailMediumURLstr: image.Thumbnail_medium_url_str,
			ImageOriginPageURLstr:      image.Origin_page_url_str,
			ImageOriginPageURLhostStr:  originPageURL.Host,
			CreationUNIXtimeStr:        strconv.FormatFloat(image.Creation_unix_time_f, 'f', 6, 64),
			FlowNameStr:                pFlowNameStr,
			OwnerUserNameStr:           userNameStr,
		}
		featuredImagesLst = append(featuredImagesLst, featured)
	}
	return featuredImagesLst, nil
}

//------------------------------------------
// POSTS
//------------------------------------------

func getFeaturedPosts(pMaxRandomCursorPositionInt int, // 500
	pElementsNumToGetInt int, // 5
	pRuntimeSys          *gf_core.RuntimeSys) ([]*GFfeaturedPost, *gf_core.GFerror) {

	//gets posts starting in some random position (time wise), 
	//and as many as specified after that random point
	postsLst, gfErr := gf_publisher_core.DBmongoGetRandomPostsRange(pElementsNumToGetInt,
		pMaxRandomCursorPositionInt,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	featuredPostsLst := postsToFeatured(postsLst, pRuntimeSys)
	return featuredPostsLst, nil
}

//------------------------------------------

func postsToFeatured(pPostsLst []*gf_publisher_core.GFpost, pRuntimeSys *gf_core.RuntimeSys) []*GFfeaturedPost {

	featuredPostsLst := []*GFfeaturedPost{}
	for _, post := range pPostsLst {
		featured          := postToFeatured(post, pRuntimeSys)
		featuredPostsLst = append(featuredPostsLst, featured)
	}

	// CAUTION!! - in some cases image_src is null or "error", in which case it should not 
	//             be included in the final output. This is due to past issues/bugs in the gf_image and 
	//             gf_publisher.
	featuredElementsWithNoErrorsLst := []*GFfeaturedPost{}
	for _, featured := range featuredPostsLst {
		
		pRuntimeSys.LogFun("INFO", "featured.Image_url_str - "+featured.ImageURLstr)

		if featured.ImageURLstr == "" || featured.ImageURLstr == "error" {
			errMsgStr := fmt.Sprintf("post with title [%s] has a image_src that is [%s]", featured.TitleStr, featured.ImageURLstr)
			pRuntimeSys.LogFun("ERROR", errMsgStr)
		} else {
			featuredElementsWithNoErrorsLst = append(featuredElementsWithNoErrorsLst, featured)
		}
	}

	return featuredElementsWithNoErrorsLst
}

//------------------------------------------

func postToFeatured(pPost *gf_publisher_core.GFpost, pRuntimeSys *gf_core.RuntimeSys) *GFfeaturedPost {

	postURLstr := fmt.Sprintf("/posts/%s", pPost.TitleStr)

	featured := &GFfeaturedPost{
		TitleStr:        pPost.TitleStr,
		ImageURLstr:     pPost.ThumbnailURLstr,
		URLstr:          postURLstr,
		ImagesNumberInt: len(pPost.ImagesIDsLst),
	}
	return featured
}