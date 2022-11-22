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
	"os"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------

func imageCleanup(pImgLocalFilePathStr string,
	pImgThumbs  *gf_images_core.GFimageThumbs,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	filesToRemoveLst := []string{
		pImgLocalFilePathStr,
	}

	if pImgThumbs != nil {
		filesToRemoveLst = append(filesToRemoveLst, pImgThumbs.Small_local_file_path_str)
		filesToRemoveLst = append(filesToRemoveLst, pImgThumbs.Medium_local_file_path_str)
		filesToRemoveLst = append(filesToRemoveLst, pImgThumbs.Large_local_file_path_str)
	}
	
	for _, fStr := range filesToRemoveLst {
		err := os.Remove(fStr)
		if err != nil {
			gf_err := gf_core.ErrorCreate("failed to cleanup a crawled image files",
				"file_remove_error",
				map[string]interface{}{"file_str": fStr,},
				err, "gf_crawl_core", pRuntimeSys)
			return gf_err
		}
	}
	return nil
}