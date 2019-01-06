





import pymongo

#----------------------------------------------
#ADD!! - figure out a smarter way to pick the right hostport from p_host_port_lst,
#        instead of just picking the first element

def get_client(p_log_fun,
			p_host_port_lst = ['127.0.0.1:27017']):
	p_log_fun('FUN_ENTER','gf_core_mongodb.get_client()')

	host_str,port_str = p_host_port_lst[0].split(':')

	mongo_client = pymongo.MongoClient(host_str,int(port_str))
	return mongo_client