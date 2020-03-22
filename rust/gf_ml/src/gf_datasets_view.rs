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

use std::collections::HashMap;
use tensorflow;
use gf_core;

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn view_tf_records(p_tfrecords_file_path_str: &str,
    p_img_target_file_path_str:   &str,
    p_tf_example__img_width_int:  u64,
    p_tf_example__img_height_int: u64,
    p_collage__img_width_int:     u64,
    p_collage__img_height_int:    u64,
    p_collage__rows_num_int:      u32,
    p_collage__columns_num_int:   u32) {

    let mut tf_records_reader         = gf_core::gf_tf::get_tf_records__reader(p_tfrecords_file_path_str);
    let mut tf_example_raw_buffer_lst = [0u8; 3000]; // buffer for individual examples read in from a .tfrecords file
    
    // IMAGE_COLLAGE
    let imgs_collage_config = gf_core::gf_image_collage::GFimageCollageConfig {
        output_img_file_path_str: (*p_img_target_file_path_str).to_string(),
        width_int:       p_collage__img_width_int,  
        height_int:      p_collage__img_height_int,
        rows_num_int:    p_collage__rows_num_int,
        columns_num_int: p_collage__columns_num_int
    };

    let mut collage_img_buff = image::ImageBuffer::new(imgs_collage_config.width_int as u32, imgs_collage_config.height_int as u32);
    let mut row_int    = 0;
    let mut column_int = 0;
    let mut img_index_to_collage_coord_map = HashMap::new();

    let mut i = 0;
    loop {

        let next = tf_records_reader.read_next(&mut tf_example_raw_buffer_lst);

        match next {

            Ok(resp) => match resp {
                Some(len) => { 
                    println!("data received - {} bytes", len);

                    let data_lst = &tf_example_raw_buffer_lst[0..len];
                    
                    // READ_TF_EXAMPLE
                    let (gf_img_buff, gf_img_label_int) = gf_core::gf_tf::read_tf_example__to_img_buffer(&data_lst,
                        p_tf_example__img_width_int,
                        p_tf_example__img_height_int);

                    // IMAGE_COLLAGE
                    let (new_row_int, new_column_int, full_bool) = gf_core::gf_image_collage::add_img_from_buffer(&gf_img_buff,
                        &mut collage_img_buff,
                        row_int,
                        column_int,
                        &imgs_collage_config);
                    
                    // memories the coordinate in the collage of this image. this is potentially needed later
                    // to query where the image was placed in a 2D matrix of the collage.
                    img_index_to_collage_coord_map.insert(i as u32, (row_int, column_int));
                    
                    if full_bool {
                        break;
                    }

                    row_int    = new_row_int;
                    column_int = new_column_int;
                    
                    i += 1;
                },
                None => break,
            }, 

            Err(tensorflow::io::RecordReadError::CorruptFile) | Err(tensorflow::io::RecordReadError::IoError { .. }) => {
                break;
            }
            _ => {}
        }
    }

    // DRAW_BORDERS
    gf_core::gf_image_collage::draw_borders(&mut collage_img_buff,
        imgs_collage_config.rows_num_int,
        imgs_collage_config.columns_num_int,
        img_index_to_collage_coord_map);

    // SAVE_FILE
    collage_img_buff.save(&imgs_collage_config.output_img_file_path_str).unwrap();
}