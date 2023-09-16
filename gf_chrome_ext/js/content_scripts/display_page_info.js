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
function display_page_info(p_page_images_infos_lst,
	p_page_videos_infos_lst,
	p_log_fun) {

	const gf_host_str = "https://gloflow.com";

	const gf_container = $(
		`<div id="page_info_container">
			<div id="parameters">
				<input id="gf_host" value="https://gloflow.com"></input>
			</div>
			<div class="flow_name_field">
				<p class="flow_name_msg"># Add flow names</p>
				<input type="text" class="flow_name" placeholder="general"></input>
			</div>
			<div id="collection_masonry">
			</div>
			<div id="selected_elements_preview"></div>
		</div>`);
		
	$("body").append(gf_container);

	$("body").css({
		"overflow":"hidden",
		"overflow-y":"hidden",
		"overflow-x":"hidden"
	});

	// $(gf_container).css('height', $(document).height());
    var current_scroll_y = window.scrollY;

	const window_height_int = window.innerHeight;
	$(gf_container).css({

		/*
		IMPORTANT!! - setting the height of the container to the height of the window.
			this solves a problem where on some sites scrolling wouldnt work properly with the
			height of the container set to 100% (with intention to cover the whole page).
		*/
		"height": `${window_height_int}px`,
		
		"overflow-y": "visible",
		"top": `${current_scroll_y}px`,
	})

	create_close_btn();

	// IMAGES
	p_page_images_infos_lst.forEach(p_image_map => {
		view_image(p_image_map);
	});

	// VIDEOS
	p_page_videos_infos_lst.forEach(p_video_map => {
		view_video(p_video_map);
	});

	//------------
	// MASONRY

	$(gf_container).find("#collection_masonry").masonry({
			columnWidth:  20,
			gutter:       10,
			itemSelector: ".image_in_page"
		});

	//------------
	// CHECK_IMAGES_EXIST
	check_images_exist_in_system(p_page_images_infos_lst, gf_host_str, p_log_fun);

	//------------

	//---------------------------------------------------
	function create_close_btn() {

		const close_btn_element = $(`<div id="close_btn"></div>`);
		$(gf_container).append(close_btn_element);

		//--------
		// CSS
		const icons_chrome_ext_url_str = 'url('+chrome.runtime.getURL('assets/icons.png')+')';
		$(close_btn_element).css('background-image', icons_chrome_ext_url_str);
		//--------
		
		$(document).on('click', '#close_btn',
			() => {
				$(gf_container).remove();
				$("body").css({
					"overflow": "overlay",
					"overflow-y":"visible",
					"overflow-x":"hidden"
				});
			});
	}

	//---------------------------------------------------
	function view_image(p_image_map) {

		const full_img_src_str = p_image_map['full_img_src_str'];
		const img_name_str     = p_image_map['img_name_str'];

		const image_container_element = $(`
			<div class="image_in_page">
				<img src="${full_img_src_str}"></img>
				<div class="tags"></div>
			</div>
		`);
		
		const img = $(image_container_element).find("img")[0];
		//-----------------
		// GIF
		if (full_img_src_str.split('.').pop() == 'gif') {
			// $(img).addClass('gf_gif');
			// $(img).attr("data-playon","hover"); //GIF_PLAYER API
			// $(img).gifplayer();
		}

		//-----------------

		

		$(gf_container).find('#collection_masonry').append(image_container_element);
			$(img).load(() => {

				//-------------------
				// UNIQUE_IMAGE_DIV_ID - every image is assigned a unique div ID, so that it can be reference without conflicts.
				/*
				IMPORTANT!! - appending time_in_ms to string that is then hashed to obtain a unique div ID, because full_img_src_str
					is not guaranteed to be unique to an image on the page. some pages on the web have the same image displayed
					in multiple places on the same page. this would lead to several images having the same ID, and GF chrome operations
					pontetially then applied multiple times. if an image that has duplicates on the page like this is selected
					by the user to be added to a flow, it would then be added multiple times leading to a duplicate in the GF system.
				*/
				const time_in_ms_f = performance.now();
				const hash_str     = hash_code(`${time_in_ms_f}_${full_img_src_str}`, p_log_fun);
				const img_element_id_str = `id_${hash_str}`;
				
				const img_element_id_clean_str = img_element_id_str.replace('-', '_');
				p_log_fun('INFO', `img_element_id_clean_str - ${img_element_id_clean_str}`);

				//-------------------

				init_image_hud(img_element_id_clean_str,
					img,
					image_container_element,
					p_image_map,
					gf_container,
					gf_host_str,
					p_log_fun);

				//-------------------
				// IMPORTANT!! - reload the masonry layout with the newly loaded image
				$(gf_container).find('#collection_masonry').masonry();
				
				//-------------------
			});

		//---------------------------------------------------
	}

	//---------------------------------------------------
	function view_video(p_video_map) {

		const video_source_str = p_video_map['video_source_str'];
		const embed_url_str    = p_video_map['embed_url_str'];

		const video_in_page_element = $(
			'<div class="video_in_page">'+
				'<iframe src="'+embed_url_str+'"></iframe>'+
			'</div>');

		$(gf_container).find('#collection_masonry').append(video_in_page_element);

		init_video_hud(video_in_page_element, p_video_map, p_log_fun);

		//-------------------
		// IMPORTANT!! - reload the masonry layout with the newly loaded image
		$(gf_container).find('#collection_masonry').masonry();

		//-------------------
	}

	//---------------------------------------------------
}

