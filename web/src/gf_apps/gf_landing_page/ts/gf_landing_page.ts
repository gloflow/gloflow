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

import * as gf_calc               from "./gf_calc";
import * as gf_email_registration from "./gf_email_registration";
import * as gf_images             from "./gf_images";
import * as gf_procedural_art     from "./gf_procedural_art";

$(document).ready(()=>{
	//-------------------------------------------------
	function log_fun(p_g,p_m) {
		var msg_str = p_g+':'+p_m
		//chrome.extension.getBackgroundPage().console.log(msg_str);

		switch (p_g) {
			case "INFO":
				console.log("%cINFO"+":"+"%c"+p_m, "color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
				break;
			case "FUN_ENTER":
				console.log("%cFUN_ENTER"+":"+"%c"+p_m, "color:yellow; background-color:lightgray","background-color:lightgray");
				break;
		}
	}
	//-------------------------------------------------
	$("time.timeago").timeago();

	init_remote(log_fun);

	//----------------------
	//IMPORTANT!! - wait for all images in the page to load first
	$(window).on("load", ()=>{
		gf_calc.run(log_fun);
	});
	//----------------------
});
//--------------------------------------------------------
export function init_remote(p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_landing_page.init_remote()');

	init(remote_register_user_email, p_log_fun);
	//--------------------------------------------------------
	function remote_register_user_email(p_inputed_email_str :string,
		p_on_complete_fun,
		p_log_fun) {
		p_log_fun('FUN_ENTER', 'gf_landing_page.init_remote().remote_register_user_email()');
		
		const url_str       = '/landing/register_invite_email';
		const data_args_map = {
			'email_str': p_inputed_email_str
		};
		
		$.ajax({
			'url':         url_str,
			'type':        'POST',
			'data':        JSON.stringify(data_args_map),
			'contentType': 'application/json',
			'success':     (p_data_map)=>{
	     		p_on_complete_fun('success', p_data_map);
			}
		});
	}
	//--------------------------------------------------------
}
//--------------------------------------------------------
function init(p_register_user_email_fun, p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_landing_page.init()');

	const featured_elements_infos_lst = load_static_data(p_log_fun);
	
	gf_procedural_art.init(p_log_fun);
	gf_email_registration.init(p_register_user_email_fun, p_log_fun);

	init_posts_img_num();
	gf_images.init(p_log_fun);

	//draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
	$(window).resize(()=>{

		//small screen widths dont display procedural_art
		if ($(window).innerWidth() > 660) {
			gf_procedural_art.init(p_log_fun);
		}
	});
	//--------------------------------------------------------
	function init_posts_img_num() {
		p_log_fun('FUN_ENTER', 'gf_landing_page.init().init_posts_img_num()');

		$('#featured_posts .post_info').each((p_i, p_post)=>{

			const post_images_number = $(p_post).find('.post_images_number')[0];
			const label_element      = $(post_images_number).find('.label');

			//HACK!! - "-1" was visually inferred
			$(post_images_number).css('right','-'+($(post_images_number).outerWidth()-1)+'px');
			$(label_element).css('left',$(post_images_number).outerWidth()+'px');

			$(p_post).mouseover((p_e)=>{
				$(post_images_number).css('visibility','visible');
			});
			$(p_post).mouseout((p_e)=>{
				$(post_images_number).css('visibility','hidden');
			});
		});
	}
	//--------------------------------------------------------
}
//--------------------------------------------------------
function load_static_data(p_log_fun) :Object[] {
	p_log_fun('FUN_ENTER', 'gf_landing_page.load_static_data()');
	
	const featured_elements_infos_lst :Object[] = []; 

	$('#posts .post_info').each((p_i)=>{
		const element = this;
		const featured_element_image_url_str  :string = $(element).find('img').attr('src');
		const featured_element_images_num_str :string = $(element).find('.post_images_number').find('.num').text();
		const featured_element_title_str      :string = $(element).find('.post_title').text();

		const featured_element_info_map :Object = {
			'element':    $(element),
			'image_src':  featured_element_image_url_str,
			'images_num': featured_element_images_num_str,
			'title_str':  featured_element_title_str
		};

		featured_elements_infos_lst.push(featured_element_info_map);
	});

	return featured_elements_infos_lst;
}
//--------------------------------------------------------
/*layout(Function p_log_fun) {
	//p_log_fun('FUN_ENTER','gf_landing_page.layout()');

	
	layout_featured_columns(p_log_fun);

	//--------------------------------------------------------
	setup_re_layout_on_resize() {
		//draw a new canvas when the view is resized, and delete the old one (with the old dimensions)
		window.onResize.listen((Event p_event) {

				//small screen widths dont display procedural_art
				if (window.innerWidth > 660) {
					gf_procedural_art.init(p_log_fun);
				}

				layout_featured_columns(p_log_fun);


				gf_email_registration.layout_email_button(p_log_fun);

				//in case the email registration form is open, this 
				//will reposition it to properly fit the new layout
				gf_email_registration.layout_email_form(p_log_fun);
			});
	}
	//--------------------------------------------------------
	setup_re_layout_on_resize();
}
//--------------------------------------------------------
function layout_featured_columns(Function p_log_fun,
					{int p_columns_distance_int:8}) {
	//p_log_fun('FUN_ENTER','gf_landing_page.layout_featured_columns()');

	//--------------------------------------------------------
	layout_featured_posts_display() {
		//p_log_fun('FUN_ENTER','gf_landing_page.layout_featured_columns().layout_featured_posts_display()');

		queryAll('.featured_posts .post_info').forEach((DivElement p_element) {

				final DivElement post_images_number = p_element.query('.post_images_number');
				final DivElement num_element        = post_images_number.query('.num');
				final DivElement label_element      = post_images_number.query('.label');

				//HACK!! - "-1" was visually inferred
				post_images_number.style.right = '-${post_images_number.offsetWidth-1}px';
				label_element.style.left       = '${post_images_number.offsetWidth}px';

				p_element.onMouseOver.listen((p_e) {
					post_images_number.style.visibility = 'visible';
				});
				p_element.onMouseOut.listen((p_e) {
					post_images_number.style.visibility = 'hidden';
				});
			});
	}
	//--------------------------------------------------------
	layout_featured_posts_display();
}*/