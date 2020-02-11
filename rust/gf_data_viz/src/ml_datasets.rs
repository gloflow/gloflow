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

extern crate cairo;
extern crate rand;

use cairo::{ ImageSurface, Format, Context };
use std::fs;
use std::fs::File;

//-------------------------------------------------
struct GfConfig {
    image_width_int     : u32,
    image_height_int    : u32,
    target_dir_path_str : String,
}

struct GFruntimeGfx {
    canvas: ImageSurface,
    ctx:    Context,
}

//-------------------------------------------------
pub fn generate(p_dataset_name_str : String,
    p_image_width_int     : u32,
    p_image_height_int    : u32,
    p_target_dir_path_str : String) {

    let gf_runtime_gfx = runtime_get_graphics(p_image_width_int,
        p_image_height_int);

    let gf_config = GfConfig{
        image_width_int:     p_image_width_int,
        image_height_int:    p_image_height_int,
        target_dir_path_str: p_target_dir_path_str,
    };

    // TRAIN_DATASET
    generate_of_type(&p_dataset_name_str,
        &"train",
        &gf_config,
        &gf_runtime_gfx);

    // VALIDATION_DATASET
    generate_of_type(&p_dataset_name_str,
        &"validation",
        &gf_config,
        &gf_runtime_gfx);
}

//-------------------------------------------------
fn generate_of_type(p_dataset_name_str :&str,
    p_type_str       : &str,
    p_gf_config      : &GfConfig,
    p_gf_runtime_gfx : &GFruntimeGfx) {


    let objs_number_int = 100;
    
    for i in 0..objs_number_int {

        draw_rect(p_gf_config, p_gf_runtime_gfx)
    }
}


//-------------------------------------------------
fn draw_rect(p_gf_config : &GfConfig,
    p_gf_runtime_gfx : &GFruntimeGfx) {

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
fn image_save_to_file(p_target_file_name_str : String,
    p_target_dir_path_str : String,
    p_surface             : &ImageSurface) {



    let filename_str = format!("{}/{}",
        p_target_dir_path_str,
        p_target_file_name_str);

    let mut file = File::create(filename_str)
        .expect("failed to create a file to FS");

    p_surface.write_to_png(&mut file)
        .expect("failed to save a PNG image to a file in FS");
}

//-------------------------------------------------
fn runtime_get_graphics(p_image_width_int : u32,
    p_image_height_int : u32) -> GFruntimeGfx {

    let surface = ImageSurface::create(Format::ARgb32,
        p_image_width_int as i32,
        p_image_height_int as i32)
        .expect("Cairo failed to create a drawing surface");

    let ctx = Context::new(&surface);

    let gf_runtime_gfx = GFruntimeGfx{
        canvas: surface,
        ctx:    ctx,
    };

    gf_runtime_gfx
}