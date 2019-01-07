import os
import logging
import time
import subprocess
import multiprocessing
import signal

cwd_str = os.path.dirname(os.path.abspath(__file__))
#--------------------------------------------------------
#FIX!! - remove p_workers_various_params_dict from the argument list of start_new_os_process() arguments,
#        because start_new_os_process() is handled generically in gf_ops_app_runner.py, and extra non-standard
#        arguments should be avoided.

def start_new_os_process(p_log_fun,
				p_port_str                                   = '3050',
				p_mongo_db_name_str                          = 'test_db',
				p_images_store_local_dir_path_str            = '%s/../tests/images/output/original'%(cwd_str),
				p_images_thumbnails_store_local_dir_path_str = '%s/../tests/images/output/thumbnails'%(cwd_str),
				p_images_s3_bucket_name_str                  = 'gf--test',
				p_service_bin_path_str                       = os.path.abspath('%s/../../../../bin/gf_images_service'%(cwd_str))):
	p_log_fun('FUN_ENTER','gf_image_service_utils.start_new_os_process()')
	assert isinstance(p_port_str            ,basestring)
	assert isinstance(p_service_bin_path_str,basestring)

	#---------------------------------------------------
	def run_process(p_child_conn,
				p_log_fun):
		p_log_fun('FUN_ENTER','gf_image_service_utils.start_new_os_process().run_process()')

		args_lst = [
			"-port=%s"%(p_port_str),
			"-mongodb_db_name=%s"%(p_mongo_db_name_str),
			"-images_store_local_dir_path=%s"%(p_images_store_local_dir_path_str),
			"-images_thumbnails_store_local_dir_path=%s"%(p_images_thumbnails_store_local_dir_path_str),
			"-images_s3_bucket_name=gf--test",
		]
		cmd_str = '%s %s'%(p_service_bin_path_str,
						' '.join(args_lst))
		p_log_fun('INFO','cmd_str - %s'%(cmd_str))

		p = subprocess.Popen(cmd_str,
						shell   = True,
						stdout  = subprocess.PIPE,
						bufsize = 1)
		assert isinstance(p,subprocess.Popen)
		#--------------------------------------------------
		#IMPORTANT!! - workers cleanup on parent shutdown

		def handle_signal_terminate(p_signum,
									p_frame):
			p_log_fun('INFO','+++ ++ -- SIGNAL SIGTERM RECEIVED - gf_images_main_service')

			p.terminate()

		import signal
		signal.signal(signal.SIGTERM, handle_signal_terminate)
		#--------------------------------------------------
		p.wait() #block this process and let the child run
	#---------------------------------------------------

	multiprocessing.log_to_stderr(logging.DEBUG)
	parent_conn, child_conn = multiprocessing.Pipe()

	#:multiprocessing.Process
	p = multiprocessing.Process(target = run_process, 
							args = (child_conn,
								p_log_fun))
	p.start()
	time.sleep(2) #time for service to startup
	
	return p