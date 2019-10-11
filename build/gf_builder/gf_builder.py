import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import delegator

sys.path.append('%s/../../meta'%(cwd_str))
import gf_meta

sys.path.append('%s/../../ops/utils'%(cwd_str))
import gf_build_changes
import gf_build



print("    ---   GF_BUILDER -------------------")
print(delegator.run("ls -al").out)
print(delegator.run("pwd").out)
print(delegator.run("whoami").out)
#--------------------------------------------------
def build_apps():
    apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']
    build_meta_map        = gf_meta.get()['build_info_map']


    #------------------------
    print("DIFF")

    #LIST_CHANGED_APPS - determine how which apps/services changed
    changed_apps_map = gf_build_changes.list_changed_apps(apps_changes_deps_map,
        p_commits_lookback_int = 1, 
        p_mark_all_bool        = True)

    #VIEW
    gf_build_changes.view_changed_apps(changed_apps_map, "go")
    gf_build_changes.view_changed_apps(changed_apps_map, "web")
    #------------------------


    if len(changed_apps_map.keys()) == 0:
        exit(0)
    else:

        #GO
        for app_name_str, v in changed_apps_map["go"].items():

            app_meta_map           = build_meta_map[app_name_str]
            app_go_path_str        = app_meta_map['go_path_str']
            app_go_output_path_str = app_meta_map['go_output_path_str']

            gf_build.run_go(app_name_str,
                app_go_path_str,
                app_go_output_path_str,

                #IMPORTANT!! - binaries are packaged in Alpine Linux, which uses a different standard library then stdlib, 
                #              so all binary dependencies are to be statically linked into the output binary 
                #              without depending on standard dynamic linking.
                p_static_bool = True, 
                
                #gf_build.run_go() should exit if the "go build" CLI run returns with a non-zero exit code.
                #gf_builder.py is meant to run in CI environments, and so we want the stage in which it runs 
                #to be marked as failed because of the non-zero exit code.
                p_exit_on_fail_bool = True)


#--------------------------------------------------
#IMPORTANT!! - get Git commit from the deployed artifact (making API call to a target service).
#              this is needed to know how far HEAD of this monorepo is ahead from the commit 
#              that was used to build a particular service, to then use that integer distance
#              as the p_commits_lookback_int when determening which apps change when calling 
#              list_changed_apps(). 

def get_deployed_commit(p_domain_str = "https://gloflow.com"):
    True

#--------------------------------------------------
build_apps()