///<reference path="../d/jquery.d.ts" />

namespace gf_post_image_view {
//------------------------------------------------
export function init(p_image_post_element, p_log_fun) {
	p_log_fun('FUN_ENTER','gf_post_image_view.init()');

	$(p_image_post_element).find('img').on('click',(p_event)=>{

		const img_medium_url_str :string = $(p_event.target).attr('src');
		const img_large_url_str  :string = img_medium_url_str.replace('medium','large');

		view_image(img_large_url_str, p_log_fun);
	});
}
//------------------------------------------------
function view_image(p_img_url_str :string, p_log_fun) {
	p_log_fun('FUN_ENTER','gf_post_image_view.view_image()');

	const image_view_element = $(`
		<div id='image_view'>
			<div id='background'></div>
			<img></img>
			<div id="close_button">&#10006;</div>
		</div>`);

	//--------------------------------------------------------
    function load_image() {
        p_log_fun('FUN_ENTER','gf_post_image_view.view_image().load_image()');

			const image :HTMLImageElement = document.createElement('img');
            image.src                     = p_img_url_str;

            $(image).on('load',(p_e)=>{
            	console.log('img-------');
                
                const image_x_int :number = ($(window).innerWidth()-$(image).width())/2;
		    	const image_y_int :number = ($(window).innerHeight()-$(image).height())/2;

		    	$(image).css('left',image_x_int+'px');
				$(image).css('top' ,image_y_int+'px');

				$(image_view_element).append(image);

				const close_btn = $(image_view_element).find('#close_button');
				$(close_btn).css('left',(image_x_int+$(image).width())+'px');
				$(close_btn).css('top' ,image_y_int+'px');
            });
    }
    //--------------------------------------------------------

    //offset the top of the image_viewer in case the user scrolled
    $(image_view_element).css('top',document.body.scrollTop+'px');
    
    $('body').append(image_view_element);

	//prevent scrolling while in image_view
	$('body').css('overflow','hidden');

	//'ESCAPE' key
	$(document).on('keyup',(p_event)=>{
		if (p_event.which == 27) {
			
			$(image_view_element).remove();
			$('body').css('overflow','auto');
		}
	});
	$(image_view_element).find('#close_button').on('click',(p_event)=>{
		$(image_view_element).remove();
		$('body').css('overflow','auto');
	});
}
//--------------------------------------------------------
}