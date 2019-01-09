///<reference path="../d/pixi.js.d.ts" />

namespace gf_domains_infos {
//-----------------------------------------------------
export function draw(p_domains_lst :Object[],
	p_width_int  :number,
	p_height_int :number,
	p_ctx_map    :Object,
	p_onPick_fun,
	p_log_fun) :PIXI.Container {
	//{int      p_width_int :200,
	//int      p_height_int:600}) :PIXI.Sprite {
	p_log_fun('FUN_ENTER','gf_domains_infos.draw()');

	const container               = new PIXI.Container();
	const domain_info_padding_int = 5;
	//--------------
	//TOP_DOMAINS
	
	var top_domains_lst :Object[];
	if (p_domains_lst.length > 30) top_domains_lst = p_domains_lst.slice(0,30);
	else                           top_domains_lst = p_domains_lst;
	//--------------
	//100:height=domain_percent_of_all:x
	//x=(height*domain_percent_of_all)/100;
	
	const domains_count_int               :number = top_domains_lst.length;
	const domain_percent_of_all           :number = 100/domains_count_int;
	const new_scaled_domain_info_height_f :number = ((p_height_int*domain_percent_of_all)/100);

	p_log_fun('INFO','>>>>>>>>>>>>>>>>>>>');
	p_log_fun('INFO','domains_count_int               - '+domains_count_int);
	p_log_fun('INFO','domain_percent_of_all           - '+domain_percent_of_all);
	p_log_fun('INFO','new_scaled_domain_info_height_f - '+new_scaled_domain_info_height_f);
	//--------------
	const domains_stats_lst :PIXI.Container[] = [];

	var last :PIXI.Container;
    for (var domain_info_map of top_domains_lst) {

    	const domain_info_height_int :number = new_scaled_domain_info_height_f - 2*domain_info_padding_int;
		const domain_info :PIXI.Container = draw_domain_info(domain_info_map,
			20, //p_width_int
			domain_info_height_int, //p_height_int
			gf_color.get_int('lightgrey'),
			p_ctx_map,
			p_onPick_fun,
			p_log_fun);

    	//domain_info.height = new_scaled_domain_info_height_f;

    	//--------------
    	//LAYOUT
    	if (last != null) {
    		domain_info.y = domain_info_padding_int + last.y + last.height + domain_info_padding_int;
    	}
    	//--------------
    	container.addChild(domain_info);

    	domain_info_map['info_container'] = domain_info;
    	last                              = domain_info;

    	domains_stats_lst.push(domain_info);
    }

    //--------------
	//SELECTOR_BAR
	const selector_bar_spr :PIXI.Container = init_selector_bar(domains_stats_lst,
		top_domains_lst,
		60,              //p_width_int
		p_height_int,    //p_height_int
		gf_color.get_int('lightgrey'), //p_color_int
		p_onPick_fun,
		p_log_fun);
	selector_bar_spr.x = 10;
	container.addChild(selector_bar_spr);
	//--------------
    return container;
}
//-----------------------------------------------------
function init_selector_bar(p_domains_stats_lst :PIXI.Container[],
	p_domains_lst :Object[],
	p_width_int   :number,
	p_height_int  :number,
	p_color_int   :number,
	p_onPick_fun,
	p_log_fun) :PIXI.Container {
	//{int      p_width_int :30,
	//int       p_height_int:100
	p_log_fun('FUN_ENTER','gf_domains_infos.init_selector_bar()');

	const container     = new PIXI.Container();
	const background_gr = new PIXI.Graphics();
	//---------------
	draw_background(background_gr,
		p_width_int,
		p_height_int,
		p_color_int,
		p_log_fun);
	container.addChild(background_gr);
	//---------------
	//SCROOL_INDICATOR
	const scroll_indicator_gr = new PIXI.Graphics();
	draw_background(scroll_indicator_gr,
		p_width_int,
		1, //p_height_int
		gf_color.get_int('gray'),
		p_log_fun);
		//p_width_int :p_width_int,
		//p_height_int:1);

	container.addChild(scroll_indicator_gr);
	//---------------

	//-----------------------------------------------------
	function get_selected_domain_stats() :Object {
		p_log_fun('FUN_ENTER','gf_domains_infos.init_selector_bar().get_selected_domain_stats()');

		for (var i=0;i<p_domains_stats_lst.length;i++) {
			var ds :PIXI.Container = p_domains_stats_lst[i];

			//--------------------
			//IMPORTANT!! - hitTestObject() - hit detection
			//                                if the horizontal selector bar crosses
			//                                the domain_stats sprite (its bounding box)
			//                                this will return True (bool)
			//if (domain_stats.hitTestObject(scroll_indicator_gr)) {

			if (hitTestObject(ds.x,
				ds.y,
				ds.width,
				ds.height,
				scroll_indicator_gr.x,
				scroll_indicator_gr.y,
				scroll_indicator_gr.width,
				scroll_indicator_gr.height)) {
				return {
					'i':           i,
					'domain_stats':ds
				};
			}
			//--------------------
		}
		//-----------------------------------------------------
		function hitTestObject(x1, y1, w1, h1, x2, y2, w2, h2) {
			if (x1 + w1 > x2) if (x1 < x2 + w2) if (y1 + h1 > y2) if (y1 < y2 + h2) return true;
			return false;
		}
		//-----------------------------------------------------
		/*function hitTestObject(r1, r2) {

			return !(r2.x > (r1.x + r1.width) || 
				(r2.x + r2.width) < r1.x  || 
				r2.y > (r1.y + r1.height) ||
				(r2.y + r2.height) < r1.y);
		}*/
	}
	//-----------------------------------------------------
	function on_pick(p_domain_info_map :Object,
		p_selected_domain_stats :PIXI.Container) {
		p_log_fun('FUN_ENTER','gf_domains_infos.init_selector_bar().on_pick()');

		const posts_count_int :number = p_domain_info_map['posts_count_int'];
		
		/*//---------------
		//DRAW POSTS COUNT TEXT
		final TextField posts_count_txt = new TextField();

		posts_count_txt.defaultTextFormat = new TextFormat('Arial', 18, Color.White);
		posts_count_txt.text     = '$posts_count_int';
		posts_count_txt.x        = 6;
		posts_count_txt.y        = 4;
		posts_count_txt.width    = 50;
		posts_count_txt.height   = 50;
		posts_count_txt.wordWrap = true;
		
		container.addChild(posts_count_txt);
		//---------------*/
		
		/*gf_domain.activate(p_domain_info_map,
						p_log_fun);*/
	}
	//-----------------------------------------------------
		
	container.interactive = true;

	var move_subscription;
	var current_i_int = 0;
	//container.addEventListener(MouseEvent.MOUSE_DOWN,
	
	container.on('mousedown',(p_e)=>{
			
			//move_subscription = container.addEventListener(MouseEvent.MOUSE_MOVE,(p_e)=>{
			container.on('mousemove',(p_e)=>{

				const mouse_y_int :number = parseInt(p_e.localY);

				//HACK!! - stagexl for some reason sometimes reports localY as 1 or 2. 
				//         so here Im filtering that out, to avoid flickering due 
				//         to sudden repositioning of scroll_indicator_gr
				if (mouse_y_int > 2) scroll_indicator_gr.y = mouse_y_int;


				const selected_map          :Object         = get_selected_domain_stats();
				const selected_domain_stats :PIXI.Container = selected_map['domain_stats'];
				const i                     :number         = selected_map['i'];

				//--------------------------
				//IMPORTANT!! - client function
					
				//only run the client function if the selected domain 
				//has actually changed... to avoid flickering/changing
				//data that has not changed
				if (i != current_i_int) {
					current_i_int = i;

					const domain_info_map :Object = p_domains_lst[i];

					on_pick(domain_info_map, selected_domain_stats);
					p_onPick_fun(domain_info_map);
				}
				//--------------------------
			});
		});

	//container.addEventListener(MouseEvent.MOUSE_UP,
	container.on('mouseup',
		(p_e)=>{
			if (move_subscription!=null) {
				move_subscription.cancel();
			}
		});
	return container;
}
//-----------------------------------------------------
function draw_domain_info(p_domain_info_map :Object,
	p_width_int  :number,
	p_height_int :number,
	p_color_int  :number,
	p_ctx_map    :Object,
	p_onPick_fun,
	p_log_fun) :PIXI.Container {
	//{int     p_width_int :20,
	//int      p_height_int:30,
	//int      p_color_int :Color.LightGray}) :PIXI.Sprite {
	//p_log_fun('FUN_ENTER','gf_domains_infos.draw_domain_info()');

	const items_count_int :number = p_domain_info_map['posts_count_int'] + p_domain_info_map['images_count_int'];
	const container               = new PIXI.Container();
	
	//---------------
	//DRAW BACKGROUND
	const graphics = new PIXI.Graphics();
	container.addChild(graphics);

	draw_background(graphics,
		p_width_int,
		p_height_int,
		p_color_int,
		p_log_fun);

	container.addChild(graphics);
	//---------------
	//SIZING
	if (items_count_int > 50) {
		container.width  = 70;
		container.height = 25;
	}
	else if (items_count_int > 10) {
		container.width  = 50;
		container.height = 20;
	}
	else {
		container.width  = 30;
		container.height = 10;
	}
	//---------------
	
	/*container.addEventListener(MouseEvent.MOUSE_OVER,
		(p_e) {
			query('#domain_stats_canvas').style.cursor = 'pointer';
		});
	container.addEventListener(MouseEvent.MOUSE_OUT,
		(p_e) {
			query('#domain_stats_canvas').style.cursor = '';
		});

	int base_x_int;
	int base_width_int  = container.width.toInt();
	int base_height_int = container.height.toInt();
	
	container.addEventListener(MouseEvent.CLICK,
		(p_e) {

			base_x_int = container.x.toInt();

			//--------------
			//DESELECT_OLD
			//deselect currently selected element
			if (p_ctx_map['selected_spr'] != null) {
				final Sprite old_selected_spr   = p_ctx_map['selected_spr'];
				final Shape  old_selected_shape = p_ctx_map['selected_shape'];

				old_selected_spr.width  = p_ctx_map['selected_spr_base_width_int'];
				old_selected_spr.height = p_ctx_map['selected_spr_base_height_int'];
				old_selected_spr.x      = p_ctx_map['selected_spr_base_x_int'];

				//draw old background back to normal color
				draw_background(old_selected_shape.graphics,
								p_color_int,
								p_log_fun,
								p_width_int :p_width_int,
								p_height_int:p_height_int);
			}
			//--------------
			//make just clicked element the selected one
			p_ctx_map['selected_spr']                 = container;
			p_ctx_map['selected_shape']               = shape;
			p_ctx_map['selected_spr_base_width_int']  = base_width_int;
			p_ctx_map['selected_spr_base_height_int'] = base_height_int;
			p_ctx_map['selected_spr_base_x_int']      = base_x_int;
			//--------------
			//set this domain_stats display at the top of the display hierarchy 
			final this_index_int = container.parent.getChildIndex(container);
			final last_index_int = container.parent.numChildren-1;
			container.parent.swapChildrenAt(this_index_int,last_index_int);
			//--------------
			//CHANGE BACKGROUND COLOR
			draw_background(shape.graphics,
							Color.Yellow,
							p_log_fun,
							p_width_int :p_width_int,
							p_height_int:p_height_int);
			//--------------
			//SCALE SIZE
			container.width  = p_width_int;
			container.height = p_height_int;


			container.x = (container.x + base_width_int)-container.width;
			//--------------
			

			p_onPick_fun(p_domain_info_map);
		},
		useCapture:true);*/

	return container;
}
//-----------------------------------------------------
function draw_background(p_graphics :PIXI.Graphics,
	p_width_int  :number,
	p_height_int :number,
	p_color_int  :number,
	p_log_fun) {
	//{int     p_width_int :20,
	//int      p_height_int:20}) {

	p_graphics.clear();
	p_graphics.moveTo(0, 0);

    p_graphics.beginFill(p_color_int); //.beginPath();
    //p_graphics.lineStyle(1,p_color_int);

    //single_page_height_px-1 - so that a little space is shown between pages
	p_graphics.drawRect(0,0, //x/y 
		p_width_int,   //p_width_px 
		p_height_int); //p_height_px

	p_graphics.endFill(); //.closePath();
	//p_graphics.strokeColor(p_color_int,1);
	//p_graphics.fillColor(p_color_int);
}
//-----------------------------------------------------
}