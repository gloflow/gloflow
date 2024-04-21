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

// use protoc_rust;

fn main() {
    
    println!("cargo running build.rs >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>");

    /*
    // TENSORFLOW_PROTOBUFF_DEFS
    let protobuff_input_lst = [
        "src/gf_protobuff/tf_feature.proto",
        "src/gf_protobuff/tf_example.proto"
    ];

    // IMPORTANT!! - generates .rs files from protobuff definitions.
    //               this is run only once durring the build stage, 
    //               and generates static .rs files, which are then used in
    //               the final build of the Rust binary.
    protoc_rust::run(protoc_rust::Args {
        out_dir:   "src/gf_protobuff",
        input:     &protobuff_input_lst,
        includes:  &["src/gf_protobuff"],
        customize: protoc_rust::Customize{
            ..Default::default()
        },

    }).expect("ERROR!! - failed to build Rust protobuffers in gf_core")
    */
}