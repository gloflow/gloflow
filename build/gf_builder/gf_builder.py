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

sys.path.append('%s/../../ops/containers'%(cwd_str))
import gf_containers

#--------------------------------------------------
def main():
    
    print("    ---   GF_BUILDER -------------------")
    print(delegator.run("ls -al").out)
    print("pwd[%s] - whoami[%s]"%(delegator.run("pwd").out.strip(), delegator.run("whoami").out.strip()))

    args_map = parse_args()

    #GET_CHANGED_APPS
    changed_apps_files_map = get_changed_apps()

    #------------------------
    # TEST
    if args_map["run"] == "test":
        test_apps(changed_apps_files_map)

    #------------------------
    # BUILD
    elif args_map["run"] == "build":

        #IMPORTANT!! - only insert Git commit hash if gf_builder.py is run in CI
        if not args_map["drone_commit_sha"] == None:
            git_commit_hash_str = args_map["drone_commit_sha"]
            paste_git_commit_hash(git_commit_hash_str)
            
        build_apps(changed_apps_files_map)

    #------------------------
    # BUILD_CONTAINERS
    elif args_map["run"] == "build_containers":

        gf_dockerhub_user_str = args_map["gf_dockerhub_user"]
        
        build_apps_containers(changed_apps_files_map,
            gf_dockerhub_user_str)

    #------------------------
    # PUBLISH_CONTAINERS
    elif args_map["run"] == "publish_containers":

        gf_dockerhub_user_str = args_map["gf_dockerhub_user"]
        gf_dockerhub_pass_str = args_map["gf_dockerhub_pass"]

        assert not gf_dockerhub_user_str == None
        assert not gf_dockerhub_pass_str == None

        publish_apps_containers(changed_apps_files_map,
            gf_dockerhub_user_str,
            gf_dockerhub_pass_str)

    #------------------------

#--------------------------------------------------
# PUBLISH_APPS_CONTAINERS
def publish_apps_containers(p_changed_apps_files_map,
    p_gf_dockerhub_user_str,
    p_gf_dockerhub_pass_str):
    assert isinstance(p_gf_dockerhub_user_str, basestring)
    assert isinstance(p_gf_dockerhub_pass_str, basestring)
    
    build_meta_map = gf_meta.get()['build_info_map']
    
    # "all" - this key holds a map with all the apps that had either their Go or Web code changed
    apps_names_lst = p_changed_apps_files_map["all"].keys()

    for app_name_str in apps_names_lst:

        assert build_meta_map.has_key(app_name_str)
        app_build_meta_map = build_meta_map[app_name_str]

        # PUBLISH
        gf_containers.publish(app_name_str,
            app_build_meta_map,
            p_gf_dockerhub_user_str,
            p_gf_dockerhub_pass_str,
            gf_log.log_fun,
            p_exit_on_fail_bool = True)

#--------------------------------------------------
# BUILD_APPS_CONTAINERS
def build_apps_containers(p_changed_apps_files_map,
    p_gf_dockerhub_user_str):
    assert isinstance(p_changed_apps_files_map, dict)

    build_meta_map = gf_meta.get()['build_info_map']
    web_meta_map   = gf_web_meta.get()

    # IMPORTANT!! - for each app that has any of its code changed rebuild both the Go and Web code,
    #               since the containers has to be fully rebuilt.
    # "all" - this key holds a map with all the apps that had either their Go or Web code changed
    apps_names_lst = p_changed_apps_files_map["all"].keys()


    for app_name_str in apps_names_lst:

        assert build_meta_map.has_key(app_name_str)
        app_build_meta_map = build_meta_map[app_name_str]

        assert web_meta_map.has_key(app_name_str)
        app_web_meta_map = web_meta_map[app_name_str]
        
        gf_containers.build(app_name_str,
            app_build_meta_map,
            app_web_meta_map,
            gf_log.log_fun,

            # DOCKERHUB_USER
            p_user_name_str = p_gf_dockerhub_user_str,

            # gf_containers.build() should exit if the docker build CLI run returns with a non-zero exit code.
            # gf_builder.py is meant to run in CI environments, and so we want the stage in which it runs 
            # to be marked as failed because of the non-zero exit code.
            p_exit_on_fail_bool = True)

