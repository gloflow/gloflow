

namespace gf_gifs {
//---------------------------------------------------
export function http__gif_get_info(p_gf_img_id_str,
					p_host_str,
					p_on_complete_fun,
					p_on_error_fun,
					p_log_fun) {
	p_log_fun('FUN_ENTER','gf_gifs.http__gif_get_info()');

	const url_str = 'http://'+p_host_str+'/images/gif/get_info?gfimg_id='+p_gf_img_id_str;
	p_log_fun('INFO','url_str - '+url_str);

	//-------------------------
	//HTTP AJAX
	$.get(url_str,
		(p_data_map) => {
			console.log('response received');
			//const data_map = JSON.parse(p_data);

			if (p_data_map["status_str"] == 'OK') {
				const gif_map = p_data_map['data']['gif_map'];
				p_on_complete_fun(gif_map);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
		});
	//-------------------------	
}
//---------------------------------------------------
}