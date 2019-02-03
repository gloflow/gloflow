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

///<reference path="../d/jquery.timeago.d.ts" />

namespace gf_crawl_stats {



//-------------------------------------------------
export function init__queries(p_log_fun) {
	p_log_fun("FUN_ENTER", "gf_crawl_stats.init__queries()");

	//---------------------
	//ERRORS

	gf_crawl_stats__errors.init(p_log_fun);
	//---------------------
	//FETCHES

	$('#get_new_fetches_per_day_stats_btn').on('click', ()=>{

		gf_stats.http__stats_query("fetches_by_days",
			(p_stats_map)=>{

				const stats__fetches_by_days_map = p_stats_map["stats__fetches_by_days_map"];
				const stats_container            = $(`<div class='stats_container'></div>`);
				const parent_element             = stats_container;
				$('body').append(stats_container);

				gf_crawl_stats__fetches.view__fetches_per_day(stats__fetches_by_days_map, parent_element, p_log_fun)
			},
			()=>{},
			p_log_fun);
	});
	//---------------------
	//LINKS
	$('#get_new_links_per_day_stats_btn').on('click', ()=>{

		gf_stats.http__stats_query("new_links_per_day",
			(p_stats_map)=>{
					
				const new_links_per_day_lst = p_stats_map['new_links_per_day_lst'];
					
				const stats_container = $(`<div class='stats_container'></div>`);
				const parent_element  = stats_container;
				$('body').append(stats_container);

				gf_crawl_stats__links.view__new_links_per_day(new_links_per_day_lst, parent_element, p_log_fun);
			},
			()=>{},
			p_log_fun);
	});
	
	$('#get_unresolved_links_stats_btn').on('click', ()=>{

		gf_stats.http__stats_query("unresolved_links",
			(p_stats_map)=>{
					
				const stat_unresolved_links_lst = p_stats_map['stat_unresolved_links_lst'];
				const stats                     = gf_crawl_stats__links.view__unresolved(stat_unresolved_links_lst, p_log_fun);

				const stats_container = $(`<div class='stats_container'></div>`);
				$(stats_container).append(stats);
				$('body').append(stats_container);
			},
			()=>{},
			p_log_fun);
	});
	
	//---------------------
	//GIF
	$('#get_gif_main_stats_btn').on('click', ()=>{

		gf_stats.http__stats_query("gif",
			(p_stats_map)=>{
				const gifs_stats = gf_crawl_stats__images.view__gif_stats(p_stats_map, p_log_fun);
				
				const stats_container = $(`<div class='stats_container'></div>`);
				$(stats_container).append(gifs_stats);
				$('body').append(stats_container);
			},
			()=>{},
			p_log_fun);
	});

	$('get_gifs_per_day_stats_btn').on('click', ()=>{

		gf_stats.http__stats_query("gifs_per_day",
			(p_stats_map)=>{

				const stats__gifs_by_days_map = p_stats_map["stats__gifs_by_days_map"];
				gf_crawl_stats__images.view__gifs_per_day_stats(stats__gifs_by_days_map, p_log_fun);
			},
			()=>{},
			p_log_fun);
	});
	//---------------------
	$('#get_crawled_images_domains_btn').on('click',()=>{

		gf_stats.http__stats_query('crawled_images_domains',
			(p_stats_map)=>{


				const crawled_images_domains_lst = p_stats_map['crawled_images_domains_lst'];
				const images_domains_e           = gf_crawl_stats__images.view__crawled_images_domains(crawled_images_domains_lst, p_log_fun);
				const stats_container            = $(`<div class='stats_container'></div>`);
				$(stats_container).append(images_domains_e);
				$('body').append(stats_container);
			},
			()=>{},
			p_log_fun);
	});
	//---------------------
}
//---------------------------------------------------
}