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

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import argparse

from colored import fg, bg, attr
import delegator

sys.path.append('%s/../meta'%(cwd_str))
import gf_meta
import gf_web_meta

sys.path.append('%s/utils'%(cwd_str))
import gf_build
import gf_build_changes

sys.path.append('%s/tests'%(cwd_str))
import gf_tests

sys.path.append('%s/aws/s3'%(cwd_str))
import gf_s3_utils

sys.path.append('%s/web'%(cwd_str))
import gf_web__build

sys.path.append('%s/containers'%(cwd_str))
import gf_containers
#--------------------------------------------------
def main():
    
    print('')
    print('                              %sGLOFLOW BUILD TOOL%s'%(fg('green'),attr(0)))
    print('')

    #--------------------------------------------------
    def log_fun(g, m):
        if g == "ERROR":
            print('%s%s%s:%s%s%s'%(bg('red'), g, attr(0), fg('red'), m, attr(0)))
        else:
            print('%s%s%s:%s%s%s'%(fg('yellow'), g, attr(0), fg('green'), m, attr(0)))
    #--------------------------------------------------
    
    build_meta_map        = gf_meta.get()['build_info_map']
    apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']
    args_map   = parse_args()
    run_str    = args_map['run']

    app_name_str = args_map['app']
    assert build_meta_map.has_key(app_name_str)

    #--------------------------------------------------
    def go_build(p_static_bool):
        app_meta_map = build_meta_map[app_name_str]
        if not app_meta_map.has_key('go_output_path_str'):
            print("not a main package")
            exit()
            
        gf_build.run_go(app_name_str,
            app_meta_map['go_path_str'],
            app_meta_map['go_output_path_str'],
            p_static_bool = p_static_bool)
    #--------------------------------------------------

    #-------------
    #BUILD
    if run_str == 'build':
        
        #build using dynamic linking, its quicker while in dev.
        go_build(False)
    #-------------
    #BUILD_WEB
    elif run_str == 'build_web':
        apps_names_lst = [app_name_str]
        web_meta_map   = gf_web_meta.get() 

        gf_web__build.build(apps_names_lst, web_meta_map, log_fun)
    #-------------
    #BUILD_CONTAINERS
    elif run_str == 'build_containers':

        #build using static linking, containers are based on Alpine linux, 
        #which has a minimal stdlib and other libraries, so we want to compile 
        #everything needed by this Go package into a single binary.
        go_build(True)
        
        web_meta_map = gf_web_meta.get()

        gf_containers.build(app_name_str, 
            build_meta_map,
            web_meta_map,
            log_fun)
    #-------------
    #TEST
    elif run_str == 'test':

        app_meta_map = build_meta_map[app_name_str]
        
        aws_creds_file_path_str = args_map['aws_creds']
        aws_creds_map           = gf_s3_utils.parse_creds(aws_creds_file_path_str)
        test_name_str           = args_map['test_name']
        
        gf_tests.run(app_name_str, test_name_str, app_meta_map, aws_creds_map)
    #-------------
    #LIST_CHANGED_APPS
    elif run_str == 'list_changed_apps':
        changed_apps_map = gf_build_changes.list_changed_apps(apps_changes_deps_map)
        gf_build_changes.view_changed_apps(changed_apps_map)
    #-------------
    else:
        print("unknown run command - %s"%(run_str))
        exit()
#--------------------------------------------------
def parse_args():

    arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

    #-------------
    #RUN
    arg_parser.add_argument('-run', action = "store", default = 'build',
        help = '''
- '''+fg('yellow')+'build'+attr(0)+'''             - build an app
- '''+fg('yellow')+'build_web'+attr(0)+'''         - build web code (ts/js/css/html) for an app
- '''+fg('yellow')+'build_containers'+attr(0)+'''  - build Docker containers for an app
- '''+fg('yellow')+'test'+attr(0)+'''              - run code tests for an app
- '''+fg('yellow')+'list_changed_apps'+attr(0)+''' - list all apps (and files) that have changed from last to the last-1 commit (for monorepo CI)

        ''')
    #-------------
    #APP
    arg_parser.add_argument('-app', action = "store", default = 'gf_images',
        help = '''
- '''+fg('yellow')+'gf_images'+attr(0)+'''
- '''+fg('yellow')+'gf_publisher'+attr(0)+'''
- '''+fg('yellow')+'gf_tagger'+attr(0)+'''
- '''+fg('yellow')+'gf_landing_page'+attr(0)+'''
- '''+fg('yellow')+'gf_analytics'+attr(0)+'''
- '''+fg('yellow')+'gf_crawl_lib'+attr(0)+'''
- '''+fg('yellow')+'gf_crawl_core'+attr(0)+'''
        ''')
    #-------------
    #AWS_S3_CREDS
    arg_parser.add_argument('-aws_creds',
        action =  "store",
        default = "%s/../../creds/aws/s3.txt"%(cwd_str),
        help =    '''path to the file containing AWS S3 credentials to be used''')
    #-------------
    #TEST_NAME
    arg_parser.add_argument('-test_name',
        action =  "store",
        default = "all",
        help =    '''if only a particular test needs to be run''')
    #-------------
    cli_args_lst   = sys.argv[1:]
    args_namespace = arg_parser.parse_args(cli_args_lst)
    args_map       = {
        "run":       args_namespace.run,
        "app":       args_namespace.app,
        "aws_creds": args_namespace.aws_creds,
        "test_name": args_namespace.test_name,
    }
    return args_map
#--------------------------------------------------
main()