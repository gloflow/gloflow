# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
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

import hashlib
import binascii

#---------------------------------------------------------------------------------
def id_create(p_unique_vals_for_id_lst,
    p_unix_time_f):

    h = hashlib.md5()

    h.update(str(p_unix_time_f).encode())

    for v in p_unique_vals_for_id_lst:
        h.update(v.encode())

    hex_str = binascii.hexlify(h.digest()).decode()
    id_str = hex_str

    return id_str