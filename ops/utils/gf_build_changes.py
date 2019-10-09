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
def list_changed_apps(p_apps_changes_deps_map):
    assert isinstance(p_apps_changes_deps_map, dict)
    assert p_apps_changes_deps_map.has_key('apps_names_map')
    assert p_apps_changes_deps_map.has_key('system_packages_lst')

    apps_names_map      = p_apps_changes_deps_map['apps_names_map']
    system_packages_lst = p_apps_changes_deps_map['system_packages_lst']
    assert isinstance(system_packages_lst, list)

    #latest_commit_hash_str = gf_cli_utils.run_cmd('git rev-parse HEAD')
    #assert len(latest_commit_hash_str) == 32

    list_st = gf_cli_utils.run_cmd('git diff --name-only HEAD HEAD~1', p_print_output_bool=False)

    changed_apps_map = {}
    #--------------------------------------------------
    def add_change_to_all_apps(p_file_changed_str):
        for a, _ in apps_names_map.items():
            if changed_apps_map.has_key(a): changed_apps_map[a].append(p_file_changed_str)
            else:                           changed_apps_map[a] = [p_file_changed_str]

    #--------------------------------------------------
    def update_dependant_apps(p_package_name_str, p_file_path_str):
        dependant_apps_lst = []
        for app_str, package_deps_lst in p_apps_changes_deps_map['apps_names_map'].items():
            if p_package_name_str in package_deps_lst:
                dependant_apps_lst.append(app_str)
                
        
        for app_str in dependant_apps_lst:
            if changed_apps_map.has_key(app_str): changed_apps_map[app_str].append(p_file_path_str)
            else:                                 changed_apps_map[app_str] = [p_file_path_str]

    #--------------------------------------------------
    
    for l in list_st.split('\n'):

        #------------------------
        #GO
        if l.startswith('go'):
            #an app itself changed
            if l.startswith('go/gf_apps'):
                package_name_str = l.split('/')[2]
                update_dependant_apps(package_name_str, l)

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
                package_name_str = l.split('/')[3]
                update_dependant_apps(package_name_str, l)

            #IMPORTATN!! - one of the web libs changed, so all apps should be rebuilt
            #FIX!!       - have a better way of determening which apps use this lib, 
            #              to avoid rebuilding unaffected apps
            elif l.startswith('web/libs'):
                add_change_to_all_apps(l)

            else:
                package_name_str = l.split('/')[2]
                if package_name_str in system_packages_lst:
                    add_change_to_all_apps(l)
        #------------------------

    return changed_apps_map

#--------------------------------------------------
def view_changed_apps(p_changed_apps_map):
    assert isinstance(p_changed_apps_map, dict)

    if len(p_changed_apps_map.items()) == 0:
        print('NO APPS CHANGED')
    else:
        for app_name_str, changed_files_lst in p_changed_apps_map.items():
            print('%s%s%s'%(fg('yellow'), app_name_str, attr(0)))
            for f in changed_files_lst:
                print('\t%s/%s%s%s'%(os.path.dirname(f), fg('green'), os.path.basename(f), attr(0)))