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
from botocore.exceptions import ClientError

#-------------------------------------------------------------
def describe_instances_by_tags(p_tags_map,
	p_region_str="us-east-1"):
	
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

#-------------------------------------------------------------
# SECURITY_GROUPS
#-------------------------------------------------------------
def get_secgr_id_by_name(p_name_str, p_region_str="us-east-1"):
	assert(isinstance(p_name_str, str))

	client = boto3.client("ec2", region_name=p_region_str)
	
	response = client.describe_security_groups(
		Filters=[{"Name": "group-name", "Values": [p_name_str]}]
	)

	if response["SecurityGroups"]:
		id_str = response["SecurityGroups"][0]["GroupId"]
		return id_str
	else:
		print(f"Security group '{p_name_str}' not found.")
		return None
	
#---------------------------------------------------------------------------------
def add_secgr_ingress_ip(p_ip_cidr_str,
	p_port_int,
	p_sec_gr_id_str,
	p_region_str="us-east-1"):
	assert(isinstance(p_ip_cidr_str, str))
	assert(isinstance(p_port_int, int))
	assert(isinstance(p_sec_gr_id_str, str))

	try:
		client = boto3.client("ec2", region_name=p_region_str)
		
		response = client.authorize_security_group_ingress(
			GroupId=p_sec_gr_id_str,
			IpPermissions=[
				{
					'IpProtocol': 'tcp',
					'FromPort': p_port_int,
					'ToPort': p_port_int,
					'IpRanges': [{'CidrIp': p_ip_cidr_str}]
				}
			]
		)
		
		if response['Return'] is not True:
			raise Exception(f"Ingress rule addition failed: {response}")

	except ClientError as e:
		raise Exception(f"Failed to add ingress rule: {e}")
	
#---------------------------------------------------------------------------------
def delete_secgr_ingress_ip(p_ip_cidr_str,
	p_port_int,
	p_sec_gr_id_str,
	p_region_str="us-east-1"):
	assert(isinstance(p_ip_cidr_str, str))
	assert(isinstance(p_port_int, int))
	assert(isinstance(p_sec_gr_id_str, str))

	try:
		client = boto3.client("ec2", region_name=p_region_str)

		response = client.revoke_security_group_ingress(
			GroupId=p_sec_gr_id_str,
			IpPermissions=[
				{
					'IpProtocol': 'tcp',
					'FromPort': p_port_int,
					'ToPort': p_port_int,
					'IpRanges': [{'CidrIp': p_ip_cidr_str}]
				}
			]
		)
		
		if response['Return'] is not True:
			raise Exception(f"Ingress rule revoking failed: {response}")
		
	except ClientError as e:
		raise Exception(f"Failed to revoke ingress rule: {e}")

#---------------------------------------------------------------------------------
def check_secgr_ingress_ip(p_ip_cidr_str,
	p_port_int,
	p_sec_gr_id_str,
	p_region_str="us-east-1"):
	assert(isinstance(p_ip_cidr_str, str))
	assert(isinstance(p_port_int, int))
	assert(isinstance(p_sec_gr_id_str, str))
	
	try:
		client   = boto3.client("ec2", region_name=p_region_str)
		response = client.describe_security_groups(GroupIds=[p_sec_gr_id_str])

		for permission in response['SecurityGroups'][0]['IpPermissions']:
			if (permission['FromPort'] == p_port_int and
				permission['ToPort'] == p_port_int and
				permission['IpProtocol'] == 'tcp'):
				for ip_range in permission['IpRanges']:
					if ip_range['CidrIp'] == p_ip_cidr_str:
						return True
		return False

	except ClientError as e:
		raise Exception(f"Failed to check ingress rule: {e}")

#---------------------------------------------------------------------------------
# check if the IP/port is allowed in the security group, and if its not add it

def ensure_secgr_ingress_ip(p_ip_cidr_str,
	p_port_int,
	p_sec_gr_id_str,
	p_region_str="us-east-1"):
	assert(isinstance(p_ip_cidr_str, str))
	assert(isinstance(p_port_int, int))
	assert(isinstance(p_sec_gr_id_str, str))

	if not check_secgr_ingress_ip(p_ip_cidr_str, p_port_int, p_sec_gr_id_str, p_region_str):
		add_secgr_ingress_ip(p_ip_cidr_str, p_port_int, p_sec_gr_id_str, p_region_str)