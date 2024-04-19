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
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import subprocess
from colored import fg, bg, attr

sys.path.append("%s/../../gf_core"%(modd_str))
import gf_core_cli

import gf_os_docker

#-------------------------------------------------------------
# PULL
def pull(p_image__full_name_str,
	p_log_fun,
	p_docker_user_str   = None,
	p_docker_pass_str   = None,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = False):

	# often times public containers are being pulled, so no login is needed for that and 
	# callers dont submit their credentials.
	if not p_docker_pass_str == None and not p_docker_pass_str == "":
		gf_os_docker.login(p_docker_user_str,
			p_docker_pass_str,
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
	p_log_fun,
	p_app_web_meta_map    = None,
	p_user_name_str       = "local",
	p_git_commit_hash_str = None,
	p_exit_on_fail_bool   = False,
	p_docker_sudo_bool    = False):
	p_log_fun("INFO", f"p_app_name_str - {p_app_name_str}")
	assert isinstance(p_app_name_str,       str)
	assert isinstance(p_app_build_meta_map, dict)

	#------------------
	# META

	is_service_bool = False
	if "service_name_str" in p_app_build_meta_map:
		assert "service_base_dir_str" in p_app_build_meta_map.keys()
		is_service_bool = True

		
	if is_service_bool:
		service_name_str     = p_app_build_meta_map["service_name_str"]
		service_base_dir_str = p_app_build_meta_map["service_base_dir_str"]
		assert os.path.isdir(service_base_dir_str)

	# service_dockerfile_path_str = "%s/Dockerfile"%(service_base_dir_str)
	service_dockerfile_path_str = get_service_dockerfile(p_app_build_meta_map)
	
	#------------------
	# COPY_FILES_TO_DIR
	if "copy_to_dir_lst" in p_app_build_meta_map.keys():

		copy_to_dir_lst = p_app_build_meta_map["copy_to_dir_lst"]
		
		copy_files(copy_to_dir_lst)

	#------------------
	# PREPARE_WEB_FILES
	if not p_app_web_meta_map == None:
		assert isinstance(p_app_web_meta_map, dict)
		assert "pages_map" in p_app_web_meta_map.keys()
		pages_map = p_app_web_meta_map["pages_map"]

		prepare_web_files(pages_map,
			service_base_dir_str,
			p_log_fun,
			p_docker_sudo_bool = p_docker_sudo_bool)

	#------------------
	# IMAGE_FULL_NAMES
		
	if is_service_bool:
		image_name_str = service_name_str
	else:
		assert "cont_image_name_str" in p_app_build_meta_map.keys()
		image_name_str = p_app_build_meta_map["cont_image_name_str"]
		
	image_full_names_lst = get_image_full_names(image_name_str,
		p_app_build_meta_map,
		p_user_name_str,
		p_git_commit_hash_str = p_git_commit_hash_str)

	#------------------
	# BUILD_ARGS
	build_args_map = {}

	# BASE_IMAGE_TAG - tag of the base image from which the main image thats being
	#                  built from is inheriting from.
	if not p_git_commit_hash_str == None:
		build_args_map["GF_BASE_IMAGE_TAG"] = p_git_commit_hash_str
	else:
		build_args_map["GF_BASE_IMAGE_TAG"] = "latest"
	
	#------------------
	# DOCKER_BUILD
	gf_os_docker.build_image(image_full_names_lst,
		service_dockerfile_path_str,
		p_log_fun,
		p_build_args_map    = build_args_map,
		p_exit_on_fail_bool = p_exit_on_fail_bool,
		p_docker_sudo_bool  = p_docker_sudo_bool)

	#------------------

#-------------------------------------------------------------
# PUBLISH
def publish(p_app_name_str,
	p_app_build_meta_map,
	p_docker_user_str,
	p_docker_pass_str,
	p_log_fun,
	p_git_commit_hash_str = None,
	p_exit_on_fail_bool   = False,
	p_docker_sudo_bool    = False):
	p_log_fun("INFO", "p_app_name_str - %s"%(p_app_name_str))
	assert isinstance(p_app_build_meta_map, dict)

	if "service_name_str" in p_app_build_meta_map.keys():
		service_name_str = p_app_build_meta_map["service_name_str"]
		image_name_str   = service_name_str
	else:
		image_name_str = p_app_name_str

	# service_version_str = p_app_build_meta_map["version_str"]
	#
	# image_tag_str = None
	# if not p_git_commit_hash_str == None:
	# 	image_tag_str = p_git_commit_hash_str
	# else:
	# 	service_version_str = p_app_build_meta_map["version_str"]
	# 	image_tag_str       = service_version_str

	image_full_names_lst = get_image_full_names(image_name_str,
		p_app_build_meta_map,
		p_docker_user_str,
		p_git_commit_hash_str = p_git_commit_hash_str)
	assert isinstance(image_full_names_lst, list)
	
	for image_full_name_str in image_full_names_lst:
		
		# DOCKER_PUSH
		gf_os_docker.push(image_full_name_str,
			p_docker_user_str,
			p_docker_pass_str,
			p_log_fun,
			p_exit_on_fail_bool = p_exit_on_fail_bool,
			p_docker_sudo_bool  = p_docker_sudo_bool)

