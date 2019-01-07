package gf_images_lib

import (
	"github.com/globalsign/mgo"
	"gf_core"
)
//--------------------------------------------------
func db_index__init(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	//---------------------
	//FIELDS  - "t","flows_names_lst","origin_url_str"
	//QUERIES - flows_db__images_exist() issues these queries

	doc_type__index := mgo.Index{
		Key:        []string{"t","flows_names_lst","origin_url_str"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := p_runtime_sys.Mongodb_coll.EnsureIndex(doc_type__index)
	if err != nil {
		gf_err := gf_core.Error__create(`failed to create db index on fields - {"t","flows_names_lst","origin_url_str"}`,
			"mongodb_ensure_index_error",nil,err,"gf_images_lib",p_runtime_sys)
		return gf_err
	}

	//DEPRECATED!! - flow_name_str field is deprecated in favor of flows_names_lst, 
	//               but some of the old img records still use it and havent been migrated yet. 
	//               so we're creating this index until migration is complete, and then 
	//               it should be removed.
	//FIELDS  - "t","flow_name_str","origin_url_str"
	//QUERIES - flows_db__images_exist() issues these queries

	doc_type__index = mgo.Index{
		Key:        []string{"t","flow_name_str","origin_url_str"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = p_runtime_sys.Mongodb_coll.EnsureIndex(doc_type__index)
	if err != nil {
		gf_err := gf_core.Error__create(`failed to create db index on fields - {"t","flow_name_str","origin_url_str"}`,
			"mongodb_ensure_index_error",nil,err,"gf_images_lib",p_runtime_sys)
		return gf_err
	}
	//---------------------

	return nil

}