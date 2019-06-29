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
import boto3

class Gf_s3_info:
    imgs__bucket            = None
    discovered_imgs__bucket = None
    s3_resource             = None
#--------------------------------------------------
def parse_creds(p_aws_creds_file_path_str):
    print(p_aws_creds_file_path_str)
    assert os.path.isfile(os.path.abspath(p_aws_creds_file_path_str))

    f                = open(p_aws_creds_file_path_str,'r')
    aws_s3_creds_map = {}
    for l in f.readlines():

        if l == '' or l == '\n': continue
        if l.startswith('#'):    continue #ignore comments
        
        k, v = l.strip().split("=")
        k    = k.strip()
        v    = v.strip()

        if k == "GF_AWS_ACCESS_KEY_ID" or \
            k == "GF_AWS_SECRET_ACCESS_KEY" or \
            k == "GF_AWS_TOKEN":
            aws_s3_creds_map[k]=v
    f.close()   
    return aws_s3_creds_map
#---------------------------------------------------
def s3_connect(p_aws_access_key_id_str,
    p_aws_secret_access_key_str):

    session = boto3.Session(
	    aws_access_key_id     = p_aws_access_key_id_str,
	    aws_secret_access_key = p_aws_secret_access_key_str,
	)
    
    s3_resource             = session.resource('s3')
    imgs__bucket            = s3_resource.Bucket('gf--img')
    discovered_imgs__bucket = s3_resource.Bucket('gf--discovered--img')

    gf_s3_info                         = Gf_s3_info()
    gf_s3_info.imgs__bucket            = imgs__bucket
    gf_s3_info.discovered_imgs__bucket = discovered_imgs__bucket
    gf_s3_info.s3_resource             = s3_resource

    return gf_s3_info