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

//-------------------------------------------------
extract__main(log_fun)

//-------------------------------------------------
function log_fun(p_g,p_m) {
	var msg_str = p_g+':'+p_m
	//chrome.extension.getBackgroundPage().console.log(msg_str);

	switch (p_g) {
		case "INFO":
			console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
			break;
		case "FUN_ENTER":
			console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
			break;
	}
}

//---------------------------------------------------
function extract__main(p_log_fun) {
	p_log_fun('FUN_ENTER','extract.extract__main()');

	chrome.extension.onMessage.addListener(
		(p_request, p_sender, p_send_response_fun) => {

			const request_source_str = p_request.source_str;
			const request_type_str   = p_request.request_type_str;

			if (request_source_str == 'popup') {
				handle_msg(request_type_str, p_log_fun);
			}
		});
}

//---------------------------------------------------
/*IMPORTANT!! - this script is run in all frames of the particular page.
it may get run in advertising iframes as well, so potentially 10's of these scripts
are run in the same page context.
when the popup sends a message to the particular tab, all these scripts receive it. 
however chrome extensions have a limit where only the first invocation of p_send_response_fun
is executed, and its response returned to the popup sender as response. All other responses
are ignored, and potentially data extracted not used.
to avoid this, p_send_response_fun is not used at all, and instead data results are sent to the 
background page for storage, and then that data is displayed to the user by the
page_element_picker content_script*/

function handle_msg(p_request_type_str, p_log_fun) {
	p_log_fun('FUN_ENTER','extract.handle_msg()');

	switch(p_request_type_str) {
		//-------------
		//GET PAGE IMAGE INFOS

		case 'get_page_img_infos':
			const new_page_img_infos_lst = get_images_info(p_log_fun);
			const msg_map = {
				'page_img_infos_lst': new_page_img_infos_lst
			};
			send_to_bg_page(msg_map, (p_resp)=>{});
			break;

		//-------------
		//GET PAGE VIDEO INFOS

		case 'get_page_videos_infos':
			const new_page_video_infos_lst = get_videos_info(p_log_fun);
			const msg_map = {
				'page_videos_infos_lst':new_page_video_infos_lst
			};
			send_to_bg_page(msg_map,(p_resp)=>{});
			break;

		//-------------
	}
}

//---------------------------------------------------
function send_to_bg_page(p_msg_map, p_on_complete_fun) {
	p_log_fun('FUN_ENTER','exctract.send_to_bg_page()');

	chrome.extension.sendRequest(p_msg_map,
		(p_resp)=>{
			p_on_complete_fun(p_resp);
		});
}

//---------------------------------------------------
//ADD!! - detect you tube embeds in other non-youtube.com pages
//        via the <embed> tag

