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
	$.each(p_page_images_infos_lst,
		(p_i, p_image_map) => {
			view_image(p_image_map);
		});

	// VIDEOS
	$.each(p_page_videos_infos_lst,
		(p_i, p_video_map) => {
			view_video(p_video_map);
		});

	//------------
	// MASONRY

	$(gf_container).find("#collection_masonry").masonry(
		{
			columnWidth:  20,
			gutter:       10,
			itemSelector: ".image_in_page"
		});

	//------------
	// CHECK_IMAGES_EXIST
	check_images_exist_in_system(p_page_images_infos_lst, p_log_fun);

	//------------
	// $(document).resize(function() {
	// 	$('#page_info_gf_container').css('width',$(document).width());
	//	$('#page_info_gf_container').css('height',$(document).height());
	// });

	//---------------------------------------------------
	function create_close_btn() {
		
		p_log_fun("FUN_ENTER", "display_page_info.display_page_info().create_close_btn()");

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

		const image_in_page_element = $(`
			<div class="image_in_page">
				
			</div>
		`);
		
		const img = $('<img></img>').attr('src', full_img_src_str);

		$(image_in_page_element).append(img);
		//-----------------
		// GIF
		if (full_img_src_str.split('.').pop() == 'gif') {
			// $(img).addClass('gf_gif');
			// $(img).attr("data-playon","hover"); //GIF_PLAYER API
			// $(img).gifplayer();
		}

		//-----------------

		$(gf_container).find('#collection_masonry').append(image_in_page_element);
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
				const img_id_str = 'id_'+hash_code(`${time_in_ms_f}_${full_img_src_str}`, p_log_fun);
				
				const img_id_clean_str = img_id_str.replace('-', '_');
				p_log_fun('INFO', 'img_id_clean_str - '+img_id_clean_str);

				init_image_hud(img_id_clean_str,
					image_in_page_element,
					p_image_map,
					gf_container,
					p_log_fun);

				//-------------------
				// IMPORTANT!! - reload the masonry layout with the newly loaded image
				$(gf_container).find('#collection_masonry').masonry();
				
				//-------------------
			});
	}

	//---------------------------------------------------
	function view_video(p_video_map) {

		const video_source_str = p_video_map['video_source_str'];
		const embed_url_str    = p_video_map['embed_url_str'];

		const video_in_page_element = $(
			'<div class="video_in_page">'+
				'<iframe src="'+embed_url_str+'"></iframe>'+
			'</div>');

		$(container).find('#collection_masonry').append(video_in_page_element);

		init_video_hud(video_in_page_element, p_video_map, p_log_fun);

		//-------------------
		// IMPORTANT!! - reload the masonry layout with the newly loaded image
		$(container).find('#collection_masonry').masonry();

		//-------------------
	}

	//---------------------------------------------------
}

//---------------------------------------------------
// CHECK_IMAGES_EXIST_IN_SYSTEM

