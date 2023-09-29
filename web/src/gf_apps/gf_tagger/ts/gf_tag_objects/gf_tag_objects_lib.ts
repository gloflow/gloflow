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

///<reference path="../../../../d/jquery.d.ts" />
///<reference path="../../../../d/masonry.layout.d.ts" />
///<reference path="../../../../d/jquery.timeago.d.ts" />

import * as gf_image_viewer from "./../../../../gf_core/ts/gf_image_viewer";
import * as gf_sys_panel    from "./../../../../gf_sys_panel/ts/gf_sys_panel";
import * as gf_utils        from "./../../../gf_images/ts/gf_images_flows_browser/gf_utils";

// GF_GLOBAL_JS_FUNCTION - included in the page from gf_core (.js file)
declare var gf_upload__init;

//-------------------------------------------------
export function init(p_log_fun) {

	// SYS_PANEL
    gf_sys_panel.init_with_auth(p_log_fun);

    //---------------------
	// MASONRY
    $('#images_container').masonry({
		itemSelector: '.gf_image',
		columnWidth:  6
	});

    $('#posts_container').masonry({
		itemSelector: '.gf_post',
		columnWidth:  6
	});


    /*
	IMPORTANT!! - as each image loads call masonry to reconfigure the view.
		this is necessary so that initial images in the page, before
		load_new_page() starts getting called, are properly laid out
		by masonry.
	*/
	$('.gf_image img').on('load', ()=>{
		$('#images_container').masonry();
	});
	$('.gf_post img').on('load', ()=>{
		$('#posts_container').masonry();
	});
    //---------------------

    init_images(p_log_fun);
}

//-------------------------------------------------
function init_images(p_log_fun) {

    $('#images_container .gf_image').each((p_i, p_e)=>{

		const image_element = p_e;
		gf_utils.init_image_date(image_element, p_log_fun);

		const image_id_str = $(image_element).data('data-img_id');
		const image_flows_names_lst = $(image_element).data('data-img_flows_names').split(",");
		const img_thumb_medium_url_str = $(image_element).find('img').data('img_thumb_medium_url');
		const img_thumb_large_url_str  = $(image_element).find('img').data('img_thumb_large_url');
		const img_format_str           = $(image_element).attr('data-img_format');
		const origin_page_url_link = $(image_element).find(".origin_page_url a")[0];

		//----------------
		// CLEANUP - for images that dont come from some origin page (direct uploads, or generated images)
		//           this origin_page_url is set to empty string. check for that and remove it.
		// FIX!! - potentially on the server/template-generation side this div node shouldnt get included
		//         at all for images that dont have an origin_page_url.
		if ($(origin_page_url_link).text().trim() == "") {
			$(image_element).find(".origin_page_url").remove();
		}

		//----------------
		// LINK_TEXT_SHORTEN - if the link text (not its href) is too long, dont display it completely in the UI
		//                     because it clutters the UI too much.
		//                     instead shorten it at its cutoff length and append "..."
		const link_text_cutoff_threshold_int = 50;
		if ($(origin_page_url_link).text().length > link_text_cutoff_threshold_int) {
			const old_link_text_str = $(origin_page_url_link).text();
			const new_link_text_str = `${old_link_text_str.slice(0, link_text_cutoff_threshold_int)}...`;
			$(origin_page_url_link).text(new_link_text_str);
		}

		//----------------
		else {

			gf_image_viewer.init(image_element,
				image_id_str,
				img_thumb_medium_url_str,
				img_thumb_large_url_str,
				image_flows_names_lst,
				p_log_fun);
		}

		//----------------
	});
}