//---------------------------------------------------
// CHECK_IMAGES_EXIST_IN_SYSTEM

function check_images_exist_in_system(p_page_images_infos_lst,
	p_gf_host_str,
	p_log_fun) {
	
	const images_extern_urls_lst = []; // :List<:String>
	p_page_images_infos_lst.forEach(p_image_map => {

			const full_img_src_str = p_image_map['full_img_src_str'];
			images_extern_urls_lst.push(full_img_src_str);
		});

	console.log(images_extern_urls_lst)

	const msg_map = {
		"source_str": "content_script",
		"type_str":   "check_images_exist",
		"images_extern_urls_lst": images_extern_urls_lst,
		"gf_host_str":            p_gf_host_str
	}

	send_msg_to_bg_page(msg_map, (p_response_map)=>{

		switch(p_response_map["status_str"]) {
			case "OK":

				const existing_images_lst = p_response_map["existing_images_lst"];

				$.each(existing_images_lst, (p_i, p_e)=>{

					const img__id_str               = p_e['id_str'];
					const img__origin_url_str       = p_e['origin_url_str'];
					const img__origin_page_url_str  = p_e['origin_page_url_str'];
					const img__creation_unix_time_f = p_e['creation_unix_time_f'];
					const img__flows_names_lst      = p_e["flows_names_lst"];

					existing_img__update_view(img__id_str,
						img__origin_url_str,
						img__origin_page_url_str,
						img__creation_unix_time_f,
						img__flows_names_lst);
				});

				break;

			case "ERROR":
				p_on_error_fun(p_response_map["data_map"]);
				break;
		}
	});

	//---------------------------------------------------
	function existing_img__update_view(p_img__system_id_str,
		p_existing_img__origin_url_str,
		p_existing_img__origin_page_url_str,
		p_existing_img__creation_unix_time_f,
		p_img__flows_names_lst) {

		const date                 = new Date(p_existing_img__creation_unix_time_f*1000);
		const data_str             = `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()} - ${date.getDate()}.${date.getMonth()}.${date.getFullYear()}`;
		const existing_img_preview = $('#page_info_container').find(`img[src="${p_existing_img__origin_url_str}"]`)[0];
		
		const flows_links_str = p_img__flows_names_lst
			.map((p_name_str)=>`<a class="flow_name" href="${p_gf_host_str}/images/flows/browser?fname=${p_name_str}" target="_blank">${p_name_str}</a>`)
			.join(",");

		const parent_element = $(existing_img_preview).parent();
		$(parent_element).append(
			`<div class="img_exists">
				<div class="exists_msg">exists in flows: <span>${flows_links_str}</span></div>
				<div class="origin_page_url"><span>origin page url</span>: <a href="${p_existing_img__origin_page_url_str}">${p_existing_img__origin_page_url_str}</div>
				<div class="creation_time">
					<span class="msg">created on:</span>
					<span class="time">${data_str}</span>
				</div>
			</div>`);

		//-------------------
		/*
		IMPORTANT!! - mark the image as existing, in the "data" property of its container div.
			this is read by other parts of the code, one which is the tagging UI control that checks if the image exists in the system,
			before the tag can be added to it.
		*/
		$(parent_element).attr("data-img_exists_in_flow_bool", true);

		// this is the GF_ID of the image
		$(parent_element).attr("data-img_system_id_str", p_img__system_id_str);

		//-------------------
	}
	//---------------------------------------------------
}

