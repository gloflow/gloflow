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

///<reference path="../d/jquery.d.ts" />

namespace gf_post_tag_mini_view {
//-----------------------------------------------------
function init_tags_mini_view(p_tags_lst :string[],
	p_domain_str :string,
	p_log_fun) {
	p_log_fun('FUN_ENTER','gf_post.gf_post_main().init_tags_mini_view()');
  
	//Element tags_container_element = $('#post_tags');
	//List    tags_html_elements_lst = $('.post_tag');
	const post_title_element = $('.post_title_str');
  
	const tags_mini_view_element = tags_mini_view.init(p_tags_lst,
		'post',
		p_domain_str,
								
		//onClose_fun
		()=>{
			$('#post_info_container').css('opacity','1.0');
			$('.post_element').each((p_i,p_element)=>{$(p_element).css('opacity','1.0');});
		},
		p_log_fun);
			
	//#post_info_container is what the post_title is attached to.
	//this is important because the tags_mini_view is being positioned relative to 
	//post title, so its important that they share the same DOM/scene-graph parent
	$('#post_info_container').append(tags_mini_view_element);
	
	//#tags_container - is the part of the tags_mini_view that is always visible, 
	//                  the part with the '#' symbol
	const tags_mini_view_tags_container_element = (tags_mini_view_element).find('#tags_container');
	
	//-----------------
	//Y position
	
	gf_vis_lib.center_vertical_element(tags_mini_view_tags_container_element, post_title_element, p_log_fun);
	//-----------------
	//X position
	
	const space_between_post_title_and_tags_mini_view :number = 10;
	const tags_mini_view_new_x                        :number = $(post_title_element).offset.left + $(post_title_element).offset.width + space_between_post_title_and_tags_mini_view;
		
	$(tags_mini_view_tags_container_element).css('left',tags_mini_view_new_x.toString()+'px');
	//-----------------
	
	//when the cursor is over the tags_mini_view lighten all of the post_element's
	$(tags_mini_view_tags_container_element).on('onmouseover',(event)=>{
		
		////--------------
		////ANALYTICS
		//gf_analytics.a_track_click('objects_with_tags_mini_view', //p_category_str
		//						   'mini_view_hover',             //p_label_str
		//						   p_log_fun);
		////---------------------
		
		//---------------------
		//ANIMATION
		//js.scoped(() {
		//	js.context.jQuery('#post_info_container').animate({
		//		'opacity': 0.3
		//	},200);
		//	
		//	js.context.jQuery('.post_element').animate({
		//		'opacity': 0.3
		//	},200);
		//});
		//---------------------
	});
	
	$(tags_mini_view_tags_container_element).on('onmouseleave',(p_event)=>{
		$('#post_info_container').css('opacity','1.0');
		$('.post_element').each((p_i,p_element)=>{$(p_element).css('opacity','1.0');});
	});

	//----------------
	const post_tags_element = $('.post_tags');
	
	//position .post_tags div vertically centered with .post_title_str
	const new_mini_view_y :number = $(post_title_element).offset().top + ($(post_title_element).offset().height - $(post_tags_element).offset().height)/2;
	
	$(post_tags_element).css('top',new_mini_view_y.toString());
	//----------------
}
//-----------------------------------------------------
}