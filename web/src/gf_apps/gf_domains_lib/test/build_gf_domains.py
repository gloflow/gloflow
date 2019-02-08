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

import envoy

output_file_str = './bin/gf_domains_browser.js'
files_lst = [
	'./src/domains_browser/gf_domains_browser.ts',
	'./src/domains_browser/gf_domain.ts',
	'./src/domains_browser/gf_domains_conn.ts',
	'./src/domains_browser/gf_domains_infos.ts',
	'./src/domains_browser/gf_domains_search.ts',
	'../../../gf_core/client/src/gf_color.ts',
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
c = 'tsc --out %s %s'%(output_file_str,
					' '.join(files_lst))

print c
r = envoy.run(c)

print r.std_out
print r.std_err