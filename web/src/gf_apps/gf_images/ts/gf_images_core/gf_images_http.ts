/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

//---------------------------------------------------
// CLASSIFY

export async function classify(p_image_id_str :string,
	p_client_type_str :string,
    p_log_fun         :Function) {
	
	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const url_str = `/v1/images/classify`;
		p_log_fun("INFO", `url_str ${url_str} - image_id ${p_image_id_str}`);

		const data_map = {
			"client_type_str": p_client_type_str,
			"images_ids_lst":  [p_image_id_str]
		};
		
		console.log(data_map)

		//-------------------------
		// HTTP AJAX
		$.post(url_str,
			JSON.stringify(data_map),
			function(p_data_map) {
				console.log("response received");
				// const data_map = JSON.parse(p_data);

				console.log(`data_map["status"] - ${p_data_map["status"]}`);
				
				if (p_data_map["status"] == "OK") {

					p_resolve_fun({

					});
				}
				else {
					p_reject_fun(p_data_map["data"]);
				}
			});

		//-------------------------
	});
}

//---------------------------------------------------
// SHARE

export async function share(p_image_id_str :string,
	p_email_address_str :string,
	p_email_subject_str :string,
	p_email_body_str    :string,
    p_log_fun           :Function) {
	
	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const url_str = `/v1/images/share`;
		p_log_fun("INFO", `url_str ${url_str} - image_id ${p_image_id_str}`);

		const data_map = {
			"image_id":      p_image_id_str,
			"email_address": p_email_address_str,
			"email_subject": p_email_subject_str,
			"email_body":    p_email_body_str
		};
		

		console.log(data_map)


		//-------------------------
		// HTTP AJAX
		$.post(url_str,
			JSON.stringify(data_map),
			function(p_data_map) {
				console.log("response received");
				// const data_map = JSON.parse(p_data);

				console.log(`data_map["status"] - ${p_data_map["status"]}`);
				
				if (p_data_map["status"] == "OK") {

					p_resolve_fun({

					});
				}
				else {
					p_reject_fun(p_data_map["data"]);
				}
			});

		//-------------------------
	});
}

//---------------------------------------------------
// GET

export async function get(p_image_id_str :string,
    p_log_fun :Function) {

	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const page_size_int = 10;
		const url_str       = `/v1/images/get?img_id=${p_image_id_str}`;
		p_log_fun("INFO", `url_str - ${url_str}`);

		//-------------------------
		// HTTP AJAX
		$.get(url_str,
			function(p_data_map) {
				console.log("response received");
				// const data_map = JSON.parse(p_data);

				console.log(`data_map["status"] - ${p_data_map["status"]}`);
				
				if (p_data_map["status"] == "OK") {

					const image_exists_bool = p_data_map["data"]["image_exists_bool"];
					const image_export_map  = p_data_map["data"]["image_export_map"];
					p_resolve_fun({
						"image_exists_bool": image_exists_bool,
						"image_export_map":  image_export_map,
					});
				}
				else {
					p_reject_fun(p_data_map["data"]);
				}
			});

		//-------------------------
	});
}

//---------------------------------------------------
// GET_PAGE

interface PageData {
    pages_lst            :any[];
    pages_user_names_lst :string[];
}

export async function get_page(p_flow_name_str :string,
	p_current_page_int :number,
	p_pages_num_int    :number,
	p_log_fun          :Function) :Promise<PageData> {

	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const page_size_int = 10;
		const url_str       = `/images/flows/browser_page?fname=${p_flow_name_str}&pg_index=${p_current_page_int}&pg_size=${page_size_int}&pg_num=${p_pages_num_int}`;
		p_log_fun("INFO", `url_str - ${url_str}`);

		//-------------------------
		// HTTP AJAX
		$.get(url_str,
			function(p_data_map) {

				console.log("response received...");
				console.log(`data_map["status"] - ${p_data_map["status"]}`);
				
				if (p_data_map["status"] == "OK") {

					const pages_lst = p_data_map["data"]["pages_lst"];
					const pages_user_names_lst = p_data_map["data"]["pages_user_names_lst"];

					p_resolve_fun({
						"pages_lst":            pages_lst,
						"pages_user_names_lst": pages_user_names_lst,
					});
				}
				else {
					p_reject_fun(p_data_map["data"]);
				}
			});

		//-------------------------
	});
}