#-------------------------------------------------------------
# GET_IMAGE_FULL_NAMES
def get_image_full_names(p_image_name_str,
	p_app_build_meta_map,
	p_user_name_str,
	p_git_commit_hash_str = None):
	assert isinstance(p_image_name_str, str)
	assert isinstance(p_app_build_meta_map, dict)
	assert isinstance(p_user_name_str, str)

	# IMAGE_TAG
	image_tag_str = None
	if not p_git_commit_hash_str == None:
		# if a git commit hash was supplied, tag the image with that
		image_tag_str = p_git_commit_hash_str
	else:
		service_version_str = p_app_build_meta_map["version_str"]

		# assert len(service_version_str.split(".")) == 4 # format x.x.x.x
		image_tag_str = service_version_str

	image_full_names_lst = []

	# standard name
	image_full_name_str = "%s/%s:%s"%(p_user_name_str, p_image_name_str, image_tag_str)
	image_full_names_lst.append(image_full_name_str)

	# IMPORTANT!! - "latest" name - its important to always havea a "latest" image that points
	#               to the most up-to-date container image for use in situations when we dont know
	#               the version number or git commit hash or some other tag.
	if not image_tag_str == "latest":
		image_full_name_latest_str = "%s/%s:latest"%(p_user_name_str, p_image_name_str)
		image_full_names_lst.append(image_full_name_latest_str)

	return image_full_names_lst

#-------------------------------------------------------------
def copy_files(p_copy_to_dir_lst):
	assert isinstance(p_copy_to_dir_lst, list)

	print("")
	print("             COPY FILES")
	for src, target_dir_str in p_copy_to_dir_lst:
		if not os.path.isdir(target_dir_str):
			gf_core_cli.run("mkdir -p %s"%(target_dir_str))

		# COPY_DIR
		if os.path.isdir(src):
			src_dir_str = src
			gf_core_cli.run("cp -r %s %s"%(src_dir_str, target_dir_str))

		# COPY_FILE
		else:
			src_file_str = src
			gf_core_cli.run("cp %s %s"%(src_file_str, target_dir_str))

#-------------------------------------------------------------
# PREPARE_WEB_FILES
def prepare_web_files(p_pages_map,
	p_service_base_dir_str,
	p_log_fun,
	p_docker_sudo_bool = False):
	p_log_fun("FUN_ENTER", "gf_containers.prepare_web_files()")
	assert isinstance(p_pages_map, dict)
	assert os.path.dirname(p_service_base_dir_str)

	for pg_name_str, pg_info_map in p_pages_map.items():
		print(f"======== {fg('green')}{'%s'%(pg_name_str)}{attr(0)}")
		assert isinstance(pg_info_map, dict)
		assert "build_dir_str" in pg_info_map.keys()
		assert os.path.isdir(pg_info_map["build_dir_str"])

		build_dir_str = os.path.abspath(pg_info_map["build_dir_str"])

		#------------------
		# CREATE_TARGET_DIR
		target_dir_str = os.path.abspath(f"{p_service_base_dir_str}/static")
		gf_core_cli.run(f"mkdir -p {target_dir_str}")

		#------------------
		# COPY_PAGE_WEB_CODE
		gf_core_cli.run(f"cp -r {build_dir_str}/* {target_dir_str}")

		#------------------
		
	#------------------
	# MOVE_TEMPLATES_OUT_OF_STATIC

	# IMPORTANT!! - templates should not be in the static/ dir, which would make them servable
	#               over HTTP which we dont want. instead its moved out of the static/ dir 
	#               to its parent dir where its private.
	#               templates are originally in the static/ dir because durring the build process they were
	#               handled together with other static content (html/css/js files) and as output moved
	#               into that static/ dir from other locations while in development.
	gf_core_cli.run("rm -rf %s/../templates"%(target_dir_str)) # remove existing templates build dir
	gf_core_cli.run("mv %s/templates %s/.."%(target_dir_str, target_dir_str))
	
	#------------------

#-------------------------------------------------------------
def get_service_dockerfile(p_app_build_meta_map):

	if "dockerfile_path_str" in p_app_build_meta_map.keys():
		service_dockerfile_path_str = p_app_build_meta_map["dockerfile_path_str"]
	else:

		service_base_dir_str = p_app_build_meta_map["service_base_dir_str"]
		assert os.path.isdir(service_base_dir_str)

		service_dockerfile_path_str = "%s/Dockerfile"%(service_base_dir_str)

	return service_dockerfile_path_str