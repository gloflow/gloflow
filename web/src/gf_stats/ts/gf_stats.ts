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

///<reference path="../../d/jquery.d.ts" />

namespace gf_stats {
//---------------------------------------------------
export function init__batch(p_log_fun) {
	p_log_fun("FUN_ENTER","gf_stats.init__batch()");
	
	const container = $(`
		<div id="batch_stats_list_container">
			<select id='batch_stats_list'></select>
		</div>
	`);

	const select_e = $(container).find('#batch_stats_list');
	//----------------
	http__batch_plots_list((p_stats_list_lst)=>{
			for (let stat_name_str of p_stats_list_lst) {
				const stat_img_name_str = stat_name_str+'.png';
				const option_e          = $('<option value="'+stat_img_name_str+'">'+stat_name_str+'</option>'); 
				$(select_e).append(option_e);
			}
		},
		()=>{},
		p_log_fun);
	//----------------

	$(select_e).change(()=>{
		

		const stat_img_name_str = $(container).find('#batch_stats_list :selected').val();

		//IMPORTANT!! - adding unix time at the end of URL since its always different and will
		//              force the browser not to use the cached version.
		const unix_time_f      = Math.floor(Date.now() / 1000);
		const stat_img_url_str = '/a/static/plots/'+stat_img_name_str+'?t='+unix_time_f;

		$('#batch_stats').append(`
			<div class='batch_stat_img'>
				<img src='`+stat_img_url_str+`'></img>
			</div>`);
	})

	return container;
}

//---------------------------------------------------
function http__batch_plots_list(p_on_complete_fun, p_on_error_fun, p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_stats.http__batch_plots_list()");

	$.get('/a/stats/batch/list', (p_data_map)=>{

		console.log('response received');
		//const data_map = JSON.parse(p_data);

		if (p_data_map["status_str"] == 'OK') {
			const stats_list_lst = p_data_map["data"]['stats_list_lst'];
			p_on_complete_fun(stats_list_lst);
		}
		else {
			p_on_error_fun(p_data_map["data"]);
		}
	});
}
//---------------------------------------------------
//IMPORTANT!! - its called run_stats_fun because its actually running a data query function on the server,
//              its not pulling pre-computed stats results from a DB, or pulling a plotted stat image 
//              for display on the client.

export function http__stats_query(p_stat_name_str :string,
	p_on_complete_fun,
	p_on_error_fun,
	p_log_fun) {
	p_log_fun("FUN_ENTER","gf_stats.http__stats_query()");

	const url_str = '/a/stats/query';
	p_log_fun('INFO','url_str - '+url_str);

	//-------------------------
	//HTTP AJAX
	$.post(url_str,
		JSON.stringify({
			"stat_name_str": p_stat_name_str
		}),
		function(p_data_map) {
			console.log('response received');
			//const data_map = JSON.parse(p_data);
			
			if (p_data_map["status_str"] == 'OK') {

				const result_data_map = p_data_map["data"]['result_data_map'];
				p_on_complete_fun(result_data_map);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
	//------------------------- 
}
//---------------------------------------------------
}