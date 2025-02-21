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

///<reference path="../../../../d/jquery.d.ts" />

declare var sigma;
declare var c3;
//---------------------------------------------------
export function view__new_links_per_day(p_stats_lst, p_parent, p_log_fun) {
    p_log_fun('FUN_ENTER','gf_crawl_stats__links.view__new_links_per_day()');

    const plot = $(`
        <div id='plots'>
            <div id='new_links_per_day__plot'>
                <svg width='800' height='600'></svg>
            </div>
            <div id='total_links_count_per_day__plot'>
                <svg width='800' height='600'></svg>
            </div>
        </div>`);

    $(p_parent).append(plot);
    //-----------------
    const daily_total_count_lst           = [];
    const daily_valid_for_crawl_total_lst = [];
    const daily_fetched_total_lst         = [];
    const total_counts_lst                = [];
    
    var i               = 0;
    var total_count_int = 0;

    for (var day_stat_map of p_stats_lst) {

        console.log('>>>>> =====================================')
        console.log(day_stat_map);

        const count_int                 = day_stat_map['total_count_int'];
        const valid_for_crawl_total_int = day_stat_map['valid_for_crawl_total_int'];
        const fetched_total_int         = day_stat_map['fetched_total_int'];

        console.log(count_int);
        console.log(valid_for_crawl_total_int)
        console.log(fetched_total_int)

        daily_total_count_lst.push(count_int);
        daily_valid_for_crawl_total_lst.push(valid_for_crawl_total_int);
        daily_fetched_total_lst.push(fetched_total_int);

        total_counts_lst.push(total_count_int); //total_count_int);

        i               += 1;
        total_count_int += count_int;
    }
    //-----------------

    daily_total_count_lst.unshift('per-day new links count'); //add column title as first element
    daily_valid_for_crawl_total_lst.unshift('per-day valid-for-crawl links count');
    daily_fetched_total_lst.unshift('per-day fetched links count');

    const chart = c3.generate({
        bindto: '#new_links_per_day__plot',
        data: {
          columns: [
            daily_total_count_lst,
            daily_valid_for_crawl_total_lst,
            daily_fetched_total_lst
            //['data1', 30, 200, 100, 400, 150, 250],
            //['data2', 50, 20, 10, 40, 15, 25]
          ]
        }
    });

    total_counts_lst.unshift('per-day total links count')
    const chart2 = c3.generate({
        bindto: '#total_links_count_per_day__plot',
        data: {
          columns: [
            total_counts_lst
          ]
        }
    });
}
//---------------------------------------------------
export function view__unresolved(p_stats_lst, p_log_fun) {
	p_log_fun('FUN_ENTER','gf_crawl_stats__links.view__unresolved()');

    const container = $(`
        <div id="unresolved_links">
            <div class="title">unresolved_links</div>
        </div>`);

    for (var d_map of p_stats_lst) {

        const origin_domain_str = d_map['origin_domain_str'];

        const domain = $(`
            <div class="origin_domain">
                <div class="title">`+origin_domain_str+`</div>
                <div class="plot_urls_references_graph_btn">graph</div>
                <div class="origin_urls"></div>
            </div>`);
        $(container).append(domain);

        const origin_urls_lst               = d_map['origin_urls_lst'];
        const a_hrefs__from_origin_urls_lst = d_map['a_hrefs__from_origin_urls_lst'];

        $(domain).find('.plot_urls_references_graph_btn').on('click',(p_e)=>{

            const graph = view__links_graph(origin_domain_str,
                origin_urls_lst,
                a_hrefs__from_origin_urls_lst,
                domain,
                p_log_fun);
        });

        for (var i=0;i<origin_urls_lst.length;i++) {

            const origin_url_str              = origin_urls_lst[i];
            const a_hrefs_from_origin_url_lst = a_hrefs__from_origin_urls_lst[i];

            $(domain).append($(`
                <div>
                    <div class="origin_url">`+origin_url_str+`</div>
                    <div class="a_hrefs_count">`+a_hrefs_from_origin_url_lst.length+`</div>
                </div>`));
        }
    }

    return container;
}
//---------------------------------------------------
export function view__links_graph(p_origin_domain_str :string,
    p_origin_urls_lst               :string[],
    p_a_hrefs__from_origin_urls_lst :string[],
    p_parent,
    p_log_fun) {
    p_log_fun('FUN_ENTER','gf_crawl_stats__links.view__links_graph()');

    const c = $(`
        <div id="graph_container">
            <div id="close_btn">x</div>
        </div>`);

    $(p_parent).append(c);
        
    const edges_lst = [];
    const nodes_lst = [];
    const nodes_map = {}; //used as a SET, to eliminate duplicates (sigmajs doesnt like duplicate nodes in the nodes_lst)

    //----------------------
    nodes_lst.push({
        id:    p_origin_domain_str,
        label: p_origin_domain_str,
        x:     0,   //Math.random(),
        y:     0,   //Math.random(),
        size:  0.5, //Math.random(),
        color: '#111'
    });
    nodes_map[p_origin_domain_str] = true;
    //----------------------

    var i=0;
    for (var origin_url_str of p_origin_urls_lst) {

        if (!(origin_url_str in nodes_map)) {
            nodes_lst.push({
                id:    origin_url_str,
                label: origin_url_str,
                x:     Math.random(),
                y:     Math.random(),
                size:  0.5, //Math.random(),
                color: '#111'
            });

            nodes_map[origin_url_str] = true;
        }
    
        const a_hrefs_from_origin_url_lst = p_a_hrefs__from_origin_urls_lst[i];

        edges_lst.push({
                id:     origin_url_str+' '+i,
                source: p_origin_domain_str,
                target: origin_url_str,
                size:   Math.random(),
                color:  '#ccc'
            }); 

        for (var j=0; j < a_hrefs_from_origin_url_lst.length; j++) {

            const a_href_from_origin_url_str = a_hrefs_from_origin_url_lst[j];
                
            if (!(a_href_from_origin_url_str in nodes_map)) {
                nodes_lst.push({
                    id:    a_href_from_origin_url_str,
                    label: 'n', //a_href_from_origin_url_str,
                    x:     Math.random(),
                    y:     Math.random(),
                    size:  0.1, //Math.random(),
                    color: '#666'
                });

                nodes_map[a_href_from_origin_url_str] = true;
            }   

            edges_lst.push({
                id:     origin_url_str + '_'+a_href_from_origin_url_str+'_'+j,
                source: origin_url_str,
                target: a_href_from_origin_url_str,
                size:   0.5, //Math.random(),
                color:  '#ccc'
            });       
        }

        i+=1;
    }

    const s = new sigma({
        graph: {
            nodes: nodes_lst,
            edges: edges_lst
        },
        container: 'graph_container'
    });

    s.startForceAtlas2({worker: true, barnesHutOptimize: false});

    //-----------------
    //CLOSE_BTN
    $(c).find('#close_btn').on('click', (p_e)=>{
        s.stopForceAtlas2({worker: true, barnesHutOptimize: false});
        $(c).remove();

    });
    //-----------------
}
//---------------------------------------------------
export function view__crawled_domains(p_domains_lst, p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_crawl_stats__links.view__crawled_domains()');

    const links_domains_e = $(`
        <div id="links_domains">
            <div class="title">crawled_links_domains</div>
        </div>`);
   
    for (var domain_map of p_domains_lst) {

        const domain_str              = domain_map['domain_str'];
        const links_count_int         = domain_map['links_count_int'];
        const creation_unix_times_lst = domain_map['creation_unix_times_lst'];
        const a_href_lst              = domain_map['a_href_lst'];
        const origin_urls_lst         = domain_map['origin_urls_lst'];
        const valid_for_crawl_lst     = domain_map['valid_for_crawl_lst'];
        const fetched_lst             = domain_map['fetched_lst'];
        const images_processed_lst    = domain_map['images_processed_lst'];

        const domain_e = $(`<div class="links_domain">
                <div class='domain_name'><a href="http://`+domain_str+`" target="_blank">`+domain_str+`</a></div>
                <div class='links_count'><span class="links_count_label">links_count - </span><span class="links_count">`+links_count_int+`</span></div>
            </div>`);
        $(links_domains_e).append(domain_e);

        for (var i=0;i<a_href_lst.length;i++) {

            const a_href_str            = a_href_lst[i];
            const valid_for_crawl_bool  = valid_for_crawl_lst[i];
            const fetched_bool          = fetched_lst[i];
            const images_processed_bool = images_processed_lst[i];

            const u = `
                <div class="a_href">
                    <a target="_blank" href="`+a_href_str+`">`+a_href_str+`</a>
                    <div class="status">
                        <span class="status_item valid_for_crawl">
                            <span style="font-weight:bold;">valid_for_crawl:</span>
                            <span class="value">`+valid_for_crawl_bool+`</span>
                        </span>
                        <span class="status_item fetched">
                            <span style="font-weight:bold;">fetched:</span>
                            <span class="value">`+fetched_bool+`</span>
                        </span>
                        <span class="status_item images_processed">
                            <span style="font-weight:bold;">images_processed:</span>
                            <span class="value">`+images_processed_bool+`</span>
                        </span>
                    </div>
                </div>`;
            $(domain_e).append(u);
        }
    }

    return links_domains_e;
}