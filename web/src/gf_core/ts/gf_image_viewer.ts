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

///<reference path="../../d/jquery.d.ts" />

//-------------------------------------------------
export function init(p_image_element,
	p_image_id_str             :string,
	p_img_thumb_medium_url_str :string,
	p_img_thumb_large_url_str  :string,
	p_flows_names_lst          :string[],
	p_log_fun) {

	$(p_image_element).find("img").click(()=>{

		console.log("click")

		const gf_link_str = `/images/v/${p_image_id_str}`;

		const image_view_element = $(`
			<div id="image_viewer">
				<div id="background"></div>
				<div id="image_detail">
					<img src="${p_img_thumb_large_url_str}"></img>
					<div class="flows_names">${p_flows_names_lst}</div>
					<div class="image_view_link">
						<a href="${gf_link_str}">gf link</a>
					</div>
				</div>
			</div>`);

		$('body').append(image_view_element);

		const image_detail_element = $(image_view_element).find("#image_detail");
		$(image_detail_element).css("position", "absolute");

		//----------------------
		// BAKCGROUND
		const bg = $(image_view_element).find("#background");

		// position the background vertically where the user has scrolled to
		$(bg).css("top", `${$(window).scrollTop()}px`);

		//----------------------
		// IMPORTANT!! - turn off vertical scrolling while viewing the image
		$("body").css("overflow", "hidden");

		//----------------------
		// IMG_ONLOAD
		$(image_view_element).find("img").on("load", ()=>{
			image_position_and_scale(image_detail_element);
		});

		//-------------------------------------------------
		function resize_handler() {
			image_position_and_scale(image_detail_element);
		}

		//-------------------------------------------------
		// reposition image detail if the window resizes
		$(window).on("resize", resize_handler);

	    //----------------------
	    $(bg).click(()=>{
	    	$(image_view_element).remove();
			$(window).off("resize", resize_handler); // stop positioning on resize

	    	// turn vertical scrolling back on when done viewing the image
	    	$("body").css("overflow", "auto");
	    });

	    //----------------------
	});
}

//-------------------------------------------------
function image_position_and_scale(p_image_detail_element) {

	const window_width_int        = $(window).width() - 100;  // some padding removed from real window width
	const window_height_int       = $(window).height() - 100; // some padding removed from real window height
	const image_detail_width_int  = $(p_image_detail_element).outerHeight();
	const image_detail_height_int = $(p_image_detail_element).outerHeight();

	//----------------------
	// WIDTH
	// image view is larger than the window height, so scale it back
	if (image_detail_width_int > window_width_int) {
		$(p_image_detail_element).css("width", "90%");
	}

	// HEIGHT
	// image view is larger than the window height, so scale it back
	if (image_detail_height_int > window_height_int) {
		$(p_image_detail_element).css("height", "90%");
	}

	//----------------------
	// POSITION
	// its important that position is calculated after width/height setting,
	// so that new dimensions are taken into account

	// Math.max() - returns the largest of zero or more numbers.
	// Math.max(10, 20);   //20
	// Math.max(-10, -20); //-10
	const image_detail_x = Math.max(0, (($(window).width() - $(p_image_detail_element).outerWidth()) / 2) + $(window).scrollLeft());
	const image_detail_y = Math.max(0, (($(window).height() - $(p_image_detail_element).outerHeight()) / 2) + $(window).scrollTop());

	$(p_image_detail_element).css("left", `${image_detail_x}px`);
	$(p_image_detail_element).css("top",  `${image_detail_y}px`);
	
	//----------------------
}