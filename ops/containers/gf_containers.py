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
cwd_str = os.path.abspath(os.path.dirname(__file__))

import subprocess
from colored import fg, bg, attr

sys.path.append('%s/../utils'%(cwd_str))
import gf_cli_utils

import gf_os_docker

#-------------------------------------------------------------
# PULL
def pull(p_image__full_name_str,
	p_log_fun,
	p_dockerhub_user_str = None,
	p_dockerhub_pass_str = None,
	p_exit_on_fail_bool  = True,
	p_docker_sudo_bool   = False):

	# often times public containers are being pulled, so no login is needed for that and 
	# callers dont submit their credentials.
	if not p_dockerhub_pass_str == None and not p_dockerhub_pass_str == "":
		gf_os_docker.login(p_dockerhub_user_str,
			p_dockerhub_pass_str,
			p_exit_on_fail_bool = True,
			p_docker_sudo_bool  = p_docker_sudo_bool)

	# DOCKER_PULL
	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
	
	cmd_lst.extend([
		"docker pull",
		p_image__full_name_str
	])
	c_pull = " ".join(cmd_lst)
	p_log_fun("INFO", "cmd - %s"%(c_pull))

	p = subprocess.Popen(c_pull,
		shell   = True,
		stdout  = subprocess.PIPE,
		bufsize = 1)

	for line in p.stdout:
		clean_line_str = line.strip()
		print(clean_line_str)

	if p_exit_on_fail_bool:
		if not p.returncode == None and not p.returncode == 0:
			exit()

