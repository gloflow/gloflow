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
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------

func pipelineRenderLandingPage(pImagesMaxRandomCursorPositionInt int, // 500
	pPostsMaxRandomCursorPositionInt int,
	pFeaturedPostsToGetInt  int, // 5
	pFeaturedImagesToGetInt int, // 10
	pTemplate               *template.Template,
	pSubtemplatesNamesLst   []string,
	pRuntimeSys             *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	//-------------------
	// FEATURED_IMAGES - two random groups of images are fetched
	featuredImages0lst, gfErr := getFeaturedImgs(pImagesMaxRandomCursorPositionInt,
		pFeaturedImagesToGetInt,
		"general",
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	featuredImages1lst, gfErr := getFeaturedImgs(pImagesMaxRandomCursorPositionInt,
		pFeaturedImagesToGetInt,
		"general",
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	//-------------------
	featuredPostsLst, gfErr := getFeaturedPosts(pPostsMaxRandomCursorPositionInt,
		pFeaturedPostsToGetInt,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	templateRenderedStr, gfErr := renderTemplate(featuredPostsLst,
		featuredImages0lst,
		featuredImages1lst,
		pTemplate,
		pSubtemplatesNamesLst,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return templateRenderedStr, nil
}