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

///<reference path="../../../d/jquery.d.ts" />
///<reference path="../../../d/jquery.timeago.d.ts" />

import * as gf_image_colors from "./../../../gf_core/ts/gf_image_colors";
import * as gf_time         from "./../../../gf_core/ts/gf_time";

//-------------------------------------------------
export function init(p_log_fun) {

	$('#featured_images_0').find('.image_info').each((p_i, p_image_info_element)=>{
		
		init_img(p_image_info_element);
	});
	$('#featured_images_1').find('.image_info').each((p_i, p_image_info_element)=>{
		
		init_img(p_image_info_element);
	});

	//-------------------------------------------------
	function init_img(p_image_info_element) {
		// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
		//           this origin_page_url is set to empty string. check for that and remove it.
		// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
		//         at all for images that dont have an origin_page_url.
		if ($(p_image_info_element).find(".origin_page_url a").text().trim() == "") {
			$(p_image_info_element).find(".origin_page_url").remove();
		}


		// FIX!! - this function has been moved to gf_core/gf_images_viewer.ts, as a general viewer,
		//         to use universaly gf_images/gf_landing_page.
		//         gf_images flows_browser already uses the version from gf_core.
		init_image_viewer(p_image_info_element, p_log_fun);

		gf_time.init_creation_date(p_image_info_element, p_log_fun);

		//----------------------
		// IMAGE_PALLETE
		const img = $(p_image_info_element).find("img")[0];

		const assets_paths_map = {
			"copy_to_clipboard_btn": "/images/static/assets/gf_copy_to_clipboard_btn.svg",
		}
		gf_image_colors.init_pallete(img,
			assets_paths_map,
			(p_color_dominant_hex_str,
			p_colors_hexes_lst)=>{



				console.log("image colors pallete", p_color_dominant_hex_str, p_colors_hexes_lst)

				$(p_image_info_element).css("background-color", p_color_dominant_hex_str);
				$(p_image_info_element).find(".image_title").css("background-color", p_color_dominant_hex_str);
				$(p_image_info_element).find(".origin_page_url").css("background-color", p_color_dominant_hex_str);

			});

		//----------------------
	}

	//-------------------------------------------------
}

//-------------------------------------------------
// DEPRECATED!!
// REMOVE!!
function init_image_viewer(p_image_element, p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_landing_page.init_image_viewer()');

	const img_thumb_medium_url = $(p_image_element).find('img').data('img_thumb_medium_url');
	const image_view           = $(`
		<div id="image_view">
			<div id="background"></div>
			<div id="image_detail">
				<img src="${img_thumb_medium_url}"></img>
			</div>
		</div>`);

	$(p_image_element).find('img').click(()=>{
		
		console.log(img_thumb_medium_url)
		$('body').append(image_view);

		//----------------------
		// BAKCGROUND
		const bg = $(image_view).find('#background');

		// position the background vertically where the user has scrolled to
		$(bg).css('top', $(window).scrollTop()+'px');

		//----------------------
		// IMPORTANT!! - turn off vertical scrolling while viewing the image
		$("body").css("overflow-y", "hidden");

		//----------------------
		const image_detail = $(image_view).find('#image_detail');
		$(image_detail).css("position", "absolute");
	    $(image_detail).css("top",  Math.max(0, (($(window).height() - $(image_detail).outerHeight()) / 2) + $(window).scrollTop())  + "px");
	    $(image_detail).css("left", Math.max(0, (($(window).width()  - $(image_detail).outerWidth()) / 2)  + $(window).scrollLeft()) + "px");
	    
		//----------------------
	    $(bg).click(()=>{
	    	$(image_view).remove();
	    	$("body").css("overflow-y", "auto"); // turn vertical scrolling back on when done viewing the image
	    });

	    //----------------------
	});
}

//-------------------------------------------------
/*function init_image_date(p_target_element, p_log_fun) {

	const creation_time_element = $(p_target_element).find('.creation_time');
	const creation_time_f       = parseFloat($(creation_time_element).text());
	const creation_date         = new Date(creation_time_f*1000);

	const date_msg_str = $.timeago(creation_date);
	$(creation_time_element).text(date_msg_str);

	const creation_date__readable_str = creation_date.toDateString();
	const creation_date__readble      = $('<div class="full_creation_date">'+creation_date__readable_str+'</div>');

	$(creation_time_element).mouseover((p_e)=>{
		$(creation_time_element).append(creation_date__readble);
	});

	$(creation_time_element).mouseout((p_e)=>{
		$(creation_date__readble).remove();
	});
}*/