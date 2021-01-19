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

#--------------------------------------------------
# RUN_GO
def run_go(p_name_str,
    p_go_dir_path_str,
    p_go_output_path_str,
    p_static_bool       = False,
    p_exit_on_fail_bool = True,
    p_dynamic_libs_dir_path_str = os.path.abspath("%s/../../rust/build"%(modd_str))):
    assert isinstance(p_static_bool, bool)
    assert os.path.isdir(p_go_dir_path_str)

    print(p_go_output_path_str)
    assert os.path.isdir(os.path.dirname(p_go_output_path_str))

    print("")
    if p_static_bool:
        print(" -- %sSTATIC BINARY BUILD%s"%(fg("yellow"), attr(0)))
        
    print(" -- build %s%s%s service"%(fg("green"), p_name_str, attr(0)))
    print(" -- go_dir_path    - %s%s%s"%(fg("green"), p_go_dir_path_str, attr(0)))  
    print(" -- go_output_path - %s%s%s"%(fg("green"), p_go_output_path_str, attr(0)))  

    cwd_str = os.getcwd()
    os.chdir(p_go_dir_path_str) # change into the target main package dir

    # GO_GET
    _, _, exit_code_int = gf_cli_utils.run_cmd("go get -u")
    print("")
    print("")

    # STATIC_LINKING - when deploying to containers it is not always guaranteed that all
    #                  required libraries are present. so its safest to compile to a statically
    #                  linked lib.
    #                  build time a few times larger then regular, so slow for dev.
    # "-ldflags '-s'" - omit the symbol table and debug information.
    # "-a" - forces all packages to be rebuilt
    if p_static_bool:
        
        # https://golang.org/cmd/link/
        # IMPORTANT!! - "CGO_ENABLED=0" and "-installsuffix cgo" no longer necessary since golang 1.10.
        #               "CGO_ENABLED=0" we also dont want to disable since Rust libs are used in Go via CGO.
        # "-extldflags flags" - Set space-separated flags to pass to the external linker
        args_lst = [
            # "CGO_ENABLED=0",
            "GOOS=linux",
            "go build",
            "-a",
            # "-installsuffix cgo",

            # LINKER_FLAGS
            # "-ldl" - "-l" provides lib path. links in  /usr/lib/libdl.so/.a
            #          this is needed to prevent Rust .a lib errors relating
            #          to undefined references to "dlsym","dladdr"
            ('''-ldflags '-s -extldflags "-static -ldl"' ''').strip(),
            
            
            "-o %s"%(p_go_output_path_str),
        ]
        c_str = " ".join(args_lst)
        
    # DYNAMIC_LINKING - fast build for dev.
    else:
        c_str = "go build -o %s"%(p_go_output_path_str)

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