package gf_crawl_lib

//--------------------------------------------------
func Get_all_crawlers() map[string]Crawler {

	crawlers_map := map[string]Crawler{
		"rhubarbes":Crawler{
			Name_str:     "gloflow.com",
			Start_url_str:"http://gloflow.com/",
		},
	}

	return crawlers_map
}


