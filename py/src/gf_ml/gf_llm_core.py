# GloFlow application and media management/publishing platform
# Copyright (C) 2024 Ivan Trajkovic
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

import json
import requests

import sentry_sdk

#--------------------------------------------
def run_llm(p_prompt_str,
	p_openai_api_key_str,
	p_model_str="gpt-4",
	p_temperature=0.7,
	p_functions_meta_map=None,
	p_verbose_bool=False):

	ctx_lst = []
	
	init_resp_str, first_func_call_map, sub_ctx_lst = model_request(p_prompt_str,
		p_openai_api_key_str,
		p_model_str=p_model_str,
		p_functions_meta_map = p_functions_meta_map)

	ctx_lst.extend(sub_ctx_lst)

	#--------------------------------------------
	def develop_context_with_model(p_first_func_call_map,
		p_ctx_map):

		if p_verbose_bool:
			print(">>>>>>>>>>>>>>>> entering function calls loop...")

		func_call_to_run_map = p_first_func_call_map

		i=0
		while True:

			if p_verbose_bool:
				print(f">>>>>>>>>>>>>>>> ---------------------------- local function execution {i}...")
			
			# dont run more than 5 function calls
			if i > 5:
				break
			
			#----------------------
			# FUNC_EXEC
			gpt_result_msg_map = func_exec(func_call_to_run_map,
				p_functions_meta_map)

			# add the function call to the conversation context
			p_ctx_map.append(gpt_result_msg_map)

			# print(gpt_result_msg_map)

			#----------------------
			# RETURN_RESULTS_TO_MODEL

			new_gen_text_str, new_func_call_map, sub_ctx_lst = model_request_with_ctx(p_ctx_map,
				p_openai_api_key_str,
				p_model_str=p_model_str,
				p_temperature=p_temperature,
				p_functions_meta_map=p_functions_meta_map,
				p_verbose_bool=p_verbose_bool)
			
			if new_func_call_map == None:
				
				# no more function calls needed by the model,
				# just return the final generated text
				return new_gen_text_str
			
			else:

				# there are more function calls to run, so run them
				func_call_to_run_map = new_func_call_map
				p_ctx_map.extend(sub_ctx_lst)

			#----------------------

			i+=1
	
	#--------------------------------------------

	if p_verbose_bool:

		print("first llm response --------------")
		print(init_resp_str)

	
	# FUNCTION_CALL
	# first reply is a function call that might result in a whole sequence of calls that a model
	# has to execute in order to gather all the context it needs to generate the final answer.
	if not first_func_call_map == None:
		
		develop_context_with_model(first_func_call_map, ctx_lst)

	# REGULAR_REPLY       
	else:

		print(init_resp_str)
		
		if not init_resp_str == "":
			result_map = json.loads(init_resp_str)

			if p_verbose_bool:
				print("LLM result:")
				print(result_map)

			return result_map
		else:
			return {}

#--------------------------------------------
def func_exec(p_func_call_map,
	p_functions_meta_map,
	p_json_output_bool=False):

	func_name_str = p_func_call_map["name"]

	# FUNC_ARGS - argument values are defined as a JSON map
	func_arguments_map = json.loads(p_func_call_map["arguments"])

	
	#----------------------
	# CALL_FUNCTION
	# apply map of arguments to the function as named arguments

	func_result = p_functions_meta_map[func_name_str]["func"](**func_arguments_map)

	#----------------------

	output_str = ""

	if p_json_output_bool:
		output_str = json.dumps(func_result)
	else:
		if not func_result == None:
			assert isinstance(func_result, str)
			output_str = func_result

	gpt_result_msg_map = {

		"role": "function",
		"name": func_name_str,
		"content": output_str
	}

	

	return gpt_result_msg_map

#--------------------------------------------
def model_request_with_ctx(p_ctx_map,
	p_api_key_str,
	p_model_str="gpt-4",
	p_temperature=0.7,
	p_functions_meta_map=None,
	p_verbose_bool=False):

	if p_verbose_bool:
		print(">>>>>>>>>>>>>>>> model request (with ctx)...")

	# GPT_MESSAGE_HISTORY
	messages_lst = p_ctx_map

	data_map = {
		"model": p_model_str,
		"messages": messages_lst,
		"temperature": p_temperature
	}

	# if function signatures metadata is supplied, pass it to the model
	# for it to be able to use it in the completion.
	if not p_functions_meta_map == None:

		# definitions of functions meant for the GPT model
		functions_gpt_defs_lst = [v["gpt_def_map"] for _, v in p_functions_meta_map.items()]

		data_map["functions"] = functions_gpt_defs_lst

	# HTTP_REQUEST
	r_map = make_request(data_map, p_api_key_str, p_verbose_bool=p_verbose_bool)

	#----------------------
	# MSG

	msg_map = r_map["choices"][0]["message"]

	p_ctx_map.append(msg_map)

	# FUNCTION_CALL
	if "function_call" in msg_map.keys():
		func_call_map = msg_map["function_call"]

		if p_verbose_bool:
			print(f"GPT function_call - {func_call_map['name']}, args: {func_call_map['arguments']}")
		
		return None, func_call_map, p_ctx_map

	# GENERATED TEXT
	else:
		
		generated_text_str = msg_map["content"]
		
		if p_verbose_bool:
			print(f"\n\nanswer: {generated_text_str}")

		return generated_text_str, None, p_ctx_map
	
	#----------------------

