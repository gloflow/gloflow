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

import hashlib
import json
import urlparse

import gf_image_verify

#---------------------------------------------------
class Image_ADT():
	def __init__(self, p_props_dict):
		
		self.id_str    = p_props_dict['id_str']
		self.title_str = p_props_dict['title_str']
		
		#---------------
		# when the image comes from an external url (as oppose to it being 
		# created internally, or uploaded directly to the system)
		self.origin_url_str = p_props_dict.get('origin_url_str', None)

		# actual path on the filesystem, of the fullsized image gotten from origin_url_str
		self.original_file_internal_uri_str = p_props_dict.get('original_file_internal_uri_str', None)
		
		#---------------
		# relative url's
		# '/images/image_name.*'
		
		self.thumbnail_small_url_str  = p_props_dict.get('thumbnail_small_url_str',  None)
		self.thumbnail_medium_url_str = p_props_dict.get('thumbnail_medium_url_str', None)
		self.thumbnail_large_url_str  = p_props_dict.get('thumbnail_large_url_str',  None)
		
		#---------------
		self.format_str = p_props_dict['format_str'] #"jpeg"|"png"|"gif"
		self.width_str  = p_props_dict.get('width_str',  None)
		self.height_str = p_props_dict.get('height_str', None)
		
		#---------------
		self.dominant_color_hex_str = p_props_dict.get('dominant_color_hex_str', None)
		
		self.tags_lst = p_props_dict.get('tags_lst', None)
		
#---------------------------------------------------
#->:Image_ADT
def create(p_image_info_map,
	p_db_context_map,
	p_log_fun):
	
	new_image_info_dict = gf_image_verify.verify_image_info(p_image_info_map,
		p_db_context_map,
		p_log_fun)

	if not p_image_info_map.has_key('id_str'):
		image_id_str = create_id(new_image_info_dict['title_str'],
			new_image_info_dict['format_str'],
			p_log_fun)
		new_image_info_dict['id_str'] = image_id_str
		
	image_adt = Image_ADT(new_image_info_dict)
	return image_adt

#---------------------------------------------------
def create_id_from_url(p_image_url_str, p_log_fun):
	p_log_fun('FUN_ENTER', 'gf_image.create_id_from_url()')
	
	# urlparse() - used so that any possible url query parameters are not used in the 
	#              os.path.basename() result
	image_path_str          = urlparse.urlparse(p_image_url_str).path
	image_file_name_str     = os.path.basename(image_path_str)
	image_title_str,ext_str = os.path.splitext(image_file_name_str)
	normalized_ext_str      = gf_image_verify.check_image_format_str(ext_str.lower().strip('.'), p_log_fun)
	
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
	p_log_fun('FUN_ENTER', 'gf_image.create_id()')
	assert isinstance(p_image_path_str,basestring)
	
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
def serialize(p_image_adt, p_log_fun):

	image_info_map = {
		'id_str':    p_image_adt.id_str,
		'title_str': p_image_adt.title_str,
		
		'origin_url_str':                 p_image_adt.origin_url_str,
		'original_file_internal_uri_str': p_image_adt.original_file_internal_uri_str,
		
		'thumbnail_small_url_str':  p_image_adt.thumbnail_small_url_str,
		'thumbnail_medium_url_str': p_image_adt.thumbnail_medium_url_str,
		'thumbnail_large_url_str':  p_image_adt.thumbnail_large_url_str,
		
		'format_str': p_image_adt.format_str,
		'width_str':  p_image_adt.width_str,
		'height_str': p_image_adt.height_str,
		
		'dominant_color_hex_str': p_image_adt.dominant_color_hex_str,
		'tags_lst':               p_image_adt.tags_lst
	}
	
	return image_info_map

#---------------------------------------------------
#->:Map(image_info_map)
def deserialize(p_raw_image_info_map, p_log_fun):
	p_log_fun('FUN_ENTER', 'gf_image.deserialize()')
	assert isinstance(p_raw_image_info_map, dict)
	
	image_info_map = p_raw_image_info_map
	return image_info_map