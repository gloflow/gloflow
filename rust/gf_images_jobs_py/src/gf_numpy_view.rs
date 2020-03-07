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
use image;
use numpy::{PyArray2, PyArray3, PyArray4};

use gf_images_jobs;

//-------------------------------------------------
// VIEW_NUMPY_ARR_2D
#[allow(non_snake_case)]
pub fn arr_2D(p_numpy_2d_lst: &PyArray2<f64>,
    p_img_target_file_path_str: String) {

    let numpy_shape_lst = p_numpy_2d_lst.shape();
    let mut img_buff = image::ImageBuffer::new(numpy_shape_lst[0] as u32, numpy_shape_lst[1] as u32);

    for ((x, y), val_f) in p_numpy_2d_lst.as_array_mut().indexed_iter_mut() {
        let pixel = img_buff.get_pixel_mut(x as u32, y as u32);

        *pixel = image::Rgba([
            (*val_f) as u8,
            (*val_f) as u8,
            (*val_f) as u8,
            0 as u8]);
    }
    img_buff.save(&p_img_target_file_path_str).unwrap();
}

//-------------------------------------------------
// VIEW_NUMPY_ARR_3D
#[allow(non_snake_case)]
pub fn arr_3D(p_numpy_3d_lst: &PyArray3<f64>,
    p_img_target_file_path_str: String) {

    let numpy_shape_lst = p_numpy_3d_lst.shape();
    assert!(numpy_shape_lst[2] == 3); // on 3rd axis the shape is always 3 (RGB)

    let mut img_buff = image::ImageBuffer::new(numpy_shape_lst[0] as u32, numpy_shape_lst[1] as u32);
    let arr          = p_numpy_3d_lst.as_array_mut();

    for y in 0..numpy_shape_lst[0] {

        // Axis(0) - gives us the first dimension of the 3D NumPy array (rows).
        let row_2d = arr.subview(ndarray::Axis(0), y);

        for x in 0..numpy_shape_lst[1] {
            
            // Axis(0) - gives us the first dimension of the 2D NumPy sub-array
            //           which here represents the individual pixel in a particular column
            //           (which itself is a 1D array of RGB values).
            let col   = row_2d.subview(ndarray::Axis(0), x);
            let pixel = img_buff.get_pixel_mut(x as u32, y as u32);
            *pixel    = image::Rgba([
                (col[0]) as u8,
                (col[1]) as u8,
                (col[2]) as u8,
                1 as u8]);
        }
    }

    img_buff.save(&p_img_target_file_path_str).unwrap();
}

//-------------------------------------------------
// VIEW_NUMPY_ARR_4D
#[allow(non_snake_case)]
pub fn arr_4D(p_numpy_4d_lst: &PyArray4<f64>,
    p_img_target_file_path_str: String,
    p_width_int:       u32,
    p_height_int:      u32,
    p_rows_num_int:    u32,
    p_columns_num_int: u32) {

    

    let numpy_shape_lst = p_numpy_4d_lst.shape();
    assert!(numpy_shape_lst[3] == 3); // on 3rd axis the shape is always 3 (RGB)



    let imgs_collage_config = gf_images_jobs::gf_image_collage::GFimageCollageConfig {
        output_img_file_path_str: p_img_target_file_path_str,
        width_int:       p_width_int,
        height_int:      p_height_int,
        rows_num_int:    p_rows_num_int,
        columns_num_int: p_columns_num_int,
    };



    let mut collage_img_buff = image::ImageBuffer::new(p_width_int, p_height_int);
    let arr                  = p_numpy_4d_lst.as_array_mut();

    let mut row_int    = 0;
    let mut column_int = 0;

    let mut img_index_to_collage_coord_map = HashMap::new();

    // multiple 3D numpy arrays are packed in sequence in an array
    for i in 0..numpy_shape_lst[0] {
        


        let image_3d = arr.subview(ndarray::Axis(0), i);
        println!("image {}", i);

        let img_width_int  = numpy_shape_lst[2]; // columns are in rows, so index is 2
        let img_height_int = numpy_shape_lst[1]; // images are packed by row, so index is 1

        // new image buffer for each 3d array (3d array is a single RGBA 2D image)
        let mut img_buff = image::ImageBuffer::new(img_width_int as u32, img_height_int as u32);
        
        
        // image is composed of an array of image row's (of pixels)
        for y in 0..img_height_int {

            // Axis(0) - gives us the first dimension of the image 3D NumPy array (rows).
            let row_2d = image_3d.subview(ndarray::Axis(0), y);
            
            // image row is composed of 1D pixels represented as arrays
            for x in 0..img_width_int {
                
                // Axis(0) - gives us the first dimension of the 2D NumPy sub-array
                //           which here represents the individual pixel in a particular column
                //           (which itself is a 1D array of RGB values).
                let col = row_2d.subview(ndarray::Axis(0), x);

                let pixel = img_buff.get_pixel_mut(x as u32, y as u32);

                
                *pixel = image::Rgba([
                    (col[0]) as u8,
                    (col[1]) as u8,
                    (col[2]) as u8,
                    255 as u8]);

            }
        }

        

        let (new_row_int, new_column_int, full_bool) = gf_images_jobs::gf_image_collage::add_img_from_buffer(&img_buff,
            &mut collage_img_buff,
            row_int,
            column_int,
            &imgs_collage_config);
        
        // memories the coordinate in the collage of this image. this is potentially needed later
        // to query where the image was placed in a 2D matrix of the collage.
        img_index_to_collage_coord_map.insert(i, (new_row_int, new_column_int));

        if full_bool {
            break;
        }

        row_int    = new_row_int;
        column_int = new_column_int;
    }





    let cell_width_int  = p_width_int / p_columns_num_int;
    let cell_height_int = p_height_int / p_rows_num_int;

    for i in 0..numpy_shape_lst[0] {

        let img_width_int  = numpy_shape_lst[2]; // columns are in rows, so index is 2
        let img_height_int = numpy_shape_lst[1]; // images are packed by row, so index is 1



        let (img_row_int, img_column_int) = img_index_to_collage_coord_map.get(&i).unwrap();
        println!("img [{}] row/column - {}/{}", i, img_row_int, img_column_int);

        for y in 0..img_height_int-1 {
            for x in 0..img_width_int-1 {

                // draw border on edges
                if (x == 0 || y == 0 || x == (img_width_int-1) || y == (img_height_int-1)) {

                    let global_x_int = img_column_int*cell_width_int + x as u32;
                    let global_y_int = img_row_int*cell_height_int + y as u32;
                    println!("img px [row/column {}/{}] global x/y - {}/{}", y, x, global_x_int, global_y_int);



                    let pixel = collage_img_buff.get_pixel_mut(global_x_int, global_y_int);
                    *pixel = image::Rgba([0, 0, 0, 255 as u8]);

                }
            }
        }

    }




    collage_img_buff.save(&imgs_collage_config.output_img_file_path_str).unwrap();
}