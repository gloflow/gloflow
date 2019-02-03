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

namespace gf_crawl_stats__errors {

//---------------------------------------------------
export function init(p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_crawl_stats__errors.init()");

	$('#view_crawl_errors_btn').on('click', (p_stats_map)=>{

		gf_stats.http__stats_query('errors',
			(p_stats_map)=>{
				const errors_stats = view__errors_stats(p_stats_map,p_log_fun);

				$('#crawl_errors #errors').append(errors_stats);
			},
			()=>{},
			p_log_fun);
	});
}
//---------------------------------------------------
function view__errors_stats(p_stats_map, p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_crawl_stats__errors.view__errors_stats()");

	const stats = $(`<div id="errors_stats"></div>`);

	const stat_errors_lst = p_stats_map['stat_errors_lst'];
	for (var crawler_errors_map of stat_errors_lst) {

		const crawler_name_str = crawler_errors_map['crawler_name_str'];
		const errors_types_lst = crawler_errors_map['errors_types_lst'];

		const crawler_errors = $(`<div class="crawler_errors">
				<div class='crawler_name'>`+crawler_name_str+`</div>
			</div>`);

		$(stats).append(crawler_errors);

		for (var errors_type_map of errors_types_lst) {
			const error_type_str = errors_type_map['type_str'];
			const count_int      = errors_type_map['count_int'];
			const urls_lst       = errors_type_map['urls_lst'];

			const errors_type = $(`
				<div class='errors_type'>
					<div><span class='label'>type:</span><span>`+error_type_str+`</span></div>
					<div><span class='label'>count:</span><span>`+count_int+`</span></div>
					<div class='urls_lst'></div>
				</div>`);
			
			for (var url_str of urls_lst) {
				$(errors_type).find('.urls_lst').append($('<div class="url"><a target="_blank" href="'+url_str+'">'+url_str+'</a></div>'))
			}

			$(crawler_errors).append(errors_type);
		}        
	}

	return stats;
}
//---------------------------------------------------
}