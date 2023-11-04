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

#---------------------------------------------------
# sends the HTTP request to the gf_images service, to process the image

#->:Map(image_processing_results_map)
def dispatch_process_extern_image(p_image_url_str,
	p_gf_images_service_host_str,
	p_log_fun,
	p_process_from_scratch_if_prexisting_bool = True):
	p_log_fun('FUN_ENTER', 'gf_images_client.dispatch_process_extern_image()')
	
	assert p_gf_images_service_host_str.startswith('http://') or \
		p_gf_images_service_host_str.startswith('https://')
				 
	# new_image_url_str = base64.b64encode(p_image_url_str)
	new_image_url_str = urllib.quote(p_image_url_str)
	url_str           = '%s/images/start_job'%(p_gf_images_service_host_str)
	headers_dict = {'accept': 'text/event-stream'} # use Server Side Events
	params_dict  = {
		'type':    'process_extern_image',
		'img_url': new_image_url_str
	}
	
	# params_dict  = {
	# 	'image_url_str':                           p_image_url_str,
	# 	'process_from_scratch_if_prexisting_bool': p_process_from_scratch_if_prexisting_str
	# }

	p_log_fun('INFO','url_str - [%s]'%(url_str))

	r = requests.get(url_str,
		params  = params_dict,
		headers = headers_dict,
		stream  = True) # HTTP Server-Side-Events
	assert r.status_code == 200

	#---------------------------------------------------
	# SSE - streaming data lines

	#->:List<:Map>
	def stream_responses():
		p_log_fun('FUN_ENTER','gf_images_client.dispatch_process_extern_image().stream_responses()')

		data_items_lst = []
		for line_str in r.iter_lines():

			p_log_fun('INFO','>>>>>>>>>>>>>>>>>>>>>>>>')
			p_log_fun('INFO',line_str)

			#filter out keep-alive new lines
			if line_str:
				if line_str.startswith('data: '):
					msg_str  = line_str.strip('data: ')
					msg_dict = json.loads(msg_str)

					assert msg_dict.has_key('status_str')
					status_str = msg_dict['status_str']

					assert isinstance(status_str,basestring)
					assert status_str == 'ok' or \
						   status_str == 'error'
					assert status_str == 'ok'

					data_dict = msg_dict['data_dict']
					assert isinstance(data_dict,dict)

					data_items_lst.append(data_dict)
		return data_items_lst

	#---------------------------------------------------
	data_items_lst = stream_responses()


	assert isinstance(data_items_lst,list)
	assert len(data_items_lst) == 1 #just the job_output_dict

	job_output_dict = data_items_lst[0]
	assert isinstance(job_output_dict,dict)


	p_log_fun('INFO', 'job_output_dict - %s'%(json.dumps(job_output_dict)))


	assert job_output_dict.has_key('image_id_str')
	image_id_str = job_output_dict['image_id_str']
	assert isinstance(image_id_str, basestring)
	
	assert job_output_dict.has_key('thumbs_info_map')
	thumbs_info_map = job_output_dict['thumbs_info_map']
	assert isinstance(thumbs_info_map, dict)

	assert thumbs_info_map.has_key('thumbnail_small_relative_url_str')
	thumbnail_small_relative_url_str = thumbs_info_map['thumbnail_small_relative_url_str']
	assert isinstance(thumbnail_small_relative_url_str, basestring)
	assert thumbnail_small_relative_url_str.startswith('/images')
	
	assert thumbs_info_map.has_key('thumbnail_medium_relative_url_str')
	thumbnail_medium_relative_url_str = thumbs_info_map['thumbnail_medium_relative_url_str']
	assert isinstance(thumbnail_medium_relative_url_str, basestring)
	assert thumbnail_medium_relative_url_str.startswith('/images')
	
	assert thumbs_info_map.has_key('thumbnail_large_relative_url_str')
	thumbnail_large_relative_url_str = thumbs_info_map['thumbnail_large_relative_url_str']
	assert isinstance(thumbnail_large_relative_url_str, basestring)
	assert thumbnail_large_relative_url_str.startswith('/images')
	
	image_processing_results_map = {
		'image_id_str':    image_id_str,
		'thumbs_info_map': {
 			'thumbnail_small_relative_url_str':  thumbnail_small_relative_url_str,
			'thumbnail_medium_relative_url_str': thumbnail_medium_relative_url_str,
			'thumbnail_large_relative_url_str':  thumbnail_large_relative_url_str
		}
	}
	
	return image_processing_results_map