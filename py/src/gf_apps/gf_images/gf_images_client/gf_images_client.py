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

import json
import urllib
import requests

default_host_str = 'https://gloflow.com'

#---------------------------------------------------
# sends the HTTP request to the gf_images service, to process the image
# ->:Map(image_processing_results_map)

def add_image(p_image_url_str,
	p_log_fun  = None,
	p_host_str = default_host_str,
	p_process_from_scratch_if_prexisting_bool = True):
	
	assert p_host_str.startswith('http://') or \
		p_host_str.startswith('https://')
				 
	# new_image_url_str = base64.b64encode(p_image_url_str)
	new_image_url_str = urllib.quote(p_image_url_str)
	url_str           = f'{p_host_str}/images/jobs/start'
	headers_map = {
		'accept': 'text/event-stream' # use Server Side Events
	}
	params_map  = {
		'type':    'process_extern_image',
		'img_url': new_image_url_str
	}
	
	# params_map  = {
	# 	'image_url_str':                           p_image_url_str,
	# 	'process_from_scratch_if_prexisting_bool': p_process_from_scratch_if_prexisting_str
	# }

	if not p_log_fun == None:
		p_log_fun('INFO', 'url_str - [%s]'%(url_str))

	r = requests.get(url_str,
		params  = params_map,
		headers = headers_map,
		stream  = True) # HTTP Server-Side-Events
	assert r.status_code == 200

	'''
	#---------------------------------------------------
	# SSE - streaming data lines

	# ->:List<:Map>
	def stream_responses():
		
		data_items_lst = []
		for line_str in r.iter_lines():

			if not p_log_fun == None:
				p_log_fun('INFO', '>>>>>>>>>>>>>>>>>>>>>>>>')
				p_log_fun('INFO', line_str)

			# filter out keep-alive new lines
			if line_str:
				if line_str.startswith('data: '):
					msg_str  = line_str.strip('data: ')
					msg_map = json.loads(msg_str)

					assert msg_map.has_key('status_str')
					status_str = msg_map['status_str']

					assert isinstance(status_str, str)
					assert status_str == 'ok' or \
						status_str == 'error'
					assert status_str == 'ok'

					data_map = msg_map['data_map']
					assert isinstance(data_map, dict)

					data_items_lst.append(data_map)
		return data_items_lst

	#---------------------------------------------------
	data_items_lst = stream_responses()
	'''



	resp_map = json.loads(r.text)
	assert(isinstance(resp_map.has_key("status"), str))
	assert(isinstance(resp_map["data"], dict))


	data_map = resp_map["data"]

	running_job_id_str = data_map["running_job_id_str"]
	job_expected_outputs_lst = data_map["job_expected_outputs_lst"]

	assert isinstance(job_expected_outputs_lst, list)
	assert len(job_expected_outputs_lst) == 1

	new_image_info_map = job_expected_outputs_lst[0]
	assert isinstance(new_image_info_map, dict)




	assert(p_image_url_str == new_image_info_map["image_source_url_str"])

	# IMAGE_ID
	assert "image_id_str" in new_image_info_map.keys()
	image_id_str = new_image_info_map['image_id_str']
	assert isinstance(image_id_str, str)




	#----------------
	# THUMBS

	assert new_image_info_map.has_key('thumbnail_small_relative_url_str')
	thumbnail_small_relative_url_str = new_image_info_map['thumbnail_small_relative_url_str']
	assert isinstance(thumbnail_small_relative_url_str, str)
	assert thumbnail_small_relative_url_str.startswith('/images')
	
	assert new_image_info_map.has_key('thumbnail_medium_relative_url_str')
	thumbnail_medium_relative_url_str = new_image_info_map['thumbnail_medium_relative_url_str']
	assert isinstance(thumbnail_medium_relative_url_str, str)
	assert thumbnail_medium_relative_url_str.startswith('/images')
	
	assert new_image_info_map.has_key('thumbnail_large_relative_url_str')
	thumbnail_large_relative_url_str = new_image_info_map['thumbnail_large_relative_url_str']
	assert isinstance(thumbnail_large_relative_url_str, str)
	assert thumbnail_large_relative_url_str.startswith('/images')
	
	#----------------

	image_results_map = {
		'image_id_str': image_id_str,
		'thumbs_info_map': {
 			'thumbnail_small_relative_url_str':  thumbnail_small_relative_url_str,
			'thumbnail_medium_relative_url_str': thumbnail_medium_relative_url_str,
			'thumbnail_large_relative_url_str':  thumbnail_large_relative_url_str
		}
	}
	
	return image_results_map