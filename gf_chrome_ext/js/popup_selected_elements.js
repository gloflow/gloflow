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
function init_selected_elements_view(p_log_fun) {
	p_log_fun('FUN_ENTER', 'popup_selected_elements.init_selected_elements_view()');

	$(document).on('click','#view_selected_elements_card .symbol',(p_e)=>{
		get_selected_elements___bckg_pg((p_selected_elements_map)=>{
				show_selected_elements(p_selected_elements_map, p_log_fun);
			},
			p_log_fun);
	});
}

//-------------------------------------------------
function get_selected_elements___bckg_pg(p_o_cComplete_fun, p_log_fun) {
	p_log_fun('FUN_ENTER','popup_selected_elements.get_selected_elements___bckg_pg()');

	const msg_map = {
		'source_str':'popup_selected_elements',
		'type_str':  'get_selected_elements'
	};
	chrome.runtime.sendMessage(msg_map,
		(p_response) => {

			const selected_elements_map = p_response.selected_elements_map;
			//p_log_fun('INFO',selected_elements_map);

			p_on_complete_fun(selected_elements_map);
		});
}

//-------------------------------------------------
//DRAW
//-------------------------------------------------
function show_selected_elements(p_selected_elements_map, p_log_fun) {
	p_log_fun('FUN_ENTER', 'popup_selected_elements.show_selected_elements()');

	const selected_images_lst = p_selected_elements_map['images_lst'];
	//const selected_videos_lst = p_selected_elements_map['videos_lst'];

	show_selected_images(selected_images_lst, p_log_fun);

	$("#card_1").animate({
		'left':"-=400"
	}, 200, (p_e)=>{});

	$("#card_2").animate({
		'left':"-=400"
	}, 200, (p_e)=>{});
}

//-------------------------------------------------
//SHOW SELECTED IMAGES

function show_selected_images(p_img_infos_lst, p_log_fun) {
	p_log_fun('FUN_ENTER','popup_selected_elements.show_selected_images()');

	const container = $(
		'<div id="selected_images">'     +
			'<div id="collection"></div>'+
		'</div>');
	$('#card_2').append(container);


	$.each(p_img_infos_lst,(p_i,p_image_info_map)=>{

		const full_img_src_str = p_image_info_map["full_img_src_str"];
		const raw_img_src_str  = p_image_info_map["raw_img_src_str"];

		//p_log_fun('INFO',p_image_info_map['full_img_src_str']);
		
		//raw_img_src_str - is included here, in <img> and <h5>
		//                  for debugging purposes, to see if the browser loads 
		//                  the raw urls extracted from the target pages DOM
		
		//"data-*" - html5 allows for these data attributes
		const thumbnail_html_str = $(
			'<div class="selected_image">'                                          +
				'<img src='+full_img_src_str+'></img>'                              +
				//'<div><span>image url</span><span>'+full_img_src_str+'</span></div>'+
			'</div>');
                         
		$(container).find('#collection').append(thumbnail_html_str);

		$(thumbnail_html_str).find('img').load((p_e)=>{

			//-------------------
			//IMPORTANT!! - reload the masonry layout with the newly loaded image
			$(container).find('#collection').masonry();
			//-------------------
		});
	});

	//------------
	//MASONRY

	$(container).find('#collection').masonry(
		{
			columnWidth:  10,
			gutter:       10,
			itemSelector: '.selected_image'
		});
	//------------
}

//-------------------------------------------------
//SHOW SELECTED VIDEOS

