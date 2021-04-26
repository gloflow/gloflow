# GloFlow application and media management/publishing platform
# Copyright (C) 2019 Ivan Trajkovic
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
import delegator

import gf_cli_utils

sys.path.append("%s/../../py/gf_core"%(modd_str))
import gf_core_cli

#--------------------------------------------------
# RUN_IN_CONTAINER
def run_in_cont(p_name_str,
    p_go_dir_path_str,
    p_go_output_path_str,
    p_static_bool = False):


    repo_local_path_str = os.path.abspath(f'{modd_str}/../../../gloflow').strip()
    cmd_lst = [
        "sudo", "docker", "run",
        "--rm", # remove after exit 
        "-v", f"{repo_local_path_str}:/home/gf", # mount repo into the container
        "glofloworg/gf_builder_go_ubuntu:latest",
        
        # "python3 -u", "/home/gf/build/gf_builder/gf_builder.py", "-run=build_go"
        "python3", "-u", "/home/gf/ops/cli__build.py", "-run=build_go", "-build_outof_cont", f"-app={p_name_str}"
    ]
    p = gf_core_cli.run__view_realtime(cmd_lst, {

        },
        "gf_build_go", "green")

    p.wait()

#--------------------------------------------------
# RUN
def run(p_name_str,
    p_go_dir_path_str,
    p_go_output_path_str,
    p_static_bool       = False,
    p_exit_on_fail_bool = True,
    p_dynamic_libs_dir_path_str = os.path.abspath("%s/../../rust/build"%(modd_str)),
    p_go_get_bool = True):
    assert isinstance(p_static_bool, bool)
    
    print("")
    if p_static_bool:
        print(" -- %sSTATIC BINARY BUILD%s"%(fg("yellow"), attr(0)))
        
    print(" -- build %s%s%s service"%(fg("green"), p_name_str, attr(0)))
    print(" -- go_dir_path    - %s%s%s"%(fg("green"), p_go_dir_path_str, attr(0)))  
    print(" -- go_output_path - %s%s%s"%(fg("green"), p_go_output_path_str, attr(0)))  

    assert os.path.isdir(p_go_dir_path_str)
    assert os.path.isdir(os.path.dirname(p_go_output_path_str))

    print("Go cache dir:")
    gf_cli_utils.run_cmd(f"go env GOCACHE")

    cwd_str = os.getcwd()
    os.chdir(p_go_dir_path_str) # change into the target main package dir



    # RUST_DYNAMIC_LIBS
    dynamic_libs_dir_path_str    = os.path.abspath(f"{modd_str}/../../rust/build")
    tf_dynamic_libs_dir_path_str = os.path.abspath(f"{modd_str}/../../rust/build/tf_lib/lib")

    print(f"dynamic libs dir - {fg('green')}{dynamic_libs_dir_path_str}{attr(0)}")
    gf_cli_utils.run_cmd(f"ls -al {dynamic_libs_dir_path_str}")


    LD_paths_lst = [
        f"LD_LIBRARY_PATH={dynamic_libs_dir_path_str}",
        f"LD_LIBRARY_PATH={tf_dynamic_libs_dir_path_str}"
    ]
    
    # GO_GET
    if p_go_get_bool:
        _, _, exit_code_int = gf_cli_utils.run_cmd(f"{' '.join(LD_paths_lst)} go get -u")
        print("")
        print("")

    

    #-----------------------------
    # STATIC_LINKING - when deploying to containers it is not always guaranteed that all
    #                  required libraries are present. so its safest to compile to a statically
    #                  linked lib.
    #                  build time a few times larger then regular, so slow for dev.
    # "-ldflags '-s'" - omit the symbol table and debug information.

    if p_static_bool:
        
        print(f"{fg('yellow')}STATIC LINKING{attr(0)} --")
        
        # https://golang.org/cmd/link/
        # IMPORTANT!! - "CGO_ENABLED=0" and "-installsuffix cgo" no longer necessary since golang 1.10.
        #               "CGO_ENABLED=0" we also dont want to disable since Rust libs are used in Go via CGO.
        
        # IMPORTANT!! - debug .a files:
        #   "ar -t libgf_images_jobs.a" - get a list of Archived object files in static .a libs.
        #                                 static library is an archive (ar) of object files.
        #                                 The object files are usually in the ELF format

        gf_cli_utils.run_cmd(f"ldconfig -v")
        # gf_cli_utils.run_cmd(f"cp {dynamic_libs_dir_path_str}/libgf_images_jobs.a /usr/lib")

        args_lst = [
            
            ' '.join(LD_paths_lst),
            # f"LD_LIBRARY_PATH={dynamic_libs_dir_path_str}",
            # f"LD_LIBRARY_PATH=/usr/lib",

            # "CGO_ENABLED=0",
            "GOOS=linux",
            "go build",

            # force rebuilding of packages that are already up-to-date.
            "-a",

            # "-installsuffix cgo",

            # LINKER_FLAGS
            # "-ldflags"    - arguments to pass on each go tool link invocation
            # "-s"          - Omit the symbol table and debug information
            # "-extldflags" - Set space-separated flags to pass to the external linker.
            #                 on Alpine builds the GCC toolchain linker "ld" is used.
            # "-static"     - On systems that support dynamic linking, this 
            #                 overrides -pie and prevents linking with the shared libraries.
            # "-ldl"        - "-l" provides lib path. links in  /usr/lib/libdl.so/.a
            #                 this is needed to prevent Rust .a lib errors relating
            #                 to undefined references to "dlsym","dladdr"
            #
            # (f'''-ldflags '-s -extldflags "-t -static -lgf_images_jobs -ldl -lglib"' ''').strip(),
            # (f'''-ldflags '-s -extldflags "-lm"' ''').strip(),
            ('''-ldflags '-s -extldflags "-static -ldl"' ''').strip(),
            
            "-o %s"%(p_go_output_path_str),
        ]
        c_str = " ".join(args_lst)
    
    #-----------------------------
    # DYNAMIC_LINKING - fast build for dev.
    else:
        print(f"{fg('yellow')}DYNAMIC LINKING{attr(0)} --")

        c_str = f"{' '.join(LD_paths_lst)} go build -o {p_go_output_path_str}"

    #-----------------------------
    
    print(c_str)
    _, _, exit_code_int = gf_cli_utils.run_cmd(c_str)

    # IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
    if not exit_code_int == 0:
        if p_exit_on_fail_bool:
            exit(exit_code_int)

    os.chdir(cwd_str) # return to initial dir

    print("build done...")