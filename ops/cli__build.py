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
#--------------------------------------------------
def main():
    
    print ''
    print '                              %sBUILD GLOFLOW%s'%(fg('green'),attr(0))
    print ''

    b_meta_map = gf_meta.get()['build_info_map']
    args_map   = parse_args()
    run_str    = args_map['run']

    if run_str == 'build' or run_str == 'test':
        app_name_str = args_map['app']
        assert b_meta_map.has_key(app_name_str)
        app_meta_map = b_meta_map[app_name_str]

    #-------------
    if run_str == 'build':
        build__go_bin(app_name_str, app_meta_map['go_path_str'], app_meta_map['go_output_path_str'])
    #-------------
    elif run_str == 'test':

        aws_s3_creds_file_path_str = args_map['aws_s3_creds']
        aws_s3_creds_map           = parse_creds(aws_s3_creds_file_path_str)

        test(app_name_str, aws_s3_creds_map)
    #-------------
#--------------------------------------------------
def build__go_bin(p_name_str,
    p_main_go_file_path_str,
    p_output_path_str):
    assert os.path.isfile(p_main_go_file_path_str)
    assert os.path.isdir(os.path.dirname(p_output_path_str))

    print ''
    print ' -- build %s%s%s service'%(fg('green'), p_name_str, attr(0))
    
    cwd_str = os.getcwd()
    os.chdir(os.path.dirname(p_main_go_file_path_str)) #change into the target main package dir

    c = 'go build -o %s'%(p_output_path_str)
    print c
    r = delegator.run(c)
    if not r.out == '': print r.out
    if not r.err == '': print '%sFAILED%s >>>>>>>\n%s'%(fg('red'),attr(0),r.err)

    os.chdir(cwd_str) #return to initial dir
#--------------------------------------------------
def test(p_name_str,
    p_aws_s3_creds_map):
    assert isinstance(p_aws_s3_creds_map,dict)
    print ''
    print ' -- test %s%s%s service'%(fg('green'), p_name_str, attr(0))

    cwd_str = os.getcwd()
    os.chdir(os.path.dirname(p_main_go_file_path_str)) #change into the target main package dir

    c = 'go test'
    print c
    r = delegator.run(c,env=p_aws_s3_creds_map)
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
- '''+fg('yellow')+'build'+attr(0)+''' - build a particular app
- '''+fg('yellow')+'test'+attr(0)+'''  - run code tessts
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
    cli_args_lst   = sys.argv[1:]
    args_namespace = arg_parser.parse_args(cli_args_lst)
    args_map       = {
        "run":args_namespace.run,
        "app":args_namespace.app,
    }
    return args_map
#--------------------------------------------------
main()