#-------------------------------------------------------------
# BUILD
def build(p_app_name_str,
	p_app_build_meta_map,
	p_app_web_meta_map,
	p_log_fun,
	p_user_name_str     = "local",
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_containers.build()")
	p_log_fun("INFO",      "p_app_name_str - %s"%(p_app_name_str))
	assert isinstance(p_app_name_str,       basestring)
	assert isinstance(p_app_build_meta_map, dict)
	assert isinstance(p_app_web_meta_map,   dict)

	#------------------
	# META

	assert "service_name_str" in p_app_build_meta_map
	assert "service_base_dir_str" in p_app_build_meta_map

	service_name_str     = p_app_build_meta_map["service_name_str"]
	service_base_dir_str = p_app_build_meta_map["service_base_dir_str"]
	assert os.path.isdir(service_base_dir_str)

	# service_dockerfile_path_str = "%s/Dockerfile"%(service_base_dir_str)
	service_dockerfile_path_str = get_service_dockerfile(p_app_build_meta_map)
	service_version_str         = p_app_build_meta_map["version_str"]
	# assert len(service_version_str.split(".")) == 4 # format x.x.x.x
	#------------------
	# COPY_FILES_TO_DIR
	if p_app_build_meta_map.has_key("copy_to_dir_lst"):
		copy_to_dir_lst = p_app_build_meta_map["copy_to_dir_lst"]
		copy_files(copy_to_dir_lst)
	#------------------
	# PREPARE_WEB_FILES
	assert p_app_web_meta_map.has_key("pages_map")
	pages_map = p_app_web_meta_map["pages_map"]

	prepare_web_files(pages_map,
		service_base_dir_str,
		p_log_fun)
	#------------------

	image_name_str = service_name_str
	image_tag_str  = service_version_str


	build_docker_image(image_name_str,
		image_tag_str,
		service_dockerfile_path_str,
		p_user_name_str,
		p_log_fun,
		p_exit_on_fail_bool = p_exit_on_fail_bool,
		p_docker_sudo_bool  = p_docker_sudo_bool)

#--------------------------------------------------
def copy_files(p_copy_to_dir_lst):
    assert isinstance(p_copy_to_dir_lst, list)

    print('')
    print('             COPY FILES')
    for src_f_str, target_dir_str in p_copy_to_dir_lst:
        if not os.path.isdir(target_dir_str):
			gf_cli_utils.run_cmd('mkdir -p %s'%(target_dir_str))
        gf_cli_utils.run_cmd('cp %s %s'%(src_f_str, target_dir_str))

#-------------------------------------------------------------
# PREPARE_WEB_FILES
def prepare_web_files(p_pages_map,
	p_service_base_dir_str,
	p_log_fun):
	p_log_fun('FUN_ENTER', 'gf_containers.prepare_web_files()')
	assert isinstance(p_pages_map, dict)
	assert os.path.dirname(p_service_base_dir_str)

	for pg_name_str, pg_info_map in p_pages_map.items():
		assert isinstance(pg_info_map, dict)
		assert pg_info_map.has_key('build_dir_str')
		assert os.path.isdir(pg_info_map['build_dir_str'])

		build_dir_str = os.path.abspath(pg_info_map['build_dir_str'])

		#------------------
		# CREATE_TARGET_DIR
		target_dir_str = os.path.abspath('%s/static'%(p_service_base_dir_str))
		gf_cli_utils.run_cmd('mkdir -p %s'%(target_dir_str))
		#------------------
		# COPY_PAGE_WEB_CODE
		gf_cli_utils.run_cmd('cp -r %s/* %s'%(build_dir_str, target_dir_str))
		#------------------

	#------------------
	# MOVE_TEMPLATES_OUT_OF_STATIC

	# IMPORTANT!! - templates should not be in the static/ dir, which would make them servable
	#               over HTTP which we dont want. instead its moved out of the static/ dir 
	#               to its parent dir where its private.
	#               templates are originally in the static/ dir because durring the build process they were
	#               handled together with other static content (html/css/js files) and as output moved
	#               into that static/ dir from other locations while in development.
	gf_cli_utils.run_cmd('rm -rf %s/../templates'%(target_dir_str)) # remove existing templates build dir
	gf_cli_utils.run_cmd('mv %s/templates %s/..'%(target_dir_str, target_dir_str))
	#------------------

#-------------------------------------------------------------
# PUBLISH
def publish(p_app_name_str,
	p_app_build_meta_map,
	p_dockerhub_user_str,
	p_dockerhub_pass_str,
	p_log_fun,
	p_exit_on_fail_bool = False):
	p_log_fun('FUN_ENTER', 'gf_containers.publish()')
	p_log_fun('INFO',      'p_app_name_str - %s'%(p_app_name_str))
	assert isinstance(p_app_build_meta_map, dict)

	service_name_str    = p_app_build_meta_map['service_name_str']
	service_version_str = p_app_build_meta_map['version_str']

	publish_docker_image(service_name_str,
		service_version_str,
		p_dockerhub_user_str,
		p_dockerhub_pass_str,
		p_log_fun,
		p_exit_on_fail_bool = p_exit_on_fail_bool)

#-------------------------------------------------------------
# PUBLISH_DOCKER_IMAGE
def publish_docker_image(p_image_name_str,
	p_image_tag_str,
	p_dockerhub_user_str,
	p_dockerhub_pass_str,
	p_log_fun,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):	

	#------------------
	# DOCKERHUB_LOGIN
	# cmd_lst = [
	# 	"docker", "login",
	# 	"-u", p_dockerhub_user_str,
	# 	"--password-stdin"
	# ]
	#
	# process = subprocess.Popen(cmd_lst, stdin = subprocess.PIPE, stdout = subprocess.PIPE)
	# process.stdin.write(p_dockerhub_pass_str) # write password on stdin of "docker login" command
	# stdout_str, stderr_str = process.communicate() # wait for command completion
	# print(stdout_str)
	# print(stderr_str)

	gf_os_docker.login(p_dockerhub_user_str,
		p_dockerhub_pass_str,
		p_exit_on_fail_bool = p_exit_on_fail_bool,
		p_docker_sudo_bool  = p_docker_sudo_bool)
	#------------------
	# DOCKER_PUSH

	full_image_name_str = '%s/%s:%s'%(p_dockerhub_user_str, p_image_name_str, p_image_tag_str)
	cmd_lst = [
		"docker push",
		full_image_name_str
	]

	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", " - %s"%(c_str))

	stdout_str, _, exit_code_int = gf_cli_utils.run_cmd(c_str)
	

	# IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)
	#------------------
	# DOCKER_LOGOUT
	stdout_str, _, _ = gf_cli_utils.run_cmd("docker logout")
	print(stdout_str)
	#------------------

#-------------------------------------------------------------
# BUILD_DOCKER_IMAGE
def build_docker_image(p_image_name_str,
	p_image_tag_str,
	p_dockerfile_path_str,
	p_user_name_str,
	p_log_fun,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_containers.build_docker_image()")
	assert os.path.isfile(p_dockerfile_path_str)
	assert "Dockerfile" in os.path.basename(p_dockerfile_path_str)

	full_image_name_str  = "%s/%s:%s"%(p_user_name_str, p_image_name_str, p_image_tag_str)
	context_dir_path_str = os.path.dirname(p_dockerfile_path_str)

	p_log_fun("INFO", "====================+++++++++++++++=====================")
	p_log_fun("INFO", "                 BUILDING DOCKER IMAGE")
	p_log_fun("INFO", "              %s"%(p_image_name_str))
	p_log_fun("INFO", "Dockerfile          - %s"%(p_dockerfile_path_str))
	p_log_fun("INFO", "full_image_name_str - %s"%(full_image_name_str))
	p_log_fun("INFO", "====================+++++++++++++++=====================")

	cmd_lst = []

	# RUN_WITH_SUDO - Docker deamon runs as root, and so for docker client to be able to connect to it
	#                 without any custom config the client needs to be run with "sudo".
	#                 if some config is in place to avoid this, set p_docker_sudo_bool to False.
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker build",
		"-f %s"%(p_dockerfile_path_str),
		"--tag=%s"%(full_image_name_str),
		context_dir_path_str
	])

	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", " - %s"%(c_str))

	# change to the dir where the Dockerfile is located, for the 'docker'
	# tool to have the proper context
	old_cwd = os.getcwd()
	os.chdir(context_dir_path_str)
	
	#---------------------------------------------------
	def get_image_id_from_line(p_stdout_line_str):
		p_lst = p_stdout_line_str.split(' ')

		assert len(p_lst) == 3
		image_id_str = p_lst[2]

		# IMPORTANT!! - check that this is a valid 12 char Docker ID
		assert len(image_id_str) == 12
		return image_id_str

	#---------------------------------------------------

	stdout_str, _, exit_code_int = gf_cli_utils.run_cmd(c_str)

	for line_str in stdout_str:
		if line_str.startswith("Successfully built"):
			image_id_str = get_image_id_from_line(line_str)
			return image_id_str

	# IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	# change back to old dir
	os.chdir(old_cwd)

#-------------------------------------------------------------
def get_service_dockerfile(p_app_build_meta_map):
	service_base_dir_str = p_app_build_meta_map["service_base_dir_str"]
	assert os.path.isdir(service_base_dir_str)

	service_dockerfile_path_str = "%s/Dockerfile"%(service_base_dir_str)
	return service_dockerfile_path_str