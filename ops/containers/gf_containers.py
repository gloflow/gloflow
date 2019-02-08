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

import subprocess
from colored import fg, bg, attr

sys.path.append('%s/../../meta'%(cwd_str))
import gf_meta
import gf_web_meta

sys.path.append('%s/../utils'%(cwd_str))
import gf_cli_utils
#-------------------------------------------------------------
def build(p_app_name_str,
	p_log_fun,
	p_user_name_str = 'local'):
	p_log_fun('FUN_ENTER','gf_containers.build()')
	assert isinstance(p_app_name_str, basestring)
	
	#------------------
	#META
	build_meta_map = gf_meta.get()['build_info_map']
	web_meta_map   = gf_web_meta.get()

	if not build_meta_map.has_key(p_app_name_str):
		p_log_fun("ERROR","supplied app (%s) does not exist in gf_meta"%(p_app_name_str))
		return
	app_meta_map = build_meta_map[p_app_name_str]

	service_name_str     = app_meta_map['service_name_str']
	service_base_dir_str = app_meta_map['service_base_dir_str']
	assert os.path.isdir(service_base_dir_str)

	service_version_str  = app_meta_map['version_str']
	assert len(service_version_str.split(".")) == 4 #format x.x.x.x

	#------------------
	#COPY_FILES_TO_DIR
	if app_meta_map.has_key('copy_to_dir_lst'):
		copy_to_dir_lst = app_meta_map['copy_to_dir_lst']
		copy_files(copy_to_dir_lst)
	#------------------
	#COPY_WEB_FILES
	if web_meta_map.has_key(p_app_name_str):

		app_web_meta_map = web_meta_map[p_app_name_str]
		assert app_web_meta_map.has_key('pages_map')
		pages_map = app_web_meta_map['pages_map']

		prepare_web_static_dir(pages_map, service_base_dir_str, p_log_fun)
	#------------------

	build_docker_container(service_name_str,
		service_base_dir_str,
		service_version_str,
		p_user_name_str,
		p_log_fun)
#--------------------------------------------------
def copy_files(p_copy_to_dir_lst):
    assert isinstance(p_copy_to_dir_lst, list)

    print('')
    print('             COPY FILES')
    for src_f_str, target_dir_str in p_copy_to_dir_lst:
        if not os.path.isdir(target_dir_str): gf_cli_utils.run_cmd('mkdir -p %s'%(target_dir_str))
        gf_cli_utils.run_cmd('cp %s %s'%(src_f_str, target_dir_str))
#-------------------------------------------------------------
def prepare_web_static_dir(p_pages_map,
	p_service_base_dir_str,
	p_log_fun):
	p_log_fun('FUN_ENTER','gf_containers.prepare_web_static_dir()')
	assert isinstance(p_pages_map, dict)
	assert os.path.dirname(p_service_base_dir_str)

	for pg_name_str, pg_info_map in p_pages_map.items():

		assert pg_info_map.has_key('build_dir_str')
		assert os.path.isdir(pg_info_map['build_dir_str'])
		build_dir_str = pg_info_map['build_dir_str']

		#------------------
		#CREATE_TARGET_DIR
		target_dir_str = '%s/static'%(p_service_base_dir_str)
		gf_cli_utils.run_cmd('mkdir -p %s'%(target_dir_str))
		#------------------
		#COPY_PAGE_WEB_CODE
		gf_cli_utils.run_cmd('cp -p -r %s/* %s'%(build_dir_str, target_dir_str))
		#------------------
#-------------------------------------------------------------
def build_docker_container(p_service_name_str,
	p_service_base_dir_str,
	p_service_version_str,
	p_user_name_str,
	p_log_fun):
	p_log_fun('FUN_ENTER','gf_containers.build_docker_container()')
	assert os.path.isdir(p_service_base_dir_str)

	image_name_str       = '%s/%s:%s'%(p_user_name_str, p_service_name_str, p_service_version_str)
	dockerfile_path_str  = '%s/Dockerfile'%(p_service_base_dir_str)
	context_dir_path_str = p_service_base_dir_str

	p_log_fun('INFO','====================+++++++++++++++=====================')
	p_log_fun('INFO','                 BUILDING PACKAGE/SERVICE IMAGE')
	p_log_fun('INFO','              %s'%(p_service_name_str))
	p_log_fun('INFO','Dockerfile     - %s'%(dockerfile_path_str))
	p_log_fun('INFO','image_name_str - %s'%(image_name_str))
	p_log_fun('INFO','====================+++++++++++++++=====================')

	cmd_lst = [
		'sudo docker build',
		'-f %s'%(dockerfile_path_str),
		'--tag=%s'%(image_name_str),
		context_dir_path_str
	]

	cmd_str = ' '.join(cmd_lst)
	p_log_fun('INFO',' - %s'%(cmd_str))

	#change to the dir where the Dockerfile is located, for the 'docker'
	#tool to have the proper context
	old_cwd = os.getcwd()
	os.chdir(context_dir_path_str)
	
	r = subprocess.Popen(cmd_str, shell = True, stdout = subprocess.PIPE, bufsize = 1)

	#---------------------------------------------------
	def get_image_id_from_line(p_stdout_line_str):
		p_lst = p_stdout_line_str.split(' ')

		assert len(p_lst) == 3
		image_id_str = p_lst[2]

		#IMPORTANT!! - check that this is a valid 12 char Docker ID
		assert len(image_id_str) == 12
		return image_id_str
	#---------------------------------------------------

	for line in r.stdout:
		line_str = line.strip() #strip() - to remove '\n' at the end of the line

		#------------------
		#display the line, to update terminal display
		print(line_str)
		#------------------

		if line_str.startswith('Successfully built'):
			image_id_str = get_image_id_from_line(line_str)
			return image_id_str

	#change back to old dir
	os.chdir(old_cwd)