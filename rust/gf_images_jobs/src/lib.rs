/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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


extern crate libc;
use std::ffi::CStr; // https://doc.rust-lang.org/1.0.0/std/ffi/struct.CString.html

use std::str::FromStr;




pub mod gf_image_collage;

mod gf_image_color;
mod gf_image_generate;
mod gf_image_io;
mod gf_image_transform;
mod gf_image;
mod gf_tf;

//-------------------------------------------------
// C_API
//-------------------------------------------------
#[no_mangle]
#[allow(non_snake_case)]
pub extern "C" fn c__run_job(p_job_name: *const libc::c_char,
    p_img_source_file_path_str: *const libc::c_char,
    p_img_target_file_path_str: *const libc::c_char) {
    
    //---------------------
    // INPUT

    let buf_job_name = unsafe {CStr::from_ptr(p_job_name).to_bytes()};
    let job_name_str = String::from_utf8(buf_job_name.to_vec()).unwrap();
    
    let buf_img_source_file_path_str = unsafe {CStr::from_ptr(p_img_source_file_path_str).to_bytes()};
    let img_source_file_path_str     = String::from_utf8(buf_img_source_file_path_str.to_vec()).unwrap();

    let buf_img_target_file_path_str = unsafe {CStr::from_ptr(p_img_target_file_path_str).to_bytes()};
    let img_target_file_path_str     = String::from_utf8(buf_img_target_file_path_str.to_vec()).unwrap();
    

    println!("RUST ---------- running job - {}", job_name_str);
    println!("img_source_file_path_str - {}", img_source_file_path_str);
    println!("img_target_file_path_str - {}", img_target_file_path_str);

    //---------------------
}

//-------------------------------------------------
#[no_mangle]
#[allow(non_snake_case)]
pub extern "C" fn c__apply_transforms(p_transforms_c_lst: Vec<*const libc::c_char>,
    p_img_source_file_path_c_str: *const libc::c_char,
    p_img_target_file_path_c_str: *const libc::c_char) {

    //---------------------
    // INPUT

    let buf_img_source_file_path_str = unsafe {CStr::from_ptr(p_img_source_file_path_c_str).to_bytes()};
    let img_source_file_path_str     = String::from_utf8(buf_img_source_file_path_str.to_vec()).unwrap();

    let buf_img_target_file_path_str = unsafe {CStr::from_ptr(p_img_target_file_path_c_str).to_bytes()};
    let img_target_file_path_str     = String::from_utf8(buf_img_target_file_path_str.to_vec()).unwrap();


    // .into_iter() - consumes the collection so that on each iteration the exact data is provided.
    //                Once the collection has been consumed it is no longer available for reuse.
    //
    // RUST_NOTE - if a plain iterator is used here the compiler gives an error:
    //             expected type `*const i8`
    //             found type `&*const i8`
    let mut loaded_transforms_lst: Vec<String> = vec![];
    for transform_c_str in p_transforms_c_lst.into_iter() {
        
        let buf_transform_str = unsafe {CStr::from_ptr(transform_c_str).to_bytes()};
        let transform_str     = String::from_utf8(buf_transform_str.to_vec()).unwrap();

        loaded_transforms_lst.push(transform_str);
    }

    //---------------------

    // APPLY
    apply_transforms(loaded_transforms_lst,
        &img_source_file_path_str,
        &img_target_file_path_str)
}

//-------------------------------------------------
#[no_mangle]
#[allow(non_snake_case)]
pub extern "C" fn c__create_collage(p_input_imgs_files_paths_c_lst: Vec<*const libc::c_char>,
    p_output_img_file_path_c_str: *const libc::c_char) {

    //---------------------
    // INPUT

    // .into_iter() - consumes the collection so that on each iteration the exact data is provided.
    //                Once the collection has been consumed it is no longer available for reuse.
    //
    // RUST_NOTE - if a plain iterator is used here the compiler gives an error:
    //             expected type `*const i8`
    //             found type `&*const i8`
    let mut loaded_imgs_files_paths_lst: Vec<String> = vec![];
    for img_file_path_c_str in p_input_imgs_files_paths_c_lst.into_iter() {
        
        let buf_img_file_path_str = unsafe {CStr::from_ptr(img_file_path_c_str).to_bytes()};
        let img_file_path_str     = String::from_utf8(buf_img_file_path_str.to_vec()).unwrap();

        loaded_imgs_files_paths_lst.push(img_file_path_str);
    }

    let buf_output_img_file_path_c_str = unsafe {CStr::from_ptr(p_output_img_file_path_c_str).to_bytes()};
    let output_img_file_path_str       = String::from_utf8(buf_output_img_file_path_c_str.to_vec()).unwrap();

    //---------------------

    let imgs_collage_config = gf_image_collage::GFimageCollageConfig {
        output_img_file_path_str: output_img_file_path_str,
        width_int:                400,
        height_int:               400,
        rows_num_int:             5,
        columns_num_int:          5,
    };

    gf_image_collage::create(loaded_imgs_files_paths_lst,
        &imgs_collage_config);
}

