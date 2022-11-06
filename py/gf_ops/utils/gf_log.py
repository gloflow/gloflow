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

from colored import fg, bg, attr

#---------------------------------------------------
def log_fun(g, m):
    if g == "ERROR":
        print("%s%s%s:%s%s%s"%(bg("red"), g, attr(0), fg("red"), m, attr(0)))
    else:
        print("%s%s%s:%s%s%s"%(fg("yellow"), g, attr(0), fg("green"), m, attr(0)))