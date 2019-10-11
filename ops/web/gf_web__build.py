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

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import os
from colored import fg,bg,attr
import BeautifulSoup as bs

sys.path.append('%s/../utils'%(cwd_str))
import gf_cli_utils as gf_u

#---------------------------------------------------
def build(p_apps_names_lst, p_apps_meta_map, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.build()")
	assert isinstance(p_apps_names_lst, list)
	assert len(p_apps_names_lst) > 0
	assert isinstance(p_apps_meta_map, dict)

	for app_str in p_apps_names_lst:
		
		#-----------------
		#META
		if not p_apps_meta_map.has_key(app_str):
			p_log_fun("ERROR", "supplied app (%s) does not exist in gf_web_meta"%(app_str))
			return

		app_map = p_apps_meta_map[app_str]
		#-----------------

		#BUILD PAGES - build each page of the app
		for page_name_str, page_info_map in app_map['pages_map'].items():

			build_dir_str = os.path.abspath(page_info_map['build_dir_str'])

			build_page(page_name_str,
				build_dir_str,
				page_info_map,
				p_log_fun)

#---------------------------------------------------
def build_page(p_page_name_str,
    p_target_build_dir_str,
    p_page_info_map,
    p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.build_page()")
	
	print('')
	print('')
	p_log_fun('INFO', '%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s'%(fg('orange_red_1'), attr(0)))
	p_log_fun('INFO', '             %sBUILD PAGE%s - %s%s%s'%(fg('cyan'), attr(0), fg('orange_red_1'), p_page_name_str, attr(0)))
	p_log_fun('INFO', '%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s'%(fg('orange_red_1'), attr(0)))
	print('')
	print('')

	p_log_fun('INFO', 'build_dir_str - %s'%(p_target_build_dir_str))
	assert isinstance(p_target_build_dir_str, basestring)
	assert os.path.isdir(p_target_build_dir_str)

	if p_page_info_map.has_key('main_html_path_str'):
		main_html_path_str = os.path.abspath(p_page_info_map['main_html_path_str'])
		assert os.path.isfile(main_html_path_str)
		assert main_html_path_str.endswith('.html')
		assert '.'.join(os.path.basename(main_html_path_str).split('.')[:-1]) == p_page_name_str
		p_log_fun('INFO', 'main_html_path_str - %s'%(main_html_path_str))

	if p_page_info_map.has_key('url_base_str'):
		url_base_str = p_page_info_map['url_base_str']
		p_log_fun('INFO', 'url_base_str - %s'%(url_base_str))

	if p_page_info_map.has_key('main_html_path_str'):	
		f = open(main_html_path_str, 'r')
		main_html_str = f.read()
		f.close()

		soup = bs.BeautifulSoup(main_html_str)
		
	#---------------------------------------------------
	def process_scripts():
		scripts_lst = soup.findAll('script')

		#if there are scripts detected in the page
		if len(scripts_lst) > 0:
			js_libs_build_dir_str = '%s/js/lib'%(p_target_build_dir_str)
			gf_u.run_cmd('mkdir -p %s'%(js_libs_build_dir_str)) #create dir and all parent dirs

		for script in scripts_lst:

			#some <script> tags might just contain source code, and not reference an external JS file
			if not script.has_key('src'):
				continue

			src_str = script['src']
			
			if src_str.startswith('http://') or src_str.startswith('https://'):
				print('EXTERNAL_URL - DO NOTHING')
				continue
			
			main_html_dir_path_str = os.path.dirname(main_html_path_str)
			assert os.path.isdir(main_html_dir_path_str)
			
			local_path_str = os.path.abspath('%s/%s'%(main_html_dir_path_str, src_str))
			print(local_path_str)
			assert os.path.isfile(local_path_str)

			#-----------------
			if local_path_str.endswith('.ts'):
				p_log_fun('INFO', '%s------------ TYPESCRIPT --------------------------------%s'%(fg('yellow'), attr(0)))
				
				#---------------------------------------------------
				def build_typescript(p_out_file_str):
					
					cmd_lst = [
						'tsc',
						'--module system', #needed with the "--out" option

						#Enables emit interoperability between CommonJS and ES Modules via creation of namespace objects for all imports
						#'--esModuleInterop',
						'--out %s'%(p_out_file_str),
						main_ts_file_str
					]
					gf_u.run_cmd(' '.join(cmd_lst))

					#minify into the same file name as the Typescript compiler output
					target_dir_str = os.path.dirname(p_out_file_str)
					minify_js(p_out_file_str, [p_out_file_str], p_log_fun)

				#---------------------------------------------------

				main_ts_file_str       = local_path_str
				minified_file_name_str = '%s.min.js'%('.'.join(os.path.basename(main_ts_file_str).split('.')[:-1]))
				minified_file_path_str = '%s/js/%s'%(p_target_build_dir_str, minified_file_name_str)

				build_typescript(minified_file_path_str)

				#HTML_MODIFY - change the src in the html tag to the minified name, and url_base (dont leave relative path)
				script['src'] = '%s/js/%s'%(url_base_str, minified_file_name_str)
			#-----------------
			elif local_path_str.endswith('.js'):
				p_log_fun('INFO', '%s------------ JAVASCRIPT --------------------------------%s'%(fg('yellow'), attr(0)))

				#IMPORTANT!! - JS files are currently used for libraries only, so just copy the JS file to the final build dir
				gf_u.run_cmd('cp %s %s'%(local_path_str, js_libs_build_dir_str))

				#HTML_MODIFY - change the src in the html tag, to include the url_base (dont leave relative path)
				script['src'] = '%s/js/lib/%s'%(url_base_str, os.path.basename(local_path_str))
			#-----------------

	#---------------------------------------------------
	def process_css():
		p_log_fun('INFO', '%s------------ CSS ---------------------------------------%s'%(fg('yellow'), attr(0)))
		css_links_lst = soup.findAll('link', {'type':'text/css'})
		
		target_dir_str = '%s/css/%s'%(p_target_build_dir_str, p_page_name_str)
		gf_u.run_cmd('mkdir -p %s'%(target_dir_str)) #create dir and all parent dirs

		for css in css_links_lst:
			src_str = css['href']
			assert src_str.endswith('.css') or src_str.endswith('.scss')

			if src_str.startswith('http://') or src_str.startswith('https://'):
				print('EXTERNAL_URL - DO NOTHING')
				continue

			#full paths are relative to the dir holding the main html file (app entry point)
			full_path_str = os.path.abspath('%s/%s'%(os.path.dirname(main_html_path_str), src_str))
			print(full_path_str)
			assert os.path.isfile(full_path_str)
			
			#SASS
			if src_str.endswith('.scss'):
				css_file_name_str = os.path.basename(src_str).replace('.scss', '.css')
				final_src_str     = '%s/%s'%(target_dir_str, css_file_name_str)
				gf_u.run_cmd('sass %s %s'%(full_path_str, final_src_str))

				#HTML_MODIFY - change the src in the html tag, to include the url_base (dont leave relative path)
				css['href'] = '%s/css/%s/%s'%(url_base_str, p_page_name_str, css_file_name_str)

			#CSS
			else:
				
				gf_u.run_cmd('cp %s %s'%(full_path_str, target_dir_str))

				#HTML_MODIFY - change the src in the html tag, to include the url_base (dont leave relative path)
				css['href'] = '%s/css/%s/%s'%(url_base_str, p_page_name_str, os.path.basename(full_path_str))

	#---------------------------------------------------

	if p_page_info_map.has_key('main_html_path_str'):
		process_scripts()
		process_css()

		#-----------------
		#CREATE_FINAL_MODIFIED_HTML - create the html template file in the build dir that contains all 
		#                             the modified urls for JS/CSS
		target_html_file_path_str = '%s/templates/%s/%s.html'%(p_target_build_dir_str, p_page_name_str, p_page_name_str)
		gf_u.run_cmd('mkdir -p %s'%(os.path.dirname(target_html_file_path_str)))

		f = open(target_html_file_path_str, 'w+')
		f.write(soup.prettify())
		f.close()
		#-----------------

	#-----------------
	#SUBTEMPLATES
	if p_page_info_map.has_key('subtemplates_lst'):
		process_subtemplates(p_page_name_str,
			p_target_build_dir_str,
			p_page_info_map,
			p_log_fun)
	#-----------------

	if p_page_info_map.has_key('files_to_copy_lst'):
		process_files_to_copy(p_page_info_map, p_log_fun)
		
#---------------------------------------------------
def process_files_to_copy(p_page_info_map, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.process_files_to_copy()")
	assert isinstance(p_page_info_map, dict)

	print('')
	p_log_fun('INFO', '%s------------ COPY_FILES --------------------------------%s'%(fg('yellow'), attr(0)))
	print('')

	files_to_copy_lst = p_page_info_map['files_to_copy_lst']
	assert isinstance(files_to_copy_lst, list)

	#COPY_FILES
	for file_to_copy_tpl in files_to_copy_lst:
		assert isinstance(file_to_copy_tpl, tuple)
		src_file_str, target_dir_str = file_to_copy_tpl

		assert os.path.isfile(src_file_str)
		assert os.path.isdir(target_dir_str)

		gf_u.run_cmd('cp %s %s'%(src_file_str, target_dir_str))
		assert os.path.isfile('%s/%s'%(target_dir_str, os.path.basename(src_file_str)))

#---------------------------------------------------
def process_subtemplates(p_page_name_str,
	p_target_build_dir_str, 
	p_page_info_map,
	p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.process_subtemplates()")
	assert isinstance(p_page_name_str, basestring)
	assert os.path.isdir(p_target_build_dir_str)
	assert p_page_info_map.has_key('subtemplates_lst')

	print('')
	p_log_fun('INFO', '%s------------ SUBTEMPLATES --------------------------------%s'%(fg('yellow'), attr(0)))
	print('')

	subtemplates_lst = p_page_info_map['subtemplates_lst']
	assert isinstance(subtemplates_lst, list)

	#SUBTEMPLATES__BUILD_DIR
	target_subtemplates_build_dir_str = '%s/templates/%s/subtemplates'%(p_target_build_dir_str, p_page_name_str)
	if not os.path.isdir(target_subtemplates_build_dir_str):
		gf_u.run_cmd('mkdir -p %s'%(target_subtemplates_build_dir_str))

	for s_path_str in subtemplates_lst:
		print(s_path_str)
		assert isinstance(s_path_str, basestring)
		assert os.path.isfile(s_path_str)
		assert s_path_str.endswith('.html')

		#SUBTEMPLATE__COPY
		gf_u.run_cmd('cp %s %s'%(s_path_str, target_subtemplates_build_dir_str))

#---------------------------------------------------
def minify_js(p_js_target_file_str,
    p_js_files_lst,
    p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.minify_js()")

	cmd_lst = [
		'uglifyjs',
		'--output %s'%(p_js_target_file_str),
		' '.join(p_js_files_lst),
	]
	gf_u.run_cmd(' '.join(cmd_lst))

#---------------------------------------------------
# def build_page(p_page_name_str,
#     p_build_dir_str,
#     p_page_info_map,
#     p_log_fun):
# 	p_log_fun("FUN_ENTER", "gf_web__build.build_page()")
# 	assert isinstance(p_build_dir_str, basestring)
# 	assert os.path.isdir(p_build_dir_str)
#
# 	print('')
# 	print('')
# 	p_log_fun('INFO', '%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s'%(fg('orange_red_1'), attr(0)))
# 	p_log_fun('INFO', '             %sBUILD PAGE%s - %s%s%s'%(fg('cyan'), attr(0), fg('orange_red_1'), p_page_name_str, attr(0)))
# 	p_log_fun('INFO', 'page type - %s'%(p_page_info_map['type_str']))
# 	p_log_fun('INFO', '%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s'%(fg('orange_red_1'), attr(0)))
# 	p_log_fun('INFO', 'build_dir_str - %s'%(p_build_dir_str))
#
# 	#---------------------------------------------------
# 	def build_typescript(p_out_file_str, p_minified_file_str, p_ts_files_lst):
# 		p_log_fun("FUN_ENTER", "gf_web__build.build_page().build_typescript()")
#
# 		cmd_lst = [
# 			'tsc',
# 			'--out %s'%(p_out_file_str),
# 			' '.join(p_ts_files_lst)
# 		]
# 		gf_u.run_cmd(' '.join(cmd_lst))
#
# 		#minify into the same file name as the Typescript compiler output
# 		minify_js(p_minified_file_str, [p_out_file_str], p_log_fun)
# 	#---------------------------------------------------
#
# 	#-----------------
# 	if p_page_info_map['type_str'] == 'ts':
# 		p_log_fun('INFO', '%s------------ TYPESCRIPT --------------------------------%s'%(fg('yellow'), attr(0)))
#
# 		out_file_str      = p_page_info_map['ts']['out_file_str']
# 		minified_file_str = p_page_info_map['ts']['minified_file_str']
# 		ts_files_lst      = p_page_info_map['ts']['files_lst']
#
# 		build_typescript(out_file_str, minified_file_str, ts_files_lst)
# 		#-----------------
# 		#COPY LIBS
#
# 		if p_page_info_map['ts'].has_key('libs_files_lst'):
# 			p_log_fun('INFO', '%s------------ TS_LIBS -----------------------------%s'%(fg('yellow'), attr(0)))
#
# 			if not os.path.isdir(p_build_dir_str): gf_u.run_cmd('mkdir -p %s/js/lib'%(p_build_dir_str))
#			
# 			for lib_file_str in p_page_info_map['ts']['libs_files_lst']:
# 				gf_u.run_cmd('cp %s %s/js/lib'%(lib_file_str, p_build_dir_str))
# 		#-----------------
# 	#-----------------
# 	#if p_page_info_map['type_str'] == 'js':
# 	#	p_log_fun('INFO','%s------------ JAVASCRIPT --------------------------------%s'%(fg('yellow'),attr(0))))
# 	#	#-----------------
# 	#	#MINIFY JS
# 	#	if p_page_info_map['js'].has_key('minified_file_str'):
# 	#		p_log_fun('INFO','------------ JS_MINIFY ---------------------------')
# 	#		minified_file_str = p_page_info_map['js']['minified_file_str']
# 	#		js_files_lst      = p_page_info_map['js']['files_lst']
# 	#
# 	#		minify_js(minified_file_str, js_files_lst, p_log_fun)
# 	#	#-----------------
# 	#	#COPY JS LIBS
# 	#
# 	#	if p_page_info_map['js'].has_key('libs_files_lst'):
# 	#		p_log_fun('INFO','------------ JS_LIBS -----------------------------')
# 	#
# 	#		for lib_file_str in p_page_info_map['js']['libs_files_lst']:
# 	#			gf_u.run_cmd('cp %s %s/js/lib'%(lib_file_str, p_build_dir_str))
# 	#-----------------
# 	#CSS
#
# 	if p_page_info_map.has_key('css'):
# 		p_log_fun('INFO', '%s------------ CSS ---------------------------------%s'%(fg('yellow'), attr(0)))
# 		css_files_lst = p_page_info_map['css']['files_lst']
#
# 		for f_tpl in css_files_lst:
# 			assert len(f_tpl) == 2
#
# 			src_file_str, dest_dir_src = f_tpl
#
# 			if not os.path.isdir(dest_dir_src): gf_u.run_cmd('mkdir -p %s'%(dest_dir_src))
# 			gf_u.run_cmd('cp %s %s'%(src_file_str, dest_dir_src))
# 	#-----------------
# 	#TEMPLATES
# 	if p_page_info_map.has_key('templates'):
# 		p_log_fun('INFO', '%s------------ TEMPLATES -----------------------------%s'%(fg('yellow'), attr(0)))
#		
# 		assert p_page_info_map['templates'].has_key('files_lst')
# 		templates_files_lst = p_page_info_map['templates']['files_lst']
# 		assert isinstance(templates_files_lst, list)
#
#
# 		#COPY_TEMPLATE_FILES - copy them from their source location to the desired build location
# 		for tmpl_file_str, tmpl_target_dir_str in templates_files_lst:
# 			print('tmpl_file_str - %s'%(tmpl_file_str))
#
# 			assert os.path.isfile(tmpl_file_str)
#
# 			#if target template dir doesnt exist, create it
# 			if not os.path.isdir(tmpl_target_dir_str): gf_u.run_cmd('mkdir -p %s'%(tmpl_target_dir_str))
#
# 			gf_u.run_cmd('cp %s %s'%(tmpl_file_str, tmpl_target_dir_str))
# 	#-----------------
# 	#FILES_TO_COPY
#
# 	if p_page_info_map.has_key('files_to_copy_lst'):
# 		p_log_fun('INFO', '%s------------ FILES_TO_COPY -----------------------%s'%(fg('yellow'), attr(0)))
# 		files_to_copy_lst = p_page_info_map['files_to_copy_lst']
#
# 		for f_tpl in files_to_copy_lst:
# 			src_file_str, dest_dir_src = f_tpl
#
# 			gf_u.run_cmd('cp %s %s'%(src_file_str, dest_dir_src))
# 	#-----------------