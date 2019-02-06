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

import * as gf_images_stats from "../stats/gf_images_stats";

$(document).ready(()=>{
    //-------------------------------------------------
    function log_fun(p_g,p_m) {
        var msg_str = p_g+':'+p_m
        //chrome.extension.getBackgroundPage().console.log(msg_str);

        switch (p_g) {
            case "INFO":
                console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
                break;
            case "FUN_ENTER":
                console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
                break;
        }
    }
    //-------------------------------------------------
    init(log_fun);
});
//-------------------------------------------------
export function init(p_log_fun) {
	p_log_fun('FUN_ENTER','gf_images_dashboard.main()');

	//-------------------------
	//COMPLETED_IMAGE_JOBS
	$('#completed_image_jobs_runtime_infos .btn').on('click',()=>{
		gf_images_stats.http__get_stat_data('completed_image_jobs_runtime_infos',
			(p_stat_data_map)=>{
				const completed_image_jobs_runtime_infos_lst = p_stat_data_map['data_map']['completed_image_jobs_runtime_infos_lst'];
				if (completed_image_jobs_runtime_infos_lst.length > 0) {
					gf_images_stats.stat_view__completed_image_jobs_runtime_infos(completed_image_jobs_runtime_infos_lst, p_log_fun);
				}
			},
			p_log_fun);
	});
	//-------------------------
	//COMPLETED_IMAGE_JOBS
	$('#image_jobs_errors .btn').on('click',()=>{
		gf_images_stats.http__get_stat_data('image_jobs_errors',
			(p_stat_data_map)=>{
				const image_jobs_errors_lst = p_stat_data_map['data_map']['image_jobs_errors_lst'];
				if (image_jobs_errors_lst.length) {
					gf_images_stats.stat_view__image_jobs_errors(image_jobs_errors_lst, p_log_fun);
				}
			},
			p_log_fun);
	});
	//-------------------------
}