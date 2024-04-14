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

///<reference path="../../../d/jquery.d.ts" />

//------------------------------------------------------------
function init_range_bar_background_canvas(p_background_element,
	p_pages_cache_map,
	p_seek_start_page_int,
	p_seek_end_page_int,
	p_log_fun,
	p_width_px              = 60,
	p_height_px             = 600,
	p_cached_page_color_hex = "0xFF8ED6FF") {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access_canvas.init_range_bar_background_canvas()');
	
	// the dimensions of the canvas must be set this way, not via CSS
	// (CSS changes would only scale the canvas)
	p_background_element.attributes['width']  = '${p_width_px}';
	p_background_element.attributes['height'] = '${p_height_px}';
	
	//------------
	// CSS
	
	p_background_element.style.position = 'relative';
	p_background_element.style.top      = '0px';

	//------------
	
	Stage      stage      = new Stage(p_background_element);
    RenderLoop renderLoop = new RenderLoop();
    
    renderLoop.addStage(stage);
  
    //------------------------------------------------------------
    function draw_cached_pages_fun() {
    	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access_canvas.init_range_bar_background_canvas().draw_cached_pages_fun()');

    	const pages_number_int      = (p_seek_start_page_int - p_seek_end_page_int).abs();
    	const single_page_height_px = (p_height_px / pages_number_int).toInt();
    	
    	final Sprite container = new Sprite();
    	final Shape shape      = new Shape();
    		
    	shape.graphics.moveTo(0, 0);
    	shape.graphics.beginPath();
    	
    	p_pages_cache_map.keys.forEach((p_page_index_str){
    			
    		const page_index_int = int.parse(p_page_index_str);
    		const page_y         = page_index_int*single_page_height_px;

    		// single_page_height_px-1 - so that a little space is shown between pages
    		shape.graphics.rect(0, page_y, // x/y 
				p_width_px,               // p_width_px 
				single_page_height_px-1); // p_height_px
    	});
    	
    	shape.graphics.closePath();
    	shape.graphics.strokeColor(0xFFFCCFFF,1);
    	shape.graphics.fillColor(p_cached_page_color_hex);
    	
    	container.addChild(shape);
    	stage.addChild(container);
    }

    //------------------------------------------------------------

    draw_cached_pages_fun();

    return draw_cached_pages_fun;
}
//------------------------------------------------------------
function init_button_seek_info_background_canvas(p_background_element,
	p_log_fun,
	p_height_px = 50,
	p_width_px  = 50,
	p_color_hex = "0xFFFFE4C4") {
	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access_canvas.init_button_seek_info_background_canvas()');

	// the dimensions of the canvas must be set this way, not via CSS
	// (CSS changes would only scale the canvas)
	p_background_element.attributes['width']  = `${p_width_px}`;
	p_background_element.attributes['height'] = `${p_height_px}`;
	
	//------------
	// CSS
	$(p_background_element).style("position", 'relative');
	$(p_background_element).style("top",      '0px');

	//------------
	   
    Stage      stage      = new Stage(p_background_element);
  
    renderLoop.addStage(stage);
  
    //------------------------------------------------------------
    function init_button() {
    	p_log_fun('FUN_ENTER', 'gf_vis_group_random_access_canvas.init_button_seek_info_background_canvas().init_button()');
    	
    	Sprite container = new Sprite();
    	Shape shape      = new Shape();
    	
    	shape.graphics.moveTo(0, 0);
    	
    	shape.graphics.beginPath();
    	shape.graphics.rect(0, 0, // x/y 
			30,  // p_width_px, 
			30); // p_height_px);
    	
    	// shape.graphics.lineTo(80,1); // x/y
    	shape.graphics.closePath();
    	
    	// shape.graphics.strokeColor(p_button_conn_color_hex,0);
    	shape.graphics.fillColor(p_color_hex);
    	
    	container.addChild(shape);
    	
    	$(container).onMouseClick.listen((MouseEvent p_event) {
    		print('click--------------------------------------------');
    	});
    	
    	return container;
    }

    //------------------------------------------------------------

    Sprite button_spr = init_button();
    button_spr.x      = 0;
    button_spr.y      = 0;
    stage.addChild(button_spr);
}