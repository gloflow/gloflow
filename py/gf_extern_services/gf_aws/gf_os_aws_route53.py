# GloFlow application and media management/publishing platform
# Copyright (C) 2020 Ivan Trajkovic
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

#--------------------------------------------------
def hosted_zone__get_id(p_name_str,
    p_aws_creds_map,
    p_region_str = "us-east-1"):
    assert isinstance(p_name_str, str)

    aws_client = boto3.client("route53",
        aws_access_key_id     = p_aws_creds_map["AWS_ACCESS_KEY_ID"],
        aws_secret_access_key = p_aws_creds_map["AWS_SECRET_ACCESS_KEY"],
        region_name           = p_region_str)

    r = aws_client.list_hosted_zones()
    
    zone_id_str = None
    for z in r["HostedZones"]:
        if f"{p_name_str}." == z["Name"]:
            zone_id_str = z["Id"]
    
    if zone_id_str == None:
        print(f"ZONE with name [{p_name_str}] doesnt exist!")

    return zone_id_str