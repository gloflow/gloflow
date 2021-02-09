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

#![allow(non_snake_case)]

use std::collections::HashMap;
use image::{GenericImageView};

//-------------------------------------------------
pub struct GFimageCollageConfig {
    pub output_img_file_path_str: String,
    pub width_int:                u64,
    pub height_int:               u64,
    pub rows_num_int:             u32,
    pub columns_num_int:          u32,
}

//-------------------------------------------------
pub fn create(p_input_imgs_files_paths_lst: Vec<String>,
    p_imgs_collage_config: &GFimageCollageConfig) {

    let mut row_int    = 0;
    let mut column_int = 0;

    // COLLAGE_IMG_BUFFER
    let mut collage_img_buff = image::ImageBuffer::new(p_imgs_collage_config.width_int as u32,
        p_imgs_collage_config.height_int as u32);

    for img_file_path_str in p_input_imgs_files_paths_lst.iter() {

        //---------------------
        // OPEN_FILE
        println!("{}", img_file_path_str);
        
        let img:         image::DynamicImage                          = image::open(img_file_path_str).unwrap();
        let gf_img_buff: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = img.to_rgba();

        // ADD_IMAGE_FROM_BUFFER
        let (new_row_int, new_column_int, full_bool) = add_img_from_buffer(&gf_img_buff,
            &mut collage_img_buff,
            row_int,
            column_int,
            p_imgs_collage_config);
        
        if full_bool {
            break;
        }
        
        row_int    = new_row_int;
        column_int = new_column_int;
        //---------------------
    }

    collage_img_buff.save(&p_imgs_collage_config.output_img_file_path_str).unwrap();
}

//-------------------------------------------------
// ADD_IMG_FROM_BUFFER
pub fn add_img_from_buffer(p_img_buff: &image::ImageBuffer<image::Rgba<u8>, Vec<u8>>,
    p_collage_img_buff:    &mut image::ImageBuffer<image::Rgba<u8>, Vec<u8>>,
    p_row_int:             u32,
    p_column_int:          u32,
    p_imgs_collage_config: &GFimageCollageConfig) -> (u32, u32, bool) {

    // CELL_DIMENSIONS
    let cell_width_int  = p_imgs_collage_config.width_int as u32 / p_imgs_collage_config.columns_num_int;
    let cell_height_int = p_imgs_collage_config.height_int as u32 / p_imgs_collage_config.rows_num_int;

    //---------------------
    // NEW_DIMENSIONS - fit an image so that the smaller dimension (width or height) is assigned the 
    //                  dimension of the collage cell, and the larger dimension (width/height) is then scaled
    //                  in proportion to maintain the original aspect ration.
    //                  proportion:
    //                  if img_width > img_height -> img_width / img_height = x / cell_height
    //                  if img_width < img_height -> img_width / img_height = cell_width / x
    let img_width_int  = p_img_buff.width();
    let img_height_int = p_img_buff.height();
    
    // get scaled down dimensions of an image. aspect ratio is preserved, and a certain heuristic
    // is used to fit image dimensions into target cell dimensions.
    let (new_width_f, new_height_f) = get_img_new_dimensions(img_width_int,
        img_height_int,
        cell_width_int,
        cell_height_int);

    //---------------------
    // RESIZE
    // 
    // FILTER_TYPES:
    // Nearest 	    31 ms - worst quality (pixelization when downsampling is visible)
    // Triangle 	414 ms
    // CatmullRom 	817 ms
    // Gaussian 	1180 ms
    // Lanczos3 	1170 ms

    let mut resized_img_buff = image::imageops::resize(p_img_buff,
        new_width_f as u32,
        new_height_f as u32,
        image::FilterType::Nearest);
    
    // CROP
    let crop_x_int = (new_width_f - cell_width_int as f32) / 2.0;
    let crop_y_int = (new_height_f - cell_height_int as f32) / 2.0;
    let croped_img_buff = image::imageops::crop(&mut resized_img_buff,
        crop_x_int as u32,
        crop_y_int as u32,
        cell_width_int,
        cell_height_int);

    // DynamicImage - resize_to_fill() - was used previously before doing resizing/croping directly
    //---------------------
    // COPY_PIXELS

    let collage_window__x_int = p_column_int * cell_width_int;
    let collage_window__y_int = p_row_int * cell_height_int;

    for x in 0..cell_width_int as u32 {
        let collage__global_x_int = collage_window__x_int + x;

        for y in 0..cell_height_int as u32 {
            
            let collage__global_y_int = collage_window__y_int + y;
            let px                    = croped_img_buff.get_pixel(x, y);

            p_collage_img_buff.put_pixel(collage__global_x_int, collage__global_y_int, px);
        }
    }

    //---------------------
    // right edge of the image has been reached
    if p_column_int == (p_imgs_collage_config.columns_num_int - 1) {
        
        let new_column_int = 0; // move back to the left of the image

        // last row has been completed, dont draw any more images, no more cells left.
        if p_row_int == p_imgs_collage_config.rows_num_int - 1 {
            return (p_row_int, new_column_int, true);
        } else {
            let new_row_int = p_row_int + 1; // move one row down
            return (new_row_int, new_column_int, false);
        }
    } 
    // still room left to move to the right
    else {
        let new_column_int = p_column_int + 1; // move one column to the right
        return (p_row_int, new_column_int, false);
    }
}

