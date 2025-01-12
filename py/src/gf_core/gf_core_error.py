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
import pprint
import sentry_sdk
from colored import fg, attr

#--------------------------------------------------------
def create(p_msg_str,
	p_type_str,
	p_data_map,
	p_exception,
	p_sybsystem_name_str,
	p_log_fun,
	p_sentry_bool=True,
	p_reraise_bool=False):
	assert isinstance(p_msg_str, str)
	assert isinstance(p_type_str, str)
	assert isinstance(p_data_map, dict) or p_data_map == None
	assert isinstance(p_exception, Exception)
	assert isinstance(p_sybsystem_name_str, str)
	
	full_msg_str = f'''
		{fg('green')}{p_msg_str}{attr('reset')}
		type:           {fg('green')}{p_type_str}{attr('reset')}
		subsys:         {fg('green')}{p_sybsystem_name_str}{attr('reset')}
		exception args: {fg('green')}{p_exception.args}{attr('reset')}
		data:           {fg('green')}{pprint.pformat(p_data_map)}{attr('reset')}
		trace: {fg('green')}{traceback.format_exc()}{attr('reset')}
	'''
	p_log_fun('ERROR', full_msg_str)

	if p_sentry_bool:
		sentry_sdk.capture_exception(p_exception)
	
	if p_reraise_bool:
		raise p_exception