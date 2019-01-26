# GloFlow media management/publishing system
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

import os
from colored import fg, bg, attr

import gf_cli_utils
#--------------------------------------------------
def run(p_name_str,
    p_go_path_str,
    p_output_path_str,
    p_copy_to_dir_lst):
    assert os.path.isdir(p_go_path_str)
    assert os.path.isdir(os.path.dirname(p_output_path_str))
    assert isinstance(p_copy_to_dir_lst, list)

    print('')
    print(' -- build %s%s%s service'%(fg('green'), p_name_str, attr(0)))

    copy_files(p_copy_to_dir_lst)
    run_go(p_name_str, p_go_path_str, p_output_path_str)
#--------------------------------------------------
def copy_files(p_copy_to_dir_lst):
    print('')
    print('             COPY FILES')
    for src_f_str, target_dir_str in p_copy_to_dir_lst:
        if not os.path.isdir(target_dir_str): gf_cli_utils.run_cmd('mkdir -p %s'%(target_dir_str))
        gf_cli_utils.run_cmd('cp %s %s'%(src_f_str, target_dir_str))
#--------------------------------------------------
def run_go(p_name_str,
    p_go_path_str,
    p_output_path_str):
    
    cwd_str = os.getcwd()
    os.chdir(p_go_path_str) #change into the target main package dir

    gf_cli_utils.run_cmd('go build -o %s'%(p_output_path_str))
    
    os.chdir(cwd_str) #return to initial dir