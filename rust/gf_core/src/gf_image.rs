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

use std::fs::File;

use image;
use cairo;

//-------------------------------------------------
// SAVE_IMAGE_BUFF
#[allow(non_snake_case)]
pub fn save_image_buff(p_img_buff: image::ImageBuffer<image::Rgba<u8>, Vec<u8>>,
    p_img_target_file_path_str: &str) {

    let img: image::DynamicImage = image::ImageRgba8(p_img_buff);
    img.save(p_img_target_file_path_str).unwrap();
}

//-------------------------------------------------
// SAVE_CAIRO
#[allow(non_snake_case)]
pub fn save_cairo(p_surface: &cairo::ImageSurface,
    p_img_target_file_path_str: &str) {

    let mut file = File::create(p_img_target_file_path_str)
        .expect("failed to create a file to FS");

    p_surface.write_to_png(&mut file)
        .expect("failed to save a PNG image to a file in FS");
}