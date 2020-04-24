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

sys.path.append("%s/../gf_core"%(modd_str))
import gf_core_cli


#----------------------------------------------
# JOBS
#----------------------------------------------
def job_run(p_config__file_path_str,
	p_acl_token_secret_id_str,
	p_host_url_str                   = "127.0.0.1:4646",
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_sudo_bool                      = False):

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"NOMAD_TOKEN='%s'"%(p_acl_token_secret_id_str),
		"nomad job run",
		"-address=https://%s"%(p_host_url_str)
	])

	if not p_ca_intermediate__file_path_str == None:
		cmd_lst.extend([
			# path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
			"-ca-cert=%s"%(p_ca_intermediate__file_path_str),
			
			# ath to a PEM encoded client certificate for TLS authentication to the Nomad server.
			"-client-cert=%s"%(p_cert_combined__file_path_str),
			"-client-key=%s"%(p_cert_key__file_path_str),
		])

	cmd_lst.extend([
		p_config__file_path_str
	])

	stdout_str, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

#----------------------------------------------
def job_status(p_name_str,
	p_acl_token_secret_id_str,
	p_host_url_str                   = "127.0.0.1:4646",
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_sudo_bool                      = False):

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"NOMAD_TOKEN='%s'"%(p_acl_token_secret_id_str),
		"nomad job status",
		"-address=https://%s"%(p_host_url_str)
	])
	
	if not p_ca_intermediate__file_path_str == None:
		cmd_lst.extend([
			# path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
			"-ca-cert=%s"%(p_ca_intermediate__file_path_str),
			
			# ath to a PEM encoded client certificate for TLS authentication to the Nomad server.
			"-client-cert=%s"%(p_cert_combined__file_path_str),
			"-client-key=%s"%(p_cert_key__file_path_str),
		])

	stdout_str, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

#----------------------------------------------
# ACL
#----------------------------------------------
# ACL_BOOTSTRAP
def acl_bootstrap(p_host_url_str = "127.0.0.1:4646",
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_sudo_bool                      = False):

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"nomad acl bootstrap"
		"-address=https://%s"%(p_host_url_str)
	])

	if not p_ca_intermediate__file_path_str == None:
		cmd_lst.extend([
			# path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
			"-ca-cert=%s"%(p_ca_intermediate__file_path_str),
			
			# ath to a PEM encoded client certificate for TLS authentication to the Nomad server.
			"-client-cert=%s"%(p_cert_combined__file_path_str),
			"-client-key=%s"%(p_cert_key__file_path_str),
		])

	stdout_str, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

#----------------------------------------------
# ACL_TOKEN_CREATE
def acl_token_create(p_name_str,
	p_output_file_path_str,
	p_acl_token_secret_id_str,
	p_policies_lst                   = [],
	p_type_str                       = "client",
	p_host_url_str                   = "127.0.0.1:4646",
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_sudo_bool = False):
	assert p_type_str == "management" or \
		p_type_str == "client"

	# "nomad acl token self" - get information about the current token

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"NOMAD_TOKEN='%s'"%(p_acl_token_secret_id_str),
		"nomad acl token create",

		"-name='%s'"%(p_name_str),
		"-type='%s'"%(p_type_str),

		# GLOBAL_TOKEN - are created in the authoritative region and then replicate to all other regions
		"-global",

		"-address=https://%s"%(p_host_url_str)
	])

	if not p_ca_intermediate__file_path_str == None:
		cmd_lst.extend([
			# path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
			"-ca-cert=%s"%(p_ca_intermediate__file_path_str),
			
			# ath to a PEM encoded client certificate for TLS authentication to the Nomad server.
			"-client-cert=%s"%(p_cert_combined__file_path_str),
			"-client-key=%s"%(p_cert_key__file_path_str),
		])

	for p in p_policies_lst:
		cmd_lst.append("-policy='%s'"%(p))

	stdout_str, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

	#-------------
	# WRITE_TO_FILE

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")
	fs_write_cmd_str = '''bash -c "echo '%s' > %s"'''%(stdout_str, p_output_file_path_str)
	cmd_lst.append(fs_write_cmd_str)

	_, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

	#-------------

#----------------------------------------------
# ACL_POLICY_APPLY
def acl_policy_apply(p_name_str,
	p_description_str,
	p_config__file_path_str,
	p_acl_token_secret_id_str,
	p_host_url_str                   = "127.0.0.1:4646",
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_sudo_bool = False):
	
	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"NOMAD_TOKEN='%s'"%(p_acl_token_secret_id_str),
		"nomad acl policy apply",
		"-description '%s'"%(p_description_str),
		"-address=https://%s"%(p_host_url_str),
	])

	if not p_ca_intermediate__file_path_str == None:
		cmd_lst.extend([
			# path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
			"-ca-cert=%s"%(p_ca_intermediate__file_path_str),
			
			# ath to a PEM encoded client certificate for TLS authentication to the Nomad server.
			"-client-cert=%s"%(p_cert_combined__file_path_str),
			"-client-key=%s"%(p_cert_key__file_path_str),
		])

	cmd_lst.extend([
		p_name_str,
		p_config__file_path_str
	])

	stdout_str, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

#----------------------------------------------
# ACL_POLICY_LIST
def acl_policy_list(p_acl_token_secret_id_str,
	p_host_url_str = "127.0.0.1:4646",
	p_ca_intermediate__file_path_str = None,
	p_cert_combined__file_path_str   = None,
	p_cert_key__file_path_str        = None,
	p_sudo_bool = False):

	cmd_lst = []
	if p_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"NOMAD_TOKEN='%s'"%(p_acl_token_secret_id_str),
		"nomad acl policy list",
		"-address=https://%s"%(p_host_url_str),
	])

	if not p_ca_intermediate__file_path_str == None:
		cmd_lst.extend([
			# path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
			"-ca-cert=%s"%(p_ca_intermediate__file_path_str),
			
			# ath to a PEM encoded client certificate for TLS authentication to the Nomad server.
			"-client-cert=%s"%(p_cert_combined__file_path_str),
			"-client-key=%s"%(p_cert_key__file_path_str),
		])

	stdout_str, _, return_code = gf_core_cli.run_cmd(" ".join(cmd_lst), p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

	policies_lst = []
	for l in stdout_str.strip().split("\n")[1:]:

		policy_name_str = l.split()[0].strip()
		policies_lst.append({
			"name_str": policy_name_str
		})
	return policies_lst