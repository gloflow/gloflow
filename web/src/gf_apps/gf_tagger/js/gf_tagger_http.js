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

//-----------------------------------------------------
async function gf_tagger__http_add_tags_to_obj(p_tags_lst,  
    p_object_system_id_str,
    p_object_type_str,
    p_meta_map,
	p_host_str,
    p_log_fun) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = `${p_host_str}/v1/tags/create`;
		const tags_str = p_tags_lst.join(' ');
		const data_map = {
			"otype": p_object_type_str,
			"o_id":  p_object_system_id_str,
			"tags":  tags_str,
			"meta_map": p_meta_map,
		};

		const response = await fetch(url_str, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data_map)
		});

		if (response.ok) {
			const responseMap = await response.json();
			const status_str = responseMap["status"];
			const data_map   = responseMap["data"];
			
			if (status_str === "OK") {
				return p_resolve_fun(data_map);
			} else {
				return p_reject_fun(data_map);
			}
		} else {
			return p_reject_fun(`Fetch failed: ${response.status} ${response.statusText}`);
		}
	});
	return p;
}

//-----------------------------------------------------
async function gf_tagger__http_get_objs_with_tag(p_tag_str, 
    p_object_type_str,
	p_host_str,
    p_log_fun) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const url_str = `${p_host_str}/v1/tags/objects?tags=${p_tag_str}&otype=${p_object_type_str}`;

		const response = await fetch(url_str, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (response.ok) {
			const data_str = await response.text();
			const data_map = JSON.parse(data_str);
			const objects_with_tags_map = data_map['objects_with_tags_map'];
			
			return Promise.resolve(objects_with_tags_map);
		} else {
			return Promise.reject(`Fetch failed: ${response.status} ${response.statusText}`);
		}
	});
	return p;
}