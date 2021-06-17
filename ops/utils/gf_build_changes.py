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

sys.path.append("%s/../../py/gf_core"%(modd_str))
import gf_core_cli

#--------------------------------------------------
# p_mark_all_bool - mark all apps as changed. used mainly for debugging.

def list_changed_apps(p_apps_changes_deps_map,
    p_commits_lookback_int = 1,
    p_mark_all_bool        = False):
    assert isinstance(p_apps_changes_deps_map, dict)
    assert "apps_gf_packages_map" in p_apps_changes_deps_map.keys()
    assert "system_packages_lst" in p_apps_changes_deps_map.keys()
    assert isinstance(p_commits_lookback_int, int)

    apps_gf_packages_map = p_apps_changes_deps_map['apps_gf_packages_map']
    system_packages_lst  = p_apps_changes_deps_map['system_packages_lst']
    assert isinstance(system_packages_lst, list)

    changed_apps_files_map  = {
        # IMPORTANT!! - these are all apps that have either "go" or "web" changed. this is 
        #               needed because when building in CI
        #               even if only Go files changed we need Web code built as well so that
        #               the final container can be built in its full form.
        "all": {},
        "go":  {},
        "web": {},
    }

    #------------------------
    # DEBUGGING - mark all apps as changed
    if p_mark_all_bool:
        for a, _ in apps_gf_packages_map.items():
            changed_apps_files_map["all"][a] = ["all"]
            changed_apps_files_map["go"][a]  = ["all"]
            changed_apps_files_map["web"][a] = ["all"]
        return changed_apps_files_map

    #------------------------

    # latest_commit_hash_str, _ = gf_core_cli.run('git rev-parse HEAD')
    # assert len(latest_commit_hash_str) == 32

    #------------------------
    # FIX!! - dont just look 1 commit back to see what changed. if localy a developer makes several commits and then uploads code
    #         got github(or other) and CI clones it this function then might miss some of the services/apps/packages that changed 
    #         several commits back.
    #         instead some mechanism for getting the number of commits that some deployment environment is behind HEAD,
    #         and then use that number for this "p_commits_lookback_int" argument (that would be >1).
    past_commit_str = "HEAD~%s"%(p_commits_lookback_int)

    #------------------------

    list_str, _, _ = gf_core_cli.run("git diff --name-only HEAD %s"%(past_commit_str), p_print_output_bool=False)

    #--------------------------------------------------
    # IMPORTANT!! - the file that changed affects all apps, so they all need to be marked as changed
    #               and this file added to the list of changed files of all apps.

    def add_change_to_all_apps(p_file_changed_str, p_type_str):
        assert p_type_str == "go" or p_type_str == "web"
        for a, _ in apps_gf_packages_map.items():

            if a in changed_apps_files_map[p_type_str].keys():
                changed_apps_files_map["all"][a].append(p_file_changed_str)
                changed_apps_files_map[p_type_str][a].append(p_file_changed_str)
            else:
                changed_apps_files_map["all"][a].append(p_file_changed_str)
                changed_apps_files_map[p_type_str][a] = [p_file_changed_str]

    #--------------------------------------------------
    # IMPORTANT!! - update only the apps that have this files package is marked as a dependancy of
    def update_dependant_apps_file_lists(p_package_name_str, p_file_path_str, p_type_str):
        assert p_type_str == "go" or p_type_str == "web"

        # build out a list of apps that this package (p_package_name_str) is a dependency of
        dependant_apps_lst = []
        for app_str, app_gf_package_lst in p_apps_changes_deps_map['apps_gf_packages_map'].items():
            if p_package_name_str in app_gf_package_lst:
                dependant_apps_lst.append(app_str)

        # for all apps that are determined to have changed (because they depend on p_package_name_str package) 
        # add this file (p_file_path_str) to those apps lists of changed files.
        for app_str in dependant_apps_lst:
            
            if app_str in changed_apps_files_map.keys():
                changed_apps_files_map["all"][app_str].append(p_file_changed_str)
                changed_apps_files_map[p_type_str][app_str].append(p_file_changed_str)
            else:  
                changed_apps_files_map["all"][app_str].append(p_file_changed_str)                                     
                changed_apps_files_map[p_type_str][app_str] = [p_file_changed_str]

    #--------------------------------------------------
    
    for l in list_str.split('\n'):

        #------------------------
        # GO
        if l.startswith('go'):
            # an app changed
            if l.startswith('go/gf_apps'):
                package_name_str = l.split('/')[2] # third element in the file path is a package name
                update_dependant_apps_file_lists(package_name_str, l, "go")

            # one of the system packages has changed
            else:
                for sys_package_str in system_packages_lst:

                    # IMPORTANT!! - one of the system packages has changed, so infer
                    #              that all apps have changed.
                    if l.startswith('go/%s'%(sys_package_str)):
                        add_change_to_all_apps(l, "go")

        #------------------------
        # WEB
        elif l.startswith('web'):
            if l.startswith('web/src/gf_apps'):
                package_name_str = l.split('/')[3] # get package_name from the path of the changed file
                update_dependant_apps_file_lists(package_name_str, l, "web")

            # IMPORTATN!! - one of the web libs changed, so all apps should be rebuilt
            # FIX!!       - have a better way of determening which apps use this lib, 
            #               to avoid rebuilding unaffected apps
            elif l.startswith('web/libs'):
                add_change_to_all_apps(l, "web")

            else:
                package_name_str = l.split('/')[2]
                if package_name_str in system_packages_lst:
                    add_change_to_all_apps(l, web)

        #------------------------

    return changed_apps_files_map

#--------------------------------------------------
def view_changed_apps(p_changed_apps_files_map, p_type_str):
    assert isinstance(p_changed_apps_files_map, dict)
    assert p_type_str == "go" or p_type_str == "web"

    print("----- %s"%(p_type_str))

    if len(p_changed_apps_files_map[p_type_str].items()) == 0:
        print('NO APPS CHANGED')
    else:
        for app_name_str, changed_files_lst in p_changed_apps_files_map[p_type_str].items():
            print('%s%s%s'%(fg('yellow'), app_name_str, attr(0)))
            
            if len(changed_files_lst) == 0:
                print('%s%s%s'%(fg('green'),
                    "ALL",
                    attr(0)))

            for f in changed_files_lst:

                print('\t%s/%s%s%s'%(os.path.dirname(f),
                    fg('green'),
                    os.path.basename(f),
                    attr(0)))