//-------------------------------------------------
// RUST_API
//-------------------------------------------------
#[allow(non_snake_case)]
pub fn apply_transforms(p_transformations_lst: Vec<String>,
    p_img_source_file_path_str: &str,
    p_img_target_file_path_str: &str) {

    /*contrast // X
    brightness
    saturation // X
    sharpen
    sepia
    noise      // X
    hue*/

    println!("apply transform -------++++++++++++++++");
    let mut gf_img = gf_image_io::native__open_image(p_img_source_file_path_str);

    for trans_str in p_transformations_lst.iter() {
        let trans_components_lst: Vec<&str> = trans_str.split(":").collect();
        let trans_name_str                  = trans_components_lst[0];

        match trans_name_str {

            "saturate" => {
                let saturate_level_str = trans_components_lst[1];
                let saturate_level_f   = f32::from_str(saturate_level_str).unwrap();


                let color_ref = gf_image::GFcolorRGB{
                    r: 140, g: 205, b: 211,
                };

                gf_image_transform::saturate(&mut gf_img,
                    &color_ref,
                    saturate_level_f);
            },
            "contrast" => {
                let contrast_level_str = trans_components_lst[1];
                let contrast_level_f   = f32::from_str(contrast_level_str).unwrap();
                gf_image_transform::contrast(&mut gf_img, contrast_level_f);
            },
            "noise" => {
                gf_image_transform::noise(&mut gf_img);
            },
            _ => {
                println!("ERROR!! - unknown transformation - {}", trans_name_str)
            }
        }
    }

    gf_image_io::native__save_image(&gf_img, p_img_target_file_path_str);
}

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn create_collage(p_input_imgs_files_paths_lst: Vec<String>,
    p_output_img_file_path_c_str: String,
    p_width_int:                  u32,
    p_height_int:                 u32,
    p_rows_num_int:               u32,
    p_columns_num_int:            u32) {

    let imgs_collage_config = gf_image_collage::GFimageCollageConfig {
        output_img_file_path_str: p_output_img_file_path_c_str,
        width_int:       p_width_int,
        height_int:      p_height_int,
        rows_num_int:    p_rows_num_int,
        columns_num_int: p_columns_num_int,
    };

    gf_image_collage::create(p_input_imgs_files_paths_lst,
        &imgs_collage_config);
}

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn generate_ml_dataset_to_tfrecords(p_dataset_name_str: String,
    p_img_width_int:  u32,
    p_img_height_int: u32,
    p_target_dir_path_str: String) {




    gf_image_generate::ml_dataset_to_tfrecords(p_dataset_name_str,
        p_img_width_int,
        p_img_height_int,
        p_target_dir_path_str);





}



//-------------------------------------------------
/*#[allow(non_snake_case)]
pub fn add_img_from_buffer_to_collage(p_img_buff: &image::ImageBuffer<image::Rgba<u8>, Vec<u8>>,
    p_collage_img_buff:    &mut image::ImageBuffer<image::Rgba<u8>, Vec<u8>>,
    p_row_int:             u32,
    p_column_int:          u32,
    p_imgs_collage_config: &gf_image_collage::GFimageCollageConfig) -> (u32, u32, bool) {

    let (new_row_int, new_column_int, full_bool) = gf_image_collage::add_img_from_buffer(p_img_buff,
        p_collage_img_buff,
        p_row_int,
        p_column_int,
        p_imgs_collage_config);

    return (new_row_int, new_column_int, full_bool);
}*/