# GloFlow application and media management/publishing platform
# Copyright (C) 2020 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

from colored import fg, bg, attr

sys.path.append("%s/../../py/gf_core"%(modd_str))
import gf_core_cli

#--------------------------------------------------
# RUST NIGHTLY
# rustup self update                 - rustup update
# rustup toolchain install nightly   - rustup install Rust nightly
# rustup run nightly rustc --version - rstup test Rust nightly
# rustup default nightly             - rustup make Rust nightly the global default

#--------------------------------------------------
# RUN
def run(p_cargo_crate_dir_path_str,
    p_static_bool       = False,
    p_exit_on_fail_bool = True,
    p_verbose_bool      = False):
    assert os.path.isdir(p_cargo_crate_dir_path_str)
    
    print(f"{fg('yellow')}BUILD{attr(0)}")
    
    cwd_str = os.getcwd()
    os.chdir(os.path.abspath(p_cargo_crate_dir_path_str)) # change into the target main package dir

    #-------------
    # "rustup update stable"
    # _, _, exit_code_int = gf_cli_utils.run_cmd("cargo clean")

    #-------------

    c_lst = [
        # 'RUSTFLAGS="$RUSTFLAGS -A warnings"', # turning off rustc warnings
        # "RUSTFLAGS='-L %s'"%(os.path.abspath("%s/../../rust/gf_images_jobs/test"%(modd_str))),
        "cargo build",
    ]

    if p_verbose_bool:
        c_lst.append("--verbose")

    # STATIC_LINKING - some outputed libs (imported by Go for example) should contain their
	#                  own versions of libs statically linked into them.
    if p_static_bool:

        #-------------
        # MUSL - staticaly compile libc compatible lib into the output binary. without MUSL
        #        rust statically compiles all program libs except the standard lib.
        # musl-gcc - musl-gcc is a wrapper around GCC that uses the musl C standard library
        #            implementation to build programs. It is well suited for being linked with other libraries
        #            into a single static executable with no shared dependencies.
        #            its used by "cargo build" if we target linux-musl.
        #            "sudo apt-get install musl-tools" - make sure its installed
        #
        # x86_64-unknown-linux-musl - for 64-bit Linux.
        #                             for this to work "rustup" has to be used to install this
        #                             build target into the Rust toolchain. 
        #                             (for GF CI this is done in the gf_builder Dockerfile__gf_builder)
        c_lst.append("--target x86_64-unknown-linux-musl")

        #-------------

    # DYNAMIC_LINKING
    else:
        c_lst.append("--release")

    # _, _, exit_code_int = gf_cli_utils.run_cmd(" ".join(c_lst))
    _, _, exit_code_int = gf_core_cli.run(" ".join(c_lst))

    # IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
    if not exit_code_int == 0:
        if p_exit_on_fail_bool:
            exit(exit_code_int)

    os.chdir(cwd_str) # return to initial dir

#--------------------------------------------------
def prepare_libs(p_name_str,
    p_cargo_crate_dir_path_str,
    p_type_str,
    p_exit_on_fail_bool = True):
    assert os.path.isdir(p_cargo_crate_dir_path_str)
    assert p_type_str == "lib_rust"


    print(f"{fg('yellow')}PREPARE LIBS{attr(0)}>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

    #-------------
    # EXTERN_LIB
    prepare_libs__extern(p_exit_on_fail_bool = p_exit_on_fail_bool)

    #-------------

    target_build_dir_path_str = os.path.abspath("%s/../../rust/build"%(modd_str))
    assert os.path.isdir(target_build_dir_path_str)

    target_lib_file_path_lst = []
    if p_type_str == "lib_rust":
        

        #-------------
        # FIX!! - dont hardcode the app_name here like this, but parse Cargo.toml to detect if 
        #         one of the Crate types is "staticlib".
        if p_name_str == "gf_images_jobs":

        
            # RUST_PY - CPYTHON_EXTENSION - this lib is Python extension written in Rust.
            #                               at the moment in GF the convention is for these Rust libs to have a postfix "_py".
            if p_cargo_crate_dir_path_str.endswith("_py"):
                
                # DYNAMIC_LIB
                # IMPORTANT!! - Rust compiles this dynamic lib with the "lib" prefix, but the Python VM
                #               requires extension libs to not have the "lib" prefix.
                source__py_lib_file_path_str = "%s/target/release/lib%s_py.so"%(p_cargo_crate_dir_path_str, p_name_str)
                target__py_lib_file_path_str = "%s/%s_py.so"%(target_build_dir_path_str, p_name_str)

                assert os.path.isfile(source__py_lib_file_path_str)
                target_lib_file_path_lst.append((source__py_lib_file_path_str, target__py_lib_file_path_str))

            else:

                # STATIC_LIB
                source__static_lib_file_path_str = "%s/target/release/lib%s.a"%(p_cargo_crate_dir_path_str, p_name_str)
                assert os.path.isfile(source__static_lib_file_path_str)
                target_lib_file_path_lst.append((source__static_lib_file_path_str, target_build_dir_path_str))

                # DYNAMIC_LIB
                source__dynamic_lib_file_path_str = "%s/target/release/lib%s.so"%(p_cargo_crate_dir_path_str, p_name_str)
                assert os.path.isfile(source__dynamic_lib_file_path_str)
                target_lib_file_path_lst.append((source__dynamic_lib_file_path_str, target_build_dir_path_str))

        #-------------
        # ALL
        else:
            target_lib_file_path_lst.append(("%s/target/release/lib%s.so"%(p_cargo_crate_dir_path_str, p_name_str), target_build_dir_path_str))

        #-------------


    # COPY_FILES
    for source_f, target_f in target_lib_file_path_lst:
        c_str = "cp %s %s"%(source_f, target_f)
        _, _, exit_code_int = gf_core_cli.run(c_str) # gf_cli_utils.run_cmd(c_str)



    # IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
    if not exit_code_int == 0:
        if p_exit_on_fail_bool:
            exit(exit_code_int)
    

#--------------------------------------------------
def prepare_libs__extern(p_exit_on_fail_bool = True):


    #-------------
    # TENSORFLOW
    print(f"{fg('green')}prepare TensorFlow lib{attr(0)}")

    tf__version_str  = "1.15.0"
    tf__filename_str = f"libtensorflow-cpu-linux-x86_64-{tf__version_str}.tar.gz"
    tf__url_str      = f"https://storage.googleapis.com/tensorflow/libtensorflow/{tf__filename_str}"



    # DOWNLOAD
    _, _, exit_code_int = gf_core_cli.run(f"curl {tf__url_str} --output tflib.tar.gz")
    if not exit_code_int == 0:
        if p_exit_on_fail_bool:
            exit(exit_code_int)

    #-------------
    

    # FIX!! - COMPLETE!!
    #         download TF lib and place it in appropriate dir, to have a fresh TF libs
    #         in the build server context, without including the lib in the repo itself.
    print("FIIIIIIIIIIINIIIIIIIISSHHH!!!")