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

//-------------------------------------------------
function gf_tagger__init_ui(p_obj_type_str,
	p_obj_element,
	p_input_element_parent_selector_str,

	p_tags_create_pre_fun,
	p_on_tags_created_fun,

	p_notes_create_pre_fun,
	p_on_notes_created_fun,

	p_on_tagging_ui_add_fun,
	p_on_tagging_ui_remove_fun,
	p_http_api_map,
	p_log_fun) {

    console.log("gf_tagger UI init...")

    const tagging_input_ui_element = gf_tagger__init_input_ui(p_obj_type_str,
		p_tags_create_pre_fun,
		p_on_tags_created_fun,
		p_on_tagging_ui_remove_fun,
		p_http_api_map,
		p_log_fun);

	const notes_input_ui_element = gf_tagger__init_notes_input_ui(p_obj_type_str,
		p_notes_create_pre_fun,
		p_on_notes_created_fun,
		p_on_tagging_ui_remove_fun,
		p_http_api_map,
		p_log_fun);
		
	const tagging_ui_element = $(`
		<div class="tagging_controls">
			<div class="add_tags_button">add tags</div>
			<div class="add_notes_button">add notes</div>
		</div>`);

	//------------------------------
	// OPEN TAG INPUT UI
	$(tagging_ui_element).find('.add_tags_button').on('click', (p_event)=>{

		p_event.stopImmediatePropagation();

		// remove the tagging_input_container if its already displayed
		// for tagging another post_element
		if ($('#tagging_input_container') != null) {
			$('#tagging_input_container').detach();
		}

		const position_relative_to_element = p_obj_element;

		gf_tagger__place_tags_input_ui(tagging_input_ui_element,
			position_relative_to_element,
			p_input_element_parent_selector_str,
			p_log_fun);

		if (p_on_tagging_ui_add_fun != null) p_on_tagging_ui_add_fun();

		// remove the initial controls when the full control opens
		$(tagging_ui_element).detach();
	});

	//------------------------------
	// OPEN NOTES INPUT UI
	$(tagging_ui_element).find('.add_notes_button').on('click', (p_event)=>{

		p_event.stopImmediatePropagation();

		// remove the tagging_input_container if its already displayed
		// for tagging another post_element
		if ($('#tagging_input_container') != null) {
			$('#tagging_input_container').detach();
		}

		const position_relative_to_element = p_obj_element;

		gf_tagger__place_tags_input_ui(notes_input_ui_element,
			position_relative_to_element,
			p_input_element_parent_selector_str,
			p_log_fun);

		if (p_on_tagging_ui_add_fun != null) p_on_tagging_ui_add_fun();

		// remove the initial controls when the full control opens
		$(tagging_ui_element).detach();
	});

	//------------------------
	// IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
	//               but unline mouseon/mouseout it will not fire if its children are entered/left.
	
	$(p_obj_element).on('mouseenter', (p_event)=>{
		$(p_obj_element).append(tagging_ui_element);
	});

	$(p_obj_element).on('mouseleave', (p_event)=>{

        // IMPORTANT!! - detaching in order to keep event handlers
		$(tagging_ui_element).detach();

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

    /*
	//'T' key - open tagging UI to the element that has the cursor 
	//          hovering over it
	final subscription = document.onKeyUp.listen((p_event) {
		if (p_event.keyCode == 84) {

			//remove the tagging_input_container if its already displayed
			//for tagging another post_element
			if (query('#tagging_input_container') != null) {
				query('#tagging_input_container').remove();
			}

			gf_tagger__place_tags_input_ui(tagging_input_ui_element,
                p_obj_element, // post_element,
                p_log_fun);

			//prevent this handler being invoked while the user
			//is typing in tags into the input field
			//subscription.pause();
		}
	});
	*/
}

//-------------------------------------------------
// NOTES
//-------------------------------------------------
function gf_tagger__init_notes_input_ui(p_obj_type_str,
	p_notes_create_pre_fun,
	p_on_notes_created_fun,
	p_on_tagging_ui_remove_fun,
	p_http_api_map,
	p_log_fun) {

	const input_ui_element = $(`
		<div id='notes_panel'>
			<div id='background'></div>

			<div id='container'>
				<div class='note_input_panel'>
					<textarea id="note_input" cols="30" rows="3"></textarea>
				</div>

				<div id='add_note_btn'>
					<div class='icon'>+</div>
				</div>
				<div id='notes'>
				</div>
			</div>
		</div>`);
	
	// 'ESCAPE' key
	$(document).on('keyup', (p_event)=>{
		if (p_event.which == 27) {

			// remove any previously present tagging_input_container's
			$(input_ui_element).detach();
			if (p_on_tagging_ui_remove_fun != null) {
				p_on_tagging_ui_remove_fun();
			}
		}
	});

	const tags_input_element = $(input_ui_element).find('#note_input');

	$(input_ui_element).find('#add_note_btn').on('click', async (p_event)=>{

		p_event.stopImmediatePropagation();
			
		const note_str = await add_note_to_obj(p_obj_type_str,
			p_notes_create_pre_fun,
			input_ui_element,
			p_http_api_map,
			p_log_fun);

		close();
		p_on_notes_created_fun(note_str);
	});

	return input_ui_element;

	//-----------------------------------------------------
	function close() {

        // clear input field before closing, so its empty next time its oepend by the user
        $(input_ui_element).find("input").val("");

		$(input_ui_element).detach();
		if (p_on_tagging_ui_remove_fun != null) {
			p_on_tagging_ui_remove_fun();
		}
	}

	//-----------------------------------------------------
}

//-------------------------------------------------
// TAGS
//-------------------------------------------------
function gf_tagger__init_input_ui(p_obj_type_str,
	p_tags_create_pre_fun,
	p_on_tags_created_fun,
	p_on_tagging_ui_remove_fun,
	p_http_api_map,
	p_log_fun) {
	

	const input_ui_element = $(`
		<div id="tagging_input_container" class="bubble-in">
			<div id="background"></div>
			<input type="text" id="tags_input" placeholder="(space) separated tags">
			<div id="submit_btn">add</div>
			<div id="close_btn">&#10006;</div>
		</div>`);
	
	const tags_input_element = $(input_ui_element).find('#tags_input');

	// 'ESCAPE' key
	$(document).on('keyup', (p_event)=>{
		if (p_event.which == 27) {

			// remove any previously present tagging_input_container's
			$(input_ui_element).detach();
			if (p_on_tagging_ui_remove_fun != null) {
				p_on_tagging_ui_remove_fun();
			}
		}
	});

	// to handlers for the same thing, one for the user clicking on the button,
	// the other for the user pressing 'enter'  
	$(tags_input_element).on('keyup', async (p_event)=>{

			// 'ENTER' key
			if (p_event.which == 13) {
				p_event.preventDefault();
				
				const tags_lst = await add_tags_to_obj(p_obj_type_str,
					input_ui_element,

					p_tags_create_pre_fun,
					p_http_api_map,
					p_log_fun);

				close();
				p_on_tags_created_fun(tags_lst);
      		}
		});
	
	$(input_ui_element).find('#submit_btn').on('click', async (p_event)=>{

			p_event.stopImmediatePropagation();

			const tags_lst = await add_tags_to_obj(p_obj_type_str,
				input_ui_element,
				p_tags_create_pre_fun,
				p_http_api_map,
				p_log_fun);

			close();
			p_on_tags_created_fun(tags_lst);
		});

	//-----------------------------------------------------
	function close() {

        // clear input field before closing, so its empty next time its oepend by the user
        $(input_ui_element).find("input").val("");

		$(input_ui_element).detach();
		if (p_on_tagging_ui_remove_fun != null) {
			p_on_tagging_ui_remove_fun();
		}
	}

	//-----------------------------------------------------
	// TAG INPUT CLOSE BUTTON
	$(input_ui_element).find('#close_btn').on('click', (p_event)=>{

		p_event.stopImmediatePropagation();

		const tagging_input_container_element = $(p_event.target).parent();

		$(tagging_input_container_element).detach();
		if (p_on_tagging_ui_remove_fun != null) {
			p_on_tagging_ui_remove_fun();
		}
	});
	
	return input_ui_element;
}

//-----------------------------------------------------
function gf_tagger__place_tags_input_ui(p_input_ui_element,
	p_position_relative_to_element,
	p_input_element_parent_selector_str,
	p_log_fun) {
	
	/*
	input element itself is attached to a different element, outside of this control. it could be "body",
	or some other parent.
	*/
	$(p_input_element_parent_selector_str).append(p_input_ui_element);

	const relative_element__width_int = $(p_position_relative_to_element).width();
	const input_ui_element__width_int = $(p_input_ui_element).width();


	console.log("input ui element width", input_ui_element__width_int)
	//------------------------
	// Y_COORDINATE

	var relative_to_element_y_int;

	/*
	IMPORTANT!! - some elements dont have the css "top" property set, its computed instead by the browser
		as a result of other styles and elements.
		for these situations the elements "top" property is either set to "auto", or when checked via style.top
		its set to "".
		in those cases the computed offset().top value has to be used.
	*/
	/*
	IMPORTANT!! - using css("top") instead of $(p_position_relative_to_element).offset().top because
		with masonry which sets the css("top") property offset().top doesnt return the correct value.
		css("top") also works correctly in the test cases, so using that for now.
	*/
	if ($(p_position_relative_to_element).css("top") == "auto" || p_position_relative_to_element.style.top == "") {
		relative_to_element_y_int = $(p_position_relative_to_element).offset().top;	
	} else {

		/*
		for other elements who's "top" is explicitly set, using css("top", 10) is the more precise way
		(it leads to accurate values that reflect actual position, while offset().top in those situations
		can be give the wrong value.
		*/
		relative_to_element_y_int = parseInt($(p_position_relative_to_element).css("top"), 10);	
	}
						
	const tagging_input_y = relative_to_element_y_int; // $(p_input_ui_element).height()/2;

	//------------------------
	// X_COORDINATE
	const relative_to_element_x_int        = $(p_position_relative_to_element).offset().left;
	const input_ui_horizontal_overflow_int = (input_ui_element__width_int - relative_element__width_int)/2;

	var tagging_input_x;

	// input_ui is wider then target element
	if (input_ui_horizontal_overflow_int > 0) {

		// input_ui is cutoff on the left side
		if ((relative_to_element_x_int - input_ui_horizontal_overflow_int) < 0) {

			console.log("left cutoff...")
			// position input_ui with its left side aligned with left edge of element to be tagged
			tagging_input_x = relative_to_element_x_int;
		}
		// input_ui is cutoff on the right side
		else if (((relative_to_element_x_int+relative_element__width_int) + input_ui_horizontal_overflow_int) > $(window).innerWidth()) {

			console.log("right cutoff...")

			// position inpout_ui with its right edge aligned with the right edge of element to be tagged
			tagging_input_x = (relative_to_element_x_int+relative_element__width_int) - input_ui_element__width_int;
		}
		// no cutoff
		else {

			console.log("no cutoff...")

			// positions that tag input container in the middle, and above, of the post_element
			tagging_input_x = relative_to_element_x_int-(input_ui_element__width_int-relative_element__width_int)/2;
		}
	}
	// input_ui is narrower then element, so just position normally
	else {

		console.log("regular positioning...")

		// positions that tag input container in the middle, and above, of the post_element
		tagging_input_x = relative_to_element_x_int-(input_ui_element__width_int-relative_element__width_int)/2;
	}

	//------------------------
	$(p_input_ui_element).css('position', 'absolute');
	$(p_input_ui_element).css('left',     `${tagging_input_x}px`);
	$(p_input_ui_element).css('top',      `${tagging_input_y}px`);
}

//-----------------------------------------------------
async function add_tags_to_obj(p_obj_type_str,
	p_tagging_ui_element,
	p_tags_create_pre_fun,
	p_http_api_map,
	p_log_fun) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const tags_str = $(p_tagging_ui_element).find('#tags_input').val();
		const tags_lst = tags_str.split(' ');
		p_log_fun('INFO', `tags_lst - ${tags_lst.toString()}`);
		
		const existing_tags_lst = [];

		$(p_tagging_ui_element).parent().find('.tags_container').find('a').each((p_i, p_tag)=>{
			const tag_str = $(p_tag).text().trim();
			existing_tags_lst.push(tag_str);
		});

		
		const new_tags_lst = [];
		for (var tag_str of tags_lst) {

			// filter out only tags that are currently not existing/attached to this object
			if (!(tag_str in existing_tags_lst)) {
				new_tags_lst.push(tag_str);
			}
		}

		//------------------------
		// CREATE_PRE_HOOK
		const object_system_id_str = await p_tags_create_pre_fun(new_tags_lst);

		//------------------------

		// ADD!! - some visual success/failure indicator
		const tags_meta_map = {};

		const data_map = await p_http_api_map["gf_tagger"]["add_tags_to_obj"](new_tags_lst,
			object_system_id_str,
			p_obj_type_str,
			tags_meta_map,
			p_log_fun);

		const added_tags_lst = data_map['added_tags_lst'];
		p_log_fun('INFO', `added_tags_lst: ${added_tags_lst}`);

		p_resolve_fun(added_tags_lst);
	});
	return p;
}

//-----------------------------------------------------
async function add_note_to_obj(p_obj_type_str,
	p_tags_create_pre_fun,
	p_tagging_ui_element,
	p_http_api_map,
	p_log_fun) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		console.log("AAAAAAAAAAAAA", p_tagging_ui_element)

		const note_str = $(p_tagging_ui_element).find('#note_input').val();
		p_log_fun('INFO', `note_str - ${note_str}`);



		p_resolve_fun(note_str);

	});
	return p;
}