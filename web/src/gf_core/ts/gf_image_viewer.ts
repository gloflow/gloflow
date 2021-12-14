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
	p_img_thumb_medium_url_str :string,
	p_img_thumb_large_url_str  :string,
	p_flow_name_str            :string,
	p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_image_viewer.init()");

	//const img_thumb_medium_url = $(p_image_element).find('img').data('img_thumb_medium_url');

	$(p_image_element).find("img").click(()=>{

		console.log("click")

		const image_view = $(`
			<div id="image_viewer">
				<div id="background"></div>
				<div id="image_detail">
					<img src="${p_img_thumb_large_url_str}"></img>
				</div>
			</div>`);

		$('body').append(image_view);

		//----------------------
		// BAKCGROUND
		const bg = $(image_view).find("#background");

		// position the background vertically where the user has scrolled to
		$(bg).css("top", $(window).scrollTop()+"px");

		//----------------------
		// IMPORTANT!! - turn off vertical scrolling while viewing the image
		$("body").css("overflow", "hidden");

		//----------------------
		// IMG_ONLOAD
		$(image_view).find("img").on("load", ()=>{

			const image_detail = $(image_view).find("#image_detail");
			$(image_detail).css("position", "absolute");

			// Math.max() - returns the largest of zero or more numbers.
			// Math.max(10, 20);   //20
			// Math.max(-10, -20); //-10
			const image_x = Math.max(0, (($(window).width() - $(image_detail).outerWidth()) / 2) + $(window).scrollLeft());
			const image_y = Math.max(0, (($(window).height() - $(image_detail).outerHeight()) / 2) + $(window).scrollTop());

			$(image_detail).css("left", image_x+"px");
		    $(image_detail).css("top",  image_y+"px");
		});

	    //----------------------
	    $(bg).click(()=>{
	    	$(image_view).remove();

	    	// turn vertical scrolling back on when done viewing the image
	    	$("body").css("overflow", "auto");
	    });

	    //----------------------
	});
}