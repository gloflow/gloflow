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

#--------------------------------------------------
# FILE

def get_from_file(p_aws_creds_file_path_str):
    print(p_aws_creds_file_path_str)
    assert os.path.isfile(os.path.abspath(p_aws_creds_file_path_str))

    f             = open(p_aws_creds_file_path_str, 'r')
    aws_creds_map = {}
    for l in f.readlines():

        if l == '' or l == '\n': continue
        if l.startswith('#'):    continue #ignore comments
        
        k, v = l.strip().split("=")
        k    = k.strip()
        v    = v.strip()

        if k == "AWS_ACCESS_KEY_ID" or \
            k == "AWS_SECRET_ACCESS_KEY":
            #k == "GF_AWS_TOKEN":
            aws_creds_map[k]=v
    f.close()   
    return aws_creds_map

#--------------------------------------------------
# ENV_VARS

def get_from_env_vars():
    aws_creds_map = {
        "aws_access_key_id_str": os.environ["AWS_ACCESS_KEY_ID"],
        "aws_secret_access_key": os.environ["AWS_SECRET_ACCESS_KEY"],
        #"aws_token_str": os.environ["GF_AWS_TOKEN"]
    }
    return aws_creds_map
    