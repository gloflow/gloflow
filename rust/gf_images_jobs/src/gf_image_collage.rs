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

use image::{GenericImageView};

use crate::gf_image_io;

//-------------------------------------------------
pub struct GFimageCollageConfig {
    pub input_imgs_files_paths_lst: Vec<String>,
    pub output_img_file_path_str:   String,
    pub width_int:                  u32,
    pub height_int:                 u32,
    pub rows_num_int:               u32,
    pub columns_num_int:            u32,
}

//-------------------------------------------------
pub fn create(p_imgs_collage_config: &GFimageCollageConfig) {

    // CELL_DIMENSIONS
    let cell_width_int:  u32 = p_imgs_collage_config.width_int / p_imgs_collage_config.columns_num_int;
    let cell_height_int: u32 = p_imgs_collage_config.height_int / p_imgs_collage_config.rows_num_int;

    let mut collage_img = image::ImageBuffer::new(p_imgs_collage_config.width_int,
        p_imgs_collage_config.height_int);

    let mut row_int    = 0;
    let mut column_int = 0;

    for img_file_path_str in p_imgs_collage_config.input_imgs_files_paths_lst.iter() {

        //---------------------
        // OPEN_FILE
        let gf_img = gf_image_io::native__open_image(img_file_path_str);

        // RESIZE
        /*Nearest 	31 ms - worst quality (pixelization when downsampling is visible)
        Triangle 	414 ms
        CatmullRom 	817 ms
        Gaussian 	1180 ms
        Lanczos3 	1170 ms*/
        let resized_img = gf_img.raw_img.resize_to_fill(cell_width_int,
            cell_height_int,
            image::FilterType::Triangle);

        let collage_window__x_int = column_int * cell_width_int;
        let collage_window__y_int = row_int * cell_height_int;

        for x in 0..cell_width_int {
            let collage__global_x_int = collage_window__x_int + x;

            for y in 0..cell_height_int {
                
                let collage__global_y_int = collage_window__y_int + y;
                let mut px                = resized_img.get_pixel(x, y);

                collage_img.put_pixel(collage__global_x_int, collage__global_y_int, px);
            }
        }
        
        //---------------------

        // right edge of the image has been reached
        if column_int == (p_imgs_collage_config.columns_num_int - 1) {
            column_int = 0; // move back to the left of the image

            // last row has been completed, dont draw any more images, no more cells left.
            if row_int == p_imgs_collage_config.rows_num_int - 1 {
                break
            } else {
                row_int += 1; // move one row down
            }
        } else {
            column_int += 1; // move one column to the right
        }
    }

    collage_img.save(&p_imgs_collage_config.output_img_file_path_str).unwrap();
}