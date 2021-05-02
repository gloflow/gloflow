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

import json
import subprocess
import base64

import delegator

sys.path.append("%s/../utils"%(cwd_str))
import gf_cli_utils

#---------------------------------------------------
def cont_is_running(p_cont_name_str,
	p_log_fun,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = True):

	sudo_str = ""
	if p_docker_sudo_bool:
		sudo_str = "sudo"

	stdout_str, stderr_str, exit_code_int = gf_cli_utils.run_cmd("%s docker ps -a | grep %s"%(sudo_str, p_cont_name_str))

	if not stderr_str == "":
		print(stderr_str)
		
	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in stderr_str:
		exit(1)

	if stdout_str == "":
		print("CONTAINER NOT RUNNING -----------------------")
		return False
	else:
		print("CONTAINER RUNNING -----------------------")
		return True

#---------------------------------------------------
def cont_is_running_remote(p_cont_name_str,
	p_fab_conn,
	p_log_fun,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = True):

	print("CHECK IF CONTAINER IS RUNNING - %s"%(p_cont_name_str))
	sudo_str = ""
	if p_docker_sudo_bool:
		sudo_str = "sudo"

	c_str = "%s docker ps -a | grep %s"%(sudo_str, p_cont_name_str)
	print(c_str)
	r = p_fab_conn.run(c_str, warn=True)
	if not r.stdout == "": print(r.stdout)
	if not r.stderr == "": print(r.stderr)



	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in r.stderr:
		print("cannot connect to Docker daemon")
		exit(1)
	
	# # IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    # #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	# #               in other scenarious its acceptable for this command to fail, and we want the caller
	# #               to keep executing.
	# if not r.return_code == 0:
	# 	if p_exit_on_fail_bool:
	# 		exit(r.return_code)

	if r.stdout == "":
		print("CONTAINER NOT RUNNING")
		return False
	else:
		print("CONTAINER RUNNING")
		return True

#-------------------------------------------------------------
# RUN
def run(p_full_image_name_str,
	p_log_fun,
	p_container_name_str = None,
	p_ports_map          = None,
	p_volumes_map        = None,
	p_hostname_str       = None,
	p_host_network_bool  = False,
	p_detached_bool      = True,
	p_exit_on_fail_bool  = False,
	p_docker_sudo_bool   = True):
	assert isinstance(p_full_image_name_str, str)

	print("")
	print("RUNNING DOCKER CONTAINER - %s"%(p_full_image_name_str))

	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"docker run",
		"--restart=always",
	])

	# CONTAINER_NAME
	if not p_container_name_str == None:
		cmd_lst.append("--name %s"%(p_container_name_str))

	# PORTS
	if not p_ports_map == None:
		for host_port_str, container_port_str in p_ports_map.items():
			# IMPORTANT!! - "-p" publish a container's port or a range of ports to the host.
			cmd_lst.append("-p %s:%s"%(host_port_str, container_port_str))

	# VOLUMES
	if not p_volumes_map == None:
		for host_dir_str, container_dir_str in p_volumes_map.items():
			# IMPORTANT!! - "-v" - mount a host directory into a particular directory path in the
			#                      container filesystem.
			cmd_lst.append("-v %s:%s"%(host_dir_str, container_dir_str))

	# HOSTNAME
	if not p_hostname_str == None:
		cmd_lst.append("-h %s"%(p_hostname_str))

	# HOST_NETWORK
	if p_host_network_bool:
		cmd_lst.append("--net=host")

	# DETACHED
	if p_detached_bool:
		cmd_lst.append("-d")


	# IMAGE_NAME
	cmd_lst.append(p_full_image_name_str)


	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", " - %s"%(c_str))

	stdout_str, stderr_str, exit_code_int = gf_cli_utils.run_cmd(c_str)

	if not stderr_str == "":
		print(stderr_str)

	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in stderr_str:
		exit(1)

	# IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	#               in other scenarious its acceptable for this command to fail, and we want the caller
	#               to keep executing.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	# CONTAINER_ID
	container_id_str = stdout_str.strip()
	assert len(container_id_str) == 64

	return container_id_str

