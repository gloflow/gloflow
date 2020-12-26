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

#--------------------------------------------------
def test__sqs_process(p_aws_region_str):

	

	bin_str = f"{modd_str}/../../build/gf_eth_monitor"
	cmd_lst = [
		bin_str,
		"test", "eth_node_event_process",
		f"--config={modd_str}/../../config/gf_eth_monitor.yaml"
	]
	geth_p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		env={
			"AWS_REGION":            p_aws_region_str,
			"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
			"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
			"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"]
		})


	read_process_stdout(geth_p.stdout, "gf_eth_monitor")


#--------------------------------------------------
def run():
	
	aws_region_str        = "us-east-1"
	test_gf_geth_path_str = f"{modd_str}/../../../gloflow_go-ethereum/build/bin/geth"
	assert os.path.isfile(test_gf_geth_path_str)



	data_dir_path_str = f"{modd_str}/test_data_geth"

	#-------------------------------------------------------------
	# START_GETH
	def start_geth():
	  
		print("%srunning GETH...%s"%(fg("yellow"), attr(0)))
		
		assert "AWS_ACCESS_KEY_ID" in os.environ.keys()
		assert "AWS_SECRET_ACCESS_KEY" in os.environ.keys()
		assert "GF_AWS_SQS_QUEUE" in os.environ.keys()
		
		cmd_lst = [
			test_gf_geth_path_str,
			f"--datadir={data_dir_path_str}",
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
				"AWS_REGION":            aws_region_str,
				"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
				"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
				"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"]
			})

		t = threading.Thread(target=read_process_stdout, args=(geth_p.stdout, "geth"))
		t.start()

		return geth_p

	#-------------------------------------------------------------

	p = start_geth()
	
	import time
	time.sleep(10) # give time to Geth to initialize and send some events

	print("done sleeping")



	test__sqs_process(aws_region_str)

	#-------------------------------------------------------------
	def cleanup():
		p.terminate()

	#-------------------------------------------------------------

	cleanup()

#-------------------------------------------------------------
def read_process_stdout(p_out, p_type_str):

	for line in iter(p_out.readline, b''):
		
		header_color_str = None
		
		if p_type_str == "geth":
			header_color_str = fg("cyan")
		elif p_type_str == "gf_eth_monitor":
			header_color_str = fg("yellow")

		line_str = line.strip().decode("utf-8")

		# ERROR
		if "ERROR" in line_str or "error" in line_str:
			print("%s%s:%s%s%s%s"%(header_color_str, p_type_str, attr(0), bg("red"), line_str, attr(0)))
		else:
			print("%s%s:%s%s"%(header_color_str, p_type_str, attr(0), line_str))

	p_out.close()