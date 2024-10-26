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
cwd_str = os.path.abspath(os.path.dirname(__file__))

import boto3

import gf_os_aws_creds

#--------------------------------------------------
class GFs3info:
    imgs__bucket            = None
    discovered_imgs__bucket = None
    s3_resource             = None

#--------------------------------------------------
# FIX!! - use gf_aws_creds.get_from_file() directly
def parse_creds(p_aws_creds_file_path_str):
    aws_s3_creds_map = gf_aws_creds.get_from_file(p_aws_creds_file_path_str)
    
    assert "AWS_ACCESS_KEY_ID" in aws_s3_creds_map.keys()
    assert "AWS_SECRET_ACCESS_KEY" in aws_s3_creds_map.keys()
    assert(len(aws_s3_creds_map["AWS_ACCESS_KEY_ID"]) == 20)
    assert(len(aws_s3_creds_map["AWS_SECRET_ACCESS_KEY"]) == 40)
    
    return aws_s3_creds_map

#---------------------------------------------------
def s3_connect(p_aws_access_key_id_str,
    p_aws_secret_access_key_str):

    session = boto3.Session(
	    aws_access_key_id     = p_aws_access_key_id_str,
	    aws_secret_access_key = p_aws_secret_access_key_str,
	)
    
    s3_resource             = session.resource("s3")
    imgs__bucket            = s3_resource.Bucket("gf--img")
    discovered_imgs__bucket = s3_resource.Bucket("gf--discovered--img")

    gf_s3_info                         = GFs3info()
    gf_s3_info.imgs__bucket            = imgs__bucket
    gf_s3_info.discovered_imgs__bucket = discovered_imgs__bucket
    gf_s3_info.s3_resource             = s3_resource

    return gf_s3_info