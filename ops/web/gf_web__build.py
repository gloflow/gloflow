# GloFlow media management/publishing system
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

sys.path.append('%s/../../meta'%(cwd_str))
import gf_web_meta

sys.path.append('%s/../utils'%(cwd_str))
import gf_cli_utils as gf_u
#---------------------------------------------------
def build(p_apps_names_lst, p_log_fun):
	p_log_fun("FUN_ENTER","gf_web__build.build()")
	assert isinstance(p_apps_names_lst, list)
	assert len(p_apps_names_lst) > 0

	apps_meta_map = gf_web_meta.get()

	for app_str in p_apps_names_lst:
		
		#-----------------
		#META
		if not apps_meta_map.has_key(app_str):
			p_log_fun("ERROR","supplied app (%s) does not exist in gf_web_meta"%(app_str))
			return
		app_map = apps_meta_map[app_str]
		#-----------------

		#BUILD PAGES - build each page of the app
		for page_name_str, page_info_map in app_map['pages_map'].items():

			build_dir_str = page_info_map['build_dir_str']

			build_page(page_name_str,
				build_dir_str,
				page_info_map,
				p_log_fun)
#---------------------------------------------------
def build_page(p_page_name_str,
    p_build_dir_str,
    p_page_info_map,
    p_log_fun):
	p_log_fun("FUN_ENTER", "gf_web__build.build_page()")
	assert os.path.isdir(p_build_dir_str)

	print('')
	print('')
	p_log_fun('INFO', '%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s'%(fg('orange_red_1'), attr(0)))
	p_log_fun('INFO', '             %sBUILD PAGE%s - %s%s%s'%(fg('cyan'), attr(0), fg('orange_red_1'), p_page_name_str, attr(0)))
	p_log_fun('INFO', 'page type - %s'%(p_page_info_map['type_str']))
	p_log_fun('INFO', '%s>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%s'%(fg('orange_red_1'), attr(0)))
	p_log_fun('INFO', 'build_dir_str - %s'%(p_build_dir_str))
	
	assert isinstance(p_build_dir_str, basestring)
	assert os.path.isdir(p_build_dir_str)

	#---------------------------------------------------
	def build_typescript(p_out_file_str, p_minified_file_str, p_ts_files_lst):
		p_log_fun("FUN_ENTER", "gf_web__build.build_page().build_typescript()")

		cmd_lst = [
			'tsc',
			'--out %s'%(p_out_file_str),
			' '.join(p_ts_files_lst)
		]
		gf_u.run_cmd(' '.join(cmd_lst))

		#minify into the same file name as the Typescript compiler output
		minify_js(p_minified_file_str, [p_out_file_str], p_log_fun)
	#---------------------------------------------------

	#-----------------
	if p_page_info_map['type_str'] == 'ts':
		p_log_fun('INFO','%s------------ TYPESCRIPT --------------------------------%s'%(fg('yellow'),attr(0)))

		out_file_str      = p_page_info_map['ts']['out_file_str']
		minified_file_str = p_page_info_map['ts']['minified_file_str']
		ts_files_lst      = p_page_info_map['ts']['files_lst']

		build_typescript(out_file_str, minified_file_str, ts_files_lst)
		#-----------------
		#COPY LIBS

		if p_page_info_map['ts'].has_key('libs_files_lst'):
			p_log_fun('INFO','%s------------ TS_LIBS -----------------------------%s'%(fg('yellow'),attr(0)))

			if not os.path.isdir(p_build_dir_str): gf_u.run_cmd('mkdir -p %s/js/lib'%(p_build_dir_str))
			
			for lib_file_str in p_page_info_map['ts']['libs_files_lst']:
				gf_u.run_cmd('cp %s %s/js/lib'%(lib_file_str, p_build_dir_str))
		#-----------------
	#-----------------
	#if p_page_info_map['type_str'] == 'js':
	#	p_log_fun('INFO','%s------------ JAVASCRIPT --------------------------------%s'%(fg('yellow'),attr(0))))
	#	#-----------------
	#	#MINIFY JS
	#	if p_page_info_map['js'].has_key('minified_file_str'):
	#		p_log_fun('INFO','------------ JS_MINIFY ---------------------------')
	#		minified_file_str = p_page_info_map['js']['minified_file_str']
	#		js_files_lst      = p_page_info_map['js']['files_lst']
	#
	#		minify_js(minified_file_str, js_files_lst, p_log_fun)
	#	#-----------------
	#	#COPY JS LIBS
	#
	#	if p_page_info_map['js'].has_key('libs_files_lst'):
	#		p_log_fun('INFO','------------ JS_LIBS -----------------------------')
	#
	#		for lib_file_str in p_page_info_map['js']['libs_files_lst']:
	#			gf_u.run_cmd('cp %s %s/js/lib'%(lib_file_str, p_build_dir_str))
	#-----------------
	#CSS

	if p_page_info_map.has_key('css'):
		p_log_fun('INFO','%s------------ CSS ---------------------------------%s'%(fg('yellow'),attr(0)))
		css_files_lst = p_page_info_map['css']['files_lst']

		for f_tpl in css_files_lst:
			assert len(f_tpl) == 2

			src_file_str, dest_dir_src = f_tpl

			if not os.path.isdir(dest_dir_src): gf_u.run_cmd('mkdir -p %s'%(dest_dir_src))
			gf_u.run_cmd('cp %s %s'%(src_file_str, dest_dir_src))
	#-----------------
	#TEMPLATES
	if p_page_info_map.has_key('templates'):
		p_log_fun('INFO','%s------------ TEMPLATES -----------------------------%s'%(fg('yellow'),attr(0)))
		
		assert p_page_info_map['templates'].has_key('files_lst')
		templates_files_lst = p_page_info_map['templates']['files_lst']
		assert isinstance(templates_files_lst, list)


		#COPY_TEMPLATE_FILES - copy them from their source location to the desired build location
		for tmpl_file_str, tmpl_target_dir_str in templates_files_lst:
			print('tmpl_file_str - %s'%(tmpl_file_str))

			assert os.path.isfile(tmpl_file_str)

			#if target template dir doesnt exist, create it
			if not os.path.isdir(tmpl_target_dir_str): gf_u.run_cmd('mkdir -p %s'%(tmpl_target_dir_str))

			gf_u.run_cmd('cp %s %s'%(tmpl_file_str, tmpl_target_dir_str))
	#-----------------
	#FILES_TO_COPY

	if p_page_info_map.has_key('files_to_copy_lst'):
		p_log_fun('INFO','%s------------ FILES_TO_COPY -----------------------%s'%(fg('yellow'),attr(0)))
		files_to_copy_lst = p_page_info_map['files_to_copy_lst']

		for f_tpl in files_to_copy_lst:
			src_file_str, dest_dir_src = f_tpl

			gf_u.run_cmd('cp %s %s'%(src_file_str, dest_dir_src))
	#-----------------
#---------------------------------------------------
def minify_js(p_js_target_file_str,
    p_js_files_lst,
    p_log_fun):
	p_log_fun("FUN_ENTER","gf_apps_build.minify_js()")

	cmd_lst = [
		'uglifyjs',
		'--output %s'%(p_js_target_file_str),
		' '.join(p_js_files_lst),
	]
	gf_u.run_cmd(' '.join(cmd_lst))