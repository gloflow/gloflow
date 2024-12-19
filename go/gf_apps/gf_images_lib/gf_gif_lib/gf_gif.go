/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_gif_lib

import (
	"fmt"
	"os"
	"context"
	"path/filepath"
	"io"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"crypto/sha256"
	"encoding/hex"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------

type GFgif struct {
	Id                         primitive.ObjectID `json:"-"                     bson:"_id,omitempty"`
	Id_str                     string        `json:"id_str"                     bson:"id_str"` 
	T_str                      string        `json:"-"                          bson:"t"` //"gif"
	Creation_unix_time_f       float64       `json:"creation_unix_time_f"       bson:"creation_unix_time_f"`
	Deleted_bool               bool          `json:"deleted_bool"               bson:"deleted_bool"`
	Deleted_unix_time_f        float64       `json:"deleted_unix_time_f"        bson:"deleted_unix_time_f"`
	
	//------------------
	// indicates if GIF data in here is valid:
	//   - urls are correct and work
	//   - gif dimensions are correct
	//   -  frames number is correct
	// this is a new field, and some old GIF's might be valid but not contain this field.
	// in scenarios where a gif is indicated as valid, but shows to have any of its data not correct
	// (or its link dont work) then this field will be set to false. 
	Valid_bool                 bool          `json:"valid_bool"                 bson:"valid_bool"`
	
	//------------------
	Title_str                  string        `json:"title_str"                  bson:"title_str"`
	Gf_url_str                 string        `json:"gf_url_str"                 bson:"gf_url_str"`
	Origin_url_str             string        `json:"origin_url_str"             bson:"origin_url_str"`         // external url from which the GIF came
	Origin_page_url_str        string        `json:"origin_page_url_str"        bson:"origin_page_url_str"`    // external url of the page from which the GIF came
	Origin_page_domain_str     string        `json:"origin_page_domain_str"     bson:"origin_page_domain_str"` // external domain of the page from which the GIF came
	Width_int                  int           `json:"width_int"                  bson:"width_int"`
	Height_int                 int           `json:"height_int"                 bson:"height_int"`
	Preview_frames_num_int     int           `json:"preview_frames_num_int"     bson:"preview_frames_num_int"`
	PreviewFramesS3urlsLst     []string      `json:"previewFramesS3urlsLst" bson:"previewFramesS3urlsLst"`
	Tags_lst                   []string      `json:"tags_lst"                   bson:"tags_lst"`
	Hash_str                   string        `json:"hash_str"                   bson:"hash_str"`
	GFimageIDstr               gf_images_core.GFimageID `json:"gf_image_id_str" bson:"gf_image_id_str"`
}

//--------------------------------------------------

func ProcessAndUpload(pImageIDstr gf_images_core.GFimageID,
	p_image_source_url_str                        string,
	p_image_origin_page_url_str                   string,
	p_gif_download_and_frames__local_dir_path_str string,
	p_image_client_type_str                       string, //what type of client is processing this gif
	p_flows_names_lst                             []string,
	p_create_new_db_img_bool                      bool,
	pMediaDomainStr                               string,
	pS3bucketNameStr                              string,
	pS3info                                       *gf_aws.GFs3Info,
	pCtx                                          context.Context,
	pRuntimeSys                                   *gf_core.RuntimeSys) (*GFgif, *gf_core.GFerror) {

	gif, local_image_file_path_str, gfErr := Process(pImageIDstr,
		p_image_source_url_str,
		p_image_origin_page_url_str,
		p_gif_download_and_frames__local_dir_path_str,
		p_image_client_type_str,
		p_flows_names_lst,
		p_create_new_db_img_bool,
		pMediaDomainStr,
		pS3bucketNameStr,
		pS3info,
		pCtx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}
	//-----------------------
	// SAVE_IMAGE TO FS (S3)

	img_title_str, gfErr := gf_images_core.GetImageTitleFromURL(p_image_source_url_str,pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	s3_target_file_path_str := fmt.Sprintf("gifs/%s.gif", img_title_str)
	s3_resp_str, gfErr      := gf_aws.S3putFile(local_image_file_path_str, //p_target_file__local_path_str string,
		s3_target_file_path_str,
		pS3bucketNameStr,
		pS3info,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	fmt.Println(s3_resp_str)

	//-----------------------
	
	return gif, nil
}

//--------------------------------------------------

func Process(pImageIDstr gf_images_core.GFimageID,   
	p_image_source_url_str                        string,   
	p_image_origin_page_url_str                   string,
	p_gif_download_and_frames__local_dir_path_str string,
	p_image_client_type_str                       string, //what type of client is processing this gif
	p_flows_names_lst                             []string,
	p_create_new_db_img_bool                      bool,
	pMediaDomainStr                            string,
	pS3bucketNameStr                          string,
	pS3info                                       *gf_aws.GFs3Info,
	pCtx                                          context.Context,
	pRuntimeSys                                   *gf_core.RuntimeSys) (*GFgif, string, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_gif.Process()")
	
	//-------------
	// FETCH
	local_image_file_path_str, _, gfErr := gf_images_core.FetcherGetExternImage(p_image_source_url_str,
		p_gif_download_and_frames__local_dir_path_str,
		false, // p_random_time_delay_bool
		pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	//-----------------------
	// IMPORTANT!! - save first N frames of the GIF, to be uploaded to S3, and 
	//               served in UI's as GIF preview animations. this is an 
	//               optimization to handle really large GIF's an in general all GIF's
	//               (to save on bandwidth and download the full GIF only when the 
	//               user explicitly wants to view the full version)

	frames_num_int, frames_s3_urls_lst, varGFerr, framesGFerrsLst := storePreviewFrames(local_image_file_path_str,
		p_gif_download_and_frames__local_dir_path_str,
		pMediaDomainStr,
		pS3bucketNameStr,
		pS3info,
		pRuntimeSys)
	if varGFerr != nil {
		return nil, "", varGFerr
	}


	for _, frameGFerr := range framesGFerrsLst {
		if frameGFerr != nil {

			// FIX!! - return all errors to the user, to know exactly which frames failed, 
			//         even though most likely all frames failed.
			return nil, "", frameGFerr
		}
	}

	//-----------------------
	// GIF_GET_DIMENSIONS
	img_width_int, img_height_int, gfErr := getDimensions(local_image_file_path_str, pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	//-----------------------
	// GIF_OBJ_CREATE
	gif, gfErr := dbMongoCreate(p_image_source_url_str,
		p_image_origin_page_url_str,
		img_width_int,
		img_height_int,
		frames_num_int,
		frames_s3_urls_lst,
		pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	//-----------------------
	// IMAGE_CREATE

	if p_create_new_db_img_bool {

		// IMAGE_ID
		var gf_image_id_str gf_images_core.GFimageID
		if pImageIDstr == "" {
			new_image_id_str, gfErr := gf_images_core.CreateIDfromURL(p_image_source_url_str, pRuntimeSys)
			if gfErr != nil {
				return nil, "", gfErr
			}
			gf_image_id_str = new_image_id_str
		} else {
			gf_image_id_str = pImageIDstr
		}

		// IMAGE_TITLE
		image_title_str, gfErr := gf_images_core.GetImageTitleFromURL(p_image_source_url_str,pRuntimeSys)
		if gfErr != nil {
			return nil, "", gfErr
		}

		gif_first_frame_str := gif.PreviewFramesS3urlsLst[0]

		//-----------------------
		// DEPRECATED!! - remove this, Image_new_info should be used only, and should be validated directly, 
		//                not via gf_images_core.Image__verify_image_info()

		gf_image_info_map := map[string]interface{}{
			"id_str":                         string(gf_image_id_str),
			"title_str":                      image_title_str,
			"image_client_type_str":          p_image_client_type_str,

			//--------------
			"flows_names_lst":                p_flows_names_lst,
			"origin_url_str":                 p_image_source_url_str, //*p_image_origin_url_str,
			"origin_page_url_str":            p_image_origin_page_url_str,
			"original_file_internal_uri_str": local_image_file_path_str,

			//--------------
			"format_str":                     "gif",
			"width_int":                      img_width_int,
			"height_int":                     img_height_int,

			//--------------
			"thumbnail_small_url_str":        gif_first_frame_str, //image_thumbs.Small_relative_url_str,
			"thumbnail_medium_url_str":       gif_first_frame_str, //image_thumbs.Medium_relative_url_str,
			"thumbnail_large_url_str":        gif_first_frame_str, //image_thumbs.Large_relative_url_str,

			//"dominant_color_hex_str":dominant_color_hex_str,
		}

		verified_image_info_map, gfErr := gf_images_core.VerifyImageInfo(gf_image_info_map, pRuntimeSys)
		if gfErr != nil {
			return nil, "", gfErr
		}
		//-----------------------
		verified_gf_image_id_str := gf_images_core.GFimageID(verified_image_info_map["id_str"].(string)) //type-casting, gf_images_core.Gf_image_id is a type (not function)
		gf_image_info := &gf_images_core.GFimageNewInfo{
			IDstr:                          verified_gf_image_id_str,                                           // image_id_str,
			Title_str:                      verified_image_info_map["title_str"].(string),                      // image_title_str,
			Flows_names_lst:                verified_image_info_map["flows_names_lst"].([]string),              // p_flows_names_lst,
			Image_client_type_str:          verified_image_info_map["image_client_type_str"].(string),          // p_image_client_type_str,
			Origin_url_str:                 verified_image_info_map["origin_url_str"].(string),                 // p_image_source_url_str,
			Origin_page_url_str:            verified_image_info_map["origin_page_url_str"].(string),            // p_image_origin_page_url_str,
			Original_file_internal_uri_str: verified_image_info_map["original_file_internal_uri_str"].(string), // image_local_file_path_str,
			ThumbnailSmallURLstr:           verified_image_info_map["thumbnail_small_url_str"].(string),        // gif_first_frame_str,
			ThumbnailMediumURLstr:          verified_image_info_map["thumbnail_medium_url_str"].(string),       // gif_first_frame_str,
			ThumbnailLargeURLstr:           verified_image_info_map["thumbnail_large_url_str"].(string),        // gif_first_frame_str,
			Format_str:                     verified_image_info_map["format_str"].(string),                     // "gif",
		}

		// IMPORTANT!! - creates a GF_Image struct and stores it in the DB.
		//               every GIF in the system has its GF_Gif DB struct and GF_Image DB struct.
		//               these two structs are related by origin_url

		_, gfErr = gf_images_core.ImageCreateNew(gf_image_info, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, "", gfErr
		}

		// link the new gf_image DB record to the gf_gif DB record
		dbMongoUpdateImageID(gif.Id_str, verified_gf_image_id_str, pRuntimeSys)
	}

	//-----------------------

	return gif, local_image_file_path_str, nil
}

//--------------------------------------------------

func storePreviewFrames(pLocalFilePathSrc string,
	p_frames_images_dir_path_str string,
	pMediaDomainStr           string, 
	pS3bucketNameStr         string,
	pS3info                      *gf_aws.GFs3Info,
	pRuntimeSys                  *gf_core.RuntimeSys) (int, []string, *gf_core.GFerror, []*gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_gif.storePreviewFrames()")

	max_num__of_preview_frames_int      := 10
	frames_images_file_paths_lst, gfErr := gifFramesSaveToFS(pLocalFilePathSrc, p_frames_images_dir_path_str, max_num__of_preview_frames_int, pRuntimeSys)
	if gfErr != nil {
		return 0, nil, gfErr, nil
	}

	fmt.Println("== - ==++++   frames_images_file_paths_lst - "+fmt.Sprint(frames_images_file_paths_lst))

	previewFramesNumInt := len(frames_images_file_paths_lst)

	// ADD!! - make thumbnails out of individual frames - to reduce/standardize their size
	//-----------------------
	// SAVE_IMAGES TO FS (S3)
	previewFramesS3urlsLst := []string{}
	gfErrorsLst            := make([]*gf_core.GFerror, len(frames_images_file_paths_lst))
	for i, frameImageFilePathStr := range frames_images_file_paths_lst {

		frameImageFileNameStr  := filepath.Base(frameImageFilePathStr)
		targetFilePathStr      := fmt.Sprintf("gifs/frames/%s", frameImageFileNameStr)
		targetFileLocalPathStr := frameImageFilePathStr

		// UPLOAD
		s3_response_str, gfErr := gf_aws.S3putFile(targetFileLocalPathStr,
			targetFilePathStr,
			pS3bucketNameStr,
			pS3info,
			pRuntimeSys)

		if gfErr != nil {
			pRuntimeSys.LogFun("ERROR","GIF FRAME S3_UPLOAD ERROR >>> "+fmt.Sprint(gfErr.Error))
			gfErrorsLst[i] = gfErr
		}

		fmt.Println(s3_response_str)

		//-----------------------

		imageURLstr := gf_images_core.ImageGetPublicURL(targetFilePathStr,
			pMediaDomainStr, // pS3bucketNameStr,
			pRuntimeSys)

		previewFramesS3urlsLst = append(previewFramesS3urlsLst, imageURLstr)
	}

	//-----------------------

	return previewFramesNumInt, previewFramesS3urlsLst, nil, gfErrorsLst
}

//--------------------------------------------------

func gifFramesSaveToFS(pLocalFilePathStr string,
	p_frames_images_dir_path_str string,
	p_frames_num_to_get_int      int,
	pRuntimeSys                  *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {

	cyan  := color.New(color.FgCyan).SprintFunc()
	black := color.New(color.FgBlack).Add(color.BgWhite).SprintFunc()

	pRuntimeSys.LogFun("INFO", "")
	pRuntimeSys.LogFun("INFO", cyan("       --- GIF")+" - "+cyan("GET_FRAMES"))
	pRuntimeSys.LogFun("INFO", black(pLocalFilePathStr))
	pRuntimeSys.LogFun("INFO", "")

	//---------------------
	// GIF_GET_DIMENSIONS
	img_width_int, img_height_int, gfErr := getDimensions(pLocalFilePathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//---------------------

	file, err := os.Open(pLocalFilePathStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("OS failed to open a GIF file to then save its frames as individual files",
			"file_open_error",
			map[string]interface{}{"local_file_path_str": pLocalFilePathStr,},
			err, "gf_gif_lib", pRuntimeSys)
		return nil, gfErr
	}

	//---------------------
	// IMPORTANT!! - gif.DecodeAll - can and will panic frequently,
	//                               because a lot of the GIF images on the internet are somewhat broken
	defer func() {
		if r := recover(); r != nil {
			_ = gf_core.ErrorCreate("Gif__frames__save_to_fs() has failed, a panic was caught, likely from gif.DecodeAll()",
				"panic_error",
				map[string]interface{}{"local_file_path_str": pLocalFilePathStr,},
				err, "gf_gif_lib", pRuntimeSys)
		}
	}()

	gif_image, gif_err := gif.DecodeAll(file)

	if gif_err != nil {
		gfErr := gf_core.ErrorCreate("gif.DecodeAll() failed to parse a gif in order to save its frames to FS",
			"gif_decoding_frames_error",
			map[string]interface{}{"local_file_path_str": pLocalFilePathStr,},
			gif_err, "gf_gif_lib", pRuntimeSys)
		return nil, gfErr
	}

	//---------------------

	overpaint_image := image.NewRGBA(image.Rect(0, 0, img_width_int, img_height_int))

	// draw first frame of the GIF to the canvas
	draw.Draw(overpaint_image,
		overpaint_image.Bounds(),
		gif_image.Image[0],
		image.ZP,
		draw.Src)

	source_file_name_str := filepath.Base(pLocalFilePathStr)
	new_files_names_lst  := []string{}

	// IMPORTANT!! - save GIF frames to .png files on local filesystem
	for i, frame_img := range gif_image.Image {

		//-------------------
		// IMPORTANT!! - if p_frames_num_to_get_int is 0, the caller wants all GIF frames, so no need 
		//               to check if the current GIF frame ("i") is larger then the max number of frames
		//               the user wants saves.
		//              
		// IMPORTANT!! - a GIF might have fewer frames then are asked for in p_frames_num_to_get_int

		if p_frames_num_to_get_int != 0 && i > p_frames_num_to_get_int {
			break //expected number of frames has been saved, so exit the loop
		}

		//-------------------

		draw.Draw(overpaint_image,
			overpaint_image.Bounds(),
			frame_img,
			image.ZP,
			draw.Over)

		//-------------------
		// save current frame
		new_file_name_str := fmt.Sprintf("%s/%s_%d.png", p_frames_images_dir_path_str, source_file_name_str, i)
		file, err         := os.Create(new_file_name_str)
		if err != nil {
			gfErr := gf_core.ErrorCreate("OS failed to create a file to save a GIF frame to FS",
				"file_create_error",
				map[string]interface{}{"new_file_name_str": new_file_name_str,},
				err, "gf_gif_lib", pRuntimeSys)
			return nil, gfErr
		}

		err = png.Encode(file, overpaint_image)
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to encode png image_byte array while saving GIF frame to FS",
				"png_encoding_error",
				map[string]interface{}{"new_file_name_str": new_file_name_str,},
				err, "gf_gif_lib", pRuntimeSys)
			return nil, gfErr
		}

		file.Close()

		//-------------------
		fmt.Sprint("++++++++  new_file_name_str - "+new_file_name_str)

		new_files_names_lst = append(new_files_names_lst, new_file_name_str)
	}

	return new_files_names_lst, nil
}

//--------------------------------------------------

func getDimensions(pLocalFilePathStr string,
	pRuntimeSys *gf_core.RuntimeSys) (int, int, *gf_core.GFerror) {

	file, err := os.Open(pLocalFilePathStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("OS failed to open a file to get image dimensions",
			"file_open_error",
			map[string]interface{}{"local_file_path_str": pLocalFilePathStr,},
			err, "gf_gif_lib", pRuntimeSys)
		return 0, 0, gfErr
	}

	//---------------------
	// IMPORTANT!! - gif.DecodeAll - can and will panic frequently,
	//                               because a lot of the GIF images on the internet are somewhat broken
	defer func() {
		if r := recover(); r != nil {
			_ = gf_core.ErrorCreate("getDimensions() has failed, a panic was caught, likely from gif.DecodeAll()",
				"panic_error",
				map[string]interface{}{"local_file_path_str": pLocalFilePathStr,},
				err, "gf_gif_lib", pRuntimeSys)
		}
	}()

	gif, gif_err := gif.DecodeAll(file)

	if gif_err != nil {
		gfErr := gf_core.ErrorCreate("gif.DecodeAll() failed to parse a gif in order to save its frames to FS",
			"gif_decoding_frames_error",
			map[string]interface{}{"local_file_path_str": pLocalFilePathStr,},
			gif_err, "gf_gif_lib", pRuntimeSys)
		return 0, 0, gfErr
	}

	//---------------------

	var lowestX  int
	var lowestY  int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY,nil
}

//--------------------------------------------------

func gifGetHash(p_image_local_file_path_str string,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	hash := sha256.New()

	f, err := os.Open(p_image_local_file_path_str)
	if err != nil {
		gfErr := gf_core.ErrorCreate("OS failed to open a GIF file to get its hash",
			"file_open_error",
			map[string]interface{}{"image_local_file_path_str": p_image_local_file_path_str,},
			err, "gf_gif_lib", pRuntimeSys)
		return "", gfErr
	}
	defer f.Close()

	io.Copy(hash,f)

	hash_str := hex.EncodeToString(hash.Sum(nil))
	return hash_str, nil
}