#-------------------------------------------------------------
# REMOVE
def remove_by_name(p_container_name_str,
	p_log_fun,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = True):

	sudo_str = ""
	if p_docker_sudo_bool:
		sudo_str = "sudo"

	cmd_str = "%s docker rm -f `%s docker ps -a | grep %s | awk '{print $1}'`"%(sudo_str, sudo_str, p_container_name_str)
	stdout_str, stderr_str, exit_code_int = gf_cli_utils.run_cmd(cmd_str)

	if not stderr_str == "":
		print(stderr_str)
		
	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in stderr_str:
		exit(1)

	# IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	#               in other scenarious its acceptable for this command to fail, and we want the caller
	#               to keep executing.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

#-------------------------------------------------------------
def remove_by_name_remote(p_container_name_str,
	p_fab_conn,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = True):

	sudo_str = ""
	if p_docker_sudo_bool:
		sudo_str = "sudo"

	cmd_str       = "%s docker rm -f `%s docker ps -a | grep %s | awk '{print $1}'`"%(sudo_str, sudo_str, p_container_name_str)
	r             = p_fab_conn.run(cmd_str)
	exit_code_int = r.return_code

	stdout_and_stderr_str = r.stdout+r.stderr
	
	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in r.stdout:
		exit(1)

	# IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	#               in other scenarious its acceptable for this command to fail, and we want the caller
	#               to keep executing.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

#-------------------------------------------------------------
# PULL_IMAGE
def pull_remote(p_cont_image_name_str,
	p_fab_conn,
	p_log_fun,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_os_docker.pull_image()")

	sudo_str = ""
	if p_docker_sudo_bool:
		sudo_str = "sudo"

	c_str = "%s docker pull %s"%(sudo_str, p_cont_image_name_str)
	p_fab_conn.run(c_str)

#-------------------------------------------------------------
# PUSH
def push(p_image_full_name_str,
	p_docker_user_str,
	p_docker_pass_str,
	p_log_fun,
	p_host_str          = None,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_os_docker.push()")
	p_log_fun("INFO",      f"image_full_name - {p_image_full_name_str}")
	assert isinstance(p_docker_user_str, str)

	#------------------
	# LOGIN
	login(p_docker_user_str,
		p_docker_pass_str,
		p_host_str          = p_host_str,
		p_exit_on_fail_bool = p_exit_on_fail_bool,
		p_docker_sudo_bool  = p_docker_sudo_bool)

	#------------------
	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"docker push",
		p_image_full_name_str
	])

	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", " - %s"%(c_str))

	stdout_str, stderr_str, exit_code_int = gf_cli_utils.run_cmd(c_str)

	if not stderr_str == "":
		print(stderr_str)
		
	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in stderr_str:
		exit(1)

	# IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	#               in other scenarious its acceptable for this command to fail, and we want the caller
	#               to keep executing.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	#------------------
	# DOCKER_LOGOUT
	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
	cmd_lst.append("docker logout")
	stdout_str, _, _ = gf_cli_utils.run_cmd(" ".join(cmd_lst))
	print(stdout_str)

	#------------------

#-------------------------------------------------------------
# BUILD_IMAGE
def build_image(p_image_names_lst,
	p_dockerfile_path_str,
	p_log_fun,
	p_build_args_map    = None,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_os_docker.build_image()")
	assert isinstance(p_image_names_lst, list)
	print(p_dockerfile_path_str)
	assert os.path.isfile(p_dockerfile_path_str)
	assert "Dockerfile" in os.path.basename(p_dockerfile_path_str)

	# full_image_name_str  = "%s/%s:%s"%(p_user_name_str, p_image_name_str, p_image_tag_str)
	context_dir_path_str = os.path.dirname(p_dockerfile_path_str)

	p_log_fun("INFO", "====================+++++++++++++++=====================")
	p_log_fun("INFO", "                 BUILDING DOCKER IMAGE")
	p_log_fun("INFO", "image_names - %s"%(p_image_names_lst))
	p_log_fun("INFO", "Dockerfile  - %s"%(p_dockerfile_path_str))
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
	])

	# TAGS - there can be multiple tags for the same image
	for n in p_image_names_lst:
		cmd_lst.append("--tag=%s"%(n))

	# BUILD_ARGS
	if not p_build_args_map == None:
		for k, v in p_build_args_map.items():
			cmd_lst.append("--build-arg %s=%s"%(k, v))

	# CONTEXT_DIR
	cmd_lst.append(context_dir_path_str)

	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", c_str)

	# change to the dir where the Dockerfile is located, for the 'docker'
	# tool to have the proper context
	old_cwd = os.getcwd()
	os.chdir(context_dir_path_str)
	
	#---------------------------------------------------
	def get_image_id_from_line(p_stdout_line_str):
		p_lst = p_stdout_line_str.split(" ")

		assert len(p_lst) == 3
		image_id_str = p_lst[2]

		# IMPORTANT!! - check that this is a valid 12 char Docker ID
		assert len(image_id_str) == 12
		return image_id_str

	#---------------------------------------------------

	p = subprocess.Popen(c_str, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, bufsize=1)
	
	image_id_str = ""
	for line in iter(p.stdout.readline, b''):	
		line_str = line.strip().decode("utf-8")
		
		print(line_str)

		# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
		if "Cannot connect to the Docker daemon" in line_str:
			exit(1)
		
		# # CONTAINER_STARTED
		# if line_str.startswith("Successfully built"):
		#
		# 	image_id_str = get_image_id_from_line(line_str)
		# 	print("image ID - %s"%(image_id_str))
	

	p.communicate()
	exit_code_int = p.returncode
	

	# IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	#               in other scenarious its acceptable for this command to fail, and we want the caller
	#               to keep executing.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	# change back to old dir
	os.chdir(old_cwd)

