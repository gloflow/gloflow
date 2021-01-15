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

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import subprocess
import threading

from colored import fg, bg, attr

import gf_test_utils

#--------------------------------------------------
def run(p_aws_region_str):
	
	print(f"    {fg('green')}TEST MASTER_GET_BLOCK{attr(0)} >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	#--------------------------------------------------
	def start_master():
		bin_str = f"{modd_str}/../../build/gf_eth_monitor"
		cmd_lst = [
			bin_str,
			"start", "service",
			f"--config={modd_str}/../../config/gf_eth_monitor.yaml"
		]



		print("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDd")
		print(os.environ["GF_SENTRY_ENDPOINT"])


		p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
			env={
				"AWS_REGION":            p_aws_region_str,
				"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
				"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
				"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"],
				"GF_SENTRY_ENDPOINT":    os.environ["GF_SENTRY_ENDPOINT"],

				"GF_EVENTS_CONSUME":        "false",
				"GF_WORKERS_AWS_DISCOVERY": "false" # use the localy started worker_inspector, at test startup
			})


		t = threading.Thread(target=gf_test_utils.read_process_stdout, args=(p.stdout, "gf_eth_monitor", "magenta"))
		t.start()


		return p

	#--------------------------------------------------


	p = start_master()

	import time
	time.sleep(10)


	#--------------------------------------------------
	def test():
		import requests



		print("MAKINT TEST CLIENT REQUEST")
		url_str = "http://127.0.0.1:4050/gfethm/v1/block?b=100"

		print(url_str)

		r = requests.get(url_str)


		print(r.text)
		print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> TEST COMPLETE ----------------")

	#--------------------------------------------------
	test()
	
	p.terminate()