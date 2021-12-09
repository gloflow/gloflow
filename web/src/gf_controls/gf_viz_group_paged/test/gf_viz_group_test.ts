
///<reference path="../../../d/jquery.d.ts" />

import * as gf_viz_group_paged from "./../ts/gf_viz_group_paged";

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

    gf_viz_group_paged.init(id_str,
        parent_id_str,
        test_elements_lst,
        initial_pages_num_int,
        assets_uris_map,
        element_create_fun,
        elements_page_get_fun);
});