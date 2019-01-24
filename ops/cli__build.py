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

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import argparse

from colored import fg,bg,attr
import delegator

sys.path.append('%s/../meta'%(cwd_str))
import gf_meta

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
    print('                              %sBUILD GLOFLOW%s'%(fg('green'),attr(0)))
    print('')

    #--------------------------------------------------
    def log_fun(g, m):
        if g == "ERROR":
            print('%s%s%s:%s%s%s'%(bg('red'), g, attr(0), fg('red'), m, attr(0)))
        else:
            print('%s%s%s:%s%s%s'%(fg('yellow'), g, attr(0), fg('green'), m, attr(0)))
    #--------------------------------------------------
    
    b_meta_map = gf_meta.get()['build_info_map']
    args_map   = parse_args()
    run_str    = args_map['run']

    app_name_str = args_map['app']
    assert b_meta_map.has_key(app_name_str)

    #-------------
    #BUILD
    if run_str == 'build':
        if not app_meta_map.has_key('go_output_path_str'):
            print("not a main package")
            exit()

        app_meta_map = b_meta_map[app_name_str]
        build__go_bin(app_name_str, app_meta_map['go_path_str'], app_meta_map['go_output_path_str'])
    #-------------
    #BUILD_WEB
    elif run_str == 'build_web':
        apps_names_lst = [app_name_str]
        gf_web__build.build(apps_names_lst, log_fun)
    #-------------
    #BUILD_CONTAINERS
    elif run_str == 'build_containers':
        gf_containers.build(log_fun)
    #-------------
    #TEST
    elif run_str == 'test':
        aws_creds_file_path_str = args_map['aws_creds']
        aws_creds_map           = gf_s3_utils.parse_creds(aws_creds_file_path_str)
        test_name_str           = args_map['test_name']
        
        gf_tests.run(app_name_str, test_name_str, app_meta_map, aws_creds_map)
    #-------------
#--------------------------------------------------
def build__go_bin(p_name_str,
    p_go_path_str,
    p_output_path_str):
    assert os.path.isdir(p_go_path_str)
    assert os.path.isdir(os.path.dirname(p_output_path_str))

    print ''
    print ' -- build %s%s%s service'%(fg('green'), p_name_str, attr(0))
    
    cwd_str = os.getcwd()
    os.chdir(p_go_path_str) #change into the target main package dir

    c = 'go build -o %s'%(p_output_path_str)
    print c
    r = delegator.run(c)
    if not r.out == '': print r.out
    if not r.err == '': print '%sFAILED%s >>>>>>>\n%s'%(fg('red'),attr(0),r.err)

    os.chdir(cwd_str) #return to initial dir

#--------------------------------------------------
def parse_args():

    arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

    #-------------
    #RUN
    arg_parser.add_argument('-run', action = "store", default = 'build',
        help = '''
- '''+fg('yellow')+'build'+attr(0)+'''            - build an app
- '''+fg('yellow')+'build_web'+attr(0)+'''        - build web code (ts/js/css/html) for an app
- '''+fg('yellow')+'build_containers'+attr(0)+''' - build Docker containers for an app
- '''+fg('yellow')+'test'+attr(0)+'''             - run code tests for an app

        ''')
    #-------------
    #APP
    arg_parser.add_argument('-app', action = "store", default = 'build',
        help = '''
- '''+fg('yellow')+'gf_images'+attr(0)+'''
- '''+fg('yellow')+'gf_publisher'+attr(0)+'''
- '''+fg('yellow')+'gf_tagger'+attr(0)+'''
- '''+fg('yellow')+'gf_landing_page'+attr(0)+'''
- '''+fg('yellow')+'gf_analytics'+attr(0)+'''
- '''+fg('yellow')+'gf_crawl_lib'+attr(0)+'''
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
        "run":      args_namespace.run,
        "app":      args_namespace.app,
        "aws_creds":args_namespace.aws_creds,
        "test_name":args_namespace.test_name,
    }
    return args_map
#--------------------------------------------------
main()