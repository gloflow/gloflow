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

use std::fs;
use std::fs::File;
use rand;
use cairo;

use gf_core;

//-------------------------------------------------
struct GFconfig {
    image_width_int:     u32,
    image_height_int:    u32,
    target_dir_path_str: String,
}

struct GFruntimeGfx {
    canvas: cairo::ImageSurface,
    ctx:    cairo::Context,
}

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn generate(p_dataset_name_str: String,
    p_image_width_int:     u32,
    p_image_height_int:    u32,
    p_target_dir_path_str: String) {

    let gf_runtime_gfx = runtime_get_graphics(p_image_width_int,
        p_image_height_int);

    let gf_config = GFconfig{
        image_width_int:     p_image_width_int,
        image_height_int:    p_image_height_int,
        target_dir_path_str: p_target_dir_path_str,
    };



    println!("==============>>>");

    // TRAIN_DATASET
    generate_for_env(&p_dataset_name_str,
        &"train",
        &gf_config,
        &gf_runtime_gfx);

    // VALIDATION_DATASET
    generate_for_env(&p_dataset_name_str,
        &"validation",
        &gf_config,
        &gf_runtime_gfx);
}

//-------------------------------------------------
#[allow(non_snake_case)]
fn generate_for_env(p_dataset_name_str: &str,
    p_env_str:        &str,
    p_gf_config:      &GFconfig,
    p_gf_runtime_gfx: &GFruntimeGfx) {

    


    println!(" generate for ENV - {}", p_env_str);


    let objs_number_int = 100;
    
    for i in 0..objs_number_int {

        println!("==============>>> 111111");
        draw_rect(p_gf_config, p_gf_runtime_gfx);


        // IMPORTANT!! - ".png" is relevant here, because currently gf_core Rust Crate
        //               Cargo.toml specifies "cairo-rs" dependency with feature for "png"
        //               enabled only. have not tested saving cairo Surfaces to .jpeg.
        let output_file_path_str = format!("{}/{}-{}.png",
            p_gf_config.target_dir_path_str,
            p_dataset_name_str,
            i);

        // SAVE_FILE
        gf_core::gf_image::save_cairo(&p_gf_runtime_gfx.canvas,
            &output_file_path_str);
    }
}

//-------------------------------------------------
// DRAW_RECT
#[allow(non_snake_case)]
fn draw_rect(p_gf_config: &GFconfig,
    p_gf_runtime_gfx: &GFruntimeGfx) {

    let gfx_ctx         = &p_gf_runtime_gfx.ctx;
    let randomize_bool  = true;
    if randomize_bool {
        gfx_ctx.set_source_rgb(rand::random::<f64>(), rand::random::<f64>(), rand::random::<f64>());
    } else {
        gfx_ctx.set_source_rgb(1.0, 1.0, 1.0);
    }
    gfx_ctx.paint();

    let w = rand::random::<f64>()*p_gf_config.image_width_int as f64;
    let h = rand::random::<f64>()*p_gf_config.image_height_int as f64;
    let x = rand::random::<f64>()*p_gf_config.image_width_int as f64 - w;
    let y = rand::random::<f64>()*p_gf_config.image_height_int as f64 - h;


    gfx_ctx.rectangle(x as f64, y as f64, w as f64, h as f64);
}

//-------------------------------------------------
#[allow(non_snake_case)]
fn runtime_get_graphics(p_image_width_int: u32,
    p_image_height_int: u32) -> GFruntimeGfx {

    let surface = cairo::ImageSurface::create(cairo::Format::ARgb32,
        p_image_width_int as i32,
        p_image_height_int as i32)
        .expect("failed to create a drawing surface with the Cairo backend");

    let ctx = cairo::Context::new(&surface);

    let gf_runtime_gfx = GFruntimeGfx{
        canvas: surface,
        ctx:    ctx,
    };

    gf_runtime_gfx
}