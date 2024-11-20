# GloFlow application and media management/publishing platform
# Copyright (C) 2024 Ivan Trajkovic
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
from icecream import ic
import pandas as pd
import gloflow as gf

table_name__extern_load_str = "gf_extern_load"

#---------------------------------------------------------------------------------
def get_load_type_counts(p_db_client):
	query_str = f'''
		SELECT load_type, COUNT(*)
		FROM {table_name__extern_load_str}
		GROUP BY load_type
		ORDER BY load_type;
	'''
	df = pd.read_sql(query_str, p_db_client)
	return df

#---------------------------------------------------------------------------------
def relate_observations(p_target_observation_id_int,
	p_related_observations_ids_lst,
	p_db_client):
	assert(isinstance(p_target_observation_id_int, int))
	assert(isinstance(p_related_observations_ids_lst, list))
	for o in p_related_observations_ids_lst:
		assert(isinstance(o, int))

	cur = p_db_client.cursor()

	query_str = f'''
		UPDATE {table_name__extern_load_str}
		SET related_observations = %s
		WHERE id = %s
	'''

	cur.execute(query_str, 
		(
			p_related_observations_ids_lst,
			p_target_observation_id_int
		))
	
	p_db_client.commit()
	cur.close()

#---------------------------------------------------------------------------------
# GET_GROUP

def get_group(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_group_id_str,
	p_db_client):
	
	cur = p_db_client.cursor()
	query_str = f'''
		SELECT
			id,
			url,
			resp_cache_file_path,
			meta_map
		FROM {table_name__extern_load_str}
		WHERE
			load_type     = %s AND
			part_key      = %s AND
			source_domain = %s AND
			group_id      = %s

		ORDER BY fetch_datetime DESC;
	'''

	params_lst = [
		p_load_type_str,
		p_part_key_str,
		p_source_domain_str,
		p_group_id_str
	]

	cur.execute(query_str, params_lst)
	results_lst = cur.fetchall()

	if results_lst is None:
		return None
	else:
		ic(results_lst)

		loads_lst = []
		for row in results_lst:
			load_map = {
				"id":                   row[0],
				"url":                  row[1],
				"resp_cache_file_path": row[2],
				"meta_map":             row[3]
			}
			loads_lst.append(load_map)

		return loads_lst

#---------------------------------------------------------------------------------
# GET_LATEST_GROUP_ID

def get_latest_group_id(p_load_type_str,
	p_source_domain_str,
	p_db_client):

	cur = p_db_client.cursor()
	query_str = f'''
		SELECT group_id
		FROM {table_name__extern_load_str}
		WHERE
			load_type     = %s AND
			source_domain = %s AND
			group_id IS NOT NULL

		ORDER BY fetch_datetime DESC
	'''

	cur.execute(query_str, [
		p_load_type_str,
		p_source_domain_str
	])
	result_tpl = cur.fetchone()
	group_id_str = result_tpl[0]

	return group_id_str

#---------------------------------------------------------------------------------
# GET_LATEST_LOAD

def get_latest(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_db_client,
	p_group_id_str=None):

	where_sql_str = "load_type = %s AND part_key = %s AND source_domain = %s"
	params_lst = [
		p_load_type_str,
		p_part_key_str,
		p_source_domain_str
	]
	if p_group_id_str is not None:
		where_sql_str += " AND group_id = %s"
		params_lst.append(p_group_id_str)

	cur = p_db_client.cursor()
	query_str = f'''
		SELECT
			id,
			url,
			resp_cache_file_path,
			meta_map
		FROM {table_name__extern_load_str}
		WHERE {where_sql_str}
		ORDER BY fetch_datetime DESC
		LIMIT 1;
	'''

	cur.execute(query_str, params_lst)
	result_tpl = cur.fetchone()

	if result_tpl is None:
		return None
	else:
		ic(result_tpl)
		load_map = {
			"id":                   result_tpl[0],
			"url":                  result_tpl[1],
			"resp_cache_file_path": result_tpl[2],
			"meta_map":             result_tpl[3]
		}
		return load_map

#---------------------------------------------------------------------------------
# INSERT
def insert(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_cache_s3_key_str,
	p_db_client,
	p_meta_map     = {},
	p_url_str      = None,
	p_group_id_str = None):
	assert(isinstance(p_meta_map, dict))

	cur = p_db_client.cursor()

	query_str = f'''INSERT INTO {table_name__extern_load_str} (
			load_type,
			part_key,
			url,
			resp_cache_file_path,
			meta_map,
			source_domain,
			group_id
		)
		VALUES (%s, %s, %s, %s, %s, %s, %s)
		RETURNING id
	'''

	cur.execute(query_str, 
		(
			p_load_type_str,
			p_part_key_str,
			p_url_str,
			p_cache_s3_key_str,
			json.dumps(p_meta_map),
			p_source_domain_str,
			p_group_id_str
		))

	id_int = cur.fetchone()[0]
	p_db_client.commit()
	cur.close()

	return id_int

#---------------------------------------------------------------------------------
# INIT
def init(p_db_client):

	cur = p_db_client.cursor()
	
	if not gf.db.table_exists(table_name__extern_load_str, cur):
		
		# source_domain  - domain on which this ad was discovered
		# fetch_datetime - time when the GF system stored this item

		sql_str = f"""
			CREATE TABLE {table_name__extern_load_str} (
			
				id SERIAL PRIMARY KEY,
				
				-- what type of extern loading is done
				-- model, model_variant, etc.
				load_type VARCHAR(255),

				-- -----------------------
				-- PARTITION_KEY
				-- string representing partition key
				part_key VARCHAR(255),

				-- -----------------------
				-- URL
				-- for html data this is the URL of the page, for json
				-- returns it might be the URL of the API endpoint.
				-- for other types of data it might be None.

				url VARCHAR(1000),
				
				-- -----------------------
				-- RELATED_OBSERVATIONS
				-- these are observations that this observations depends on or is related to.
				-- there can be multiple observations that can be used by this observation.

				related_observations INT[],

				-- -----------------------
				-- RESPONSE
				resp_cache_file_path TEXT,

				-- -----------------------
				-- META
				-- various metadata that can be attached to a load event.
				meta_map JSONB,
				
				-- -----------------------
				-- GROUP_ID
				-- allows for grouping of observations.
				group_id VARCHAR(255),

				-- -----------------------

				source_domain  VARCHAR(255),
				fetch_datetime TIMESTAMP DEFAULT NOW()
			);
		"""
		cur.execute(sql_str)
		p_db_client.commit()