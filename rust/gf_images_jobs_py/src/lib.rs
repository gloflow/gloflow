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

// PyObject - a reference to a Python object
#[macro_use] extern crate cpython;
use cpython::{PyObject, PyResult, Python}; // PyTuple, PyDict};
use std::collections::HashMap;

use gf_images_jobs;

//-------------------------------------------------
// PY_C_API
//-------------------------------------------------
py_module_initializer!(gf_images_jobs_py, initgf_images_jobs_py, PyInit_gf_images_jobs_py, |py, m| {

    m.add(py, "__doc__", "GloFlow images_jobs")?;

    m.add(py, "apply_transforms", py_fn!(py, 
        py__apply_transforms(a: Vec<String>,
            b: String,
            c: String)))?;

    m.add(py, "create_collage", 
        py_fn!(py, py__create_collage(a: Vec<String>,
            b: String,
            c: u32,
            d: u32,
            e: u32,
            f: u32)))?;
    Ok(())
});

//-------------------------------------------------
#[allow(non_snake_case)]
fn py__apply_transforms(p_py: Python,
    p_transforms_lst:           Vec<String>,
    p_img_source_file_path_str: String,
    p_img_target_file_path_str: String) -> PyResult<PyObject> {

    gf_images_jobs::apply_transforms(p_transforms_lst,
        &p_img_source_file_path_str,
        &p_img_target_file_path_str);

    Ok(p_py.None())
}

//-------------------------------------------------
#[allow(non_snake_case)]
fn py__create_collage(p_py: Python,
    p_img_file_paths_lst:         Vec<String>,
    p_output_img_file_path_c_str: String,
    p_width_int:                  u32,
    p_height_int:                 u32,
    p_rows_num_int:               u32,
    p_columns_num_int:            u32) -> PyResult<PyObject> {
    

    gf_images_jobs::create_collage(p_img_file_paths_lst,
        p_output_img_file_path_c_str,
        p_width_int,
        p_height_int,
        p_rows_num_int,
        p_columns_num_int);

    Ok(p_py.None())
}