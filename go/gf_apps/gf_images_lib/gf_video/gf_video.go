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

package gf_video

import (
	"fmt"
	"os/exec"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// ffmpeg is not bundled in the service containers, so it has to be checked for
// before any CLI invocation that depends on it.

func checkFFmpegInstalled(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		gfErr := gf_core.ErrorCreate("ffmpeg is not installed (or not in PATH), which is required for video operations",
			"cli_tool_not_installed",
			map[string]interface{}{"cli_tool": "ffmpeg"},
			err, "gf_video", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func GetVideoFrameFromURL(pURLstr string,
	pImageLocalFilePathStr string,
	pFrameIndexInt         int,
	pRuntimeSys            *gf_core.RuntimeSys) *gf_core.GFerror {

	gfErr := checkFFmpegInstalled(pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	formatStr := "mjpeg"
	cmdLst := []string{
		"ffmpeg",
		"-seekable", fmt.Sprint(pFrameIndexInt),
		"-i", pURLstr,
		"-ss", "00:00:01.000",
		"-vframes", "1",
		"-f", formatStr,
		pImageLocalFilePathStr,
	}



	_, _, gfErr = gf_core.CLIrunStandard(cmdLst,
		nil,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	return nil
}