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
import threading

from colored import fg, bg, attr

import gf_test_utils
import gf_test__sqs_process
import gf_test__master_get_block

#--------------------------------------------------
def run_all():
	
	assert "AWS_ACCESS_KEY_ID" in os.environ.keys()
	assert "AWS_SECRET_ACCESS_KEY" in os.environ.keys()
	assert "GF_AWS_SQS_QUEUE" in os.environ.keys()
	assert "GF_SENTRY_ENDPOINT" in os.environ.keys()


	aws_region_str               = "us-east-1"
	test_gf_geth_path_str        = f"{modd_str}/../../../gloflow_go-ethereum/build/bin/geth"
	data_dir_path_str            = f"{modd_str}/test_data_geth"
	py_plugins_base_dir_path_str = f"{modd_str}/../plugins"

	assert os.path.isfile(test_gf_geth_path_str)
	assert os.path.isdir(py_plugins_base_dir_path_str)


	



	# INIT
	p, p__worker_inspector = init(test_gf_geth_path_str,
		data_dir_path_str,
		py_plugins_base_dir_path_str,
		aws_region_str)







	# gf_test__sqs_process.run(aws_region_str)

	gf_test__master_get_block.run(aws_region_str)






	#-------------------------------------------------------------
	def cleanup():
		p.terminate()
		p__worker_inspector.terminate()

	#-------------------------------------------------------------
	cleanup()

#--------------------------------------------------
def init(p_geth_bin_path_str,
	p_geth_data_dir_path_str,
	p_py_plugins_base_dir_path_str,
	p_aws_region_str):
	assert os.path.isdir(p_py_plugins_base_dir_path_str)

	#-------------------------------------------------------------
	# START_GETH
	def start_geth():
	  
		print("%srunning GETH...%s"%(fg("yellow"), attr(0)))
		
		assert "AWS_ACCESS_KEY_ID" in os.environ.keys()
		assert "AWS_SECRET_ACCESS_KEY" in os.environ.keys()
		assert "GF_AWS_SQS_QUEUE" in os.environ.keys()
		
		cmd_lst = [
			p_geth_bin_path_str,
			f"--datadir={p_geth_data_dir_path_str}",
			"--syncmode=fast",
			"--http",
			"--maxpeers=10"
		]
		c_str = " ".join(cmd_lst)
		print(c_str)

		# When shell=True the shell is the child process, and the commands are its children.
		# So any SIGTERM or SIGKILL will kill the shell but not its child processes.
		# The best way I can think of is to use shell=False, otherwise when you kill
		# the parent shell process, it will leave a defunct shell process.
		# CMD also has to be a list here, since its not being passed in as a string
		# to the child shell.
		geth_p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
			env={
				"AWS_REGION":            p_aws_region_str,
				"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
				"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
				"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"]
			})

		t = threading.Thread(target=gf_test_utils.read_process_stdout, args=(geth_p.stdout, "geth", "cyan"))
		t.start()

		return geth_p

	#-------------------------------------------------------------
	def start_worker_inspector():
		cmd_lst = [
			f"{modd_str}/../../build/gf_eth_monitor_worker_inspector"
		]
		p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
			env={
				"GF_SENTRY_ENDPOINT":          os.environ["GF_SENTRY_ENDPOINT"],
				"GF_PY_PLUGINS_BASE_DIR_PATH": p_py_plugins_base_dir_path_str,
			})

		t = threading.Thread(target=gf_test_utils.read_process_stdout, args=(p.stdout, "gf_eth_monitor__worker_inspector", "yellow"))
		t.start()

		return p

	#-------------------------------------------------------------


	p = start_geth()
	p__worker_inspector = start_worker_inspector()

	import time
	time.sleep(10) # give time to Geth to initialize and send some events

	print("done sleeping")

	return p, p__worker_inspector

