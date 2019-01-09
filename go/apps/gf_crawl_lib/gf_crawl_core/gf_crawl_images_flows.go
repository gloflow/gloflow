package gf_crawl_core

import (
	"fmt"
	"path/filepath"
	"github.com/davecgh/go-spew/spew"
	"gf_core"
	"apps/gf_images_lib"
	"apps/gf_images_lib/gf_images_utils"
)
//--------------------------------------------------
func Flows__add_extern_image(p_crawler_page__gf_image_id_str string,
	p_flows_names_lst                   []string,
	p_crawled_images_s3_bucket_name_str string,
	p_gf_images_s3_bucket_name_str      string,
	p_runtime                           *Crawler_runtime,
	p_runtime_sys                       *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_flows.Flows__add_extern_image()")


	//this is used temporarily to donwload images to, before upload to S3
	images_store_local_dir_path_str := "."

	fmt.Println("p_crawler_page__gf_image_id_str - "+p_crawler_page__gf_image_id_str)
	fmt.Println("p_flows_names_lst               - "+fmt.Sprint(p_flows_names_lst))

	gf_page_img,gf_err := image__db_get(p_crawler_page__gf_image_id_str,p_runtime,p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}


	gf_image_id_str := gf_page_img.Image_id_str
	gf_images_s3_bucket__is_uploaded_to__bool := false

	//IMPORTANT!! - some crawler_page_images dont have their gf_image_id_str set,
	//              which means that they dont have their corresponding gf_image.
	if gf_image_id_str == "" {



		p_runtime_sys.Log_fun("INFO","")
		p_runtime_sys.Log_fun("INFO","CRAWL_PAGE_IMAGE MISSING ITS GF_IMAGE --- STARTING_PROCESSING")
		p_runtime_sys.Log_fun("INFO","")

		//S3_UPLOAD - images__process_crawler_page_image() uploads image and its thumbs to S3 
		//            after it finishes processing it.
		gf_image,gf_image_thumbs,local_image_file_path_str,gf_err := images_pipe__single_simple(gf_page_img,
																				images_store_local_dir_path_str,
																				p_crawled_images_s3_bucket_name_str,
																				p_runtime,
																				p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		gf_image_id_str = gf_image.Id_str
		//-------------------
		//S3_UPLOAD_TO_GF_IMAGES_BUCKET
		//IMPORTANT!! - crawler_page_image and its thumbs are uploaded to the crawled images S3 bucket,
		//              but gf_images service /images/d endpoint redirects users to the gf_images
		//              S3 bucket (gf--img).
		//              so we need to upload the new image to that gf_images S3 bucket as well.
		//FIX!! - too much uploading, very inefficient, figure out a better way!

		gf_err = gf_images_utils.S3__store_gf_image(local_image_file_path_str,
										gf_image_thumbs,
										p_gf_images_s3_bucket_name_str,
										p_runtime.S3_info,
										p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}

		//IMPORTANT!! - gf_images_utils.Trans__s3_store_image() uploads the image and its thumbs to S3, 
		//              to indicate that we dont need to upload it later again.
		gf_images_s3_bucket__is_uploaded_to__bool = true
		//-------------------
		//CLEANUP

		gf_err = image__cleanup(local_image_file_path_str,gf_image_thumbs,p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
		//-------------------
	}
	//--------------------------
	//ADD_FLOWS_TO_DB

	//IMPORTANT!! - for each flow que up a DB operation. they will most likely
	//              be sent out to the DB server together once the buffer fills up.
	for _,flow_name_str := range p_flows_names_lst {

		gf_err := gf_images_lib.Flows_db__add_flow_to_image(flow_name_str,
													gf_image_id_str,
													p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
	}
	//--------------------------
	//S3_COPY_BETWEEN_BUCKETS - gf--discovered--img -> gf--img
	//                          only for gf_images that have not already been uploaded to the gf--img bucket
	//                          because they needed to be reprecossed and were downloaded from a URL onto
	//                          the local FS first.

	if !gf_images_s3_bucket__is_uploaded_to__bool {

		gf_image,gf_err := gf_images_utils.DB__get_image(gf_image_id_str,p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}


		fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
		spew.Dump(gf_image)




		/*S3__get_image_original_file_s3_filepath is wrong!! FIXX!!!
		the path of the originl_file that its returning is of a file named by its gf_img ID, which is wrong. 
		that filename should be of the file as it was original found in a html page or elsewhere.
		all images added via browser extension, or added by crawler, are named with original file name, not with ID.
		figure out if fixing this is going to break already added images (images added to a flow here from crawled images), 
		since they're all named by ID now (which is a bug)*/

		original_file_s3_filepath_str                                := gf_images_utils.S3__get_image_original_file_s3_filepath(gf_image,p_runtime_sys)
		t_small_s3_path_str,t_medium_s3_path_str,t_large_s3_path_str := gf_images_utils.S3__get_image_thumbs_s3_filepaths(gf_image,p_runtime_sys)



		fmt.Println("original_file_s3_filepath_str - "+original_file_s3_filepath_str)
		fmt.Println("t_small_s3_path_str           - "+t_small_s3_path_str)
		fmt.Println("t_medium_s3_path_str          - "+t_medium_s3_path_str)
		fmt.Println("t_large_s3_path_str           - "+t_large_s3_path_str)




		//ADD!! - copy t_small_s3_path_str first, and then copy original_file_s3_filepath_str and medium/large thumb in separate goroutines
		//        (in parallel and after the response returns back to the user). 
		//        this is critical to improve perceived user response time, since the small thumb is necessary to view an image in flows, 
		//        but the original_file and medium/large thumbs are not (and can take much longer to S3 copy without the user noticing).
		files_to_copy_lst := []string{
			original_file_s3_filepath_str,
			t_small_s3_path_str, 
			t_medium_s3_path_str,
			t_large_s3_path_str,
		}
		


		for _,s3_path_str := range files_to_copy_lst {

			//IMPORTANT!! - the Crawler_page_img has alread been uploaded to S3, so we dont need 
			//              to download it from S3 and reupload to gf_images S3 bucket. Instead we do 
			//              a file copy operation within the S3 system without downloading here.

			source_bucket_and_file__s3_path_str := filepath.Clean(fmt.Sprintf("/%s/%s",p_crawled_images_s3_bucket_name_str,s3_path_str))


			//DEBUGGING
			fmt.Println("")
			fmt.Println("==========================")
			fmt.Println(p_gf_images_s3_bucket_name_str)
			fmt.Println(source_bucket_and_file__s3_path_str)
			fmt.Println(s3_path_str)
			fmt.Println("")

			gf_err := gf_core.S3__copy_file(p_gf_images_s3_bucket_name_str, //p_target_bucket_name_str,
				source_bucket_and_file__s3_path_str, //p_source_file__s3_path_str
				s3_path_str,                         //p_target_file__s3_path_str
				p_runtime.S3_info,
				p_runtime_sys)
			if gf_err != nil {
				return gf_err
			}
		}
	}
	//--------------------------
	
	return nil
}