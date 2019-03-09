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

import os
from colored import fg, bg, attr
import delegator

import gf_cli_utils
#--------------------------------------------------
def run_go(p_name_str,
    p_go_dir_path_str,
    p_output_path_str,
    p_static_bool = False):
    assert isinstance(p_static_bool, bool)

    print(p_go_dir_path_str)
    
    assert os.path.isdir(p_go_dir_path_str)



    print(p_output_path_str)
    assert os.path.isdir(os.path.dirname(p_output_path_str))

    print('')
    if p_static_bool: print(' -- %sSTATIC BINARY BUILD%s'%(fg('yellow'), attr(0)))
    print(' -- build %s%s%s service'%(fg('green'), p_name_str, attr(0)))

    cwd_str = os.getcwd()
    os.chdir(p_go_dir_path_str) #change into the target main package dir

    #STATIC_LINKING - when deploying to containers it is not always guaranteed that all
    #                 required libraries are present. so its safest to compile to a statically
    #                 linked lib.
    #                 build time a few times larger then regular, so slow for dev.
    if p_static_bool:
        args_lst = [
            'CGO_ENABLED=0',
            'GOOS=linux',
            'go build',
            '-ldflags',
            "-s",
            '-a',
            '-installsuffix cgo',
            '-o %s'%(p_output_path_str),
        ]
        c_str = ' '.join(args_lst)
        
    #DYNAMIC_LINKING - fast build for dev.
    else:
        c_str = 'go build -o %s'%(p_output_path_str)

    gf_cli_utils.run_cmd(c_str)
    
    os.chdir(cwd_str) #return to initial dir
#--------------------------------------------------
def list_changed_apps():

    apps_names_lst = [
        'gf_images',
        'gf_analytics',
        'gf_publisher',
        'gf_tagger',
        'gf_landing_page',
    ]

    system_packages_lst = [
        'gf_core',
        'gf_rpc_lib',
        'gf_stats'
    ]

    #latest_commit_hash_str = gf_cli_utils.run_cmd('git rev-parse HEAD')
    #assert len(latest_commit_hash_str) == 32

    list_st = gf_cli_utils.run_cmd('git diff --name-only HEAD HEAD~1', p_print_output_bool=False)

    changed_apps_map = {}

    #--------------------------------------------------
    def add_change_to_all_apps(p_file_changed_str):
        for a in apps_names_lst:
            if changed_apps_map.has_key(a): changed_apps_map[a].append(p_file_changed_str)
            else:                           changed_apps_map[a] = [p_file_changed_str]
    #--------------------------------------------------
    
    for l in list_st.split('\n'):

        #------------------------
        #GO
        if l.startswith('go'):
            #an app itself changed
            if l.startswith('go/gf_apps'):
                app_name_str = l.split('/')[2]
                assert app_name_str in apps_names_lst
                
                if changed_apps_map.has_key(app_name_str): changed_apps_map[app_name_str].append(l)
                else:                                      changed_apps_map[app_name_str] = [l]

            #if one of the system packages has changed
            else:
                for sys_package_str in system_packages_lst:

                    #IMPORTANT!! - one of the system packages has changed, so infer
                    #              that all apps have changed.
                    if l.startswith('go/%s'%(sys_package_str)):
                        add_change_to_all_apps(l)
        #------------------------
        #WEB
        elif l.startswith('web'):
            if l.startswith('web/src/gf_apps'):
                app_name_str = s.split('/')[3]
                assert app_name_str in apps_names_lst

                if changed_apps_map.has_key(app_name_str): changed_apps_map[app_name_str].append(l)
                else:                                      changed_apps_map[app_name_str] = [l]

            #IMPORTATN!! - one of the web libs changed, so all apps should be rebuilt
            #FIX!!       - have a better way of determening which apps use this lib, 
            #              to avoid rebuilding unaffected apps
            elif l.startswith('web/libs'):
                add_change_to_all_apps(l)
        #------------------------

    return changed_apps_map
#--------------------------------------------------
def view_changed_apps(p_changed_apps_map):
    assert isinstance(p_changed_apps_map, dict)

    if len(p_changed_apps_map.items()) == 0:
        print('NO APPS CHANGED')
    else:
        for k, changed_files_lst in p_changed_apps_map.items():
            print('%s%s%s'%(fg('yellow'), k, attr(0)))
            for f in changed_files_lst:
                print('\t%s/%s%s%s'%(os.path.dirname(f), fg('green'), os.path.basename(f), attr(0)))