//---------------------------------------------------
// HUD
//---------------------------------------------------
function init_image_hud(p_image_element_id_str,
	p_image,
	p_image_container_element,
	p_image_info_map,
	p_gf_container_element,
	p_gf_host_str,
	p_log_fun) {

	const full_img_src_str = p_image_info_map['full_img_src_str'];

	const img_height = p_image_info_map['img_height'];
	const img_width  = p_image_info_map['img_width'];
	const hud        = $(
		`<div id="${p_image_element_id_str}" class="hud">
			<div class="background"></div>
			<div class="full_img_src">${full_img_src_str}</div>
			<div id="actions">
				<!-- ADD_TO_POST BTN -->
				<div class="add_to_post_btn" title="add image to post">
					<div class="symbol">
						<div class="icon"></div>
					</div>
				</div>
			</div>
		</div>`);

	img_ext_str = full_img_src_str.split('.').pop();
	//--------------	
	if (img_ext_str == 'gif') {
		//-------------------
		// ADD_TO_GIF_BTN
		const gif_btn = $(
			`<div id="${p_image_element_id_str}" class="add_to_gif_flow_btn" title="add gif to gif-flow">
				<div class="symbol">
					<div class="icon">GIF</div>
				</div>
			</div>`);

		//-------------------
		$(hud).find('#actions').append(gif_btn);
	} else {
		//-------------------
		// ADD_TO_IMAGE_FLOW BTN
		const flow_btn = $(
				`<div class="add_to_image_flow_btn" title="add image to image-flow">
					<div class="symbol">
						<div class="icon"></div>
					</div>
				</div>`);

		//-------------------
		$(hud).find('#actions').append(flow_btn);
	}

	//--------------
	// IMPORTANT!! - testing for image dimensions, so that this info is not displayed if
	//               the image is too small, since it will obstruct actions div.
	if (img_height > 120) {
		$(hud).append(`<div class="img_height">height: <span>${img_height}</span>px</div>`);
		$(hud).append(`<div class="img_width">width: <span>${img_width}</span>px</div>`);	
	}
	
	//--------------

	const icons_chrome_ext_url_str = `url(${chrome.runtime.getURL('assets/icons.png')})`;
	$(hud).find('.add_to_image_flow_btn .icon').css('background-image', icons_chrome_ext_url_str);
	$(hud).find('.add_to_post_btn .icon').css('background-image', icons_chrome_ext_url_str);
	//------------
	// ADD_TO_GIF_FLOW

	/*
	const add_to_gif_flow__selector_str = '#'+p_image_element_id_str+' .add_to_gif_flow_btn';
	$(document).on('click', add_to_gif_flow__selector_str, ()=>{

		const image_flows_names_lst = ["general", "gifs"];
		const gf_host_str           = $(p_gf_container_element).find("input#gf_host").val();

		add_image_to_flow(full_img_src_str,
			image_flows_names_lst,
			gf_host_str,
			()=>{

				//-------------------
				// IMPORTANT!! - adding the .btn_ok class activates the CSS animation
				$(hud).find('.add_to_image_flow_btn .icon').addClass('btn_ok');

				//-------------------
				$(hud).find('.add_to_image_flow_btn').css('pointer-events', 'none');
			},
			(p_error_data_map)=>{},
			p_log_fun);
	});
	*/

	//------------
	// ADD_TO_IMAGE_FLOW
	const add_to_image_flow__selector_str = `#${p_image_element_id_str} .add_to_image_flow_btn`;
	$(document).on("click", add_to_image_flow__selector_str, async (p_event)=>{
		p_event.stopImmediatePropagation();

		await add_to_image_flow_btn_handler();
	});
	
	var image_added_to_flow_bool = false;

	//---------------------------------------------------
	async function add_to_image_flow_btn_handler() {
		const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

			const flows_names_str 		= $("input.flow_name").val();
			const final_flows_names_lst = []

			if(flows_names_str == ""){
				final_flows_names_lst.push("general");
			}else{
				const flows_names_lowercased_str = flows_names_str.toLowerCase()
				const flows_names_lst 			 = flows_names_lowercased_str.split(" ");
				const flows_names_filtered_lst   = flows_names_lst.filter(n => n) // removes empty strings from array
				final_flows_names_lst.push(...flows_names_filtered_lst);
			}

			const gf_host_str = $(p_gf_container_element).find("input#gf_host").val();

			//-------------------
			// DISABLE_EVENTS
			// remove the event handler so that clicking on the add_to_flow button no longer functions,
			// since the image is in the process of been added and shouldnt be add-able anymore.
			$(document).off("click", add_to_image_flow__selector_str, add_to_image_flow_btn_handler);

			//-------------------
			// DISPLAY_MESSAGE - to let the user know that adding the image to flow is in progress
			$(add_to_image_flow__selector_str).append("<div class='adding_in_progress'>adding to flow in progress...</div>");

			//-------------------

			add_image_to_flow(full_img_src_str,
				final_flows_names_lst,
				gf_host_str,

				//---------------------------------------------------
				// on_complete_fun
				(p_image_id_str)=>{
					//-------------------
					// IMPORTANT!! - adding the .btn_ok class activates the CSS animation
					$(hud).find(".add_to_image_flow_btn .icon").addClass("btn_ok");

					//-------------------

					$(hud).find(".add_to_image_flow_btn").css("pointer-events", "none");

					// REMOVE_MESSAGE
					$(add_to_image_flow__selector_str).find(".adding_in_progress").remove();

					image_added_to_flow_bool = true;

					p_resolve_fun(p_image_id_str);
				},

				//---------------------------------------------------
				// on_error_fun
				(p_data_map)=>{

					//-------------------
					// ENABLE_EVENTS
					// adding an image failed for some reason, so allow the user to attempt to add it again.
					$(document).on("click", add_to_image_flow__selector_str, add_to_image_flow_btn_handler);

					//-------------------

					// REMOVE_MESSAGE
					$(add_to_image_flow__selector_str).find(".adding_in_progress").remove();

					p_reject_fun();
				},

				//---------------------------------------------------
				p_log_fun);
		});
		return p;
	}

	//---------------------------------------------------
	//------------

	//------------
	// ADD_TO_POST
	const add_to_post__selector_str = `#${p_image_element_id_str} .add_to_post_btn`;
	$(document).on('click', add_to_post__selector_str, (p_event)=>{
		p_event.stopImmediatePropagation();

		add_image_to_post(p_image_info_map, p_log_fun);
	});

	//------------
	// GIF
	if (full_img_src_str.split('.').pop() == 'gif') {

		//$(p_image_container_element).find('.gf_gif').gifplayer();

		/*
		const gif_element = $(p_image_container_element).find('img')[0];
		console.log(gif_element)
		var sup1 = new SuperGif({ gif: gif_element } );
		sup1.load();
		sup1.play();
		*/
	}

	//------------

	var hud_attached_bool = false;
	$(p_image_container_element).mouseenter((p_event)=>{
		// p_event.stopImmediatePropagation();

		if (!hud_attached_bool) {
			$(p_image_container_element).append(hud);
			hud_attached_bool = true;
		}
		else {
			$(hud).css('visibility', 'visible');
		}
	});
	$(p_image_container_element).mouseleave((p_event)=>{
		// p_event.stopImmediatePropagation();

		$(hud).css('visibility', 'hidden');
	});

	//------------
	// TAGGER
	/*
	IMPORTANT!! - need to set the width of the tags container to be the same as the image to which
		it is attached, so that as tags are added the tags container scales in height, not in width
		beyond the width of the image to which it is attached.
	*/
	$(p_image_container_element).find(".tags")[0].style.width = `${p_image.width}px`;


	init_tagging(p_image_container_element,
		p_gf_container_element,

		// tags_create_pre_fun
		// called before a tag is about to be added to an image
		async ()=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				var img_system_id_str;

				/*
				IMPORTANT!! - check that image parent contains the data property first, which indicates
					if the image already exists in the system or not.
					this is being checked here via the data property, because initially image_added_to_flow_bool
					is only being set to true when user manually adds image to flow, not if the image is 
					previously existing in the system.
				*/
				if (!image_added_to_flow_bool) {

					if ($(p_image_container_element).attr("data-img_exists_in_flow_bool") !== undefined) {
						const img_exists_in_flow_str = $(p_image_container_element).attr("data-img_exists_in_flow_bool");
						if (img_exists_in_flow_str === "true") {

							console.log("tags_create_pre_fun - image already exists in the system from past additions...")
							image_added_to_flow_bool = true;


							img_system_id_str = $(p_image_container_element).attr("data-img_system_id_str");
							
						}
					}
				}

				// definitive check if the image exists in the flow, after all possible updates to image_added_to_flow_bool value
				if (!image_added_to_flow_bool) {

					console.log("tags_create_pre_fun - image definitelly doesnt exist in the system...");

					// if image doesnt exist in the flow, first add it.
					img_system_id_str = await add_to_image_flow_btn_handler();
				}

				// image exists in the flow, so just get its data system_id
				else {
					img_system_id_str = $(p_image_container_element).attr("data-img_system_id_str");
				}

				p_resolve_fun(img_system_id_str);
			});
			return p;
		},
		p_gf_host_str,
		p_log_fun);

	//------------
}

