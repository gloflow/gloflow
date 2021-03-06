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

/*$(document).ready(()=>{
	main(log_fun)
});*/
main(log_fun);
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
//-------------------------------------------------
function main(p_log_fun) {
	p_log_fun('FUN_ENTER','popup.main()');
	
	init_selected_elements_view(p_log_fun);
	init_buttons(p_log_fun);
}
//-------------------------------------------------
function init_buttons(p_log_fun) {
	p_log_fun('FUN_ENTER','popup.init_buttons()');
	
	//-----------------------
	//CREATE POST
	
	$(document).on('click','#create_post_btn',()=>{

		//-------------
		//TARGET_HOST

		//const host_str = 'http://gloflow.com';
		const target_host_str = $('#target_host input').val();
		p_log_fun('INFO','target_host_str - '+target_host_str);
		//-------------

		$('#create_post_btn').css('background-color','yellow');
		$('#create_post_btn').css('color'           ,'black');

		get__selected_elements((p_selected__post_elements_lst)=>{
				get__post_origin_page_url((p_post_origin_page_url_str)=>{
					//-----------------------
					//CREATE_POST
					http__create_post(p_selected__post_elements_lst,
						p_post_origin_page_url_str,
						target_host_str,
						(p_images_job_id_str) => {

							$('#create_post_btn').css('background-color','green');

							//ADD!! - some visual indicator of success
							$('#create_post_btn').css('background-color','green');
							p_log_fun('INFO', p_images_job_id_str);

							//-------------------
							//IMAGE_JOB_STATUS
							post_images_job_status(p_images_job_id_str, target_host_str, p_log_fun);
							//-------------------
						},
						(p_err_data)=>{
							$('#create_post_btn').css('background-color','red');
						},
						p_log_fun);
					//-----------------------
				}, p_log_fun);
			}, p_log_fun);
	});
	//-----------------------
	//TELL CONTENT-SCRIPT TO GET IMAGES INFO
	
	$(document).on('click', '#get_tab_page_images_btn', (p_e)=>{

		//POPUP->CONTENT_SCRIPT
		get_page_img_infos__from_content_scr((p_img_infos_lst)=>{
				display_page_info_in_content_scr(()=>{}, p_log_fun);
			},
			p_log_fun);
	})
	//-----------------------
	//TELL CONTENT-SCRIPT TO GET VIDEOS INFO
	
	$(document).on('click', '#get_tab_page_videos_btn', (p_e)=>{
		//POPUP->CONTENT_SCRIPT
		get_page_video_infos__from_content_scr((p_videos_infos_lst)=>{
				display_page_info_in_content_scr(()=>{}, p_log_fun);
			},
			p_log_fun);
	})
	//-----------------------
	//SHOW SELECTED ASSETS

	$(document).on('click','#show_selected_assets_btn',() => {
		$('body').css('background-color','red');

		const selected_assets_str = chrome.extension.getURL('html/selected_elements_ui.html');
		chrome.tabs.create({
				'url':selected_assets_str
			},
			() => {});
	});
	//-----------------------
}
//---------------------------------------------------
//BACKGROUND_PAGE COMM
//---------------------------------------------------
function clear__selected_elements(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER','popup.get__selected_elements()');

	const msg_map = {
		'source_str':'popup',
		'type_str':  'clear__selected_elements',
	};
	chrome.extension.sendRequest(msg_map,
		(p_response) => {
			p_on_complete_fun();
		});
}
//---------------------------------------------------
function get__selected_elements(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER','popup.get__selected_elements()');

	const msg_map = {
		'source_str':'popup',
		'type_str':  'get__selected_elements',
	};
	chrome.extension.sendRequest(msg_map,
		(p_response) => {
			const selected_elements_map = p_response.selected_elements_map;
			const selected_images_lst   = selected_elements_map['images_lst'];
			const selected_videos_lst   = selected_elements_map['videos_lst'];

			const selected_post_elements_lst = [];

			//IMAGES
			$.each(selected_images_lst,
				(p_i,p_image_map) => {

					selected_post_elements_lst.push({
						'type_str':      'image',
						'extern_url_str':p_image_map['full_img_src_str']
					});
				});

			//VIDEOS
			$.each(selected_videos_lst,
				(p_i,p_video_map) => {
					selected_post_elements_lst.push({
						'type_str':      'video',
						'extern_url_str':p_video_map['full_img_src_str']
					});
				});

			p_on_complete_fun(selected_post_elements_lst);
		});

	/*const selected_lst = [];
	$('#picked_page_assets_lst').find('a').each((p_i,p_element) => {
		const asset_url_str    = $(p_element).attr('href');
		const post_element_info_map = {
			'type_str'      :$(p_element).attr('data-type_str'),                  
			'extern_url_str':asset_url_str
		};
			
		selected_lst.push(post_element_info_map);
	});*/
}
//-------------------------------------------------
//CONTENT_SCRIPT COMM
//-------------------------------------------------
function get_page_img_infos__from_content_scr(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER','popup.get_page_img_infos__from_content_scr()');

	chrome.tabs.getSelected(null,(p_tab)=>{

		//IMPORTANT!! - popup just signals to the content script thats running in the tab
		//              to get page videos (get_page_videos_infos msg), without expecting results back.
		//              instead content scripts (multiple running in iframes of the tab) send that data
		//              to the background page (due to Chrome tabs.sendMessage() limitations)
		const msg_info_map = {
			'source_str':'popup',
			'type_str':  'get_page_img_infos'
		};

		//send a message to the particular tab where the content-script is running
		chrome.tabs.sendMessage(p_tab.id, msg_info_map, {},
			(p_response) => {
				//const page_img_infos_map = p_response.page_img_infos_map;
				//p_on_complete_fun(page_img_infos_map);
				p_on_complete_fun();
			});
	});
}
//-------------------------------------------------
function get_page_video_infos__from_content_scr(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER', 'popup.get_page_video_infos__from_content_scr()');

	chrome.tabs.getSelected(null,(p_tab)=>{

		//IMPORTANT!! - popup just signals to the content script thats running in the tab
		//              to get page videos (get_page_videos_infos msg), without expecting results back.
		//              instead content scripts (multiple running in iframes of the tab) send that data
		//              to the background page (due to Chrome tabs.sendMessage() limitations)
		const msg_info_map = {
			'source_str':'popup',
			'type_str':  'get_page_videos_infos'
		};
		chrome.tabs.sendMessage(p_tab.id, msg_info_map, {},
			(p_response) => {
				//const page_videos_infos_map = p_response.page_videos_infos_lst;
				//p_on_complete_fun(page_videos_infos_map);
				p_on_complete_fun();
			});
	});
}
//-------------------------------------------------
function get__post_origin_page_url(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER', 'popup.get__post_origin_page_url()');

	chrome.tabs.getSelected(null, (p_tab) => {
		const msg_info_map = {
			'source_str':'popup',
			'type_str':  'get_post_origin_page_url'
		};
		chrome.tabs.sendMessage(p_tab.id, msg_info_map, {},
			(p_response) => {
				const post_origin_page_url_str = p_response;
				p_on_complete_fun(post_origin_page_url_str);
			});
	});
}
//-------------------------------------------------
function display_page_info_in_content_scr(p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER','popup.display_page_info_in_content_scr()');

	chrome.tabs.getSelected(null,(p_tab) => {
		const msg_info_map = {
			'source_str':'popup',
			'type_str':  'display_page_info'
		};
		chrome.tabs.sendMessage(p_tab.id, msg_info_map, {},
			(p_response) => {
				p_on_complete_fun();
			});
	});
}
//-------------------------------------------------
//VAR
//-------------------------------------------------
function run_script_in_tab(p_script_code_str, p_tab_id, p_log_fun) {
	p_log_fun('FUN_ENTER','popup.run_script_in_tab()')	
	const details_map = {
		'code':p_script_code_str
	};
	chrome.tabs.executeScript(p_tab_id, details_map, ()=>{});
}