function show_selected_videos(p_videos_infos_lst, p_log_fun) {
	p_log_fun('FUN_ENTER','popup_selected_elements.show_selected_videos()')
	
	for (var i=0;i<p_videos_infos_lst.length;i++) {

		const video_info_map = p_videos_infos_lst[i];
		p_log_fun('INFO','video_info_map:'+video_info_map);

		//---------------------------------
		//YOUTUBE
		
		//{
		//	'youtube_video_embed_url_str':youtube_video_embed_url_str
		//}
		
		if ('youtube_video_embed_url_str' in video_info_map) {
			const youtube_video_embed_url_str = video_info_map['youtube_video_embed_url_str'];
			
			//ADD!! - //'webkitAllowFullScreen mozallowfullscreen allowFullScreen' to the iframe 
			const video_embed_html_str = 
				'<div>'                      +
					'<iframe id    ="video"' +
					'        width ="420"'   +
					'        height="315"'   + 
					'        src   ="'+youtube_video_embed_url_str+'"'         + 
					'        frameborder="0"></iframe>'                        +
					//'<button class            ="btn btn-success btn-small use_as_post_element_video_btn" ' +
					//				'data-video_embed_url_str="'+youtube_video_embed_url_str+'" '          +
					//				'>use as post element video</button>'                                  +
				'</div>';
                         
			$('#page_content_display_container').append(video_embed_html_str)
		}
		//---------------------------------
		//VIMEO
		
		else if ('vimeo_video_embed_url_str' in video_info_map ) {
			const vimeo_video_embed_url_str = video_info_map['vimeo_video_embed_url_str'];
			p_log_fun('INFO','vimeo_video_embed_url_str:'+vimeo_video_embed_url_str);
			
			//ADD!! - //'webkitAllowFullScreen mozallowfullscreen allowFullScreen' to the iframe 
			const video_embed_html_str = 
				'<div>'                       +
					'<iframe id    ="video"'  +
					'        width ="500"'    +
					'        height="281" '   +
					'        src   ="'+vimeo_video_embed_url_str+'"' + 
					'        frameborder="0"></iframe>'              +
					//'<button class            ="btn btn-success btn-small use_as_post_element_video_btn" ' +
					//		'data-video_embed_url_str="'+vimeo_video_embed_url_str+'" '                    +
					//		'>use as post element video</button>'                                          +
				'</div>';
																 
			$('#page_content_display_container').append(video_embed_html_str);
		}
		
		else if ('vimeo_flash_vars_str' in video_info_map) {
			const vimeo_flash_vars_str = video_info_map['vimeo_flash_vars_str'];
			p_log_fun('INFO','vimeo_flash_vars_str:'+vimeo_flash_vars_str);
			
			const video_object_html_str = 
				'<object type   = "application/x-shockwave-flash"'                                        +
						'id     = "player7174318_2088053176"'                                             +
						'name   = "player7174318_2088053176"'                                             + 
						'class  = "'                                                                      +
						'data   = "http://a.vimeocdn.com/p/flash/moogaloop/5.2.39/moogaloop.swf?v=1.0.0"' +
						'width  = "100%"'                                                                 +
						'height = "100%"'                                                                 + 
						'style  = "visibility: visible;">'                                                +
					'<param name="allowscriptaccess" value="always">'                   +
					'<param name="allowfullscreen"   value="true">'                     +
					'<param name="scalemode"         value="noscale">'                  +
					'<param name="quality"           value="high">'                     +
					'<param name="wmode"             value="opaque">'                   +
					'<param name="bgcolor"           value="#000000">'                  +
					'<param name="flashvars"         value="'+vimeo_flash_vars_str+'">' + 
				'</object>';
																 
			$('#page_content_display_container').append(video_object_html_str);
		}
		//---------------------------------
		//FIX!! - DOES NOT WORK!! - ooyala does not have an iframe embed like this, 
		//NOTE - I hardly ever run into ooyala player (July2014)

		//OOYALA
		else if ('ooyala_video_embed_url_str' in video_info_map ) {
			const ooyala_video_embed_url_str = video_info_map['ooyala_video_embed_url_str'];
			
			//'<iframe id    ="video"'  +
			//'        width ="500"'    +
			//'        height="281" '   + 
			//'webkitAllowFullScreen mozallowfullscreen allowFullScreen'
			//'        src   ="'+ooyala_video_embed_url_str+'"'  + 
			//'        frameborder="0"></iframe>' +
		 
			const video_embed_html_str = 
				'<div>'                                                                        +
					'<script src  ='+ooyala_video_embed_url_str+'></script>'                   +
					//'<button class="btn btn-success btn-small use_as_post_element_video_btn" ' +
					//		'data-video_embed_url_str="'+vimeo_video_embed_url_str+'" '        +
					//		'>use as post element image</button>'                              +
				'</div>';
																 
			$('#page_content_display_container').append(video_embed_html_str);
		}
		//---------------------------------
	}
	
	$('.use_as_post_element_video_btn').click(function(p_e) {
		const video_embed_url_str = $(p_e.target).attr('data-video_embed_url_str');
		p_log_fun('INFO',video_embed_url_str);
		
		//data-type_str - attribute gets sent to the server
		const picked_video_listing_node = $(
			'<span class="picked_video_listing">'                                                   +
				'<i class="icon-remove" style="cursor:pointer"></i>'                                +
				'<a  data-type_str="video" href='+video_embed_url_str+'>'+video_embed_url_str+'</a>'+
				'<br>'                                                                              +
			'</span>');
		$('#picked_page_assets_lst').append(picked_video_listing_node);
		
		//----------------
		//ICON-REMOVE - little 'x' next to the image url
		
		//when the icon-remove is clicked the picked_image_listing_node 
		//should also get removed from the DOM (and subsequently not 
		//loaded for post creation)
		$(picked_video_listing_node).find('i').click(function(p_e) {
			//$(picked_image_listing_node).remove()
			
			
			//$(p_e.target).parent() - reffers to <span class='picked_img_listing'>
			$(p_e.target).parent().remove();
		});
		//----------------
	});
}