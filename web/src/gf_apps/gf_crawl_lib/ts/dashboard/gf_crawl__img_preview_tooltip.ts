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

///<reference path="../../../../d/jquery.d.ts" />

namespace gf_crawl__img_preview_tooltip {

//---------------------------------------------------
export function init(p_url_str,
	p_crawler_page_image_id_str,
	p_origin_page_url_str,
	p_url_element,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_crawl__img_preview_tooltip.init()');

	const flows_names_lst = ['general'];

	var tooltip_element;
	var img_visible_bool = false;
	var i                = 0;
	$(p_url_element).on('mouseover', ()=>{

		if (!img_visible_bool) {
		
			tooltip_element = $(`
				<div class='img_preview_tooltip'>
					<div class='add_to_image_flow_btn'>
						<div class='symbol'><div class='icon'></div></div>
					</div>
					<div class='origin_page_url'>
						<a href='`+p_origin_page_url_str+`'>`+p_origin_page_url_str+`</a>
					</div>
					<img src='`+p_url_str+`'></img>
				</div>`);

			p_url_element.append(tooltip_element);	

			const img_element = $(tooltip_element).find('img');
			$(img_element).on('load',(p_e)=>{
				
				const image_height_int = this.naturalHeight;
				const image_width_int  = this.naturalWidth;

				$(tooltip_element).append(`
					<div class='image_dimensions'>
						<span>`+image_height_int+`</span>px x <span>`+image_width_int+`</span>px
					</div>`);

				i++;
			});

			//--------------------
			$(tooltip_element).find('.add_to_image_flow_btn').on('click',()=>{

				http__add_to_image_flows(p_crawler_page_image_id_str,
					flows_names_lst,
					function() {

						console.log('FLOW_DONE>>>>>>>>>>>>>>>>>>>')

						//-------------------
						//IMPORTANT!! - adding the .btn_ok class activates the CSS animation
						$(tooltip_element).find('.add_to_image_flow_btn .icon').addClass('btn_ok');
						//-------------------
					},
					p_log_fun);
			});
			//--------------------

			img_visible_bool = true;
		}
	});

	$(p_url_element).on('mouseleave',()=>{
		$(tooltip_element).remove();
		img_visible_bool = false;
	});
}
//---------------------------------------------------
function http__add_to_image_flows(p_crawler_page_image_id_str,
	p_flows_names_lst,
	p_on_complete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_crawl__img_preview_tooltip.http__add_to_image_flows()');

	const url_str = '/a/crawl/image/add_to_flow';
	p_log_fun('INFO','url_str - '+url_str);

	const data_map = {
		'crawler_page_image_id_str': p_crawler_page_image_id_str,
		'flows_names_lst':           p_flows_names_lst,
	};    
	//-------------------------
	//HTTP AJAX
	$.post(url_str,
		JSON.stringify(data_map),
		(p_data_map) => {
			p_on_complete_fun(p_data_map);
		});
	//-------------------------
}
//---------------------------------------------------
}