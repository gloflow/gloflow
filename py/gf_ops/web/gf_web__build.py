# GloFlow application and media management/publishing platform
# Copyright (C) 2019 Ivan Trajkovic
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
modd_str = os.path.abspath(os.path.dirname(__file__))

import os
from colored import fg, bg, attr
from bs4 import BeautifulSoup

sys.path.append("%s/../../gf_core"%(modd_str))
import gf_core_cli

#--------------------------------------------------
# RUN_IN_CONTAINER
def run_in_cont(p_app_str,
	p_page_name_str=None):

	repo_local_path_str = os.path.abspath(f'{modd_str}/../../../../gloflow').strip()

	py_cmd_lst = [
		"python3", "-u", "/home/gf/ops/cli__build.py", "-run=build_web", "-build_outof_cont",
		f"-app={p_app_str}"
	]
	if not p_page_name_str == None:
		py_cmd_lst.append(f"-page_name={p_page_name_str}")

	cmd_lst = [
		"sudo", "docker", "run",
		"--rm", # remove after exit 
		"-v", f"{repo_local_path_str}:/home/gf", # mount repo into the container
		"glofloworg/gf_builder_web:latest",
	]
	cmd_lst.extend(py_cmd_lst)
	
	p = gf_core_cli.run__view_realtime(cmd_lst, {},
		"gf_build_web", "green")

	p.wait()

#---------------------------------------------------
def build(p_apps_names_lst,
	p_apps_meta_map,
	p_log_fun,
	p_page_name_str=None):
	assert isinstance(p_apps_names_lst, list)
	assert len(p_apps_names_lst) > 0
	assert isinstance(p_apps_meta_map, dict)

	#---------------------------------------------------
	def individual_page(p_page_name_str, p_page_info_map):
		print(f"page name {fg('yellow')}{p_page_name_str}{attr(0)}")

		build_dir_str = os.path.abspath(p_page_info_map["build_dir_str"])

		build_copy_dir_str     = p_page_info_map.get("build_copy_dir_str", None)
		build_copy_bas_dir_str = None
		if not build_copy_dir_str == None:
			build_copy_bas_dir_str = os.path.abspath(build_copy_dir_str)

		build_page(p_page_name_str,
			build_dir_str,
			build_copy_bas_dir_str,
			p_page_info_map,
			p_log_fun)

	#---------------------------------------------------
	for app_str in p_apps_names_lst:
		
		#-----------------
		# META
		if not app_str in p_apps_meta_map.keys():
			p_log_fun("ERROR", f"supplied app ({app_str}) does not exist in gf_web_meta")
			return

		app_map = p_apps_meta_map[app_str]

		#-----------------
		# BUILD PAGES - build each page of the app.
		#               no specific page was picked for build
		if p_page_name_str == None:
			
			for page_name_str, page_info_map in app_map["pages_map"].items():
				individual_page(page_name_str, page_info_map)
				
		else:
			page_info_map = app_map["pages_map"][p_page_name_str]

			individual_page(p_page_name_str, page_info_map)

#---------------------------------------------------
# BUILD_PAGE

