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

use image;

use image::{GenericImageView};
use crate::gf_image;
use gf_core;

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn native__open_image(p_img_source_file_path_str: &str) -> gf_image::GFimage {

    // DynamicImage - enumeration over all supported ImageBuffer<P> types
    let img: image::DynamicImage = image::open(p_img_source_file_path_str).unwrap();
    let (width_int, height_int)  = img.dimensions();
    
    let gf_img = gf_image::GFimage{
        width_int:  width_int,
        height_int: height_int,
        raw_img:    img,
    };

    return gf_img;
}

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn native__save_image(p_gf_img: &gf_image::GFimage,
    p_img_target_file_path_str: &str) {

    let img_data_lst = p_gf_img.raw_img.to_rgba().to_vec();
    let img_buffer   = image::ImageBuffer::from_vec(p_gf_img.width_int,
        p_gf_img.height_int,
        img_data_lst).unwrap();

    gf_core::gf_image::save_image_buff(img_buffer, &p_img_target_file_path_str);
}