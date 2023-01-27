# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
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
	
	# CLI
	args_map = parse_args()

	#-----------------------
	# LOAD_DATA
	serialized_output_file_str = args_map["serialized_output_file_str"]
	assert not serialized_output_file_str == None

	state_history_file_str = args_map["state_history_file_str"]
	assert not state_history_file_str == None

	#--------------------------------------------------
	def load(p_path_str):
		f = open(p_path_str, "r")
		data = json.loads(f.read())
		f.close()
		return data
	
	#--------------------------------------------------

	serialized_output_lst = load(serialized_output_file_str)
	state_history_lst     = load(state_history_file_str)

	#-----------------------



	


	i=0
	for program_state_history_lst in state_history_lst:

		for s in program_state_history_lst:
			print(f"{i} >> x {s['x_f']} time {s['creation_unix_time_f']}")
			i+=1
	
	j=0
	for program_output_lst in serialized_output_lst:

		for o in program_output_lst:
			if "type_str" in o.keys() and o["type_str"] == "cube":
				print(f"{j} cube x {o['props_map']['x_f']}")
				j+=1


	print(f"states # {len(state_history_lst)}")
	print(f"cubes #  {j}")

	print("done...")

	#-----------------------
	# OUTPUT
	out_map = {}
	print(f"GF_OUT:{json.dumps(out_map)}")

	#-----------------------

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
	#----------------------------
	arg_parser.add_argument("-serialized_output_file", action = "store", default=None,
		help = "file path for the serialized_output json file")

	arg_parser.add_argument("-state_history_file", action = "store", default=None,
		help = "file path for the state_history json file")

	#----------------------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"serialized_output_file_str": args_namespace.serialized_output_file,
		"state_history_file_str":     args_namespace.state_history_file
	}

#--------------------------------------------------
main()