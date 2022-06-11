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

///<reference path="../../../../d/jquery.d.ts" />

import * as gf_tagger_client from "./gf_tagger_client";

//-----------------------------------------------------
// in gf_post view
export function init_tag_input(p_obj_id_str :string,
	p_obj_type_str :string,
	p_obj_element,
	p_on_tags_created_fun,
	p_on_tag_ui_add_fun,
	p_on_tag_ui_remove_fun,
	p_http_api_map,
	p_log_fun) {
	
	const tagging_input_ui_element = init_tagging_input_ui_element(p_obj_id_str,
		p_obj_type_str,
		p_on_tags_created_fun,
		p_on_tag_ui_remove_fun,
		p_http_api_map,
		p_log_fun);
		
	const tagging_ui_element = $(`
		<div class="post_element_controls">
			<div class="add_tags_button">add tags</div>
		</div>`);

	// OPEN TAG INPUT UI
	$(tagging_ui_element).find('.add_tags_button').on('click', (p_event)=>{

		// remove the tagging_input_container if its already displayed
		// for tagging another post_element
		if ($('#tagging_input_container') != null) {
			$('#tagging_input_container').remove();
		}

	
		place_tagging_input_ui_element(tagging_input_ui_element,
			p_obj_element,
			p_log_fun);

		if (p_on_tag_ui_add_fun != null) p_on_tag_ui_add_fun();
	});

	//------------------------

	/*//'T' key - open tagging UI to the element that has the cursor 
	//          hovering over it
	final subscription = document.onKeyUp.listen((p_event) {
		if (p_event.keyCode == 84) {

			//remove the tagging_input_container if its already displayed
			//for tagging another post_element
			if (query('#tagging_input_container') != null) {
				query('#tagging_input_container').remove();
			}

			place_tagging_input_ui_element(tagging_input_ui_element,
										   p_obj_element, //post_element,
										   p_log_fun);

			//prevent this handler being invoked while the user
			//is typing in tags into the input field
			//subscription.pause();
		}
	});*/

	//------------------------
	// IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
	//               but unline mouseon/mouseout it will not fire if its children are entered/left.
	
	$(p_obj_element).on('mouseenter', (p_event)=>{
		$(p_obj_element).append(tagging_ui_element);
	});

	$(p_obj_element).on('mouseleave', (p_event)=>{
		$(tagging_ui_element).remove();

		// // relatedTarget - The relatedTargert property can be used with the mouseover 
		// //                 event to indicate the element the cursor just exited, 
		// //                 or with the mouseout event to indicate the element the 
		// //                 cursor just entered.
		// if (p_event.relatedTarget != null && 
		// 	!p_event.relatedTarget.classes.contains('add_tags_button')) {
		//	tagging_ui_element.remove();
		// }
	});

	//------------------------
}

//-----------------------------------------------------
// TAGS UI UTILS
//-----------------------------------------------------
export function init_tagging_input_ui_element(p_obj_id_str :string,
	p_obj_type_str :string,
	p_on_tags_created_fun,
	p_on_tag_ui_remove_fun,
	p_http_api_map,
	p_log_fun) {
	
	const tagging_input_ui_element = $(`
		<div id="tagging_input_container">
			<div id="background"></div>
			<input type="text" id="tags_input" placeholder="(space) separated tags">
			<div id="submit_tags_button">add</div>
			<div id="close_tagging_input_container_button">&#10006;</div>
		</div>`);
	
	const tags_input_element = $(tagging_input_ui_element).find('#tags_input');

	// 'ESCAPE' key
	$(document).on('keyup', (p_event)=>{
		if (p_event.which == 27) {
			//remove any previously present tagging_input_container's
			$(tagging_input_ui_element).remove();
			if (p_on_tag_ui_remove_fun != null) {
				p_on_tag_ui_remove_fun();
			}
		}
	});

	// to handlers for the same thing, one for the user clicking on the button,
	// the other for the user pressing 'enter'  
	$(tags_input_element).on('keyup', async (p_event)=>{

			// 'ENTER' key
			if (p_event.which == 13) {
				p_event.preventDefault();
				
				const tags_lst = await add_tags_to_obj(p_obj_id_str,
					p_obj_type_str,
					tagging_input_ui_element,
					p_http_api_map,
					p_log_fun);

				close();
				p_on_tags_created_fun(tags_lst);
      		}
		});
	
	$(tagging_input_ui_element).find('#submit_tags_button').on('click', async (p_event)=>{

			const tags_lst = await add_tags_to_obj(p_obj_id_str,
				p_obj_type_str,
				tagging_input_ui_element,
				p_http_api_map,
				p_log_fun);

			close();
			p_on_tags_created_fun(tags_lst);
		});

	//-----------------------------------------------------
	function close() {
		$(tagging_input_ui_element).remove();
		if (p_on_tag_ui_remove_fun != null) {
			p_on_tag_ui_remove_fun();
		}
	}

	//-----------------------------------------------------
	// TAG INPUT CLOSE BUTTON
	$(tagging_input_ui_element).find('#close_tagging_input_container_button').on('click',(p_event)=>{

		const tagging_input_container_element = $(p_event.target).parent();

		$(tagging_input_container_element).remove();
		if (p_on_tag_ui_remove_fun != null) {
			p_on_tag_ui_remove_fun();
		}
	});
	
	return tagging_input_ui_element;
}

