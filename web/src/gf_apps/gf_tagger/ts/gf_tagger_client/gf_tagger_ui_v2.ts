/*
GloFlow application and media management/publishing platform
Copyright (C) 2025 Ivan Trajkovic

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
export function init(p_obj_id_str,
	p_obj_type_str,
	p_obj_elem,
	p_obj_parent_elem,
    
    p_callbacks_map,
    p_http_api_map,
	p_log_fun) {

	const tagging_input_ui_element = init_input_ui(p_obj_id_str,
		p_obj_type_str,
		p_callbacks_map,
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

		const position_relative_to_element = p_obj_elem;

		place_tags(tagging_input_ui_element,
			position_relative_to_element,
			p_obj_parent_elem,
			p_log_fun);


		if ("on_tagging_ui_add_fun" in p_callbacks_map) {
			p_callbacks_map["on_tagging_ui_add_fun"]();
		}

		// remove the initial controls when the full control opens
		$(tagging_ui_element).detach();
	});

	//------------------------
	// IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
	//               but unline mouseon/mouseout it will not fire if its children are entered/left.
	
	$(p_obj_elem).on('mouseenter', (p_event)=>{
		$(p_obj_elem).append(tagging_ui_element);
	});

	$(p_obj_elem).on('mouseleave', (p_event)=>{

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
}

//-------------------------------------------------
// TAGS
//-------------------------------------------------
function init_input_ui(p_obj_id_str,
	p_obj_type_str,
	p_callbacks_map,
	p_http_api_map,
	p_log_fun) {
	

	// IMPORTANT!! - "autofocus" added to input field, so that it receives keyboard focus as soon
	// 				 as its displayed.

	const input_ui_element = $(`
		<div id="tagging_input_container" class="bubble-in">
			<div id="background"></div>

			<input type="text" id="tags_input" placeholder="(space) separated tags" autofocus>
			<div id="submit_btn">+</div>


			<div id="generate_btn">G</div>
			<div id="generated_tags">
				
			</div>

			<div id="tagging_access_control">
				<div id="public_btn">public</div>
				<div id="private_btn">private</div>
			</div>
		</div>`);
	
	const tags_input_element = $(input_ui_element).find('#tags_input');

	//---------------------------
	// GENERATE_BTN

	const user_selected_gen_tags_lst: string[] = [];
	var generated_bool = false;

	$(input_ui_element).find("#generate_btn").on('click', async ()=>{

		// if tags have already been generated for the image, dont do it again
		if (generated_bool) {
			return;
		}

		const generated_tags_lst = await generate_tags(p_obj_id_str,
			p_callbacks_map,
			p_http_api_map,
			p_log_fun);

		const gen_tags_container_element = $(input_ui_element).find("#generated_tags");
		
		for (var tag_str of generated_tags_lst) {

			((p_tag_str)=>{
				const tag_element = $(`
					<div class="gf_tag tag_gen">
						${p_tag_str}
					</div>`);


				$(gen_tags_container_element).append(tag_element);
				
				var selected_bool = false;
				$(tag_element).on("click", (p_event)=>{

					if (!selected_bool) {

						// user has selected this tag to associate with obj
						user_selected_gen_tags_lst.push(p_tag_str);

						// mark element as selected
						$(tag_element).css('background-color', 'green');

						selected_bool = true;
					}
					else {
						
						selected_bool = false;
						$(tag_element).css('background-color', 'lightgray');
					}
				});
			})(tag_str);
		}

		generated_bool = true;
	});

	//---------------------------

	// 'ESCAPE' key
	$(document).on('keyup', (p_event)=>{
		if (p_event.which == 27) {

			// remove any previously present tagging_input_container's
			$(input_ui_element).detach();
			if ("on_tagging_ui_remove_fun" in p_callbacks_map) {
				p_callbacks_map["on_tagging_ui_remove_fun"]();
			}
		}
	});

	// to handlers for the same thing, one for the user clicking on the button,
	// the other for the user pressing 'enter'  
	$(tags_input_element).on('keyup', async (p_event)=>{

			// 'ENTER' key
			if (p_event.which == 13) {
				p_event.preventDefault();
				
				const tags_lst = await add_tags_to_obj_v2(p_obj_id_str,
					p_obj_type_str,
					user_selected_gen_tags_lst,
					input_ui_element,

					p_callbacks_map,
					p_http_api_map,
					p_log_fun);

				close();

				if ("tags_created_fun" in p_callbacks_map) {
					p_callbacks_map["tags_created_fun"](tags_lst);
				}
      		}
		});
	
	$(input_ui_element).find('#submit_btn').on('click', async (p_event)=>{

			p_event.stopImmediatePropagation();

			const tags_lst = await add_tags_to_obj_v2(p_obj_id_str,
				p_obj_type_str,
				user_selected_gen_tags_lst,
				input_ui_element,
				p_callbacks_map,
				p_http_api_map,
				p_log_fun);

			close();

			if ("tags_created_fun" in p_callbacks_map) {
				p_callbacks_map["tags_created_fun"](tags_lst);
			}
		});

	//-----------------------------------------------------

	$('body').on('click', (p_event)=>{

		p_event.stopImmediatePropagation();		
		close();
	});



	// stop propagation of events that reach the base container of the ui control
	$(input_ui_element).on('click', (p_event)=>{
		p_event.stopImmediatePropagation();
	});
	
	//-----------------------------------------------------
	function close() {

        // clear input field before closing, so its empty next time its oepend by the user
        $(input_ui_element).find("input").val("");

		$(input_ui_element).detach();

		if ("on_tagging_ui_remove_fun" in p_callbacks_map) {
			p_callbacks_map["on_tagging_ui_remove_fun"]();
		}
	}
	
	//-----------------------------------------------------
	return input_ui_element;
}

//-----------------------------------------------------
async function generate_tags(p_obj_id_str: string,
	p_callbacks_map: any,
	p_http_api_map: any,
	p_log_fun: (level: string, message: string) => void): Promise<string[]> {
	
	const p = new Promise<string[]>(async function(p_resolve_fun, p_reject_fun) {

		// HTTP
		const data_map = await p_http_api_map["gf_images"]["classify_image"](p_obj_id_str,
			p_log_fun);

		const generated_tags_lst: string[] = data_map['generated_tags_lst'];
		p_log_fun('INFO', `generated_tags_lst: ${generated_tags_lst}`);

		p_resolve_fun(generated_tags_lst);

	});
	return p;
}

//-----------------------------------------------------
async function add_tags_to_obj_v2(p_obj_id_str,
	p_obj_type_str,
	p_gen_tags_lst,
	p_tagging_ui_element,
	p_callbacks_map,
	p_http_api_map,
	p_log_fun) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		//------------------------
		// INPUT_FIELD_TAGS

		const input_field_tags_str = $(p_tagging_ui_element).find('#tags_input').val();
		var input_field_tags_lst = [];

		// split tags from input field only if there are any
		if (input_field_tags_str != undefined && input_field_tags_str != "") {
			input_field_tags_lst = input_field_tags_str.split(' ');
			
			p_log_fun('INFO', `input_field_tags_lst - ${input_field_tags_lst.toString()}`);
		}

		//------------------------
		// EXISTING_TAGS
		const existing_tags_lst = [];
		$(p_tagging_ui_element).parent().find('.tags_container').find('a').each((p_i, p_tag)=>{
			const tag_str = $(p_tag).text().trim();
			existing_tags_lst.push(tag_str);
		});

		//------------------------
		// FILTER_EXISTING_TAGS - only tags that are currently not existing/attached to this object

		const new_tags_lst = [];
		for (var tag_str of input_field_tags_lst) {
			if (!(tag_str in existing_tags_lst)) {
				new_tags_lst.push(tag_str);
			}
		}
		for (var tag_str of p_gen_tags_lst) {
			if (!(tag_str in existing_tags_lst)) {
				new_tags_lst.push(tag_str);
			}
		}

		//------------------------
		// RUN CREATE_PRE_HOOK
		const object_system_id_str = await p_callbacks_map["tags_pre_create_fun"](new_tags_lst);

		//------------------------
		// ADD!! - some visual success/failure indicator

		const tags_meta_map = {};

		// HTTP
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
function place_tags(p_input_ui_element,
	p_position_relative_to_element,
	p_obj_parent_elem,
	p_log_fun) {
	
	/*
	input element itself is attached to a different element, outside of this control. it could be "body",
	or some other parent.
	*/
	$(p_obj_parent_elem).append(p_input_ui_element);

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
	if ($(p_position_relative_to_element).css("top") == "auto" || $(p_position_relative_to_element).css("top") == "") {
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


	// INPUT_FOCUS - give browser focus to the input element, so that the user can
	//				 start entering tags right away.
	$(p_input_ui_element).find("input").focus();
}