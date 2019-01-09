///<reference path="../d/pixi.js.d.ts" />

namespace gf_domains_conn {
//-----------------------------------------------------
/*activate_conn(Function p_log_fun) {
	p_log_fun('FUN_ENTER','domains_conn.activate_conn()');
}*/
//-----------------------------------------------------
export function draw_connectivity(p_domains_lst :Object[],
	p_domains_infos_spr :PIXI.Container,
	p_width_int  :number,
	p_height_int :number,
	p_color_int  :number,
	p_log_fun) :PIXI.Container {
	//{int     p_width_int :300,
	//int      p_height_int:300,
	//int      p_color_int :0xF31111FF}) :PIXI.Sprite {
	p_log_fun('FUN_ENTER','gf_domains_conn.draw_connectivity()');

	const container = new PIXI.Container();
	//------------------
	//BACKGROUND
	/*shape.graphics.moveTo(0, 0);
    shape.graphics.beginPath();
    
    //single_page_height_px-1 - so that a little space is shown between pages
	shape.graphics.rect(0,0,         //x/y 
						p_width_int,   //p_width_px 
						p_height_int); //p_height_px

    //shape.graphics.strokeColor(Color.Blue, 5);

    shape.graphics.closePath();
	shape.graphics.strokeColor(p_color_int,1);
	shape.graphics.fillColor(p_color_int);*/

	const graphics = new PIXI.Graphics();

	const random_background_color = p_color_int+Math.floor(Math.random()*200);
	graphics.beginFill(random_background_color);
	graphics.lineStyle(1,p_color_int);

	//single_page_height_px-1 - so that a little space is shown between pages
	graphics.drawRect(0,0, //x/y 
		p_width_int,   //p_width_px 
		p_height_int); //p_height_px

	container.addChild(graphics);
	//------------------


	for (var domain_info_map of p_domains_lst) {

  		const container :PIXI.Sprite = domain_info_map['container'];

    	if ('info_container' in domain_info_map) {
	    	const domain_info_container :PIXI.Container = domain_info_map['info_container'];

	    	const start_x_int :number = container.x           + container.width;
	    	const start_y_int :number = container.y           + container.height/2;
	    	const end_x_int   :number = p_domains_infos_spr.x + domain_info_container.x;
	    	const end_y_int   :number = p_domains_infos_spr.y + domain_info_container.y + domain_info_container.height/2;

	    	graphics.beginFill(gf_color.get_int('black')); //.beginPath();

	    	graphics.moveTo(start_x_int, start_y_int);
			graphics.lineTo(end_x_int,end_y_int);
			
			graphics.endFill(); //.closePath();

			//graphics.strokeColor(gf_color.get_hex('lightgray'), 1);
			graphics.lineStyle(1, //lineWidth
				gf_color.get_int('lightgrey'), //color
				1);                            //alpha
		}
    }
	return container;
}
//-----------------------------------------------------
}