#--------------------------------------------------
def test_apps(p_changed_apps_files_map):
    assert isinstance(p_changed_apps_files_map, dict)

    print("\n\n TEST APPS ----------------------------------------------------- \n\n")

    build_meta_map        = gf_meta.get()['build_info_map']
    apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']
    apps_gf_packages_map  = apps_changes_deps_map["apps_gf_packages_map"]

    # AWS_CREDS
    aws_creds_map = gf_aws_creds.get_from_env_vars()
    assert isinstance(aws_creds_map, dict)

    # nothing changed
    if len(p_changed_apps_files_map.keys()) == 0:
        return
    else:

        #------------------------
        # GO
        print("\nGO--------\n")
        for app_name_str, v in p_changed_apps_files_map["go"].items():
            
            test_name_str = "all"

            # IMPORTANT!! - get all packages that are dependencies of this app, so that 
            #               tests for all these packages can be run. dont just run the tests of the app
            #               but of all the packages that are its dependencies as well.
            app_gf_packages_lst = apps_gf_packages_map[app_name_str]

            # RUN_TESTS_FOR_ALL_APP_PACKAGES
            for app_gf_package_name_str in app_gf_packages_lst:

                print("about to run tests for - %s"%(app_gf_package_name_str))
                if build_meta_map.has_key(app_gf_package_name_str):
                    gf_package_meta_map = build_meta_map[app_gf_package_name_str]

                    gf_tests.run(app_gf_package_name_str,
                        test_name_str,
                        gf_package_meta_map,
                        aws_creds_map,

                        # IMPORTANT!! - in case the tests that gf_test.run() executes fail, 
                        #               run() should call exit() and force this whole process to exit, 
                        #               so that CI marks the build as failed.
                        p_exit_on_fail_bool = True)
        #------------------------
    
#--------------------------------------------------
# BUILD_APPS
def build_apps(p_changed_apps_files_map):
    assert isinstance(p_changed_apps_files_map, dict)

    print("\n\n BUILD APPS ----------------------------------------------------- \n\n")

    build_meta_map = gf_meta.get()['build_info_map']
    
    # nothing changed
    if len(p_changed_apps_files_map.keys()) == 0:
        return
    else:

        # IMPORTANT!! - for each app that has any of its code changed rebuild both the Go and Web code,
        #               since the containers has to be fully rebuilt.
        # "all" - this key holds a map with all the apps that had either their Go or Web code changed
        apps_names_lst = p_changed_apps_files_map["all"].keys()

        #------------------------
        # WEB
        print("\n\nWEB--------\n\n")
            
        web_meta_map = gf_web_meta.get()

        gf_web__build.build(apps_names_lst,
            web_meta_map,
            gf_log.log_fun)
        
        #------------------------
        # GO
        print("\n\nGO--------\n\n")
        
        for app_name_str in apps_names_lst:

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
# GET_CHANGED_APPS
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
# PASTE_GIT_COMMIT_HASH
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

    #-------------
    # RUN
    arg_parser.add_argument('-run', action = "store", default = 'build',
        help = '''
- '''+fg('yellow')+'test'+attr(0)+'''               - run app code tests
- '''+fg('yellow')+'build'+attr(0)+'''              - build app golang/web code
- '''+fg('yellow')+'build_containers'+attr(0)+'''   - build app Docker containers
- '''+fg('yellow')+'publish_containers'+attr(0)+''' - publish app Docker containers
        ''')

    #-------------
    # # MONGODB_HOST
    # arg_parser.add_argument('-mongodb_host',
    #     action =  "store",
    #     default = "all",
    #     help =    '''mongodb host to connect to (for testing, etc.)''')

    #-------------
    # ENV_VARS
    drone_commit_sha_str  = os.environ.get("DRONE_COMMIT_SHA", None) # Drone defined ENV var
    gf_dockerhub_user_str = os.environ.get("GF_DOCKERHUB_USER", None)
    gf_dockerhub_pass_str = os.environ.get("GF_DOCKERHUB_P", None)

    print(gf_dockerhub_user_str)
    print(gf_dockerhub_pass_str)
    #-------------
    cli_args_lst   = sys.argv[1:]
    args_namespace = arg_parser.parse_args(cli_args_lst)
    return {
        "run":               args_namespace.run,
        "drone_commit_sha":  drone_commit_sha_str,
        "gf_dockerhub_user": gf_dockerhub_user_str,
        "gf_dockerhub_pass": gf_dockerhub_pass_str,
        # "mongodb_host": args_namespace.mongodb_host,
    }

#--------------------------------------------------
main()