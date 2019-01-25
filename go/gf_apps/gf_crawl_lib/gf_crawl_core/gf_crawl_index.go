/*
GloFlow media management/publishing system
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_crawl_core

import (
	"time"
	"fmt"
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/olivere/elastic"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
type Gf_index__query_run struct {
	Id                   bson.ObjectId `bson:"_id,omitempty"`
	Id_str               string        `bson:"id_str"`
	T_str                string        `bson:"t"` //"index__query_run"
	Run_time_milisec_int int64         `bson:"run_time_milisec_int"`
	Hits_total_int       int64         `bson:"hits_total_int"`
	Hits_scores_lst      []float64     `bson:"hits_scores_lst"`
	Hits_score_max_f     float64       `bson:"hits_score_max_f"`
	Hits_urls_lst        []string      `bson:"hits_urls_lst"`
}
//--------------------------------------------------
func index__get_stats(p_runtime *Gf_crawler_runtime, p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime.Esearch_client.IndexStats("gf_crawl_pages")
}
//--------------------------------------------------
func Index__query(p_term_str string,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_index.Index__query()")


	//ADD!! - use termquery result relevance
	//      - a terms query would incorporate the percentage of terms that were found
	//      - https://www.elastic.co/guide/en/elasticsearch/guide/current/relevance-intro.html

	index_name_str := "gf_crawl_pages"
	field_name_str := "page_text_str"
	term_query     := elastic.NewTermQuery(field_name_str,p_term_str)

	ctx := context.Background()

	query_result,err := p_runtime.Esearch_client.Search().
	    Index(index_name_str).                                   //search in index "twitter"
	    Query(term_query).                                       //specify the query
	    Highlight(elastic.NewHighlight().Field(field_name_str)). //result HIGHLIGHTING
	    Sort("user", true).                                      //sort by "user" field, ascending
	    From(0).Size(10).                                        //take documents 0-9
	    Pretty(true).                                            //pretty print request and response JSON
	    Do(ctx)                                                  //execute

	if err != nil {
		gf_err := gf_core.Error__create("failed to issue a elasticsearch index query - "+p_term_str,
			"elasticsearch_query_index",
			&map[string]interface{}{
				"term_str":      p_term_str,
				"index_name_str":index_name_str,
				"field_name_str":field_name_str,
			},
			err, "gf_crawl_lib", p_runtime_sys)
		return gf_err
	}

	query_run_time_milisec_int := query_result.TookInMillis

	//----------------
	//HITS
	var search_hits *elastic.SearchHits = query_result.Hits
	total_hits_int                     := search_hits.TotalHits
	hits_score_max_f                   := search_hits.MaxScore
	hits_lst                           := search_hits.Hits


	hits_scores_lst := []float64{}
	hits_urls_lst   := []string{}
	for _,search_hit := range hits_lst {

		//-------------------
		//relevance - the algorithm that we use to calculate how similar 
		//            the contents of a full-text field are to a full-text query string.
		//          - standard similarity algorithm used in Elasticsearch is known as 
		//            term frequency/inverse document frequency, or TF/IDF
		//          - Individual queries may combine the TF/IDF score with other 
		//            factors such as the term proximity in phrase queries, 
		//            or term similarity in fuzzy queries
		//          
		//          - "_score" - field in results
		hit_score_f    := search_hit.Score
		hits_scores_lst = append(hits_scores_lst,*hit_score_f)
		//-------------------

		var highlights_map map[string][]string = search_hit.Highlight
		hit_explanation_str                   := search_hit.Explanation.Description

		fmt.Println(highlights_map)
		fmt.Println(hit_explanation_str)

		//result_doc_source__json_str := string(search_hit.Source)
		var hit__doc_fields_map map[string]interface{} = search_hit.Fields
		hit__url_str                                  := hit__doc_fields_map["Url_str"].(string)

		hits_urls_lst = append(hits_urls_lst,hit__url_str)
	}
	//----------------

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("crawler_page_img:%f",creation_unix_time_f)

	query_run := &Gf_index__query_run{
		Id_str:              id_str,
		T_str:               "index__query_run",
		Run_time_milisec_int:query_run_time_milisec_int,
		Hits_total_int:      total_hits_int,
		Hits_scores_lst:     hits_scores_lst,
		Hits_score_max_f:    *hits_score_max_f,
		Hits_urls_lst:       hits_urls_lst,
	}

	err = p_runtime_sys.Mongodb_coll.Insert(query_run)
	if err != nil {
		gf_err := gf_core.Error__create("failed to insert a index__query_run into mongodb for a elasticsearch index query",
			"mongodb_insert_error",
			&map[string]interface{}{
				"term_str":      p_term_str,
				"total_hits_int":total_hits_int,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}

	return nil
}
//--------------------------------------------------
func index__add_to__of_url_fetch(p_url_fetch *Gf_crawler_url_fetch,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_index.index__add_to__of_url_fetch()")

	index_name_str     := "gf_crawl_pages"
	es_record_type_str := "crawl_page"
	ctx                := context.Background()

	_, err := p_runtime.Esearch_client.Index().
		Index(index_name_str).
		Type(es_record_type_str).
		BodyJson(p_url_fetch).
		//Refresh(true). //refresh this index after completing this Index() operation
		Do(ctx)
	if err != nil {
		err_msg_str := fmt.Sprintf("failed to add/index a url_fetch record (es type - %s) to the elasticsearch index - %s", es_record_type_str, index_name_str)
		gf_err := gf_core.Error__create(err_msg_str,
			"elasticsearch_add_to_index",
			&map[string]interface{}{
				"url_fetch_url_str": p_url_fetch.Url_str,
				"index_name_str":    index_name_str,
				"es_record_type_str":es_record_type_str,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}
	return nil
}