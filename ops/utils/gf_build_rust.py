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
# RUN_IN_CONTAINER
def run_in_cont():

    repo_local_path_str = os.path.abspath(f'{modd_str}/../../../gloflow').strip()
    cmd_lst = [
        "sudo", "docker", "run",
        "--rm", # remove after exit 
        "-v", f"{repo_local_path_str}:/home/gf", # mount repo into the container
        "glofloworg/gf_builder_rust_ubuntu:latest",

        # FIX!! - stop using gf_builder.py!! use "/home/gf/ops/cli__build.py" instead!
        "python3", "-u", "/home/gf/build/gf_builder/gf_builder.py", "-run=build_rust"
    ]
    p = gf_core_cli.run__view_realtime(cmd_lst, {},
        "gf_build_rust", "green")

    p.wait()

#--------------------------------------------------
# RUN
def run(p_cargo_crate_dir_path_str,
    p_static_bool       = False,
    p_exit_on_fail_bool = True,
    p_verbose_bool      = False):
    assert os.path.isdir(p_cargo_crate_dir_path_str)
    
    print(f"{fg('yellow')}BUILD{attr(0)}")
    print(f"crate dir - {fg('yellow')}{p_cargo_crate_dir_path_str}{attr(0)}")

    cwd_str = os.getcwd()
    os.chdir(os.path.abspath(p_cargo_crate_dir_path_str)) # change into the target main package dir

    #-------------
    # "rustup update stable"
    # _, _, exit_code_int = gf_cli_utils.run_cmd("cargo clean")

    #-------------

    c_lst = []

    if p_static_bool:

        # DOCUMENT!! - without this the py extension wont compile. 
        #              complaining that target for musl from gf_images_job lib cant be used, since this py extension
        #              package is marked as a dynamic lib (which it has to be to be importable by the Py VM).
        if os.path.basename(p_cargo_crate_dir_path_str) == "gf_images_jobs_py":
            c_lst.append("RUSTFLAGS='-C target-feature=-crt-static'")

    c_lst.extend([
        # 'RUSTFLAGS="$RUSTFLAGS -A warnings"', # turning off rustc warnings
        # "RUSTFLAGS='-L %s'"%(os.path.abspath("%s/../../rust/gf_images_jobs/test"%(modd_str))),

        # if compiling on Ubuntu for Alpine for example, this ENV var should be set
        # "PKG_CONFIG_ALLOW_CROSS=1",
        
        "cargo build",
    ])

    if p_verbose_bool:
        # c_lst.append("--verbose")
        c_lst.append("-vv") # very verbose


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


    cmd_str = " ".join(c_lst)
    print(cmd_str)
    _, _, exit_code_int = gf_core_cli.run(cmd_str)

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

    

    target_build_dir_path_str = os.path.abspath("%s/../../rust/build"%(modd_str))
    assert os.path.isdir(target_build_dir_path_str)

    target_lib_file_path_lst = []
    if p_type_str == "lib_rust":

        release_dir_str = f"{p_cargo_crate_dir_path_str}/target/release"
        gf_core_cli.run(f"ls -al {release_dir_str}")    

        #-------------
        # FIX!! - dont hardcode the app_name here like this, but parse Cargo.toml to detect if 
        #         one of the Crate types is "staticlib".
        if p_name_str == "gf_images_jobs":

            #-------------
            # EXTERN_LIB
            target_lib_dir_str = f"{modd_str}/../../rust/build" # f"{modd_str}/../../build/gf_apps/gf_images/tf_lib"

            prepare_libs__extern(target_lib_dir_str,
                p_tf_libs_bool=True,
                p_exit_on_fail_bool=p_exit_on_fail_bool)

            #-------------
            
            
            # RUST_PY - CPYTHON_EXTENSION - this lib is Python extension written in Rust.
            #                               at the moment in GF the convention is for these Rust libs to have a postfix "_py".
            if p_cargo_crate_dir_path_str.endswith("_py"):
                
                # DYNAMIC_LIB
                # IMPORTANT!! - Rust compiles this dynamic lib with the "lib" prefix, but the Python VM
                #               requires extension libs to not have the "lib" prefix.
                source__py_lib_file_path_str = f"{release_dir_str}/lib{p_name_str}_py.so"
                target__py_lib_file_path_str = f"{target_build_dir_path_str}/{p_name_str}_py.so"

                assert os.path.isfile(source__py_lib_file_path_str)
                target_lib_file_path_lst.append((source__py_lib_file_path_str, target__py_lib_file_path_str))

            else:

                # STATIC_LIB
                source__static_lib_file_path_str = f"{release_dir_str}/lib{p_name_str}.a"
                assert os.path.isfile(source__static_lib_file_path_str)
                target_lib_file_path_lst.append((source__static_lib_file_path_str, target_build_dir_path_str))

                # DYNAMIC_LIB
                source__dynamic_lib_file_path_str = f"{release_dir_str}/lib{p_name_str}.so"
                assert os.path.isfile(source__dynamic_lib_file_path_str)
                target_lib_file_path_lst.append((source__dynamic_lib_file_path_str, target_build_dir_path_str))

        #-------------
        # ALL
        else:
            target_lib_file_path_lst.append((f"{release_dir_str}/lib{p_name_str}.so"), target_build_dir_path_str)

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
def prepare_libs__extern(p_target_lib_dir_str,
    p_tf_libs_bool      = False,
    p_exit_on_fail_bool = True):

    #-------------
    # TENSORFLOW
    # IMPORTANT!! - download TF lib and place it in appropriate dir, to have a fresh TF libs
    #               in the build server context, without including the lib in the repo itself.
    if p_tf_libs_bool:

        print(f"{fg('green')}prepare TensorFlow lib{attr(0)}")

        lib_file_name_str  = "tflib.tar.gz"
        tf__version_str    = "1.15.0"
        tf__filename_str   = f"libtensorflow-cpu-linux-x86_64-{tf__version_str}.tar.gz"
        tf__url_str        = f"https://storage.googleapis.com/tensorflow/libtensorflow/{tf__filename_str}"

        gf_core_cli.run(f"mkdir -p {p_target_lib_dir_str}/tf_lib")

        # DOWNLOAD
        _, _, exit_code_int = gf_core_cli.run(f"curl {tf__url_str} --output {lib_file_name_str}")
        if not exit_code_int == 0:
            if p_exit_on_fail_bool:
                exit(exit_code_int)

        # UNPACK
        gf_core_cli.run(f"mv {lib_file_name_str} {p_target_lib_dir_str}/{lib_file_name_str}")
        gf_core_cli.run(f"tar -xvzf {p_target_lib_dir_str}/{lib_file_name_str} -C {p_target_lib_dir_str}/tf_lib")

    #-------------