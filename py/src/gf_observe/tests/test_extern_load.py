import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir
sys.path.append(f"{modd_str}/../..")

from icecream import ic

import gloflow as gf
from gf_extern_services.gf_aws import gf_aws_s3

#---------------------------------------------------------------------------------
def main():

	env_str	= "dev"

	ic(gf.version)

	# DB
	db_name_str = "gloflow"
	db_client = gf.db_init_client(db_name_str, env_str)
	gf.observe_init(db_client)



	part_dim_1_str = "dim1"
	part_dim_2_str = "dim2"
	part_dim_3_str = "dim3"

	domain_str = "test_domain.org"
	url_str = f"https://{domain_str}/page1"
	s3_bucket_name_str = "gf-jobs-dagster-sink--gf-core-dev-us-east-1"

	runtime_map = {
		"db_client": db_client,
		"s3_data_sink_bucket_str": s3_bucket_name_str
	}

	#-------------------
	# OBSERVE_EXTERN_LOAD
	
	part_key_str  = f"{part_dim_1_str}__{part_dim_2_str}__{part_dim_3_str}"
	result_s3_key_str = gf.observe_ext_load("load_type_A",
		part_key_str,
		domain_str,
		runtime_map,
		p_meta_map = {
			"some": "meta",
			"and":  "other"
		},
		p_url_str = url_str,
		p_resp_type_str = "html",
		p_resp_data_html_str = "<html>...</html>",
		p_resp_store_file_name_str = "test_resp_file.html")
	
	#-------------------


	assert(gf_aws_s3.file_exists(s3_bucket_name_str, result_s3_key_str))

	

#---------------------------------------------------------------------------------




main()