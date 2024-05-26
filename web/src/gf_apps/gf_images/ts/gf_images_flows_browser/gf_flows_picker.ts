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

// ///<reference path="../../../../d/jquery.d.ts" />

import * as gf_images_events from "./../gf_images_core/gf_events";
import * as gf_user_events from "./../../../../gf_events/ts/gf_user_events";

//---------------------------------------------------
export async function init(p_events_enabled_bool :boolean,
    p_host_str :string,
    p_log_fun) {

    const all_flows_lst = await http__get_all_flows(p_host_str, p_log_fun) as {}[];


    // <div id="flows_experimental_label">experimental:</div>
    const all_flows_container = $(`
        <div id="flows_picker">
            <div id="expand_btn"></div>
            <div id="flows">
            </div>

            
            <div id="flows_experimental">
            </div>
        </div>`);
    $('body').append(all_flows_container);

    //------------------
    // allow for the flow-picker to be toggled in visibility,
    // displayed/hidden by clicking this button.

    // flow-picker not initially visible
    var visible_bool = false;
    $(all_flows_container).find("#flows").css("display", "none");
    $(all_flows_container).find("#flows_experimental").css("display", "none");

    //------------------

    $(all_flows_container).find("#expand_btn").click(()=>{


        if (visible_bool) {

            // hide
            $(all_flows_container).find("#flows").css("display", "none");
            $(all_flows_container).find("#flows_experimental").css("display", "none");
            visible_bool = false;
        }
        else {

            // show
            $(all_flows_container).find("#flows").css("display", "block");
            $(all_flows_container).find("#flows_experimental").css("display", "block");
            visible_bool = true;

            //------------------
            // EVENTS
            if (p_events_enabled_bool) {
                
                const event_meta_map = {

                };
                gf_user_events.send_event_http(gf_images_events.GF_IMAGES_FLOWS_PICKER_OPEN,
                    "browser",
                    event_meta_map,
                    p_host_str)
            }

            //------------------
        }
    });



    const experimental_flows_lst = [
        "discovered",
        "gifs"
    ];
    for (const flow_map of all_flows_lst ) {
        const flow_name_str = flow_map["flow_name_str"];

        // FIX!! - allow access to these flows only to logged in users, ton of content there
        //         but not filtered yet for NSFW.
        if (flow_name_str == "discovered" || flow_name_str == "gifs") {
            continue;
        }

        const flow_imgs_count_int :number = flow_map["flow_imgs_count_int"];
        const flow_url_str        :string = `${p_host_str}/images/flows/browser?fname=${flow_name_str}`;

        var target_container_id_str :string;
        if (experimental_flows_lst.includes(flow_name_str)) {
            target_container_id_str = "flows_experimental";
        } else {
            target_container_id_str = "flows";
        }

        $(all_flows_container).find(`#${target_container_id_str}`).append(`
            <div class="flow_info">
                <div class="flow_imgs_count">${flow_imgs_count_int}</div>
                <div class="flow_name">
                    <a href="${flow_url_str}">${flow_name_str}</a>
                </div>
            </div>
        `);
    }
}

//---------------------------------------------------
async function http__get_all_flows(p_host_str :string,  p_log_fun) {
    return new Promise(function(p_resolve_fun, p_reject_fun) {

        const url_str = `${p_host_str}/v1/images/flows/all`;
        p_log_fun('INFO', `url_str - ${url_str}`);

        //-------------------------
        // HTTP AJAX
        $.get(url_str,
            function(p_data_map) {
                if (p_data_map["status"] == 'OK') {
                    const all_flows_lst = p_data_map['data']['all_flows_lst'];
                    p_resolve_fun(all_flows_lst);
                }
                else {
                    p_reject_fun(p_data_map["data"]);
                }
            });

        //-------------------------	
    });
}