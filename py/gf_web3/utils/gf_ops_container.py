# GloFlow application and media management/publishing platform
# Copyright (C) 2021 Ivan Trajkovic
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

import gf_core_cli

#--------------------------------------------------
def build(p_cont_image_name_str,
	p_dockerfile_path_str,
	p_docker_sudo_bool=False):
	
	docker_context_dir_str = f"{modd_str}/../.."

	print("BUILDING CONTAINER -----------=========================")
	print(f"container image name - {p_cont_image_name_str}")
	print(f"dockerfile           - {p_dockerfile_path_str}")
	
	assert os.path.isfile(p_dockerfile_path_str)

	c_lst = []
	if p_docker_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		"docker build",
		f"-f {p_dockerfile_path_str}",
		f"--tag={p_cont_image_name_str}",
		docker_context_dir_str
	])

	c_str = " ".join(c_lst)
	print(c_str)

	_, _, exit_code_int = gf_core_cli.run(c_str)

	if not exit_code_int == 0:
		exit(1)

#--------------------------------------------------
def publish(p_cont_image_name_str,
	p_docker_user_str,
	p_docker_pass_str,
	p_docker_sudo_bool=False):

	print("PUBLISHING CONTAINER -----------=========================")
	print(f"container image name - {p_cont_image_name_str}")

	# LOGIN
	docker_login(p_docker_user_str,
		p_docker_pass_str,
		p_docker_sudo_bool = p_docker_sudo_bool)

	#------------------------
	c_lst = []
	if p_docker_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		f"docker push {p_cont_image_name_str}"
	])

	c_str = " ".join(c_lst)
	print(c_str)

	_, _, exit_code_int = gf_core_cli.run(c_str)

	if not exit_code_int == 0:
		exit(1)

	#------------------------

#--------------------------------------------------
def run(p_cont_image_name_str,
	p_docker_ports_lst=[],
	p_docker_sudo_bool=False):

	c_lst = []
	if p_docker_sudo_bool:
		c_lst.append("sudo")

	ports_str = ' '.join([f'-p {p}:{p2}' for p, p2 in p_docker_ports_lst])
	c_lst.extend([
		f"docker run {ports_str} {p_cont_image_name_str}"
	])

	c_str = " ".join(c_lst)
	print(c_str)

	_, _, exit_code_int = gf_core_cli.run(c_str)

	if not exit_code_int == 0:
		exit(1)

#-------------------------------------------------------------
# DOCKER_LOGIN
def docker_login(p_docker_user_str,
	p_docker_pass_str,
	p_docker_sudo_bool = False):
	assert isinstance(p_docker_user_str, str)
	assert isinstance(p_docker_pass_str, str)

	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker", "login",
		"-u", p_docker_user_str,
		"--password-stdin"
	])
	print(" ".join(cmd_lst))

	p = subprocess.Popen(cmd_lst, stdin = subprocess.PIPE, stdout = subprocess.PIPE, stderr = subprocess.PIPE)
	p.stdin.write(bytes(p_docker_pass_str.encode("utf-8"))) # write password on stdin of "docker login" command
	
	stdout, stderr = p.communicate() # wait for command completion
	stdout_str = stdout.decode("ascii")
	stderr_str = stderr.decode("ascii")

	if not stdout_str == "":
		print(stdout_str)
	if not stderr_str == "":
		print(stderr_str)

	if not p.returncode == 0:
		exit(1)

	# ERROR
	if "Error" in stderr_str or "unauthorized" in stderr_str:
		print("failed to Docker login")
		exit(1)