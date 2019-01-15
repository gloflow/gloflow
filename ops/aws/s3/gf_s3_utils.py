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


import boto3
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

    return imgs__bucket, discovered_imgs__bucket, s3_resource