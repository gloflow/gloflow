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

import os
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
# GET_SELF_IP

def get_self_ip():

	#---------------------------------------------------
	# IMPORTANT!! - this approach works most of the time.
	#               in some providers or networks DNS traffic might be routed from different exit routers
	#               (access points) than regular http traffic. 
	#               in thoise situations DNS nameserver that is pinged will return IP that is not appropriate.
	def dns_method():
		target_namespace_server_str = "ns1.google.com"

		# "-4" - force discovery of ipv4 address, since some ISP's by default return ipv6
		cmd_str = '''dig -4 TXT +short o-o.myaddr.l.google.com @%s | awk -F'"' '{ print $2}' '''%(target_namespace_server_str)
		print(cmd_str)

		r           = delegator.run(cmd_str)
		self_ip_str = r.out.strip()
		# assert len(self_ip_str.split(".")) == 4
		return self_ip_str

	#---------------------------------------------------
	def extern_service_ipinfoio_method():

		#---------------------------------------------------
		def get_remote():

			cmd_str = "curl http://ipinfo.io"
			print(cmd_str)
			
			r = delegator.run(cmd_str)

			self_ip_str = json.loads(r.out)["ip"]

			f=open("self_ip_cache.txt", "w")
			f.write(self_ip_str)


			return self_ip_str

		#---------------------------------------------------

		ip_str = get_remote()
		return ip_str
	
		

	#---------------------------------------------------
	def extern_service_ifconfigme_method():
		r = delegator.run("curl ifconfig.me")
		ip_str = r.out
		return ip_str

	#---------------------------------------------------
	def discover_ip():


		# IMPORTANT!! - in some networks in some countries many exit points are used dynamically for the same
		#               access point (mobile router); and each of these methods determines a different IP.
		#               so we're executing them all here and returning them all.
		dns__self_ip_str        = dns_method()
		ipinfo__self_ip_str     = extern_service_ipinfoio_method()
		ifconfigme__self_ip_str = extern_service_ifconfigme_method()
		print(f"DISCOVERED - self IPs:")
		print(f"dns         - {dns__self_ip_str}")
		print(f"ipinfo.io   - {ipinfo__self_ip_str}")
		print(f"ifconfig.me - {ifconfigme__self_ip_str}")
		
		# some users of this function expect the list of self-ips to be
		# unique and will cause errors otherwise.
		ips_no_dups_lst = list(set([
			dns__self_ip_str,
			ipinfo__self_ip_str,
			ifconfigme__self_ip_str
		]))

		return ips_no_dups_lst
	
	#---------------------------------------------------
	


	# CACHE
	cache_path_str = "self_ip_cache.txt"
	if os.path.isfile(cache_path_str):
		f=open("self_ip_cache.txt", "r")
		ips_str = f.read()
		f.close()



		ips_lst = []

		for ip_str in ips_str.split("\n"):

			# if text is a valid IP 
			if len(ip_str.split(".")) == 4:
				ips_lst.append(ip_str)

		print("CACHED - self IP:")
		print(ips_lst)

		return ips_lst

	# DISCOVER
	else:
		ips_lst = discover_ip()
			
		return ips_lst



	