//->:List<:Dict(video_info_map)>
function get_videos_info(p_log_fun) {
	p_log_fun('FUN_ENTER', 'extract.get_videos_info()')

	const page_url_str    = window.location.toString();
	const videos_info_lst = [];
	
	//------------------------------------
	//YOUTUBE.COM DOMAIN
	
	//the user is currently watching a video on youtube.com
	//"\/\/" is "//" escaped
	if (page_url_str.match("^http:\/\/www.youtube.com\/watch") || page_url_str.match("^https:\/\/www.youtube.com\/watch")) {

		const youtube_video_embed_url_str = $('link[itemprop="embedURL"]').attr('href');
		p_log_fun('INFO', 'youtube_video_embed_url_str:'+youtube_video_embed_url_str);
		
		const video_info_map = {
			'type_str':        'video',
			'page_url_str':    page_url_str,
			'video_source_str':'youtube',
			'embed_url_str':   youtube_video_embed_url_str
		};
		
		videos_info_lst.push(video_info_map);
	}
	//------------------------------------
	//if its any other page, search all elements 'src' attribute
	//and see if it contains "http://www.youtube.com", "http://player.vimeo.com"
	else {
		
		//"*=" - Attribute Contains Selector [name*="value"]
		//Selects elements that have the specified attribute 
		//with a value containing the a given substring.
		//------------------------------------
		//YOUTUBE - IFRAME EMBED
		
		$('*[src*="https://www.youtube.com"]').each((p_i,p_element) => {
			p_log_fun('INFO', 'YOUTUBE IFRAME EMBED++++++++++++++++++++++++++++++++');
			p_log_fun('INFO', $(p_element).attr('src'));
			
			const video_info_map = {
				'type_str':        'video',
				'page_url_str':    page_url_str,
				'video_source_str':'youtube',
				'embed_url_str':   $(p_element).attr('src')
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
		//VIMEO - IFRAME EMBED
		
		$('*[src*="http://player.vimeo.com"]').each((p_i,p_element) => {
			p_log_fun('INFO', 'VIMEO IFRAME EMBED++++++++++++++++++++++++++++++++');
			p_log_fun('INFO', $(p_element).attr('src'));
			
			const video_info_map = {
				'type_str':        'video',
				'page_url_str':    page_url_str,
				'video_source_str':'vimeo',
				'embed_url_str':   $(p_element).attr('src')
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
		//VIMEO - FLASH PLAYER (<OBJECT> TAG)
		
		$('object[data*="http://a.vimeocdn.com"]').each((p_i,p_element) => {
			p_log_fun('INFO','VIMEO FLASH PLAYER OBJECT TAG++++++++++++++++++++++++++++++++');
			p_log_fun('INFO',$(p_element).attr('data'));
			
			const flash_vars_str  = $(p_element).find('*[name*="flashvars"]').attr('value');
			const video_info_map = {
				'type_str':            'video',
				'page_url_str':        page_url_str,
				'video_source_str':    'vimeo',
				'vimeo_flash_vars_str':flash_vars_str
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
		//OOYALA - IFRAME EMBED
		
		$('*[src*="http://player.ooyala.com"]').each((p_i,p_element) => {
			p_log_fun('INFO', 'OOYALA IFRAME EMBED++++++++++++++++++++++++++++++++');
			p_log_fun('INFO', $(p_element).attr('src'));
			
			const video_info_map = {
				'type_str':        'video',
				'page_url_str':    page_url_str,
				'video_source_str':'ooyala',
				'embed_url_str':   $(p_element).attr('src')
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
	}
	//------------------------------------

	return videos_info_lst;
}

//---------------------------------------------------
//->:List<:Dict()>(img_infos_lst)
function get_images_info(p_log_fun) {
	p_log_fun('FUN_ENTER', 'extract.get_images_info()')
	
	const page_url_str        = window.location.toString();
	const min_image_dimension = 20;

	//---------------------------------------------------
	function get_image_info(p_jq_element) {
		p_log_fun('FUN_ENTER', 'extract.get_images_info().get_image_info()')
		
		//".src" instead of ".attr('src')" - gets the fully resolved url (including the host)
		//                                   and not just the value thats stored in the "src" html attr
		const full_img_src_str = $(p_jq_element)[0].src;
		const img_name_str     = full_img_src_str.split('/').pop();
		const img_width        = $(p_jq_element).width();
		const img_height       = $(p_jq_element).height();

		const img_info_map = {
			'type_str':        'image',
			'page_url_str':    page_url_str,
			'full_img_src_str':full_img_src_str,
			'img_name_str':    img_name_str,

			'img_width' :img_width,
			'img_height':img_height
		};
		console.log(img_info_map);

		return img_info_map;
	}
	//---------------------------------------------------

	const img_infos_lst = [];
	$('img').each((p_i, p_element) => {
		const img_info_map = get_image_info(p_element);
		
		if (img_info_map['img_width']  > min_image_dimension && img_info_map['img_height'] > min_image_dimension) {

			//only use the image if both of its dimensions are larger
			//then the minimum treshold
			img_infos_lst.push(img_info_map);
		}
	});

	return img_infos_lst;
}