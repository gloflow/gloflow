# GloFlow application and media management/publishing platform
# Copyright (C) 2022 Ivan Trajkovic
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

import urllib.parse
import requests

#--------------------------------------------------
# RUN
def run(p_notify_completion_url_str,
	p_git_commit_hash_str = None,
	p_app_name_str        = None):
	
	url_str = None

	# add git_commit_hash as a querystring argument to the notify_completion URL.
	# the entity thats receiving the completion notification needs to know what the tag
	# is of the newly created container.
	if not p_git_commit_hash_str == None:
		
		url = urllib.parse.urlparse(p_notify_completion_url_str)
		
		# QUERY_STRING
		qs_lst = urllib.parse.parse_qsl(url.query)
		qs_lst.append(("base_img_tag", p_git_commit_hash_str)) # .parse_qs() places all values in lists

		qs_str = "&".join(["%s=%s"%(k, v) for k, v in qs_lst])

		# _replace() - "url" is of type ParseResult which is a subclass of namedtuple;
		#              _replace is a namedtuple method that:
		#              "returns a new instance of the named tuple replacing
		#              specified fields with new values".
		url_new = url._replace(query=qs_str)
		url_str = url_new.geturl()
	else:
		url_str = p_notify_completion_url_str

	print("NOTIFY_COMPLETION - HTTP REQUEST - %s"%(url_str))
	print(f"GIT commit_hash - {p_git_commit_hash_str}")

	#--------------------------
	# HTTP_POST

	data_map = {}
	if not p_app_name_str == None:
		data_map["app_name"] = p_app_name_str

	r = requests.post(url_str, json=data_map)
	print(r.text)

	if not r.status_code == 200:
		print("notify_completion http request failed")
		exit(1)
	
	#--------------------------