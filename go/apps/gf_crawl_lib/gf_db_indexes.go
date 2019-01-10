package gf_crawl_lib

import (
	"github.com/globalsign/mgo"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_crawl_lib/gf_crawl_core"
)
//--------------------------------------------------
func db_index__init(p_runtime *gf_crawl_core.Crawler_runtime,
			p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {



	//---------------------
	//FIELDS - "t"
	doc_type__index := mgo.Index{
	    Key:        []string{"t"}, //all stat queries first match on "t"
	    Unique:     true,
	    DropDups:   true,
	    Background: true,
	    Sparse:     true,
	}

	err := p_runtime_sys.Mongodb_coll.EnsureIndex(doc_type__index)
	if err != nil {
		gf_err := gf_core.Error__create(`failed to create db index on fields - {"t"}`,
			"mongodb_ensure_index_error",nil,err,"gf_crawl_lib",p_runtime_sys)
		return gf_err
	}
	//---------------------
	//FIELDS - "t","hash_str"
	doc_type__index = mgo.Index{
	    Key:        []string{"t","hash_str"}, //all stat queries first match on "t"
	    Unique:     true,
	    DropDups:   true,
	    Background: true,
	    Sparse:     true,
	}

	err = p_runtime_sys.Mongodb_coll.EnsureIndex(doc_type__index)
	if err != nil {
		gf_err := gf_core.Error__create(`failed to create db index on fields - {"t","hash_str"}`,
			"mongodb_ensure_index_error",nil,err,"gf_crawl_lib",p_runtime_sys)
		return gf_err
	}
	//---------------------


	return nil

}