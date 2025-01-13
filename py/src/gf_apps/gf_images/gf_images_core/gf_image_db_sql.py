# GloFlow application and media management/publishing platform
# Copyright (C) 2025 Ivan Trajkovic
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

import json
from gf_core import gf_core_utils
from . import gf_image

#---------------------------------------------------
# p_mongo_source_bool - temporary, used during mongo->sql migration

def put_images(p_image_adt_lst,
	p_db_client,
	p_mongo_source_bool: bool = False):
	assert all(isinstance(img, gf_image.GFimage) for img in p_image_adt_lst), \
		"All items must be instances of gf_image.GFimage"

	table_name_str = "gf_images"
	cur = p_db_client.cursor()

	query_str = f'''INSERT INTO {table_name_str} (
			v,
			id,
			creation_time,
			user_id,

			client_type,
			title,
			flows_names,

			origin_url,
			origin_page_url,

			original_file_int_url,

			thumb_small_url,
			thumb_medium_url,
			thumb_large_url,

			format,
			width,
			height,

			dominant_color_hex,
			palette_colors_hex,

			meta_map,
			tags_lst,

			mongo_source
		)

		VALUES (
			%s, %s, %s, %s,
			%s, %s, %s,
			%s, %s,
			%s,
			%s, %s, %s,
			%s, %s, %s,
			%s, %s,
			%s, %s,
			%s
		)
	'''

	# prepare data for batch insertion
	insertion_data_lst = []
	for p_image_adt in p_image_adt_lst:
		
		sql_timestamp_str = gf_core_utils.unix_to_sql_timestamp(p_image_adt.creation_unix_time_f)
		insertion_data_lst.append((
			"0.1",  # v
			p_image_adt.id_str,
			sql_timestamp_str,
			p_image_adt.user_id_str,

			p_image_adt.client_type_str,
			p_image_adt.title_str,
			p_image_adt.flows_names_lst,

			p_image_adt.origin_url_str,
			p_image_adt.origin_page_url_str,

			p_image_adt.original_file_int_url_str,

			p_image_adt.thumb_small_url_str,
			p_image_adt.thumb_medium_url_str,
			p_image_adt.thumb_large_url_str,

			p_image_adt.format_str,
			p_image_adt.width_int,
			p_image_adt.height_int,

			p_image_adt.dominant_color_hex_str,
			p_image_adt.palette_colors_hex_lst,

			json.dumps(p_image_adt.meta_map),
			p_image_adt.tags_lst,

			p_mongo_source_bool
		))

	# batch insert
	cur.executemany(query_str, insertion_data_lst)
	
	p_db_client.commit()
	cur.close()

#---------------------------------------------------------------------------------
	
def check_images_exist(p_ids_lst, p_db_client):
	assert isinstance(p_ids_lst, list) and \
		all(isinstance(p, str) for p in p_ids_lst), "Input must be a list of strings"

	table_name_str = "gf_images"
	cur = p_db_client.cursor()


	# createa list with as many "%s" (sql var placeholder) as there are IDs in p_ids_lst
	# and then join it into a single string with ", " as separator
	ids_placeholder_str = ', '.join(['%s'] * len(p_ids_lst))

	query_str = f'''
		SELECT id
		FROM {table_name_str}
		WHERE id IN ({ids_placeholder_str})
	'''

	cur.execute(query_str, tuple(p_ids_lst))

	# get a list of ID's that already exist in the database	
	existing_ids_lst = cur.fetchall()
	existing_ids_set = set([row[0] for row in existing_ids_lst])

	# compose a list of bools, same length as p_ids_lst, that for each
	# ID indicates if its already present or not.
	exists_bool_lst = [id_str in existing_ids_set for id_str in p_ids_lst]

	cur.close()

	return exists_bool_lst