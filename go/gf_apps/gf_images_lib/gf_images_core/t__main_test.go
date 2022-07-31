/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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
	"os"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

var logFun func(string,string)
var cliArgsMap map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	logFun     = gf_core.Init_log_fun()
	cliArgsMap = CLI__parse_args(logFun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func TestImageTransform(pTest *testing.T) {
	
	

	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_images_core_tests",
		Log_fun:          logFun,
	}



	testImageLocalFilePathStr       := "./../tests_data/test_image_03.jpeg"
	testImageOutputLocalFilePathStr := "./../tests_data/transform/test_image_03_resized.jpeg"
	testImageSizePxInt              := 600

	normalizedExtStr := "jpeg"
	img, gfErr := ImageLoadFile(testImageLocalFilePathStr, normalizedExtStr, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}


	
	
	// RESIZE_IMAGE
	gfErr = resizeImage(img,
		testImageOutputLocalFilePathStr,
		testImageSizePxInt,
		runtimeSys)

	if gfErr != nil {
		pTest.FailNow()
	}
}