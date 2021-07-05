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

import * as gf_image_process from "./../../../gf_core/ts/gf_image_process";

//-------------------------------------------------
export function init(p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_images.init()');

	$('#featured_images').find('.image_info').each((p_i, p_image_info_element)=>{

		// FIX!! - this function has been moved to gf_core/gf_images_viewer.ts, as a general viewer,
		//         to use universaly gf_images/gf_landing_page.
		//         gf_images flows_browser already uses the version from gf_core.
		init_image_viewer(p_image_info_element, p_log_fun);

		init_image_date(p_image_info_element, p_log_fun);



		//-------------
		// NEW - move into its own function
		var image_colors_shown_bool = false;
		$(p_image_info_element).find("img").on("mouseover", (p_event)=>{



			if (!image_colors_shown_bool) {
				
				const image        = p_event.target;
				const image_colors = gf_image_process.get_colors(image);


				const color_info_element = $(`<div class="colors_info">
					<div class="color_dominant"></div>
					<div class="color_pallete"></div>
				</div>`);
				
				color_info_element.insertAfter(image);

				//-------------
				// COLOR_DOMINANT
				const color_dominant_element = $(color_info_element).find(".color_dominant");
				$(color_dominant_element).css("background-color", image_colors.color_hex_str);

				//-------------
				// COLOR_PALLETE
				const color_pallete_element = $(p_image_info_element).find(".colors_info .color_pallete");
				// const color_pallete_sub_lst = image_colors.color_palette_lst.slice(1, 6);
				image_colors.color_palette_lst.forEach((p_color_hex_str)=>{

					console.log("-------------")
					console.log(p_color_hex_str);

					$(color_pallete_element).append(`<div class="color" style="background-color:#${p_color_hex_str};"></div>`)
				})

				//-------------

				var color_dominant_label_element = $(`<div class="color_dominant_label">color dominant</div>`);
				$(color_dominant_element).on("mouseover", ()=>{
					color_info_element.append(color_dominant_label_element);
				});
				$(color_dominant_element).on("mouseout", ()=>{
					$(color_dominant_label_element).remove();
				});


				var color_pallete_label_element = $(`<div class="color_pallete_label">color pallete</div>`);
				$(color_pallete_element).on("mouseover", ()=>{
					color_info_element.append(color_pallete_label_element);
				});
				$(color_pallete_element).on("mouseout", ()=>{
					$(color_pallete_label_element).remove();
				});
				


				image_colors_shown_bool = true;
			}
		});


		//-------------

	});
}

//-------------------------------------------------
// DEPRECATED!!
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
function init_image_date(p_image_element, p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_images.init_image_date()');

	const creation_time_element = $(p_image_element).find('.creation_time');
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
}