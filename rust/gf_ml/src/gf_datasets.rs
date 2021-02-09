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

use std::f64::consts::PI;
use std::fs;
// use std::fs::File;
use rand;
use image;
use cairo;
use tensorflow;

use gf_core;
use crate::gf_ml_client;

//-------------------------------------------------
struct GFdatasetConfig {
    name_str:            String,
    elements_num_int:    u64,
    image_width_int:     u64,
    image_height_int:    u64,
    target_dir_path_str: String,
    save_img_files_bool: bool,
}

struct GFruntimeGfx {
    surface: cairo::ImageSurface,
}

//-------------------------------------------------
// GENERATE_AND_REGISTER - generates a dataset and registers the newly generated dataset with
//                         a remote GF ML server.

pub fn generate_and_register(p_dataset_name_str: String,
    p_classes_lst:         Vec<String>,
    p_elements_num_int:    u64,
    p_image_width_int:     u64,
    p_image_height_int:    u64,
    p_target_dir_path_str: String,
    p_gf_ml_host_str:      String) {



    


    generate(p_dataset_name_str,
        p_classes_lst,
        p_elements_num_int,
        p_image_width_int,
        p_image_height_int,
        p_target_dir_path_str);



    let url_str = format!("{}/ml/dataset/create", p_gf_ml_host_str);
    gf_ml_client::get_blocking(url_str.as_ref());

}

//-------------------------------------------------
// GENERATE
pub fn generate(p_dataset_name_str: String,
    p_classes_lst:         Vec<String>,
    p_elements_num_int:    u64,
    p_image_width_int:     u64,
    p_image_height_int:    u64,
    p_target_dir_path_str: String) {

    let mut gf_runtime_gfx = runtime_get_graphics(p_image_width_int,
        p_image_height_int);

    let gf_dataset_config = GFdatasetConfig{
        name_str:            p_dataset_name_str,
        elements_num_int:    p_elements_num_int,
        image_width_int:     p_image_width_int,
        image_height_int:    p_image_height_int,
        target_dir_path_str: p_target_dir_path_str,
        save_img_files_bool: true,
    };

    // TFRECORDS_FILE_PATHS
    let tf_records_file_path__train_str    = format!("{}/tfrecords/{}__train.tfrecords", &gf_dataset_config.target_dir_path_str, &gf_dataset_config.name_str);
    let tf_records_file_path__validate_str = format!("{}/tfrecords/{}__validate.tfrecords", &gf_dataset_config.target_dir_path_str, &gf_dataset_config.name_str);

    // TRAIN_DATASET
    generate_for_env(&"train",
        &p_classes_lst,
        &tf_records_file_path__train_str,
        &gf_dataset_config,
        &mut gf_runtime_gfx);

    // VALIDATION_DATASET
    generate_for_env(&"validate",
        &p_classes_lst,
        &tf_records_file_path__validate_str,
        &gf_dataset_config,
        &mut gf_runtime_gfx);
}

