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

use std::io;
use std::collections::{HashMap};
use image;
use tensorflow;
use protobuf;
use protobuf::Message;

use png;
use png::HasParameters; // needed for png_encoder.set() call

use crate::gf_protobuff::tf_feature::{Features, Feature, Int64List, BytesList};
use crate::gf_protobuff::tf_example::{Example};

use crate::gf_image_io;

//-------------------------------------------------
pub fn write_file(p_output_file_path_str: &str) {
    
    let label_int         = 0 as i64;
    let img_file_path_str = "data/output_ml/generated/train/rect/test-rect-0.png";
    
    
    
    let f = ::std::fs::OpenOptions::new()
        .write(true)
        .create(true)
        .open(p_output_file_path_str)
        .unwrap();

    
    // TFRecord file is a sequence of strings
    let mut record_writer = tensorflow::io::RecordWriter::new(::std::io::BufWriter::new(f));



    //-----------------
    // IMAGE_BUFFER
    
    let gf_img            = gf_image_io::native__open_image(img_file_path_str);
    let gf_img_buff: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = gf_img.raw_img.to_rgba();

    let mut buf_writer = Vec::new(); // std::io::BufWriter::new(vec![]);

    {
        let mut png_encoder = png::Encoder::new(&mut buf_writer,
            gf_img_buff.width(),
            gf_img_buff.height());

        png_encoder.set(png::ColorType::RGBA).set(png::BitDepth::Eight);
        
        let mut png_encoder_writer = png_encoder.write_header().unwrap();

        // let data = [255, 0, 0, 255, 0, 0, 0, 255]; // An array containing a RGBA sequence
        let o = png_encoder_writer.write_image_data(&gf_img_buff.into_raw()).unwrap();
    }

    let img_png_encoded_data_lst: &[u8] = &buf_writer;

    // Vec<Vec<u8>> - used because protobuf::RepeatedField::from_vec() requires a 2D array.
    let img_bytes_lst: Vec<Vec<u8>> = vec![img_png_encoded_data_lst.to_vec()]; // vec![gf_img_buff.into_raw()];
    




    //-----------------
    // FEATURE_LABEL
    let mut tf_feature_label = Feature::new();
    let mut tf_label_int     = Int64List::new();
    tf_label_int.set_value(vec![label_int]);
    tf_feature_label.set_int64_list(tf_label_int);

    //-----------------
    // FEATURE_IMG
    let mut tf_feature_img = Feature::new();
    let mut tf_img_bytes   = BytesList::new();

    tf_img_bytes.set_value(protobuf::RepeatedField::from_vec(img_bytes_lst));
    tf_feature_img.set_bytes_list(tf_img_bytes);

    //-----------------
    // FEATURES
    
    let mut tf_feature_map = HashMap::new();
    tf_feature_map.insert("label".to_string(), tf_feature_label);
    tf_feature_map.insert("img".to_string(),   tf_feature_img);

    let mut tf_features = Features::new();
    tf_features.set_feature(tf_feature_map);

    //-----------------
    // EXAMPLE
    let mut tf_example = Example::new();
    tf_example.set_features(tf_features);

    let tf_example_bytes_lst = tf_example.write_to_bytes().unwrap();

    //-----------------

    record_writer.write_record(&tf_example_bytes_lst).unwrap();

    for x in 0..20 {

        /*// IMAGE_BUFFER
        let img_file_path_str = "test__numpy_3d.jpeg";
        let gf_img            = gf_image_io::native__open_image(img_file_path_str);
        let gf_img_buff: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = gf_img.raw_img.to_rgba();

        //-----------------
        // FEATURES
        let img_label_int: u8      = 3;
        let img_bytes_lst: Vec<u8> = gf_img_buff.into_raw();

        // create a record as a list of bytes (u8), and concatenate each "feature/dimension"
        // together into that single list of bytes.
        let mut record_bytes_lst: Vec<u8> = Vec::new();
        record_bytes_lst.push(img_label_int);
        record_bytes_lst.extend(img_bytes_lst);*/

        //-----------------
        // WRITE_RECORD
        
        println!("===")
    }

    
    

}

