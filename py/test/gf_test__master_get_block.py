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

# import subprocess
# import threading

from colored import fg, bg, attr

sys.path.append("%s/../utils"%(modd_str))
import gf_core_cli

#--------------------------------------------------
def run(p_test_ci_bool, p_aws_region_str):
	
	print(f"    {fg('green')}TEST MASTER_GET_BLOCK{attr(0)} >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	#--------------------------------------------------
	def start_master():
		bin_str = f"{modd_str}/../../build/gf_eth_monitor"
		cmd_lst = [
			bin_str,
			"start", "service",
			f"--config={modd_str}/../../config/gf_eth_monitor.yaml"
		]



		
		print(f'GF_SENTRY_ENDPOINT - {os.environ["GF_SENTRY_ENDPOINT"]}')
		print(f'GF_AWS_SQS_QUEUE   - {os.environ["GF_AWS_SQS_QUEUE"]}')

		# p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		# 	env={
		# 		"AWS_REGION":            p_aws_region_str,
		# 		"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
		# 		"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
		# 		"GF_SENTRY_ENDPOINT":    os.environ["GF_SENTRY_ENDPOINT"],
		# 		"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"],
		# 		"GF_EVENTS_CONSUME":        "false",
		# 		"GF_WORKERS_AWS_DISCOVERY": "false" # "true" # "false" # use the localy started worker_inspector, at test startup
		# 	})
		# t = threading.Thread(target=gf_test_utils.read_process_stdout, args=(p.stdout, "gf_eth_monitor", "magenta"))
		# t.start()


		p = gf_core_cli.run__view_realtime(cmd_lst, {
				"AWS_REGION":            p_aws_region_str,
				"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
				"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
				"GF_SENTRY_ENDPOINT":    os.environ["GF_SENTRY_ENDPOINT"],
				"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"],

				"GF_EVENTS_CONSUME":        "false",
				"GF_WORKERS_AWS_DISCOVERY": "false" # "true" # "false" # use the localy started worker_inspector, at test startup
			},
			"gf_eth_monitor", "green")

		return p

	#--------------------------------------------------
	
	p = start_master()

	# give external services time to process, only in local test mode
	if not p_test_ci_bool:
		import time
		time.sleep(10)

	#--------------------------------------------------
	def test__master():
		import requests



		block_num_int = 10

		print("MAKING TEST CLIENT REQUEST")
		url_str = f"http://127.0.0.1:4050/gfethm/v1/block?b={block_num_int}"

		print(url_str)

		r = requests.get(url_str)


		print(r.text)


		import json
		r_map = json.loads(r.text)

		assert r_map["status_str"] == "OK"
		assert "data" in r_map.keys()
		assert "block_from_workers_map" in r_map["data"]
		assert "miners_map" in r_map["data"]
		assert isinstance(r_map["data"]["miners_map"], dict)
		
		for worker_host_str, block_info_map in r_map["data"]["block_from_workers_map"].items():
			assert isinstance(worker_host_str, str)
			assert isinstance(block_info_map, dict)


			print("BLOCK INFO ----------------")
			print(block_info_map)

			assert "block_num_int"     in block_info_map.keys()
			assert "gas_used_uint"     in block_info_map.keys()
			assert "gas_limit_uint"    in block_info_map.keys()
			assert "coinbase_addr_str" in block_info_map.keys()
			assert "txs_lst"           in block_info_map.keys()
			assert "txs_hashes_lst"    in block_info_map.keys()
			assert "time_uint"         in block_info_map.keys()
			assert "block"             in block_info_map.keys()
		



			# assert block_num_int == block_info_map["block_num_int"]


		print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> TEST COMPLETE ----------------")

	#--------------------------------------------------
	test__master()


	# give external services time to process, only in local test mode
	if not p_test_ci_bool:
		time.sleep(20)

	p.terminate()