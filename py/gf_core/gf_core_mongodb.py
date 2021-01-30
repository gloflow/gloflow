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

import pymongo

#----------------------------------------------
# ADD!! - figure out a smarter way to pick the right hostport from p_host_port_lst,
#         instead of just picking the first element

def get_client(p_log_fun, p_host_port_lst = ['127.0.0.1:27017']):
	p_log_fun('FUN_ENTER', 'gf_core_mongodb.get_client()')

	host_str,port_str = p_host_port_lst[0].split(':')

	mongo_client = pymongo.MongoClient(host_str, int(port_str))
	return mongo_client