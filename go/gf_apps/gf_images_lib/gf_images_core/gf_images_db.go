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
	"math/rand"
	"time"
)

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
