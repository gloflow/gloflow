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
	"fmt"
	"os"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
)

var logFun func(string,string)
var cliArgsMap map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	logFun, _  = gf_core.LogsInit()
	cliArgsMap = CLIparseArgs(logFun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func TestImageTransform(pTest *testing.T) {

	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_images_core_tests",
		LogFun:           logFun,
	}

	testImageLocalFilePathStr       := "./../tests_data/test_image_03.jpeg"
	testImageOutputLocalFilePathStr := "./../tests_data/transform/test_image_03_resized.jpeg"
	testImageThumbSizePxInt         := 600

	normalizedExtStr := "jpeg"
	
	image, gfErr := ImageLoadFile(testImageLocalFilePathStr, normalizedExtStr, runtimeSys)
	if gfErr != nil {
		pTest.FailNow()
	}

	largerDimensionPxInt, largerDimensionNameStr := GetImageLargerDimension(image, runtimeSys)
	
 
	// RESIZE_IMAGE
	thumbWidthPxInt, thumbHeightPxInt := ThumbsGetSizeInPx(testImageThumbSizePxInt,
		largerDimensionPxInt,
		largerDimensionNameStr)

	fmt.Printf("larger dimension (px) - %d %s\n", largerDimensionPxInt, largerDimensionNameStr)
	fmt.Printf("thumb size (px)       - %d %d\n", thumbWidthPxInt, thumbHeightPxInt)

	gfErr = resizeImage(image,
		testImageOutputLocalFilePathStr,
		thumbWidthPxInt,
		thumbHeightPxInt,
		runtimeSys)

	if gfErr != nil {
		pTest.FailNow()
	}
}