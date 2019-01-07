package gf_images_lib 

import (
	"apps/gf_images_lib/gf_images_utils"
	"gf_core"
)
//---------------------------------------------------
func Add_tags_to_image(p_image *gf_images_utils.Gf_image,
					p_tags_lst    []string,
					p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_image.Add_tags_to_image()")
	
	if len(p_tags_lst) > 0 {

		//add all new tags with the current tags associated with an image,
		//with possible duplicates existing
		p_image.Tags_lst = append(p_image.Tags_lst,p_tags_lst...)

		//-----------
		set := map[string]bool{}
		for _,t_str := range p_image.Tags_lst {
			set[t_str]=true
		}
		//-----------
		list_no_duplicates_lst := []string{}
		for k_str,_ := range set {
			list_no_duplicates_lst = append(list_no_duplicates_lst,k_str)
		}

		//eliminate duplicates from the list
		p_image.Tags_lst = list_no_duplicates_lst
		//-----------
	}
}