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
	p_log_fun('FUN_ENTER', 'display_page_info.display_page_info()');

	const gf_container = $(
		`<div id="page_info_container">
			<div id="parameters">
				<input id="gf_host" value="https://gloflow.com"></input>
			</div>
			<div class="flow_name_field">
				<p class="flow_name_msg"># Add flow names</p>
				<input type="text" class="flow_name" placeholder="general"></input>
			</div>
			<div id="collection"></div>
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

	$(gf_container).css({
		"height":	"100%",
		"overflow-y": "visible",
		"top":	current_scroll_y+"px",
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

	$(gf_container).find("#collection").masonry(
		{
			columnWidth:  20,
			gutter:       10,
			itemSelector: ".image_in_page"
		});

	//------------

	check_images_exist_in_system(p_page_images_infos_lst, p_log_fun);

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
		const icons_chrome_ext_url_str = 'url('+chrome.extension.getURL('assets/icons.png')+')';
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
		//p_log_fun('FUN_ENTER','display_page_info.display_page_info().view_image()');

		const full_img_src_str = p_image_map['full_img_src_str'];
		const img_name_str     = p_image_map['img_name_str'];

		//p_log_fun('INFO','+++-------------------+++');
		//p_log_fun('INFO',JSON.stringify(p_image_map));
		//p_log_fun('INFO',full_img_src_str);
		//p_log_fun('INFO',img_name_str);

		const image_in_page_element = $(`
			<div class="image_in_page">
				
			</div>
		`);
		
		const img = $('<img></img>').attr('src', full_img_src_str);

		$(image_in_page_element).append(img);
		//-----------------
		//GIF
		if (full_img_src_str.split('.').pop() == 'gif') {
			// $(img).addClass('gf_gif');
			// $(img).attr("data-playon","hover"); //GIF_PLAYER API
			// $(img).gifplayer();
		}
		//-----------------

		$(gf_container).find('#collection').append(image_in_page_element);
			$(img).load(() => {
				const img_id_str = 'id_'+hash_code(full_img_src_str, p_log_fun);
				
				const img_id_clean_str = img_id_str.replace('-', '_');
				p_log_fun('INFO', 'img_id_clean_str - '+img_id_clean_str);

				init_image_hud(img_id_clean_str,
					image_in_page_element,
					p_image_map,
					gf_container,
					p_log_fun);

				//-------------------
				// IMPORTANT!! - reload the masonry layout with the newly loaded image
				$(gf_container).find('#collection').masonry();
				//-------------------
			});
	}

	//---------------------------------------------------
	function view_video(p_video_map) {
		//p_log_fun('FUN_ENTER','display_page_info.display_page_info().view_video()');

		const video_source_str = p_video_map['video_source_str'];
		const embed_url_str    = p_video_map['embed_url_str'];

		const video_in_page_element = $(
			'<div class="video_in_page">'+
				'<iframe src="'+embed_url_str+'"></iframe>'+
			'</div>');

		$(container).find('#collection').append(video_in_page_element);

		init_video_hud(video_in_page_element, p_video_map, p_log_fun);

		//-------------------
		//IMPORTANT!! - reload the masonry layout with the newly loaded image
		$(container).find('#collection').masonry();
		//-------------------
	}

	//---------------------------------------------------
}

//---------------------------------------------------
function check_images_exist_in_system(p_page_images_infos_lst, p_log_fun) {
	p_log_fun('FUN_ENTER','display_page_info.check_images_exist_in_system()');

	const images_extern_urls_lst = []; //:List<:String>
	$.each(p_page_images_infos_lst,
		(p_i,p_image_map) => {
			const full_img_src_str = p_image_map['full_img_src_str'];
			images_extern_urls_lst.push(full_img_src_str);
		});

	//-------------------
	//IMPORTANT!! - since request to host_str is made from the context of the page in which the 
	//              content is located, browser security imposes that the same protocol (http|https)
	//              is used to communicate with host_str as with the origin-domain of the page
	const origin_url_str = window.location.href;
	const protocol_str   = origin_url_str.split('://')[0];
	const host_str       = protocol_str+'://gloflow.com'; //127.0.0.1:3050';
	//-------------------

	console.log(images_extern_urls_lst)

	http__check_imgs_exist_in_flow(images_extern_urls_lst,
		host_str,
		(p_existing_images_lst)=>{

			console.log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>");
			console.log(p_existing_images_lst);

			$.each(p_existing_images_lst,(p_i,p_e)=>{

				const img__id_str               = p_e['id_str'];
				const img__origin_url_str       = p_e['origin_url_str'];
				const img__origin_page_url_str  = p_e['origin_page_url_str'];
				const img__creation_unix_time_f = p_e['creation_unix_time_f'];

				existing_img__update_view(img__origin_url_str, img__origin_page_url_str, img__creation_unix_time_f);
			});
		},
		(p_error_data_map)=>{},
		p_log_fun);

	//---------------------------------------------------
	function existing_img__update_view(p_existing_img__origin_url_str,
		p_existing_img__origin_page_url_str,
		p_existing_img__creation_unix_time_f) {
		p_log_fun('FUN_ENTER','display_page_info.check_images_exist_in_system().existing_img__update_view()');

		const date                 = new Date(p_existing_img__creation_unix_time_f*1000);
		const data_str             = date.getHours()+':'+date.getMinutes()+':'+date.getSeconds()+' - '+date.getDate()+'.'+date.getMonth()+'.'+date.getFullYear();
		const existing_img_preview = $('#page_info_container').find('img[src="'+p_existing_img__origin_url_str+'"]')[0];
		$(existing_img_preview).parent().append(
			'<div class="img_exists">'+
				'<div class="exists_msg">added</div>'+
				'<div class="origin_page_url"><span>origin page url</span><a href="'+p_existing_img__origin_page_url_str+'">'+p_existing_img__origin_page_url_str+'</div>'+
				'<div class="creation_time">'+
					'<span class="msg">created on:</span>'+
					'<span class="time">'+data_str+'</span>'+
				'</div>'+
			'</div>');
	}
	//---------------------------------------------------
}

//---------------------------------------------------
//HUD
//---------------------------------------------------
function init_image_hud(p_image_id_str,
	p_image_in_page_element,
	p_image_info_map,
	p_gf_container_element,
	p_log_fun) {
	//p_log_fun('FUN_ENTER','display_page_info.init_image_hud()');

	const full_img_src_str = p_image_info_map['full_img_src_str'];

	const img_height = p_image_info_map['img_height'];
	const img_width  = p_image_info_map['img_width'];
	const hud        = $(
		'<div id="'+p_image_id_str+'"class="hud">'+
			'<div class="background"></div>'+
			'<div class="full_img_src">'+full_img_src_str+'</div>'+
			'<div id="actions">'+
				/*//-------------------
				//ADD_TO_IMAGE_FLOW BTN
				'<div class="add_to_image_flow_btn" title="add image to image-flow">'+
					'<div class="symbol">'+
						'<div class="icon"></div>'+
					'</div>'+
				'</div>'+*/
				//-------------------
				//ADD_TO_POST BTN
				'<div class="add_to_post_btn" title="add image to post">'+
					'<div class="symbol">'+
						'<div class="icon"></div>'+
					'</div>'+
				'</div>'+
				//-------------------
			'</div>'+
		'</div>');

	img_ext_str = full_img_src_str.split('.').pop();
	//--------------	
	if (img_ext_str == 'gif') {
		//-------------------
		// ADD_TO_GIF_BTN
		const gif_btn = $(
				'<div id="'+p_image_id_str+'" class="add_to_gif_flow_btn" title="add gif to gif-flow">'+
					'<div class="symbol">'+
						'<div class="icon">GIF</div>'+
					'</div>'+
				'</div>');
		//-------------------
		$(hud).find('#actions').append(gif_btn);
	} else {
		//-------------------
		// ADD_TO_IMAGE_FLOW BTN
		const flow_btn = $(
				'<div class="add_to_image_flow_btn" title="add image to image-flow">'+
					'<div class="symbol">'+
						'<div class="icon"></div>'+
					'</div>'+
				'</div>');
		//-------------------
		$(hud).find('#actions').append(flow_btn);
	}
	//--------------
	// IMPORTANT!! - testing for image dimensions, so that this info is not displayed if
	//               the image is too small, since it will obstruct actions div.
	if (img_height > 120) {
		$(hud).append('<div class="img_height">height: <span>'+img_height+'</span>px</div>');
		$(hud).append('<div class="img_width">width: <span>'+img_width+'</span>px</div>');	
	}
	//--------------

	const icons_chrome_ext_url_str = 'url('+chrome.extension.getURL('assets/icons.png')+')';
	$(hud).find('.add_to_image_flow_btn .icon').css('background-image',icons_chrome_ext_url_str);
	$(hud).find('.add_to_post_btn .icon').css('background-image',icons_chrome_ext_url_str);
	//------------
	// ADD_TO_GIF_FLOW

	const add_to_gif_flow__selector_str = '#'+p_image_id_str+' .add_to_gif_flow_btn';
	$(document).on('click', add_to_gif_flow__selector_str,()=>{

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
	//------------
	// ADD_TO_IMAGE_FLOW
	const add_to_image_flow__selector_str = '#'+p_image_id_str+' .add_to_image_flow_btn';
	$(document).on("click", add_to_image_flow__selector_str, ()=>{
		
		/*const image_origin_page_url_str = window.location.href;
		//-------------------
		// IMPORTANT!! - since request to host_str is made from the context of the page in which the 
		//               content is located, browser security imposes that the same protocol (http|https)
		//               is used to communicate with host_str as with the origin-domain of the page
		const origin_url_str = window.location.href;
		const protocol_str   = origin_url_str.split('://')[0];
		// const host_str       = protocol_str+'://127.0.0.1:3050';
		const host_str       = protocol_str+'://gloflow.com';
		//-------------------

		http__add_image_to_flow(full_img_src_str,
			image_origin_page_url_str,
			host_str,
			(p_images_job_id_str,
			p_image_id_str,
			p_thumbnail_small_relative_url_str)=>{

				//-------------------
				// IMPORTANT!! - adding the .btn_ok class activates the CSS animation
				$(hud).find('.add_to_image_flow_btn .icon').addClass('btn_ok');
				//-------------------

				$(hud).find('.add_to_image_flow_btn').css('pointer-events','none');
			},
			p_log_fun);*/

		// IN GF EXTENTION ADDED INPUT FIELD THAT DIRECTS IMAGES TO FLOWS VIA STRING
		const flows_names_str 			 = $("input.flow_name").val();
		const flows_names_lowercased_str = flows_names_str.toLowerCase()
		const flows_names_lst 			 = flows_names_lowercased_str.split(" ");
		const flows_names_filtered_lst   = flows_names_lst.filter(n => n) // removes empty strings from array
		const image_flows_names_lst 	 = []

		if(flows_names_str == ""){
			image_flows_names_lst.push("general");
		}else{
			image_flows_names_lst.push(...flows_names_filtered_lst);
		}

		//const image_flows_names_lst = ["general"];
		const gf_host_str           = $(p_gf_container_element).find("input#gf_host").val();
		add_image_to_flow(full_img_src_str,
			image_flows_names_lst,
			gf_host_str,
			()=>{
				//-------------------
				//IMPORTANT!! - adding the .btn_ok class activates the CSS animation
				$(hud).find(".add_to_image_flow_btn .icon").addClass("btn_ok");
				//-------------------

				$(hud).find(".add_to_image_flow_btn").css("pointer-events", "none");
			},
			(p_error_data_map)=>{},
			p_log_fun);
	});
	//------------

	//------------
	// ADD_TO_POST
	const add_to_post__selector_str = '#'+p_image_id_str+' .add_to_post_btn';
	$(document).on('click', add_to_post__selector_str,()=>{
		add_image_to_post(p_image_info_map,p_log_fun);
	});
	//------------
	// GIF
	if (full_img_src_str.split('.').pop() == 'gif') {

		//$(p_image_in_page_element).find('.gf_gif').gifplayer();

		/*const gif_element = $(p_image_in_page_element).find('img')[0];
		console.log(gif_element)
		var sup1 = new SuperGif({ gif: gif_element } );
		sup1.load();
		sup1.play();
		console.log('22222222222222222222')*/
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
		$(hud).css('visibility', 'hidden'); //.remove();
	});
}

//---------------------------------------------------
function init_video_hud(p_video_in_page_element, p_video_info_map, p_log_fun) {
	//p_log_fun('FUN_ENTER','display_page_info.init_video_hud()');

	const hud = $(
		'<div class="hud">'+
			'<div class="add_to_post_btn">'+
				'<div class="symbol">'+
					'<div class="icon"></div>'+
				'</div>'+
			'</div>'+
		'</div>');
	const icons_chrome_ext_url_str = 'url('+chrome.extension.getURL('assets/icons.png')+')';
	$(hud).find('.icon').css('background-image', icons_chrome_ext_url_str);

	$(p_video_in_page_element).mouseenter(function(p_e) {
		$(p_video_in_page_element).append(hud);
	});

	$(p_video_in_page_element).mouseleave(function(p_e) {
		$(hud).remove();
	});
}

//---------------------------------------------------
function add_image_to_flow(p_full_img_src_str,
	p_images_flows_names_lst,
	p_gf_host_str,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun("FUN_ENTER", "display_page_info.add_image_to_flow()");

	const image_origin_page_url_str = window.location.href;
	// //-------------------
	// // IMPORTANT!! - since request to host_str is made from the context of the page in which the 
	// //               content is located, browser security imposes that the same protocol (http|https)
	// //               is used to communicate with host_str as with the origin-domain of the page
	// const origin_url_str = window.location.href;
	// const protocol_str   = origin_url_str.split("://")[0];
	// const host_str       = `${protocol_str}://gloflow.com`;
	// //-------------------

	http__add_image_to_flow(p_full_img_src_str,
		image_origin_page_url_str,
		p_images_flows_names_lst,
		p_gf_host_str, // host_str,
		(p_images_job_id_str, p_image_id_str, p_thumbnail_small_relative_url_str)=>{

			console.log("image added")
			console.log(`image job ID    - ${p_images_job_id_str}`)
			console.log(`image ID        - ${p_image_id_str}`)
			console.log(`thumb small URL - ${p_thumbnail_small_relative_url_str}`)

			p_on_complete_fun();
		},
		(p_error_data_map)=>{p_on_error_fun(p_error_data_map)},
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