def build_page(p_page_name_str,
	p_build_dir_str,
	p_build_copy_dir_str,
	p_page_info_map,
	p_log_fun):
	
	print("")
	p_log_fun("INFO", "%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s"%(fg("orange_red_1"), attr(0)))
	p_log_fun("INFO", "             %sBUILD PAGE%s - %s%s%s"%(fg("cyan"), attr(0), fg("orange_red_1"), p_page_name_str, attr(0)))
	p_log_fun("INFO", "%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s"%(fg("orange_red_1"), attr(0)))
	print("")

	p_log_fun("INFO", f"build_dir_str - {p_build_dir_str}")
	assert isinstance(p_build_dir_str, str)

	# make build dir if it doesnt exist
	if not os.path.isdir(p_build_dir_str):
		gf_core_cli.run(f"mkdir -p {p_build_dir_str}")

	if "main_html_path_str" in p_page_info_map.keys():

		main_html_path_str = os.path.abspath(p_page_info_map["main_html_path_str"])
		
		print(f"main html path - {main_html_path_str}")

		assert os.path.isfile(main_html_path_str)
		assert main_html_path_str.endswith(".html")
		assert ".".join(os.path.basename(main_html_path_str).split(".")[:-1]) == p_page_name_str
		p_log_fun("INFO", f"main_html_path_str - {main_html_path_str}")

	# URL_BASE
	url_base_str = ""
	if "url_base_str" in p_page_info_map.keys():
		url_base_str = p_page_info_map["url_base_str"]
		p_log_fun("INFO", f"url_base_str - {url_base_str}")

	if "main_html_path_str" in p_page_info_map.keys():	
		f = open(main_html_path_str, "r")
		main_html_str = f.read()
		f.close()

		soup = BeautifulSoup(main_html_str)
		
	#---------------------------------------------------
	# JS
	
	def process_scripts():
		scripts_dom_nodes_lst = soup.findAll("script")

		# if there are scripts detected in the page
		if len(scripts_dom_nodes_lst) > 0:
			js_libs_build_dir_str = f"{p_build_dir_str}/js/lib"
			gf_core_cli.run(f"mkdir -p {js_libs_build_dir_str}") # create dir and all parent dirs

		main_html_dir_path_str = os.path.dirname(main_html_path_str)
		assert os.path.isdir(main_html_dir_path_str)

		for script_dom_node in scripts_dom_nodes_lst:

			# some <script> tags might just contain source code, and not reference an external JS file.
			# using .get("src") instead of ["src"] because the "src" DOM attribute might not be present.
			if script_dom_node.get("src") == None:
				continue

			src_str = script_dom_node["src"]

			if src_str.startswith("http://") or src_str.startswith("https://"):
				print("EXTERNAL_URL - DO NOTHING")
				continue
			

			
			
			local_path_str = os.path.abspath(f"{main_html_dir_path_str}/{src_str}")
			print(local_path_str)
			assert os.path.isfile(local_path_str)

			#-----------------
			if local_path_str.endswith(".ts"):
				p_log_fun("INFO", "%s------------ TYPESCRIPT --------------------------------%s"%(fg("yellow"), attr(0)))
				
				#---------------------------------------------------
				def build_typescript(p_out_file_str):
					
					cmd_lst = [
						
						"tsc",
						"--module system", # needed with the "--out" option
						"--target ES2020", # "--target es2017", # "--target es6",

						# Enables emit interoperability between CommonJS and ES Modules via creation of namespace objects for all imports
						# '--esModuleInterop',
						f"--outFile {p_out_file_str}",
						main_ts_file_str
					]
					cmd_str = " ".join(cmd_lst)
					print(cmd_str)
					
					_, _, return_code_int = gf_core_cli.run(cmd_str)

					if return_code_int > 0:
						print("ERROR!! - TypeScript Compilation failed!")
						exit(-1)

					# minify into the same file name as the Typescript compiler output
					target_dir_str = os.path.dirname(p_out_file_str)
					minify_js(p_out_file_str, [p_out_file_str], p_log_fun)

				#---------------------------------------------------

				main_ts_file_str       = local_path_str
				minified_file_name_str = "%s.min.js"%(".".join(os.path.basename(main_ts_file_str).split(".")[:-1]))
				minified_file_path_str = f"{p_build_dir_str}/js/{minified_file_name_str}"

				build_typescript(minified_file_path_str)

				# HTML_MODIFY - change the src in the html tag to the minified name, and url_base (dont leave relative path)
				script_dom_node["src"] = f"{url_base_str}/js/{minified_file_name_str}"

			#-----------------
			elif local_path_str.endswith(".js"):
				
				p_log_fun("INFO", "%s------------ JAVASCRIPT --------------------------------%s"%(fg("yellow"), attr(0)))

				# IMPORTANT!! - just copy the JS file to the final build dir
				gf_core_cli.run(f"cp {local_path_str} {js_libs_build_dir_str}")

				# HTML_MODIFY - change the src in the html tag, to include the url_base (dont leave relative path)
				script_dom_node["src"] = f"{url_base_str}/js/lib/{os.path.basename(local_path_str)}"

			#-----------------

	#---------------------------------------------------
	# CSS
	
	def process_css():

		p_log_fun("INFO", "%s------------ CSS ---------------------------------------%s"%(fg("yellow"), attr(0)))
		css_links_lst = soup.findAll("link", {"type": "text/css"})
		
		target_dir_str = f'{p_build_dir_str}/css/{p_page_name_str}'
		gf_core_cli.run(f'mkdir -p {target_dir_str}') #create dir and all parent dirs

		for css in css_links_lst:
			src_str = css["href"]
			
			assert src_str.endswith('.css') or src_str.endswith('.scss')

			if src_str.startswith('http://') or src_str.startswith('https://'):
				print('EXTERNAL_URL - DO NOTHING')
				continue

			# full paths are relative to the dir holding the main html file (app entry point)
			full_path_str = os.path.abspath(f'{os.path.dirname(main_html_path_str)}/{src_str}')
			print(full_path_str)
			assert os.path.isfile(full_path_str)
			
			# SASS
			if src_str.endswith('.scss'):
				css_file_name_str = os.path.basename(src_str).replace('.scss', '.css')
				final_src_str     = f'{target_dir_str}/{css_file_name_str}'
				gf_core_cli.run(f'sass {full_path_str} {final_src_str}')

				# HTML_MODIFY - change the src in the html tag, to include the url_base
				#               (dont leave relative path)
				css["href"] = f'{url_base_str}/css/{p_page_name_str}/{css_file_name_str}'

			# CSS
			else:
				
				gf_core_cli.run(f'cp {full_path_str} {target_dir_str}')

				# HTML_MODIFY - change the src in the html tag, to include the url_base (dont leave relative path)
				css["href"] = f'{url_base_str}/css/{p_page_name_str}/{os.path.basename(full_path_str)}'

	#---------------------------------------------------

	if "main_html_path_str" in p_page_info_map.keys():

		process_scripts()

		#-----------------
		# CSS
		css_process_bool = p_page_info_map.get("css_process_bool", True)
		if css_process_bool:
			process_css()

		#-----------------
		# CREATE_FINAL_MODIFIED_HTML - create the html template file in the build dir that contains all 
		#                              the modified urls for JS/CSS
		target_html_file_path_str = f'{p_build_dir_str}/templates/{p_page_name_str}/{p_page_name_str}.html'
		gf_core_cli.run(f'mkdir -p {os.path.dirname(target_html_file_path_str)}')

		f = open(target_html_file_path_str, 'w+')
		f.write(soup.prettify())
		f.close()

		#-----------------

	#-----------------
	# SUBTEMPLATES
	if "subtemplates_lst" in p_page_info_map.keys():
		process_subtemplates(p_page_name_str,
			p_build_dir_str,
			p_page_info_map,
			p_log_fun)

	#-----------------
	# IMPORTANT!! - do after build_copy_dir is created
	if "files_to_copy_lst" in p_page_info_map.keys():
		process_files_to_copy(p_page_info_map, p_log_fun)

	#-----------------
	# BUILD_COPY - this propety allows for the build dir of a page to be copied to some other dir after the build is complete.
	#              this has to run after all other build steps complete, so that it includes all the build artifacts.
	#
	# IMPORTANT!! - only some pages in some apps define this. gf_solo is one of these apps, it adds this property
	#               to the page defs of all other apps (since gf_solo includes all apps).
	if not p_build_copy_dir_str == None:
		print(f"copying {fg('green')}build{attr(0)} dir ({p_build_dir_str}) to {fg('yellow')}{p_build_copy_dir_str}{attr(0)}")
		
		gf_core_cli.run(f'mkdir -p {p_build_copy_dir_str}')
		gf_core_cli.run(f'cp -r {p_build_dir_str} {p_build_copy_dir_str}')

	#-----------------
	print("")
	p_log_fun("INFO", "%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s END"%(fg("orange_red_1"), attr(0)))
	print("")

