/*
GloFlow application and media management/publishing platform
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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/olivere/elastic"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------

type Gf_index__query_run struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               string        `bson:"id_str"`
	T_str                string        `bson:"t"` //"index__query_run"
	Run_time_milisec_int int64         `bson:"run_time_milisec_int"`
	Hits_total_int       int64         `bson:"hits_total_int"`
	Hits_scores_lst      []float64     `bson:"hits_scores_lst"`
	Hits_score_max_f     float64       `bson:"hits_score_max_f"`
	Hits_urls_lst        []string      `bson:"hits_urls_lst"`
}

//--------------------------------------------------

func index__get_stats(pRuntime *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) {
	pRuntime.EsearchClient.IndexStats("gf_crawl_pages")
}

//--------------------------------------------------

func IndexQuery(p_term_str string,
	pRuntime     *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {


	// ADD!! - use termquery result relevance
	//       - a terms query would incorporate the percentage of terms that were found
	//       - https://www.elastic.co/guide/en/elasticsearch/guide/current/relevance-intro.html

	index_name_str := "gf_crawl_pages"
	field_name_str := "page_text_str"
	term_query     := elastic.NewTermQuery(field_name_str,p_term_str)

	ctx := context.Background()

	query_result,err := pRuntime.EsearchClient.Search().
	    Index(index_name_str).                                   // search in index "twitter"
	    Query(term_query).                                       // specify the query
	    Highlight(elastic.NewHighlight().Field(field_name_str)). // result HIGHLIGHTING
	    Sort("user", true).                                      // sort by "user" field, ascending
	    From(0).Size(10).                                        // take documents 0-9
	    Pretty(true).                                            // pretty print request and response JSON
	    Do(ctx)                                                  // execute

	if err != nil {
		gf_err := gf_core.ErrorCreate("failed to issue a elasticsearch index query - "+p_term_str,
			"elasticsearch_query_index",
			map[string]interface{}{
				"term_str":       p_term_str,
				"index_name_str": index_name_str,
				"field_name_str": field_name_str,
			},
			err, "gf_crawl_lib", pRuntimeSys)
		return gf_err
	}

	query_run_time_milisec_int := query_result.TookInMillis

	//----------------
	// HITS
	var search_hits *elastic.SearchHits = query_result.Hits
	total_hits_int                     := search_hits.TotalHits
	hits_score_max_f                   := search_hits.MaxScore
	hits_lst                           := search_hits.Hits


	hits_scores_lst := []float64{}
	hits_urls_lst   := []string{}
	for _, search_hit := range hits_lst {

		//-------------------
		// relevance - the algorithm that we use to calculate how similar 
		//             the contents of a full-text field are to a full-text query string.
		//           - standard similarity algorithm used in Elasticsearch is known as 
		//             term frequency/inverse document frequency, or TF/IDF
		//           - Individual queries may combine the TF/IDF score with other 
		//             factors such as the term proximity in phrase queries, 
		//             or term similarity in fuzzy queries
		//          
		//           - "_score" - field in results
		hit_score_f    := search_hit.Score
		hits_scores_lst = append(hits_scores_lst, *hit_score_f)

		//-------------------

		var highlights_map map[string][]string = search_hit.Highlight
		hit_explanation_str                   := search_hit.Explanation.Description

		fmt.Println(highlights_map)
		fmt.Println(hit_explanation_str)

		//result_doc_source__json_str := string(search_hit.Source)
		var hit__doc_fields_map map[string]interface{} = search_hit.Fields
		hit__url_str                                  := hit__doc_fields_map["Url_str"].(string)

		hits_urls_lst = append(hits_urls_lst, hit__url_str)
	}
	
	//----------------

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("crawler_page_img:%f", creation_unix_time_f)

	query_run := &Gf_index__query_run{
		Id_str:               id_str,
		T_str:                "index__query_run",
		Run_time_milisec_int: query_run_time_milisec_int,
		Hits_total_int:       total_hits_int,
		Hits_scores_lst:      hits_scores_lst,
		Hits_score_max_f:     *hits_score_max_f,
		Hits_urls_lst:        hits_urls_lst,
	}

	coll_name_str := "gf_crawl"
	gf_err        := gf_core.MongoInsert(query_run,
		coll_name_str,
		map[string]interface{}{
			"term_str":           p_term_str,
			"total_hits_int":     total_hits_int,
			"caller_err_msg_str": "failed to insert a index__query_run into the DB for a elasticsearch index query",
		},
		ctx,
		pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//--------------------------------------------------

func indexAddToOfURLfetch(pURLfetch *GFcrawlerURLfetch,
	pRuntime    *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	indexNameStr    := "gf_crawl_pages"
	esRecordTypeStr := "crawl_page"
	ctx             := context.Background()

	_, err := pRuntime.EsearchClient.Index().
		Index(indexNameStr).
		Type(esRecordTypeStr).
		BodyJson(pURLfetch).
		// Refresh(true). //refresh this index after completing this Index() operation
		Do(ctx)
	if err != nil {
		errMsgStr := fmt.Sprintf("failed to add/index a url_fetch record (es type - %s) to the elasticsearch index - %s", esRecordTypeStr, indexNameStr)
		gfErr := gf_core.ErrorCreate(errMsgStr,
			"elasticsearch_add_to_index",
			map[string]interface{}{
				"url_fetch_url_str":  pURLfetch.Url_str,
				"index_name_str":     indexNameStr,
				"es_record_type_str": esRecordTypeStr,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return gfErr
	}
	return nil
}