#--------------------------------------------
def model_request(p_prompt_str,
	p_api_key_str,
	p_model_str="gpt-4",
	p_functions_meta_map=None,
	p_verbose_bool=False):

	if p_verbose_bool:
		print(">>>>>>>>>>>>>>>> model request...")

	ctx_lst = [] 
	
	if p_verbose_bool:
		print(f"""
		
		prompt: 
			{p_prompt_str}
		
		""")
	
	gpt_msg_map = {
		"role": "user",
		"content": p_prompt_str
	}
	ctx_lst.append(gpt_msg_map)

	data_map = {
		"model": p_model_str, # "gpt-3.5-turbo",
		"messages": [
			gpt_msg_map
		],
		"temperature": 0.7
	}

	# if function signatures metadata is supplied, pass it to the model
	# for it to be able to use it in the completion.
	if not p_functions_meta_map == None:

		# definitions of functions meant for the GPT model
		functions_gpt_defs_lst = [v["gpt_def_map"] for _, v in p_functions_meta_map.items()]

		data_map["functions"] = functions_gpt_defs_lst

	'''
	payload = {
		"prompt": text_prompt,
		"max_tokens": 100,
	}
	'''
	
	r_map = make_request(data_map, p_api_key_str)



	#----------------------
	# MSG

	reply_gpt_msg_map = r_map["choices"][0]["message"]

	ctx_lst.append(reply_gpt_msg_map)

	# FUNCTION_CALL
	if "function_call" in reply_gpt_msg_map.keys():
		func_call_map = reply_gpt_msg_map["function_call"]

		if p_verbose_bool:
			print(f"GPT function_call - {func_call_map['name']}, args: {func_call_map['arguments']}")
		
		return None, func_call_map, ctx_lst

	# GENERATED TEXT
	else:
		
		generated_text_str = reply_gpt_msg_map["content"]
		
		if p_verbose_bool:
			print(f"\n\nanswer: {generated_text_str}")

		return generated_text_str, None, ctx_lst
	
	#----------------------

#--------------------------------------------
def make_request(p_data_map,
	p_api_key_str,
	p_verbose_bool=False):

	# Set the API endpoint for GPT-4 (if available)
	# api_endpoint = "https://api.openai.com/v1/engines/gpt-4/completions"
	api_endpoint = "https://api.openai.com/v1/chat/completions"

	headers = {
		"Content-Type": "application/json",
		"Authorization": f"Bearer {p_api_key_str}",
	}

	# HTTP_POST
	response = requests.post(api_endpoint, json=p_data_map, headers=headers)

	if response.status_code == 200:

		r_map = response.json()

		if p_verbose_bool:
			print("GPT response:")
			print(r_map)

		#----------------------
		# METADATA
		query_run_unix_time_int = r_map["created"]
		model_used_str = r_map["model"]

		#----------------------
		# METRICS
		total_tokens_consumed_int = r_map["usage"]["total_tokens"]
		prompt_tokens_int         = r_map["usage"]["prompt_tokens"]
		completion_tokens_int     = r_map["usage"]["completion_tokens"]


		return r_map
	
	else:
		print("Error:", response.status_code, response.text)

		with sentry_sdk.push_scope() as scope:
			scope.set_extra("status_code", response.status_code)
			scope.set_extra("response_text", response.text)
	
			raise Exception("Error:", response.status_code, response.text)
		

#--------------------------------------------
def dalle3_request(p_prompt_str,
	p_api_key_str):

	url_str = "https://api.openai.com/v1/engines/dall-e-3/completions"

	data_map = {
		"prompt": p_prompt_str,
		"max_tokens": 256,
	}

	headers = {
		"Content-Type": "application/json",
		"Authorization": f"Bearer {p_api_key_str}",
	}

	resp = requests.post(url_str, json=data_map, headers=headers)

	if resp.status_code == 200:
		image_url_str = resp.json()["choices"][0]["text"]
		return image_url_str

		# Optionally, you can download the image using additional code
		# For example, using the requests module:
		# image_response = requests.get(image_url)
		# with open("generated_image.png", "wb") as image_file:
		#     image_file.write(image_response.content)

	else:
		print("Error:", resp.status_code, resp.text)