#---------------------------------------------------
def process_files_to_copy(p_page_info_map, p_log_fun):
	assert isinstance(p_page_info_map, dict)

	print("")
	p_log_fun("INFO", "%s------------ COPY_FILES --------------------------------%s"%(fg("yellow"), attr(0)))
	print("")

	files_to_copy_lst = p_page_info_map["files_to_copy_lst"]
	assert isinstance(files_to_copy_lst, list)


	# COPY_FILES
	for file_to_copy_tpl in files_to_copy_lst:
		assert isinstance(file_to_copy_tpl, tuple)
		src_file_str, target_dir_str = file_to_copy_tpl

		assert os.path.isfile(src_file_str)
		
		if not os.path.isdir(target_dir_str):
			gf_core_cli.run(f'mkdir -p {target_dir_str}')


		gf_core_cli.run(f"cp {src_file_str} {target_dir_str}")

		final_path_str = f"{target_dir_str}/{os.path.basename(src_file_str)}"
		assert os.path.isfile(final_path_str)

#---------------------------------------------------
def process_subtemplates(p_page_name_str,
	p_build_dir_str,
	p_page_info_map,
	p_log_fun):
	assert isinstance(p_page_name_str, str)
	assert os.path.isdir(p_build_dir_str)
	assert "subtemplates_lst" in p_page_info_map.keys()

	print("")
	p_log_fun("INFO", "%s------------ SUBTEMPLATES --------------------------------%s"%(fg("yellow"), attr(0)))
	print("")

	subtemplates_lst = p_page_info_map["subtemplates_lst"]
	assert isinstance(subtemplates_lst, list)

	# SUBTEMPLATES__BUILD_DIR
	target_subtemplates_build_dir_str = f"{p_build_dir_str}/templates/{p_page_name_str}/subtemplates"
	if not os.path.isdir(target_subtemplates_build_dir_str):
		gf_core_cli.run(f"mkdir -p {target_subtemplates_build_dir_str}")

	for s_path_str in subtemplates_lst:
		print(s_path_str)
		assert isinstance(s_path_str, str)
		assert os.path.isfile(s_path_str)
		assert s_path_str.endswith(".html")

		# SUBTEMPLATE__COPY
		gf_core_cli.run(f"cp {s_path_str} {target_subtemplates_build_dir_str}")

#---------------------------------------------------
def minify_js(p_js_target_file_str,
	p_js_files_lst,
	p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.minify_js()")

	cmd_lst = [
		"uglifyjs",
		f"--output {p_js_target_file_str}",
		" ".join(p_js_files_lst),
	]
	gf_core_cli.run(" ".join(cmd_lst))

#---------------------------------------------------