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

package gf_images_lib

//-------------------------------------------------
type Config struct {
	Flow_to_s3bucket_map map[string]string
}
//-------------------------------------------------
func Config__get() Config {

	flow_to_s3bucket_map := map[string]string{
		"general":   "gf--img",
		"discovered":"gf--img--discover",
	}

	config := Config{
		Flow_to_s3bucket_map:flow_to_s3bucket_map,
	}

	return config
}