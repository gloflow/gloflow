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

"""
extern_load 
	- helps with tracking and monitoring of external data loading.
	- tracks loads of external data as events in the DB.
	- allows for historic tracking of external data loading.
	- simple mechanism, py native.
	- results of loads are stored in the DB in string form
		- results can be html or json.
		- if json then stored in postgres jsonb type.
	- the partition key is stored for each load event.
		- parititon keys can be multidimensional
		- each dim separated by "__".
"""

import os, json
import boto3
from icecream import ic
import gloflow as gf

# stores data on loading/processing of external models info
# such as from automani or some other source.
table_name__extern_load_str = "gf_extern_load"
default_s3_bucket_name_str  = "gf"

#---------------------------------------------------------------------------------
# PUBLIC
#---------------------------------------------------------------------------------
# OBSERVE
# p_resp_store_file_path_str - allow external users of this API to determine the file path
# 	                           to be used for caching the response data. user knows best what
#                              naming structure makes sense for their domain.
#                              files are stored in some block storage, local or remote (s3, etc.).
#                              filepath should be relative.
# p_related_observations_ids_lst - list of observation ids that this observation depends on.
#                                  that dependance is usually in the form of a data dependency.

def observe(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_runtime_map,
	p_meta_map={},
	p_url_str=None,
	p_resp_data_map=None,
	p_cache_file_path_str=None,
	p_related_observations_ids_lst=None):
	
	assert(isinstance(p_load_type_str, str))
	if p_related_observations_ids_lst is not None:
		assert(isinstance(p_related_observations_ids_lst, list))

	if p_cache_file_path_str is None:
		store_result_bool = False
	else:
		store_result_bool = True

	
	#-----------------------
	# RESP_FILE_STORAGE
	cache_s3_key_str = None
	if store_result_bool:

		assert(isinstance(p_resp_data_map, dict))
		resp_data_str = json.dumps(p_resp_data_map)
			
		cache_s3_key_str = upload_cache(resp_data_str,
			p_load_type_str,
			p_cache_file_path_str,
			p_source_domain_str,
			p_runtime_map)

	#-----------------------
	# DB
	observation_id_str = db_insert(p_load_type_str,
		p_part_key_str,
		p_source_domain_str,
		cache_s3_key_str,
		p_runtime_map["db_client"],
		p_meta_map=p_meta_map,
		p_url_str=p_url_str,
		p_related_observations_ids_lst=p_related_observations_ids_lst)

	ic(observation_id_str)

	#-----------------------

	return observation_id_str

#---------------------------------------------------------------------------------
# RELATE
def relate(p_observation_id_str,
	p_related_observations_ids_lst,
	p_runtime_map):

	assert(isinstance(p_observation_id_str, str))
	assert(isinstance(p_related_observations_ids_lst, list))

	db_relate_observations(p_observation_id_str,
		p_related_observations_ids_lst,
		p_runtime_map["db_client"])
	
#---------------------------------------------------------------------------------
# GET_CACHE

def get_cache(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_runtime_map):
	assert(isinstance(p_load_type_str, str))

	#-----------------------
	# DB
	latest_load_map = db_get_latest_load(p_load_type_str,
		p_part_key_str,
		p_source_domain_str,
		p_runtime_map["db_client"])

	#-----------------------
	# S3

	ic(latest_load_map)
	
	bucket_name_str = p_runtime_map.get("s3_data_sink_bucket_str", default_s3_bucket_name_str)
	ic(bucket_name_str)

	s3 = boto3.client('s3')
	resp = s3.get_object(Bucket=bucket_name_str,
		Key=latest_load_map["resp_cache_file_path"])

	r = resp['Body']
	
	resp_str = resp['Body'].read().decode('utf-8')
	resp_map = json.loads(resp_str)
	assert(isinstance(resp_map, dict))

	#-----------------------
	
	observation_id_int = latest_load_map["id"]
	assert(isinstance(observation_id_int, int))

	return resp_map, observation_id_int

#---------------------------------------------------------------------------------
def init(p_db_client):

	db_init(p_db_client)

#---------------------------------------------------------------------------------
# CACHE
#---------------------------------------------------------------------------------
def upload_cache(p_data_str,
	p_load_type_str,
	p_cache_file_path_str,
	p_source_domain_str,
	p_runtime_map):
	assert(isinstance(p_data_str, str))
	assert(isinstance(p_load_type_str, str))
	bucket_name_str = p_runtime_map.get("s3_data_sink_bucket_str", "gf")
	ic(bucket_name_str)

	s3 = boto3.client('s3')

	#-----------------------
	file_path_norm_str = os.path.normpath(p_cache_file_path_str)
	dir_str            = f"gf/ext_load/{p_source_domain_str}/{p_load_type_str}"
	s3_key_str         = f'{dir_str}/{file_path_norm_str}'
	ic(s3_key_str)

	#-----------------------
	# S3
	s3.put_object(Bucket=bucket_name_str,
		Key=s3_key_str,
		Body=p_data_str)

	#-----------------------

	return s3_key_str

#---------------------------------------------------------------------------------
# DB
#---------------------------------------------------------------------------------
def db_relate_observations(p_target_observation_id_str,
	p_related_observations_ids_lst,
	p_db_client):
	assert(isinstance(p_target_observation_id_str, str))
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
			p_target_observation_id_str
		))
	
	p_db_client.commit()
	cur.close()

#---------------------------------------------------------------------------------
# GET_LATEST_LOAD

def db_get_latest_load(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
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
			
			load_type = %s
			AND
			part_key = %s
			AND
			source_domain = %s
		
		ORDER BY fetch_datetime DESC
		LIMIT 1;
	'''

	cur.execute(query_str,
		(
			p_load_type_str,
			p_part_key_str,
			p_source_domain_str
		))
	
	result_tpl = cur.fetchone()
	load_map = {
		"id":                   result_tpl[0],
		"url":                  result_tpl[1],
		"resp_cache_file_path": result_tpl[2],
		"meta_map":             result_tpl[3]
	}
	return load_map

#---------------------------------------------------------------------------------
# INSERT
def db_insert(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_cache_s3_key_str,
	p_db_client,
	p_meta_map={},
	p_url_str=None):
	assert(isinstance(p_meta_map, dict))

	cur = p_db_client.cursor()

	query_str = f'''INSERT INTO {table_name__extern_load_str} (
			load_type,
			part_key,
			url,
			resp_cache_file_path,
			meta_map,
			source_domain
		)
		VALUES (%s, %s, %s, %s, %s, %s)
		RETURNING id
	'''

	cur.execute(query_str, 
		(
			p_load_type_str,
			p_part_key_str,
			p_url_str,
			p_cache_s3_key_str,
			json.dumps(p_meta_map),
			p_source_domain_str
		))

	id_int = cur.fetchone()[0]
	p_db_client.commit()
	cur.close()

	return id_int

#---------------------------------------------------------------------------------
# INIT
def db_init(p_db_client):

	cur = p_db_client.cursor()
	
	if not gf.db_table_exists(table_name__extern_load_str, cur):
		
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

				source_domain  VARCHAR(255),
				fetch_datetime TIMESTAMP DEFAULT NOW()
			);
		"""
		cur.execute(sql_str)
		p_db_client.commit()