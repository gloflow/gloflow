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
func image__cleanup(p_img_local_file_path_str string,
	p_img_thumbs  *gf_images_core.GFimageThumbs,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_images_utils.image__cleanup()")

	files_to_remove_lst := []string{
		p_img_local_file_path_str,
	}

	if p_img_thumbs != nil {
		files_to_remove_lst = append(files_to_remove_lst, p_img_thumbs.Small_local_file_path_str)
		files_to_remove_lst = append(files_to_remove_lst, p_img_thumbs.Medium_local_file_path_str)
		files_to_remove_lst = append(files_to_remove_lst, p_img_thumbs.Large_local_file_path_str)
	}
	
	for _, f_str := range files_to_remove_lst {
		err := os.Remove(f_str)
		if err != nil {
			gf_err := gf_core.Error__create("failed to cleanup a crawled image files",
				"file_remove_error",
				map[string]interface{}{"file_str": f_str,},
				err, "gf_crawl_core", p_runtime_sys)
			return gf_err
		}
	}
	return nil
}