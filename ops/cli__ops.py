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

import argparse
from colored import fg, bg, attr
import delegator

sys.path.append('%s/../meta'%(modd_str))
import gf_meta

sys.path.append('%s/aws/s3'%(modd_str))
import gf_s3_data_info
import gf_s3_utils

#--------------------------------------------------
def main():
    
    print ''
    print '                   -------------  %sOPS%s %sGLOFLOW%s  -------------'%(bg('dark_orange_3a'), attr(0), bg('cyan'), attr(0))
    print ''

    b_meta_map = gf_meta.get()['build_info_map']
    args_map   = parse_args()
    run_str    = args_map['run']

    aws_creds_file_path_str = args_map['aws_creds']
    aws_creds_map           = gf_s3_utils.parse_creds(aws_creds_file_path_str)
    assert isinstance(aws_creds_map,dict)

    #-------------
    if run_str == 's3_data_info':
        aws_access_key_id_str     = aws_creds_map['GF_AWS_ACCESS_KEY_ID']
        aws_secret_access_key_str = aws_creds_map['GF_AWS_SECRET_ACCESS_KEY']
        gf_s3_data_info.stats__image_buckets_general(aws_access_key_id_str, aws_secret_access_key_str)
        
    #-------------

#--------------------------------------------------
def parse_args():

    arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

    #-------------
    #RUN
    arg_parser.add_argument('-run', action = "store", default = 'build',
        help = '''
- '''+fg('yellow')+'s3_data_info'+attr(0)+''' - view AWS S3 data information summaries of files used by GF
        ''')
    #-------------
    #AWS_S3_CREDS
    arg_parser.add_argument('-aws_creds', action = "store", default = '%s/../../creds/aws/s3.txt'%(modd_str), help = '''path to the file containing AWS S3 credentials to be used''')
    #-------------
    cli_args_lst   = sys.argv[1:]
    args_namespace = arg_parser.parse_args(cli_args_lst)
    args_map       = {
        "run":       args_namespace.run,
        "aws_creds": args_namespace.aws_creds,
    }
    return args_map

#--------------------------------------------------
main()