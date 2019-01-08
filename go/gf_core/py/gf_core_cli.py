# GloFlow media management/publishing system
# Copyright (C) 2019 Ivan Trajkovic
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

import sys
import argparse
#-----------------------------------------------------
#IMPORTANT!! - these arguments(service_info parameters) have precendence over
#              what is returned by services get_service_info() function

#->:Dict(dynamic_service_info_dict)
def parse_args(p_cmd_line_args_defs_map, p_log_fun):
	p_log_fun('FUN_ENTER','gf_core_cli.parse_args()')
	assert isinstance(p_cmd_line_args_defs_map,dict)
	
	passed_in_args_lst = sys.argv[1:]
	p_log_fun('INFO','passed in args:%s'%(passed_in_args_lst))
	
	#RawTextHelpFormatter - so that newlines in the help text are rendered when "-h" option 
	#                       is passed on the command line
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#load up all command line argument definitions
	for arg_name_str,arg_def_map in p_cmd_line_args_defs_map.items():
		arg_default  = arg_def_map['default']
		arg_help_str = arg_def_map['help']
		
		arg_parser.add_argument('-%s'%(arg_name_str), 
								action  = "store",
								default = arg_default,
								help    = arg_help_str)
	#:Namespace
	args_namespace = arg_parser.parse_args(passed_in_args_lst)
	
	#extracts command line arguments from args_namespace 
	#(only the expected arguments), and repacks them into dynamic_service_info_dict
	extracted_args_map = {}
	for arg_name_str,_ in p_cmd_line_args_defs_map.items():
		extracted_args_map[arg_name_str] = getattr(args_namespace,arg_name_str)
		
	return extracted_args_map
#-----------------------------------------------------
#->:Bool
def confirm(p_prompt_str, p_resp=False):
    	
    if p_prompt_str is None:
        p_prompt_str = 'Confirm'

    if p_resp:
        p_prompt_str = '%s %s|%s: ' % (p_prompt_str, 'y', 'n')
    else:
        p_prompt_str = '%s %s|%s: ' % (p_prompt_str, 'n', 'y')
        
    while True:
        answer_str = raw_input(p_prompt_str)
        if not answer_str:
            return p_resp
        if answer_str not in ['y', 'Y', 'n', 'N']:
            print 'please enter y or n.'
            continue
        if answer_str == 'y' or answer_str == 'Y':
            return True
        if answer_str == 'n' or answer_str == 'N':
            return False