#-------------------------------------------------------------
# LOGIN
def login(p_docker_user_str,
	p_docker_pass_str,
	p_host_str          = None,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = False):
	assert isinstance(p_docker_user_str, str)
	assert isinstance(p_docker_pass_str, str)



	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker login",
		"-u", p_docker_user_str,
		"--password-stdin"
	])

	# HOST
	if not p_host_str == None:
		cmd_lst.append(p_host_str)

	print(" ".join(cmd_lst))
	process = subprocess.Popen(cmd_lst, shell=True, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
	
	# STDIN_WRITE
	process.stdin.write(bytes(p_docker_pass_str.encode())) # write password on stdin of "docker login" command
	
	stdout, stderr = process.communicate() # wait for command completion
	stdout_str = stdout.decode()
	stderr_str = stderr.decode()
	print(stdout_str)

	if not stderr_str == "":
		print(stderr_str)
		
	# IMPORTANT!! - failure to reach Dcoerk daemon should always exit. its not a expected failure.
	if "Cannot connect to the Docker daemon" in stderr_str:
		exit(1)

	# IMPORTANT!! - if command returns a non-zero exit code in some environments (CI) we
    #               want to fail with that a non-zero exit code - this way CI will flag builds as failed.
	#               in other scenarious its acceptable for this command to fail, and we want the caller
	#               to keep executing.
	if p_exit_on_fail_bool:
		if not process.returncode == 0:
			exit()

#---------------------------------------------------
# LOGIN__REMOTE
def login__remote(p_docker_user_str,
	p_docker_pass_str,
	p_fab_conn,
	p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.login__remote()")
	assert isinstance(p_docker_user_str, str)
	assert isinstance(p_docker_pass_str, str)

	import fabric
	
	#---------------------------
	# UPLOAD_PASS_FILE

	pass_f_str = "tmp_file"
	f = open(pass_f_str, "w")
	f.write(p_docker_pass_str)
	f.close()
	
	fabric.api.put(pass_f_str) # upload password file
	
	#---------------------------
	# IMPORTANT!! - specify pasword from stdin so that it doesnt show up
	#               as a part of the final command (in logs)

	c_str = "cat %s | sudo docker login -u %s --password-stdin"%(pass_f_str, p_docker_user_str)
	p_fab_conn.run(c_str)
	# fabric.api.run(c_str)

	#---------------------------
	p_fab_conn.run("rm %s"%(pass_f_str))
	delegator.run("rm %s"%(pass_f_str)) # clean local tmp_file that holds the dockerhub password

#---------------------------------------------------
# LOGIN__REMOTE_FROM_FILE
def login__remote_from_file(p_docker_user_str,
	p_docker_pass_str,
	p_fab_conn,
	p_log_fun,
	p_ecr_bool           = False,
	p_aws_account_id_str = None,
	p_region_str         = "us-east-1"):
	p_log_fun("FUN_ENTER", "gf_os_docker.login__remote_from_file()")
	assert isinstance(p_docker_user_str, str)
	assert isinstance(p_docker_pass_str, str)

	#---------------------------
	# UPLOAD_PASS_FILE

	pass_f_str = "tmp_file"
	f = open(pass_f_str, "w")
	f.write(p_docker_pass_str)
	f.close()

	p_fab_conn.put(pass_f_str) # upload password file

	#---------------------------
	# IMPORTANT!! - specify pasword from stdin so that it doesnt show up
	#               as a part of the final command (in logs)

	if p_ecr_bool:
		assert not p_aws_account_id_str == None
		ecr_acc_str = f"{p_aws_account_id_str}.dkr.ecr.{p_region_str}.amazonaws.com"
		p_fab_conn.run(f"cat {pass_f_str} | sudo docker login -u {p_docker_user_str} --password-stdin {ecr_acc_str}")
	else:
		p_fab_conn.run(f"cat {pass_f_str} | sudo docker login -u {p_docker_user_str} --password-stdin")
	
	#---------------------------
	p_fab_conn.run("rm %s"%(pass_f_str))
	delegator.run("rm %s"%(pass_f_str)) # clean local tmp_file that holds the dockerhub password

#---------------------------------------------------
# CLEAN_STOP__CONTAINERS
def clean_stop__containers(p_cont_image_name_str,
	p_fab_conn,
	p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.clean_stop__containers()")

	#--------------------
	# STOP_CURRENT_CONTAINERS
	image_ids_str = p_fab_conn.run("sudo docker ps -a | grep %s | awk '{print $1}'"%(p_cont_image_name_str))
	print("image_ids_str - %s"%(image_ids_str))

	if not image_ids_str == "":
		print("    >>  image already running - %s"%(p_cont_image_name_str))
		print("    >>  stoping containers    - %s"%(p_cont_image_name_str))

		for l in image_ids_str.split("\n"):
			image_id_str = l
			p_fab_conn.run("sudo docker stop %s"%(image_id_str)) # stop first
			p_fab_conn.run("sudo docker rm %s"%(image_id_str))   # remove, to not conflict with new ones

	#--------------------

#---------------------------------------------------
def install_base_docker(p_fab_conn,
	p_log_fun,
	p_os_type_str = "ubuntu"):
	p_log_fun("FUN_ENTER", "gf_os_docker.install_base_docker()")
	assert p_os_type_str == "ubuntu" or p_os_type_str == "amazon_linux_2"
	
	if p_os_type_str == "ubuntu":
		p_fab_conn.run("sudo apt-get clean")
		p_fab_conn.run("sudo apt-get update")
		p_fab_conn.run("sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common")
		p_fab_conn.run("sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -")
		

		# last one validated for Kubernetes installations
		# docker_version_str = "17.06.0~ce-0~ubuntu"

		##FIX!! - hardcoding to "zesty" ubuntu version (17.04) because in 17.10 at the moment (dec 10 2017) there is no docker-ce package
		##        so Im hardcdoing 17.04 just for the moment so that the compatible docker-ce package is used
		##p_fab_api.run('sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"')
		#p_fab_api.run('sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu zesty stable"')

		p_fab_conn.run("sudo apt-get update")
		p_fab_conn.run("sudo apt-get install -y \
			apt-transport-https \
			ca-certificates \
			curl \
			gnupg-agent \
			software-properties-common")
		p_fab_conn.run("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -")
		p_fab_conn.run('sudo add-apt-repository \
			"deb [arch=amd64] https://download.docker.com/linux/ubuntu \
			$(lsb_release -cs) \
			stable"')
		p_fab_conn.run("sudo apt-get update")
		p_fab_conn.run("sudo apt-get install -y docker-ce docker-ce-cli containerd.io")

	elif p_os_type_str == "amazon_linux_2":

		p_fab_conn.run("sudo yum -y install docker")
		p_fab_conn.run("sudo systemctl start docker")


#---------------------------------------------------
def dockerhub__get_auth_config_json(p_docker_user_str,
	p_docker_pass_str,
	p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.dockerhub__get_auth_config_json()")

	auth_str             = base64.b64encode("%s:%s"%(p_docker_user_str, p_docker_pass_str))
	auth_config_map      = {"auths": {"https://index.docker.io/v1/": {"auth": auth_str}}}
	auth_config_json_str = json.dumps(auth_config_map)
	return auth_config_json_str