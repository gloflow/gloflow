// SPDX-License-Identifier: GPL-2.0
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

package gf_images_jobs_client

import (
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------

// called "expected" because jobs are long-running processes, and they might fail at various stages
// of their processing. in that case some of these result values will be satisfied, others will not.
type GFjobExpectedOutput struct {
	Image_id_str                      gf_images_core.GFimageID `json:"image_id_str"`
	Image_source_url_str              string                   `json:"image_source_url_str"`
	Thumbnail_small_relative_url_str  string                   `json:"thumbnail_small_relative_url_str"`
	Thumbnail_medium_relative_url_str string                   `json:"thumbnail_medium_relative_url_str"`
	Thumbnail_large_relative_url_str  string                   `json:"thumbnail_large_relative_url_str"`
}
