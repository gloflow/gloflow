import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import delegator

sys.path.append('%s/../../meta'%(cwd_str))
import gf_meta

sys.path.append('%s/../../ops/utils'%(cwd_str))
import gf_build_changes



print("    ---   GF_BUILDER -------------------")
print(delegator.run("ls -al").out)
print(delegator.run("pwd").out)

#--------------------------------------------------
def build_apps():
    apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']
    build_meta_map        = gf_meta.get()['build_info_map']


    print("DIFF")
    changed_apps_map = gf_build_changes.list_changed_apps(apps_changes_deps_map)
    gf_build_changes.view_changed_apps(changed_apps_map)




    if len(changed_apps_map.keys()) == 0:
        exit(0)
    else:



        for app_name_str, v in changed_apps_map.items():

            app_meta_map           = build_meta_map[app_name_str]
            app_go_path_str        = app_meta_map['go_path_str']
            app_go_output_path_str = app_meta_map['go_output_path_str']

            gf_build.run_go(app_name_str,
                app_go_path_str,
                app_go_output_path_str,

                #IMPORTANT!! - binaries are packaged in Alpine Linux, which uses a different standard library then stdlib, 
                #              so all binary dependencies are to be statically linked into the output binary 
                #              without depending on standard dynamic linking.
                p_static_bool = True)

#--------------------------------------------------
build_apps()