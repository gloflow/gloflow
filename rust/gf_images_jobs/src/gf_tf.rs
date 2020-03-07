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
use tensorflow;
use crate::gf_image_io;

//-------------------------------------------------
pub fn write_file() {
    
    
    
    let output_file_path_str = "gf_test.tfrecord";
    


    let f = ::std::fs::OpenOptions::new()
        .write(true)
        .create(true)
        .open(output_file_path_str)
        .unwrap();

    
    // TFRecord file is a sequence of strings
    let mut record_writer = tensorflow::io::RecordWriter::new(::std::io::BufWriter::new(f));





    for x in 0..10 {

        // IMAGE_BUFFER
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
        record_bytes_lst.extend(img_bytes_lst);

        //-----------------
        // WRITE_RECORD
        record_writer.write_record(&record_bytes_lst).unwrap(); // test_data_str.as_bytes()).unwrap();
    
    }

    
    

}

