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

	for line in iter(p.stdout.readline, b''):	
		line_str = line.strip().decode("utf-8")
		print(line_str)

	for line in iter(p.stderr.readline, b''):	
		line_str = line.strip().decode("utf-8")
		print(line_str)

	p.communicate()
	
	return "", "", p.returncode