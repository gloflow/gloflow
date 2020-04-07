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

import time
from colored import fg, bg, attr

sys.path.append("%s/../../../go/gf_core/py"%(modd_str))
import gf_core_cli

#---------------------------------------------------
# ROOT
def generate__ca_root(p_output_files_base_str,
	p_config__file_path_str,
	p_sudo_bool = False):
	assert os.path.isfile(p_config__file_path_str)
	assert p_config__file_path_str.endswith(".json")
	
	print("%sGENERATE ROOT_CA%s"%(fg("yellow"), attr(0)))

	c_lst = []

	#-----------------
	# GENERATE
	if p_sudo_bool:
		c_lst.append("sudo")
	
	c_lst.extend([
		"cfssl gencert",
		"-initca", # "-initca" - initialise new CA
		p_config__file_path_str
	])

	#-----------------
	c_lst.append("|")

	#-----------------
	# SAVE_TO_FS
	if p_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		# "-bare" - the response from CFSSL is not wrapped in the API standard response
		"cfssljson -bare %s"%(p_output_files_base_str)
	])

	#-----------------
	c_str = " ".join(c_lst)
	_, _, return_code = gf_core_cli.run_cmd(c_str, p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

	print("%sdone...%s"%(fg("green"), attr(0)))

#-------------------------------------------------------------
# INTERMEDIATE
def generate__ca_intermediate(p_output_files_base_str,
	p_root_ca_base_str,
	p_config__file_path_str,
	p_config_csr__file_path_str,
	p_profile_name_str = None,
	p_sudo_bool        = False):
	assert os.path.isfile(p_config__file_path_str)
	assert os.path.isfile(p_config_csr__file_path_str)

	print("%sGENERATE INTERMEDIATE_CA%s"%(fg("yellow"), attr(0)))

	# ROOT_CA
	root_ca_cert__file_path_str = "%s.pem"%(p_root_ca_base_str)
	root_ca_key__file_path_str  = "%s-key.pem"%(p_root_ca_base_str)

	c_lst = []
	
	#-----------------
	# GENERATE
	if p_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		"cfssl gencert",

		# ROOT_CA
		"-ca %s"%(root_ca_cert__file_path_str),
		"-ca-key %s"%(root_ca_key__file_path_str),

		"-config %s"%(p_config__file_path_str)
	])

	if not p_profile_name_str == None:
		c_lst.append("-profile %s"%(p_profile_name_str))

	c_lst.append(p_config_csr__file_path_str)

	#-----------------
	c_lst.append("|")

	#-----------------
	# SAVE_TO_FS
	if p_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		"cfssljson -bare %s"%(p_output_files_base_str)
	])
	
	#-----------------
	c_str = " ".join(c_lst)
	_, _, return_code = gf_core_cli.run_cmd(c_str, p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

	print("%sdone...%s"%(fg("green"), attr(0)))

#-------------------------------------------------------------
# LEAF
def generate__cert_leaf(p_output_files_base_str,
	p_ca_base_str,
	p_config__file_path_str):

	# INTERMEDIATE_CA
	ca_cert__file_path_str = "%s.pem"%(p_ca_base_str)
	ca_key__file_path_str  = "%s-key.pem"%(p_ca_base_str)

	c_lst = []
	
	#-----------------
	# GENERATE
	if p_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		"cfssl gencert",

		# INTERMEDIATE_CA
		"-ca=%s"%(ca_cert__file_path_str),
		"-ca-key=%s"%(ca_key__file_path_str),
		"-config=%s"%(p_config__file_path_str),
		'''-hostname="server.global.nomad,localhost,127.0.0.1"''', "-",
	])

	#-----------------
	c_lst.append("|")

	#-----------------
	# SAVE_TO_FS
	if p_sudo_bool:
		c_lst.append("sudo")

	c_lst.append("cfssljson -bare %s"%(p_output_files_base_str))

	#-----------------
	c_str = " ".join(c_lst)
	_, _, return_code = gf_core_cli.run_cmd(c_str, p_env_map = None)
	if not return_code == 0:
		print("CLI failed...")
		exit()

	print("%sdone...%s"%(fg("green"), attr(0)))

#-------------------------------------------------------------
def archive_if_exists(p_files_base_str, p_sudo_bool = False):
	
	dir_str       = os.path.abspath(os.path.dirname(p_files_base_str))
	file_base_str = os.path.basename(p_files_base_str)

	if p_sudo_bool: sudo_str = "sudo"
	else:           sudo_str = ""

	# list all files in target dir
	# "-1"  - force output to be one entry per line
	# "^%s" - pattern matches the file_base only at the start of the line
	stdout_str, _, return_code = gf_core_cli.run_cmd("%s ls -1 %s | grep '^%s'"%(sudo_str, dir_str, file_base_str),
		p_env_map = None,
		p_print_output_bool = True)
	
	if stdout_str == "":
		return True

	stdout_clean_str = stdout_str.strip()
	lines_lst        = stdout_clean_str.split("\n")

	if len(lines_lst) > 0:

		# IMPORTANT!! - ask use if they want to recreate/archive existing certs. if they dont
		#               dont archive and return False
		print("CERT ALREADY EXISTS")
		if not gf_core_cli.confirm("recreate cert (and archive old)?"):
			return False

		archive_time = time.time()

		# process each file that needs to be archived
		for l in lines_lst:
			file_name_str = l.split()[-1:][0]
			file_path_str = "%s/%s"%(dir_str, file_name_str)

			# ARCHIVE_FILE - rename the file
			c = "%s mv %s %s/old_%s__%s"%(sudo_str, file_path_str, dir_str, archive_time, file_name_str)
			_, _, return_code = gf_core_cli.run_cmd(c, p_env_map = None)
			if not return_code == 0:
				print("CLI failed...")
				exit()
	
	return True