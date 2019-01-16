

import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))



def get():

    meta_map = {
        'build_info_map':{
            #-------------
            #MAIN
            'gf_images':{
                'go_path_str':       '%s/../go/apps/gf_images'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_images_service'%(cwd_str),
            },
            #-------------
            #MAIN
            'gf_publisher':{
                'go_path_str':       '%s/../go/apps/gf_publisher'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_publisher_service'%(cwd_str),
            },
            #-------------
            #LIB
            'gf_publisher_lib':{
                'go_path_str':'%s/../go/apps/gf_publisher_lib'%(cwd_str)
            },
            #-------------
            #MAIN
            'gf_tagger':{
                'go_path_str':       '%s/../go/apps/gf_tagger'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_tagger_service'%(cwd_str),
            },
            #-------------
            #MAIN
            'gf_landing_page':{
                'go_path_str':       '%s/../go/apps/gf_landing_page'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_landing_page_service'%(cwd_str),
            },
            #-------------
            #MAIN
            'gf_analytics':{
                'go_path_str':       '%s/../go/apps/gf_analytics'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_analytics_service'%(cwd_str),
            },
            #-------------
            #LIB
            'gf_crawl_lib':{
                'go_path_str':'%s/../go/apps/gf_crawl_lib'%(cwd_str),
            },
            #-------------
        }
    }
    return meta_map