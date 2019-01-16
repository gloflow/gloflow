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

package gf_crawl_lib

//--------------------------------------------------
func Get_all_crawlers() map[string]Gf_crawler {

	crawlers_map := map[string]Gf_crawler{
		"gloflow.com":Gf_crawler{
			Name_str:     "gloflow.com",
			Start_url_str:"http://gloflow.com/",
		},
	}
	return crawlers_map
}


