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

// ///<reference path="../../d/jquery.d.ts" />

import * as gf_events from "./gf_events";
import * as gf_user_events from "./../../../../gf_events/ts/gf_user_events";

//-------------------------------------------------
export function init(p_image_element :HTMLElement,
	p_image_id_str             :string,
	p_img_thumb_medium_url_str :string,
	p_img_thumb_large_url_str  :string,
	
	p_flows_names_lst      :string[],
	p_tags_lst             :string[],
	p_events_enabled_bool  :boolean,
	p_host_str			   :string,
	p_plugin_callbacks_map :any,
	p_log_fun              :Function) {

	$(p_image_element).find("img").click(()=>{

		const gf_link_str = `/images/v/${p_image_id_str}`;

		const image_view_element = $(`
			<div id="image_viewer" class="gf_center">
				<div id="background"></div>
				<div id="main">
					<div id="image_container">
						<img src="${p_img_thumb_large_url_str}"></img>
						<div id="tags">
						
						</div>
					</div>
					<div id="image_details">
						
						<div id="flows_names" class="gf_center">
							
						</div>
						
						<div id="image_view_link" class="gf_center">
							<a href="${gf_link_str}">gf link</a>
						</div>
					</div>
				</div>
			</div>`)[0];

		$('body').append(image_view_element);

		const bg                      = $(image_view_element).find("#background");
		const main_element            = $(image_view_element).find("#main")[0];
		const image_container_element = $(image_view_element).find("#image_container")[0];
		const flows_names_element     = $(image_view_element).find("#flows_names")[0];
		const tags_element            = $(image_view_element).find("#tags")[0];

		p_flows_names_lst.forEach((p_flow_name_str :string)=>{
			$(flows_names_element).append(`
				<div class="flow_name">
					<a href="https://gloflow.com/images/flows/browser?fname=${p_flow_name_str}">${p_flow_name_str}</a>
				</div>
			`);
		});

		p_tags_lst.forEach((p_tag_str :string)=>{
			$(tags_element).append(`
				<a class="tag" href="https://gloflow.com/v1/tags/objects?tag=${p_tag_str}&otype=image">#${p_tag_str}</a>
			`);
		});

		// position the image vertically where the user has scrolled to
		$(image_view_element).css("top", `${$(window).scrollTop()}px`);

		//----------------------
		// IMPORTANT!! - turn off vertical scrolling while viewing the image
		$("body").css("overflow", "hidden");

		//----------------------
		// IMG_ONLOAD
		$(image_view_element).find("img").on("load", ()=>{
			$(main_element).css("visibility", "visible");
		});

	    //----------------------
	    $(bg).click(async ()=>{
	    	$(image_view_element).remove();

	    	// turn vertical scrolling back on when done viewing the image
	    	$("body").css("overflow", "auto");

			//------------------
            // EVENTS
            if (p_events_enabled_bool) {
                
                const event_meta_map = {
					"image_id": p_image_id_str,
                };
                await gf_user_events.send_event_http(gf_events.GF_IMAGES_IMAGE_VIEWER_CLOSE,
                    "browser",
                    event_meta_map,
                    p_host_str)
            }

            //------------------
	    });

	    //----------------------

		//--------------------------
		// EVENTS
		if (p_events_enabled_bool) {
			
			const event_meta_map = {
				"image_id": p_image_id_str,
			};
			gf_user_events.send_event_http(gf_events.GF_IMAGES_IMAGE_VIEWER_OPEN,
				"browser",
				event_meta_map,
				p_host_str)
		}

		//--------------------------
		// PLUGIN_CALLBACK

		if ("image_viewer_open" in p_plugin_callbacks_map) {
			p_plugin_callbacks_map["image_viewer_open"](image_view_element);
		}

		//--------------------------
	});
}

//-------------------------------------------------