//-------------------------------------------------
fn get_img_new_dimensions(p_img_width_int: u32,
    p_img_height_int:  u32,
    p_cell_width_int:  u32,
    p_cell_height_int: u32) -> (f32, f32) {

    // DIMENSIONS_RATIO - ratio of width to height of both the image and target cell.
    let img_dim_ratio_f  = p_img_width_int as f32 / p_img_height_int as f32;
    let cell_dim_ratio_f = p_cell_width_int as f32 / p_cell_height_int as f32;
    
    // println!("img_ratio  - {}", img_dim_ratio_f);
    // println!("cell_ratio - {}", cell_dim_ratio_f);

    // SCALE - get the scale value used to scale the image dimensions down to fit the cell dimensions.
    //         dimension ratios of the image and target cell are compared.
    //         - if img_dim_ratio ratio is higher then cell_dim_ratio then the ratio of image/cell heights
    //           is used as the scaling factor for fitting the image into a cell.
    //         - if img_dim_ratio ratio is lower then cell_dim_ratio then the ratio of image/cell widths is used.
    let scale_f = if img_dim_ratio_f > cell_dim_ratio_f {
        p_cell_height_int as f32 / p_img_height_int as f32
    } else {
        p_cell_width_int as f32 / p_img_width_int as f32
    };

    let new_width_f  = p_img_width_int as f32 * scale_f;
    let new_height_f = p_img_height_int as f32 * scale_f;


    // println!("scale_f     - {}", scale_f);
    // println!("img_width   - {}", p_img_width_int);
    // println!("img_height  - {}", p_img_height_int);
    // println!("new width   - {}", new_width_f);
    // println!("new height  - {}", new_height_f);
    // println!("cell width  - {}", p_cell_width_int);
    // println!("cell height - {}", p_cell_height_int);


    return (new_width_f, new_height_f);
}

//-------------------------------------------------
// DRAW_BORDERS
pub fn draw_borders(p_collage_img_buff: &mut image::ImageBuffer<image::Rgba<u8>, Vec<u8>>, 
    p_rows_num_int:                   u32,
    p_columns_num_int:                u32,
    p_img_index_to_collage_coord_map: HashMap<u32, (u32, u32)>) {

    let cells_num_int   = p_rows_num_int * p_columns_num_int;
    let cell_width_int  = p_collage_img_buff.width() / p_columns_num_int;
    let cell_height_int = p_collage_img_buff.height() / p_rows_num_int;

    for i in 0..cells_num_int {

        let (img_row_int, img_column_int) = p_img_index_to_collage_coord_map.get(&i).unwrap();

        // HORIZONTAL_BORDERS (TOP/BOTTM)
        let global_y_top_int    = img_row_int * cell_height_int + 0 as u32;
        let global_y_bottom_int = img_row_int * cell_height_int + cell_height_int-1 as u32;
        for x in 0..cell_width_int {
            let global_x_int = img_column_int * cell_width_int + x as u32;

            // TOP
            let pixel = p_collage_img_buff.get_pixel_mut(global_x_int, global_y_top_int);
            *pixel = image::Rgba([0, 0, 0, 255 as u8]);

            // BOTTOM
            let pixel = p_collage_img_buff.get_pixel_mut(global_x_int, global_y_bottom_int);
            *pixel = image::Rgba([0, 0, 0, 255 as u8]);
        }

        // VERTICAL_BORDRS (LEFT/RIGHT)
        let global_x_left_int  = img_column_int * cell_width_int + 0 as u32;
        let global_x_right_int = img_column_int * cell_width_int + cell_width_int-1 as u32;
        for y in 0..cell_height_int {
            let global_y_int = img_row_int * cell_height_int + y as u32;

            // LEFT
            let pixel = p_collage_img_buff.get_pixel_mut(global_x_left_int, global_y_int);
            *pixel = image::Rgba([0, 0, 0, 255 as u8]);

            // RIGHT
            let pixel = p_collage_img_buff.get_pixel_mut(global_x_right_int, global_y_int);
            *pixel = image::Rgba([0, 0, 0, 255 as u8]);
        }
    }
}