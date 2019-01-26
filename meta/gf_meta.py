

import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))
#-------------------------------------------------------------
def get():

    meta_map = {
        'build_info_map':{
            #-------------
            #MAIN
            #GF_IMAGES
            'gf_images':{
                'version_str':         '0.7.3.7',
                'go_path_str':         '%s/../go/apps/gf_images'%(cwd_str),
                'go_output_path_str':  '%s/../bin/gf_apps/gf_images/gf_images_service'%(cwd_str),
                'service_name_str':    'gf_images_service',
                'service_base_dir_str':'%s/../bin/gf_apps/gf_images'%(cwd_str),
            },
            
            #LIB
            #GF_IMAGES_LIB
            'gf_images_lib':{
                'go_path_str':               '%s/../go/apps/gf_images_lib'%(cwd_str),
                'test_data_to_serve_dir_str':'%s/../go/apps/gf_images_lib/tests_data'%(cwd_str), #for tests serve data over http from this dir
            },
            #-------------
            #MAIN
            #GF_ANALYTICS
            'gf_analytics':{
                'version_str':         '0.7.3.16',
                'go_path_str':         '%s/../go/gf_apps/gf_analytics'%(cwd_str),
                'go_output_path_str':  '%s/../bin/gf_apps/gf_analytics/gf_analytics_service'%(cwd_str),
                'service_name_str':    'gf_analytics_service',
                'service_base_dir_str':'%s/../bin/gf_apps/gf_analytics'%(cwd_str),
                'copy_to_dir_lst':[
                    ('%s/../go/gf_stats/py/cli_stats.py'%(cwd_str),                                                    '%s/../bin/gf_apps/gf_analytics/py'%(cwd_str)),
                    ('%s/../go/gf_core/py/stats/gf_errors__counts_by_day.py'%(cwd_str),                                '%s/../bin/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_imgs__counts_by_day.py'%(cwd_str),           '%s/../bin/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__counts_by_day.py'%(cwd_str), '%s/../bin/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__null_breakdown.py'%(cwd_str),'%s/../bin/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__per_crawler.py'%(cwd_str),   '%s/../bin/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_url_fetches__counts_by_day.py'%(cwd_str),         '%s/../bin/gf_apps/gf_analytics/py/stats'%(cwd_str))
                ]
            },
            #-------------
            #LIB
            #GF_CRAWL_LIB
            'gf_crawl_lib':{
                'go_path_str':'%s/../go/apps/gf_crawl_lib'%(cwd_str),
            },
            #-------------
            #MAIN
            #GF_PUBLISHER
            'gf_publisher':{
                'version_str':         '0.6.1.0',
                'go_path_str':         '%s/../go/apps/gf_publisher'%(cwd_str),
                'go_output_path_str':  '%s/../bin/gf_apps/gf_publisher/gf_publisher_service'%(cwd_str),
                'service_name_str':    'gf_publisher_service',
                'service_base_dir_str':'%s/../bin/gf_apps/gf_publisher'%(cwd_str),
            },
            
            #LIB
            #GF_PUBLISHER_LIB
            'gf_publisher_lib':{
                'go_path_str':'%s/../go/apps/gf_publisher_lib'%(cwd_str),

                #for tests serve data over http from this dir.
                #gf_publisher test runs an gf_images jobs_mngr to test post_creation, and jobs_mngr
                #needs to be able to fetch images over http that come from this dir.
                'test_data_to_serve_dir_str':'%s/../go/apps/gf_images_lib/tests_data'%(cwd_str),
            },
            #-------------
            #MAIN
            #GF_LANDING_PAGE
            'gf_landing_page':{
                'version_str':         '0.6.9.0',
                'go_path_str':         '%s/../go/gf_apps/gf_landing_page'%(cwd_str),
                'go_output_path_str':  '%s/../bin/gf_apps/gf_landing_page/gf_landing_page_service'%(cwd_str),
                'service_name_str':    'gf_landing_page_service',
                'service_base_dir_str':'%s/../bin/gf_apps/gf_landing_page'%(cwd_str),
            },
            #-------------
            #MAIN
            #GF_TAGGER
            'gf_tagger':{
                'version_str':         '0.6.1.0',
                'go_path_str':         '%s/../go/gf_apps/gf_tagger'%(cwd_str),
                'go_output_path_str':  '%s/../bin/gf_apps/gf_tagger/gf_tagger_service'%(cwd_str),
                'service_name_str':    'gf_tagger_service',
                'service_base_dir_str':'%s/../bin/gf_apps/gf_tagger'%(cwd_str),
            },
            #-------------
        }
    }
    return meta_map