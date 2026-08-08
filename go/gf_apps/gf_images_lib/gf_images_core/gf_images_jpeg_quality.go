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

package gf_images_core

import (
	"image/jpeg"
)

//-------------------------------------------------
// JPEG_QUALITY
//
// thumbnails (small/medium/large) generated for images being processed
// are always encoded as JPEG (see resizeImage() in gf_images_transformer.go).
// callers can now optionally pass an explicit JPEG quality value (1-100)
// all the way down the pipeline (TransformImage/TransformProcessImage/
// CreateThumbnails/resizeImage). if not explicitly set (nil) the default
// below is used, preserving prior behavior.

// JPEGqualityDefaultInt - matches the implicit Go stdlib default
// (image/jpeg's jpeg.Encode() with nil *jpeg.Options uses jpeg.DefaultQuality).
// used when no explicit quality is passed in by a caller.
const JPEGqualityDefaultInt = jpeg.DefaultQuality // 75

// JPEGqualityHighInt - a higher-quality/lower-compression setting, producing
// sharper thumbnails at the cost of larger file sizes. intended for content
// like infographics/plots/diagrams where the default JPEG compression
// introduces visible fuzziness/artifacts around sharp edges and text.
const JPEGqualityHighInt = 95

//-------------------------------------------------
