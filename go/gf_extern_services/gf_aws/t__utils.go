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

package gf_aws

import (
	// "os"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFs3TestInfo struct {
	GFs3Info              *GFs3Info
	// AWSaccessKeyIDstr     string
	// AWSsecretAccessKeyStr string
	// AWStokenStr           string
}

//---------------------------------------------------

func TgetS3info(pRuntimeSys *gf_core.RuntimeSys) *GFs3TestInfo {

	/*
	awsAccessKeyIDstr     := os.Getenv("GF_AWS_ACCESS_KEY_ID")
	awsSecretAccessKeyStr := os.Getenv("GF_AWS_SECRET_ACCESS_KEY")
	awsTokenStr           := os.Getenv("GF_AWS_TOKEN")

	if awsAccessKeyIDstr == "" || awsSecretAccessKeyStr == "" {
		panic("test AWS credentials were not supplied")
	}
	*/

	s3info, gfErr := S3init(// awsAccessKeyIDstr, awsSecretAccessKeyStr, awsTokenStr,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	s3testInfo := &GFs3TestInfo{
		GFs3Info:              s3info,
		// AWSaccessKeyIDstr:     awsAccessKeyIDstr,
		// AWSsecretAccessKeyStr: awsSecretAccessKeyStr,
		// AWStokenStr:           awsTokenStr,
	}
	return s3testInfo
}