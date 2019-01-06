import traceback
#--------------------------------------------------------
#p_surrounding_context_attribs_tpl - order matters

def handle_exception(p_exception,
				p_formated_msg_str,
				p_surrounding_context_attribs_tpl,
				p_log_fun):
	p_log_fun('FUN_ENTER','gf_error.handle_exception()')
	assert isinstance(p_exception,Exception)

	if p_formated_msg_str                == None or \
	   p_surrounding_context_attribs_tpl == None:
	
		msg_str = ''
	else:
		assert isinstance(p_surrounding_context_attribs_tpl,tuple)
		msg_str = p_formated_msg_str%(p_surrounding_context_attribs_tpl)
		
	#print p_exception.message
	#print p_pymods_dict[p_pymods_dict.keys()[0]]
	#traceback_mod_ref = p_pymods_dict[p_pymods_dict.keys()[0]]
	
	#print traceback_mod_ref
	#print traceback_mod_ref.__name__
	#print dir(traceback_mod_ref)
	
	p_log_fun('INFO','p_exception.message:%s'%(p_exception.message))
	
	with_trace_msg_str = '''
		%s
		exception args:%s
		exception msg :%s
		trace         :%s'''%(msg_str,
			                  p_exception.args,
			                  p_exception.message,
						      traceback.format_exc())

	p_log_fun('ERROR',with_trace_msg_str)