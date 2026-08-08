/*
GloFlow application and media management/publishing platform
Copyright (C) 2026 Ivan Trajkovic

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
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// IMAGE_OUTPUT_FORMAT
//
// controls the encoding format used for generated thumbnails
// (small/medium/large) in CreateThumbnails()/resizeImage().
//
// IMPORTANT!! - this is distinct from the *source*/original image format
//               (see NormalizeImageFormat()/CheckImageFormat() in
//               gf_images_utils.go), which describes the format of the
//               image being uploaded/fetched, not the format thumbnails
//               are encoded/written in.
//
// historically thumbnails were always encoded as JPEG, which is well
// suited for photographic images but introduces visible fuzziness/
// artifacts around sharp edges and text - a bad fit for infographics/
// plots/screenshots. callers can now explicitly request PNG output for
// these cases.

type GFimageOutputFormat string

const GFimageOutputFormatJPEG GFimageOutputFormat = "jpeg"
const GFimageOutputFormatPNG  GFimageOutputFormat = "png"

// GFimageOutputFormatDefault - preserves prior behavior (thumbnails always
// encoded as JPEG) when no explicit output format is requested.
const GFimageOutputFormatDefault = GFimageOutputFormatJPEG

//-------------------------------------------------
// OUTPUT_FORMAT_FROM_STR - validates/normalizes a user/caller-supplied
// output format string (e.g. from a job message, an MCP tool arg, or a
// CLI flag) into a GFimageOutputFormat.
//
// an empty string is treated as "use the default" (nil-like), matching
// the pattern used for pJPEGqualityInt (nil -> JPEGqualityDefaultInt)
// elsewhere in this package.
func OutputFormatFromStr(pFormatStr string,
	pRuntimeSys *gf_core.RuntimeSys) (GFimageOutputFormat, *gf_core.GFerror) {

	normalizedStr := strings.ToLower(strings.TrimSpace(pFormatStr))

	switch normalizedStr {
	case "":
		return GFimageOutputFormatDefault, nil
	case string(GFimageOutputFormatJPEG), "jpg":
		return GFimageOutputFormatJPEG, nil
	case string(GFimageOutputFormatPNG):
		return GFimageOutputFormatPNG, nil
	default:
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("unsupported image output format - %s (supported values - jpeg, png)", pFormatStr),
			"verify__invalid_image_output_format_error",
			map[string]interface{}{"output_format_str": pFormatStr,},
			nil, "gf_images_core", pRuntimeSys)
		return "", gfErr
	}
}

//-------------------------------------------------
// FILE_EXT - the file extension (without leading dot) used for a given
// output format, e.g. for naming generated thumbnail files.
func (p GFimageOutputFormat) FileExt() string {
	switch p {
	case GFimageOutputFormatPNG:
		return "png"
	default:
		return "jpeg"
	}
}
