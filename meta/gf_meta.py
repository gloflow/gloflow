

import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))
#-------------------------------------------------------------
def get():
    
    #IMPORTANT!! - dependency graph between go/web packages and apps, used to know 
    #              which app containers to rebuild (CI/CD tools) in this monorepo.
    #FIX!!       - have an automated way of determening this graph (no time for that right now).
    apps_changes_deps_map = {
        'apps_names_map': {
            'gf_images':       ['gf_images',     'gf_images_lib'],
            'gf_analytics':    ['gf_analytics',  'gf_crawl_lib',     'gf_domains_lib'],
            'gf_publisher':    ['gf_images_lib', 'gf_publisher',     'gf_publisher_lib'],
            'gf_tagger':       ['gf_images_lib', 'gf_publisher_lib', 'gf_tagger'],
            'gf_landing_page': ['gf_images_lib', 'gf_publisher_lib', 'gf_landing_page'],
        },

        'system_packages_lst': [
            'gf_core',
            'gf_rpc_lib',
            'gf_stats'
        ]
    }

    meta_map = {
        'apps_changes_deps_map': apps_changes_deps_map,
        'build_info_map':{
            #------------------------
            #GF_SOLO
            'gf_solo':{
                #'type_str':           'main',
                'version_str':        '0.8.0.0',
                'go_output_path_str': '%s/../build/gf_apps/apps/gf_solo/gf_solo'%(cwd_str),
                'copy_to_dir_lst': [
                    ('%s/../go/src/apps/gf_solo/gf_php_lib/php/index.php'%(cwd_str), '%s/../go/build/apps/gf_solo/php'%(cwd_str)),
                ]
            },
            #-------------
            #MAIN
            #GF_IMAGES
            'gf_images':{
                'version_str':          '0.8.0.8',
                'go_path_str':          '%s/../go/gf_apps/gf_images'%(cwd_str),
                'go_output_path_str':   '%s/../build/gf_apps/gf_images/gf_images_service'%(cwd_str),
                'service_name_str':     'gf_images_service',
                'service_base_dir_str': '%s/../build/gf_apps/gf_images'%(cwd_str),
            },
            
            #LIB
            #GF_IMAGES_LIB
            'gf_images_lib':{
                'go_path_str':                '%s/../go/gf_apps/gf_images_lib'%(cwd_str),
                'test_data_to_serve_dir_str': '%s/../go/gf_apps/gf_images_lib/tests_data'%(cwd_str), #for tests serve data over http from this dir
            },
            #-------------
            #MAIN
            #GF_ANALYTICS
            'gf_analytics':{
                'version_str':          '0.8.0.1',
                'go_path_str':          '%s/../go/gf_apps/gf_analytics'%(cwd_str),
                'go_output_path_str':   '%s/../build/gf_apps/gf_analytics/gf_analytics_service'%(cwd_str),
                'service_name_str':     'gf_analytics_service',
                'service_base_dir_str': '%s/../build/gf_apps/gf_analytics'%(cwd_str),
                'copy_to_dir_lst': [
                    ('%s/../go/gf_stats/py/cli_stats.py'%(cwd_str),                                                    '%s/../build/gf_apps/gf_analytics/py'%(cwd_str)),
                    ('%s/../go/gf_core/py/stats/gf_errors__counts_by_day.py'%(cwd_str),                                '%s/../build/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_imgs__counts_by_day.py'%(cwd_str),           '%s/../build/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__counts_by_day.py'%(cwd_str), '%s/../build/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__null_breakdown.py'%(cwd_str),'%s/../build/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__per_crawler.py'%(cwd_str),   '%s/../build/gf_apps/gf_analytics/py/stats'%(cwd_str)),
                    ('%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_url_fetches__counts_by_day.py'%(cwd_str),         '%s/../build/gf_apps/gf_analytics/py/stats'%(cwd_str))
                ]
            },
            #-------------
            #LIB
            #GF_CRAWL_LIB
            'gf_crawl_lib':{
                'go_path_str': '%s/../go/gf_apps/gf_crawl_lib'%(cwd_str),
            },
            #-------------
            #MAIN
            #GF_PUBLISHER
            'gf_publisher':{
                'version_str':          '0.8.0.1',
                'go_path_str':          '%s/../go/gf_apps/gf_publisher'%(cwd_str),
                'go_output_path_str':   '%s/../build/gf_apps/gf_publisher/gf_publisher_service'%(cwd_str),
                'service_name_str':     'gf_publisher_service',
                'service_base_dir_str': '%s/../build/gf_apps/gf_publisher'%(cwd_str),
            },
            
            #LIB
            #GF_PUBLISHER_LIB
            'gf_publisher_lib':{
                'go_path_str':'%s/../go/gf_apps/gf_publisher_lib'%(cwd_str),

                #for tests serve data over http from this dir.
                #gf_publisher test runs an gf_images jobs_mngr to test post_creation, and jobs_mngr
                #needs to be able to fetch images over http that come from this dir.
                'test_data_to_serve_dir_str':'%s/../go/gf_apps/gf_images_lib/tests_data'%(cwd_str),
            },
            #-------------
            #MAIN
            #GF_LANDING_PAGE
            'gf_landing_page':{
                'version_str':          '0.8.0.11',
                'go_path_str':          '%s/../go/gf_apps/gf_landing_page'%(cwd_str),
                'go_output_path_str':   '%s/../build/gf_apps/gf_landing_page/gf_landing_page_service'%(cwd_str),
                'service_name_str':     'gf_landing_page_service',
                'service_base_dir_str': '%s/../build/gf_apps/gf_landing_page'%(cwd_str),
            },
            #-------------
            #MAIN
            #GF_TAGGER
            'gf_tagger':{
                'version_str':          '0.8.0.1',
                'go_path_str':          '%s/../go/gf_apps/gf_tagger'%(cwd_str),
                'go_output_path_str':   '%s/../build/gf_apps/gf_tagger/gf_tagger_service'%(cwd_str),
                'service_name_str':     'gf_tagger_service',
                'service_base_dir_str': '%s/../build/gf_apps/gf_tagger'%(cwd_str),
            },
            #-------------
        }
    }
    return meta_map