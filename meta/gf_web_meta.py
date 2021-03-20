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

	apps_map = {
		#-----------------------------
		"gf_landing_page":{
			"pages_map":{
				"gf_landing_page":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_landing_page"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_landing_page/templates/gf_landing_page/gf_landing_page.html"%(modd_str),
					"url_base_str":       "/landing/static",
				}
			}
		},

		#-----------------------------
		"gf_images":{
			"pages_map":{
				#-------------
				#IMAGES_FLOWS_BROWSER
				"gf_images_flows_browser":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_images"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_images/templates/gf_images_flows_browser/gf_images_flows_browser.html"%(modd_str),
					"url_base_str":       "/images/static",
					#"type_str":           "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_images/js/gf_images_flows_browser.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_images/js/gf_images_flows_browser.min.js"%(modd_str),
					# 	"files_lst":[
					# 		"%s/../web/src/gf_apps/gf_images/ts/flows_browser/gf_images_flows_browser.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_gifs.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_gifs_viewer.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_image_viewer.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_sys_panel.ts"%(modd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[
					# 		"%s/../web/libs/js/masonry.pkgd.min.js"%(modd_str),
					# 		"%s/../web/libs/js/jquery.timeago.js"%(modd_str),
					# 	]
					# 	#-------------
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_images/css/gf_images_flows_browser.css"%(modd_str), "%s/../web/build/gf_apps/gf_images/css/flows_browser"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_gifs_viewer.css"%(modd_str),                    "%s/../web/build/gf_apps/gf_images/css/flows_browser"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_image_viewer.css"%(modd_str),                   "%s/../web/build/gf_apps/gf_images/css/flows_browser"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_sys_panel.css"%(modd_str),                      "%s/../web/build/gf_apps/gf_images/css/flows_browser"%(modd_str)),
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_images/templates/flows_browser/gf_images_flows_browser.html"%(modd_str), "%s/../web/build/gf_apps/gf_images/templates/flows_browser"%(modd_str)),
					# 	]
					# }
				},

				#-------------
				# IMAGES_DASHBOARD
				"gf_images_dashboard":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_images"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_images/templates/gf_images_dashboard/gf_images_dashboard.html"%(modd_str),
					"url_base_str":       "/images/static",
					#"type_str":           "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_images/js/dashboard__ff0099__ooo.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_images/js/dashboard__ff0099__ooo.min.js"%(modd_str),
					# 	"files_lst": [
					# 		"%s/../web/src/gf_apps/gf_images/ts/dashboard/gf_images_dashboard.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_images/ts/stats/gf_images_stats.ts"%(modd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[
					# 		"%s/../web/libs/js/d3.v3.js"%(modd_str),
					# 		"%s/../web/libs/js/nv.d3_1.8.3.js"%(modd_str),
					# 	]
					# 	#-------------
					# },
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_images/css/dashboard/gf_dashboard.css"%(modd_str), "%s/../web/build/gf_apps/gf_images/css/dashboard"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_sys_panel.css"%(modd_str),                     "%s/../web/build/gf_apps/gf_images/css/dashboard"%(modd_str)),
					# 		("%s/../web/libs/css/nv.d3.css"%(modd_str),                                   "%s/../web/build/gf_apps/gf_images/css/dashboard"%(modd_str)),
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_images/templates/dashboard/gf_images_dashboard.html"%(modd_str), "%s/../web/build/gf_apps/gf_images/templates/dashboard"%(modd_str)),
					# 	]
					# }
				}
				#-------------
			}
		},

		#-----------------------------
		"gf_publisher":{
			"pages_map":{
				#-------------
				# GF_POST
				"gf_post":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_publisher"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_publisher/templates/gf_post/gf_post.html"%(modd_str),
					"url_base_str":       "/posts/static",
					#"type_str":      "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_publisher/js/gf_post.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_publisher/js/gf_post.min.js"%(modd_str),
					# 	"files_lst":[
					# 		"%s/../web/src/gf_apps/gf_publisher/ts/gf_post/gf_post.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_publisher/ts/gf_post/gf_post_image_view.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_publisher/ts/gf_post/gf_post_tag_mini_view.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_client.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_input_ui.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_sys_panel.ts"%(modd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[]
					# 	#-------------
					# },
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_publisher/css/gf_post.css"%(modd_str),         "%s/../web/build/gf_apps/gf_publisher/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_publisher/css/gf_post_tagging.css"%(modd_str), "%s/../web/build/gf_apps/gf_publisher/css"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_sys_panel.css"%(modd_str),                 "%s/../web/build/gf_apps/gf_publisher/css"%(modd_str)),
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_publisher/templates/gf_post/gf_post.html"%(modd_str), "%s/../web/build/gf_apps/gf_publisher/templates/gf_post"%(modd_str)),
					# 	]
					# }
				},

				#-------------
				# GF_POSTS_BROWSER
				"gf_posts_browser":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_publisher"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_publisher/templates/gf_posts_browser/gf_posts_browser.html"%(modd_str),
					"url_base_str":       "/posts/static",
					# "type_str":      "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_publisher/js/gf_posts_browser.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_publisher/js/gf_posts_browser.min.js"%(modd_str),
					# 	"files_lst":[
					# 		"%s/../web/src/gf_apps/gf_publisher/ts/gf_posts_browser/gf_posts_browser.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_publisher/ts/gf_posts_browser/gf_posts_browser_view.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_publisher/ts/gf_posts_browser/gf_posts_browser_client.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_client.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_input_ui.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_notes_ui.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_sys_panel.ts"%(modd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[
					# 		"%s/../web/libs/js/masonry.pkgd.min.js"%(modd_str),
					# 		"%s/../web/libs/js/jquery.timeago.js"%(modd_str),
					# 	]
					# 	#-------------
					# },
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_publisher/gf_posts_browser.css"%(modd_str),         "%s/../web/build/gf_apps/gf_publisher/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_publisher/gf_posts_browser_tagging.css"%(modd_str), "%s/../web/build/gf_apps/gf_publisher/static/css"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_sys_panel.css"%(modd_str),                      "%s/../web/build/gf_apps/gf_publisher/static/css"%(modd_str)),
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_publisher/templates/gf_posts_browser/gf_posts_browser.html"%(modd_str), "%s/../web/build/gf_apps/gf_publisher/templates/gf_posts_browser"%(modd_str)),
					# 	]
					# }
				}

				#-------------
			}
		},

		#-----------------------------
		"gf_analytics":{
			"pages_map":{
				#-------------
				# DASHBOARD
				"gf_analytics_dashboard":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_analytics"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_analytics/templates/gf_analytics_dashboard/gf_analytics_dashboard.html"%(modd_str),
					"url_base_str":       "/posts/static",
					# "type_str":      "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_analytics/js/gf_analytics_dashboard.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_analytics/js/gf_analytics_dashboard.min.js"%(modd_str),
					# 	"files_lst":[
					# 		"%s/../web/src/gf_apps/gf_analytics/ts/dashboard/gf_analytics_dashboard.ts"%(modd_str),
					# 		#-------------
					# 		#STATS__GF_IMAGES
					# 		"%s/../web/src/gf_apps/gf_images/ts/stats/gf_images_stats.ts"%(modd_str),
					# 		#-------------
					# 		#STATS__GF_CRAWL
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__errors.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__fetches.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__images.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__links.ts"%(modd_str),
					# 		#-------------
					# 		"%s/../web/src/gf_stats/ts/gf_stats.ts"%(modd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[
					# 		"%s/../web/libs/js/d3.v3.js"%(modd_str),
					# 		"%s/../web/libs/js/nv.d3_1.8.3.js"%(modd_str),
					# 		"%s/../web/libs/js/jquery.timeago.js"%(modd_str),
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js"%(modd_str),
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js"%(modd_str),
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/pixi.min.js"%(modd_str),
					# 	]
					# 	#-------------
					# },
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_analytics/css/dashboard/gf_analytics_dashboard.css"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_stats/css/gf_stats.css"%(modd_str),                                     "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#-------------
					# 		#STATS__GF_CRAWL
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl.css"%(modd_str),         "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_errors.css"%(modd_str),  "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_fetches.css"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_images.css"%(modd_str),  "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_links.css"%(modd_str),   "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#-------------
					# 		#STATS__GF_IMAGES
					# 		("%s/../web/src/gf_apps/gf_images/css/stats/stats__images.css"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#-------------
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_analytics/templates/dashboard/gf_analytics_dashboard.html"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/templates/gf_analytics_dashboard"%(modd_str)),
					# 	]
					# }
				},

				#-------------
				# CRAWL_DASHBOARD
				"gf_crawl_dashboard":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_analytics"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_crawl_lib/templates/gf_crawl_dashboard/gf_crawl_dashboard.html"%(modd_str),
					"url_base_str":       "/a/static",
					# "type_str":      "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_analytics/js/gf_crawl_dashboard.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_analytics/js/gf_crawl_dashboard.min.js"%(modd_str),
					# 	"files_lst":[
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl__img_preview_tooltip.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl_dashboard.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl_events.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl_images_browser.ts"%(modd_str),							
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[
					# 		"%s/../web/libs/js/d3.v3.js"%(modd_str),
					# 		"%s/../web/libs/js/c3.min.js"%(modd_str),
					# 		"%s/../web/libs/js/sigma.1.2.0.layout.forceAtlas2.min.js"%(modd_str),
					# 		"%s/../web/libs/js/sigma.1.2.0.min.js"%(modd_str),
					# 		"%s/../web/libs/js/jquery.timeago.js"%(modd_str),
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js"%(modd_str),
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js"%(modd_str),
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/pixi.min.js"%(modd_str),
					# 	]
					# 	#-------------
					# },
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/browser.css"%(modd_str),            "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/errors.css"%(modd_str),             "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/gf_crawl_dashboard.css"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/fetches.css"%(modd_str),            "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),							
					# 		#("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/gifs.css"%(modd_str),               "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/images.css"%(modd_str),             "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#("%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/links.css"%(modd_str),              "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		#("%s/../web/src/gf_apps/gf_crawl_lib/css/lib/c3.min.css"%(modd_str),                   "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_crawl_lib/templates/dashboard/gf_crawl_dashboard.html"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/templates/gf_crawl_dashboard"%(modd_str)),
					# 	]
					# },
					# "files_to_copy_lst":[
					# 	("%s/../web/src/gf_apps/gf_crawl_lib/assets/icons.png"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/assets"%(modd_str),),
					# ]
				},

				#-------------
				# DOMAINS_BROWSER

				# IMPORTANT!! - this is in analytics, because domains are sources for images/posts, and so dont 
				#               belong to neither gf_images nor gf_publisher. maybe it should be its own app?
				"gf_domains_browser":{
					"build_dir_str":      "%s/../web/build/gf_apps/gf_analytics"%(modd_str),
					"main_html_path_str": "%s/../web/src/gf_apps/gf_domains_lib/templates/gf_domains_browser/gf_domains_browser.html"%(modd_str),
					"url_base_str":       "/a/static",
					# "type_str":      "ts",
					# "ts":{
					# 	"out_file_str":      "%s/../web/build/gf_apps/gf_analytics/js/gf_domains_browser.js"%(modd_str),
					# 	"minified_file_str": "%s/../web/build/gf_apps/gf_analytics/js/gf_domains_browser.min.js"%(modd_str),
					# 	"files_lst":[
					# 		"%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domain.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_browser.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_conn.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_infos.ts"%(modd_str),
					# 		"%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_search.ts"%(modd_str),
					# 		"%s/../web/src/gf_core/ts/gf_color.ts"%(modd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	"libs_files_lst":[
					# 		#"%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js"%(modd_str),
					# 		"%s/../web/libs/js/jquery.autocomplete.min.js"%(modd_str),
					# 		"%s/../web/libs/js/pixi.min.js"%(modd_str),
					# 	]
					# 	#-------------
					# },
					# "css":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_domains_lib/css/gf_domains_browser.css"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 		("%s/../web/src/gf_core/css/gf_sys_panel.css"%(modd_str),                      "%s/../web/build/gf_apps/gf_analytics/css"%(modd_str)),
					# 	]
					# },
					# "templates":{
					# 	"files_lst":[
					# 		("%s/../web/src/gf_apps/gf_domains_lib/templates/domains_browser/gf_domains_browser.html"%(modd_str), "%s/../web/build/gf_apps/gf_analytics/templates/gf_domains_browser"%(modd_str)),
					# 	]
					# }
				}

				#-------------
			}
		},

		#-----------------------------
		# "gf_user":{
		# 	"pages_map":{
		# 		"gf_user_profile":{
		# 			"type_str":      "ts",
		# 			"build_dir_str": "%s/../web/build/gf_apps/gf_user"%(modd_str),
		# 			"ts":{
		# 				"out_file_str":      "%s/../web/build/gf_apps/gf_user/js/gf_user_profile.js"%(modd_str),
		# 				"minified_file_str": "%s/../web/build/gf_apps/gf_user/js/gf_user_profile.min.js"%(modd_str),
		# 				"files_lst":[
		# 					"%s/../web/src/gf_apps/gf_user/gf_user_profile.ts"%(modd_str),
		# 					"%s/../web/src/gf_core/gf_sys_panel.ts"%(modd_str),
		# 				],
		# 				#-------------
		# 				#LIBS
		# 				"libs_files_lst":[]
		# 				#-------------
		# 			},
		# 			"css":{
		# 				"files_lst":[
		# 					("%s/../web/src/gf_apps/gf_user/css/gf_user_profile.css"%(modd_str), "%s/../web/build/gf_apps/gf_user/css"%(modd_str))
		# 				]
		# 			},
		#
		# 			#static files to copy without change
		# 			"files_to_copy_lst":[
		# 				("%s/../web/src/gf_apps/gf_user/gf_user_profile.html"%(modd_str), "%s/../web/build/gf_apps/gf_user"%(modd_str),)
		# 			]
		# 		}
		# 	}
		# },
		#-----------------------------
		# "gf_tagger":{
		# 	"pages_map":{
		# 		"gf_tag_objects":{
		# 			"type_str":          "dart",
		# 			"code_root_dir_str": "%s/../web/src/gf_apps/gf_tagger/gf_tag_objects"%(modd_str),
		# 			"target_deploy_dir": "%s/../web/build/gf_apps/gf_tagger/static"%(modd_str),
		# 		},
		# 		##IMPORTANT!! - not a page itself, instead its code being used by other pages, but its included here
		# 		##              so that its code gets built when pages for this app are built
		# 		#"gf_tagger_client":{
		# 		#	"code_root_dir_str":"%s/../src/apps/gf_tagger/client/gf_tagger_client"%(modd_str),
		# 		#	"target_deploy_dir":"%s/../bin/apps/gf_tagger/static"%(modd_str),
		# 		#}
		# 	}
		# }
		#-----------------------------
	}
	return apps_map