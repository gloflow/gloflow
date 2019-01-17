/*
GloFlow media management/publishing system
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

package gf_core

import (
	"os"
)
//---------------------------------------------------
func T__get_s3_info(p_runtime_sys *Runtime_sys) *Gf_s3_info {

	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")

	if aws_access_key_id_str == "" || aws_secret_access_key_str == "" {
		panic("test AWS credentials were not supplied")
	}
	
	s3_info, gf_err := S3__init(aws_access_key_id_str, aws_secret_access_key_str, aws_token_str, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	return s3_info
}