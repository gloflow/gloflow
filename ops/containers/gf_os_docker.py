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

import subprocess

#-------------------------------------------------------------
def login(p_dockerhub_user_str,
	p_dockerhub_pass_str,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = False):
	assert isinstance(p_dockerhub_user_str, basestring)
	assert isinstance(p_dockerhub_pass_str, basestring)

	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker", "login",
		"-u", p_dockerhub_user_str,
		"--password-stdin"
	])
	print(" ".join(cmd_lst))

	process = subprocess.Popen(cmd_lst, stdin = subprocess.PIPE, stdout = subprocess.PIPE)
	process.stdin.write(p_dockerhub_pass_str) # write password on stdin of "docker login" command
	stdout_str, stderr_str = process.communicate() # wait for command completion
	print(stdout_str)
	print(stderr_str)

	if p_exit_on_fail_bool:
		if not process.returncode == 0:
			exit()