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
	p_log_fun('FUN_ENTER','gf_image_db.image_exists()')

	#-----------
	#ADD!! - use redis as a cache for mongo data
	if p_db_type_str == 'redis':
		key_str = 'img:%s'%(p_image_id_str)
		if p_db_context_map['redis_client'].exists(key_str):
			return True
		else:
			return False
	#-----------
	#MONGO
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
	p_log_fun('FUN_ENTER','gf_image_db.db_get()')
	
	#---------------
	#MONGO
	if p_db_type_str == 'mongo':
		
		#SCALING!! - image_exists() does a full query to mongo
		#            investigate further
		if image_exists(p_image_id_str,
						p_db_context_map,
						p_log_fun,
						p_db_type_str = 'mongo'):

			mongo_client        = p_db_context_map['mongodb_client']
			gf_images_coll      = mongo_client[p_mongo_db_name_str][p_mongo_coll_name_str]
			raw_image_info_dict = gf_images_coll.find({"id_str":p_image_id_str})[0]
			image_info_dict     = gf_image.deserialize(raw_image_info_dict,
													   p_log_fun)
		else:
			return None
	#---------------

	#create() - does verification and adt construction
	image_adt = gf_image.create(image_info_dict,
		                        p_db_context_map,
								p_log_fun)
	assert isinstance(image_adt,gf_image.Image_ADT)
	
	return image_adt
#---------------------------------------------------
def db_put(p_image_adt,
		p_db_context_map,
		p_log_fun,
		p_db_type_str         = 'mongo',
		p_mongo_db_name_str   = 'prod_db',
		p_mongo_coll_name_str = 'images'):
	p_log_fun('FUN_ENTER','gf_image_db.db_put()')
	assert isinstance(p_image_adt,gf_image.Image_ADT)
	
	image_info_dict = gf_image.serialize(p_image_adt,
									p_log_fun)
	#---------------
	#MONGO
	if p_db_type_str == 'mongo':
		mongo_client = p_db_context_map['mongodb_client']
		data_coll    = mongo_client[p_mongo_db_name_str][p_mongo_coll_name_str]

		#spec          - a dict specifying elements which must be present for a document to be updated
		#upsert = True - insert doc if it doesnt exist, else just update
		p_mongo_data_collection.update({'id_str':p_image_adt.id_str}, #spec
									image_info_dict, 
									upsert = True)
	#---------------
#---------------------------------------------------	
#CAUTION!! - this is a very expensive operation for large image DB's

#->:List<:Image_ADT>
def db_get_all(p_db_context_map,
		p_log_fun,
		p_db_type_str         = 'mongo',
		p_mongo_db_name_str   = 'prod_db',
		p_mongo_coll_name_str = 'images'):
	p_log_fun('FUN_ENTER','gf_image_db.db_get_all()')

	#----------
	#MONGO
	if p_db_type_str == 'mongo':
		True
	#----------