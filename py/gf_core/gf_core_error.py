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

import traceback

#--------------------------------------------------------
# p_surrounding_context_attribs_tpl - order matters

def handle_exception(p_exception,
	p_formated_msg_str,
	p_surrounding_context_attribs_tpl,
	p_log_fun):
	p_log_fun('FUN_ENTER', 'gf_error.handle_exception()')
	assert isinstance(p_exception,Exception)

	if p_formated_msg_str                == None or \
	   p_surrounding_context_attribs_tpl == None:
	
		msg_str = ''
	else:
		assert isinstance(p_surrounding_context_attribs_tpl, tuple)
		msg_str = p_formated_msg_str%(p_surrounding_context_attribs_tpl)
	
	p_log_fun('INFO', 'p_exception.message:%s'%(p_exception.message))
	
	
	with_trace_msg_str = '''
		%s
		exception args:%s
		exception msg :%s
		trace         :%s'''%(msg_str,
			p_exception.args,
			p_exception.message,
			traceback.format_exc())

	p_log_fun('ERROR', with_trace_msg_str)