///<reference path="../d/jquery.d.ts" />

namespace gf_crawl_dashboard {

declare var EventSource;

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

	gf_crawl_dashboard.init(log_fun);
	
});
//-------------------------------------------------
export function init(p_log_fun) {
	gf_crawl_events.init_SSE(p_log_fun);

	//---------------------
	//IMAGES
	$('#get_recent_images_btn').on('click',()=>{

		gf_crawl_images_browser.init__recent_images(p_log_fun);
	});
	//---------------------
}
//---------------------------------------------------
}