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

import subprocess

sys.path.append("%s/../../meta"%(modd_str))
import gf_meta

sys.path.append("%s/../utils"%(modd_str))
import gf_core_cli

#--------------------------------------------------
def start(p_aws_creds_map,
	p_docker_sudo_bool = False):
	assert isinstance(p_aws_creds_map, dict)
	assert "GF_AWS_ACCESS_KEY_ID"     in p_aws_creds_map.keys()
	assert "GF_AWS_SECRET_ACCESS_KEY" in p_aws_creds_map.keys()


	# META
	meta_map = gf_meta.get()
	assert "local_cluster_config_dir_path_str" in meta_map.keys()
	local_cluster_config_dir_path_str = meta_map["local_cluster_config_dir_path_str"]
	assert os.path.isdir(local_cluster_config_dir_path_str)

	# local_cluster used for dev has its own S3 bucket, independent from the public or testing buckets.
	images_s3_bucket_name_str = meta_map["aws_s3_map"]["images_s3_bucket_map"]["local_cluster"]

	# IMPORTANT!! - needed for Elasticsearch to bootup correctly. otherwise error is raised:
	#               "max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]"
	sudo_str = ""
	if p_docker_sudo_bool:
		sudo_str = "sudo"
	_, stderr_str, return_code_int = gf_core_cli.run("%s sysctl -w vm.max_map_count=262144"%(sudo_str))
	if not return_code_int == 0:
		print("GF - failed to increase kernel memory setting for Elasticsearch")
		print(stderr_str)
		exit(1)




	cmd_lst = []
	

	if p_docker_sudo_bool:
		cmd_lst.extend([
			"sudo",

			# IMPORTANT!! - "sudo" starts its own environment, separate from the callers environment.
			#               so in order to pass in AWS credentials to docker-compose we have to
			#               name vars we want to pass in to "sudo" with the "--preserve-env" sudo CLI arg.
			"--preserve-env=GF_AWS_ACCESS_KEY_ID",
			"--preserve-env=GF_AWS_SECRET_ACCESS_KEY",
			"--preserve-env=GF_IMAGES_S3_BUCKET_NAME"
		])

	cmd_lst.extend([
		"docker-compose", "up",
	])

	# env_map = {
	# 	"GF_AWS_ACCESS_KEY_ID":     p_aws_creds_map["GF_AWS_ACCESS_KEY_ID"],
	# 	"GF_AWS_SECRET_ACCESS_KEY": p_aws_creds_map["GF_AWS_SECRET_ACCESS_KEY"],
	# }
	os.environ["GF_AWS_ACCESS_KEY_ID"]     = p_aws_creds_map["GF_AWS_ACCESS_KEY_ID"]
	os.environ["GF_AWS_SECRET_ACCESS_KEY"] = p_aws_creds_map["GF_AWS_SECRET_ACCESS_KEY"]
	os.environ["GF_IMAGES_S3_BUCKET_NAME"] = images_s3_bucket_name_str

	cwd = os.getcwd()
	os.chdir(local_cluster_config_dir_path_str)

	c_str = " ".join(cmd_lst)
	print(c_str)

	# docker-compose is a long-running process, and we want the output of the containers
	# that it starts to be passed onto us before waiting for the processes to complete.
	r = subprocess.Popen(c_str, shell = True, stdout = subprocess.PIPE, bufsize = 1) # env = env_map)
	for l in r.stdout:
		print(l.strip())

	# _, _, return_code_int = gf_core_cli.run(" ".join(cmd_lst),
	# 	p_env_map = env_map)

	os.chdir(cwd)


	