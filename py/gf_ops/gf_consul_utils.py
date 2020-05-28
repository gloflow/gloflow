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

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import subprocess

# sys.path.append("%s/../gf_core"%(modd_str))
# import gf_core_cli

#-------------------------------------------------------------
# SERVICES__CATALOG
def services__catalog(p_auth_token_str,
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_host_str                       = "127.0.0.1:8501"):

	# cmd_lst = []
	# if p_sudo_bool:
	# 	cmd_lst.append("sudo")
	#
	# cmd_lst.extend([
	# 	"CONSUL_HTTP_TOKEN=%s"%(p_auth_token_str),
	# 	"consul catalog services",
	# 	"-tags",
	# 	"-http-addr=https://%s"%(p_host_str),
	# ])
	
	# CONSUL_HTTP_API
	assert os.path.isfile(p_ca_intermediate__file_path_str)
	assert os.path.isfile(p_cert_combined__file_path_str)
	assert os.path.isfile(p_cert_key__file_path_str)

	url_str = "https://%s/v1/catalog/services"%(p_host_str)
	r = requests.get(url_str,
		verify=p_ca_intermediate__file_path_str,
		cert=(p_cert_combined__file_path_str, p_cert_key__file_path_str))


	print(r.text)
	
#-------------------------------------------------------------
# AGENT__START
def agent__start(p_name_str,
	p_type_str                 = "server",
	p_container_image_name_str = "consul:1.7.2",
	p_config__file_path_str    = None,
	p_external_ip_str          = "127.0.0.1",
	p_root_agent_ip_str        = "127.0.0.1",
	p_bootstrap_server_bool    = False,
	p_datacenter_str           = "us-east-1",
	p_volume_mounts_lst        = None,
	p_sudo_bool                = False):
	assert p_type_str == "client" or p_type_str == "server"
	assert os.path.isfile(p_config__file_path_str)

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"docker run",
		"--restart=always",
		"--net=host",
		"--name=%s"%(p_name_str),

		# turn off the Consul dashboard
		'''-e "CONSUL_UI=false"'''
	])

	#-------------
	# VOLUMES
	if not p_config__file_path_str == None:
		cmd_lst.append("-v %s:/consul/config/%s"%(p_config__file_path_str, os.path.basename(p_config__file_path_str)))

	if not p_volume_mounts_lst == None:
		assert isinstance(p_volume_mounts_lst, list)
		for source_str, target_str in p_volume_mounts_lst:
			cmd_lst.append("-v %s:%s"%(source_str, target_str))

	#-------------

	if p_type_str == "client":

		# ENV
		cmd_lst.extend([
			# '''-e "CONSUL_LOCAL_CONFIG={'leave_on_terminate': true}"''' # its default since Consul 0.7
		])
		
	# MAIN
	cmd_lst.extend([
		p_container_image_name_str,

		"agent",
		"-datacenter=%s"%(p_datacenter_str),
	])

	if p_type_str == "server":
		cmd_lst.extend([

			# BIND_IP
			# this is the address advertised to the rest of the cluster.
			# should be reachable by all other nodes in the cluster.
			# Consul uses both TCP and UDP and the same port for both.
			"-bind=%s"%(p_external_ip_str),
			
			# specifies how many server agents to watch for before bootstrapping the cluster for the first time
			"-bootstrap-expect=1"
		])

		if p_bootstrap_server_bool:

			# flag is used to control if a server is in "bootstrap" mode. It is important that
			# no more than one server per datacenter be running in this mode. Technically, a server in 
			# bootstrap mode is allowed to self-elect as the Raft leader.
			cmd_lst.append("-bootstrap")

	if p_type_str == "client":

		cmd_lst.extend([
			# address of the agent to use to join the cluster
			"-retry-join=%s"%(p_root_agent_ip_str),

			# DNS
			# "-dns-port=53"
		])

	c_str = " ".join(cmd_lst)
	print(c_str)
	consul_agent_p = subprocess.Popen(c_str, shell=True, stdout=subprocess.PIPE, bufsize=1)
	
	return consul_agent_p