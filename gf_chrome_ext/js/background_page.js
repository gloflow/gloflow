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

importScripts('utils/var.js');
importScripts('utils/image_utils.js');
importScripts('build/gf_tagger_http.js');

main(log_fun);

//---------------------------------------------------
function main(p_log_fun) {
	
	p_log_fun('INFO', `background_page started...`);
	
	const main_domain_str = "gloflow.com";
	const ctx_map = {
		'logged_in_bool': false,
		'selected_elements_map': {},
		// 'selected_images_lst' :[],
		// 'selected_videos_lst' :[]
	};
	
	chrome.runtime.onMessage.addListener((p_request, p_sender, p_send_response_fun)=>{

		on_request_received_fun(p_request, p_sender, p_send_response_fun);

		/*
		IMPORTANT!! - if some asynchronous code is being run in the background script,
			(as is the case here, in the handle functions above)
			returning true right away (while async ops are running) from the message listener is keeping
			the message port open until the async operation is done and for the sendResponse() call. 
		*/
		return true;
	});
	
	//---------------------------------------------------	
	function on_request_received_fun(p_request, p_sender, p_send_response_fun) {		

		p_log_fun('INFO', `background_page MSG RECEIVED ------------ ${p_request.source_str} | ${p_request.type_str}`);

		switch (p_request.source_str) {
			
			// AUTH
			case 'popup_auth':
				handle_auth_msg(p_request.type_str, p_send_response_fun)
					.then(() => {
						
					})
					.catch((error) => {
						console.log("Error caught:", error);
					});

				break;

			case 'popup':
				handle_popup_msg(p_request.type_str, p_request);
				break;
			case 'popup_selected_elements':
				handle_popup_selected_elements(p_request.type_str, p_log_fun);
				break;

			case 'content_script':
				handle_content_script_msg(p_request.type_str, p_send_response_fun, p_request);
				break;
			
		}
		
		//---------------------------------------------------
		// AUTH
		//---------------------------------------------------
		// HANDLE_AUTH_MSG
		function handle_auth_msg(p_msg_type_str, p_send_response_fun) {
			
			return new Promise(async (p_resolve_fun, p_reject_fun) => {

				switch (p_msg_type_str) {

					//----------------
					// LOGIN
					case "login":

						try {
							await Promise.all([
								check_cookie("gf_sess"),
								check_cookie("Authorization")
							]);

							// mark the extension state as logged in
							ctx_map["logged_in_bool"] = true;

							p_send_response_fun({"status_str": "OK"});
							p_resolve_fun();

						} catch (error) {
							p_send_response_fun({"status_str": "ERROR"});
							p_reject_fun(error);
						}

						
						break;
					
					//----------------
					// LOGGED_IN
					case "logged_in":
						
						p_send_response_fun({
							"logged_in_bool": ctx_map["logged_in_bool"]
						});

						p_resolve_fun();
						break;

					//----------------
				}
			});

			//---------------------------------------------------
			function check_cookie(p_name_str) {
				return new Promise((p_resolve_fun, p_reject_fun) => {

					// GET_COOKIE
					chrome.cookies.get({ url: `https://${main_domain_str}`, name: p_name_str },
						function(p_cookie) {
							if (p_cookie) {
								p_resolve_fun("cookie found");
								
							} else {
								// console.log(`Cookie ${p_name_str} not found`);
								p_reject_fun("cookie not found");
							}
						});
					});
			}
			
			//---------------------------------------------------
		}

		//---------------------------------------------------
		// DETECT_LOGOUT
		chrome.cookies.onChanged.addListener(function(p_change_info) {

			// cookie deletion
			if (p_change_info.removed) {

				const cookie_name_str = p_change_info.cookie.name;

				// check if cookies related to GF auth have been deleted, to infer that user has logged out
				if (cookie_name_str == "gf_sess" || cookie_name_str == "Authorization") {

					// mark the extension state as logged out
					ctx_map["logged_in_bool"] = false;
				}
			}
		});

		//---------------------------------------------------
		// VAR
		//---------------------------------------------------	
		function handle_popup_selected_elements(p_msg_type_str, p_request) {
			switch(p_msg_type_str) {
				//----------------
				// GET SELECTED ASSETS

				case "get_selected_elements":
					get__selected_elements(ctx_map,
						(p_selected_elements_map) => {
							console.log(p_selected_elements_map);
							const msg_map = {
								"selected_elements_map": p_selected_elements_map
							};
							p_send_response_fun(msg_map);
						},
						p_log_fun);
					break;

				//----------------
			}

			
		}

		//---------------------------------------------------	
		function handle_popup_msg(p_msg_type_str, p_request) {
			switch (p_msg_type_str) {
				//----------------
				// LOG_MSG
				case 'log_msg':
					console.log("POPUP:${p_request.msg_str}");
					break;

				//----------------
				// GET_SELECTED_ELEMENTS
				case "get__selected_elements":
					get__selected_elements(ctx_map,
						(p_selected_elements_map)=>{
							const msg_map = {
								"selected_elements_map": p_selected_elements_map
							};
							p_send_response_fun(msg_map);
						},
						p_log_fun);

				//----------------
				case "clear__selected_elements":
					clear__selected_elements(ctx_map,
						()=>{
							const msg_map = {};
							p_send_response_fun(msg_map);
						},
						p_log_fun);

				//----------------
			}
		}

		//---------------------------------------------------
		async function handle_content_script_msg(p_msg_type_str, p_send_response_fun, p_request) {

			switch (p_msg_type_str) {

				//----------------
				// ADD_TAGS_TO_IMAGE
				case "add_tags_to_image":
					{
						const image_system_id_str = p_request["image_system_id_str"];
						const tags_lst            = p_request["tags_lst"];
						const gf_host_str         = p_request["gf_host_str"];

						const object_type_str = "image";
						const meta_map = {};
						
						const response_map = await gf_tagger__http_add_tags_to_obj(tags_lst,
							image_system_id_str,
							object_type_str,
							meta_map,
							gf_host_str,
							p_log_fun);
					}

					break;

				//----------------
				// LOG_MSG
				case "log_msg":
					console.log(`CONTENT_SCR:${p_request.msg_str}`);
					break;

				//----------------
				// ADD_IMAGE_TO_FLOW
				case "add_image_to_flow":

					const full_img_src_str          = p_request["full_img_src_str"];
					const image_origin_page_url_str = p_request["image_origin_page_url_str"];
					const images_flows_names_lst    = p_request["images_flows_names_lst"];
					const gf_host_str               = p_request["gf_host_str"];

					add_image_to_flow(full_img_src_str,
						image_origin_page_url_str,
						images_flows_names_lst,
						gf_host_str,

						// p_on_complete_fun
						(p_image_id_str)=>{
							p_send_response_fun({
								"status_str":   "OK",
								"image_id_str": p_image_id_str
							});
						},

						// p_on_error_fun
						(p_data_map)=>{
							p_send_response_fun({
								"status_str": "ERROR",
								"data_map":   p_data_map,
							});
						},

						p_log_fun);
					
					break;
				
				//----------------
				// CHECK_IMAGES_EXIST
				
				case "check_images_exist":
					{
						const images_extern_urls_lst = p_request["images_extern_urls_lst"];
						const gf_host_str            = p_request["gf_host_str"];

						/*
						IMPORTANT!! - we want to check if an image exists in any of the flows,
							for a given user or "anon", instead of targeting the check at
							a specific flow.
						*/
						const flow_name_str = "all";

						check_images_exist_in_flow(images_extern_urls_lst,
							flow_name_str,
							gf_host_str,

							// on_complete_fun
							(p_existing_images_lst)=>{
								p_send_response_fun({
									"status_str": "OK",
									"existing_images_lst": p_existing_images_lst
								});
							},
							// on_error_fun
							()=>{
								p_send_response_fun({
									"status_str": "ERROR",
									"data_map":   p_data_map,
								});
							},
							p_log_fun);
					}
					break;
					
				//----------------
				/*
				// ADD_ELEMENT_TO_POST
				case "add_element_to_post":
					var element_info_map = p_request["element_info_map"];

					add_element_to_post(element_info_map,
						ctx_map,
						(p_status_str) => {
							const msg_map = {
								"status_str": p_status_str
							};
							p_send_response_fun(msg_map);
						},
						p_log_fun);
					break;

				//----------------
				// REMOVE_ELEMENT_FROM_POST
				case "remove_element_from_post":
					var element_info_map = p_request["element_info_map"];
					remove_element_from_post(element_info_map, ctx_map, p_log_fun);
					break;
				*/
				//----------------
			}
		}

		//---------------------------------------------------
	}

	//---------------------------------------------------
}

