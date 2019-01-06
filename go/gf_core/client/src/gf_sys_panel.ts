

namespace gf_sys_panel {
//-----------------------------------------------------
export function init(p_log_fun) {
	p_log_fun('FUN_ENTER','gf_sys_panel.init()');

	const sys_panel_element = $(
		`<div id="sys_panel">
			<div id="view_handle"></div>
			<div id="home_btn">
				'<img src="/images/d/gf_header_logo.png"></img>
			</div>
			<div id="images_app_btn"><a href="/images/flows/browser">Images</a></div>
			<div id="publisher_app_btn"><a href="/posts/browser">Posts</a></div>
			<div id="get_invited_btn">get invited</div>
			<div id="login_btn">login</div>
		</div>`);

	$('body').append(sys_panel_element);

	$(sys_panel_element).find('#view_handle').on('mouseover',(p_e)=>{

		$(sys_panel_element).animate({
			top:0 //move it
		},
		200,
		()=>{
			
			$(sys_panel_element).find('#view_handle').css('visibility','hidden');
		});
	});
}
//-----------------------------------------------------
}