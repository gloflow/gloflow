/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

//-------------------------------------------------
export function send_event_http(p_event_type_str :string,
    p_source_type_str :string,
    p_data_map        :object,
    p_host_str        :string) {

    return new Promise(function(p_resolve_fun, p_reject_fun) {
        const data_map = {
            "type_str":        p_event_type_str,
            "source_type_str": p_source_type_str,
            "data_map":        p_data_map
        };

        //-------------------
        // URL
        const url_str = `${p_host_str}/v1/a/ue`;

        //-------------------
        
        /*
        if (navigator.sendBeacon) {
            
            const data = JSON.stringify(data_map);
            const blob = new Blob([data], { type: 'application/json' });

            navigator.sendBeacon(url_str, blob);

            const resp_data_map = {};
            p_resolve_fun(resp_data_map);
        }
        */

        
        $.ajax({
            'url':         url_str,
            'type':        'POST',
            'data':        JSON.stringify(data_map),
            'contentType': 'application/json',

            /*
            The XMLHttpRequest.withCredentials property is a Boolean that indicates
            whether or not cross-site Access-Control requests should be made using
            credentials such as cookies, authorization headers or TLS client certificates.
            Setting withCredentials has no effect on same-site requests.
            if withCredentials is set to true in an XHR request,
            the Origin header will be included in the request for cross-origin requests.
            */
            xhrFields: {
                withCredentials: true
            },

            // HEADERS
            headers: {
                
            },

            'success': (p_response_map)=>{
                
                const status_str = p_response_map["status"];
                const data_map   = {}; // p_response_map["data"];

                if (status_str == "OK") {
                    p_resolve_fun(data_map);
                } else {
                    p_reject_fun(data_map);
                }
            },
            'error': (jqXHR, p_text_status_str :string)=>{
                p_reject_fun(p_text_status_str);
            }
        });
    });
}