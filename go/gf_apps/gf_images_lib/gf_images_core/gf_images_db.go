/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_images_core

import (
	"context"
	"math/rand"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// GET_IMAGE

func DBgetImage(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	var image *GFimage
	var gfErr *gf_core.GFerror

	// SQL
	image, gfErr = DBsqlGetImage(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// if there is no image found with desired ID in SQL, try to get it from MongoDB
	if image == nil {

		// MONGODB
		image, gfErr = DBmongoGetImage(pImageIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	return image, nil
}

//---------------------------------------------------

func DBaddTagToImage(pImageIDstr GFimageID,
	pTagsLst  []string,
	pCtx      context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// SQL
	existsBool, gfErr := DBsqlImageExistsByID(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	if existsBool {
		
		// SQL
		gfErr = DBsqlAddTagsToImage(pImageIDstr, pTagsLst, pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	} else {

		// MONGO
		gfErr = DBmongoAddTagsToImage(pImageIDstr, pTagsLst, pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	return nil
}

//---------------------------------------------------
// IMAGE_EXISTS
//---------------------------------------------------
// IMAGE_EXISTS_BY_ID

func DBimageExistsByID(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	// SQL
	existsBool, gfErr := DBsqlImageExistsByID(pImageIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	// if there is no image found with desired ID in SQL, try to get it from MongoDB
	if !existsBool {

		// MONGODB
		existsBool, gfErr = DBmongoImageExistsByID(pImageIDstr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}
	}

	return existsBool, nil
}

//---------------------------------------------------
// IMAGE_EXISTS_BY_URLS

func DBimageExistsByURLs(pImagesExternURLsLst []string,
	pFlowNameStr   string,
	pClientTypeStr string,
	pUserID        gf_core.GF_ID,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	// SQL
	sqlExistingImagesLst, gfErr := DBsqlImagesExistByURLs(pImagesExternURLsLst,
		pFlowNameStr,
		// pClientTypeStr,
		pUserID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	// MONGO
	mongoExistingImagesLst, gfErr := DBmongoImagesExistByURLs(pImagesExternURLsLst,
		pFlowNameStr,
		pClientTypeStr,
		pUserID,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	imagesLst := MergeImageMapsLists(mongoExistingImagesLst, sqlExistingImagesLst)
	return imagesLst, nil
}

//---------------------------------------------------
// MERGE
//---------------------------------------------------

func MergeImagesLists(pMongoImagesLst, pSQLimagesLst []*GFimage) []*GFimage {
	
	imageMap := make(map[string]*GFimage)

	// add SQL images to the map
	for _, img := range pSQLimagesLst {
		imageMap[string(img.IDstr)] = img
	}

	// add MongoDB images to the map if they don't already exist
	for _, img := range pMongoImagesLst {
		if _, exists := imageMap[string(img.IDstr)]; !exists {
			imageMap[string(img.IDstr)] = img
		}
	}

	imagesLst := make([]*GFimage, 0, len(imageMap))
	for _, img := range imageMap {
		imagesLst = append(imagesLst, img)
	}

	// RANDOMIZE_ORDER
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(imagesLst), func(i, j int) {
		imagesLst[i], imagesLst[j] = imagesLst[j], imagesLst[i]
	})

	return imagesLst
}

//---------------------------------------------------

func MergeImageMapsLists(pMongoImagesLst, pSQLimagesLst []map[string]interface{}) []map[string]interface{} {
	
	imageMap := make(map[string]map[string]interface{})

	// add SQL images to the map (SQL takes precedence)
	for _, img := range pSQLimagesLst {
		id, ok := img["id_str"].(string)
		if ok {
			imageMap[id] = img
		}
	}

	// add MongoDB images to the map if they don't already exist
	for _, img := range pMongoImagesLst {
		id, ok := img["id_str"].(string)
		if ok {
			if _, exists := imageMap[id]; !exists {
				imageMap[id] = img
			}
		}
	}

	imagesLst := make([]map[string]interface{}, 0, len(imageMap))
	for _, img := range imageMap {
		imagesLst = append(imagesLst, img)
	}

	return imagesLst
}
