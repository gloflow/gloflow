/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

main(log_fun)
extractor_init(log_fun);

//---------------------------------------------------
function main(p_log_fun) {
				
	// IMPORTANT!! - maintaining this state so that once the info is parsed
	//               it is cached 
	var page_img_infos_lst    = [];
	var page_videos_infos_lst = [];

	chrome.runtime.onMessage.addListener(
		(p_request,
		p_sender,
		p_send_response_fun) => {
			handle_msg(p_request, p_sender, p_send_response_fun);

			/*
			IMPORTANT!! - if some asynchronous code is being run in the content script,
				(as is the case here, in the handle functions above)
				returning true right away (while async ops are running) from the message listener is keeping
				the message port open until the async operation is done and for the sendResponse() call. 
			*/
			return true;
		});
	
	//---------------------------------------------------
	function handle_msg(p_request, p_sender, p_send_response_fun) {

		// p_log_fun('INFO', '============================================');
		// p_log_fun('INFO', 'page_element_picker received message');
		// p_log_fun('INFO','p_request.source_str:'+p_request.source_str)
		// p_log_fun('INFO','p_request.type_str  :'+p_request.type_str)
		
		const request_source_str = p_request.source_str;
		const request_type_str   = p_request.type_str;

		//-------------
		// MESSAGES FROM POPUP

		if (request_source_str == "popup") {
			p_log_fun('INFO', 'POPUP MSG');
			p_log_fun('INFO', 'request_type_str - '+request_type_str);

			var msg_map;
			switch(request_type_str) {
				//-------------
				// GET PAGE IMAGE INFOS

				case 'get_page_img_infos':
					const new_page_img_infos_lst = get_images_info(p_log_fun);
					page_img_infos_lst           = new_page_img_infos_lst;
					
					msg_map = {
						'page_img_infos_lst':page_img_infos_lst
					};
					p_send_response_fun(msg_map);
					break;
				//-------------
				// GET PAGE VIDEO INFOS

				case 'get_page_videos_infos':

					const new_page_video_infos_lst = get_videos_info(p_log_fun);
					page_videos_infos_lst          = new_page_video_infos_lst;

					msg_map = {
						'page_videos_infos_lst':page_videos_infos_lst
					};
					p_send_response_fun(msg_map);
					break;
				//-------------
				// get the url of the page from which images/videos are being extracted
				case 'get_post_source_page_url':
					const post_source_page_url  = window.location.href;
					p_send_response_fun(post_source_page_url);
					break;
					
				//-------------
				// DISPLAY PAGE INFO
				case 'display_page_info':

					p_log_fun('INFO', '========_____===');
					p_log_fun('INFO', page_img_infos_lst);
					p_log_fun('INFO', page_img_infos_lst.length);

					// determine whether you are a top level frame
					if (window.parent == window) {
						display_page_info(page_img_infos_lst, page_videos_infos_lst, p_log_fun);
					}
					break;
			}
		}
		//-------------
	}
	//---------------------------------------------------
}