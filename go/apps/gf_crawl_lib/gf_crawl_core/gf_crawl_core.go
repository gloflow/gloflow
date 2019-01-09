package gf_crawl_core

import (
	"github.com/olivere/elastic"
	"gf_core"
)
//--------------------------------------------------
type Crawler_runtime struct {
	Events_ctx            *gf_core.Events_ctx
	Esearch_client        *elastic.Client
	S3_info               *gf_core.Gf_s3_info
	Cluster_node_type_str string
}