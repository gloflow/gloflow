

import logging
import subprocess
import signal
import multiprocessing
#---------------------------------------------------
def run_cmd_in_os_proc(p_cmd_str,
					p_log_fun):
	p_log_fun('FUN_ENTER','gf_core_utils.run_cmd_in_os_proc()')
	#---------------------------------------------------
	def run_process():
		p_log_fun('FUN_ENTER','gf_core_utils.run_cmd_in_os_proc().run_process()')

		p_log_fun("INFO","RUN CMD IN OS_PROCESS")
		p_log_fun('INFO','p_cmd_str - %s'%(p_cmd_str))

		p = subprocess.Popen(p_cmd_str,
						shell   = True,
						stdout  = subprocess.PIPE,
						bufsize = 1)
		assert isinstance(p,subprocess.Popen)
		#---------------------------------------------------
		#IMPORTANT!! - workers cleanup on parent shutdown

		def handle_signal_terminate(p_signum,
								p_frame):
			p_log_fun('INFO','+++ ++ -- SIGNAL SIGTERM RECEIVED - gf_images_main_service')

			p.terminate()

		import signal
		signal.signal(signal.SIGTERM, handle_signal_terminate)
		#---------------------------------------------------



		bin_str = os.path.basename(p_cmd_str.split(' ')[0])
		print envoy.run('ps -e | grep %s'%(bin_str)).std_out


		p.wait() #block this process and let the child run
	#---------------------------------------------------

	multiprocessing.log_to_stderr(logging.DEBUG)
	parent_conn, child_conn = multiprocessing.Pipe()
	process                 = multiprocessing.Process(target = run_process)
	process.start()