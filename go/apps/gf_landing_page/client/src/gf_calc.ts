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

///<reference path="./d/jquery.d.ts" />

namespace gf_calc {
//--------------------------------------------------------
declare var Fingerprint;
declare var ColorThief;
//--------------------------------------------------------
export function run(p_log_fun) {
	p_log_fun('FUN_ENTER','gf_calc.run()')

	const fingerprint = new Fingerprint({canvas: true}).get();
	const colorThief  = new ColorThief();

	var results_lst = []

	//------------------------
	//FEATURED_POSTS IMAGES
	$.each($('#featured_posts .post_image img'),(i,p_img)=>{

		const browser_run__job_result_map = process_image(p_img);
		const hex_color_str               = browser_run__job_result_map['c'];
		results_lst.push(browser_run__job_result_map);

		//-------------------
		//DISPLAY IMAGE DOMINANT_COLOR

		const parent_div = $(p_img).parent().parent()[0];
		const img_dominant_color_e = $('<div class="img_dominant_color"><div class="color"></div></div>');
		$(img_dominant_color_e).find('.color').css('background-color','#'+hex_color_str);

		$(parent_div).append(img_dominant_color_e);
		//-------------------
	});
	//------------------------
	//FEATURED_IMAGES
	$.each($('#featured_images .image img'),(i,p_img)=>{

		const browser_run__job_result_map = process_image(p_img);
		const hex_color_str               = browser_run__job_result_map['c'];
		results_lst.push(browser_run__job_result_map);



		//-------------------
		//DISPLAY IMAGE DOMINANT_COLOR

		const parent_div = $(p_img).parent().parent()[0];
		const img_dominant_color_e = $('<div class="img_dominant_color"><div class="color"></div></div>');
		$(img_dominant_color_e).find('.color').css('background-color','#'+hex_color_str);

		$(parent_div).append(img_dominant_color_e);
		//-------------------
	});
	//------------------------

	const job_results_map = {
		'jr':results_lst
	};
	send_calc_results(job_results_map,
					p_log_fun);
	//--------------------------------------------------------
	function process_image(p_img) {
		p_img.crossOrigin = '';
		//-------------------------
		//GET IMAGE GF_ID
		//image src example:
		//https://s3-us-west-1.amazonaws.com/gf--prod/thumbnails/d35afafdf2a83c53a2b93b8dbc99fb9f_thumb_medium.jpeg

		const img_url_str = p_img.src;
		const l           = img_url_str.split("/");

		const img_name_str = l[l.length-1];
		const img_id_str   = img_name_str.split('_')[0];
		//-------------------------
		//JOB_CALC
		var start = Date.now()/1000; //unix time in seconds

		//----------------------
		//IMPORTANT!! - for colorThief to work, using canvas and getImageData(), cross-origin policy 
		//              on both the image, and the S3 bucket from which the image is served, have to be turned on
		//              <img src="{{.Image_url_str}}" crossOrigin="anonymous"></img>

		//get a dominant color for an image
		var color_lst   = colorThief.getColor(p_img);   //DOMINANT COLOR
		var palette_lst = colorThief.getPalette(p_img); //COLOR PALLETE

		var end = Date.now()/1000; //unix time in seconds
		//-------------------------
		//COLORS RGB->HEX

		const hex_color_str = rgb_to_hex(color_lst[0],color_lst[1],color_lst[2]);
		p_log_fun('INFO','hex_color_str - '+hex_color_str);

		const hex_pallete_lst = [];
		for (var e of palette_lst) {

			const hex_str = rgb_to_hex(e[0],e[1],e[2]);
			hex_pallete_lst.push(hex_str);
		}
		//-------------------------

		const browser_run__job_result_map = {
			'i' :img_id_str,
			'c' :hex_color_str,
			'p' :hex_pallete_lst,
			'st':start,
			'et':end,
			'f' :fingerprint
		};
		//console.log(browser_run__job_result_map);
		return browser_run__job_result_map;
	}
	//--------------------------------------------------------
	function to_hex(c) {
	    var hex = c.toString(16);
	    return hex.length == 1 ? "0" + hex : hex;
	}
	//--------------------------------------------------------
	function rgb_to_hex(r,g,b) {
	    return to_hex(r) + to_hex(g) + to_hex(b);
	}
	//--------------------------------------------------------
}
//--------------------------------------------------------
function send_calc_results(p_job_results_map,
						p_log_fun) {
	p_log_fun('FUN_ENTER','gf_calc.send_calc_results()')
 	
	console.log(p_job_results_map);

	$.ajax({
		url     :'/images/c',
		type    :'POST',
		data    :JSON.stringify(p_job_results_map),
		dataType:'json',
		success :()=>{

		}
	});
}
//--------------------------------------------------------
}