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

#[macro_use] extern crate cpython;
use cpython::{PyObject, PyResult, Python, PyTuple, PyDict};

mod ml_datasets;

//-------------------------------------------------
py_module_initializer!(gf_data_viz, initgf_data_viz, PyInit_gf_data_viz, |py, m| {

    m.add(py, "__doc__", "GloFlow data visualization")?;
    m.add(py, "ml_datasets__generate", py_fn!(py, gf_ml_datasets_generate(a: String, b: u32, c: u32, d: String)))?;
    Ok(())
});

//-------------------------------------------------
fn gf_ml_datasets_generate(p_py :Python,
    p_dataset_name_str    : String,
    p_img_width_int       : u32,
    p_img_height_int      : u32,
    p_target_dir_path_str : String) -> PyResult<PyObject> {

    ml_datasets::generate(p_dataset_name_str,
        p_img_width_int,
        p_img_height_int,
        p_target_dir_path_str);

    Ok(p_py.None())
}