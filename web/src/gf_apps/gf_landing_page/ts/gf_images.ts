/*
GloFlow media management/publishing system
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

///<reference path="./d/jquery.timeago.d.ts" />

namespace gf_images {
//-------------------------------------------------
export function init(p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_images.init()');

	$('#featured_images').find('.image_info').each((p_i, p_e)=>{
		init_image_viewer(p_e, p_log_fun);
		init_image_date(p_e, p_log_fun);
	});
}
//-------------------------------------------------
//REMOVE!! - this function has been moved to flows_browser/gf_images_viewer.ts, as a general viewer,
//           to use universaly gf_images/gf_landing_page

function init_image_viewer(p_image_element, p_log_fun) {
	p_log_fun('FUN_ENTER','gf_images.init_image_viewer()');

	const img_thumb_medium_url = $(p_image_element).find('img').data('img_thumb_medium_url');
	const image_view           = $(`
		<div id="image_view">
			<div id="background"></div>
			<div id="image_detail">
				<img src="`+img_thumb_medium_url+`"></img>
			</div>
		</div>`);

	$(p_image_element).find('img').click(()=>{

		console.log(img_thumb_medium_url)
		$('body').append(image_view);

		//----------------------
		//BAKCGROUND
		const bg = $(image_view).find('#background');

		//position the background vertically where the user has scrolled to
		$(bg).css('top',$(window).scrollTop()+'px');
		//----------------------
		//IMPORTANT!! - turn off vertical scrolling while viewing the image
		$("body").css("overflow","hidden");
		//----------------------
		const image_detail = $(image_view).find('#image_detail');
		$(image_detail).css("position","absolute");
	    $(image_detail).css("top",  Math.max(0, (($(window).height() - $(image_detail).outerHeight()) / 2) + $(window).scrollTop())  + "px");
	    $(image_detail).css("left", Math.max(0, (($(window).width()  - $(image_detail).outerWidth()) / 2)  + $(window).scrollLeft()) + "px");
	    //----------------------
	    $(bg).click(()=>{
	    	$(image_view).remove();
	    	$("body").css("overflow", "auto"); //turn vertical scrolling back on when done viewing the image
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
//-------------------------------------------------
}