package gf_crawl_core

import (
	"os"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_utils"
)

//--------------------------------------------------
func image__update_after_process(p_page_img *Crawler_page_img,
							p_gf_image_id_str string,
							p_runtime_sys     *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_utils.image__update_after_process()")

	p_page_img.Valid_for_usage_bool = true
	p_page_img.Image_id_str         = p_gf_image_id_str
	
	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
					"t":     "crawler_page_img",
					"id_str":p_page_img.Id_str,
				},
				bson.M{"$set":bson.M{
						//IMPORTANT!! - gf_image has been created for this page_image, and so the appropriate
						//              image_id_str needs to be set in the page_image DB record
						"image_id_str":p_gf_image_id_str,

						//IMPORTANT!! - image has been transformed, and is ready to be used further
						//              by other apps/services, either for display, or further calculation
						"valid_for_usage_bool":true,
					},
				})

	if err != nil {
		gf_err := gf_core.Error__create("failed to update an crawler_page_img valid_for_usage flag and its image_id (Gf_image) by its ID",
			"mongodb_update_error",
			&map[string]interface{}{
				"id_str":         p_page_img.Id_str,
				"gf_image_id_str":p_gf_image_id_str,
			},err,"gf_crawl_core",p_runtime_sys)
		return gf_err
	}
	return nil
}
//--------------------------------------------------
func image__cleanup(p_img_local_file_path_str string,
			p_img_thumbs  *gf_images_utils.Gf_image_thumbs,
			p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_utils.image__cleanup()")

	files_to_remove_lst := []string{
		p_img_local_file_path_str,
	}

	if p_img_thumbs != nil {
		files_to_remove_lst = append(files_to_remove_lst,p_img_thumbs.Small_local_file_path_str)
		files_to_remove_lst = append(files_to_remove_lst,p_img_thumbs.Medium_local_file_path_str)
		files_to_remove_lst = append(files_to_remove_lst,p_img_thumbs.Large_local_file_path_str)
	}
	
	for _,f_str := range files_to_remove_lst {
		err := os.Remove(f_str)
		if err != nil {
			gf_err := gf_core.Error__create("failed to cleanup a crawled image files",
				"file_remove_error",
				&map[string]interface{}{"file_str":f_str,},
				err,"gf_crawl_core",p_runtime_sys)
			return gf_err
		}
	}
	return nil
}