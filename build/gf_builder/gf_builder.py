import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import argparse
from colored import fg, bg, attr
import delegator

sys.path.append('%s/../../meta'%(cwd_str))
import gf_meta
import gf_web_meta

sys.path.append('%s/../../ops/utils'%(cwd_str))
import gf_build_changes
import gf_build
import gf_log

sys.path.append('%s/../../ops/tests'%(cwd_str))
import gf_tests

sys.path.append('%s/../../ops/web'%(cwd_str))
import gf_web__build

sys.path.append('%s/../../ops/aws'%(cwd_str))
import gf_aws_creds

#--------------------------------------------------
def main():
    
    print("    ---   GF_BUILDER -------------------")
    print(delegator.run("ls -al").out)
    print("pwd[%s] - whoami[%s]"%(delegator.run("pwd").out.strip(), delegator.run("whoami").out.strip()))

    args_map = parse_args()

    #GET_CHANGED_APPS
    changed_apps_files_map = get_changed_apps()

    #------------------------
    #TEST
    if args_map["run"] == "test":
        test_apps(changed_apps_files_map)

    #------------------------
    #BUILD
    elif args_map["run"] == "build":

        #IMPORTANT!! - only insert Git commit hash if gf_builder.py is run in CI
        if "DRONE_COMMIT_SHA" in os.environ:
            git_commit_hash_str = os.environ["DRONE_COMMIT_SHA"]
            paste_git_commit_hash(git_commit_hash_str)
            
        build_apps(changed_apps_files_map)

    #------------------------

#--------------------------------------------------
def test_apps(p_changed_apps_files_map):
    assert isinstance(p_changed_apps_files_map, dict)

    print("\n\n TEST APPS ----------------------------------------------------- \n\n")

    build_meta_map        = gf_meta.get()['build_info_map']
    apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']
    apps_gf_packages_map  = apps_changes_deps_map["apps_gf_packages_map"]

    #AWS_CREDS
    aws_creds_map = gf_aws_creds.get_from_env_vars()
    assert isinstance(aws_creds_map, dict)

    #nothing changed
    if len(p_changed_apps_files_map.keys()) == 0:
        return
    else:

        #------------------------
        # GO
        print("\nGO--------\n")
        for app_name_str, v in p_changed_apps_files_map["go"].items():
            
            test_name_str = "all"

            #IMPORTANT!! - get all packages that are involved in tis app, so that 
            #              tests for all these packages can be run.
            app_gf_packages_lst = apps_gf_packages_map[app_name_str]

            #RUN_TESTS_FOR_ALL_APP_PACKAGES
            for app_gf_package_name_str in app_gf_packages_lst:

                assert build_meta_map.has_key(app_gf_package_name_str)
                gf_package_meta_map  = build_meta_map[app_gf_package_name_str]

                gf_tests.run(app_gf_package_name_str,
                    test_name_str,
                    gf_package_meta_map,
                    aws_creds_map,

                    #IMPORTANT!! - in case the tests that gf_test.run() executes fail, 
                    #              run() should call exit() and force this whole process to exit, 
                    #              so that CI marks the build as failed.
                    p_exit_on_fail_bool = True)
        #------------------------
    
#--------------------------------------------------
def build_apps(p_changed_apps_files_map):
    assert isinstance(p_changed_apps_files_map, dict)

    print("\n\n BUILD APPS ----------------------------------------------------- \n\n")

    build_meta_map = gf_meta.get()['build_info_map']
    
    #nothing changed
    if len(p_changed_apps_files_map.keys()) == 0:
        return
    else:
        #------------------------
        # WEB
        print("\n\nWEB--------\n\n")
        web_meta_map   = gf_web_meta.get()
        apps_names_lst = []
        for app_name_str, v in p_changed_apps_files_map["web"].items():
            apps_names_lst.append(app_name_str)

        gf_web__build.build(apps_names_lst, web_meta_map, gf_log.log_fun)
        #------------------------
        # GO
        print("\n\nGO--------\n\n")
        for app_name_str, v in p_changed_apps_files_map["go"].items():

            app_meta_map           = build_meta_map[app_name_str]
            app_go_path_str        = app_meta_map['go_path_str']
            app_go_output_path_str = app_meta_map['go_output_path_str']

            gf_build.run_go(app_name_str,
                app_go_path_str,
                app_go_output_path_str,

                # IMPORTANT!! - binaries are packaged in Alpine Linux, which uses a different standard library then stdlib, 
                #               so all binary dependencies are to be statically linked into the output binary 
                #               without depending on standard dynamic linking.
                p_static_bool = True, 
                
                # gf_build.run_go() should exit if the "go build" CLI run returns with a non-zero exit code.
                # gf_builder.py is meant to run in CI environments, and so we want the stage in which it runs 
                # to be marked as failed because of the non-zero exit code.
                p_exit_on_fail_bool = True)
        #------------------------

#--------------------------------------------------
def get_changed_apps():
    print("DIFF")
    apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']

    # LIST_CHANGED_APPS - determine how which apps/services changed
    changed_apps_files_map = gf_build_changes.list_changed_apps(apps_changes_deps_map,
        p_commits_lookback_int = 1, 
        p_mark_all_bool        = True)

    # VIEW
    gf_build_changes.view_changed_apps(changed_apps_files_map, "go")
    gf_build_changes.view_changed_apps(changed_apps_files_map, "web")
    return changed_apps_files_map

#--------------------------------------------------
def paste_git_commit_hash(p_git_commit_hash_str):
    print("PASTE_GIT_COMMIT_HASH - %s"%(p_git_commit_hash_str))

    golang_sys_release_info_file_path_str = "%s/../../go/gf_core/gf_sys_release_info.go"%(cwd_str)
    assert os.path.isfile(golang_sys_release_info_file_path_str)
    
    original_word_regex_str = 'Git_commit_str: "",' #this is the original line of Go code
    new_word_regex_str      = 'Git_commit_str: "%s",'%(p_git_commit_hash_str)

    #------------------------
    # IMPORTANT!! - "sed" - Stream EDitor.
    #               "-i" - in-place, save to original file
    #               command string:
    #                 "s" - the substitute command
    #                 "g" - global, replace all not just first instance
    c = "sed -i 's/%s/%s/g' %s"%(original_word_regex_str, new_word_regex_str, golang_sys_release_info_file_path_str)
    print(c)
    #------------------------

    r = delegator.run(c)
    print(r.out)
    print(r.err)

#--------------------------------------------------
# IMPORTANT!! - get Git commit from the deployed artifact (making API call to a target service).
#               this is needed to know how far HEAD of this monorepo is ahead from the commit 
#               that was used to build a particular service, to then use that integer distance
#               as the p_commits_lookback_int when determening which apps change when calling 
#               list_changed_apps(). 

def get_deployed_commit(p_domain_str = "https://gloflow.com"):
    True

#--------------------------------------------------
def parse_args():
    arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
    arg_parser.add_argument('-run', action = "store", default = 'build',
        help = '''
- '''+fg('yellow')+'build'+attr(0)+'''            - build app golang/web code
- '''+fg('yellow')+'build_containers'+attr(0)+''' - build app Docker containers
- '''+fg('yellow')+'test'+attr(0)+'''             - run app code tests
        ''')

    cli_args_lst   = sys.argv[1:]
    args_namespace = arg_parser.parse_args(cli_args_lst)
    return {
        "run": args_namespace.run,
    }

#--------------------------------------------------
main()