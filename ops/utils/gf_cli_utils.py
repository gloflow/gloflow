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

from colored import fg, bg, attr
import delegator

sys.path.append("%s/../../go/gf_core/py"%(modd_str))
import gf_core_lib

#---------------------------------------------------
def run_cmd(p_cmd_str,
	p_env_map           = None,
	p_print_output_bool = True):
	
	return gf_core_cli.run_cmd(p_cmd_str,
		p_env_map           = p_env_map,
		p_print_output_bool = p_print_output_bool)