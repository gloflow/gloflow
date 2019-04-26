import os
cwd_str = os.path.abspath(os.path.dirname(__file__))
#-------------------------------------------------------------
def get():

	apps_map = {
		#-----------------------------
		'gf_images':{
			'pages_map':{
				#-------------
				#IMAGES_FLOWS_BROWSER
				'gf_images_flows_browser':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_images'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_images/templates/gf_images_flows_browser/gf_images_flows_browser.html'%(cwd_str),
					'url_base_str':       '/images/static',
					#'type_str':           'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_images/js/gf_images_flows_browser.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_images/js/gf_images_flows_browser.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_images/ts/flows_browser/gf_images_flows_browser.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_gifs.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_gifs_viewer.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_image_viewer.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_sys_panel.ts'%(cwd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/masonry.pkgd.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/jquery.timeago.js'%(cwd_str),
					# 	]
					# 	#-------------
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_images/css/gf_images_flows_browser.css'%(cwd_str), '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_gifs_viewer.css'%(cwd_str),                    '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_image_viewer.css'%(cwd_str),                   '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_sys_panel.css'%(cwd_str),                      '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_images/templates/flows_browser/gf_images_flows_browser.html'%(cwd_str), '%s/../web/build/gf_apps/gf_images/templates/flows_browser'%(cwd_str)),
					# 	]
					# }
				},
				#-------------
				#IMAGES_DASHBOARD
				'gf_images_dashboard':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_images'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_images/templates/gf_images_dashboard/gf_images_dashboard.html'%(cwd_str),
					'url_base_str':       '/images/static',
					#'type_str':           'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_images/js/dashboard__ff0099__ooo.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_images/js/dashboard__ff0099__ooo.min.js'%(cwd_str),
					# 	'files_lst': [
					# 		'%s/../web/src/gf_apps/gf_images/ts/dashboard/gf_images_dashboard.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_images/ts/stats/gf_images_stats.ts'%(cwd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/d3.v3.js'%(cwd_str),
					# 		'%s/../web/libs/js/nv.d3_1.8.3.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_images/css/dashboard/gf_dashboard.css'%(cwd_str), '%s/../web/build/gf_apps/gf_images/css/dashboard'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_sys_panel.css'%(cwd_str),                     '%s/../web/build/gf_apps/gf_images/css/dashboard'%(cwd_str)),
					# 		('%s/../web/libs/css/nv.d3.css'%(cwd_str),                                   '%s/../web/build/gf_apps/gf_images/css/dashboard'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_images/templates/dashboard/gf_images_dashboard.html'%(cwd_str), '%s/../web/build/gf_apps/gf_images/templates/dashboard'%(cwd_str)),
					# 	]
					# }
				}
				#-------------
			}
		},
		#-----------------------------
		'gf_publisher':{
			'pages_map':{
				#-------------
				#GF_POST
				'gf_post':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_publisher'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_publisher/templates/gf_post/gf_post.html'%(cwd_str),
					'url_base_str':       '/posts/static',
					#'type_str':      'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_publisher/js/gf_post.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_publisher/js/gf_post.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_publisher/ts/gf_post/gf_post.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_publisher/ts/gf_post/gf_post_image_view.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_publisher/ts/gf_post/gf_post_tag_mini_view.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_client.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_input_ui.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_sys_panel.ts'%(cwd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_publisher/css/gf_post.css'%(cwd_str),         '%s/../web/build/gf_apps/gf_publisher/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_publisher/css/gf_post_tagging.css'%(cwd_str), '%s/../web/build/gf_apps/gf_publisher/css'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_sys_panel.css'%(cwd_str),                 '%s/../web/build/gf_apps/gf_publisher/css'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_publisher/templates/gf_post/gf_post.html'%(cwd_str), '%s/../web/build/gf_apps/gf_publisher/templates/gf_post'%(cwd_str)),
					# 	]
					# }
				},
				#-------------
				#GF_POSTS_BROWSER
				'gf_posts_browser':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_publisher'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_publisher/templates/gf_posts_browser/gf_posts_browser.html'%(cwd_str),
					'url_base_str':       '/posts/static',
					# 'type_str':      'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_publisher/js/gf_posts_browser.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_publisher/js/gf_posts_browser.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_publisher/ts/gf_posts_browser/gf_posts_browser.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_publisher/ts/gf_posts_browser/gf_posts_browser_view.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_publisher/ts/gf_posts_browser/gf_posts_browser_client.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_client.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_input_ui.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_tagger/ts/gf_tagger_client/gf_tagger_notes_ui.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_sys_panel.ts'%(cwd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/masonry.pkgd.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/jquery.timeago.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_publisher/gf_posts_browser.css'%(cwd_str),         '%s/../web/build/gf_apps/gf_publisher/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_publisher/gf_posts_browser_tagging.css'%(cwd_str), '%s/../web/build/gf_apps/gf_publisher/static/css'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_sys_panel.css'%(cwd_str),                      '%s/../web/build/gf_apps/gf_publisher/static/css'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_publisher/templates/gf_posts_browser/gf_posts_browser.html'%(cwd_str), '%s/../web/build/gf_apps/gf_publisher/templates/gf_posts_browser'%(cwd_str)),
					# 	]
					# }
				}
				#-------------
			}
		},
		#-----------------------------
		'gf_analytics':{
			'pages_map':{
				#-------------
				#DASHBOARD
				'gf_analytics_dashboard':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_analytics'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_analytics/templates/gf_analytics_dashboard/gf_analytics_dashboard.html'%(cwd_str),
					'url_base_str':       '/posts/static',
					# 'type_str':      'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_analytics/js/gf_analytics_dashboard.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_analytics/js/gf_analytics_dashboard.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_analytics/ts/dashboard/gf_analytics_dashboard.ts'%(cwd_str),
					# 		#-------------
					# 		#STATS__GF_IMAGES
					# 		'%s/../web/src/gf_apps/gf_images/ts/stats/gf_images_stats.ts'%(cwd_str),
					# 		#-------------
					# 		#STATS__GF_CRAWL
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__errors.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__fetches.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__images.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/stats/gf_crawl_stats__links.ts'%(cwd_str),
					# 		#-------------
					# 		'%s/../web/src/gf_stats/ts/gf_stats.ts'%(cwd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/d3.v3.js'%(cwd_str),
					# 		'%s/../web/libs/js/nv.d3_1.8.3.js'%(cwd_str),
					# 		'%s/../web/libs/js/jquery.timeago.js'%(cwd_str),
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js'%(cwd_str),
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js'%(cwd_str),
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/pixi.min.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_analytics/css/dashboard/gf_analytics_dashboard.css'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_stats/css/gf_stats.css'%(cwd_str),                                     '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#-------------
					# 		#STATS__GF_CRAWL
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl.css'%(cwd_str),         '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_errors.css'%(cwd_str),  '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_fetches.css'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_images.css'%(cwd_str),  '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/stats/stats__crawl_links.css'%(cwd_str),   '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#-------------
					# 		#STATS__GF_IMAGES
					# 		('%s/../web/src/gf_apps/gf_images/css/stats/stats__images.css'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#-------------
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_analytics/templates/dashboard/gf_analytics_dashboard.html'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/templates/gf_analytics_dashboard'%(cwd_str)),
					# 	]
					# }
				},
				#-------------
				#CRAWL_DASHBOARD
				'gf_crawl_dashboard':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_analytics'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_crawl_lib/templates/gf_crawl_dashboard/gf_crawl_dashboard.html'%(cwd_str),
					'url_base_str':       '/a/static',
					# 'type_str':      'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_analytics/js/gf_crawl_dashboard.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_analytics/js/gf_crawl_dashboard.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl__img_preview_tooltip.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl_dashboard.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl_events.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_crawl_lib/ts/dashboard/gf_crawl_images_browser.ts'%(cwd_str),							
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/d3.v3.js'%(cwd_str),
					# 		'%s/../web/libs/js/c3.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/sigma.1.2.0.layout.forceAtlas2.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/sigma.1.2.0.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/jquery.timeago.js'%(cwd_str),
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js'%(cwd_str),
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js'%(cwd_str),
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/pixi.min.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/browser.css'%(cwd_str),            '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/errors.css'%(cwd_str),             '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/gf_crawl_dashboard.css'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/fetches.css'%(cwd_str),            '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),							
					# 		#('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/gifs.css'%(cwd_str),               '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/images.css'%(cwd_str),             '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#('%s/../web/src/gf_apps/gf_crawl_lib/css/dashboard/links.css'%(cwd_str),              '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		#('%s/../web/src/gf_apps/gf_crawl_lib/css/lib/c3.min.css'%(cwd_str),                   '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_crawl_lib/templates/dashboard/gf_crawl_dashboard.html'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/templates/gf_crawl_dashboard'%(cwd_str)),
					# 	]
					# },
					# 'files_to_copy_lst':[
					# 	('%s/../web/src/gf_apps/gf_crawl_lib/assets/icons.png'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/assets'%(cwd_str),),
					# ]
				},
				#-------------
				#DOMAINS_BROWSER

				#IMPORTANT!! - this is in analytics, because domains are sources for images/posts, and so dont 
				#              belong to neither gf_images nor gf_publisher. maybe it should be its own app?
				'gf_domains_browser':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_analytics'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_domains_lib/templates/gf_domains_browser/gf_domains_browser.html'%(cwd_str),
					'url_base_str':       '/a/static',
					# 'type_str':      'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_analytics/js/gf_domains_browser.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_analytics/js/gf_domains_browser.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domain.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_browser.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_conn.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_infos.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_domains_lib/ts/domains_browser/gf_domains_search.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_color.ts'%(cwd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		#'%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/jquery.autocomplete.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/pixi.min.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_domains_lib/css/gf_domains_browser.css'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_sys_panel.css'%(cwd_str),                      '%s/../web/build/gf_apps/gf_analytics/css'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_domains_lib/templates/domains_browser/gf_domains_browser.html'%(cwd_str), '%s/../web/build/gf_apps/gf_analytics/templates/gf_domains_browser'%(cwd_str)),
					# 	]
					# }
				}
				#-------------
			}
		},
		#-----------------------------
		'gf_landing_page':{
			'pages_map':{
				'gf_landing_page':{
					'build_dir_str':      '%s/../web/build/gf_apps/gf_landing_page'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_apps/gf_landing_page/templates/gf_landing_page/gf_landing_page.html'%(cwd_str),
					'url_base_str':       '/landing/static',
					# 'type_str':      'ts',
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_landing_page/js/gf_landing_page.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_landing_page/js/gf_landing_page.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_landing_page/ts/gf_calc.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_landing_page/ts/gf_email_registration.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_landing_page/ts/gf_images.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_landing_page/ts/gf_landing_page.ts'%(cwd_str),
					# 		'%s/../web/src/gf_apps/gf_landing_page/ts/gf_procedural_art.ts'%(cwd_str)
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/jquery.timeago.js'%(cwd_str),
					# 		'%s/../web/libs/js/color-thief.min.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_landing_page/css/domains.css'%(cwd_str),         '%s/../web/build/gf_apps/gf_landing_page/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_landing_page/css/gf_landing_page.css'%(cwd_str), '%s/../web/build/gf_apps/gf_landing_page/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_landing_page/css/images.css'%(cwd_str),          '%s/../web/build/gf_apps/gf_landing_page/css'%(cwd_str)),
					# 		('%s/../web/src/gf_apps/gf_landing_page/css/posts.css'%(cwd_str),           '%s/../web/build/gf_apps/gf_landing_page/css'%(cwd_str))
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_landing_page/templates/gf_landing_page.html'%(cwd_str), '%s/../bin/gf_apps/gf_landing_page/templates'%(cwd_str)),
					# 	]
					# }
				}
			}
		},
		#-----------------------------
		# 'gf_user':{
		# 	'pages_map':{
		# 		'gf_user_profile':{
		# 			'type_str':      'ts',
		# 			'build_dir_str': '%s/../web/build/gf_apps/gf_user'%(cwd_str),
		# 			'ts':{
		# 				'out_file_str':      '%s/../web/build/gf_apps/gf_user/js/gf_user_profile.js'%(cwd_str),
		# 				'minified_file_str': '%s/../web/build/gf_apps/gf_user/js/gf_user_profile.min.js'%(cwd_str),
		# 				'files_lst':[
		# 					'%s/../web/src/gf_apps/gf_user/gf_user_profile.ts'%(cwd_str),
		# 					'%s/../web/src/gf_core/gf_sys_panel.ts'%(cwd_str),
		# 				],
		# 				#-------------
		# 				#LIBS
		# 				'libs_files_lst':[]
		# 				#-------------
		# 			},
		# 			'css':{
		# 				'files_lst':[
		# 					('%s/../web/src/gf_apps/gf_user/css/gf_user_profile.css'%(cwd_str), '%s/../web/build/gf_apps/gf_user/css'%(cwd_str))
		# 				]
		# 			},

		# 			#static files to copy without change
		# 			'files_to_copy_lst':[
		# 				('%s/../web/src/gf_apps/gf_user/gf_user_profile.html'%(cwd_str), '%s/../web/build/gf_apps/gf_user'%(cwd_str),)
		# 			]
		# 		}
		# 	}
		# },
		#-----------------------------
		# 'gf_tagger':{
		# 	'pages_map':{
		# 		'gf_tag_objects':{
		# 			'type_str':          'dart',
		# 			'code_root_dir_str': '%s/../web/src/gf_apps/gf_tagger/gf_tag_objects'%(cwd_str),
		# 			'target_deploy_dir': '%s/../web/build/gf_apps/gf_tagger/static'%(cwd_str),
		# 		},
		# 		##IMPORTANT!! - not a page itself, instead its code being used by other pages, but its included here
		# 		##              so that its code gets built when pages for this app are built
		# 		#'gf_tagger_client':{
		# 		#	'code_root_dir_str':'%s/../src/apps/gf_tagger/client/gf_tagger_client'%(cwd_str),
		# 		#	'target_deploy_dir':'%s/../bin/apps/gf_tagger/static'%(cwd_str),
		# 		#}
		# 	}
		# }
		#-----------------------------
	}
	return apps_map