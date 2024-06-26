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

use pyo3::wrap_pyfunction;
use pyo3::prelude::*;
use numpy::{PyArray2, PyArray3, PyArray4};

// use numpy::{IntoPyArray, PyArrayDyn};
// use std::collections::HashMap;
// use image::{GenericImageView};
// use ndarray::{ArrayD, ArrayViewD, ArrayViewMutD};
// use gf_core;

// use gf_ml;
use gf_images_jobs;

mod gf_numpy_view;

//-------------------------------------------------
// PY_C_API
//-------------------------------------------------
#[pymodule]
fn gf_images_jobs_py(_py: Python, m: &PyModule) -> PyResult<()> {
    m.add_wrapped(wrap_pyfunction!(apply_transforms))?;
    m.add_wrapped(wrap_pyfunction!(create_collage))?;
    m.add_wrapped(wrap_pyfunction!(view_numpy_arr_2D))?;
    m.add_wrapped(wrap_pyfunction!(view_numpy_arr_3D))?;
    m.add_wrapped(wrap_pyfunction!(view_numpy_arr_4D))?;
    m.add_wrapped(wrap_pyfunction!(generate_ml_dataset))?;
    // m.add_wrapped(wrap_pyfunction!(view_ml_dataset))?;

    Ok(())
}

//-------------------------------------------------
// APPLY_TRANSFORMS
#[pyfunction]
fn apply_transforms(// p_py: Python,
    p_transforms_lst:           Vec<String>,
    p_img_source_file_path_str: String,
    p_img_target_file_path_str: String) -> PyResult<()> {

    gf_images_jobs::apply_transforms(p_transforms_lst,
        &p_img_source_file_path_str,
        &p_img_target_file_path_str);

    Ok(())
}

//-------------------------------------------------
// CREATE_COLLAGE
#[pyfunction]
fn create_collage(// p_py: Python,
    p_img_file_paths_lst:         Vec<String>,
    p_output_img_file_path_c_str: String,
    p_width_int:                  u64,
    p_height_int:                 u64,
    p_rows_num_int:               u32,
    p_columns_num_int:            u32) -> PyResult<()> {
    
    gf_images_jobs::create_collage(p_img_file_paths_lst,
        p_output_img_file_path_c_str,
        p_width_int,
        p_height_int,
        p_rows_num_int,
        p_columns_num_int);

    Ok(())
}

//-------------------------------------------------
// VIEW_NUMPY_ARR_2D
#[pyfunction]
fn view_numpy_arr_2D(p_numpy_2d_lst: &PyArray2<f64>,
    p_img_target_file_path_str: String) -> PyResult<()> {

    gf_numpy_view::arr_2D(p_numpy_2d_lst,
        p_img_target_file_path_str);

    Ok(())
}

//-------------------------------------------------
// VIEW_NUMPY_ARR_3D
#[pyfunction]
fn view_numpy_arr_3D(p_numpy_3d_lst: &PyArray3<f64>,
    p_img_target_file_path_str: String) -> PyResult<()> {

    gf_numpy_view::arr_3D(p_numpy_3d_lst,
        p_img_target_file_path_str);

    Ok(())
}

//-------------------------------------------------
// VIEW_NUMPY_ARR_4D
#[pyfunction]
fn view_numpy_arr_4D(p_numpy_4d_lst: &PyArray4<f64>,
    p_img_target_file_path_str: String,
    p_width_int:       u64,
    p_height_int:      u64,
    p_rows_num_int:    u32,
    p_columns_num_int: u32) -> PyResult<()> {

    gf_numpy_view::arr_4D(p_numpy_4d_lst,
        p_img_target_file_path_str,
        p_width_int,
        p_height_int,
        p_rows_num_int,
        p_columns_num_int);

    Ok(())
}

//-------------------------------------------------
// GENERATE_ML_DATASET_TO_TFRECORDS
#[pyfunction]
fn generate_ml_dataset(p_dataset_name_str: String,
    p_classes_lst:         Vec<String>,
    p_elements_num_int:    u64,
    p_img_width_int:       u64,
    p_img_height_int:      u64,
    p_target_dir_path_str: String) -> PyResult<()> {

    gf_images_jobs::generate_ml_dataset(p_dataset_name_str,
        p_classes_lst,
        p_elements_num_int,
        p_img_width_int,
        p_img_height_int,
        p_target_dir_path_str);

    Ok(())
}

//-------------------------------------------------
// fn generate_and_register_ml_dataset() {
// 
// }

//-------------------------------------------------
/*
// VIEW_ML_DATASET_FROM_TFRECORDS
#[pyfunction]
fn view_ml_dataset(p_tfrecords_file_path_str: String,
    p_img_target_file_path_str:   String,
    p_tf_example__img_width_int:  u64,
    p_tf_example__img_height_int: u64,
    p_collage__img_width_int:     u64,
    p_collage__img_height_int:    u64,
    p_collage__rows_num_int:      u32,
    p_collage__columns_num_int:   u32) -> PyResult<()> {

    gf_ml::gf_datasets_view::view_tf_records(&p_tfrecords_file_path_str,
        &p_img_target_file_path_str,
        p_tf_example__img_width_int,
        p_tf_example__img_height_int,
        p_collage__img_width_int,
        p_collage__img_height_int,
        p_collage__rows_num_int,
        p_collage__columns_num_int);

    Ok(())
}
*/