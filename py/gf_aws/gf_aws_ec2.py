# GloFlow application and media management/publishing platform
# Copyright (C) 2021 Ivan Trajkovic
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

#-------------------------------------------------------------
def describe_instances__by_tags(p_tags_map, p_region_str):
    

    client = boto3.client("ec2", region_name=p_region_str)

    custom_filter = []
    for k, v in p_tags_map.items():
        custom_filter.append({
            "Name":   f"tag:{k}", 
            "Values": [v]
        })
    
    response = client.describe_instances(Filters=custom_filter)

    aws_instances_lst = []
    for i in response["Reservations"][0]["Instances"]:
        aws_instances_lst.append(i)

    return aws_instances_lst