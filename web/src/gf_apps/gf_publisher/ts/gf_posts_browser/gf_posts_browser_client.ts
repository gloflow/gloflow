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

namespace gf_posts_browser_client {
//-----------------------------------------------------
export function get_page(p_page_index_int :number,
    p_page_elements_num_int :number,
    p_on_complete_fun,
    p_on_error_fun,
    p_log_fun) {
	p_log_fun('FUN_ENTER','gf_posts_browser_client.get_page()');

	const url_str  = '/posts/browser_page';
	const data_map = {
		'pg_index': p_page_index_int,
		'pg_size':  p_page_elements_num_int
	};

    $.ajax({
        'url':         url_str,
        'type':        'GET',
        'data':        data_map,
        'contentType': 'application/json',
        'success':     (p_response_str)=>{

            const response_map = JSON.parse(p_response_str);
            const status_str   = response_map['status_str'];
            const page_lst :Object[] = response_map['data'];

            p_on_complete_fun(page_lst);
        },
        'error':(jqXHR, p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}
//-----------------------------------------------------
}