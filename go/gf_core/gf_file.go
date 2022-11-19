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

package gf_core

import (
	"io/ioutil"
)

//---------------------------------------------------

func FileCopy(pSourceFileLocalPathStr string,
	pTargetFileLocalPathStr string,
	pRuntimeSys             *RuntimeSys) *GFerror {


	sourceFileBytesLst, err := ioutil.ReadFile(pSourceFileLocalPathStr)

	if err != nil {
		gfErr := ErrorCreate("failed to read local file in order to copy it",
			"file_read_error",
			map[string]interface{}{"source_file_local_path": pSourceFileLocalPathStr,},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}

	err = ioutil.WriteFile(pTargetFileLocalPathStr, sourceFileBytesLst, 0644)
	if err != nil {
		
		gfErr := ErrorCreate("failed to write local file in order to copy it",
			"file_write_error",
			map[string]interface{}{"target_file_local_path": pTargetFileLocalPathStr,},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}

	return nil
}