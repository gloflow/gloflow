///<reference path="../d/jquery.d.ts" />
///<reference path="../d/jqueryui.d.ts" />
///<reference path="../d/pixi.js.d.ts" />

namespace gf_domains_browser {

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

    //-----------------
    //LOAD_DATA
    const domains_infos_lst = [];
    $('#domains .domain').each((p_i,p_e)=>{

    	const domain_name_str  = $(p_e).find('.domain_name').text();
    	const domain_url_str   = $(p_e).find('.domain_url').text();
    	const posts_count_int  = parseInt($(p_e).find('.posts_count').text());
    	const images_count_int = parseInt($(p_e).find('.images_count').text());
    	domains_infos_lst.push({
    		'name_str'        :domain_name_str,
    		'url_str'         :domain_url_str,
    		'posts_count_int' :posts_count_int,
    		'images_count_int':images_count_int,
    	});
    });
    //-----------------

    gf_domains_browser.init(domains_infos_lst,
    					log_fun);
});
//-----------------------------------------------------
export function init(p_domains_infos_lst :Object[],
				p_log_fun) {
	p_log_fun('FUN_ENTER','gf_domains_browser.init()');

	$('#viz_container').on('click',()=>{
		init_viz(p_domains_infos_lst,p_log_fun);
	});
}
//-----------------------------------------------------
export function init_viz(p_domains_infos_lst :Object[],
					p_log_fun) {
	p_log_fun('FUN_ENTER','gf_domains_browser.init_viz()');

	//---------------------
	const element = `<div id='control'>
		<div id='domain_search'>
			<input type="text" id='query_input' placeholder='search domains'/>
		</div>
		<div id='selected_domain_info'>
			<div id='url'></div>
		</div>
		<canvas id='domain_stats_canvas' width="800" height="654"></canvas>

		<div id='domain_posts'>
			<div id='posts'></div>
		</div>
	</div>`;

	$('body').append(element);
	//---------------------

	const background_color_int :number = 0xFFE598; //gf_color.get_hex('whitesmoke');//0xFFFFE598;

    //const renderer = PIXI.autoDetectRenderer(1000, 2000, {antialias:true,backgroundColor:0x1099bb});
	//const renderer = new PIXI.WebGLRenderer(800, 654, {
	const renderer = new PIXI.CanvasRenderer(800,654,{
								//antialias      : true,
								backgroundColor:background_color_int, //backgroundColor:gf_color.get_hex("green"), //0x1099bb,
								view           :<HTMLCanvasElement> $('#domain_stats_canvas')[0]});

    const width_int  :number = $('#domain_stats_canvas').width();
    const height_int :number = $('#domain_stats_canvas').height();

    //-----------------
	//DOMAIN SEARCH
	
	gf_domains_search.init_domain_search(p_domains_infos_lst,
			(p_domain_info_map :Object)=>{
				pick_domain(p_domain_info_map);
			},
			p_log_fun);
	//-----------------
	//SORT DATA
	p_domains_infos_lst.sort((a,b)=>{
		const a_total = a['posts_count_int'] + a['images_count_int'];
		const b_total = b['posts_count_int'] + b['images_count_int'];
		if (a_total > b_total) {
			return -1;
		}
		else if (a_total < b_total) {
			return 1;
		}
		else {
			return 0;
		}
	});
	//-----------------

	const ctx_map          = {};
	const info_map :Object = draw_domains_stats(p_domains_infos_lst,
										width_int-300, //max_width
										100,           //max_height
										ctx_map,
										//domain_onPick_fun,
										(p_domain_info_map :Object)=>{
											pick_domain(p_domain_info_map);
										},
										p_log_fun);

    const domains_stats :PIXI.Container = info_map['container'];
    const domains_infos :PIXI.Container = info_map['domains_infos'];

    domains_infos.x = 733;

    const connectivity_spr :PIXI.Container = gf_domains_conn.draw_connectivity(p_domains_infos_lst,
																domains_infos,
																width_int,
																height_int,
																background_color_int,
																p_log_fun);

    const stage = new PIXI.Container();
    stage.addChild(connectivity_spr);
    stage.addChild(domains_stats);

    renderer.render(stage);

    //-----------------------------------------------------
    function create_post(p_post_title_str :string) :HTMLDivElement {
    	//p_log_fun('FUN_ENTER','gf_domains_browser.init().create_post()');

    	const url_encoded_title_str :string = encodeURIComponent(p_post_title_str);
    	const post                          = $(`<div id='post'></div>`);

    	//--------------------
    	//Dart doesnt allow <a> tags to be included in html text, so <a> has
    	//to be created/added manually
    	const a = new HTMLAnchorElement();
    	$(a).attr('href','http://www.gloflow.com/posts/'+url_encoded_title_str);
    	$(a).text(p_post_title_str);

    	//'_blank' - open the link in a new window
    	$(a).attr('target','_blank');
    	//--------------------
    	
    	$(post).append(a);
    	return <HTMLDivElement> post.get(0);
    }
    //-----------------------------------------------------
	function pick_domain(p_domain_info_map :Object) {
		p_log_fun('FUN_ENTER','gf_domains_browser.init().pick_domain()');


		const posts_count_int :number = p_domain_info_map['posts_count_int'];

		
		const posts_count_element = $('#control #selected_domain_info #posts_count');
		if (posts_count_element != null) posts_count_element.remove();

		$('#control #selected_domain_info').append($(`
			<div id='posts_count'>
				<span>posts #</span>`+posts_count_int+`
			</div>`));

		//-------------
		//LINK
		
		var a = $('#control #selected_domain_info #url a');
		if (a != null) a.remove();

		const domain_url_str :string = p_domain_info_map['url_str'];

		const domain_a = $('<a href='+domain_url_str+' target="_blank">'+domain_url_str+'</a>');
		$('#control #selected_domain_info #url').append(domain_a);
		//-------------
		
    	const domain_posts_element = $('#control').find('#domain_posts');
    	$(domain_posts_element).find('#posts').remove();

    	const new_posts_element = $('<div id="posts"></div>');
    	$(domain_posts_element).append(new_posts_element);

    	
    	for (var i=0;i<posts_count_int;i++) {
    		
    		const post_title_str :string         = p_domain_info_map['posts_titles_lst'][i];
    		const post_element   :HTMLDivElement = create_post(post_title_str);

    		$(new_posts_element).append(post_element);
    	}
    }
    //-----------------------------------------------------
}
//-----------------------------------------------------
function draw_domains_stats(p_domains_lst :Object[],
						p_item_max_width_int   :number,
						p_items_max_height_int :number,
						p_ctx_map              :Object,
						p_onPick_fun,
						p_log_fun) :Object {
	p_log_fun('FUN_ENTER','gf_domains_browser.draw_domains_stats()');

	const container :PIXI.Container = new PIXI.Container();

	//----------------------
	//MAX POSTS COUNT
	//get the number of items in the domain with the most items.
	//number of items is the sum of number of posts and images

	var max_items_count_int :number = 0;
	for (var domain_map of p_domains_lst) {
		const c :number = domain_map['posts_count_int'] + domain_map['images_count_int'];
		if (c > max_items_count_int) {
			max_items_count_int = c;
		}
	}
    //----------------------
    //DRAW DOMAINS
    
    var i                      :number = 0;
	var last_spr               :PIXI.Container;
	var horizontal_domains_int :number = 0;
	var conseq_right_moves_int :number = 0;
	var conseq_left_moves_int  :number = 0;
	var conseq_down_moves_int  :number = 0;

	function random_bool(){return Math.random() >= 0.5}

	console.log('=============================')
	console.log(p_items_max_height_int)
	console.log(max_items_count_int)

    for (var p_domain_info_map of p_domains_lst) {
		const domain_name_str  :string = p_domain_info_map['name_str'];
  		const posts_count_int  :number = p_domain_info_map['posts_count_int'];
  		const images_count_int :number = p_domain_info_map['images_count_int'];
  		const items_count_int  :number = posts_count_int + images_count_int;

  		//------------------
  		//DOMAIN HEIGHT
  		//calculate the height of the domain, relative to the maximum (most popular) domain height

    	//100:max_items_count_int=x:items_count_int
    	//x = (100*items_count_int)/max_items_count_int
    	//final int    relative_width_int  = ((p_item_max_width_int*items_count_int)/max_items_count_int).floor();
		const domain_relative_height_int :number = Math.floor((p_items_max_height_int*items_count_int)/max_items_count_int);
		const domain_relative_width_int  :number = domain_relative_height_int;
		//------------------
		//DOMAIN IMAGES INDICATOR WIDTH
		//width of the domain is the relative with of the items count (both images and posts).
		//here the width of just images is calculated, so that the user can view 
		//what percentage of the domain items is in images

		//domain_relative_width_int:items_count_int=x:images_count_int
		//domain_relative_width_int*images_count_int = items_count_int*x
		//x=(domain_relative_width_int*images_count_int)/items_count_int
		const domain_images_indicator_height_int :number = Math.floor((domain_relative_width_int*images_count_int)/items_count_int);
		//------------------
    	//COLOR
    	var color_int :number;
    	if (domain_relative_height_int > 50) {
    		color_int = gf_color.get_int('lightblue') + i*10
    	}
    	else if (domain_relative_height_int > 1) {
    		color_int = gf_color.get_int('green') + i*2;
    	}
    	else {
    		color_int = gf_color.get_int('orange');
    	}
    	//--------------
		const domain_spr :PIXI.Container = gf_domain.draw(domain_name_str,
													domain_relative_width_int,
													domain_relative_height_int,
													domain_images_indicator_height_int,
													color_int, //0xFF8851+i*10,
													p_log_fun); 
    	//--------------
    	//LAYOUT

    	if (i>0) {


    		if (domain_spr.width > 30) {

	    		if ((last_spr.width+domain_spr.width) < 400) {

	    			if (conseq_right_moves_int < 1) { 
		    			domain_spr.x = last_spr.x + last_spr.width + 5; //MOVE RIGHT
		    			domain_spr.y = last_spr.y;

		    			conseq_right_moves_int += 1;
		    			conseq_left_moves_int   = 0;
		    			conseq_down_moves_int   = 0;
		    		}
		    		else {
		    			domain_spr.x = last_spr.x; //- (domain_spr.width + 5); 
		    			domain_spr.y = last_spr.y + last_spr.height + 5; //MOVE DOWN

		    			conseq_right_moves_int = 0;
		    			conseq_left_moves_int  = 0;
		    			conseq_down_moves_int += 1;
		    		}
	    		}
	    		else {
	    			domain_spr.x = last_spr.x;
	    			domain_spr.y = last_spr.y + last_spr.height + 5; //MOVE DOWN

	    			conseq_right_moves_int = 0;
	    			conseq_left_moves_int  = 0;
	    			conseq_down_moves_int += 1;
	    		}
	    	}
	    	else if (domain_spr.width > 10) {
	    		if (conseq_right_moves_int > 2) {
	    			domain_spr.x = last_spr.x - (domain_spr.width + 5); //MOVE LEFT
	    			domain_spr.y = last_spr.y + 5;

					conseq_right_moves_int  = 0;
					conseq_left_moves_int  += 1;
					conseq_down_moves_int   = 0;
	    		}
	    		else {
	    			domain_spr.x = last_spr.x - domain_spr.width - 5;


	    			if (random_bool()) {
	    				domain_spr.x = last_spr.x - (domain_spr.width+5); //MOVE LEFT
	    				domain_spr.y = last_spr.y; //+ domain_spr.height + 5;

	    				conseq_right_moves_int  = 0;
						conseq_left_moves_int  += 1;
						conseq_down_moves_int   = 0;
	    			}
	    			else {

	    				if (random_bool()) {
		    				domain_spr.x = last_spr.x;
		    				domain_spr.y = last_spr.y + last_spr.height + 5; //MOVE DOWN

		    				conseq_right_moves_int = 0; 
							conseq_left_moves_int  = 0;
							conseq_down_moves_int += 1;
						}
						else {
							domain_spr.x = last_spr.x;
		    				domain_spr.y = last_spr.y + (last_spr.height + 5); //MOVE DOWN

		    				conseq_right_moves_int = 0; 
							conseq_left_moves_int  = 0;
							conseq_down_moves_int  = 0;
						}
	    			}
	    		}
	    	}
	    	else {

	    		//randomly enter this path
	    		if (random_bool()) {

	    			if (random_bool()) {
			    		domain_spr.x = last_spr.x;
			    		domain_spr.y = last_spr.y + last_spr.height + 5; //MOVE DOWN

			    		conseq_right_moves_int  = 0; 
						conseq_left_moves_int   = 0;
						conseq_down_moves_int  += 1;
			    	}
			    	else {
			    		domain_spr.x = last_spr.x;
			    		domain_spr.y = last_spr.y + (last_spr.height + 5); //MOVE DOWN

			    		conseq_right_moves_int = 0; 
						conseq_left_moves_int  = 0;
						conseq_down_moves_int  = 0;
			    	}
		    	}
		    	else {
		    		//if there are less the 5 consequent left moves
		    		//if (random_bool() && conseq_left_moves_int < 5) {
		    		if (conseq_left_moves_int < 5) {
			    		domain_spr.x = last_spr.x + last_spr.width + 5; //MOVE RIGHT
			    		domain_spr.y = last_spr.y;

			    		conseq_right_moves_int += 1;
						conseq_left_moves_int   = 0;
						conseq_down_moves_int   = 0;
			    	}
			    	else {
			    		domain_spr.x = last_spr.x - (last_spr.width + 5); //MOVE LEFT
			    		domain_spr.y = last_spr.y;

			    		conseq_right_moves_int  = 0; //+=1; //3rd horizontal element
						conseq_left_moves_int  += 1;
						conseq_down_moves_int   = 0;
			    	}
		    	}
	    	}
    	}
    	//--------------

    	container.addChild(domain_spr);
    	
    	i+=1;

    	p_domain_info_map['container'] = domain_spr;
    	last_spr                       = domain_spr;
    }
	//----------------------
	//DRAW DOMAINS INFOS	
	const domains_infos :PIXI.Container = gf_domains_infos.draw(p_domains_lst,
													200, //p_width_int
													600, //p_height_int
													p_ctx_map,
													p_onPick_fun,
													p_log_fun);
	container.addChild(domains_infos);
	//domains_infos.x = p_domain_infos_x_int;
	//----------------------
    var selected_spr :PIXI.Sprite;
    for (var p_domain_info_map of p_domains_lst) {
		const container = p_domain_info_map['stats_container'];


	}
    //----------------------
    return {
    	'container'    :container,
    	'domains_infos':domains_infos
    };
}
//-----------------------------------------------------
}