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

use gf_ml;
use gf_core;

//-------------------------------------------------
#[allow(non_snake_case)]
pub fn ml_dataset_to_tfrecords(p_dataset_name_str: String,
    p_classes_lst:         Vec<String>,
    p_elements_num_int:    u64,
    p_img_width_int:       u64,
    p_img_height_int:      u64,
    p_target_dir_path_str: String) {


    println!("ML generating in Rust - {} - {}/{} - {}",
        p_dataset_name_str,
        p_img_width_int,
        p_img_height_int,
        p_target_dir_path_str);

    gf_ml::gf_datasets::generate(p_dataset_name_str,
        p_classes_lst,
        p_elements_num_int,
        p_img_width_int,
        p_img_height_int,
        p_target_dir_path_str);


    let output_file_path_str = "./data/output_ml/gf_rust_test.tfrecords";
    gf_core::gf_tf::get_tf_records__writer(output_file_path_str);


}