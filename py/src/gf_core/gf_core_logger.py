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


from colored import fg, attr
import time

#----------------------------------------------
def get_log_fun(p_log_fun):

    def log_color_display_fun(p_group, p_msg):
        t = str(time.time())
        
        if p_group == 'FUN_ENTER':
            print(f"{t}:{fg('yellow')}FUN_ENTER{attr('reset')}:{p_msg}")
        
        elif p_group == 'ERROR':
            print(f"{t}:{fg('red')}{p_group}{attr('reset')}:{p_msg}")
        
        elif p_group == 'INFO':
            print(f"{t}:{fg('green')}{p_group}{attr('reset')}:{fg('green')}{p_msg}{attr('reset')}")
        
        elif p_group == 'INFO_USR':
            print(f"{t}:{fg('magenta')}{p_group}{attr('reset')}:{fg('magenta')}{p_msg}{attr('reset')}")

		# log message is in some way related to external-systems data
        elif p_group == 'EXTERN':
            print(f"{t}:{fg('blue')}{p_group}{attr('reset')}:{p_msg}")

		# if 'TEST' is anywhere in the group string
		# elif p_group.find('TEST'):
		#	clint.textui.puts(clint.textui.colored.cyan(p_group)+':'+p_msg)
        
		#--------------
		# JAVASCRIPT FORMATING
        elif p_group == 'JS:FUN_ENTER':
            print(f"{t}:       {fg('yellow')}FUN_ENTER{attr('reset')}:{p_msg}")
        
        elif p_group == 'JS:INFO':
            print(f"{t}:       {fg('green')}{p_group}{attr('reset')}:{fg('green')}{p_msg}{attr('reset')}")

		#--------------
            
    return log_color_display_fun