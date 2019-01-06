import time
import clint
#----------------------------------------------
def get_log_fun(p_log_fun):
	p_log_fun('FUN_ENTER','gf_core_logger.get_log_fun()')
	
	#----------------------------------------------	
	def log_color_display_fun(p_group,
							p_msg):
		
		t = str(time.time())			 
		if p_group == 'FUN_ENTER':
			clint.textui.puts(t+':'+clint.textui.colored.yellow('FUN_ENTER')+':'+p_msg)
			
		elif p_group == 'ERROR':
			clint.textui.puts(t+':'+clint.textui.colored.red(p_group)+':'+p_msg)
			
		elif p_group == 'INFO':
			clint.textui.puts(t+':'+clint.textui.colored.green(p_group)+':'+ \
												clint.textui.colored.green(p_msg))
			
		elif p_group == 'INFO_USR':
			clint.textui.puts(t+':'+clint.textui.colored.magenta(p_group)+':'+ \
												clint.textui.colored.magenta(p_msg))
			
		#log message is in some way related to external-systems data
		elif p_group == 'EXTERN':
			clint.textui.puts(t+':'+clint.textui.colored.blue(p_group)+':'+p_msg)
			
		#if 'TEST' is anywhere in the group string
		#elif p_group.find('TEST'):
		#	clint.textui.puts(clint.textui.colored.cyan(p_group)+':'+p_msg)
			
		#--------------
		#JAVASCRIPT FORMATING
		elif p_group == 'JS:FUN_ENTER':
			clint.textui.puts(t+':'+'       '+clint.textui.colored.yellow('FUN_ENTER')+':'+p_msg)
			
		elif p_group == 'JS:INFO':
			clint.textui.puts(t+':'+'       '+clint.textui.colored.green(p_group)+':'+ \
													clint.textui.colored.green(p_msg))
	#----------------------------------------------
			
	return log_color_display_fun