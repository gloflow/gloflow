# GloFlow application and media management/publishing platform
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

import sys, os
import argparse
import subprocess
import threading

from colored import fg, bg, attr
import delegator

#---------------------------------------------------
def run__view_realtime(p_cmd_lst,
	p_env_map,
	p_view__type_str,
	p_view__color_str):

	print(" ".join(p_cmd_lst))
	
	# When shell=True the shell is the child process, and the commands are its children.
	# So any SIGTERM or SIGKILL will kill the shell but not its child processes.
	# The best way I can think of is to use shell=False, otherwise when you kill
	# the parent shell process, it will leave a defunct shell process.
	# CMD also has to be a list here, since its not being passed in as a string
	# to the child shell.
	p = subprocess.Popen(p_cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		env=p_env_map)

	t = threading.Thread(target=read_process_std_stream, args=(p.stdout, p_view__type_str, p_view__color_str))
	t.start()

	return p

#---------------------------------------------------
# RUN
def run(p_cmd_str,
	p_env_map = {},
	p_exit_on_fail_bool = True,
	p_print_cmd_str     = True):

	# env map has to contains all the parents ENV vars as well
	p_env_map.update(os.environ)

	if p_print_cmd_str:
		print(f"{fg('yellow')}cmd{attr(0)} >>> {fg('green')}{p_cmd_str}{attr(0)}")

	p = subprocess.Popen(p_cmd_str,
		env     = p_env_map,
		shell   = True,
		stdout  = subprocess.PIPE,
		stderr  = subprocess.PIPE,
		bufsize = 1)

	t_o = threading.Thread(target=read_process_std_stream, args=(p.stdout, "stdout", "green"))
	t_o.start()

	t_e = threading.Thread(target=read_process_std_stream, args=(p.stderr, "stderr", "yellow"))
	t_e.start()

	p.wait()
	
	# wait for stdout/stderr printing threads to complete as well before returning from this function
	t_o.join()
	t_e.join()

	if p_exit_on_fail_bool:
		if not p.returncode == 0:

			print(f"ERROR!! - shell CMD ({p_cmd_str}) failed!")
			exit(p.returncode)

	return "", "", p.returncode

#-------------------------------------------------------------
def read_process_std_stream(p_std_stream,
	p_view_type_str,
	p_view_color_str):

	for line in iter(p_std_stream.readline, b''):
		
		header_color_str = fg(p_view_color_str)
		line_str         = line.strip().decode("utf-8")

		if "ERROR" in line_str or "error" in line_str:
			print("%s%s:%s%s%s%s"%(header_color_str, p_view_type_str, attr(0), bg("red"), line_str, attr(0)))
		else:
			print("%s%s:%s%s"%(header_color_str, p_view_type_str, attr(0), line_str))

	p_std_stream.close()

#---------------------------------------------------
# DEPRECATED!! - move all users over to run()
def run_cmd(p_cmd_str,
	p_env_map           = None,
	p_print_output_bool = True):
	
	if p_print_output_bool:
		print(p_cmd_str)
	
	if not p_env_map == None:
		assert isinstance(p_env_map, dict)
		r = delegator.run(p_cmd_str, env=p_env_map)
	else:
		r = delegator.run(p_cmd_str)

	o = ""
	e = ""

	# sometimes commands dont return any stdout
	if not r.out == "":
		o = r.out
		if p_print_output_bool:
			print(o)

	# sometimes commands dont return any stderr
	if not r.err == "":
		e = r.err
		if p_print_output_bool: print(e)
	
	return o, e, r.return_code

#-----------------------------------------------------
# DEPRECATED!! - all users of this function will be migrated to using Pythons stdlib argparse
#                directly. once they all migrate to this remove this function and this abstraction.
# IMPORTANT!! - these arguments(service_info parameters) have precendence over
#               what is returned by services get_service_info() function

def parse_args(p_cmd_line_args_defs_map, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_core_cli.parse_args()")
	assert isinstance(p_cmd_line_args_defs_map, dict)
	
	passed_in_args_lst = sys.argv[1:]
	p_log_fun("INFO", "passed in args:%s"%(passed_in_args_lst))
	
	# RawTextHelpFormatter - so that newlines in the help text are rendered when "-h" option 
	#                        is passed on the command line
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#load up all command line argument definitions
	for arg_name_str, arg_def_map in p_cmd_line_args_defs_map.items():
		arg_default  = arg_def_map["default"]
		arg_help_str = arg_def_map["help"]
		

		arg_parser.add_argument("-%s"%(arg_name_str), 
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
def confirm(p_prompt_str, p_resp=False):
	prompt_str = None
	if p_prompt_str is None:
		prompt_str = "Confirm"

	if p_resp:
		prompt_str = "%s %s|%s: "%(p_prompt_str, "y", "n")
	else:
		prompt_str = "%s %s|%s: "%(p_prompt_str, "n", "y")
		
	while True:

		answer_str = input(prompt_str)
		if not answer_str:
			return p_resp
		if answer_str not in ["y", "Y", "n", "N"]:
			print("please enter y or n.")
			continue
		if answer_str == "y" or answer_str == "Y":
			return True
		if answer_str == "n" or answer_str == "N":
			return False