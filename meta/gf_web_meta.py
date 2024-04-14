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
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_landing_page",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_landing_page/templates/gf_landing_page/gf_landing_page.html",
			"url_base_str":       "/landing/static",
		}
	}

	gf_images__pages_map = {
		#-------------
		# IMAGES_VIEW
		"gf_images_view": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_images",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_images/templates/gf_images_view/gf_images_view.html",
			"url_base_str":       "/images/static",
		},

		#-------------
		# IMAGES_FLOWS_BROWSER
		"gf_images_flows_browser": {
			"build_dir_str":      "%s/../web/build/gf_apps/gf_images"%(modd_str),
			"main_html_path_str": "%s/../web/src/gf_apps/gf_images/templates/gf_images_flows_browser/gf_images_flows_browser.html"%(modd_str),
			"url_base_str":       "/images/static",
		},

		#-------------
		# PAGE_PICKER
		# FIX!! - figure out some general solution for build_dir (not gf_images), since its not just images that
		#         are manipulated with this bookmarklet but also bookmarks/tags/etc.
		#         this should possibly go into the tagger
		"gf_page_picker": {
			"build_dir_str": f"{modd_str}/../web/build/gf_apps/gf_images",
			"files_to_copy_lst": [
				(f"{modd_str}/../web/src/gf_apps/gf_images/ts/gf_page_picker/gf_page_picker.js", f"{modd_str}/../web/build/gf_apps/gf_images/js"),
				
				# FIX!! - these icons are used by all apps. figure a different place to insert them separate of any app.
				(f"{modd_str}/../web/assets/gf_bar_handle_btn.svg",        f"{modd_str}/../web/build/gf_apps/gf_images/assets"),
				(f"{modd_str}/../web/assets/gf_close_btn_small.svg",       f"{modd_str}/../web/build/gf_apps/gf_images/assets"),
				(f"{modd_str}/../web/assets/gf_metamask_icon.svg",         f"{modd_str}/../web/build/gf_apps/gf_images/assets"),
				(f"{modd_str}/../web/assets/gf_copy_to_clipboard_btn.svg", f"{modd_str}/../web/build/gf_apps/gf_images/assets"),
				(f"{modd_str}/../web/assets/gf_add_btn.svg",               f"{modd_str}/../web/build/gf_apps/gf_images/assets"),
				(f"{modd_str}/../web/assets/gf_confirm_btn.svg",           f"{modd_str}/../web/build/gf_apps/gf_images/assets")
			]
		},

		#-------------
		# CODE_EDITOR
		

		"gf_code_editor": {
			"build_dir_str": f"{modd_str}/../web/build/gf_apps/gf_images",
			"files_to_copy_lst": [
				(f"{modd_str}/../web/src/gf_apps/gf_code_editor/templates/code_editor.html", f"{modd_str}/../web/build/gf_apps/gf_images"),
			]
		}

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
		# DOMAINS_BROWSER

		# IMPORTANT!! - this is in analytics, because domains are sources for images/posts, and so dont 
		#               belong to neither gf_images nor gf_publisher. maybe it should be its own core app?
		"gf_domains_browser": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_analytics",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_domains_lib/templates/gf_domains_browser/gf_domains_browser.html",
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
		"gf_bookmarks": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_tagger",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_tagger/templates/gf_bookmarks/gf_bookmarks.html",
			"url_base_str":       "/tags/static",
		},

		#-------------
	}


	gf_admin__pages_map = {
		#-------------
		"gf_admin_dashboard": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_admin",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_admin/templates/gf_admin_dashboard/gf_admin_dashboard.html",
			"url_base_str":       "/v1/admin/static",
		},

		#-------------
		"gf_admin_login": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_admin",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_admin/templates/gf_admin_login/gf_admin_login.html",
			"url_base_str":       "/v1/admin/static",
		}

		#-------------
	}

	gf_home__pages_map = {
		#-------------
		"gf_home_main": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_apps/gf_home",
			"main_html_path_str": f"{modd_str}/../web/src/gf_apps/gf_home/templates/gf_home_main/gf_home_main.html",
			"url_base_str":       "/v1/home/static",
		},

		#-------------
	}

	gf_identity__pages_map = {
		#-------------
		"gf_login": {
			"build_dir_str":      f"{modd_str}/../web/build/gf_identity",
			"main_html_path_str": f"{modd_str}/../web/src/gf_identity/templates/gf_login/gf_login.html",
			"url_base_str":       "/v1/identity/static",
		},

		#-------------
	}

	apps_map = {
		#-----------------------------
		# GF_SOLO
		"gf_solo": {},

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
	}


	#-----------------------------
	# GF_SOLO - contains pages of all apps

	import copy # IMPORTANT!! - do a deepcopy of pages_map, because for gf_solo we modify them

	gf_solo__pages_map = {}

	gf_solo__pages_map.update(copy.deepcopy(gf_identity__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_admin__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_home__pages_map))

	gf_solo__pages_map.update(copy.deepcopy(gf_landing_page__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_images__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_publisher__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_analytics__pages_map))
	gf_solo__pages_map.update(copy.deepcopy(gf_tagger__pages_map))

	for _, page_info_map in gf_solo__pages_map.items():

		# only for gf_solo pages is this dir defined. so that the build_dir_str of other 
		# apps/pages still gets coppied into the build dir of gf_solo.
		page_info_map["build_copy_dir_str"] = f"{modd_str}/../web/build/gf_apps/gf_solo"

	apps_map["gf_solo"] = {"pages_map": gf_solo__pages_map}

	#-----------------------------


	return apps_map