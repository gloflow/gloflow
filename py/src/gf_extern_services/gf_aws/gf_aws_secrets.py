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
def init_client(p_region_str):
    
    aws_client = boto3.client("secretsmanager",
        region_name = p_region_str)
    return aws_client

#--------------------------------------------------
def get_secret(p_secret_name_str,
    p_client):

    secret_value = p_client.get_secret_value(SecretId=p_secret_name_str)
    secret_value_str = secret_value["SecretString"]
    return secret_value_str

#--------------------------------------------------
def create_or_update_batch(p_secrets_map,
    p_region_str,
    p_aws_creds_map):
    assert isinstance(p_secrets_map, dict)
    assert isinstance(p_aws_creds_map, dict)

    aws_client = boto3.client("secretsmanager",
        aws_access_key_id     = p_aws_creds_map["AWS_ACCESS_KEY_ID"],
        aws_secret_access_key = p_aws_creds_map["AWS_SECRET_ACCESS_KEY"],
        region_name           = p_region_str)

    r = aws_client.list_secrets()
    existing_secrets_lst = [secret_map["Name"] for secret_map in r["SecretList"]]

    for secret_name_str, secret_value_str in p_secrets_map.items():

        # UPDATE - if the secret already exists
        if secret_name_str in existing_secrets_lst:
            aws_client.update_secret(SecretId = secret_name_str,
                SecretString = secret_value_str)

        # CREATE - create a secret if it doesnt exist yet
        else:
            aws_client.create_secret(Name = secret_name_str,
                SecretString = secret_value_str)