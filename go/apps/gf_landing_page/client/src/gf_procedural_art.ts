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

///<reference path="./d/jquery.d.ts" />

namespace gf_procedural_art {
//--------------------------------------------------------
export function init(p_log_fun) :HTMLCanvasElement {
	p_log_fun('FUN_ENTER','gf_landing_page.init()');

	//final CanvasElement canvas = new CanvasElement();
	const canvas = <HTMLCanvasElement> $('#randomized_art #randomized_canvas')[0];
	$(canvas).addClass('randomized_canvas'); //FIX!! - needed?

	canvas.width  = $('#randomized_art').width();
	canvas.height = $('#randomized_art').height();
	//var ctx       = canvas.context2D;

	draw_randomized_squares(canvas,
					canvas.width,
					canvas.height,
					p_log_fun);
	return canvas;
}
//-------------------------------------------------
function draw_randomized_squares(p_canvas :HTMLCanvasElement,
					p_width_int  :number,
					p_height_int :number,
					p_log_fun) {
	p_log_fun('FUN_ENTER','gf_procedural_art.draw_randomized_squares()');
	
	const ctx = p_canvas.getContext('2d');
	ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
	// flip context horizontally

	const dots_num_int :number = Math.floor(Math.random()*30);

	const random_background_color_str = get_random_color();
	ctx.fillStyle = random_background_color_str;
	
	ctx.fillRect(0, //i+20, 
			0, //i+30, 
			p_width_int,
			p_height_int);
	
	for (var i=0;i < dots_num_int;i++) {
		draw_simple_square();
		draw_complex_square();
	}
	
	//draw_border();

	//-------------------------------------------------
	function draw_simple_square() {
		const x_int :number = Math.floor(Math.random()*p_width_int);
		const y_int :number = Math.floor(Math.random()*p_height_int);
			
		ctx.fillStyle = get_random_color();

		const random_square_scale :number = Math.floor(Math.random()*30);
		ctx.fillRect(x_int, //i+20, 
				y_int, //i+30, 
				random_square_scale,
				random_square_scale);
	}
	//-------------------------------------------------
	function draw_complex_square() {
		const x_int :number = Math.floor(Math.random()*p_width_int);
		const y_int :number = Math.floor(Math.random()*p_height_int);
			
		//-----------
		//RANDOM_COLOR
		const random_r_int :number = Math.floor(Math.random()*255);
		const random_g_int :number = Math.floor(Math.random()*255);
		const random_b_int :number = Math.floor(Math.random()*255);
		ctx.fillStyle = "rgba("+random_r_int+","+random_g_int+","+random_b_int+",255)";
		//-----------

		const random_square_scale :number = Math.floor(Math.random()*50);
		ctx.fillRect(x_int, //i+20, 
				y_int, //i+30, 
				random_square_scale,
				random_square_scale);

		//nested square 
		if (random_square_scale > 10) {

			ctx.fillStyle = get_random_color();
			ctx.fillRect(x_int+2, //i+20, 
					y_int+2, //i+30, 
					random_square_scale-5,
					random_square_scale-5);

			//nested square's satellites
			if (Math.random() >= 0.5) {
				ctx.fillStyle = "rgba("+(random_r_int+15)+","+(random_g_int+15)+","+(random_b_int+15)+",255)";

				draw_satelites(x_int,
							y_int,
							random_square_scale);
			}
		}	
		//-------------------------------------------------
		function draw_satelites(p_target_x_int :number,
						p_target_y_int :number,
						p_scale_int    :number) {

			ctx.fillRect(p_target_x_int+6, //i+20, 
					p_target_y_int+6, //i+30, 
					p_scale_int-3,
					p_scale_int-3);

			ctx.fillRect(p_target_x_int+8, //i+20, 
					p_target_y_int+8, //i+30, 
					p_scale_int*0.6,
					p_scale_int*0.6);

			ctx.fillRect(p_target_x_int-4, //i+20, 
					p_target_y_int-4, //i+30, 
					p_scale_int*0.5,
					p_scale_int*0.4);
		}
		//-------------------------------------------------	
	}
	//-------------------------------------------------
	function draw_border() {
		//create a border after everythings else, so that it overlaps other things
		//previously drawn on the canvas
		ctx.lineWidth   = 1;
		ctx.strokeStyle = "rgba(184,86,40,255)";
		ctx.strokeRect(0,0,
				p_width_int,
				p_height_int);
	}
	//-------------------------------------------------
	function get_random_color() :string {
		const random_r_int    :number = Math.floor(Math.random()*255);
		const random_g_int    :number = Math.floor(Math.random()*255);
		const random_b_int    :number = Math.floor(Math.random()*255);
		const random_rgba_str :string = "rgba("+random_r_int+","+random_g_int+","+random_b_int+",255)";
		return random_rgba_str;
	}
	//-------------------------------------------------
}
//--------------------------------------------------------
}