

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
		'pg_index':p_page_index_int,
		'pg_size': p_page_elements_num_int
	};

    $.ajax({
        'url':        url_str,
        'type':       'GET',
        'data':       data_map,
        'contentType':'application/json',
        'success':    (p_response_str)=>{

            const response_map = JSON.parse(p_response_str);
            const status_str   = response_map['status_str'];
            const page_lst :Object[] = response_map['data'][];

            p_on_complete_fun(page_lst);
        },
        'error':(jqXHR,p_text_status_str)=>{
            p_on_error_fun(p_text_status_str);
        }
    });
}
//-----------------------------------------------------
}