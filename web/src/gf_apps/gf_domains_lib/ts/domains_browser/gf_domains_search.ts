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

namespace gf_domains_search {
//-----------------------------------------------------
export function init_domain_search(p_domains_infos_lst :Object[],
	p_onPick_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER','gf_domains_search.init_domain_search()');

	const suggestions_lst :Object[] = [];
	for (var domain_info_map of p_domains_infos_lst) {
		suggestions_lst.push({
			'value': domain_info_map['url_str'],
			'data':  JSON.stringify(domain_info_map)
		});
	}

	console.log('>>>>>>>>>>>>>>>>>>>>>>>>')
	console.log(suggestions_lst)
	//----------------
	//JS - QUERY INPUT FIELD 
    //https://www.devbridge.com/sourcery/components/jquery-autocomplete/
    
    const config_map = {
            //'lookup':[
			//	{'value':'test' ,'data':'AE'},
			//	{'value':'test2','data':'AE2'},
			//	{'value':'test4','data':'AE3'}
			//],
			'lookup':   suggestions_lst,
			'onSelect': (p_suggestion)=>{
				
				const domain_url_str  :string = p_suggestion['value'];
				const domain_info_map :Object = JSON.stringify(p_suggestion['data']);
				p_log_fun('INFO','domain_info_map - '+domain_info_map);

				p_onPick_fun(domain_info_map);
		    }
        };
	$('#domain_search #query_input').autocomplete(<any>config_map);
    //----------------
}
//-----------------------------------------------------
}