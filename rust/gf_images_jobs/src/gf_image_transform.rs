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

use crate::gf_image_color;
use crate::gf_image;

use image::{GenericImageView, GenericImage, Pixel};
use rand::Rng;

//-------------------------------------------------
// TRANSFORMATION__SATURATION
pub fn saturate(p_gf_img: &mut gf_image::GFimage,
    p_color_ref:          &gf_image::GFcolorRGB,
    p_saturation_level_f: f32) {




    gf_image_color::saturate_selective(p_gf_img,
        p_color_ref,
        p_saturation_level_f);
}

//-------------------------------------------------
// TRANSFORMATION__BRIGHTNESS
// pub fn brightness(p_gf_img: &mut gf_image::GFimage) {
//
// }

//-------------------------------------------------
// TRANSFORMATION__CONTRAST
// http://thecryptmag.com/Online/56/imgproc_5.html
// https://math.stackexchange.com/questions/906240/algorithms-to-increase-or-decrease-the-contrast-of-an-image

// p_contrast_level_f - input contrast_level used by the caller of this function
//                      to modify the contrast in an image.
pub fn contrast(p_gf_img: &mut gf_image::GFimage,
    p_contrast_level_f: f32) {
    
    // input contrast_level has to be in the range -255 to +255
    let contrast_level_clamped_f = num::clamp(p_contrast_level_f, -255.0, 255.0);

    // contrast correction factor
    // fac - contrast correction factor
    // c   - desired level of contrast
    // fac = (259 * (c + 255)) / (255 * (259 - c))
    let contrast_correction_factor_f = (259.0 * (contrast_level_clamped_f + 255.0)) / (255.0 * (259.0 - contrast_level_clamped_f));

    //---------------------
    // LOOKUP_TABLE
    // IMPORTANT!! - 1D precomputed lookup table of pixel values with their contrast adjusted
    //               (contrast_correction_factor_f) applied. these are values for any color component (R/G/B) value.
    let mut new_px_color_component__lookup_lst = vec![0; 256];

    for px_color_component_val_int in 0..256 {

        // basic contrast/brightness linear transformation:
        // f(x) = fac * x + b
        //
        // fac - contrast control factor
        // x   - color component value (R, G, or B)
        // b   -  brightness
        // f(x) = fac * (x - 128) + 128 + b
        let new_px_channel_val_f         = contrast_correction_factor_f * (px_color_component_val_int as f32 - 128.0) + 128.0;
        let new_px_channel_val_clamped_f = num::clamp(new_px_channel_val_f, 0.0, 255.0);

        new_px_color_component__lookup_lst[px_color_component_val_int] = new_px_channel_val_clamped_f as u8;
    }
    
    //---------------------

    for x in 0..p_gf_img.width_int {
        for y in 0..p_gf_img.height_int {

            let mut px = p_gf_img.raw_img.get_pixel(x, y);
            let px_r_int = px.data[0];
            let px_g_int = px.data[1];
            let px_b_int = px.data[2];

            // lookup pre-computed color values with contrast_correction applied to them
            px.data[0] = new_px_color_component__lookup_lst[px_r_int as usize];
            px.data[1] = new_px_color_component__lookup_lst[px_g_int as usize];
            px.data[2] = new_px_color_component__lookup_lst[px_b_int as usize];

            p_gf_img.raw_img.put_pixel(x, y, px);
        }
    }
}

//-------------------------------------------------
// TRANSFORMATION__NOISE
pub fn noise(p_gf_img: &mut gf_image::GFimage) {

    let mut rng = rand::thread_rng();


    
    for x in 0..p_gf_img.width_int {
        for y in 0..p_gf_img.height_int {

            let random_px_increment_int = rng.gen_range(0, 127);
            let px                      = p_gf_img.raw_img.get_pixel(x, y);

            // img.get_pixel(x, y) - returns a 3 element array (r,g,b)
            let new_px = px.map(|p_px_ch_int| 
                if p_px_ch_int <= 255 - random_px_increment_int {
                    p_px_ch_int + random_px_increment_int
                } else {
                    255
                });

            p_gf_img.raw_img.put_pixel(x, y, new_px);
        }
    }
}