function check_images_exist_in_system(p_page_images_infos_lst, p_log_fun) {

	const gf_host_str = "https://gloflow.com"
	const images_extern_urls_lst = []; // :List<:String>
	$.each(p_page_images_infos_lst,
		(p_i, p_image_map) => {
			const full_img_src_str = p_image_map['full_img_src_str'];
			images_extern_urls_lst.push(full_img_src_str);
		});

	//-------------------
	// IMPORTANT!! - since request to host_str is made from the context of the page in which the 
	//               content is located, browser security imposes that the same protocol (http|https)
	//               is used to communicate with host_str as with the origin-domain of the page
	const origin_url_str = window.location.href;
	const protocol_str   = origin_url_str.split('://')[0];
	const host_str       = `${protocol_str}://gloflow.com`;
	
	//-------------------

	console.log(images_extern_urls_lst)

	const msg_map = {
		"source_str": "content_script",
		"type_str":   "check_images_exist",
		"images_extern_urls_lst": images_extern_urls_lst,
		"gf_host_str":            host_str
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

					existing_img__update_view(img__origin_url_str,
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
	function existing_img__update_view(p_existing_img__origin_url_str,
		p_existing_img__origin_page_url_str,
		p_existing_img__creation_unix_time_f,
		p_img__flows_names_lst) {

		const date                 = new Date(p_existing_img__creation_unix_time_f*1000);
		const data_str             = `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()} - ${date.getDate()}.${date.getMonth()}.${date.getFullYear()}`;
		const existing_img_preview = $('#page_info_container').find(`img[src="${p_existing_img__origin_url_str}"]`)[0];
		
		const flows_links_str = p_img__flows_names_lst
			.map((p_name_str)=>`<a class="flow_name" href="${gf_host_str}/images/flows/browser?fname=${p_name_str}" target="_blank">${p_name_str}</a>`)
			.join(",");

		$(existing_img_preview).parent().append(
			`<div class="img_exists">
				<div class="exists_msg">exists in flows: <span>${flows_links_str}</span></div>
				<div class="origin_page_url"><span>origin page url</span>: <a href="${p_existing_img__origin_page_url_str}">${p_existing_img__origin_page_url_str}</div>
				<div class="creation_time">
					<span class="msg">created on:</span>
					<span class="time">${data_str}</span>
				</div>
			</div>`);
	}
	//---------------------------------------------------
}

//---------------------------------------------------
// HUD
//---------------------------------------------------
function init_image_hud(p_image_id_str,
	p_image_in_page_element,
	p_image_info_map,
	p_gf_container_element,
	p_log_fun) {

	const full_img_src_str = p_image_info_map['full_img_src_str'];

	const img_height = p_image_info_map['img_height'];
	const img_width  = p_image_info_map['img_width'];
	const hud        = $(
		`<div id="${p_image_id_str}" class="hud">
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
			`<div id="${p_image_id_str}" class="add_to_gif_flow_btn" title="add gif to gif-flow">
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
	const add_to_gif_flow__selector_str = '#'+p_image_id_str+' .add_to_gif_flow_btn';
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
	const add_to_image_flow__selector_str = `#${p_image_id_str} .add_to_image_flow_btn`;
	$(document).on("click", add_to_image_flow__selector_str, add_to_image_flow_btn_handler);
	
	//---------------------------------------------------
	function add_to_image_flow_btn_handler() {
		
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

			()=>{
				//-------------------
				// IMPORTANT!! - adding the .btn_ok class activates the CSS animation
				$(hud).find(".add_to_image_flow_btn .icon").addClass("btn_ok");

				//-------------------

				$(hud).find(".add_to_image_flow_btn").css("pointer-events", "none");

				// REMOVE_MESSAGE
				$(add_to_image_flow__selector_str).find(".adding_in_progress").remove();
			},
			(p_data_map)=>{

				//-------------------
				// ENABLE_EVENTS
				// adding an image failed for some reason, so allow the user to attempt to add it again.
				$(document).on("click", add_to_image_flow__selector_str, add_to_image_flow_btn_handler);

				//-------------------

				// REMOVE_MESSAGE
				$(add_to_image_flow__selector_str).find(".adding_in_progress").remove();
			},
			p_log_fun);
	}

	//---------------------------------------------------
	//------------

	//------------
	// ADD_TO_POST
	const add_to_post__selector_str = `#${p_image_id_str} .add_to_post_btn`;
	$(document).on('click', add_to_post__selector_str, ()=>{
		add_image_to_post(p_image_info_map, p_log_fun);
	});

	//------------
	// GIF
	if (full_img_src_str.split('.').pop() == 'gif') {

		//$(p_image_in_page_element).find('.gf_gif').gifplayer();

		/*
		const gif_element = $(p_image_in_page_element).find('img')[0];
		console.log(gif_element)
		var sup1 = new SuperGif({ gif: gif_element } );
		sup1.load();
		sup1.play();
		*/
	}

	//------------

	var hud_attached_bool = false;
	$(p_image_in_page_element).mouseenter((p_e)=>{

		if (!hud_attached_bool) {
			$(p_image_in_page_element).append(hud);
			hud_attached_bool = true;
		}
		else {
			$(hud).css('visibility', 'visible');
		}
	});
	$(p_image_in_page_element).mouseleave((p_e)=>{
		$(hud).css('visibility', 'hidden');
	});
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

				p_on_complete_fun();
				break;

			case "ERROR":
				p_on_error_fun(p_response_map["data_map"]);
				break;
		}
	});
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