//---------------------------------------------------
function init_video_hud(p_video_in_page_element, p_video_info_map, p_log_fun) {

	const hud = $(
		`<div class="hud">
			<div class="add_to_post_btn">
				<div class="symbol">
					<div class="icon"></div>
				</div>
			</div>
		</div>`);
	const icons_chrome_ext_url_str = `url(${chrome.runtime.getURL('assets/icons.png')})`;
	$(hud).find('.icon').css('background-image', icons_chrome_ext_url_str);

	$(p_video_in_page_element).mouseenter(function(p_e) {
		$(p_video_in_page_element).append(hud);
	});

	$(p_video_in_page_element).mouseleave(function(p_e) {
		$(hud).remove();
	});
}

//---------------------------------------------------
// ADD_IMAGE_TO_FLOW

function add_image_to_flow(p_full_img_src_str,
	p_images_flows_names_lst,
	p_gf_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {

	const image_origin_page_url_str = window.location.href;
	const msg_map = {
		"source_str":                "content_script",
		"type_str":                  "add_image_to_flow",
		"full_img_src_str":          p_full_img_src_str,
		"image_origin_page_url_str": image_origin_page_url_str,
		"images_flows_names_lst":    p_images_flows_names_lst,
		"gf_host_str":               p_gf_host_str
	};

	send_msg_to_bg_page(msg_map, (p_response_map)=>{

		switch(p_response_map["status_str"]) {
			case "OK":

				const image_id_str = p_response_map["image_id_str"];
				p_on_complete_fun(image_id_str);
				break;

			case "ERROR":
				p_on_error_fun(p_response_map["data_map"]);
				break;
		}
	});
}

//---------------------------------------------------
// TAGGING_UI

function init_tagging(p_image_container_element,
	p_gf_container,
	p_tags_create_pre_fun,
	p_gf_host_str,
	p_log_fun) {
	
	var image_system_id_str;

	const http_api_map = {
		"gf_tagger": {
			"add_tags_to_obj": (p_new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				p_tags_meta_map,
				p_log_fun)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
					

					add_tags_to_image(p_new_tags_lst,
						image_system_id_str,
						p_gf_host_str,
						p_log_fun);

					p_resolve_fun({
						"added_tags_lst": p_new_tags_lst,
					});
				});
				return p;
			}
		}
	};

	const obj_type_str = "image";
	const input_element_parent_selector_str = "#page_info_container";

	gf_tagger__init_ui(obj_type_str,
		p_image_container_element,
		input_element_parent_selector_str,

		//---------------------------------------------------
		// tags_create_pre_fun
		async (p_tags_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// p_tags_create_pre_fun resolves the system_id of the item being tagged
				image_system_id_str = await p_tags_create_pre_fun(p_tags_lst);

				p_resolve_fun();
			});
			return p;
		},

		//---------------------------------------------------
		// on_tags_created_fun
		(p_tags_lst)=>{

			console.log("added tags >>>>>>>>>>>", p_tags_lst)
			console.log("DDDDDDDDDDDDDDDDDD", p_image_container_element, "ggg", $(p_image_container_element).find(".tags"))

			p_tags_lst.forEach(t_str=>{

				const tag_link_url_str = `${p_gf_host_str}/v1/tags/objects?tag=${t_str}&otype=image`

				const element = $(`
					<div class='bubble-in auto-width tag'>
						<a href="${tag_link_url_str}" target="_blank" style="text-decoration: none;color: inherit;">
							${t_str}
						</a>
					</div>`);
				$(p_image_container_element).find(".tags").append(element);

				// start the css animation
				element.addClass('animate');

				// IMPORTANT!! - reload the masonry layout with the newly added tag
				$(p_gf_container).find('#collection_masonry').masonry();
			})
		},

		//---------------------------------------------------
		()=>{}, // on_tag_ui_add_fun
		()=>{}, // on_tag_ui_remove_fun
		http_api_map,
		p_log_fun);
}

