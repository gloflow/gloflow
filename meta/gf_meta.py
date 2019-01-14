

import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))



def get():

    meta_map = {
        'build_info_map':{
            'gf_images':{
                'go_path_str':       '%s/../go/apps/gf_images/gf_images_service.go'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_images_service'%(cwd_str),
            },

            'gf_publisher':{
                'go_path_str':       '%s/../go/apps/gf_publisher/gf_publisher_service.go'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_publisher_service'%(cwd_str),
            },

            'gf_tagger':{
                'go_path_str':       '%s/../go/apps/gf_tagger/gf_tagger_service.go'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_tagger_service'%(cwd_str),
            },

            'gf_landing_page':{
                'go_path_str':       '%s/../go/apps/gf_landing_page/gf_landing_page_service.go'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_landing_page_service'%(cwd_str),
            },

            'gf_analytics':{
                'go_path_str':       '%s/../go/apps/gf_analytics/gf_analytics_service.go'%(cwd_str),
                'go_output_path_str':'%s/../bin/gf_analytics_service'%(cwd_str),
            },
        }
    }
    return meta_map