//---------------------------------------------------
// FLOW_OPS
//---------------------------------------------------
// ADD_IMAGE_TO_FLOW

function add_image_to_flow(p_full_img_src_str,
	p_image_origin_page_url_str,
	p_images_flows_names_lst,
	p_gf_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {

	http__add_image_to_flow(p_full_img_src_str,
		p_image_origin_page_url_str,
		p_images_flows_names_lst,
		p_gf_host_str,

		(p_images_job_id_str, p_image_id_str, p_thumbnail_small_relative_url_str)=>{

			console.log("image added")
			console.log(`image job ID    - ${p_images_job_id_str}`)
			console.log(`image ID        - ${p_image_id_str}`)
			console.log(`thumb small URL - ${p_thumbnail_small_relative_url_str}`)

			p_on_complete_fun(p_image_id_str);
		},
		(p_data_map)=>{p_on_error_fun(p_data_map)},
		p_log_fun);
}

//---------------------------------------------------
function check_images_exist_in_flow(p_images_extern_urls_lst,
	p_flow_name_str,
	p_gf_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {

	http__check_imgs_exist_in_flow(p_images_extern_urls_lst,
		p_flow_name_str,
		p_gf_host_str,

		// on_complete_fun
		(p_existing_images_lst)=>{
			p_on_complete_fun(p_existing_images_lst);
		},
		// on_error_fun
		(p_data_map)=>{
			p_on_error_fun(p_data_map);
		},
		p_log_fun);
}

//---------------------------------------------------
// POST_OPS
//---------------------------------------------------
function clear__selected_elements(p_ctx_map, p_on_complete_fun, p_log_fun) {

	// IMPORTANT!! - clear all currently selected elements
	p_ctx_map['selected_elements_map'] = {}; 
}

//---------------------------------------------------
function get__selected_elements(p_ctx_map, p_on_complete_fun, p_log_fun) {

	//-------
	const images_lst = [];
	const videos_lst = [];
	var   element_info_map;

	for (k in p_ctx_map['selected_elements_map']) {

		element_info_map = p_ctx_map['selected_elements_map'][k];

		switch(element_info_map['type_str']) {
			case 'image':
				images_lst.push(element_info_map);
				break;
			case 'video':
				videos_lst.push(element_info_map);
				break;
		}
	}
	//-------

	const selected_elements_map = {
		'images_lst': images_lst,
		'videos_lst': videos_lst
	};

	p_on_complete_fun(selected_elements_map);
}

//---------------------------------------------------
function add_element_to_post(p_element_info_map, p_ctx_map, p_on_complete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER', 'background_page.add_element_to_post()');

	const element_type_str      = p_element_info_map['type_str'];
	const selected_elements_map = p_ctx_map['selected_elements_map'];
	const add_datetime_str      = Date.now();

	switch(element_type_str) {
		//----------------
		// IMAGE

		case 'image':
			const url_str = p_element_info_map['full_img_src_str'];

			// check if an element with URL has already been added
			if (url_str in selected_elements_map) {
				var status_str = 'exists';
				p_on_complete_fun(status_str);
			}
			//if not then add it
			else {
				p_element_info_map['add_datetime_str'] = add_datetime_str;
				selected_elements_map[url_str]         = p_element_info_map;

				var status_str = 'success';
				p_on_complete_fun(status_str);
			}
			break;

		//----------------
		// VIDEOS

		case 'video':
			p_element_info_map['add_datetime_str'] = add_datetime_str;
			break;

		//----------------
		default:
			false;
	}
}

//---------------------------------------------------
function remove_element_from_post(p_element_info_map, p_ctx_map, p_log_fun) {

	const element_type_str      = p_element_info_map['type_str'];
	const selected_elements_map = p_ctx_map['selected_elements_map'];

	switch(element_type_str) {
		case 'image':
			const url_str = p_element_info_map['full_img_src_str'];
			delete selected_elements_map[url_str];
			break;
		case 'video':
			true;
			break;
		default:
			false;
	}
}

//-------------------------------------------------