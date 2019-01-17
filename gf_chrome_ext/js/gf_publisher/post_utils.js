/*
GloFlow media management/publishing system
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
function http__create_post(p_selected__post_elements_lst,
	p_post_origin_page_url_str,
	p_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER','post_utils.http__create_post()');

	//---------------------------------------------------
	//->:Map
	function extract_post_form_info() {
		p_log_fun('FUN_ENTER','post_utils.http__create_post().extract_post_form_info()');
		
		const title_str            = $('#post_title_str').val();
		const description_str      = $('#post_description_str').val();
		const tags_str             = $('#post_tags_str').val();
		const poster_user_name_str = $('#poster_user_name_str').val();
		
		p_log_fun('INFO','title_str           :'+title_str);
		p_log_fun('INFO','description_str     :'+description_str);
		p_log_fun('INFO','poster_user_name_str:'+poster_user_name_str);
		//------------------------
		//POST ELEMENTS
		
		//post elements examples:
		//{'type_str':'image','url_str'    :'http://gloflow.com/i.png'},
		//{'type_str':'text' ,'content_str':'this is a block of text'},
		//{'type_str':'link' ,'url_str'    :'http://www.yahoo.com'},
		const post_elements_lst = [];
		
		//a link goes first
		post_elements_lst.push({
			'type_str':      'link',
			'extern_url_str':p_post_origin_page_url_str
		});

		//add selected elements to post_elements_lst
		Array.prototype.push.apply(post_elements_lst, p_selected__post_elements_lst);
		//------------------------
		
		const post_info_map = {
			'title_str':               title_str,
			'client_type_str':         'gchrome_ext',
			'tags_str':                tags_str, //comma "," separated tags
			'description_str':         description_str,
			'poster_user_name_str':    poster_user_name_str,
			'post_elements_lst':       post_elements_lst,
			'post_origin_page_url_str':p_post_origin_page_url_str
		}
		
		return post_info_map
	}
	//---------------------------------------------------
	
	const post_info_map = extract_post_form_info();
	const post_info_str = JSON.stringify(post_info_map);
	const url_str       = p_host_str+'/posts/create';
	
	p_log_fun('INFO','------------------------------');
	p_log_fun('INFO','sending data to url:'+url_str);
	p_log_fun('INFO','post_info_str:'+JSON.stringify(post_info_str));
	
	//-------------------------
	//HTTP AJAX
	$.post(url_str,
		post_info_str,
		(p_data_map) => {
			console.log('response received');
			//const data_map = JSON.parse(p_data);

			if (p_data_map["status_str"] == 'OK') {

				const images_job_id_str = p_data_map['data']['images_job_id_str'];
				p_on_complete_fun(images_job_id_str);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
	//-------------------------
}
//---------------------------------------------------
function publish_post(p_post_title_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER','post_utils.publish_post()');
	
	const target_post_info_map = {
		'post_title_str':post_title_str
	};
	
	const data_str = JSON.stringify(target_post_info_map);
	const url_str  = p_host_str+'/posts/publish';
	
	$.post(url_str,
		{'data_str':data_str},
		(p_data_str) => {
			p_log_fun('INFO','p_data_str:'+p_data_str);
			
			const data_map = JSON.parse(p_data_str);
			
			if (data_map.status == 'OK') {
				p_on_complete_fun(data_map.data);
			}
			else {
				p_on_error_fun(data_map.data);
			}
		});
}