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



import subprocess

from colored import fg, bg, attr

import gf_test_utils

#--------------------------------------------------
def run(p_aws_region_str):

	
	print(f"    {fg('green')}TEST SQS_PROCESS{attr(0)} >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")







	bin_str = f"{modd_str}/../../build/gf_eth_monitor"
	cmd_lst = [
		bin_str,
		"test", "worker_event_process",
		f"--config={modd_str}/../../config/gf_eth_monitor.yaml"
	]
	p = subprocess.Popen(cmd_lst, shell=False, stdout=subprocess.PIPE, bufsize=1,
		env={
			"AWS_REGION":            p_aws_region_str,
			"AWS_ACCESS_KEY_ID":     os.environ["AWS_ACCESS_KEY_ID"],
			"AWS_SECRET_ACCESS_KEY": os.environ["AWS_SECRET_ACCESS_KEY"],
			"GF_AWS_SQS_QUEUE":      os.environ["GF_AWS_SQS_QUEUE"]
		})


	gf_test_utils.read_process_stdout(p.stdout, "gf_eth_monitor", "green")