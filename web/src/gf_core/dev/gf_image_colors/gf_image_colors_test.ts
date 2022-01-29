







///<reference path="./../../../d/jquery.d.ts" />

import * as gf_image_colors  from "./../../ts/gf_image_colors";

//-------------------------------------------------
$(document).ready(()=>{



    const img = $(".image_info").find("img")[0];
    const assets_paths_map = {
        "copy_to_clipboard_btn": "./../../../../assets/gf_copy_to_clipboard_btn.svg"
    };
	gf_image_colors.init_pallete(img,
        assets_paths_map,
        (p_color_dominant_hex_str,
        p_colors_hexes_lst)=>{
            


            // $(".image_info").css("background-color", p_color_dominant_hex_str)




        });


});