//-------------------------------------------------
// GENERATE_FOR_ENVIRONMENT
fn generate_for_env(p_env_str: &str,
    p_classes_lst:                    &Vec<String>,
    p_tfrecords_output_file_path_str: &str,
    p_gf_dataset_config:              &GFdatasetConfig,
    p_gf_runtime_gfx:                 &mut GFruntimeGfx) {

    println!(" generate for ENV - {}", p_env_str);


    // TF_RECORDS_WRITER
    // same writer used for all classes, so that they're all mixed up. there is even further
    // shuffling done by TensorFlow in its data pipelines durring training after these tfrecords
    // files are loaded.
    let mut tf_records_writer = gf_core::gf_tf::get_tf_records__writer(p_tfrecords_output_file_path_str);


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
            

            "circle" => {
                
                draw_circles(class_target_dir_str,
                    &mut tf_records_writer,
                    p_gf_dataset_config,
                    p_gf_runtime_gfx);
            }, 


            "rect" => {

                draw_rects(class_target_dir_str,
                    &mut tf_records_writer,
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
// DRAW_CIRCLES
fn draw_circles(p_target_dir_str: String,
    p_tf_records_writer: &mut tensorflow::io::RecordWriter<std::io::BufWriter<std::fs::File>>,
    p_gf_dataset_config: &GFdatasetConfig,
    p_gf_runtime_gfx:    &mut GFruntimeGfx) {

    
    // CONFIG
    let randomize_bool = true;
    let class_str      = "circle";
    let label_int      = 1;
    let img_width_int  = p_gf_dataset_config.image_width_int;
    let img_height_int = p_gf_dataset_config.image_height_int;
    

    for i in 0..p_gf_dataset_config.elements_num_int {

        // CAIRO_CONTEXT
        // IMPORTANT!! - create a context per image generated. mainly so that at the end of drawing
        //               the context can be droped (memory freed) and the reference that it holds to the
        //               Cairo surface released.
        //               this has to be done to avoid a runtime error when doing a surface.get_data()
        let ctx = cairo::Context::new(&p_gf_runtime_gfx.surface);

        // BACKGROUND_COLOR
        if randomize_bool {
            ctx.set_source_rgb(rand::random::<f64>(), rand::random::<f64>(), rand::random::<f64>());
        } else {
            ctx.set_source_rgb(1.0, 1.0, 1.0);
        }
        ctx.paint();



        let x = rand::random::<f64>() * img_width_int as f64;
        let y = rand::random::<f64>() * img_height_int as f64;
        let r = rand::random::<f64>() * img_width_int as f64;

        ctx.arc(x, y, r, 0.0, 2.0*PI);


        ctx.set_source_rgb(rand::random::<f64>(),
            rand::random::<f64>(),
            rand::random::<f64>());


        ctx.fill();



        //-----------------
        // SAVE_FILE
        if p_gf_dataset_config.save_img_files_bool {

            // IMPORTANT!! - ".png" is relevant here, because currently gf_core Rust Crate
            //               Cargo.toml specifies "cairo-rs" dependency with feature for "png"
            //               enabled only. have not tested saving cairo Surfaces to .jpeg.
            let output_file_path_str = format!("{}/{}-{}-{}.png",
                p_target_dir_str,
                p_gf_dataset_config.name_str,
                &class_str,
                i);

            gf_core::gf_image::save_cairo(&p_gf_runtime_gfx.surface,
                &output_file_path_str);
        }

        //-----------------
        // IMG_BUFFER

        // required before accessing the pixel data to ensure that all pending drawing operations are finished
        p_gf_runtime_gfx.surface.flush();
        
        // IMPORTANT!! - cairo::Surface.get_data() - Get a pointer to the data of the image surface,
        //                                           for direct inspection or modification.
        //               https://gtk-rs.org/docs/cairo/struct.Surface.html
        //               drop() - critical for the context to be dropped (its memory released) so that its
        //                        reference to the surface is released. without this surface.get_data() 
        //                        will cause a runtime error (but it will compile):
        //                        "thread '<unnamed>' panicked at 'called `Result::unwrap()` on an `Err` value: NonExclusive'"
        drop(ctx);
        let surface_data          = p_gf_runtime_gfx.surface.get_data().unwrap();
        let surface_data: Vec<u8> = surface_data.to_vec();
        let img_buffer: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = image::ImageBuffer::from_vec(img_width_int as u32,
            img_height_int as u32,
            surface_data).unwrap();
        
        // WRITE_TF_RECORD
        gf_core::gf_tf::write_tf_example__from_img_buffer(img_buffer,
            label_int,
            p_tf_records_writer);
            
        //-----------------
    }
}

//-------------------------------------------------
// DRAW_RECTANGLES
fn draw_rects(p_target_dir_str: String,
    p_tf_records_writer: &mut tensorflow::io::RecordWriter<std::io::BufWriter<std::fs::File>>,
    p_gf_dataset_config: &GFdatasetConfig,
    p_gf_runtime_gfx:    &mut GFruntimeGfx) {

    // CONFIG
    let randomize_bool = true;
    let class_str      = "rect";
    let label_int      = 0;
    let img_width_int  = p_gf_dataset_config.image_width_int;
    let img_height_int = p_gf_dataset_config.image_height_int;

    for i in 0..p_gf_dataset_config.elements_num_int {
    
        // CAIRO_CONTEXT
        // IMPORTANT!! - create a context per image generated. mainly so that at the end of drawing
        //               the context can be droped (memory freed) and the reference that it holds to the
        //               Cairo surface released.
        //               this has to be done to avoid a runtime error when doing a surface.get_data()
        let ctx = cairo::Context::new(&p_gf_runtime_gfx.surface);

        // BACKGROUND_COLOR
        if randomize_bool {
            ctx.set_source_rgb(rand::random::<f64>(), rand::random::<f64>(), rand::random::<f64>());
        } else {
            ctx.set_source_rgb(1.0, 1.0, 1.0);
        }
        ctx.paint();

        // RECTANGLE
        let rect_width_f  = rand::random::<f64>() * img_width_int as f64;
        let rect_height_f = rand::random::<f64>() * img_height_int as f64;

        let x = rand::random::<f64>() * img_width_int as f64 - rect_width_f;
        let y = rand::random::<f64>() * img_height_int as f64 - rect_height_f;

        ctx.rectangle(x, y,
            rect_width_f,
            rect_height_f);

        // RECTANGLE_COLOR
        if randomize_bool {
            ctx.set_source_rgb(rand::random::<f64>(), rand::random::<f64>(), rand::random::<f64>());
        } else {
            ctx.set_source_rgb(1.0, 1.0, 1.0);
        }
        ctx.fill();

        //-----------------
        // SAVE_FILE
        if p_gf_dataset_config.save_img_files_bool {

            // IMPORTANT!! - ".png" is relevant here, because currently gf_core Rust Crate
            //               Cargo.toml specifies "cairo-rs" dependency with feature for "png"
            //               enabled only. have not tested saving cairo Surfaces to .jpeg.
            let output_file_path_str = format!("{}/{}-{}-{}.png",
                p_target_dir_str,
                p_gf_dataset_config.name_str,
                &class_str,
                i);

            gf_core::gf_image::save_cairo(&p_gf_runtime_gfx.surface,
                &output_file_path_str);
        }

        //-----------------
        // IMG_BUFFER

        // required before accessing the pixel data to ensure that all pending drawing operations are finished
        p_gf_runtime_gfx.surface.flush();
        
        // IMPORTANT!! - cairo::Surface.get_data() - Get a pointer to the data of the image surface,
        //                                           for direct inspection or modification.
        //               https://gtk-rs.org/docs/cairo/struct.Surface.html
        //               drop() - critical for the context to be dropped (its memory released) so that its
        //                        reference to the surface is released. without this surface.get_data() 
        //                        will cause a runtime error (but it will compile):
        //                        "thread '<unnamed>' panicked at 'called `Result::unwrap()` on an `Err` value: NonExclusive'"
        drop(ctx);
        let surface_data          = p_gf_runtime_gfx.surface.get_data().unwrap();
        let surface_data: Vec<u8> = surface_data.to_vec();
        let img_buffer: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = image::ImageBuffer::from_vec(img_width_int as u32,
            img_height_int as u32,
            surface_data).unwrap();
        
        // WRITE_TF_RECORD
        gf_core::gf_tf::write_tf_example__from_img_buffer(img_buffer,
            label_int,
            p_tf_records_writer);

        //-----------------
    }
}

//-------------------------------------------------
// RUNTIME_GET_GRAPHICS
fn runtime_get_graphics(p_image_width_int: u64,
    p_image_height_int: u64) -> GFruntimeGfx {

    /*let surface = cairo::ImageSurface::create(cairo::Format::ARgb32,
        p_image_width_int as i32,
        p_image_height_int as i32)
        .expect("failed to create a drawing surface with the Cairo backend");*/
        
    let buff: Vec<u8> = vec![0; (p_image_width_int * p_image_height_int * 4) as usize];

    // the number of bytes between the start of rows in the buffer as allocated.
    // this value should always be computed by cairo_format_stride_for_width() before allocating the data buffer.
    // however it seems that Rust Cairo-rs lib doesnt have the stride_for_width() function, so doing it manually here.
    let stride_int: i32 = p_image_width_int as i32 * 4;

    // creates an image surface for the provided pixel data
    let surface = cairo::ImageSurface::create_for_data(buff,
        cairo::Format::ARgb32,
        p_image_width_int as i32,
        p_image_height_int as i32,
        stride_int).unwrap();

    let gf_runtime_gfx = GFruntimeGfx{
        surface: surface,
    };

    return gf_runtime_gfx;
}