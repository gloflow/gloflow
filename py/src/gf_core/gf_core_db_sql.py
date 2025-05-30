# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
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
import json
import boto3
import psycopg2

#---------------------------------------------------------------------------------
def init_db_client(p_db_name_str,
	p_env_str,
	p_from_aws_secrets_bool=False):
	
	if p_from_aws_secrets_bool:
		# SECRETS
		print("getting meta from secrets...")
		db_user_str, db_pass_str, db_host_str, db_port_str = get_meta_from_secrets(p_env_str)
	else:
		# ENV
		print("getting meta from env...")
		db_user_str, db_pass_str, db_host_str, db_port_str = get_meta_from_env(p_env_str)

	print("connect start...")
	db_client = psycopg2.connect(
		host     = db_host_str,
		port     = int(db_port_str),
		database = p_db_name_str,
		user     = db_user_str,
		password = db_pass_str)
	print("connect ok...")

	return db_client

#---------------------------------------------------------------------------------
def get_meta_from_env(p_db_env_str,
	p_prefix_str="GF_DB"):

	print("getting db meta from env...")

	db_host_port_str = os.environ.get(f"{p_prefix_str}_HOST")
	db_user_str = os.environ.get(f"{p_prefix_str}_CREDS_USER")
	db_pass_str = os.environ.get(f"{p_prefix_str}_CREDS_PASS")

	db_host_str, db_port_str = db_host_port_str.split(":")

	return db_user_str, db_pass_str, db_host_str, db_port_str

#---------------------------------------------------------------------------------
def get_meta_from_secrets(p_db_env_str,
	p_prefix_str="gf_rds"):

	print("getting db meta from aws secrets_manager...")

	secrets_client = boto3.client('secretsmanager',
		region_name="us-east-1")
	
	db_host_port_str = secrets_client.get_secret_value(SecretId=f"{p_prefix_str}_host_{p_db_env_str}")["SecretString"]
	db_creds_str     = secrets_client.get_secret_value(SecretId=f"{p_prefix_str}_creds_{p_db_env_str}")["SecretString"]
	print("db RDS creds fetched from aws secrets_manager...")

	db_creds_map = json.loads(db_creds_str)
	db_user_str = db_creds_map["username"]
	db_pass_str = db_creds_map["password"]

	db_host_str, db_port_str = db_host_port_str.split(":")

	return db_user_str, db_pass_str, db_host_str, db_port_str

#---------------------------------------------------------------------------------
def table_exists(p_table_name_str,
	p_db_cursor):
	p_db_cursor.execute(f"SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = '{p_table_name_str}')")
	exists_bool = p_db_cursor.fetchone()[0]
	return exists_bool