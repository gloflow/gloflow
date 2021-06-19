/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_image_storage

import (
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GF_image_storage struct {
	Type_str string // "local_fs" | "s3" | "ipfs"
	S3       GF_image_storage__s3
}

type GF_image_storage__s3 struct {
	Bucket_name_str string
	Info            *gf_core.GF_s3_info
}

//---------------------------------------------------
// PUT_FILE
func Put_file(p_source_file_path_str string,
	p_target_file_path_str string,
	p_storage              *GF_image_storage,
	p_runtime_sys          *gf_core.Runtime_sys) *gf_core.Gf_error {

	if p_storage.Type_str == "s3" {
		gf_err := gf_core.S3__upload_file(p_source_file_path_str,
			p_target_file_path_str,
			p_storage.S3.Bucket_name_str,
			p_storage.S3.Info,
			p_runtime_sys)
			
		if gf_err != nil {
			return gf_err
		}
	}

	return nil

}

//---------------------------------------------------
// GET_FILE
func Get_file(p_storage *GF_image_storage,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {




	if p_storage.Type_str == "s3" {
		
	}


	return nil

}