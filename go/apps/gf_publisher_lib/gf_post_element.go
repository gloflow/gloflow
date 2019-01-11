package gf_publisher_lib

import (
	"errors"
	"fmt"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
)
//---------------------------------------------------
type Post_element struct {

	Id_str string `bson:"id_str"`

	//type_str - "link"|"image"|"video"|"iframe"|"text"
	Type_str        string `bson:"type_str"`
	Description_str string `bson:"description_str"`

	//post_elements can be created after/before their hosting post has been created
	//so their creation datetimes might be different then the post creation datetime
	Creation_datetime_str string `bson:"creation_datetime_str"`

	//FIX!! - if type_str == "image" this url_str should be the external source link if 
	//        the PostElement_ADT is a composite of external stuff
	//----------------------
	//if type_str == "link"|"image"|"video"|"iframe" then PostElement_ADT has 
	//a external url associated with it
	Extern_url_str string `bson:"extern_url_str"`
	//----------------------
	//if type_str == "image"|"video" then source_page_url_str represents the URL of
	//the page from which this post_element was extracted from (if it wasnt uploaded directly)
	Source_page_url_str string `bson:"source_page_url_str"`
	//----------------------
	//GEOMETRIC PROPS
		
	//this is the index unique for the element, in a maximum of 3d space
	//lower orders (1d,2d) are done in the 3d tuple (x,y,0)
	//this is used for graphical/positioning ops
	//FIX!! - postfix is "_tpl" for legacy reasons. should be "_lst"
	Post_index_3_lst []int `bson:"post_index_3_lst"`
	Width_int        int   `bson:"width_int"`  //in pixels
	Height_int       int   `bson:"height_int"` //in pixels
	//----------------------
	//IMAGE - if type_str == "image"

	Image_id_str string

	//only thumbnail urls are tracked here in the Post_ADT, not the full-size (which is tracked
	//in Image_ADT), since the fullsize internal url is never used (that would be copyright infringement).
	//using thumbnails falls into fair-use
	Img_thumbnail_small_url_str  string `bson:"img_thumbnail_small_url_str"`
	Img_thumbnail_medium_url_str string `bson:"img_thumbnail_medium_url_str"`
	Img_thumbnail_large_url_str  string `bson:"img_thumbnail_large_url_str"`
	//----------------------

	Tags_lst   []string               `bson:"tags_lst"`
	Colors_lst []string               `bson:"colors_lst"`
	Meta_map   map[string]interface{} `bson:"meta_map"`
}
//---------------------------------------------------
func create_post_elements(p_post_elements_infos_lst []interface{},
	p_post_title_str string,
	p_runtime_sys    *gf_core.Runtime_sys) []*Post_element {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_element.create_post_elements()")
	p_runtime_sys.Log_fun("INFO",     "p_post_elements_infos_lst - "+fmt.Sprint(p_post_elements_infos_lst))

	post_elements_lst := []*Post_element{}
	for i,post_element := range p_post_elements_infos_lst {

		post_element_map := post_element.(map[string]interface{})
		//--------------------
		//REMOVE!! - 3d positioning will be a part of gf_3d data structures
		
		//1d index stored in the 3d index slot
		//this is the placement order for the post element
		post_index_3_lst := []int{i,0,0,}
		//post_element_map["post_index_3_lst"] = []int{i,0,0,}
		//--------------------
		//POST_ELEMENT_ID
	
		//FIX!! - post_index_3_str shouldnt be serialized in ID as a list
		//example ID output - "pub_pe:test_post_"title:[0, 0, 0]"
		post_index_3_str    := fmt.Sprint(post_index_3_lst) //post_element_map["post_index_3_lst"])
		post_element_id_str := fmt.Sprintf("pub_pe:%s:%s", p_post_title_str, post_index_3_str)
		//------------------			
		creation_datetime_str := time.Now().String()
		extern_url_str        := post_element_map["extern_url_str"].(string)

		p_runtime_sys.Log_fun("INFO","post_element extern_url_str - "+fmt.Sprint(extern_url_str))
		//------------------
		post_element := &Post_element{
			Id_str:               post_element_id_str,
			Type_str:             post_element_map["type_str"].(string),
			Creation_datetime_str:creation_datetime_str,
			Extern_url_str:       extern_url_str,
			Source_page_url_str:  post_element_map["source_page_url_str"].(string),
			Post_index_3_lst:     post_index_3_lst,
			//Description_str      :post_element_map["description_str"].(string),
		}
		
		post_elements_lst = append(post_elements_lst,post_element)
	}

	return post_elements_lst
}
//---------------------------------------------------
func get_first_image_post_element(p_post *Post, p_runtime_sys *gf_core.Runtime_sys) *Post_element {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_element.get_first_image_post_element()")

	for _,post_element := range p_post.Post_elements_lst {
		if post_element.Type_str == "image" {
			return post_element
		}
	}
	return nil //post has no image post_element
}
//---------------------------------------------------
func get_post_elements_of_type(p_post *Post,
	p_type_str    string,
	p_runtime_sys *gf_core.Runtime_sys) ([]*Post_element,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_element.get_post_elements_of_type()")
	
	if !(p_type_str == "image"  ||
		p_type_str == "link"   ||
		p_type_str == "video"  ||
		p_type_str == "iframe" ||
		p_type_str == "text") {
		gf_err := gf_core.Error__create(fmt.Sprintf("post_element type_str not of value image|link|video|iframe|text - instead its - ",p_type_str),
			"verify__invalid_value_error",
			&map[string]interface{}{"type_str":p_type_str,},
			nil, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	post_elements_lst := []*Post_element{}
	for _,post_element := range p_post.Post_elements_lst {
		if post_element.Type_str == p_type_str {
			post_elements_lst = append(post_elements_lst,post_element)
		}
	}
	return post_elements_lst,nil
}
//---------------------------------------------------
/*func create_extern_post_element(p_post_element_info_map map[string]interface{},
					p_post_title_str                  *string,
					p_gf_images_main_service_host_str *string,
					p_mongodb_coll                    *mgo.Collection,
					p_log_fun                         func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_element.create_extern_post_element()")
	
	post_element := create_post_elements([p_post_element_info_dict], //p_post_elements_lst,
								p_post_title_str,
								p_log_fun)[0]
	
	post,err := db__get_post(p_post_title_str,
						p_mongodb_coll,
						p_log_fun)
	if err != nil {
		return err
	}

	if post_element.Type_str == "image" {
		init_image_post_element(post_element,
							post,
							p_gf_images_main_service_host_str,
							p_log_fun)
	}

	post.Post_elements_lst = append(post.Post_elements_lst,post_element)
	
	err := db__update_post(post,
				p_mongodb_coll,
				p_log_fun)
	if err != nil {
		return err
	}
}*/
//---------------------------------------------------
/*func init_image_post_element(p_post_element *Post_element,
					p_post                            *Post,
					p_gf_images_main_service_host_str *string,
					p_log_fun                         func(string,string)) (*Post_element,error) {
	p_log_fun("FUN_ENTER","gf_post_element.init_image_post_element()")
	

		image_url_str := ""


		if p_post_element.Type_str == "image" {
			return nil,errors.New("post_element is not an 'image' type - "+p_post_element.Type_str)
		}

		assert(p_post_element.extern_url_str != null);
		
		image_url_str = p_post_element.extern_url_str;
		
		//gf_images_client.dispatch_process_extern_image() - sends out a HTTP client to the gf_images service
		//                                                   to dispatch the processing of this external image
		
		final f = gf_images_lib.Client__dispatch_process_extern_images(image_url_str,
														p_log_fun,
														p_reprocess_if_prexisting_bool   :false,
														p_gf_images_main_service_host_str:p_gf_images_main_service_host_str)
		return f;
	})
	.then((Map p_result_dict) {
		p_log_fun("INFO","result_dict - $p_result_dict");

		final String new_image_id_str   = p_result_dict["image_id_str"];
		p_post_element.image_id_str = new_image_id_str;
	
		//images_ids_lst - holds id"s of all images in the post
		p_post_adt.images_ids_lst.add(p_post_element.image_id_str);
		
		p_post_element.img_thumbnail_small_url_str  = p_result_dict["thumbnail_small_relative_url_str"];
		p_post_element.img_thumbnail_medium_url_str = p_result_dict["thumbnail_medium_relative_url_str"];
		p_post_element.img_thumbnail_large_url_str  = p_result_dict["thumbnail_large_relative_url_str"];

		completer.complete(p_post_element);
	})
	//------------------
	//ERROR HANDLING
	.catchError((p_error) {
		p_log_fun("ERROR","PROCESSING IMAGE POST_ELEMENT FAILED!! [$image_url_str]");
		p_log_fun("INFO" ,p_error.toString());

		//ADD!! - in case of an error in generating a thumbnail, these properties should 
		//        not be set to "error" but instead to some generic thumbnail image url (so that 
		//        in the final html rendering of the post_element a valid image is shown regardless, 
		//        and not a broken image link)
		
		p_post_element.image_id_str                 = "error";
		p_post_element.img_thumbnail_small_url_str  = "error";
		p_post_element.img_thumbnail_medium_url_str = "error";
		p_post_element.img_thumbnail_large_url_str  = "error";
		
		p_post_adt.images_ids_lst.add(p_post_element.image_id_str);

		//IMPORTANT!! - im not raising any errors here because I dont want the failure to process
		//              any one of the images to prevent the creation of entire post_element. 
		//              The processing of the image can be restarted later.
		completer.complete(p_post_element);
	});
	//------------------
}*/