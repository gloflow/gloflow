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

import os
import subprocess
import threading
from colored import fg, bg, attr


#---------------------------------------------------
def run__view_realtime(p_cmd_lst,
	p_env_map,
	p_view__type_str,
	p_view__color_str):

	# When shell=True the shell is the child process, and the commands are its children.
	# So any SIGTERM or SIGKILL will kill the shell but not its child processes.
	# The best way I can think of is to use shell=False, otherwise when you kill
	# the parent shell process, it will leave a defunct shell process.
	# CMD also has to be a list here, since its not being passed in as a string
	# to the child shell.
	p = subprocess.Popen(p_cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		env=p_env_map)

	t = threading.Thread(target=read_process_stdout, args=(p.stdout, p_view__type_str, p_view__color_str))
	t.start()

	return p



#---------------------------------------------------
# RUN
def run(p_cmd_str,
	p_env_map = {}):

	# env map has to contains all the parents ENV vars as well
	p_env_map.update(os.environ)

	p = subprocess.Popen(p_cmd_str,
		env     = p_env_map,
		shell   = True,
		stdout  = subprocess.PIPE,
		stderr  = subprocess.PIPE,
		bufsize = 1)

	stdout_lst = []
	for line in iter(p.stdout.readline, b''):	
		line_str = line.strip().decode("utf-8")
		print(line_str)
		stdout_lst.append(line_str)
	stderr_lst = []
	for line in iter(p.stderr.readline, b''):	
		line_str = line.strip().decode("utf-8")
		print(line_str)
		stderr_lst.append(line_str)

	p.communicate()
	
	return stdout_lst, stderr_lst, p.returncode

#-------------------------------------------------------------
def read_process_stdout(p_out,
	p_view_type_str,
	p_view_color_str):

	for line in iter(p_out.readline, b''):
		
		header_color_str = fg(p_view_color_str)
		line_str         = line.strip().decode("utf-8")

		# ERROR
		if "ERROR" in line_str or "error" in line_str:
			print("%s%s:%s%s%s%s"%(header_color_str, p_view_type_str, attr(0), bg("red"), line_str, attr(0)))
		else:
			print("%s%s:%s%s"%(header_color_str, p_view_type_str, attr(0), line_str))

	p_out.close()