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
// use std::fs::File;
use rand;
use cairo;

use gf_core;

//-------------------------------------------------
struct GFdatasetConfig {
    name_str:            String,
    elements_num_int:    u64,
    image_width_int:     u64,
    image_height_int:    u64,
    target_dir_path_str: String,
}

struct GFruntimeGfx {
    canvas: cairo::ImageSurface,
    ctx:    cairo::Context,
}

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn generate(p_dataset_name_str: String,
    p_classes_lst:         Vec<String>,
    p_elements_num_int:    u64,
    p_image_width_int:     u64,
    p_image_height_int:    u64,
    p_target_dir_path_str: String) {

    let gf_runtime_gfx = runtime_get_graphics(p_image_width_int,
        p_image_height_int);

    let gf_dataset_config = GFdatasetConfig{
        name_str:            p_dataset_name_str,
        elements_num_int:    p_elements_num_int,
        image_width_int:     p_image_width_int,
        image_height_int:    p_image_height_int,
        target_dir_path_str: p_target_dir_path_str,
    };

    // TRAIN_DATASET
    generate_for_env(&"train",
        &p_classes_lst,
        &gf_dataset_config,
        &gf_runtime_gfx);

    // VALIDATION_DATASET
    generate_for_env(&"validation",
        &p_classes_lst,
        &gf_dataset_config,
        &gf_runtime_gfx);
}

//-------------------------------------------------
#[allow(non_snake_case)]
fn generate_for_env(p_env_str: &str,
    p_classes_lst:       &Vec<String>,
    p_gf_dataset_config: &GFdatasetConfig,
    p_gf_runtime_gfx:    &GFruntimeGfx) {

    println!(" generate for ENV - {}", p_env_str);

    for class_str in p_classes_lst {

        // CLASS_DIR
        let class_target_dir_str = format!("{}/{}/{}",
            p_gf_dataset_config.target_dir_path_str,
            p_env_str,
            class_str);

        // fs::create_dir_all() - creates target dir and all parent dirs.
        fs::create_dir_all(&class_target_dir_str)
            .expect(&format!("ERROR!! - creation of the class [{}] dir failed", class_str));

        match class_str.as_ref() {

            "rect" => {

                draw_rects(class_target_dir_str,
                    p_gf_dataset_config,
                    p_gf_runtime_gfx);
            },

            _ => {
                println!("ERROR!! - this class [{}] is not supported yet!", class_str);
            }
        }
    }
}

//-------------------------------------------------
// DRAW_RECT
#[allow(non_snake_case)]
fn draw_rects(p_target_dir_str: String,
    p_gf_dataset_config: &GFdatasetConfig,
    p_gf_runtime_gfx: &GFruntimeGfx) {

    let gfx_ctx         = &p_gf_runtime_gfx.ctx;
    let randomize_bool  = true;
    let class_str = "rect";


    for i in 0..p_gf_dataset_config.elements_num_int {

        // BACKGROUND_COLOR
        if randomize_bool {
            gfx_ctx.set_source_rgb(rand::random::<f64>(), rand::random::<f64>(), rand::random::<f64>());
        } else {
            gfx_ctx.set_source_rgb(1.0, 1.0, 1.0);
        }
        gfx_ctx.paint();

        // RECTANGLE
        let rect_width_f  = rand::random::<f64>() * p_gf_dataset_config.image_width_int as f64;
        let rect_height_f = rand::random::<f64>() * p_gf_dataset_config.image_height_int as f64;

        let x = rand::random::<f64>() * p_gf_dataset_config.image_width_int as f64 - rect_width_f;
        let y = rand::random::<f64>() * p_gf_dataset_config.image_height_int as f64 - rect_height_f;

        gfx_ctx.rectangle(x, y,
            rect_width_f,
            rect_height_f);

        // RECTANGLE_COLOR
        if randomize_bool {
            gfx_ctx.set_source_rgb(rand::random::<f64>(), rand::random::<f64>(), rand::random::<f64>());
        } else {
            gfx_ctx.set_source_rgb(1.0, 1.0, 1.0);
        }
        gfx_ctx.fill();

        // SAVE_FILE
        // IMPORTANT!! - ".png" is relevant here, because currently gf_core Rust Crate
        //               Cargo.toml specifies "cairo-rs" dependency with feature for "png"
        //               enabled only. have not tested saving cairo Surfaces to .jpeg.
        let output_file_path_str = format!("{}/{}-{}-{}.png",
            p_target_dir_str,
            p_gf_dataset_config.name_str,
            &class_str,
            i);

        gf_core::gf_image::save_cairo(&p_gf_runtime_gfx.canvas,
            &output_file_path_str);
    }
}

//-------------------------------------------------
#[allow(non_snake_case)]
fn runtime_get_graphics(p_image_width_int: u64,
    p_image_height_int: u64) -> GFruntimeGfx {

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