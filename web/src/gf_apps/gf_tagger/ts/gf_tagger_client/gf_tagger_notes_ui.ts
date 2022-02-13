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
export function init(p_obj_id_str :string,
	p_obj_type_str :string,
	p_obj_element,
	p_log_fun) {

	const notes_panel_btn = $(`
		<div id='notes_panel_btn'>
			<div class='icon'>notes</div>
		</div>`);
	$(p_obj_element).append(notes_panel_btn);


	const notes_panel = $(`
		<div id='notes_panel'>
			<div id='background'></div>

			<div id='container'>
				<div id='add_note_btn'>
					<div class='icon'>+</div>
				</div>
				<div id='notes'>
				</div>
			</div>
		</div>`);
		
	const background   = $(notes_panel).find('#background');
	const add_note_btn = $(notes_panel).find('#add_note_btn');

	var notes_open_bool = false;
	var notes_init_bool = false;
	$(notes_panel_btn).on('click', (p_event)=>{

		if (notes_open_bool) {

		}
		else {
			$(p_obj_element).append(notes_panel);

			//------------------------
			// GET
			if (!notes_init_bool) {

				//------------
				// NOTE_INPUT_PANEL
				const note_input_panel = $(`
					<div class='note_input_panel'>
						<textarea name="note_input" cols="30" rows="3"></textarea>
					</div>`);
				notes_panel.append(note_input_panel);
				
				//------------
				
				get_notes(p_obj_id_str,
					p_obj_type_str,
					notes_panel,
					()=>{
						notes_init_bool = true;
					},
					p_log_fun);				
			}

			//------------------------
			$(notes_panel).css('visibility',"visible");
		}		
	});

	//------------------------
	add_note_btn.on('click', (p_event)=>{
		run__remote_add_note(p_obj_id_str,
			p_obj_type_str,
			notes_panel,

			// p_onComplete_fun,
			()=>{
				//----------------------
				// GROW BACKGROUND
				const background_padding_size_int :number = 30;
				const notes_height_int            :number = $(notes_panel).find('#notes').height();
				const notes_y_int                 :number = $(notes_panel).find('#notes').offset().top;
				const new_height_int              :number = notes_y_int + notes_height_int + 2*background_padding_size_int;
				$(background).css('height', new_height_int+'px');
				
				//----------------------
			},
			p_log_fun);
	});

	//------------------------
	// IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
	//               but unline mouseon/mouseout it will not fire if its children are entered/left.
	$(p_obj_element).on('mouseenter', (p_event)=>{
		$(notes_panel_btn).css('visibility', 'visible');
	});

	$(p_obj_element).on('mouseleave', (p_event)=>{
		$(notes_panel_btn).css('visibility', 'hidden');
	});

	//------------------------
	// 'ESCAPE' key
	$(document).on('keyup', (p_event)=>{
		if (p_event.which == 27) {
			// remove any previously present note_input_container's
			$(notes_panel).remove();
		}
	});

	//------------------------
}

//-----------------------------------------------------
function get_notes(p_obj_id_str :string,
	p_obj_type_str :string,
	p_notes_panel,
	p_on_complete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER','gf_tagger_notes_ui.get_notes()');

	//------------------------
	// IMPORTANT!! - get notes via HTTP from backend gf_tagger_service
	gf_tagger_client.get_notes(p_obj_id_str,
		p_obj_type_str,
		// p_on_complete_fun
		(p_notes_lst :Object[])=>{

			for (var note_map of p_notes_lst) {

				const user_id_str :string = note_map['user_id_str'];
				const body_str    :string = note_map['body_str'];

				add_note_view(body_str,
					user_id_str,
					p_notes_panel,
					p_log_fun);
			}
			p_on_complete_fun();
		},
		()=>{}, // p_onError_fun
		p_log_fun);

	//------------------------	
}

//-----------------------------------------------------
function run__remote_add_note(p_obj_id_str :string,
	p_obj_type_str :string,
	p_notes_panel,
	p_on_complete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_tagger_notes_ui.run__remote_add_note()');

	const text_element          = $(p_notes_panel).find('.note_input_panel textarea');
	const note_body_str :string = $(text_element).val();
	p_log_fun('INFO', 'note_body_str        - $note_body_str');
	p_log_fun('INFO', 'note_body_str.length - ${note_body_str.length}');

	if (note_body_str.length > 0) {

		// ADD!! - some visual success/failure indicator
		gf_tagger_client.add_note_to_obj(note_body_str,
			p_obj_id_str,
			p_obj_type_str,
			()=>{
				add_note_view(note_body_str,
					'anonymouse', // p_anonymous_user_str,
					p_notes_panel,
					p_log_fun);
				
				$(text_element).val(''); //reset text
			},
			()=>{}, // p_onError_fun
			p_log_fun);
	}
}

//-----------------------------------------------------
function add_note_view(p_body_str :string,
	p_user_id_str :string,
	p_notes_panel,
	p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_tagger_notes_ui.add_note_view()');

	if (p_body_str.length > 20 ) console.log(p_body_str.substring(0,20)+'...');

	var short_body_str :string;
	if (p_body_str.length > 20 ) short_body_str = p_body_str.substring(0,20)+'...';
	else                         short_body_str = p_body_str;

	const new_note_element = $(`
		<div class='note'>
			<div class='icon'>n</div>
			<div class='details'>
				<div class='user'>`+p_user_id_str+`</div>
				<div class='body'>`+short_body_str+`</div>
			</div>
		</div>`);

	// other notes already exist
	if ($(p_notes_panel).find('#notes').children().length > 0) {
		const latest_note = $(p_notes_panel).find('#notes').children()[0];

		// insertBefore() - makes the new_note the first element in the list,
		//                  because the newest notes are at the top.
		/*$(p_notes_panel).find('#notes').insertBefore(new_note_element, //new_child
												latest_note);        //ref_child*/
		$(new_note_element).insertBefore(latest_note);
	}
	else {
		$(p_notes_panel).find('#notes').append(new_note_element);
	}

	$(new_note_element).css('opacity','0.0');

	const duration_int = 300;
	$(new_note_element).animate({'opacity':1.0}, duration_int, ()=>{});
}