main(log_fun);
//-------------------------------------------------
function log_fun(p_g,p_m) {
	const msg_str = p_g+':'+p_m;
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
function main(p_log_fun) {
	p_log_fun('FUN_ENTER','background_page.main()');
	
	const ctx_map = {
		'selected_elements_map':{},
		//'selected_images_lst'   :[],
		//'selected_videos_lst'   :[]
	};

	chrome.extension.onRequest.addListener(on_request_received_fun);
	//---------------------------------------------------	
	function on_request_received_fun(p_request, 
								p_sender, 
								p_send_response_fun) {		

		p_log_fun('INFO','background_page MSG RECEIVED ------------');
		
		/*const msg_source_str = p_request.source_str;
		const msg_type_str   = p_request.type_str;
		p_log_fun('INFO','msg_source_str - '+msg_source_str);
		p_log_fun('INFO','msg_type_str   - '+msg_type_str);*/

		switch (p_request.source_str) {
			case 'popup':
				handle_popup_msg(p_request.type_str,
							p_request);
				break;
			case 'popup_selected_elements':
				handle_popup_selected_elements(p_request.type_str,
										p_log_fun);
				break;
			case 'content_script':
				handle_content_script_msg(p_request.type_str,
									p_send_response_fun,
									p_request);
				break;
		}
		//---------------------------------------------------	
		function handle_popup_selected_elements(p_msg_type_str,
											p_request) {
			switch(p_msg_type_str) {
				//----------------
				//GET SELECTED ASSETS

				case 'get_selected_elements':
					get__selected_elements(ctx_map,
							(p_selected_elements_map) => {

								console.log(p_selected_elements_map);

								const msg_map = {
									'selected_elements_map':p_selected_elements_map
								};

								p_send_response_fun(msg_map);
							},
							p_log_fun);
					break;
				//----------------
				default:
					p_log_fun('INFO','--------------------------------------');
					p_log_fun('INFO','background_page received unknonwn SELECTED_ASSETS_UI msg');
					p_log_fun('INFO',p_request);
					p_log_fun('INFO',"p_request['type_str']:"+p_request['type_str']);
					break;
			}
		}
		//---------------------------------------------------	
		function handle_popup_msg(p_msg_type_str,
							p_request) {
			switch (p_msg_type_str) {
				//----------------
				//LOG_MSG
				case 'log_msg':
					console.log('POPUP:'+p_request.msg_str);
					break;
				//----------------
				//GET_SELECTED_ELEMENTS
				case 'get__selected_elements':
					get__selected_elements(ctx_map,
						(p_selected_elements_map)=>{
							const msg_map = {
								'selected_elements_map':p_selected_elements_map
							};
							p_send_response_fun(msg_map);
						},
						p_log_fun);
				//----------------
				case 'clear__selected_elements':
					clear__selected_elements(ctx_map,
									()=>{
										const msg_map = {};
										p_send_response_fun(msg_map);
									},
									p_log_fun);
				//----------------
				default:
					p_log_fun('INFO','--------------------------------------');
					p_log_fun('INFO','background_page received unknonwn POPUP msg');
					p_log_fun('INFO',p_request);
					p_log_fun('INFO',"p_request['type_str']:"+p_request['type_str']);
					break;
			}
		}
		//---------------------------------------------------
		function handle_content_script_msg(p_msg_type_str,
									p_send_response_fun,
									p_request) {
			switch (p_msg_type_str) {
				//----------------
				//LOG_MSG
				case 'log_msg':
					console.log('CONTENT_SCR:'+p_request.msg_str);
					break;
				//----------------
				//ADD_ELEMENT_TO_POST
				case 'add_element_to_post':
					var element_info_map = p_request['element_info_map'];

					add_element_to_post(element_info_map,
							ctx_map,
							(p_status_str) => {
								const msg_map = {
									'status_str':p_status_str
								};
								p_send_response_fun(msg_map);
							},
							p_log_fun);
					break;
				//----------------
				//REMOVE_ELEMENT_FROM_POST
				case 'remove_element_from_post':
					var element_info_map = p_request['element_info_map'];
					remove_element_from_post(element_info_map,
										ctx_map,
										p_log_fun);
					break;
				//----------------
				default:
					p_log_fun('INFO','--------------------------------------');
					p_log_fun('INFO','background_page received unknown CONTENT_SCRIPT msg');
					p_log_fun('INFO',conso);
					p_log_fun('INFO',"p_request['type_str']:"+p_request['type_str']);
					break;
			}
		}
		//---------------------------------------------------	
	}
	//---------------------------------------------------
}
//---------------------------------------------------
//POST_OPS
//---------------------------------------------------
function clear__selected_elements(p_ctx_map,
					p_on_complete_fun,
					p_log_fun) {
	p_log_fun('FUN_ENTER','background_page.clear__selected_elements()');

	//IMPORTANT!! - clear all currently selected elements
	p_ctx_map['selected_elements_map'] = {}; 
}
//---------------------------------------------------
function get__selected_elements(p_ctx_map,
					p_on_complete_fun,
					p_log_fun) {
	p_log_fun('FUN_ENTER','background_page.get__selected_elements()');

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
		'images_lst':images_lst,
		'videos_lst':videos_lst
	};

	p_on_complete_fun(selected_elements_map);
}
//---------------------------------------------------
function add_element_to_post(p_element_info_map,
						p_ctx_map,
						p_on_complete_fun,
						p_log_fun) {
	p_log_fun('FUN_ENTER','background_page.add_element_to_post()');

	const element_type_str      = p_element_info_map['type_str'];
	const selected_elements_map = p_ctx_map['selected_elements_map'];
	const add_datetime_str      = Date.now();

	switch(element_type_str) {
		//----------------
		//IMAGE

		case 'image':
			const url_str = p_element_info_map['full_img_src_str'];

			//check if an element with URL has already been added
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
		//VIDEOS

		case 'video':
			p_element_info_map['add_datetime_str'] = add_datetime_str;
			break;
		//----------------
		default:
			false;
	}
}
//---------------------------------------------------
function remove_element_from_post(p_element_info_map,
							p_ctx_map,
							p_log_fun) {
	p_log_fun('FUN_ENTER','background_page.remove_element_from_post()');

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