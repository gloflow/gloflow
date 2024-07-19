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

import json
import gf_image

#---------------------------------------------------
#->:Bool
def image_exists(p_image_id_str,
	p_db_context_map,
	p_log_fun,
	p_db_type_str         = 'mongo',
	p_mongo_db_name_str   = 'prod_db',
	p_mongo_coll_name_str = 'images'):

	#-----------
	# ADD!! - use redis as a cache for mongo data
	if p_db_type_str == 'redis':
		key_str = 'img:%s'%(p_image_id_str)
		if p_db_context_map['redis_client'].exists(key_str):
			return True
		else:
			return False

	#-----------
	# MONGO
	elif p_db_type_str == 'mongo':
		mongo_client = p_db_context_map['mongodb_client']
		data_coll    = mongo_client[p_mongo_db_name_str][p_mongo_coll_name_str] #db/collection name
		
		#find().limit() - drastically faster for existence checking then findOne()
		#{"_id":1}      - only return the "_id" field
		if data_coll.find({"id_str":p_image_id_str},{"_id":1}).limit(1).count() == 0: return False
		else                                                                        : return True
	
	#-----------

#---------------------------------------------------
#->:Image_ADT|None
def db_get(p_image_id_str,
	p_db_context_map,
	p_log_fun,
	p_db_type_str         = 'mongo',
	p_mongo_db_name_str   = 'prod_db',
	p_mongo_coll_name_str = 'images'):
	
	#---------------
	# DB
	if p_db_type_str == 'mongo':
		
		# SCALING!! - image_exists() does a full query to mongo
		#             investigate further
		if image_exists(p_image_id_str,
			p_db_context_map,
			p_log_fun,
			p_db_type_str = 'mongo'):

			mongo_client        = p_db_context_map['mongodb_client']
			gf_images_coll      = mongo_client[p_mongo_db_name_str][p_mongo_coll_name_str]
			raw_image_info_dict = gf_images_coll.find({"id_str":p_image_id_str})[0]
			image_info_dict     = gf_image.deserialize(raw_image_info_dict, p_log_fun)
		else:
			return None

	#---------------

	# create() - does verification and adt construction
	image_adt = gf_image.create(image_info_dict, p_db_context_map, p_log_fun)
	assert isinstance(image_adt, gf_image.Image_ADT)
	
	return image_adt

#---------------------------------------------------
def db_put(p_image_adt,
	p_db_context_map,
	p_log_fun,
	p_db_type_str         = 'mongo',
	p_mongo_db_name_str   = 'prod_db',
	p_mongo_coll_name_str = 'images'):
	assert isinstance(p_image_adt, gf_image.Image_ADT)
	
	image_info_map = gf_image.serialize(p_image_adt,
		p_log_fun)
	#---------------
	# DB
	if p_db_type_str == 'mongo':
		mongo_client = p_db_context_map['mongodb_client']
		data_coll    = mongo_client[p_mongo_db_name_str][p_mongo_coll_name_str]

		# spec          - a dict specifying elements which must be present for a document to be updated
		# upsert = True - insert doc if it doesnt exist, else just update
		data_coll.update({'id_str':p_image_adt.id_str}, # spec
			image_info_map, 
			upsert = True)
		
	#---------------
	
#---------------------------------------------------	
#->:List<:Image_ADT>
def db_get_all(p_db_context_map,
	p_log_fun,
	p_db_type_str         = 'mongo',
	p_mongo_db_name_str   = 'prod_db',
	p_mongo_coll_name_str = 'images'):

	#----------
	# MONGO
	if p_db_type_str == 'mongo':
		True
		
	#----------