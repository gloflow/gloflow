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

import os, sys
modd_dir = os.path.abspath(os.path.dirname(os.path.abspath(__file__)))

from gf_core import gf_core_error

#---------------------------------------------------
#->:Map(new_image_info_map)
def verify_image_info(p_image_info_map,
	p_db_context_map,
	p_log_fun):
	p_log_fun('INFO', 'p_image_info_map:%s'%(p_image_info_map))
	assert isinstance(p_image_info_map, dict)
	
	max_title_characters = 100
	#-------------------
	# 'id_str' - is None if image_info_dict comes from the outside of the system
	#            and the ID has not yet been assigned. 
	#            ID is determined if image_info_dict comes from the DB or someplace
	#            else within the system
	
	id_str = p_image_info_map.get('id_str', None)

	#-------------------
	# TITLE
	
	title_str = None
	try:
		assert p_image_info_map.has_key('title_str')
		title_str = p_image_info_map['title_str']
		
		assert isinstance(p_image_info_map['title_str'], str)
		assert len(p_image_info_map['title_str']) < max_title_characters

	except Exception as e:
		msg_str = '''title_str [%s] is missing, invalid, or has more then max_length [%s] characters'''%(title_str, max_title_characters)
		gf_core_error.handle_exception(e,
			msg_str, #p_formated_msg_str,
			(),      #p_surrounding_context_attribs_tpl,

			p_db_context_map,
			p_log_fun,
			p_app_name_str = 'gf_image')
		
		raise Exception(msg_str)
	
	#-------------------
	# FORMAT 
	
	format_str = None
	try:
		assert p_image_info_map.has_key('format_str')
		normalized_format_str = check_image_format(p_image_info_map['format_str'].lower(), p_log_fun)
		format_str = normalized_format_str
			
	except Exception as e:
		msg_str = '''format_str [%s] is missing or is invalid'''%(format_str)
		
		gf_core_error.handle_exception(e,
			msg_str,
			(), # p_surrounding_context_attribs_tpl,

			p_db_context_map,
			p_log_fun,
			p_app_name_str = 'gf_image')
		
		raise Exception(msg_str)
	
	#-------------------
	# WIDTH/HEIGHT
	
	width_str  = None
	height_str = None
	
	try:
		assert p_image_info_map.has_key('width_str')
		width_str = p_image_info_map['width_str']
		
		assert p_image_info_map.has_key('height_str')
		height_str = p_image_info_map['height_str']
		
		assert width_str.isdigit()
		assert height_str.isdigit()
		       
	except Exception as e:
		msg_str = '''width_str/height [%s/%s] is missing or is invalid'''%(width_str, height_str)
		
		gf_core_error.handle_exception(e,
			msg_str,
			(), # p_surrounding_context_attribs_tpl,

			p_db_context_map,
			p_log_fun,
			p_app_name_str = 'gf_image')
		raise Exception(msg_str)
	
	#-------------------
	# ORIGIN URL
	
	origin_url_str = None
	
	try:
		assert p_image_info_map.has_key('origin_url_str')
		origin_url_str = p_image_info_map['origin_url_str']
	except Exception as e:
		msg_str = '''origin_url_str [%s] is missing or is invalid'''%(origin_url_str)
		
		gf_core_error.handle_exception(e,
			msg_str,
			(), # p_surrounding_context_attribs_tpl

			p_db_context_map,
			p_log_fun,
			p_app_name_str = 'gf_image')
		
		raise Exception(msg_str)
	
	#-------------------
	original_file_internal_uri_str = p_image_info_map.get('original_file_internal_uri_str',None)
	thumbnail_small_url_str        = p_image_info_map.get('thumbnail_small_url_str' ,None)
	thumbnail_medium_url_str       = p_image_info_map.get('thumbnail_medium_url_str',None)
	thumbnail_large_url_str        = p_image_info_map.get('thumbnail_large_url_str' ,None)

	#-------------------
	
	new_image_info_map = {
		'id_str': id_str,
		'title_str': title_str,
		
		'origin_url_str': origin_url_str,
		'original_file_internal_uri_str': original_file_internal_uri_str,
		
		'thumbnail_small_url_str':  thumbnail_small_url_str,
		'thumbnail_medium_url_str': thumbnail_medium_url_str,
		'thumbnail_large_url_str':  thumbnail_large_url_str,
		
		'format_str': format_str,
		'width_str':  width_str,
		'height_str': height_str,
		
		'dominant_color_hex_str':p_image_info_map['dominant_color_hex_str']
	}
	
	return new_image_info_map

#---------------------------------------------------	
#->:String(normalized_format_str)

def check_image_format(p_format_str, p_log_fun):
	
	normalized_format_str = None
	
	assert isinstance(p_format_str, str)
	assert p_format_str == 'jpeg' or \
		p_format_str == 'jpg'  or \
		p_format_str == 'gif'  or \
		p_format_str == 'png'
				 
	# normalize "jpg" (variation on jpeg) to "jpeg"
	if p_format_str == 'jpg':
		normalized_format_str = 'jpeg'
	else:
		normalized_format_str = p_format_str
		
	return normalized_format_str