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

///<reference path="../../../../d/jquery.d.ts" />
///<reference path="../../../../d/jquery.timeago.d.ts" />

import "./gf_crawl__img_preview_tooltip";

namespace gf_crawl_images_browser {
//---------------------------------------------------
export function init__recent_images(p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_crawl_images_browser.init__recent_images()");

	http__get_recent_images((p_recent_images_lst)=>{

			const recent_images__browser = view__recent_images(p_recent_images_lst, p_log_fun);
			$('#browser').append(recent_images__browser);
		},
		()=>{},
		p_log_fun);
}
//---------------------------------------------------
function view__recent_images(p_recent_images_lst, p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_crawl_images_browser.view__recent_images()");

	const browser  = $(`
		<div id="browser_container">
			<div id="recent_images"></div>
		</div>`);

	for (var domain_map of p_recent_images_lst) { 

		const domain_str               = domain_map['domain_str'];
		const imgs_count_int           = domain_map['imgs_count_int'];
		const crawler_page_img_ids_lst = domain_map['crawler_page_img_ids_lst'];
		const creation_times_lst       = domain_map['creation_times_lst'];
		const urls_lst                 = domain_map['urls_lst'];
		const nsfv_ls                  = domain_map['nsfv_ls'];
		const origin_page_urls_lst     = domain_map['origin_page_urls_lst'];

		const domain_e = $(`
			<div class="domain">
				<div class="title"><a target="_blank" href=http://"`+domain_str+`">`+domain_str+`</a></div>
				<div class="imgs_count"><span>imgs_count: </span>`+imgs_count_int+`</div>
				<div class="urls"></div>
			</div>`);
		$(browser).find('#recent_images').append(domain_e);

		var i=0;
		for (var url_str of urls_lst) {

			//------------------
			//CREATION_TIME
			const creation_unix_time_f = creation_times_lst[i];
			const creation_time_f      = parseFloat(creation_unix_time_f);
			const creation_date        = new Date(creation_time_f*1000);
			const date_msg_str         = $.timeago(creation_date);
			//------------------

			const url_e = $(`
				<div class="url" data-url="`+url_str+`">
					<div class="url_a">
						<a target="_blank" href="`+url_str+`">`+url_str+`</a>
					</div>
					<div class="creation_time">`+date_msg_str+`</div>
				</div>`);

			//------------------
			//PREVIEW_TOOLTIP

			const crawler_page_image_id_str = crawler_page_img_ids_lst[i];
			const origin_page_url_str       = origin_page_urls_lst[i];
			gf_crawl__img_preview_tooltip.init(url_str,
				crawler_page_image_id_str,
				origin_page_url_str,
				url_e,
				p_log_fun);
			//------------------

			$(domain_e).find('.urls').append(url_e);
			i++;
		}

		//----------------------
		//IMPORTANT!! - mark all images that are listed for a domain, as potentially already
		//              existing in the target flow.
		check_imgs_existence(urls_lst,
			(p_img_existance__views_lst)=>{
				for (const [url_str,image_existence_e] of p_img_existance__views_lst) {
					$(domain_e).find('.urls').find(`[data-url="`+url_str+`"]`).append($(image_existence_e));
				}
			},
			(p_error_e)=>{
				$(domain_e).append(p_error_e);
			},
			p_log_fun);
		//----------------------
	}

	return browser;
}
//---------------------------------------------------
function check_imgs_existence(p_urls_lst,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_crawl_images_browser.check_imgs_existence()');

	http__check_imgs_exist_in_flow(p_urls_lst,
		(p_existing_images_lst)=>{

			const img_existance__views_lst = [];
			$.each(p_existing_images_lst, (p_i, p_e)=>{

				const img__id_str               = p_e['id_str'];
				const img__origin_url_str       = p_e['origin_url_str'];      //image url that was found in a page
				const img__origin_page_url_str  = p_e['origin_page_url_str']; //page url from which the image url came
				const img__creation_unix_time_f = p_e['creation_unix_time_f'];

				const element = get_view(img__origin_url_str,
					img__origin_page_url_str,
					img__creation_unix_time_f);

				img_existance__views_lst.push([img__origin_url_str, element]);
			});

			p_on_complete_fun(img_existance__views_lst);
		},
		(p_error_data_map)=>{
			const error_e = $('<div>check_img_exists_in_flow - ERROR</div>');
			p_on_error_fun(error_e);
		},
		p_log_fun);

	//---------------------------------------------------
	function get_view(p_existing_img__origin_url_str,
		p_existing_img__origin_page_url_str,
		p_existing_img__creation_unix_time_f) {
		//p_log_fun('FUN_ENTER','gf_crawl_images_browser.check_imgs_existence().get_view()');

		const date = new Date(p_existing_img__creation_unix_time_f*1000);

		const data_str             = date.getHours()+':'+date.getMinutes()+':'+date.getSeconds()+' - '+date.getDate()+'.'+date.getMonth()+'.'+date.getFullYear();
		const existing_img_preview = $('#page_info_container').find('img[src="'+p_existing_img__origin_url_str+'"]')[0];
		const element              = $(`
			<div class="img_exists">
				<div class="exists_msg">added</div>
				<div class="origin_page_url">`+p_existing_img__origin_page_url_str+`</div>
				<div class="creation_time">
					<span class="msg">created on:</span>
					<span class="time">`+data_str+`</span>
				</div>
			</div>`);

		return element;
	}
	//---------------------------------------------------
}
//---------------------------------------------------
//HTTP
//---------------------------------------------------
function http__get_recent_images(p_on_complete_fun, p_on_error_fun, p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_crawl_images_browser.http__get_recent_images()');

	const url_str = '/a/crawl/image/recent';
	$.get(url_str,
		(p_data_map)=>{
			console.log('response received');
			//const data_map = JSON.parse(p_data);
			
			if (p_data_map["status_str"] == 'OK') {
				const recent_images_lst = p_data_map['data']['recent_images_lst'];
				p_on_complete_fun(recent_images_lst);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
}
//---------------------------------------------------
function http__check_imgs_exist_in_flow(p_images_extern_urls_lst,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_crawl_images_browser.http__check_imgs_exist_in_flow()');

	const url_str = '/images/flows/imgs_exist';
	//p_log_fun('INFO','url_str - '+url_str);

	const data_map = {
		'images_extern_urls_lst': p_images_extern_urls_lst,
		'flow_name_str':          'general', //check if image exists in specific flow
		'client_type_str':        'gchrome_ext'
	};

	//-------------------------
	//HTTP AJAX
	$.post(url_str,
		JSON.stringify(data_map),
		(p_data_map) => {
			//const data_map = JSON.parse(p_data);
			if (p_data_map["status_str"] == 'OK') {

				var existing_images_lst = p_data_map['data']['existing_images_lst'];

				//FIX!! - sometimes the backend returns existing_images_lst as null
				//        when there are no images, instead of []. look into that
				if (existing_images_lst == null) {
					existing_images_lst = [];
				}
				p_on_complete_fun(existing_images_lst);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
	//-------------------------    
}
//---------------------------------------------------
}