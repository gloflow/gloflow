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

from gf_observe import gf_extern_load_db as db

# stores data on loading/processing of external models info
# such as from automani or some other source.
default_s3_bucket_name_str  = "gf"

#---------------------------------------------------------------------------------
# PUBLIC
#---------------------------------------------------------------------------------
# OBSERVE
# p_url_str - url from which the observation is coming, if relevant.
# p_resp_store_file_path_str - allow external users of this API to determine the file path
# 	                           to be used for caching the response data. user knows best what
#                              naming structure makes sense for their domain.
#                              files are stored in some block storage, local or remote (s3, etc.).
#                              filepath should be relative.
# p_related_observations_ids_lst - list of observation ids that this observation depends on.
#                                  that dependance is usually in the form of a data dependency.
# p_group_id_str - user can supply assign observations to a group, if needed.
#                  useful if multiple observations are related and need to be grouped.

def observe(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_runtime_map,
	p_group_id_str        = None,
	p_meta_map            = {},
	p_url_str             = None,
	p_resp_data_map       = None,
	p_cache_file_path_str = None):
	
	assert(isinstance(p_load_type_str, str))

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
	observation_id_str = db.insert(p_load_type_str,
		p_part_key_str,
		p_source_domain_str,
		cache_s3_key_str,
		p_runtime_map["db_client"],
		p_meta_map     = p_meta_map,
		p_url_str      = p_url_str,
		p_group_id_str = p_group_id_str)

	ic(observation_id_str)

	#-----------------------
	# MONITORING
	if "monitoring_fun" in p_runtime_map:
		assert(callable(p_runtime_map["monitoring_fun"]))
		p_runtime_map["monitoring_fun"](p_event_map={
				"type":  "gf.observe.extern_load",
				"level": "info",

				# APP_LEVEL
				"load_type":     p_load_type_str,
				"source_domain": p_source_domain_str,
			})

	#-----------------------
		
	return observation_id_str

#---------------------------------------------------------------------------------
# RELATE
def relate(p_observation_id_int,
	p_related_observations_ids_lst,
	p_runtime_map):
	
	assert(isinstance(p_observation_id_int, int))
	assert(isinstance(p_related_observations_ids_lst, list))

	print(f"relate observation {p_observation_id_int} to observations {p_related_observations_ids_lst}")
	
	db.relate_observations(p_observation_id_int,
		p_related_observations_ids_lst,
		p_runtime_map["db_client"])

#---------------------------------------------------------------------------------
# GET_CACHED

def get_cached(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_runtime_map):
	assert(isinstance(p_load_type_str, str))
	assert(isinstance(p_part_key_str, str))
	assert(isinstance(p_runtime_map, dict))
	assert(isinstance(p_source_domain_str, str) or p_source_domain_str is None)

	#-----------------------
	# DB
	latest_load_map = db.get_latest(p_load_type_str,
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
	info_map = json.loads(resp_str)
	assert(isinstance(info_map, dict))

	#-----------------------
	
	observation_id_int = latest_load_map["id"]
	assert(isinstance(observation_id_int, int))

	return info_map, observation_id_int

#---------------------------------------------------------------------------------
# GET_CACHED_GROUP
def get_cached_group(p_load_type_str,
	p_part_key_str,
	p_source_domain_str,
	p_runtime_map):

	#-----------------------
	# DB
	latest_group_id_str = db.get_latest_group_id(p_load_type_str,
		p_source_domain_str,
		p_runtime_map["db_client"])

	#-----------------------

	ic(latest_group_id_str)



	loads_lst = db.get_group(p_load_type_str,
		p_part_key_str,
		p_source_domain_str,
		latest_group_id_str,
		p_runtime_map["db_client"])
	assert(isinstance(loads_lst, list))

	ic(loads_lst)


	#-----------------------
	# S3

	bucket_name_str = p_runtime_map.get("s3_data_sink_bucket_str", default_s3_bucket_name_str)
	ic(bucket_name_str)

	s3 = boto3.client('s3')

	#---------------------------------------------------------------------------------
	def get_file(p_load_map):
		
		resp = s3.get_object(Bucket=bucket_name_str,
			Key=p_load_map["resp_cache_file_path"])

		r = resp['Body']
		
		resp_str = resp['Body'].read().decode('utf-8')
		info_map = json.loads(resp_str)
		assert(isinstance(info_map, dict))
		return info_map
	
	#---------------------------------------------------------------------------------
	
	
	infos_lst = []
	for l_map in loads_lst:
		info_map = get_file(l_map)
		infos_lst.append(info_map)

	#-----------------------
	
	return infos_lst

#---------------------------------------------------------------------------------
def init(p_db_client):

	db.init(p_db_client)

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