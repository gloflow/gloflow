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

import pprint
from gf_core import gf_core_error, gf_core_utils
from . import gf_image_utils

#---------------------------------------------------
def on_map(p_image_info_map,
	p_log_fun):
	assert isinstance(p_image_info_map, dict)
	
	# pprint.pprint(p_image_info_map)

	max_title_characters = 500
	#-------------------
	# 'id_str' - is None if image_info_dict comes from the outside of the system
	#            and the ID has not yet been assigned. 
	#            ID is determined if image_info_dict comes from the DB or someplace
	#            else within the system
	
	id_str = p_image_info_map.get('id_str', None)
	if id_str is not None:
		assert isinstance(id_str, str)
		assert id_str != ""

	#-------------------
	# CREATION TIME
	try:
		if 'creation_unix_time_f' in p_image_info_map.keys():
			assert isinstance(p_image_info_map['creation_unix_time_f'], float)

	except Exception as e:
		msg_str = '''creation_unix_time_f is invalid (not a float)'''
		gf_core_error.create(msg_str,
			"image_verification",
			{"id_str": id_str}, e, "gf_images",
			p_log_fun)
		raise Exception(msg_str)
	#-------------------
	# TITLE
	
	title_str = None
	try:
		assert 'title_str' in p_image_info_map.keys()
		title_str = p_image_info_map['title_str']
		
		assert isinstance(p_image_info_map['title_str'], str)
		assert len(p_image_info_map['title_str']) < max_title_characters

	except Exception as e:
		msg_str = '''title_str [%s] is missing, invalid, or has more then max_length [%s] characters'''%(title_str, max_title_characters)
		gf_core_error.create(msg_str,
			"image_verification",
			{"id_str": id_str}, e, "gf_images",
			p_log_fun)
		raise Exception(msg_str)

	#-------------------
	# FLOWS_NAMES

	try:

		assert 'flows_names_lst' in p_image_info_map.keys()
		flows_names_lst = p_image_info_map['flows_names_lst']
		assert isinstance(flows_names_lst, list)
		assert len(flows_names_lst) > 0

		for flow_name_str in flows_names_lst:
			assert isinstance(flow_name_str, str)
			assert flow_name_str != ""

	except Exception as e:
		msg_str = '''flows_names_lst [%s] is missing or is invalid'''
		gf_core_error.create(msg_str,
			"image_verification",
			{"id_str": id_str}, e, "gf_images",
			p_log_fun)
		raise Exception(msg_str)
	
	#-------------------
	# FORMAT 
	
	format_str = None
	try:
		assert "format_str" in p_image_info_map.keys()
		assert isinstance(p_image_info_map['format_str'], str)

		normalized_format_str = check_image_format(p_image_info_map['format_str'].lower(), p_log_fun)
		format_str = normalized_format_str
			
	except Exception as e:
		msg_str = '''format_str [%s] is missing or is invalid'''%(format_str)
		gf_core_error.create(msg_str,
			"image_verification",
			{"id_str": id_str}, e, "gf_images",
			p_log_fun,
			p_reraise_bool=True)
	
	#-------------------
	# WIDTH/HEIGHT
	
	try:
		assert 'width_int' in p_image_info_map.keys()
		width_int = p_image_info_map['width_int']
		
		assert 'height_int' in p_image_info_map.keys()
		height_int = p_image_info_map['height_int']
		
		assert isinstance(width_int, int)
		assert isinstance(height_int, int)
		       
	except Exception as e:
		msg_str = '''image width/height is missing or is invalid'''
		gf_core_error.create(msg_str,
			"image_verification",
			{"id_str": id_str}, e, "gf_images",
			p_log_fun,
			p_reraise_bool=True)
	
	#-------------------
	# ORIGIN URL
	
	origin_url_str = None
	try:
		assert 'origin_url_str' in p_image_info_map.keys()
		origin_url_str = p_image_info_map['origin_url_str']

		gf_core_utils.is_valid_url(origin_url_str)

	except Exception as e:
		msg_str = '''origin_url_str [%s] is missing or is invalid'''%(origin_url_str)
		gf_core_error.create(msg_str,
			"image_verification",
			{"id_str": id_str}, e, "gf_images",
			p_log_fun,
			p_reraise_bool=True)
	
	#-------------------
	# THUMBS
	thumbnail_small_url_str  = p_image_info_map.get('thumb_small_url_str',  None)
	thumbnail_medium_url_str = p_image_info_map.get('thumb_medium_url_str', None)
	thumbnail_large_url_str  = p_image_info_map.get('thumb_large_url_str',  None)

	#---------------------------------------------------
	def check_thumb(p_url_str):
		if p_url_str is not None:
			assert isinstance(p_url_str, str)
			assert p_url_str != ""

	#---------------------------------------------------
	
	check_thumb(thumbnail_small_url_str)
	check_thumb(thumbnail_medium_url_str)
	check_thumb(thumbnail_large_url_str)

	#-------------------

#---------------------------------------------------
def check_image_format(p_format_str, p_log_fun):
	
	normalized_format_str = None
	
	assert isinstance(p_format_str, str)
	assert p_format_str in gf_image_utils.IMAGE_EXTENSIONS	
				 
	# normalize "jpg" (variation on jpeg) to "jpeg"
	if p_format_str == 'jpg':
		normalized_format_str = 'jpeg'
	else:
		normalized_format_str = p_format_str
		
	return normalized_format_str