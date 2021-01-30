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

import time
import clint

#----------------------------------------------
def get_log_fun(p_log_fun):
	# p_log_fun("FUN_ENTER", "gf_core_logger.get_log_fun()")
	
	#----------------------------------------------	
	def log_color_display_fun(p_group, p_msg):
		
		t = str(time.time())			 
		if p_group == 'FUN_ENTER':
			clint.textui.puts(t+':'+clint.textui.colored.yellow('FUN_ENTER')+':'+p_msg)
			
		elif p_group == 'ERROR':
			clint.textui.puts(t+':'+clint.textui.colored.red(p_group)+':'+p_msg)
			
		elif p_group == 'INFO':
			clint.textui.puts(t+':'+clint.textui.colored.green(p_group)+':'+ clint.textui.colored.green(p_msg))
			
		elif p_group == 'INFO_USR':
			clint.textui.puts(t+':'+clint.textui.colored.magenta(p_group)+':'+ clint.textui.colored.magenta(p_msg))
			
		# log message is in some way related to external-systems data
		elif p_group == 'EXTERN':
			clint.textui.puts(t+':'+clint.textui.colored.blue(p_group)+':'+p_msg)
			
		# if 'TEST' is anywhere in the group string
		# elif p_group.find('TEST'):
		#	clint.textui.puts(clint.textui.colored.cyan(p_group)+':'+p_msg)
			
		#--------------
		# JAVASCRIPT FORMATING
		elif p_group == 'JS:FUN_ENTER':
			clint.textui.puts(t+':'+'       '+clint.textui.colored.yellow('FUN_ENTER')+':'+p_msg)
			
		elif p_group == 'JS:INFO':
			clint.textui.puts(t+':'+'       '+clint.textui.colored.green(p_group)+':'+ clint.textui.colored.green(p_msg))

	#----------------------------------------------
			
	return log_color_display_fun