//---------------------------------------------------
function view_gif_info(p_full_img_src_str, p_host_str, p_log_fun) {
	p_log_fun("FUN_ENTER", "display_page_info.view_gif_info()");

	http__gif_get_info(p_full_img_src_str,
		p_host_str,
		(p_gif_map)=>{

			const preview_frames_s3_urls_lst = p_gif_map["preview_frames_s3_urls_lst"];

			console.log("GIF FRAMES URLS ------------------");
			console.log(preview_frames_s3_urls_lst);
		},
		(p_error_data_map)=>{},
		p_log_fun);
}

//---------------------------------------------------
// ADD_TAGS_TO_IMAGE

function add_tags_to_image(p_tags_lst,
	p_image_system_id_str,
	p_gf_host_str,
	p_log_fun) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {
		const image_origin_page_url_str = window.location.href;
		const msg_map = {
			"source_str": "content_script",
			"type_str":   "add_tags_to_image",
			"image_system_id_str": p_image_system_id_str,
			"tags_lst":            p_tags_lst,
			"gf_host_str":         p_gf_host_str
		};

		send_msg_to_bg_page(msg_map, (p_response_map)=>{

			switch(p_response_map["status_str"]) {
				case "OK":

					p_resolve_fun();
					break;

				case "ERROR":
					p_reject_fun(p_response_map["data_map"]);
					break;
			}
		});
	});
	return p;
}