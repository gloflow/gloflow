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

import logging
import subprocess
import signal
import multiprocessing

import json
import delegator

#---------------------------------------------------
def run_cmd_in_os_proc(p_cmd_str, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_core_utils.run_cmd_in_os_proc()")
	
	#---------------------------------------------------
	def run_process():
		p_log_fun("FUN_ENTER", "gf_core_utils.run_cmd_in_os_proc().run_process()")

		p_log_fun("INFO", "RUN CMD IN OS_PROCESS")
		p_log_fun("INFO", "p_cmd_str - %s"%(p_cmd_str))

		p = subprocess.Popen(p_cmd_str,
			shell   = True,
			stdout  = subprocess.PIPE,
			bufsize = 1)
		assert isinstance(p,subprocess.Popen)
		#---------------------------------------------------
		# IMPORTANT!! - workers cleanup on parent shutdown

		def handle_signal_terminate(p_signum, p_frame):
			p_log_fun("INFO", "+++ ++ -- SIGNAL SIGTERM RECEIVED - gf_images_main_service")

			p.terminate()

		import signal
		signal.signal(signal.SIGTERM, handle_signal_terminate)

		#---------------------------------------------------

		bin_str = os.path.basename(p_cmd_str.split(" ")[0])
		print(envoy.run("ps -e | grep %s"%(bin_str)).std_out)

		p.wait() # block this process and let the child run
	#---------------------------------------------------

	multiprocessing.log_to_stderr(logging.DEBUG)
	parent_conn, child_conn = multiprocessing.Pipe()
	process                 = multiprocessing.Process(target = run_process)
	process.start()

#---------------------------------------------------
def get_self_ip():

	#---------------------------------------------------
	# IMPORTANT!! - this approach works most of the time.
	#               in some providers or networks DNS traffic might be routed from different exit routers
	#               (access points) than regular http traffic. 
	#               in thoise situations DNS nameserver that is pinged will return IP that is not appropriate.
	def dns_method():
		target_namespace_server_str = "ns1.google.com"
		cmd_str = '''dig TXT +short o-o.myaddr.l.google.com @%s | awk -F'"' '{ print $2}' '''%(target_namespace_server_str)
		print(cmd_str)

		r           = delegator.run(cmd_str)
		self_ip_str = r.out.strip()
		assert len(self_ip_str.split(".")) == 4
		return self_ip_str

	#---------------------------------------------------
	def extern_service_method():

		#---------------------------------------------------
		def get_remote():

			cmd_str = "curl http://ipinfo.io"
			print(cmd_str)
			
			r           = delegator.run(cmd_str)
			self_ip_str = json.loads(r.out)["ip"]

			f=open("self_ip_cache.txt", "w")
			f.write(self_ip_str)


			return self_ip_str

		#---------------------------------------------------
		f=open("self_ip_cache.txt", "r")
		ip_str = f.read()
		f.close()

		# if text is a valid IP 
		if len(ip_str.split(".")) == 4:
			return ip_str
		else:
			ip_str = get_remote()
			return ip_str

	#---------------------------------------------------
	
	# self_ip_str = dns_method()
	self_ip_str = extern_service_method()
	print(f"self IP: {self_ip_str}")

	return self_ip_str