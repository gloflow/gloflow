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
function add_image_to_post(p_image_info_map, p_log_fun) {

	// first send the newly added post to the background_page
	add_element_to_post___bckg_pg(p_image_info_map,
		(p_response) => {
			switch(p_response.status_str) {
				//------------
				// only draw the image if it was added to the Post, who's state
				// is maintaned in the background page
				case 'success':
					draw();
					break;
				//------------
				// if this has already been added then do nothing
				case 'exists':
					break;
				//------------
			}
		},
		p_log_fun);

	//---------------------------------------------------
	function draw() {
		p_log_fun('FUN_ENTER','page_element_picker.add_image_to_post().draw()');

		const images_to_post_block_start_y    = 80;
		const all_previous_images_to_post_lst = $('body').find('.image_to_post');
		const previous_image_to_post          = all_previous_images_to_post_lst[all_previous_images_to_post_lst.length - 1];
		const img_name_str                    = p_image_info_map['img_name_str'];
		
		const image_to_post = $(
			'<div class="image_to_post">'+
				'<div class="close_btn"></div>'+
				'<div class="image_name">'+img_name_str+'</div>'+
			'</div>');

		//---------------------------------------------------
		// CLOSE_BTN

		function init_close_btn() {
			const icons_atlas_url_str = chrome.runtime.getURL("assets/icons.png");
			const close_btn           = $(image_to_post).find('.close_btn')[0];

			//--------
			// CSS
			const icons_chrome_ext_url_str = 'url('+chrome.runtime.getURL('assets/icons.png')+')';
			$(close_btn).css('background-image', icons_chrome_ext_url_str);
			//--------

			$(image_to_post).on('click','.close_btn',
				() => {
					remove_element_from_post_bckg_pg(p_image_info_map, ()=>{}, p_log_fun);
				});
		}

		//---------------------------------------------------
		function init_preview() {

			const image_url_str   = p_image_info_map['full_img_src_str'];
			const preview_element = $(
				'<div class="preview">'                +
					'<img src='+image_url_str+'></img>'+
				'</div>');

			$(image_to_post).mouseover(() => {
				$(image_to_post).append(preview_element);
			});

			$(image_to_post).mouseout(() => {
				$(preview_element).remove();
			});
		}

		//---------------------------------------------------
		init_close_btn();
		init_preview();

		//if there is at least one image_to_post
		if ($('body').find('.image_to_post').length > 0) {
			const previous_image_y = parseInt($(previous_image_to_post).css('top').replace('px',''));
			const new_y            = previous_image_y + parseInt($(previous_image_to_post).height()) + 10;

			$(image_to_post).css('top',new_y+'px');
		}
		else {
			$(image_to_post).css('top', images_to_post_block_start_y+'px');	
		}

		$('#page_info_container #selected_elements_preview').append(image_to_post);
	}
	//---------------------------------------------------
}

//---------------------------------------------------
// BACKGROUND_PAGE COMM
//---------------------------------------------------
function add_element_to_post___bckg_pg(p_element_info_map, p_on_complete_fun, p_log_fun) {	

	const msg_map = {
		'source_str':       'content_script',
		'type_str':         'add_element_to_post',
		'element_info_map': p_element_info_map
	};
	send_msg_to_bg_page(msg_map, p_on_complete_fun);
}

//---------------------------------------------------
function remove_element_from_post_bckg_pg(p_element_info_map, p_on_complete_fun, p_log_fun) {

	const msg_map = {
		'source_str':       'content_script',
		'type_str':         'remove_element_from_post',
		'element_info_map': p_element_info_map
	};
	send_msg_to_bg_page(msg_map, p_on_complete_fun);
}