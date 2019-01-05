//---------------------------------------------------
//ADD!! - detect you tube embeds in other non-youtube.com pages
//        via the <embed> tag

//->:List<:Dict(video_info_map)>
function get_videos_info(p_log_fun) {
	p_log_fun('FUN_ENTER','page_info_extraction.get_videos_info()')

	const page_url_str    = window.location.toString();
	const videos_info_lst = [];
	
	//------------------------------------
	//YOUTUBE.COM DOMAIN
	
	//the user is currently watching a video on youtube.com
	//"\/\/" is "//" escaped
	if (page_url_str.match("^http:\/\/www.youtube.com\/watch") || 
		page_url_str.match("^https:\/\/www.youtube.com\/watch")) {

		const youtube_video_embed_url_str = $('link[itemprop="embedURL"]').attr('href');
		p_log_fun('INFO','youtube_video_embed_url_str:'+youtube_video_embed_url_str);
		
		const video_info_map = {
			'type_str'        :'video',
			'page_url_str'    :page_url_str,
			'video_source_str':'youtube',
			'embed_url_str'   :youtube_video_embed_url_str
		};
		
		videos_info_lst.push(video_info_map);
	}
	//------------------------------------
	//if its any other page, search all elements 'src' attribute
	//and see if it contains "http://www.youtube.com", "http://player.vimeo.com"
	else {
		
		//"*=" - Attribute Contains Selector [name*="value"]
		//Selects elements that have the specified attribute 
		//with a value containing the a given substring.
		//------------------------------------
		//YOUTUBE - IFRAME EMBED
		
		$('*[src*="https://www.youtube.com"]').each((p_i,p_element) => {
			p_log_fun('INFO','YOUTUBE IFRAME EMBED++++++++++++++++++++++++++++++++');
			p_log_fun('INFO',$(p_element).attr('src'));
			
			const video_info_map = {
				'type_str'        :'video',
				'page_url_str'    :page_url_str,
				'video_source_str':'youtube',
				'embed_url_str'   :$(p_element).attr('src')
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
		//VIMEO - IFRAME EMBED
		
		$('*[src*="http://player.vimeo.com"]').each((p_i,p_element) => {
			p_log_fun('INFO','VIMEO IFRAME EMBED++++++++++++++++++++++++++++++++');
			p_log_fun('INFO',$(p_element).attr('src'));
			
			const video_info_map = {
				'type_str'        :'video',
				'page_url_str'    :page_url_str,
				'video_source_str':'vimeo',
				'embed_url_str'   :$(p_element).attr('src')
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
		//VIMEO - FLASH PLAYER (<OBJECT> TAG)
		
		$('object[data*="http://a.vimeocdn.com"]').each((p_i,p_element) => {
			p_log_fun('INFO','VIMEO FLASH PLAYER OBJECT TAG++++++++++++++++++++++++++++++++');
			p_log_fun('INFO',$(p_element).attr('data'));
			
			const flash_vars_str  = $(p_element).find('*[name*="flashvars"]').attr('value');
			const video_info_map = {
				'type_str'            :'video',
				'page_url_str'        :page_url_str,
				'video_source_str'    :'vimeo',
				'vimeo_flash_vars_str':flash_vars_str
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
		//OOYALA - IFRAME EMBED
		
		$('*[src*="http://player.ooyala.com"]').each((p_i,p_element) => {
			p_log_fun('INFO','OOYALA IFRAME EMBED++++++++++++++++++++++++++++++++');
			p_log_fun('INFO',$(p_element).attr('src'));
			
			const video_info_map = {
				'type_str'        :'video',
				'page_url_str'    :page_url_str,
				'video_source_str':'ooyala',
				'embed_url_str'   :$(p_element).attr('src')
			};
			
			videos_info_lst.push(video_info_map);
		});
		//------------------------------------
	}
	//------------------------------------

	return videos_info_lst;
}
//---------------------------------------------------
//->:List<:Dict()>(img_infos_lst)
function get_images_info(p_log_fun) {
	p_log_fun('FUN_ENTER','page_info_extraction.get_images_info()')
	
	const page_url_str        = window.location.toString();
	const min_image_dimension = 20;

	//---------------------------------------------------
	function get_image_info(p_jq_element) {
		p_log_fun('FUN_ENTER','page_info_extraction.get_images_info().get_image_info()')
		
		//".src" instead of ".attr('src')" - gets the fully resolved url (including the host)
		//                                   and not just the value thats stored in the "src" html attr
		const full_img_src_str = $(p_jq_element)[0].src;
		const img_name_str     = full_img_src_str.split('/').pop();
		const img_width        = $(p_jq_element).width();
		const img_height       = $(p_jq_element).height();

		const img_info_map = {
			'type_str'        :'image',
			'page_url_str'    :page_url_str,
			'full_img_src_str':full_img_src_str,
			'img_name_str'    :img_name_str,
			'img_width'       :img_width,
			'img_height'      :img_height
		};
		console.log(img_info_map);

		return img_info_map;
	}
	//---------------------------------------------------

	const img_infos_lst = [];
	$('img').each((p_i,p_element) => {
		const img_info_map = get_image_info(p_element);
		
		if (img_info_map['img_width']  > min_image_dimension &&
			img_info_map['img_height'] > min_image_dimension) {

			//only use the image if both of its dimensions are larger
			//then the minimum treshold
			img_infos_lst.push(img_info_map);
		}
	});

	return img_infos_lst;
}