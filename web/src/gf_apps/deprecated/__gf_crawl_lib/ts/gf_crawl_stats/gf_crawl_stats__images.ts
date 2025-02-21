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

///<reference path="../../../../d/jquery.timeago.d.ts" />

//---------------------------------------------------
export function view__crawled_images_domains(p_domains_lst, p_log_fun) {
    p_log_fun('FUN_ENTER','gf_crawl_stats__images.view__crawled_images_domains()');

    const images_domains_e = $(`
        <div id="images_domains">
            <div class="title">crawled_images_domains</div>
        </div>`);
    
    //---------------------------------------------------
    function view_urls(p_urls_lst,
        p_valid_for_usage_lst,
        p_downloaded_lst,
        p_s3_stored_lst,
        p_creation_unix_times_lst,
        p_domain_e) {

        for (var i=0; i < p_urls_lst.length; i++) {
            const u_str                = p_urls_lst[i];
            const valid_for_usage_bool = p_valid_for_usage_lst[i];
            const downloaded_bool      = p_downloaded_lst[i];
            const s3_stored_bool       = p_s3_stored_lst[i];
            
            //------------------
            //CREATION_TIME
            const creation_unix_time_f = p_creation_unix_times_lst[i];
            const creation_time_f      = parseFloat(creation_unix_time_f);
            const creation_date        = new Date(creation_time_f*1000);
            const date_msg_str         = $.timeago(creation_date);
            //------------------

            const u = `
                <div class="url">
                    <div class="url_a"><a target="_blank" href="`+u_str+`">`+u_str+`</a></div>
                    <div class="creation_time">`+date_msg_str+`</div>
                    <div class="status">
                        <span class="status_item valid_for_usage">
                            <span style="font-weight:bold;">valid:</span>
                            <span class="value">`+valid_for_usage_bool+`</span>
                        </span>
                        <span class="status_item downloaded">
                            <span style="font-weight:bold;">downloaded:</span>
                            <span class="value">`+downloaded_bool+`</span>
                        </span>
                        <span class="status_item s3_stored">
                            <span style="font-weight:bold;">s3:</span>
                            <span class="value">`+s3_stored_bool+`</span>
                        </span>
                    </div>
                </div>`;
            
            $(p_domain_e).append(u);
        }
    }
    //---------------------------------------------------

    for (var domain_map of p_domains_lst) {

        const domain_str              = domain_map['domain_str'];
        const imgs_count_int          = domain_map['imgs_count_int'];
        const creation_unix_times_lst = domain_map['creation_unix_times_lst'];
        const origin_urls_lst         = domain_map['origin_urls_lst'];
        const urls_lst                = domain_map['urls_lst'];
        const valid_for_usage_lst     = domain_map['valid_for_usage_lst'];
        const downloaded_lst          = domain_map['downloaded_lst'];
        const s3_stored_lst           = domain_map['s3_stored_lst'];

        const domain_e = $(`<div class="images_domain">
                <div class='domain_name'><a href="http://`+domain_str+`" target="_blank">`+domain_str+`</a></div>
                <div class='imgs_count'><span class="imgs_count_label">imgs_count - </span><span class="imgs_count">`+imgs_count_int+`</span></div>
            </div>`);
        
        $(images_domains_e).append(domain_e);


        //IMPORTANT!! - readibility by the user, and all the DOM appends performance, 
        //              require to only display the first X urls, and a "more_btn".
        //              user can then click "more_btn" to get the remaining urls.

        if (urls_lst.length>100) {
            const urls_to_view_lst = urls_lst.slice(0,100); //from 0 to 100th element
            const urls_rest_lst    = urls_lst.slice(100);   //from 100th element to the end

            view_urls(urls_to_view_lst,
                valid_for_usage_lst,
                downloaded_lst,
                s3_stored_lst,
                creation_unix_times_lst,
                domain_e);


            //MORE_BTN
            $(domain_e).append(`<div class='more_btn'>more</div>`);
            $(domain_e).find('.more_btn').on('click',()=>{

                //IMPORTANT!! - when more_btn is clicked add all the other urls
                //              that have not been displayed yet.
                view_urls(urls_rest_lst,
                    valid_for_usage_lst,
                    downloaded_lst,
                    s3_stored_lst,
                    creation_unix_times_lst,
                    domain_e);

            });
        } 
        else {
            view_urls(urls_lst,
                valid_for_usage_lst,
                downloaded_lst,
                s3_stored_lst,
                creation_unix_times_lst,
                domain_e);
        }   
    }

    const falses_lst = $('span.value:contains("false")');
    $(falses_lst).css('background-color','red');
    $(falses_lst).css('color'           ,'white');
    $(falses_lst).css('font-weight'     ,'bold');
    $(falses_lst).css('font-size'       ,'12px');

    return images_domains_e;
}
//---------------------------------------------------
export function view__gifs_per_day_stats(p_stats_map, p_log_fun) {
    p_log_fun("FUN_ENTER","gf_crawl_stats__images.view__gifs_per_day_stats()");

    console.log('>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>');
    console.log(p_stats_map);

    const stat__crawled_gifs_lst = p_stats_map['stat__crawled_gifs_lst'];
}
//---------------------------------------------------
export function view__gif_stats(p_stats_map, p_log_fun) {
    p_log_fun("FUN_ENTER","gf_crawl_stats__images.view__gif_stats()");

    console.log('>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>');
    console.log(p_stats_map);

    const stat__crawled_gifs_lst = p_stats_map['stat__crawled_gifs_lst'];
    const gifs_stats             = $(`<div id="gifs_stats"></div>`);

    for (var e_map of stat__crawled_gifs_lst) {

        const domain_str             = e_map['domain_str'];
        const imgs_count_int         = e_map['imgs_count_int'];
        const urls_by_origin_url_lst = e_map['urls_by_origin_url_lst'];

        const gifs_domain = $(`
            <div class="gifs_domain">
                <div class="title"><a target="_blank" href=http://"`+domain_str+`">`+domain_str+`</a></div>
                <div class="imgs_count">`+imgs_count_int+`</div>
                <div class="origin_urls"></div>
            </div>`);
        $(gifs_stats).append(gifs_domain);

        for (var u_map  of urls_by_origin_url_lst) {

            const origin_url_str     = u_map['origin_url_str'];
            const creation_times_lst = u_map['creation_times_lst'];
            const urls_lst           = u_map['urls_lst'];
            const nsfv_lst           = u_map['nsfv_lst'];
            
            const gifs_from_origin_e = $(`
                <div class="gifs_from_origin_url">
                    <div class="title"><span>origin_url: </span><a target="_blank" href="`+origin_url_str+`">`+origin_url_str+`</a></div>
                    <div class="urls"></div>
                </div>`);
            $(gifs_domain).find('.origin_urls').append(gifs_from_origin_e);

            var i = 0;
            for (var url_str of urls_lst) {

                //------------------
                //CREATION_TIME
                const creation_unix_time_f = creation_times_lst[i];
                const creation_time_f      = parseFloat(creation_unix_time_f);
                const creation_date        = new Date(creation_time_f*1000);
                const date_msg_str         = $.timeago(creation_date);
                //------------------

                const nsfv_bool = nsfv_lst[i];

                $(gifs_from_origin_e).find('.urls').append($(`
                    <div class="url">
                        <div class="url_a"><a target="_blank" href="`+url_str+`">`+url_str+`</a></div>
                        <div class="creation_time">`+date_msg_str+`</div>
                        <div class="nsfv">`+nsfv_bool+`</div>
                    </div>`));

                i+=1;
            }
        }
    }

    const nsfv_lst = $(gifs_stats).find('.nsfv:contains("true")');
    $(nsfv_lst).css('background-color', 'red');
    $(nsfv_lst).css('color',            'white');
    $(nsfv_lst).css('font-weight',      'bold');
    return gifs_stats
}