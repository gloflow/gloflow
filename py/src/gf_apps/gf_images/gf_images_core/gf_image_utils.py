# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA


from PIL import Image

#--------------------------------------------
def get_image_metadata(p_image_path_str):
	assert(os.path.isfile(p_image_path_str))
	
	with Image.open(p_image_path_str) as img:
		width, height = img.size
		img_type_str = img.format

	meta_map = {
		"width_int":  width,
		"height_int": height,
		"format_str": img_type_str
	}

	return meta_map