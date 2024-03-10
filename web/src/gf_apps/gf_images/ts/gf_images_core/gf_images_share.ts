/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

import * as gf_images_http from "./gf_images_http";

//---------------------------------------------------
export function init(p_image_id_str,
	p_image_container_element,
	p_gf_host_str,
	p_log_fun) {



	const http_api_map = {
		"gf_images": {
			"share_image": async (p_image_id_str,
				p_email_address_str,
                p_email_subject_str,
                p_email_body_str,
				p_log_fun)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

					await gf_images_http.share(p_image_id_str,
                        p_email_address_str,
                        p_email_subject_str,
                        p_email_body_str,
						p_gf_host_str,
						p_log_fun);

					p_resolve_fun({
						
					});
				});
				return p;
			}
		}
	};

    // UI
    const control_element = init_ui(p_image_id_str,
		p_image_container_element,
		http_api_map,
        p_log_fun);


    $(p_image_container_element).append(control_element);
}

//---------------------------------------------------

function init_ui(p_image_id_str,
	p_image_container_element,
	p_http_api_map,
    p_log_fun) {





    const control_element = $(`
        <div class="gf_images_share">
            <div class="email_btn">
				S
            </div>
        </div>`);


	var opened_bool = false;
	var sharing_dialog_element;

	$(control_element).click(()=>{

		if(opened_bool){
			$(sharing_dialog_element).remove();
			opened_bool = false;
		}
		else {
			sharing_dialog_element = email_share_dialog();

			$(control_element).append(sharing_dialog_element);
			opened_bool = true;
		}

	});

	//------------------------
	// IMPORTANT!! - onMouseEnter/onMouseLeave fire when the target element is entered/left, 
	//               but unline mouseon/mouseout it will not fire if its children are entered/left.
	
	$(p_image_container_element).on('mouseenter', (p_event)=>{
		$(p_image_container_element).append(control_element);
	});

	$(p_image_container_element).on('mouseleave', (p_event)=>{

        // IMPORTANT!! - detaching in order to keep event handlers
		$(control_element).detach();
	});

	//---------------------------------------------------
	function email_share_dialog() {

		const dialog = $(`
		<div class="gf_email_dialog">
		
			<div class="email_address">
				<input type="text" placeholder="email address">
			</div>
			<div class="email_subject">
				<input type="text" placeholder="email subject">
			</div>
			<div class="email_body">
				<textarea placeholder="email body"></textarea>
			</div>
			<div class="share_btn">
				share
			</div>
		</div>`);


		$(dialog).find(".share_btn").click(async ()=>{

			const email_address_str = $(dialog).find(".email_address input").val();
			const email_subject_str = $(dialog).find(".email_subject input").val();
			const email_body_str = $(dialog).find(".email_body textarea").val();

			// HTTP
			await p_http_api_map.gf_images.share_image(p_image_id_str,
				email_address_str,
				email_subject_str,
				email_body_str,
				p_log_fun);
			
			// mark the button as complete
			$(dialog).find(".share_btn").css("background-color", "green");

			setTimeout(()=>{
					$(dialog).remove();
				}, 1000);
		});


		return dialog;
	}

	//---------------------------------------------------

    return control_element;
}