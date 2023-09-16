/*
GloFlow application and media management/publishing platform
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

//-----------------------------------------------------
export function draw(p_name_str :string,
	p_width_int                   :number,
	p_height_int                  :number,
	p_images_indicator_height_int :number,
	p_color_int                   :number,
	p_log_fun) :PIXI.Container {
	//p_log_fun('FUN_ENTER','gf_domain.draw()');

	const color_black   = 0x000000;
	const container_spr = new PIXI.Container();
	//---------------
	const graphics = new PIXI.Graphics();
	container_spr.addChild(graphics);

	const random_background_color = p_color_int+Math.floor(Math.random()*200);
	graphics.beginFill(random_background_color);
	//graphics.lineStyle(1,p_color_int);

	//single_page_height_px-1 - so that a little space is shown between pages
	graphics.drawRect(0,0, //x/y 
		p_width_int,   //p_width_px 
		p_height_int); //p_height_px

	//---------------
	//DRAW IMAGE INDICATOR
	if (p_width_int > 5) {
		graphics.beginFill(0xFFFFFF,0.2);
		graphics.drawRect(0,0, //x/y 
			p_images_indicator_height_int, //p_width_px 
			p_height_int);                 //p_height_px
	}
	//---------------
	//TEXT

	var posts_count_txt :PIXI.Text;
	if (p_width_int > 120) {

		/*final TextField posts_count_txt = new TextField();
		posts_count_txt.defaultTextFormat = new TextFormat('Arial', 18, Color.Black);
		posts_count_txt.text     = p_name_str;
		posts_count_txt.x        = 10;
		posts_count_txt.y        = 20;
		posts_count_txt.width    = 120;
		//posts_count_txt.height   = 10;
		posts_count_txt.wordWrap = true;
		container_spr.addChild(posts_count_txt);*/

		posts_count_txt = new PIXI.Text(p_name_str,{ //text
				font:  '18px Arial',
				fill:  color_black,
				align: 'center'
			});
		posts_count_txt.x = 10;
		posts_count_txt.y = 20;
		//posts_count_txt.width = 120;
		container_spr.addChild(posts_count_txt);
	}
	else if (p_width_int > 100) {
		/*final TextField posts_count_txt = new TextField();
		posts_count_txt.defaultTextFormat = new TextFormat('Arial', 12, Color.Black);
		posts_count_txt.text     = p_name_str;
		posts_count_txt.x        = 4;
		posts_count_txt.y        = 2;
		posts_count_txt.width    = 80;
		//posts_count_txt.height   = 10;
		posts_count_txt.wordWrap = true;
		
		container_spr.addChild(posts_count_txt);*/

		posts_count_txt = new PIXI.Text(p_name_str,{
				font:  '12px Arial',
				fill:  color_black,
				align: 'center'
			});

		posts_count_txt.x = 4;
		posts_count_txt.y = 2;
		//posts_count_txt.width = 80;
		container_spr.addChild(posts_count_txt);
	}
	else if (p_width_int > 70) {
		/*final TextField posts_count_txt = new TextField();

		posts_count_txt.defaultTextFormat = new TextFormat('Arial', 10, Color.Black);
		posts_count_txt.text     = p_name_str;
		posts_count_txt.x        = 2;
		posts_count_txt.y        = 2;
		posts_count_txt.width    = 80;
		posts_count_txt.height   = 10;
		//posts_count_txt.wordWrap = true;
		
		container_spr.addChild(posts_count_txt);*/

		posts_count_txt = new PIXI.Text(p_name_str,{
				font:  '10px Arial',
				fill:  color_black,
				align: 'center'
			});
		posts_count_txt.x = 2;
		posts_count_txt.y = 2;
		//posts_count_txt.width = 80;
		//posts_count_txt.height = 10;
		container_spr.addChild(posts_count_txt);
	}

	return container_spr;
	
	/*//DRAW RANDOM PARTICLES
	if (p_width_int > 50) {
		final int sub_squares_num_int = random_gen.nextInt(10);

		for (var i=0;i<sub_squares_num_int;i++) {
			shape.graphics.moveTo(0,0);
		    shape.graphics.beginPath();

		    int random_width_int = random_gen.nextInt(20);

		    //single_page_height_px-1 - so that a little space is shown between pages
			shape.graphics.rect(random_gen.nextInt(p_width_int)-30,  //x
								random_gen.nextInt(p_height_int)-30, //y
								random_width_int,           //p_width_px 
								random_width_int); //p_height_px


			shape.graphics.closePath();
			//shape.graphics.strokeColor(p_color_int,1);
			shape.graphics.fillColor(p_color_int+i*random_gen.nextInt(100));
		}
	}*/
	//---------------
}