# GloFlow application and media management/publishing platform
# Copyright (C) 2019 Ivan Trajkovic
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
import hashlib
from . import gf_image_verify

#---------------------------------------------------
class GFimage():
	def __init__(self, p_img_map):
		
		self.id_str               = p_img_map['id_str']
		self.creation_unix_time_f = p_img_map.get('creation_unix_time_f', None)
		self.user_id_str          = p_img_map.get('user_id_str', None)


		self.client_type_str = p_img_map.get('client_type_str', None)
		self.title_str       = p_img_map['title_str']
		self.flows_names_lst = p_img_map.get('flows_names_lst', [])

		#---------------
		# when the image comes from an external url (as oppose to it being 
		# created internally, or uploaded directly to the system)
		self.origin_url_str = p_img_map.get('origin_url_str', None)

		# if the image is extracted from a page, this holds the page_url
		self.origin_page_url_str = p_img_map.get('origin_page_url_str', None)
		
		#---------------
		# internal URL (relative) of the original image, in its full size, unprocessed.
		# this is the image that is stored in the system, and is used to generate thumbs;
		# never served directly to the user.
		self.original_file_int_url_str = p_img_map.get('original_file_int_url_str')

		# THUMBS
		self.thumb_small_url_str  = p_img_map.get('thumb_small_url_str',  None)
		self.thumb_medium_url_str = p_img_map.get('thumb_medium_url_str', None)
		self.thumb_large_url_str  = p_img_map.get('thumb_large_url_str',  None)
		
		#---------------
		self.format_str = p_img_map['format_str'] #"jpeg"|"png"|"gif"
		self.width_str  = p_img_map.get('width_str',  None)
		self.height_str = p_img_map.get('height_str', None)
		
		#---------------
		self.dominant_color_hex_str = p_img_map.get('dominant_color_hex_str', None)
		self.pallete_str            = p_img_map.get('pallete_str', None)

		self.meta_map = p_img_map.get('meta_map', {})
		self.tags_lst = p_img_map.get('tags_lst', [])
		
#---------------------------------------------------
#->:GFimage

def load_adt(p_image_info_map,
	p_log_fun):
	
	gf_image_verify.on_map(p_image_info_map,
		p_log_fun)

	if not "id_str" in p_image_info_map.keys():
		image_id_str = create_id(p_image_info_map['title_str'],
			p_image_info_map['format_str'],
			p_log_fun)
		p_image_info_map['id_str'] = image_id_str
		
	gf_image = GFimage(p_image_info_map)
	return gf_image

#---------------------------------------------------
def create_id_from_url(p_image_url_str, p_log_fun):
	
	# urlparse() - used so that any possible url query parameters are not used in the 
	#              os.path.basename() result
	image_path_str           = urlparse.urlparse(p_image_url_str).path
	image_file_name_str      = os.path.basename(image_path_str)
	image_title_str, ext_str = os.path.splitext(image_file_name_str)
	normalized_ext_str       = gf_image_verify.check_image_format_str(ext_str.lower().strip('.'), p_log_fun)
	
	image_id_str = create_id(image_path_str,
		normalized_ext_str, # p_image_format_str,
		p_log_fun)
	return image_id_str

#---------------------------------------------------
# p_image_type_str - :String - 'jpeg'|'gif'|'png'

#->:String(image_id_str)
def create_id(p_image_path_str,
	p_image_format_str,
	p_log_fun):
	assert isinstance(p_image_path_str, str)
	
	assert p_image_format_str == 'jpeg' or \
		p_image_format_str == 'jpg'  or \
		p_image_format_str == 'gif'  or \
		p_image_format_str == 'png'

	m = hashlib.md5()
	m.update(p_image_path_str)
	m.update(p_image_format_str)
	
	#-------------------
	# hexdigest() - Like digest() except the digest is returned as a string of 
	#               double length, containing only hexadecimal digits. This may be used to 
	#               exchange the value safely in email or other non-binary environments.
	image_id_str = m.hexdigest()

	#-------------------
	
	return image_id_str

#---------------------------------------------------
#->:Map(image_info_map)
def serialize(p_gf_image, p_log_fun):

	image_info_map = {
		'id_str':    p_gf_image.id_str,
		'title_str': p_gf_image.title_str,
		
		'origin_url_str':                 p_gf_image.origin_url_str,
		'original_file_internal_uri_str': p_gf_image.original_file_internal_uri_str,
		
		'thumbnail_small_url_str':  p_gf_image.thumbnail_small_url_str,
		'thumbnail_medium_url_str': p_gf_image.thumbnail_medium_url_str,
		'thumbnail_large_url_str':  p_gf_image.thumbnail_large_url_str,
		
		'format_str': p_gf_image.format_str,
		'width_str':  p_gf_image.width_str,
		'height_str': p_gf_image.height_str,
		
		'dominant_color_hex_str': p_gf_image.dominant_color_hex_str,
		'tags_lst':               p_gf_image.tags_lst
	}
	
	return image_info_map

#---------------------------------------------------
#->:Map(image_info_map)

def deserialize(p_raw_image_info_map, p_log_fun):
	assert isinstance(p_raw_image_info_map, dict)
	
	image_info_map = p_raw_image_info_map
	return image_info_map