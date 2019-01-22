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
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../web/gf_apps/gf_images/client/src/flows_browser'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_images/static'%(cwd_str),
					'ts':{
						'out_file_str'     :'%s/../web/build/apps/gf_images/static/js/gf_images_flows_browser.js'%(cwd_str),
						'minified_file_str':'%s/../web/build/apps/gf_images/static/js/gf_images_flows_browser.min.js'%(cwd_str),
						'files_lst':[
							'%s/../web/gf_apps/gf_images/client/src/flows_browser/gf_images_flows_browser.ts'%(cwd_str),
							'%s/../web/gf_core/client/src/gf_gifs.ts'%(cwd_str),
							'%s/../web/gf_core/client/src/gf_gifs_viewer.ts'%(cwd_str),
							'%s/../web/gf_core/client/src/gf_image_viewer.ts'%(cwd_str),
							'%s/../web/gf_core/client/src/gf_sys_panel.ts'%(cwd_str),
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							'%s/../web/apps/gf_images/client/src/lib/masonry.pkgd.min.js'%(cwd_str),
							'%s/../web/apps/gf_images/client/src/lib/jquery.timeago.js'%(cwd_str),
						]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../web/apps/gf_images/client/css/gf_images_flows_browser.css'%(cwd_str),'%s/../go/bin/apps/gf_images/static/css'%(cwd_str)),
							('%s/../web/gf_core/client/css/gf_gifs_viewer.css'%(cwd_str)                ,'%s/../go/bin/apps/gf_images/static/css'%(cwd_str)),
							('%s/../web/gf_core/client/css/gf_image_viewer.css'%(cwd_str)               ,'%s/../go/bin/apps/gf_images/static/css'%(cwd_str)),
							('%s/../web/gf_core/client/css/gf_sys_panel.css'%(cwd_str)                  ,'%s/../go/bin/apps/gf_images/static/css'%(cwd_str)),
						]
					},
				},
				#-------------
				#IMAGES_DASHBOARD
				'gf_images_dashboard':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../web/apps/gf_images/client/src/dashboard'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_images/static'%(cwd_str),
					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_images/static/js/dashboard__ff0099__ooo.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_images/static/js/dashboard__ff0099__ooo.min.js'%(cwd_str),
						'files_lst':[
							'%s/../web/apps/gf_images/client/src/dashboard/gf_images_dashboard.ts'%(cwd_str),
							'%s/../web/apps/gf_images/client/src/stats/gf_images_stats.ts'%(cwd_str),
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							'%s/../web/apps/gf_images/client/src/lib/d3.v3.js'%(cwd_str),
							'%s/../web/apps/gf_images/client/src/lib/nv.d3_1.8.3.js'%(cwd_str),
						]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../web/apps/gf_images/client/css/lib/nv.d3.css'%(cwd_str)          ,'%s/../go/bin/apps/gf_images/static/css/lib'%(cwd_str)),
							('%s/../web/apps/gf_images/client/css/dashboard/dashboard.css'%(cwd_str),'%s/../go/bin/apps/gf_images/static/css/dashboard'%(cwd_str))
						]
					},
					#static files to copy without change
					'files_to_copy_lst':[
						('%s/../web/apps/gf_images/client/dashboard__ff0099__ooo.html'%(cwd_str),'%s/../go/bin/apps/gf_images/static'%(cwd_str),)
					]
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
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../web/apps/gf_analytics/client'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_analytics/static'%(cwd_str),

					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_analytics/static/js/gf_analytics_dashboard.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_analytics/static/js/gf_analytics_dashboard.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_analytics/client/src/dashboard/gf_analytics_dashboard.ts'%(cwd_str),
							#-------------
							#STATS__GF_IMAGES
							'%s/../go/src/apps/gf_images/client/src/stats/gf_images_stats.ts'%(cwd_str),
							#-------------
							#STATS__GF_CRAWL
							'%s/../go/src/apps/gf_crawl_lib/client/src/stats/gf_crawl_stats.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/stats/gf_crawl_stats__errors.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/stats/gf_crawl_stats__fetches.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/stats/gf_crawl_stats__images.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/stats/gf_crawl_stats__links.ts'%(cwd_str),
							#-------------
							'%s/../src/gf_stats/client/src/gf_stats.ts'%(cwd_str)
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							'%s/../go/src/apps/gf_analytics/client/lib/d3.v3.js'%(cwd_str),
							'%s/../go/src/apps/gf_analytics/client/lib/nv.d3_1.8.3.js'%(cwd_str),
							'%s/../go/src/apps/gf_analytics/client/lib/jquery.timeago.js'%(cwd_str),
							#'%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js'%(cwd_str),
							#'%s/../src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js'%(cwd_str),
							#'%s/../src/apps/gf_domains_lib/client/lib/pixi.min.js'%(cwd_str),
						]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_analytics/client/css/dashboard/gf_analytics_dashboard.css'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/gf_stats/client/css/gf_stats.css'%(cwd_str)                                 ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							#-------------
							#STATS__GF_CRAWL
							('%s/../go/src/apps/gf_crawl_lib/client/css/stats/stats__crawl.css'%(cwd_str),        '%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/stats/stats__crawl_errors.css'%(cwd_str), '%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/stats/stats__crawl_fetches.css'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/stats/stats__crawl_images.css'%(cwd_str), '%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/stats/stats__crawl_links.css'%(cwd_str),  '%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							#-------------
							#STATS__GF_IMAGES
							('%s/../go/src/apps/gf_images/client/css/stats/stats__images.css'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							#-------------
						]
					},

					#static files to copy without change
					'files_to_copy_lst':[
						('%s/../go/src/apps/gf_analytics/client/src/dashboard/analytics_dashboard__ff0099__ooo.html'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static'%(cwd_str),),
					]
				},
				#-------------
				#CRAWL_DASHBOARD
				'gf_crawl_dashboard':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../go/src/apps/gf_crawl_lib/client'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_analytics/static'%(cwd_str),

					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_analytics/static/js/gf_crawl_dashboard.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_analytics/static/js/gf_crawl_dashboard.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_crawl_lib/client/src/dashboard/gf_crawl__img_preview_tooltip.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/dashboard/gf_crawl_dashboard.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/dashboard/gf_crawl_events.ts'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/src/dashboard/gf_crawl_images_browser.ts'%(cwd_str),

							
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							'%s/../go/src/apps/gf_crawl_lib/client/lib/d3.v3.js'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/lib/c3.min.js'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/lib/sigma.1.2.0.layout.forceAtlas2.min.js'%(cwd_str),
							'%s/../go/src/apps/gf_crawl_lib/client/lib/sigma.1.2.0.min.js'%(cwd_str),
							'%s/../go/src/apps/gf_analytics/client/lib/jquery.timeago.js'%(cwd_str),
							#'%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js'%(cwd_str),
							#'%s/../src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js'%(cwd_str),
							#'%s/../src/apps/gf_domains_lib/client/lib/pixi.min.js'%(cwd_str),
						]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/browser.css'%(cwd_str)           ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/errors.css'%(cwd_str)            ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/fetches.css'%(cwd_str)           ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/gf_crawl_dashboard.css'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/gifs.css'%(cwd_str)              ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/images.css'%(cwd_str)            ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/dashboard/links.css'%(cwd_str)             ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_crawl_lib/client/css/lib/c3.min.css'%(cwd_str)                  ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							
						]
					},
					'files_to_copy_lst':[
						('%s/../go/src/apps/gf_crawl_lib/client/src/dashboard/crawl_dashboard_ff2___1112_29.html'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static'%(cwd_str),),
						('%s/../go/src/apps/gf_crawl_lib/client/assets/icons.png'%(cwd_str)                                ,'%s/../go/bin/apps/gf_analytics/static'%(cwd_str),),
					]
				},
				#-------------
				#DOMAINS_BROWSER

				#IMPORTANT!! - this is in analytics, because domains are sources for images/posts, and so dont 
				#              belong to neither gf_images nor gf_publisher. maybe it should be its own app?
				'gf_domains_browser':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../go/src/apps/gf_domains_lib/client'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_analytics/static'%(cwd_str),

					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_analytics/static/js/gf_domains_browser.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_analytics/static/js/gf_domains_browser.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_domains_lib/client/src/domains_browser/gf_domain.ts'%(cwd_str),
							'%s/../go/src/apps/gf_domains_lib/client/src/domains_browser/gf_domains_browser.ts'%(cwd_str),
							'%s/../go/src/apps/gf_domains_lib/client/src/domains_browser/gf_domains_conn.ts'%(cwd_str),
							'%s/../go/src/apps/gf_domains_lib/client/src/domains_browser/gf_domains_infos.ts'%(cwd_str),
							'%s/../go/src/apps/gf_domains_lib/client/src/domains_browser/gf_domains_search.ts'%(cwd_str),
							'%s/../go/src/gf_core/client/src/gf_color.ts'%(cwd_str),
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							#'%s/../src/apps/gf_domains_lib/client/lib/jquery-3.1.0.min.js'%(cwd_str),
							'%s/../go/src/apps/gf_domains_lib/client/lib/jquery.autocomplete.min.js'%(cwd_str),
							'%s/../go/src/apps/gf_domains_lib/client/lib/pixi.min.js'%(cwd_str),
						]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_domains_lib/client/css/gf_domains_browser.css'%(cwd_str),'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
							('%s/../go/src/gf_core/client/css/gf_sys_panel.css'%(cwd_str)                  ,'%s/../go/bin/apps/gf_analytics/static/css'%(cwd_str)),
						]
					},
				}
				#-------------
			}
		},
		#-----------------------------
		'gf_user':{
			'pages_map':{
				'gf_user_profile':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../go/src/apps/gf_user/client'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_user/static'%(cwd_str),
					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_user/static/js/gf_user_profile.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_user/static/js/gf_user_profile.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_user/client/src/gf_user_profile.ts'%(cwd_str),
							'%s/../go/src/gf_core/client/src/gf_sys_panel.ts'%(cwd_str),
						],
						#-------------
						#LIBS
						'libs_files_lst':[]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_user/client/css/gf_user_profile.css'%(cwd_str),'%s/../go/bin/apps/gf_user/static/css'%(cwd_str))
						]
					},

					#static files to copy without change
					'files_to_copy_lst':[
						('%s/../go/src/apps/gf_user/client/gf_user_profile.html'%(cwd_str),'%s/../go/bin/apps/gf_user/static'%(cwd_str),)
					]
				}
			}
		},
		#-----------------------------
		'gf_landing_page':{
			'pages_map':{
				'main':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../go/src/apps/gf_landing_page/client'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_landing_page/static'%(cwd_str),

					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_landing_page/static/js/gf_landing_page.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_landing_page/static/js/gf_landing_page.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_landing_page/client/src/gf_calc.ts'%(cwd_str),
							'%s/../go/src/apps/gf_landing_page/client/src/gf_email_registration.ts'%(cwd_str),
							'%s/../go/src/apps/gf_landing_page/client/src/gf_images.ts'%(cwd_str),
							'%s/../go/src/apps/gf_landing_page/client/src/gf_landing_page.ts'%(cwd_str),
							'%s/../go/src/apps/gf_landing_page/client/src/gf_procedural_art.ts'%(cwd_str)
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							'%s/../go/src/apps/gf_landing_page/client/lib/jquery.timeago.js'%(cwd_str),
							'%s/../go/src/apps/gf_landing_page/client/lib/color-thief.min.js'%(cwd_str),
						]
						#-------------
					},
					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_landing_page/client/css/domains.css'%(cwd_str)        ,'%s/../go/bin/apps/gf_landing_page/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_landing_page/client/css/gf_landing_page.css'%(cwd_str),'%s/../go/bin/apps/gf_landing_page/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_landing_page/client/css/images.css'%(cwd_str)         ,'%s/../go/bin/apps/gf_landing_page/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_landing_page/client/css/posts.css'%(cwd_str)          ,'%s/../go/bin/apps/gf_landing_page/static/css'%(cwd_str))
						]
					},
				}
			}
		},
		#-----------------------------
		'gf_publisher':{
			'pages_map':{
				'gf_post':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../go/src/apps/gf_publisher/client/src/gf_post'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_publisher/static'%(cwd_str),

					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_publisher/static/js/gf_post.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_publisher/static/js/gf_post.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_publisher/client/src/gf_post/gf_post.ts'%(cwd_str),
							'%s/../go/src/apps/gf_publisher/client/src/gf_post/gf_post_image_view.ts'%(cwd_str),
							'%s/../go/src/apps/gf_publisher/client/src/gf_post/gf_post_tag_mini_view.ts'%(cwd_str),
							'%s/../go/src/apps/gf_tagger/client/src/gf_tagger_client/gf_tagger_client.ts'%(cwd_str),
							'%s/../go/src/apps/gf_tagger/client/src/gf_tagger_client/gf_tagger_input_ui.ts'%(cwd_str),
							'%s/../go/src/gf_core/client/src/gf_sys_panel.ts'%(cwd_str)
						],
						#-------------
						#LIBS
						'libs_files_lst':[
						]
						#-------------
					},

					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_publisher/client/gf_post/web/css/gf_post.css'%(cwd_str)        ,'%s/../go/bin/apps/gf_publisher/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_publisher/client/gf_post/web/css/gf_post_tagging.css'%(cwd_str),'%s/../go/bin/apps/gf_publisher/static/css'%(cwd_str)),
							('%s/../go/src/gf_core/client/css/gf_sys_panel.css'%(cwd_str)                         ,'%s/../go/bin/apps/gf_publisher/static/css'%(cwd_str)),
						]
					}
				},
				#-------------
				'gf_posts_browser':{
					'type_str'         :'ts',
					'code_root_dir_str':'%s/../go/src/apps/gf_publisher/client/src/gf_posts_browser'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_publisher/static'%(cwd_str),

					'ts':{
						'out_file_str'     :'%s/../go/bin/apps/gf_publisher/static/js/gf_posts_browser.js'%(cwd_str),
						'minified_file_str':'%s/../go/bin/apps/gf_publisher/static/js/gf_posts_browser.min.js'%(cwd_str),
						'files_lst':[
							'%s/../go/src/apps/gf_publisher/client/src/gf_posts_browser/gf_posts_browser.ts'%(cwd_str),
							'%s/../go/src/apps/gf_publisher/client/src/gf_posts_browser/gf_posts_browser_view.ts'%(cwd_str),
							'%s/../go/src/apps/gf_publisher/client/src/gf_posts_browser/gf_posts_browser_client.ts'%(cwd_str),
							'%s/../go/src/apps/gf_tagger/client/src/gf_tagger_client/gf_tagger_client.ts'%(cwd_str),
							'%s/../go/src/apps/gf_tagger/client/src/gf_tagger_client/gf_tagger_input_ui.ts'%(cwd_str),
							'%s/../go/src/apps/gf_tagger/client/src/gf_tagger_client/gf_tagger_notes_ui.ts'%(cwd_str),
							'%s/../go/src/gf_core/client/src/gf_sys_panel.ts'%(cwd_str)
						],
						#-------------
						#LIBS
						'libs_files_lst':[
							'%s/../go/src/apps/gf_publisher/client/lib/masonry.pkgd.min.js'%(cwd_str),
							'%s/../go/src/apps/gf_publisher/client/lib/jquery.timeago.js'%(cwd_str),
						]
						#-------------
					},

					'css':{
						'files_lst':[
							('%s/../go/src/apps/gf_publisher/client/gf_posts_browser/web/css/gf_posts_browser.css'%(cwd_str)        ,'%s/../go/bin/apps/gf_publisher/static/css'%(cwd_str)),
							('%s/../go/src/apps/gf_publisher/client/gf_posts_browser/web/css/gf_posts_browser_tagging.css'%(cwd_str),'%s/../go/bin/apps/gf_publisher/static/css'%(cwd_str)),
							('%s/../go/src/gf_core/client/css/gf_sys_panel.css'%(cwd_str)                                           ,'%s/../go/bin/apps/gf_publisher/static/css'%(cwd_str)),
						]
					}
				}
				#-------------
			}
		},
		#-----------------------------
		'gf_tagger':{
			'pages_map':{
				'gf_tag_objects':{
					'type_str'         :'dart',
					'code_root_dir_str':'%s/../go/src/apps/gf_tagger/client/gf_tag_objects'%(cwd_str),
					'target_deploy_dir':'%s/../go/bin/apps/gf_tagger/static'%(cwd_str),
				},
				##IMPORTANT!! - not a page itself, instead its code being used by other pages, but its included here
				##              so that its code gets built when pages for this app are built
				#'gf_tagger_client':{
				#	'code_root_dir_str':'%s/../src/apps/gf_tagger/client/gf_tagger_client'%(cwd_str),
				#	'target_deploy_dir':'%s/../bin/apps/gf_tagger/static'%(cwd_str),
				#}
			}
		}
		#-----------------------------
	}
	return apps_map