//-----------------------------------------------------
export function place_tagging_input_ui_element(p_tagging_input_ui_element,
	p_relative_to_element,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_tagger_input_ui.place_tagging_input_ui_element()');
	
	$('body').append(p_tagging_input_ui_element);

	const relative_element__width_int = $(p_relative_to_element).width();
	const input_ui_element__width_int = $(p_tagging_input_ui_element).width();

	// p_tagging_input_ui_element.query('input').focus();
	//------------------------
	// Y_COORDINATE
	// document.body.scrollTop - is added to get the 'y' coord relative to the whole doc, regardless of amount of scrolling done
	// const relative_to_element_y_int :number = $(p_relative_to_element).offset().top + $('body').scrollTop(); //p_relative_to_element.getClientRects()[0].top.toInt() +	
	const relative_to_element_y_int :number = $(p_relative_to_element).offset().top;						
	
	//------------------------
	// X_COORDINATE
	const relative_to_element_x_int        :number = $(p_relative_to_element).offset().left;
	const input_ui_horizontal_overflow_int :number = (input_ui_element__width_int - relative_element__width_int)/2;

	var tagging_input_x :number;

	// input_ui is wider then target element
	if (input_ui_horizontal_overflow_int > 0) {

		// input_ui is cutoff on the left side
		if ((relative_to_element_x_int - input_ui_horizontal_overflow_int) < 0) {

			// position input_ui with its left side aligned with left edge of element to be tagged
			tagging_input_x = relative_to_element_x_int;
		}
		// input_ui is cutoff on the right side
		else if (((relative_to_element_x_int+relative_element__width_int) + input_ui_horizontal_overflow_int) > $(window).innerWidth()) {

			// position inpout_ui with its right edge aligned with the right edge of element to be tagged
			tagging_input_x = (relative_to_element_x_int+relative_element__width_int) - input_ui_element__width_int;
		}
		// no cutoff
		else {
			// positions that tag input container in the middle, and above, of the post_element
			tagging_input_x = relative_to_element_x_int-(input_ui_element__width_int-relative_element__width_int)/2;
		}
	}
	// input_ui is narrower then element, so just position normally
	else {
		// positions that tag input container in the middle, and above, of the post_element
		tagging_input_x = relative_to_element_x_int-(input_ui_element__width_int-relative_element__width_int)/2;
	}

	const tagging_input_y :number = relative_to_element_y_int - $(p_tagging_input_ui_element).height()/2;

	$(p_tagging_input_ui_element).css('position', 'absolute');
	$(p_tagging_input_ui_element).css('left',     tagging_input_x+'px');
	$(p_tagging_input_ui_element).css('top',      tagging_input_y+'px');
}

//-----------------------------------------------------
// TAGS SENDING TO SERVER
//-----------------------------------------------------
async function add_tags_to_obj(p_obj_id_str :string,
	p_obj_type_str :string,
	p_tagging_ui_element,
	p_http_api_map,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_tagger_input_ui.add_tags_to_obj()');
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const tags_str :string   = $(p_tagging_ui_element).find('#tags_input').val();
		const tags_lst :string[] = tags_str.split(' ');
		p_log_fun('INFO', 'tags_lst - '+tags_lst.toString());
		

		const existing_tags_lst :string[] = [];

		$(p_tagging_ui_element).parent().find('.tags_container').find('a').each((p_i, p_tag)=>{
			const tag_str = $(p_tag).text().trim();
			existing_tags_lst.push(tag_str);
		});

		
		const new_tags_lst :string[] = [];
		for (var tag_str of tags_lst) {

			// filter out only tags that are currently not existing/attached to this object
			if (!(tag_str in existing_tags_lst)) {
				new_tags_lst.push(tag_str);
			}
		}

		// ADD!! - some visual success/failure indicator
		const tags_meta_map = {};

		var data_map;
		if (p_http_api_map == null) {

			data_map = await gf_tagger_client.add_tags_to_obj(new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				tags_meta_map,
				p_log_fun);

		} else {

			data_map = await p_http_api_map["gf_tagger"]["add_tags_to_obj"](new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				tags_meta_map,
				p_log_fun);
		}

		const added_tags_lst :string[] = data_map['added_tags_lst'];
		p_log_fun('INFO', 'added_tags_lst:'+added_tags_lst);

		p_resolve_fun(added_tags_lst);
	});
	return p;
}