/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

//---------------------------------------------------
function http__check_imgs_exist_in_flow(p_images_extern_urls_lst,
	p_flow_name_str,
	p_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {

	const url_str = `${p_host_str}/v1/images/flows/imgs_exist`;
	p_log_fun('INFO', `url_str - ${url_str}`);

	const data_map = {
		'images_extern_urls_lst': p_images_extern_urls_lst,
		'flow_name_str':          p_flow_name_str,
		'client_type_str':        'gchrome_ext'
	};

	//-------------------------
	// HTTP

	fetch(url_str, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(data_map)
	  })
	  .then(response => response.json())
	  .then(p_data_map => {
		console.log('response received');
		console.log(`status - ${p_data_map["status"]}`);
	  
		if (p_data_map["status"] === 'OK') {
			var existing_images_lst = p_data_map['data']['existing_images_lst'];

			// FIX!! - sometimes the backend returns existing_images_lst as null
			//         when there are no images, instead of []. look into that
			if (existing_images_lst == null) {
				existing_images_lst = [];
			}
			p_on_complete_fun(existing_images_lst);
		}
		else {
			p_on_error_fun(p_data_map["data"]);
		}
	  })
	.catch(error => {
		console.log('An error occurred:', error);
	});

	//-------------------------	
}

//---------------------------------------------------
function http__add_image_to_flow(p_image_extern_url_str,
	p_image_origin_page_url_str,
	p_flows_names_lst,
	p_gf_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'image_utils.http__add_image_to_flow()');

	const url_str = `${p_gf_host_str}/v1/images/flows/add_img?auth_r=0`;
	p_log_fun('INFO', 'url_str - '+url_str);

	const data_map = {
		'image_extern_url_str':      p_image_extern_url_str,
		"image_origin_page_url_str": p_image_origin_page_url_str,
		'flows_names_lst':           p_flows_names_lst, // ['general'],
		'client_type_str':           'gchrome_ext',
	};

	//-------------------------
	// HTTP
	
	fetch(url_str, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(data_map)
	  })
	  .then(response => response.json())
	  .then(p_data_map => {
		console.log('response received');
		console.log(`status - ${p_data_map["status"]}`);
	  
		if (p_data_map["status"] === 'OK') {
			const images_job_id_str                = p_data_map['data']['images_job_id_str'];
			const image_id_str                     = p_data_map['data']['image_id_str'];
			const thumbnail_small_relative_url_str = p_data_map['data']['thumbnail_small_relative_url_str'];
		  
			p_on_complete_fun(images_job_id_str, image_id_str, thumbnail_small_relative_url_str);
		}
		else {
			p_on_error_fun(p_data_map["data"]);
		}
	  })
	.catch(error => {
		console.log('An error occurred:', error);
	});

	//-------------------------	
}

//---------------------------------------------------
function http__gif_get_info(p_gf_img_id_str,
	p_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'image_utils.http__gif_get_info()');

	const url_str = p_host_str+'/images/gif/get_info?gfimg_id='+p_gf_img_id_str;
	p_log_fun('INFO', 'url_str - '+url_str);

	//-------------------------
	// HTTP AJAX
	$.get(url_str,
		(p_data_map) => {
			console.log('response received');
			
			if (p_data_map["status"] == 'OK') {

				const gif_map = p_data_map['data']['gif_map'];
				p_on_complete_fun(gif_map);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
	//-------------------------	
}

//-------------------------------------------------
function post_images_job_status(p_images_job_id_str, p_host_str, p_log_fun) {
	p_log_fun('FUN_ENTER', 'image_utils.post_images_job_status()');

	const url_str = p_host_str+'/images/jobs/status?images_job_id_str='+p_images_job_id_str;
	p_log_fun('INFO', 'url_str - '+url_str);

	const source = new EventSource(url_str);

	source.onopen = (p_event) => {
		console.log("EVENT_SOURCE OPEN")
	}

	source.onerror = (p_event) => {
		console.log("EVENT_SOURCE ERROR")

		if (p_event.readyState == EventSource.CLOSED) {
			// Connection was closed.
		}
	}

	source.onmessage = (p_event) => {
		const data_map = JSON.parse(p_event.data);
		console.log(">>>> -- ");
		console.log(data_map);

		const image_id_str         = data_map['image_id_str'];
		const image_source_url_str = data_map['image_source_url_str'];

		const status_type_str = data_map['type_str'];
		switch (status_type_str) {
			case 'fetch_ok':
				break;
			case 'transform_error':
				const err_str = data_map['err_str'];
				break;
		}
	}
}