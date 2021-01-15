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

from colored import fg, bg, attr

#-------------------------------------------------------------
def read_process_stdout(p_out, p_type_str, p_color_str):

	for line in iter(p_out.readline, b''):
		
		header_color_str = fg(p_color_str)
		line_str         = line.strip().decode("utf-8")

		# ERROR
		if "ERROR" in line_str or "error" in line_str:
			print("%s%s:%s%s%s%s"%(header_color_str, p_type_str, attr(0), bg("red"), line_str, attr(0)))
		else:
			print("%s%s:%s%s"%(header_color_str, p_type_str, attr(0), line_str))

	p_out.close()