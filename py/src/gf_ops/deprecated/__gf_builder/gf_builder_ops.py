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

sys.path.append("%s/../containers"%(modd_str))
import gf_os_docker
import gf_containers

sys.path.append("%s/../../meta"%(modd_str))
import gf_meta

#---------------------------------------------------

# DEPRECATED!! - not being used anymore
'''
def cont__publish(p_app_name_str,
	p_app_build_meta_map,
	p_gf_dockerhub_user_str,
	p_gf_dockerhub_pass_str,
	p_log_fun,
	p_git_commit_hash_str = None,
	p_docker_sudo_bool    = False):
	assert isinstance(p_gf_dockerhub_user_str, str)

	# PUBLISH
	gf_containers.publish(p_app_name_str,
		p_app_build_meta_map,
		p_gf_dockerhub_user_str,
		p_gf_dockerhub_pass_str,
		p_log_fun,
		p_git_commit_hash_str = p_git_commit_hash_str, 
		p_exit_on_fail_bool   = True,
		p_docker_sudo_bool    = p_docker_sudo_bool)
'''

#---------------------------------------------------
# DEPRECATED?? - is this still being used?

def cont__build(p_app_name_str,
	p_dockerhub_user_name_str,
	p_log_fun,
	p_docker_sudo_bool = False):
	assert isinstance(p_dockerhub_user_name_str, str)

	build_meta_map = gf_meta.get()["build_info_map"]
	assert p_app_name_str in build_meta_map.keys()

	gf_builder_meta_map            = build_meta_map[p_app_name_str]
	cont_image_name_str            = gf_builder_meta_map["cont_image_name_str"]
	cont_image_version_str         = gf_builder_meta_map["version_str"]
	cont_image_dockerfile_path_str = os.path.abspath(gf_builder_meta_map["dockerfile_path_str"])
	assert os.path.isfile(cont_image_dockerfile_path_str)

	# DOCKER_BUILD
	image_full_name_str = "%s/%s:%s"%(p_dockerhub_user_name_str,
		cont_image_name_str,
		cont_image_version_str)

	gf_os_docker.build_image([image_full_name_str],
		cont_image_dockerfile_path_str,
		p_log_fun,
		p_exit_on_fail_bool = True,
		p_docker_sudo_bool  = p_docker_sudo_bool)