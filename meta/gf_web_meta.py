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

	gf_landing_page__pages_map = {
		"gf_landing_page": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_landing_page"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_landing_page/templates/gf_landing_page/gf_landing_page.html"%(modd_str),
			"url_base_str":       "/landing/static",
		}
	}

	gf_images__pages_map = {
		#-------------
		# IMAGES_FLOWS_BROWSER
		"gf_images_flows_browser": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_images"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_images/templates/gf_images_flows_browser/gf_images_flows_browser.html"%(modd_str),
			"url_base_str":       "/images/static",
		},

		#-------------
		# IMAGES_DASHBOARD
		"gf_images_dashboard": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_images"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_images/templates/gf_images_dashboard/gf_images_dashboard.html"%(modd_str),
			"url_base_str":       "/images/static",
		},

		#-------------
	}
	gf_publisher__pages_map = {
		#-------------
		# GF_POST
		"gf_post": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_publisher"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_publisher/templates/gf_post/gf_post.html"%(modd_str),
			"url_base_str":       "/posts/static",
		},

		#-------------
		# GF_POSTS_BROWSER
		"gf_posts_browser": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_publisher"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_publisher/templates/gf_posts_browser/gf_posts_browser.html"%(modd_str),
			"url_base_str":       "/posts/static",
		}

		#-------------
	}




	gf_analytics__pages_map = {
		#-------------
		# DASHBOARD
		"gf_analytics_dashboard":{
			"build_dir_str":      "%s/../web/build/gf_apps/gf_analytics"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_analytics/templates/gf_analytics_dashboard/gf_analytics_dashboard.html"%(modd_str),
			"url_base_str":       "/posts/static",
		},

		#-------------
		# CRAWL_DASHBOARD
		"gf_crawl_dashboard": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_analytics"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_crawl_lib/templates/gf_crawl_dashboard/gf_crawl_dashboard.html"%(modd_str),
			"url_base_str":       "/a/static",
		},

		#-------------
		# DOMAINS_BROWSER

		# IMPORTANT!! - this is in analytics, because domains are sources for images/posts, and so dont 
		#               belong to neither gf_images nor gf_publisher. maybe it should be its own app?
		"gf_domains_browser": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_analytics"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_domains_lib/templates/gf_domains_browser/gf_domains_browser.html"%(modd_str),
			"url_base_str":       "/a/static",
		}

		#-------------
	}

	gf_tagger__pages_map = {
		#-------------
		"gf_tag_objects": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_tagger",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_tagger/templates/gf_tag_objects/gf_tag_objects.html",
			"url_base_str":       "/tags/static",
		},

		#-------------
		# # IMPORTANT!! - not a page itself, instead its code being used by other pages, but its included here
		# #               so that it gets built when pages for this app are built
		# "gf_tagger_client": {
		# 	"code_root_dir_str": "%s/../src/apps/gf_tagger/client/gf_tagger_client"%(modd_str),
		# 	"target_deploy_dir": "%s/../bin/apps/gf_tagger/static"%(modd_str),
		# }

		#-------------
	}

	apps_map = {
		#-----------------------------
		# GF_SOLO
		"gf_solo": {
			
		},

		#-----------------------------
		# GF_LANDING_PAGE
		"gf_landing_page": {
			"pages_map": gf_landing_page__pages_map
		},

		#-----------------------------
		# GF_IMAGES
		"gf_images": {
			"pages_map": gf_images__pages_map
		},

		#-----------------------------
		# GF_PUBLISHER
		"gf_publisher": {
			"pages_map": gf_publisher__pages_map
		},

		#-----------------------------
		# GF_ANALYTICS
		"gf_analytics": {
			"pages_map": gf_analytics__pages_map
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
		#
		#-----------------------------
	}


	#-----------------------------
	# GF_SOLO - contains pages of all apps

	import copy # IMPORTANT!! - do a deepcopy of pages_map, because for gf_solo we modify them

	gf_solo__pages_map = {

		#-------------
		# BOOKMARKLET
		# FIX!! - build this as a page outside of gf_solo, because if gf_solo (monolith service)
		#         is not how GF is deployed, then gf_bookmarklet wont get built.
		#         figure out some general solution for build_dir (not gf_images), since its not just images that
		#         are manipulated with this bookmarklet but also bookmarks/tags/etc.
		"gf_bookmarklet": {
			"build_dir_str": "%s/../web/build/gf_apps/gf_images"%(modd_str),
			"files_to_copy_lst": [
				(f"{modd_str}/../web/src/gf_apps/gf_images/ts/gf_bookmarklet/gf_bookmarklet.js", f"{modd_str}/../web/build/gf_apps/gf_solo/gf_images/js")
			]
		}

		#-------------
	}
	
	gf_solo__pages_map.update(copy.deepcopy(gf_landing_page__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_images__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_publisher__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_analytics__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_tagger__pages_map))

	for _, page_info_map in gf_solo__pages_map.items():
		page_info_map["build_copy_dir_str"] = f"{modd_str}/../web/build/gf_apps/gf_solo"

	apps_map["gf_solo"] = {"pages_map": gf_solo__pages_map}

	#-----------------------------


	return apps_map