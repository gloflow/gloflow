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

use image::{GenericImageView, GenericImage};
use palette::{Srgb, Srgba, Lab, Lch, Pixel, Saturate};

use crate::gf_image;

//-------------------------------------------------
/*LAB_COLOR_SPACE - https://docs.rs/palette/0.5.0/palette/struct.Lab.html

CIE L*a*b* - device independent color space.
    includes all perceivable colors.
    used to convert between color spaces, because of ability to represent all their colors,
    and in color manipulation, because of its perceptual uniformity (perceptual difference between 
    two colors is equal to their numerical difference).*/

//-------------------------------------------------
pub fn saturate_selective(p_gf_img: &mut gf_image::GFimage,
    p_color_ref:          &gf_image::GFcolorRGB,
    p_saturation_level_f: f32) {


    // LAB_COLOR
    // Srgb::new() - takes in RGB components as values normalized in 0.0-1.0 range.
    let color_ref_lab: Lab = Srgb::new(p_color_ref.r as f32 / 255.0,
        p_color_ref.g as f32 / 255.0,
        p_color_ref.b as f32 / 255.0).into();

    for x in 0..p_gf_img.width_int {
        for y in 0..p_gf_img.height_int {


            let px = p_gf_img.raw_img.get_pixel(x, y);



            let px_lab: Lab = Srgb::new(px.data[0] as f32 / 255.0,
                px.data[1] as f32 / 255.0,
                px.data[2] as f32 / 255.0).into();



            let distance_f = distance(color_ref_lab, px_lab);
            if distance_f < 40.0 {


                // into_lenear() - linear color space is needed to be able to apply
                //                 arithmetic operations to color values.
                //                 (color gamma of 1.0).
                //                 in nature light behaves linearly.
                // https://matt77hias.github.io/blog/2018/07/01/linear-gamma-and-sRGB-color-spaces.html
                let lch_color: Lch = Srgb::from_raw(&px.data)
                    .into_format()
                    .into_linear()
                    .into();


                // SATURATE
                let new_lch_color = lch_color.saturate(p_saturation_level_f);


                let new_rgba_color = Srgba::from_linear(new_lch_color.into()).into_format().into_raw();
                let new_px = image::Rgba{
                    data: new_rgba_color
                };
                
                p_gf_img.raw_img.put_pixel(x, y, new_px);       
            }
        }
    }
}

//-------------------------------------------------
// CIE76 - 1976 formula that related a measured color difference to a known set of CIELAB coordinates.
//         standard Euclidian distance in LAB space.
// FIX!! - implement the 1994 and 2000 formulas for dinstance.
pub fn distance(p_color_1: Lab, p_color_2: Lab) -> f32 {

    let l_delta_sq_f   = (p_color_2.l - p_color_1.l).powf(2.0);
    let a_delta_sq_f   = (p_color_2.a - p_color_1.a).powf(2.0);
    let b_delta_sq_f   = (p_color_2.b - p_color_1.b).powf(2.0);
    let deltas_sum_f   = l_delta_sq_f + a_delta_sq_f + b_delta_sq_f;
    let lab_distance_f = deltas_sum_f.sqrt();

    return lab_distance_f;
}