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

import argparse
import json

#--------------------------------------------------
def main():
	

	args_map = parse_args()



	contract_addr_str = args_map["contract_addr_str"]
	assert not contract_addr_str == None


	
	geth__port_int = 8545
	geth__host_str = "127.0.0.1"


	print(f"contract_addr - {contract_addr_str}")



	out_map = {}
	print(f"GF_OUT:{json.dumps(out_map)}")


#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
	#----------------------------
	# CONTRACT_ADDRESS
	arg_parser.add_argument("-contract_addr", action = "store", default=None,
		help = "address of the target contract")

	#-------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"contract_addr_str": args_namespace.contract_addr
	}

#--------------------------------------------------
main()