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
cwd_str = os.path.abspath(os.path.dirname(__file__))

sys.path.append("%s/../meta"%(cwd_str))
import gf_cli_utils

#--------------------------------------------------
# BUILD
def build(p_cargo_crate_dir_path_str,
    p_exit_on_fail_bool = True):
    assert os.path.isdir(p_cargo_crate_dir_path_str)
    print("BUILD...")
    cwd_str = os.getcwd()
    os.chdir(p_cargo_crate_dir_path_str) # change into the target main package dir

    # "rustup update stable"
    c_str = "cargo build --release"

    _, _, exit_code_int = gf_cli_utils.run_cmd(c_str)
    
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
    print("PREPARE LIBS...")


    target_build_dir_path_str = "%s/../../rust/build"%(cwd_str)

    target_lib_file_path_str = None
    if p_type_str == "lib_rust":
        target_lib_file_path_str = "%s/target/release/lib%s.so"%(p_cargo_crate_dir_path_str,
            p_name_str)
    


    c_str = "cp %s %s"%(target_lib_file_path_str, target_build_dir_path_str)
    _, _, exit_code_int = gf_cli_utils.run_cmd(c_str)



    # IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
    if not exit_code_int == 0:
        if p_exit_on_fail_bool:
            exit(exit_code_int)
    


