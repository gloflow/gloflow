# GloFlow application and media management/publishing platform
# Copyright (C) 2020 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

#-------------------------------------------------------------
def get():
    

    # IMPORTANT!! - dependency graph between go/web packages and apps, used to know 
    #               which app containers to rebuild (CI/CD tools) in this monorepo.
    #               in "apps_gf_packages_map" the keys are names of applications, and values are lists of packages
    #               that are dependencies for that app. if those packages changed that app will be marked
    #               as changed and will be rebuilt.
    # 
    # FIX!! - have an automated way of determening this graph (no time for that right now).
    apps_changes_deps_map = {
        "apps_gf_packages_map": {

            "gf_images":[
                "gf_images",
                "gf_images_lib",
            ],

            "gf_analytics": [
                "gf_analytics",
                "gf_crawl_lib",
                "gf_domains_lib"
            ],
            "gf_publisher": [
                "gf_publisher",
                "gf_publisher_lib",
                "gf_images_lib"
            ],
            "gf_landing_page": [
                "gf_landing_page",
                "gf_images_lib",
                "gf_publisher_lib"
            ],
            "gf_tagger": [
                "gf_images_lib",
                "gf_publisher_lib",
                "gf_tagger"
            ],
        },

        "system_packages_lst": [
            "gf_core",
            "gf_rpc_lib",
            "gf_stats"
        ]
    }

    # AWS_S3
    aws_s3_map = {
        "images_s3_bucket_map": {
            "tests":         "gf--test--img",
            "local_cluster": "gf--local--cluster--img"
        }
    }


    meta_map = {
        "apps_changes_deps_map":             apps_changes_deps_map,
        "local_cluster_config_dir_path_str": "%s/../ops/tests/test_cluster"%(modd_str),
        "aws_s3_map":                        aws_s3_map,
        "build_info_map": {
            #------------------------
            # GF_SOLO
            "gf_solo": {
                "type_str":             "main_go",
                "version_str":          "latest",
                "go_path_str":          f"{modd_str}/../go/gf_apps/gf_solo",
                "go_output_path_str":   f"{modd_str}/../build/gf_apps/gf_solo/gf_solo",
                "service_name_str":     "gf_solo",
                "service_base_dir_str": f"{modd_str}/../build/gf_apps/gf_solo",
                "service_dockerfile_path_str": f"{modd_str}/../build/gf_apps/gf_solo/Dockerfile_ubuntu",
                "copy_to_dir_lst":    [
                    
                    #------------------------
                    # TENSORFLOW C_LIBS
                    (f"{modd_str}/../rust/build/tf_lib/lib/libtensorflow.so",           f"{modd_str}/../build/gf_apps/gf_solo/libs"),
                    (f"{modd_str}/../rust/build/tf_lib/lib/libtensorflow_framework.so", f"{modd_str}/../build/gf_apps/gf_solo/libs"),

                    #------------------------
                    # GF_IMAGES
                    
                    (f"{modd_str}/../rust/build/libgf_images_jobs.so", f"{modd_str}/../build/gf_apps/gf_solo/libs"),

                    # PY_PLUGINS                
                    (f"{modd_str}/../py/gf_apps/gf_images/plugins/gf_images_plugins_main.py", f"{modd_str}/../build/gf_apps/gf_solo/gf_images/plugins"),
                    # (f"{modd_str}/../py/gf_apps/gf_images/gf_images_palette/gf_color_palette.py", f"{modd_str}/../build/gf_apps/gf_solo/gf_images/plugins"),

                    #------------------------
                    # GF_ML_WORKER
                    ("%s/../py/gf_apps/gf_ml_worker/gf_ml_data.py"%(modd_str),      "%s/../build/gf_apps/gf_solo/gf_ml_worker/py"%(modd_str)),
                    ("%s/../py/gf_apps/gf_ml_worker/gf_plot.py"%(modd_str),         "%s/../build/gf_apps/gf_solo/gf_ml_worker/py"%(modd_str)),
                    ("%s/../py/gf_apps/gf_ml_worker/gf_simple_model.py"%(modd_str), "%s/../build/gf_apps/gf_solo/gf_ml_worker/py"%(modd_str)),
                    ("%s/../py/gf_apps/gf_ml_worker/requirements.txt"%(modd_str),   "%s/../build/gf_apps/gf_solo/gf_ml_worker/py"%(modd_str)),

                    # C_LIBS
                    # gf_images_jobs_py.so - gf_images_jobs Rust Python extension
                    # libtensorflow.so     - TensorFlow C lib
                    ("%s/../rust/build/gf_images_jobs_py.so"%(modd_str), "%s/../build/gf_apps/gf_solo/gf_ml_worker/py"%(modd_str)),
                    
                    #------------------------
                    # GF_ANALYTICS
                    
                    ("%s/../go/gf_stats/py/cli_stats.py"%(modd_str),                                                     "%s/../build/gf_apps/gf_solo/gf_analytics/py"%(modd_str)),
                    ("%s/../py/gf_stats/gf_errors__counts_by_day.py"%(modd_str),                                         "%s/../build/gf_apps/gf_solo/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_imgs__counts_by_day.py"%(modd_str),            "%s/../build/gf_apps/gf_solo/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__counts_by_day.py"%(modd_str),  "%s/../build/gf_apps/gf_solo/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__null_breakdown.py"%(modd_str), "%s/../build/gf_apps/gf_solo/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__per_crawler.py"%(modd_str),    "%s/../build/gf_apps/gf_solo/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_url_fetches__counts_by_day.py"%(modd_str),          "%s/../build/gf_apps/gf_solo/gf_analytics/py/stats"%(modd_str)),

                    #------------------------
                    # ASSETS

                    # icons png file used by gf_chrome_ext and gf_solo
                    # (f"{modd_str}/../gf_chrome_ext/assets/icons.png", f"{modd_str}/../build/gf_apps/gf_solo/assets"),

                    #------------------------
                ]
            },
            
            #------------------------
            # GF_P2P_TESTER
            "gf_p2p_tester": {
                "type_str":    "main_go",
                "version_str": "latest",
            },

            #------------------------
            # GF_ML_WORKER
            "gf_ml_worker": {
                "type_str":             "main_py",
                "version_str":          "latest",
                "service_name_str":     "gf_ml_worker",
                "service_base_dir_str": "%s/../build/gf_apps/gf_ml_worker"%(modd_str),
                "copy_to_dir_lst": [
                    ("%s/../py/gf_apps/gf_ml_worker/gf_ml_data.py"%(modd_str),      "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str)),
                    ("%s/../py/gf_apps/gf_ml_worker/gf_plot.py"%(modd_str),         "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str)),
                    ("%s/../py/gf_apps/gf_ml_worker/gf_simple_model.py"%(modd_str), "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str)),
                    ("%s/../py/gf_apps/gf_ml_worker/requirements.txt"%(modd_str),   "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str)),

                    # C_LIBS
                    # gf_images_jobs_py.so - gf_images_jobs Rust Python extension
                    # libtensorflow.so     - TensorFlow C lib
                    ("%s/../rust/build/gf_images_jobs_py.so"%(modd_str),       "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str)),
                    ("%s/../rust/build/libtensorflow.so"%(modd_str),           "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str)),
                    ("%s/../rust/build/libtensorflow_framework.so"%(modd_str), "%s/../build/gf_apps/gf_ml_worker/py"%(modd_str))
                ]
            },

            #-------------
            # GF_IMAGES_JOBS
            "gf_images_jobs": {
                "type_str":    "lib_rust",
                "version_str": "latest",
                "cargo_crate_specs_lst": [
                    {"dir_path_str": "%s/../rust/gf_images_jobs"%(modd_str), "static_bool": False}, # True},
                    {"dir_path_str": "%s/../rust/gf_images_jobs_py"%(modd_str)},
                ]
            },

            #-------------
            # MAIN
            # GF_IMAGES
            "gf_images": {
                "type_str":             "main_go",
                "version_str":          "latest", # "0.8.0.10",
                "go_path_str":          "%s/../go/gf_apps/gf_images"%(modd_str),
                "go_output_path_str":   "%s/../build/gf_apps/gf_images/gf_images_service"%(modd_str),
                "service_name_str":     "gf_images_service",
                "service_base_dir_str": "%s/../build/gf_apps/gf_images"%(modd_str),
            },
            
            # LIB
            # GF_IMAGES_LIB
            "gf_images_lib": {
                "type_str":                   "lib_go",
                "go_path_str":                "%s/../go/gf_apps/gf_images_lib"%(modd_str),
                "test_data_to_serve_dir_str": "%s/../go/gf_apps/gf_images_lib/tests_data"%(modd_str), #for tests serve data over http from this dir
            },

            #-------------
            # MAIN
            # GF_ANALYTICS
            "gf_analytics": {
                "type_str":             "main_go",
                "version_str":          "latest", # "0.8.0.7",
                "go_path_str":          "%s/../go/gf_apps/gf_analytics"%(modd_str),
                "go_output_path_str":   "%s/../build/gf_apps/gf_analytics/gf_analytics_service"%(modd_str),
                "service_name_str":     "gf_analytics_service",
                "service_base_dir_str": "%s/../build/gf_apps/gf_analytics"%(modd_str),
                "copy_to_dir_lst": [
                    ("%s/../go/gf_stats/py/cli_stats.py"%(modd_str),                                                     "%s/../build/gf_apps/gf_analytics/py"%(modd_str)),
                    ("%s/../py/gf_stats/gf_errors__counts_by_day.py"%(modd_str),                                         "%s/../build/gf_apps/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_imgs__counts_by_day.py"%(modd_str),            "%s/../build/gf_apps/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__counts_by_day.py"%(modd_str),  "%s/../build/gf_apps/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__null_breakdown.py"%(modd_str), "%s/../build/gf_apps/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_page_outgoing_links__per_crawler.py"%(modd_str),    "%s/../build/gf_apps/gf_analytics/py/stats"%(modd_str)),
                    ("%s/../go/gf_apps/gf_crawl_lib/py/stats/crawler_url_fetches__counts_by_day.py"%(modd_str),          "%s/../build/gf_apps/gf_analytics/py/stats"%(modd_str))
                ]
            },

            #-------------
            # LIB
            # GF_CRAWL_LIB
            "gf_crawl_lib": {
                "type_str":    "lib_go",
                "go_path_str": "%s/../go/gf_apps/gf_crawl_lib"%(modd_str),
            },
            "gf_crawl_core": {
                "type_str":    "lib_go",
                "go_path_str": "%s/../go/gf_apps/gf_crawl_lib/gf_crawl_core"%(modd_str),
            },

            #-------------
            # MAIN
            # GF_PUBLISHER
            "gf_publisher": {
                "type_str":             "main_go",
                "version_str":          "latest", # "0.8.0.4",
                "go_path_str":          "%s/../go/gf_apps/gf_publisher"%(modd_str),
                "go_output_path_str":   "%s/../build/gf_apps/gf_publisher/gf_publisher_service"%(modd_str),
                "service_name_str":     "gf_publisher_service",
                "service_base_dir_str": "%s/../build/gf_apps/gf_publisher"%(modd_str),
            },
            
            # LIB
            # GF_PUBLISHER_LIB
            "gf_publisher_lib": {
                "type_str":    "lib_go",
                "go_path_str": "%s/../go/gf_apps/gf_publisher_lib"%(modd_str),

                # for tests serve data over http from this dir.
                # gf_publisher test runs an gf_images jobs_mngr to test post_creation, and jobs_mngr
                # needs to be able to fetch images over http that come from this dir.
                "test_data_to_serve_dir_str":"%s/../go/gf_apps/gf_images_lib/tests_data"%(modd_str),
            },

            #-------------
            # MAIN
            # GF_LANDING_PAGE
            "gf_landing_page": {
                "type_str":             "main_go",
                "version_str":          "latest", # "0.8.0.11",
                "go_path_str":          "%s/../go/gf_apps/gf_landing_page"%(modd_str),
                "go_output_path_str":   "%s/../build/gf_apps/gf_landing_page/gf_landing_page_service"%(modd_str),
                "service_name_str":     "gf_landing_page_service",
                "service_base_dir_str": "%s/../build/gf_apps/gf_landing_page"%(modd_str),
            },

            #-------------
            # MAIN
            # GF_TAGGER
            "gf_tagger": {
                "type_str":             "main_go",
                "version_str":          "latest", # "0.8.0.1",
                "go_path_str":          "%s/../go/gf_apps/gf_tagger"%(modd_str),
                "go_output_path_str":   "%s/../build/gf_apps/gf_tagger/gf_tagger_service"%(modd_str),
                "service_name_str":     "gf_tagger_service",
                "service_base_dir_str": "%s/../build/gf_apps/gf_tagger"%(modd_str),
            },

            #-------------
            # GF_BUILDER_WEB
            "gf_builder_web": {
                "type_str":            "custom",
                "version_str":         "latest",
                "cont_image_name_str": "gf_builder_web",
                "image_tag_str":       "latest",
                "dockerfile_path_str": "%s/../Dockerfile__gf_builder_web"%(modd_str)
            },

            #-------------
            # GF_BUILDER_GO_UBUNTU
            "gf_builder_go_ubuntu": {
                "type_str":            "custom",
                "version_str":         "latest",
                "cont_image_name_str": "gf_builder_go_ubuntu",
                "image_tag_str":       "latest",
                "dockerfile_path_str": "%s/../Dockerfile__gf_builder_go__ubuntu"%(modd_str)
            },

            #-------------
            # GF_BUILDER_RUST_UBUNTU
            "gf_builder_rust_ubuntu": {
                "type_str":            "custom",
                "version_str":         "latest",
                "cont_image_name_str": "gf_builder_rust_ubuntu",
                "image_tag_str":       "latest",
                "dockerfile_path_str": "%s/../Dockerfile__gf_builder_rust__ubuntu"%(modd_str)
            },

            #-------------
        }
    }

    return meta_map