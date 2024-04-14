
///<reference path="../../../d/jquery.d.ts" />

import * as gf_viz_group         from "../ts/gf_viz_group";
import * as gf_viz_group_random_access from "../ts/gf_viz_group_random_access";

//-------------------------------------------------
$(document).ready(()=>{



    const test_elements_lst = [
        // page 1 (10 items)
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/b1b448df22b2767a8769f644f5f9e719_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/b1b448df22b2767a8769f644f5f9e719_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/b1b448df22b2767a8769f644f5f9e719_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },

        // page 2 (10 items)
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
        {
            "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
        },
    ];
    //-------------------------------------------------
    function element_create_fun(p_element_map) {

        const img_url_str = p_element_map["img_url_str"];
        
        // console.log(img_url_str)
        // console.log(p_element_container);

        const element = $(`<div><img src='${img_url_str}'></img></div>`);

        return element;
    }

    //-------------------------------------------------
    function elements_page_get_fun(p_page_index_int :number,
        p_pages_to_get_num_int :number) {
            
        const p = new Promise(function(p_resolve_fun, p_reject_fun) {

            const page_elements_lst = [
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/e217878d1817d7a314306aae2bf58abb_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
                {
                    "img_url_str": "https://gloflow.com/images/d/thumbnails/7d72ab16e6829ba4cb1fd866a898f625_thumb_medium.jpeg"
                },
            ];
            p_resolve_fun(page_elements_lst);
        });
        return p;
    }

    //-------------------------------------------------

    const id_str        = "test_viz_group";
    const parent_id_str = "test_parent";

    // number of initial pages that are supplied to gf_viz_group to display
    // before it has to initiate its own page fetching logic.
    const initial_pages_num_int = 2;
    
    const assets_uris_map = {
        "gf_bar_handle_btn": "./../../../../assets/gf_bar_handle_btn.svg",
    };


    
    const viz_props :gf_viz_group.GF_viz_props = {
        seeker_container_height_px: $(window).height(), // 500,
        seeker_container_width_px:  100,
        seeker_bar_width_px:        50, 
        seeker_range_bar_width:     30,
        seeker_range_bar_color_str: "red",
        assets_uris_map:            assets_uris_map,
        // seeker_range_bar_height: 500,
    }


    const props :gf_viz_group.GF_props = {

        container_id_str:        id_str,
        parent_container_id_str: parent_id_str,

        start_page_int:   0,
        end_page_int:     20,
        initial_page_int: 0,
        assets_uris_map:  assets_uris_map,
        viz_props: viz_props,
    };


    const seeker__container_element = gf_viz_group.init(test_elements_lst,
        props,
        element_create_fun,
        elements_page_get_fun);

    $(seeker__container_element).css("position", "fixed");
});