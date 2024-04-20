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

declare var gf_tagger__init_ui_v2;
declare var gf_tagger__http_add_tags_to_obj;

//---------------------------------------------------
// TAGGING_UI

export function init_tagging(p_post_id_str :string,
	p_post_element :HTMLDivElement,
	p_http_api_map,
	p_log_fun) {

	/*
	const http_api_map = {

		// GF_TAGGER
		"gf_tagger": {
			"add_tags_to_obj": async (p_new_tags_lst,
				p_obj_id_str,
				p_obj_type_str,
				p_tags_meta_map,
				p_log_fun)=>{
				const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

					await gf_tagger__http_add_tags_to_obj(p_new_tags_lst,
						p_obj_id_str,
						p_obj_type_str,
						{}, // meta_map
						p_gf_host_str,
						p_log_fun);

					p_resolve_fun({
						"added_tags_lst": p_new_tags_lst,
					});
				});
				return p;
			}
		}
	};
	*/
	
	const obj_type_str = "post";

	const callbacks_map = {

		//---------------------------------------------------
		// TAGS
		//---------------------------------------------------
		"tags_pre_create_fun": async (p_tags_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// passing the image_id to the gf_tagger control via this callback allows for
				// customization of the image_id fetching mechanism (whether its in the template,
				// or fetched via rest api, etc., or pulled from some internal browser/web DB).
				p_resolve_fun(p_post_id_str);
			});
			return p;
		},
		
		//---------------------------------------------------
		"tags_created_fun": (p_tags_lst)=>{

			console.log("added tags >>>>>>>>>>>", p_tags_lst);

			p_tags_lst.forEach(p_tag_str=>{

				tag_display(p_tag_str);
			});
		},

		//---------------------------------------------------
		// NOTES
		//---------------------------------------------------
		"notes_pre_create_fun": (p_notes_lst)=>{
			const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

				// passing the image_id to the gf_tagger control via this callback allows for
				// customization of the image_id fetching mechanism (whether its in the template,
				// or fetched via rest api, etc., or pulled from some internal browser/web DB).
				p_resolve_fun(p_post_id_str);
			});
			return p;
		},

		//---------------------------------------------------
		"notes_created_fun": (p_notes_lst)=>{

			console.log("added notes >>>>>>>>>>>", p_notes_lst)
		}

		//---------------------------------------------------
	}

	gf_tagger__init_ui_v2(p_post_id_str,
		obj_type_str,
		p_post_element,
		$("body"),
		callbacks_map,
		p_http_api_map,
		p_log_fun);

	//-------------------------------------------------
	function tag_display(p_tag_str) {

		/*
		check if the tags_container div exists, if not create it.
		the backend template has a div with class "tags_container" in the image container only if the image
		has tags. if it does not, the .tags_container div is not created
		*/
		if ($(p_post_element).find(".tags_container").length == 0) {
			$(p_post_element).append("<div class='tags_container'></div>");
		}

		$(p_post_element)
			.find(".tags_container")
			.append(`<a class='gf_image_tag' href='/v1/tags/objects?tag=${p_tag_str}&otype=image'>#${p_tag_str}</a>`)
	}

	//-------------------------------------------------
}