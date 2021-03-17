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

from colored import fg, bg, attr

sys.path.append("%s/../utils"%(modd_str))
import gf_core_cli

import gf_test__sqs_process
import gf_test__master_get_block

#--------------------------------------------------
def run_go(p_test_ci_bool):

	aws_region_str = "us-east-1"
	py_plugins_base_dir_path_str = f"{modd_str}/../plugins"
	
	# CI
	if p_test_ci_bool:
		
		p__worker_inspector = init_ci(py_plugins_base_dir_path_str,
			aws_region_str)

	#--------------------------------------------------
	def run__tests():

		cmd_lst = [
			"go test", "-v"
		]
		p = gf_core_cli.run__view_realtime(cmd_lst, {
				"GF_GETH_HOST":                os.environ["GF_GETH_HOST"],
				"GF_SENTRY_ENDPOINT":          os.environ["GF_SENTRY_ENDPOINT"],
				"GF_PY_PLUGINS_BASE_DIR_PATH": p_py_plugins_base_dir_path_str,
			},
			"gf_eth_monitor__worker_inspector", "cyan")
	
	#--------------------------------------------------
	run__tests()

#--------------------------------------------------
def run_py(p_test_ci_bool):
	
	assert "AWS_ACCESS_KEY_ID" in os.environ.keys()
	assert "AWS_SECRET_ACCESS_KEY" in os.environ.keys()
	assert "GF_AWS_SQS_QUEUE" in os.environ.keys()
	assert "GF_SENTRY_ENDPOINT" in os.environ.keys()


	aws_region_str               = "us-east-1"
	py_plugins_base_dir_path_str = f"{modd_str}/../plugins"

	# CI
	if p_test_ci_bool:
		
		p__worker_inspector = init_ci(py_plugins_base_dir_path_str,
			aws_region_str)

	# LOCAL
	else:

		test_gf_geth_path_str = f"{modd_str}/../../../gloflow_go-ethereum/build/bin/geth"
		data_dir_path_str     = f"{modd_str}/test_data_geth"
		assert os.path.isfile(test_gf_geth_path_str)
		assert os.path.isdir(py_plugins_base_dir_path_str)

		p__geth, p__worker_inspector = init_local(test_gf_geth_path_str,
			data_dir_path_str,
			py_plugins_base_dir_path_str,
			aws_region_str)


	#------------------------
	# TESTS
	gf_test__sqs_process.run(p_test_ci_bool, aws_region_str)
	gf_test__master_get_block.run(p_test_ci_bool, aws_region_str)

	#------------------------


	# CI
	if p_test_ci_bool:
		p__worker_inspector.terminate()

	# LOCAL
	else:
		p__geth.terminate()
		p__worker_inspector.terminate()

#--------------------------------------------------
def init(p_py_plugins_base_dir_path_str,
	p_aws_region_str):
	True

#--------------------------------------------------
def init_ci(p_py_plugins_base_dir_path_str,
	p_aws_region_str,
	p_worker_inspector_bool = True):

	# DB
	init_db()

	if p_worker_inspector_bool:
		#-------------------------------------------------------------
		def start_worker_inspector():
			cmd_lst = [
				f"{modd_str}/../../build/gf_eth_monitor_worker_inspector"
			]
			p = gf_core_cli.run__view_realtime(cmd_lst, {
					"GF_GETH_HOST":                os.environ["GF_GETH_HOST"],
					"GF_SENTRY_ENDPOINT":          os.environ["GF_SENTRY_ENDPOINT"],
					"GF_PY_PLUGINS_BASE_DIR_PATH": p_py_plugins_base_dir_path_str,
				},
				"gf_eth_monitor__worker_inspector", "yellow")

			return p

		#-------------------------------------------------------------

		p__worker_inspector = start_worker_inspector()
		return p__worker_inspector

#--------------------------------------------------
def init_local(p_geth_bin_path_str,
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

		# # When shell=True the shell is the child process, and the commands are its children.
		# # So any SIGTERM or SIGKILL will kill the shell but not its child processes.
		# # The best way I can think of is to use shell=False, otherwise when you kill
		# # the parent shell process, it will leave a defunct shell process.
		# # CMD also has to be a list here, since its not being passed in as a string
		# # to the child shell.
		# geth_p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		# 	env={
		# 		"AWS_REGION":            p_aws_region_str,
		# 		"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
		# 		"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
		# 		"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"]
		# 	})
		# t = threading.Thread(target=gf_test_utils.read_process_stdout, args=(geth_p.stdout, "geth", "cyan"))
		# t.start()

		geth_p = gf_core_cli.run__view_realtime(cmd_lst, {
				"AWS_REGION":            p_aws_region_str,
				"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
				"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
				"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"]
			},
			"geth", "cyan")

		return geth_p

	#-------------------------------------------------------------
	def start_worker_inspector():
		cmd_lst = [
			f"{modd_str}/../../build/gf_eth_monitor_worker_inspector"
		]

		# p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		# 	env={
		# 		"GF_GETH_HOST":                "127.0.0.1",
		# 		"GF_SENTRY_ENDPOINT":          os.environ["GF_SENTRY_ENDPOINT"],
		# 		"GF_PY_PLUGINS_BASE_DIR_PATH": p_py_plugins_base_dir_path_str,
		# 	})
		# t = threading.Thread(target=gf_test_utils.read_process_stdout, args=(p.stdout, "gf_eth_monitor__worker_inspector", "yellow"))
		# t.start()

		gf_core_cli.run__view_realtime(cmd_lst, {
				"GF_GETH_HOST":                "127.0.0.1",
				"GF_SENTRY_ENDPOINT":          os.environ["GF_SENTRY_ENDPOINT"],
				"GF_PY_PLUGINS_BASE_DIR_PATH": p_py_plugins_base_dir_path_str,
			},
			"gf_eth_monitor__worker_inspector", "yellow")

		return p

	#-------------------------------------------------------------


	p = start_geth()
	p__worker_inspector = start_worker_inspector()

	import time
	time.sleep(10) # give time to Geth to initialize and send some events

	print("done sleeping")

	return p, p__worker_inspector

#--------------------------------------------------
def init_db():
	

	ec2_abi_str = '''[
		{
			"constant": true,
			"inputs": [],
			"name": "name",
			"outputs": [
				{
					"name": "",
					"type": "string"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_spender",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "approve",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "totalSupply",
			"outputs": [
				{
					"name": "",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_from",
					"type": "address"
				},
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transferFrom",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "decimals",
			"outputs": [
				{
					"name": "",
					"type": "uint8"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
				{
					"name": "_owner",
					"type": "address"
				}
			],
			"name": "balanceOf",
			"outputs": [
				{
					"name": "balance",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "symbol",
			"outputs": [
				{
					"name": "",
					"type": "string"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transfer",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
				{
					"name": "_owner",
					"type": "address"
				},
				{
					"name": "_spender",
					"type": "address"
				}
			],
			"name": "allowance",
			"outputs": [
				{
					"name": "",
					"type": "uint256"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"payable": true,
			"stateMutability": "payable",
			"type": "fallback"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "owner",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "spender",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Approval",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "from",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "to",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "Transfer",
			"type": "event"
		}
	]'''

	import json
	def_lst = json.loads(ec2_abi_str)
	erc20_abi_map = {
		"type_str": "erc20",
		"def_lst":  def_lst
	}



	import pymongo
	db_name_str   = "gf_eth_monitor"
	coll_name_str = "gf_eth_meta__contracts_abi"
	client = pymongo.MongoClient('localhost', 27017)
	client[db_name_str][coll_name_str].insert(erc20_abi_map)