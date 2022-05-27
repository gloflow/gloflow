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

import svgwrite
from svgwrite import cm, mm
from svgwrite.container import SVG

#--------------------------------------------------
def main():
	
	print("PY_PLUGIN - PLOT_TX_TRACE")
	#----------------------------
	# INPUT
	args_map = parse_args()
	


	l = sys.stdin.readline()
	print(">>>>>>>>>>>>>>>>>>>>>>>>>>")
	# print(l.strip())



	tx_id_str    = args_map["tx_id_str"]
	tx_trace_map = json.loads(l)
	assert isinstance(tx_trace_map, dict)

	#----------------------------



	plot_y = 50

	dwg      = svgwrite.Drawing(filename=f"{modd_str}/test.svg", debug=True)
	plot_svg = SVG((50, plot_y))
	dwg.add(plot_svg)


	plot_width_mm_int = 500
	
	x_ops_base      = 10
	x_stack_base    = 30
	x_memory_base   = x_stack_base + 18 
	x_gas_cost_base = x_memory_base + 10

	#----------------------------
	# LEGEND
	legend_svg = SVG((50, 10))
	dwg.add(legend_svg)

	legend_svg.add(dwg.text(f"tx ID - {tx_id_str}", (0*mm, 4*mm), font_size=8))

	legend_svg.add(dwg.text(f"stack",  (x_stack_base*mm, plot_y-12), font_size=8))
	legend_svg.add(dwg.text(f"memory", ((x_memory_base-3)*mm, plot_y-12), font_size=8))
	legend_svg.add(dwg.text(f"gas cost", (x_gas_cost_base*mm, plot_y-12), font_size=8))

	#----------------------------



	i=0

	memory_ops_lst = []
	call_ops_lst   = []
	logs_ops_lst   = []

	for op_map in tx_trace_map["opcodes_lst"]:
		

		# print("--------------------")
		# print(op_map)

		op_str       = op_map["op_str"].strip()
		gas_cost_int = int(op_map["gas_cost_uint"])
		# print(f"{op_str}-{gas_cost_int}")



		stack_lst  = op_map["stack_lst"]
		memory_lst = op_map["memory_lst"]



		

		x1 = x_gas_cost_base
		x2 = x1+gas_cost_int
		y  = i*8 # 2.4
		

		op_svg = SVG((0, y))
		plot_svg.add(op_svg)

		#----------------------------
		# OP_GAS_COST__LINE
		op_svg.add(dwg.line(start=(x1*mm, 1.25*mm), end=(x2*mm, 1.25*mm),
			stroke="green",
			stroke_width=3))

		#----------------------------
		# DEBUGGING - alignment line. used to align other elements in line.
		op_svg.add(dwg.line(start=((x_stack_base-1)*mm, 1.25*mm), end=(x1*mm, 1.25*mm),
			stroke="black",
			stroke_width=0.5))

		#----------------------------
		# OP__TEXT - text local coordinate system is at lower left corner,
		#            not upper-left like everything else.
		#            so positioning it a bit lower from 0,0 (in this case 0,2)
		op_txt = dwg.text(op_str, (x_ops_base*mm, 2*mm),
			font_size=8)
		op_svg.add(op_txt)

		
		# MSTORE
		if op_str == "MSTORE":
			y__local = 0.8
			op_svg.add(dwg.rect(insert=((x_ops_base-1.5)*mm, y__local*mm), size=(1*mm, 1*mm),
				fill='blue',
				stroke='black',
				stroke_width=0.5))


			y__global = y+y__local
			memory_ops_lst.append(y__global)


		# CALLDATASIZE/CALLVALUE/CALLER
		ops_call_lst = ["CALLDATASIZE", "CALLVALUE", "CALLER"]
		if op_str in ops_call_lst:
			
			x_call_rect__global = x_ops_base-2
			y__local            = 0.6
			op_svg.add(dwg.rect(insert=((x_call_rect__global)*mm, y__local*mm), size=(1.2*mm, 1.2*mm),
				fill='red',
				stroke='black',
				stroke_width=0.5))

			y__global = y+y__local
			call_ops_lst.append(y__global)

		# LOG
		if op_str.startswith("LOG"):
			
			x_log_rect__global = x_ops_base-2
			y__local = 0.6
			op_svg.add(dwg.rect(insert=(x_log_rect__global*mm, y__local*mm), size=(1*mm, 1*mm),
				fill='cyan',
				stroke='black',
				stroke_width=0.5))
			
			y__global = y+y__local
			logs_ops_lst.append(y__global)

		#----------------------------
		# STACK
		op_stack_g = op_svg.add(dwg.g(id='op_stack', stroke='blue'))
		j=0
		for s in stack_lst:

			x = x_stack_base + j*2
			op_stack_g.add(dwg.rect(insert=(x*mm, 0.5*mm), size=(1.5*mm, 1.5*mm),
				fill='yellow',
				stroke='black',
				stroke_width=0.5))

			j+=1

		#----------------------------
		# MEMORY
		op_memory_g = op_svg.add(dwg.g(id='op_memory', stroke='blue'))
		j=0
		for s in memory_lst:

			x = x_memory_base + j*2
			op_memory_g.add(dwg.rect(insert=(x*mm, 0.5*mm), size=(1.5*mm, 1.5*mm),
				fill='orange',
				stroke='black',
				stroke_width=0.5))

			j+=1

		#----------------------------

		i+=1

	
	#--------------------------------------------------
	# ARCHS
	def draw_archs():
		
		x_call_rect__global_px = 30

		i=0
		for y__global in call_ops_lst[:-1]:

			y__global__next = call_ops_lst[i+1]

			y__global__start_str = f"{x_call_rect__global_px} {int(y__global+4.4)}"
			y__global__end_str   = f"{x_call_rect__global_px} {int(y__global__next+4.4)}"
			y__global__control_point_str = f"0 {int(y__global+(y__global__next-y__global)/2)}"

			path_str = f"M {y__global__start_str} Q {y__global__control_point_str} {y__global__end_str}"

			plot_svg.add(dwg.path(d=path_str, fill="none", stroke="red", stroke_width=0.5))
			
			i+=1

	#--------------------------------------------------
	draw_archs()

	print("done drawing...")

	# FILE_SAVE
	if args_map["stdout_bool"]:

		svg_str = dwg.tostring()
		out_map = {"svg_str": svg_str}
		print(f"GF_OUT:{json.dumps(out_map)}")
	else:
		dwg.save()

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
	#----------------------------
	# TX_ID
	arg_parser.add_argument("-tx_id", action = "store", default=None,
		help = "hex ID of the target transaction")

	#----------------------------
	# STDOUT
	arg_parser.add_argument('-stdout', action='store_true')

	#----------------------------

	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"tx_id_str":   args_namespace.tx_id,
		"stdout_bool": args_namespace.stdout
	}

#--------------------------------------------------
main()