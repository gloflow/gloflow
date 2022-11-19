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

package gf_core

import (
	"os"
)

//---------------------------------------------------

type Gf_s3_test_info struct {
	Gf_s3_info                *GFs3Info
	Aws_access_key_id_str     string
	Aws_secret_access_key_str string
	Aws_token_str             string
}

//---------------------------------------------------

func T__get_s3_info(p_runtime_sys *RuntimeSys) *Gf_s3_test_info {

	aws_access_key_id_str     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	aws_secret_access_key_str := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	aws_token_str             := os.Getenv("GF_AWS_TOKEN")

	if aws_access_key_id_str == "" || aws_secret_access_key_str == "" {
		panic("test AWS credentials were not supplied")
	}
	
	gf_s3_info, gf_err := S3init(aws_access_key_id_str, aws_secret_access_key_str, aws_token_str, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	gf_s3_test_info := &Gf_s3_test_info{
		Gf_s3_info:                gf_s3_info,
		Aws_access_key_id_str:     aws_access_key_id_str,
		Aws_secret_access_key_str: aws_secret_access_key_str,
		Aws_token_str:             aws_token_str,
	}
	return gf_s3_test_info
}