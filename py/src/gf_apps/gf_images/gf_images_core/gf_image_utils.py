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

import os
from urllib.parse import urlparse
from pathlib import Path
from PIL import Image
import requests

IMAGE_EXTENSIONS = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff']

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

#--------------------------------------------
def check_image_format(p_url_str):
	
	# do a simple file-path check if its an image
	is_image_bool, format_str = is_image_url(p_url_str)
	if is_image_bool:
		return True, format_str
	
	# if the extension doesn't indicate an image, check the Content-Type using a GET request
	try:
		# using HEAD request to minimize data transfer
		response = requests.head(p_url_str)
		if 'image' in response.headers.get('Content-Type', '').lower():
			
			# check the format based on the Content-Type
			content_type_str = response.headers['Content-Type'].split('/')[1]
			return True, content_type_str
		
	except requests.exceptions.RequestException:
		return False, None
	
	return False, None

#--------------------------------------------
def is_image_url(p_url_str):
	
	path = urlparse(p_url_str).path
	
	# check if the URL ends with a valid image file extension
	ext_str = Path(path).suffix.lower().strip(".")
	if ext_str in IMAGE_EXTENSIONS:
		
		format_str = ext_str
		return True, format_str
	else:
		return False, None