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

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import argparse
import urllib.parse
from colored import fg, bg, attr
import delegator
import requests

sys.path.append('%s/../../meta'%(modd_str))
import gf_meta
import gf_web_meta

sys.path.append('%s/../../ops/utils'%(modd_str))
import gf_build_changes
import gf_build_go
import gf_build_rust
import gf_log

sys.path.append('%s/../../ops/tests'%(modd_str))
import gf_tests

sys.path.append('%s/../../ops/web'%(modd_str))
import gf_web__build

sys.path.append('%s/../../py/gf_aws'%(modd_str))
import gf_os_aws_creds

sys.path.append('%s/../../ops/containers'%(modd_str))
import gf_os_docker
import gf_containers

#--------------------------------------------------
def main():
	
	print("    ---   GF_BUILDER -------------------")
	print(delegator.run("ls -al").out)
	print("pwd[%s] - whoami[%s]"%(delegator.run("pwd").out.strip(), delegator.run("whoami").out.strip()))

	def log_fun(m, g): print("%s: %s"%(m, g))
	args_map = parse_args()

	# GET_CHANGED_APPS
	changed_apps_files_map = get_changed_apps()

	#------------------------
	# TEST
	if args_map["run"] == "test":

		# TEST_SERVICES_RUN
		if args_map["test_services"]:
			test_services_run(log_fun,
				p_docker_sudo_bool = args_map["docker_sudo"])

		test_apps(changed_apps_files_map)

		# TEST_SERVICES_STOP
		if args_map["test_services"]:
			test_services_stop(log_fun,
				p_docker_sudo_bool = args_map["docker_sudo"])

	#------------------------
	# BUILD
	elif args_map["run"] == "build":

		# IMPORTANT!! - only insert Git commit hash if gf_builder.py is run in CI
		if not args_map["drone_commit"] == None:
			git_commit_hash_str = args_map["drone_commit"]
			paste_git_commit_hash(git_commit_hash_str)
			
		build_apps(changed_apps_files_map,
			p_static_bool=True)

	#------------------------
	# # BUILD_GO
	# elif args_map["run"] == "build_go":
	#
	# 	# IMPORTANT!! - only insert Git commit hash if gf_builder.py is run in CI
	# 	if not args_map["drone_commit"] == None:
	# 		git_commit_hash_str = args_map["drone_commit"]
	# 		paste_git_commit_hash(git_commit_hash_str)
	#
	# 	build_apps(changed_apps_files_map,
	# 		p_build_web_bool = False,
	# 		p_build_go_bool  = True,
	# 		p_static_bool    = True)

	#------------------------
	# # BUILD_RUST
	# elif args_map["run"] == "build_rust":
	#
	# 	build_rust()

	#------------------------
	# # BUILD_WEB
	# elif args_map["run"] == "build_web":
	#
	# 	# IMPORTANT!! - only insert Git commit hash if gf_builder.py is run in CI
	# 	if not args_map["drone_commit"] == None:
	# 		git_commit_hash_str = args_map["drone_commit"]
	# 		paste_git_commit_hash(git_commit_hash_str)
	#
	# 	build_apps(changed_apps_files_map,
	# 		p_build_web_bool = True,
	# 		p_build_go_bool  = False)

	#------------------------
	# BUILD_CONTAINERS
	elif args_map["run"] == "build_containers":

		gf_dockerhub_user_str = args_map["gf_dockerhub_user"]
		
		# GIT_COMMIT_HASH
		git_commit_hash_str = None
		if "DRONE_COMMIT" in os.environ.keys():
			git_commit_hash_str = os.environ["DRONE_COMMIT"]

		build_apps_containers(changed_apps_files_map,
			gf_dockerhub_user_str,
			p_git_commit_hash_str = git_commit_hash_str,
			p_docker_sudo_bool    = args_map["docker_sudo"])

	#------------------------
	# PUBLISH_CONTAINERS
	elif args_map["run"] == "publish_containers":

		gf_dockerhub_user_str = args_map["gf_dockerhub_user"]
		gf_dockerhub_pass_str = args_map["gf_dockerhub_pass"]

		assert not gf_dockerhub_user_str == None
		assert not gf_dockerhub_pass_str == None

		# GIT_COMMIT_HASH
		git_commit_hash_str = None
		if "DRONE_COMMIT" in os.environ.keys():
			git_commit_hash_str = os.environ["DRONE_COMMIT"]

		publish_apps_containers(changed_apps_files_map,
			gf_dockerhub_user_str,
			gf_dockerhub_pass_str,
			p_git_commit_hash_str = git_commit_hash_str,
			p_docker_sudo_bool    = args_map["docker_sudo"])

	#------------------------
	# NOTIFY_COMPLETION
	elif args_map["run"] == "notify_completion":

		gf_notify_completion_url_str = args_map["gf_notify_completion_url"]
		assert not gf_notify_completion_url_str == None

		# GIT_COMMIT_HASH
		git_commit_hash_str = None
		if "DRONE_COMMIT" in os.environ.keys():
			git_commit_hash_str = os.environ["DRONE_COMMIT"]

			p_git_commit_hash_str = git_commit_hash_str)

	#------------------------

#--------------------------------------------------
# NOTIFY_COMPLETION
def notify_completion(p_gf_notify_completion_url_str,
	p_git_commit_hash_str = None,
	p_app_name_str        = None):
	
	url_str = None

	# add git_commit_hash as a querystring argument to the notify_completion URL.
	# the entity thats receiving the completion notification needs to know what the tag
	# is of the newly created container.
	if not p_git_commit_hash_str == None:
		url = urllib.parse.urlparse(p_gf_notify_completion_url_str)
		
		# QUERY_STRING
		qs_lst = urllib.parse.parse_qsl(url.query)
		qs_lst.append(("base_img_tag", p_git_commit_hash_str)) # .parse_qs() places all values in lists

		qs_str = "&".join(["%s=%s"%(k, v) for k, v in qs_lst])

		# _replace() - "url" is of type ParseResult which is a subclass of namedtuple;
		#              _replace is a namedtuple method that:
		#              "returns a new instance of the named tuple replacing
		#              specified fields with new values".
		url_new = url._replace(query=qs_str)
		url_str = url_new.geturl()
	else:
		url_str = p_gf_notify_completion_url_str

	print("NOTIFY_COMPLETION - HTTP REQUEST - %s"%(url_str))

	# HTTP_POST

	data_map = {}
	if not p_app_name_str == None:
		data_map["app_name"] = p_app_name_str

	r = requests.post(url_str, json=data_map)
	print(r.text)

	if not r.status_code == 200:
		print("notify_completion http request failed")
		exit(1)
		
#--------------------------------------------------
# PUBLISH_APPS_CONTAINERS
def publish_apps_containers(p_changed_apps_files_map,
	p_gf_dockerhub_user_str,
	p_gf_dockerhub_pass_str,
	p_git_commit_hash_str = None,
	p_docker_sudo_bool    = True):
	assert isinstance(p_gf_dockerhub_user_str, str)
	assert isinstance(p_gf_dockerhub_pass_str, str)
	
	build_meta_map = gf_meta.get()["build_info_map"]
	
	# "all" - this key holds a map with all the apps that had either their Go or Web code changed
	apps_names_lst = p_changed_apps_files_map["all"].keys()

	for app_name_str in apps_names_lst:

		assert build_meta_map.has_key(app_name_str)
		app_build_meta_map = build_meta_map[app_name_str]

		# PUBLISH
		gf_containers.publish(app_name_str,
			app_build_meta_map,
			p_gf_dockerhub_user_str,
			p_gf_dockerhub_pass_str,
			gf_log.log_fun,
			p_git_commit_hash_str = p_git_commit_hash_str, 
			p_exit_on_fail_bool   = True,
			p_docker_sudo_bool    = p_docker_sudo_bool)

#--------------------------------------------------
# BUILD_APPS_CONTAINERS
def build_apps_containers(p_changed_apps_files_map,
	p_gf_dockerhub_user_str,
	p_git_commit_hash_str = None,
	p_docker_sudo_bool    = True):
	assert isinstance(p_changed_apps_files_map, dict)

	build_meta_map = gf_meta.get()["build_info_map"]
	web_meta_map   = gf_web_meta.get()

	# IMPORTANT!! - for each app that has any of its code changed rebuild both the Go and Web code,
	#               since the containers has to be fully rebuilt.
	# "all" - this key holds a map with all the apps that had either their Go or Web code changed
	apps_names_lst = p_changed_apps_files_map["all"].keys()


	for app_name_str in apps_names_lst:

		assert build_meta_map.has_key(app_name_str)
		app_build_meta_map = build_meta_map[app_name_str]

		assert web_meta_map.has_key(app_name_str)
		app_web_meta_map = web_meta_map[app_name_str]
		
		gf_containers.build(app_name_str,
			app_build_meta_map,
			app_web_meta_map,
			gf_log.log_fun,

			p_git_commit_hash_str = p_git_commit_hash_str,

			# DOCKERHUB_USER
			p_user_name_str = p_gf_dockerhub_user_str,

			# gf_containers.build() should exit if the docker build CLI run returns with a non-zero exit code.
			# gf_builder.py is meant to run in CI environments, and so we want the stage in which it runs 
			# to be marked as failed because of the non-zero exit code.
			p_exit_on_fail_bool = True,
			p_docker_sudo_bool  = p_docker_sudo_bool)

#--------------------------------------------------
def test_services_run(p_log_fun,
	p_docker_sudo_bool = True):
	
	p_log_fun("INFO", "TEST SERVICES RUN -----------------------------------------------------")
	p_log_fun("INFO", "mongodb")
	p_log_fun("INFO", "elasticsearch")


	mongo_cont_name_str  = "test_mongo"
	mongo_image_name_str = "mongo"

	search_cont_name_str  = "test_elasticsearch"
	search_image_name_str = "elasticsearch:5-alpine"

	# IMPORTANT!! - HOST_NETWORKING - test_services containers are started with host networking.
	#                                 the testing stage container itself should also be run with host networking enabled,
	#                                 so that test_services and test stage share the same network address space and can reach each other.
	# FIX!! - find a better way to do this networking without using host networking.
	#------------------------
	# MONGODB
	# remove container if its running
	if gf_os_docker.cont_is_running(mongo_cont_name_str, p_log_fun, p_docker_sudo_bool = p_docker_sudo_bool):
		gf_os_docker.remove_by_name(mongo_cont_name_str, p_log_fun, p_docker_sudo_bool = p_docker_sudo_bool)

	gf_os_docker.run(mongo_image_name_str, p_log_fun,
		p_container_name_str = mongo_cont_name_str,
		p_hostname_str       = "mongo",
		p_host_network_bool  = True,
		p_ports_map          = {"27017": "27017"},
		p_docker_sudo_bool   = p_docker_sudo_bool)

	#------------------------
	# ELASTICSEARCH
	# remove container if its running
	if gf_os_docker.cont_is_running(search_cont_name_str, p_log_fun, p_docker_sudo_bool = p_docker_sudo_bool):
		gf_os_docker.remove_by_name(search_cont_name_str, p_log_fun, p_docker_sudo_bool = p_docker_sudo_bool)

	gf_os_docker.run(search_image_name_str, p_log_fun,
		p_container_name_str = search_cont_name_str,
		p_hostname_str       = "elasticsearch",
		p_host_network_bool  = True,
		p_ports_map          = {"9200": "9200"},
		p_docker_sudo_bool   = p_docker_sudo_bool)

	#------------------------

#--------------------------------------------------
def test_services_stop(p_log_fun,
	p_docker_sudo_bool = True):
	
	p_log_fun("INFO", "TEST SERVICES STOP -----------------------------------------------------")
	p_log_fun("INFO", "mongodb")
	p_log_fun("INFO", "elasticsearch")

	mongo_cont_name_str  = "test_mongo"
	search_cont_name_str = "test_elasticsearch"
	gf_os_docker.remove_by_name(mongo_cont_name_str, p_log_fun, p_docker_sudo_bool = p_docker_sudo_bool)
	gf_os_docker.remove_by_name(search_cont_name_str, p_log_fun, p_docker_sudo_bool = p_docker_sudo_bool)

#--------------------------------------------------
def test_apps(p_changed_apps_files_map):
	assert isinstance(p_changed_apps_files_map, dict)

	print("\n\n TEST APPS ----------------------------------------------------- \n\n")

	build_meta_map        = gf_meta.get()['build_info_map']
	apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']
	apps_gf_packages_map  = apps_changes_deps_map["apps_gf_packages_map"]

	# AWS_CREDS
	aws_creds_map = gf_os_aws_creds.get_from_env_vars()
	assert isinstance(aws_creds_map, dict)

	# nothing changed
	if len(p_changed_apps_files_map.keys()) == 0:
		return
	else:

		#------------------------
		# GO
		print("\nGO--------\n")
		for app_name_str, v in p_changed_apps_files_map["go"].items():
			
			test_name_str = "all"

			# IMPORTANT!! - get all packages that are dependencies of this app, so that 
			#               tests for all these packages can be run. dont just run the tests of the app
			#               but of all the packages that are its dependencies as well.
			app_gf_packages_lst = apps_gf_packages_map[app_name_str]

			# RUN_TESTS_FOR_ALL_APP_PACKAGES
			for app_gf_package_name_str in app_gf_packages_lst:

				print("about to run tests for - %s"%(app_gf_package_name_str))
				if build_meta_map.has_key(app_gf_package_name_str):
					gf_package_meta_map = build_meta_map[app_gf_package_name_str]

					gf_tests.run(app_gf_package_name_str,
						test_name_str,
						gf_package_meta_map,
						aws_creds_map,

						# IMPORTANT!! - in case the tests that gf_test.run() executes fail, 
						#               run() should call exit() and force this whole process to exit, 
						#               so that CI marks the build as failed.
						p_exit_on_fail_bool = True)
						
		#------------------------

#--------------------------------------------------
# BUILD_RUST
def build_rust():

	print("building RUST...")
	
	apps_names_lst = [
		# "gf_data_viz",
		"gf_images_jobs"
	]
	build_meta_map = gf_meta.get()["build_info_map"]
	
	for app_name_str in apps_names_lst:
		print("\n\nRUST-------- %s\n\n"%(app_name_str))
			
		app_meta_map = build_meta_map[app_name_str]
		assert app_meta_map["type_str"] == "lib_rust"



		for cargo_crate_map in app_meta_map["cargo_crate_specs_lst"]:


			cargo_crate_dir_path_str = cargo_crate_map["dir_path_str"]
			assert os.path.dirname(cargo_crate_dir_path_str)
			
			# IMPORTANT!! - it can be specified on a per-crate basis if it should
			#               be compiled statically or not.
			#               for gf_images_jobs_py for example we specifically want to
			#               build a dynamic lib (even in Alpine) for importing into the Py VM.
			static_bool = cargo_crate_map.get("static_bool", False)

			# RUN
			gf_build_rust.run(cargo_crate_dir_path_str, 
				p_static_bool=static_bool)

			# PREPARE_LIBS
			gf_build_rust.prepare_libs(app_name_str,
				cargo_crate_dir_path_str,
				app_meta_map["type_str"])

#--------------------------------------------------
# BUILD_APPS
def build_apps(p_changed_apps_files_map,
	p_build_web_bool = True,
	p_build_go_bool  = True,
	p_static_bool    = True):
	assert isinstance(p_changed_apps_files_map, dict)

	print("\n\n BUILD APPS ----------------------------------------------------- \n\n")

	build_meta_map = gf_meta.get()["build_info_map"]
	
	# nothing changed
	if len(p_changed_apps_files_map.keys()) == 0:
		return
	else:

		# IMPORTANT!! - for each app that has any of its code changed rebuild both the Go and Web code,
		#               since the containers has to be fully rebuilt.
		# "all" - this key holds a map with all the apps that had either their Go or Web code changed
		apps_names_lst = list(p_changed_apps_files_map["all"].keys())

		#------------------------
		# WEB
		if p_build_web_bool:
			print("\n\nWEB--------\n\n")

			web_meta_map = gf_web_meta.get()

			gf_web__build.build(apps_names_lst,
				web_meta_map,
				gf_log.log_fun)
		
		#------------------------
		# GO
		if p_build_go_bool:
			for app_name_str in apps_names_lst:

				app_meta_map = build_meta_map[app_name_str]

				print("\n\nGO--------\n\n")
				app_go_path_str        = app_meta_map["go_path_str"]
				app_go_output_path_str = app_meta_map["go_output_path_str"]

				gf_build_go.run(app_name_str,
					app_go_path_str,
					app_go_output_path_str,

					# IMPORTANT!! - binaries are packaged in Alpine Linux, which uses a different standard library then stdlib, 
					#               so all binary dependencies are to be statically linked into the output binary 
					#               without depending on standard dynamic linking.
					p_static_bool = p_static_bool, 
					
					# gf_build.run_go() should exit if the "go build" CLI run returns with a non-zero exit code.
					# gf_builder.py is meant to run in CI environments, and so we want the stage in which it runs 
					# to be marked as failed because of the non-zero exit code.
					p_exit_on_fail_bool = True)

		#------------------------

#--------------------------------------------------
# GET_CHANGED_APPS
def get_changed_apps():
	print("DIFF")
	apps_changes_deps_map = gf_meta.get()["apps_changes_deps_map"]

	# LIST_CHANGED_APPS - determine how which apps/services changed
	changed_apps_files_map = gf_build_changes.list_changed_apps(apps_changes_deps_map,
		p_commits_lookback_int = 1, 
		p_mark_all_bool        = True)

	# VIEW
	gf_build_changes.view_changed_apps(changed_apps_files_map, "go")
	gf_build_changes.view_changed_apps(changed_apps_files_map, "web")
	
	return changed_apps_files_map

#--------------------------------------------------
# PASTE_GIT_COMMIT_HASH
def paste_git_commit_hash(p_git_commit_hash_str):
	print("PASTE_GIT_COMMIT_HASH - %s"%(p_git_commit_hash_str))

	golang_sys_release_info_file_path_str = "%s/../../go/gf_core/gf_sys_release_info.go"%(modd_str)
	assert os.path.isfile(golang_sys_release_info_file_path_str)
	
	original_word_regex_str = 'Git_commit_str: "",' #this is the original line of Go code
	new_word_regex_str      = 'Git_commit_str: "%s",'%(p_git_commit_hash_str)

	#------------------------
	# IMPORTANT!! - "sed" - Stream EDitor.
	#               "-i" - in-place, save to original file
	#               command string:
	#                 "s" - the substitute command
	#                 "g" - global, replace all not just first instance
	c = "sed -i 's/%s/%s/g' %s"%(original_word_regex_str, new_word_regex_str, golang_sys_release_info_file_path_str)
	print(c)
	
	#------------------------

	r = delegator.run(c)
	print(r.out)
	print(r.err)

#--------------------------------------------------
# IMPORTANT!! - get Git commit from the deployed artifact (making API call to a target service).
#               this is needed to know how far HEAD of this monorepo is ahead from the commit 
#               that was used to build a particular service, to then use that integer distance
#               as the p_commits_lookback_int when determening which apps change when calling 
#               list_changed_apps(). 

def get_deployed_commit(p_domain_str = "https://gloflow.com"):
	True

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#-------------
	# RUN
	arg_parser.add_argument("-run", action = "store", default = "build",
		help = '''
- '''+fg('yellow')+'test'+attr(0)+'''               - run app code tests
- '''+fg('yellow')+'build'+attr(0)+'''              - build app golang/web code
- '''+fg('yellow')+'build_rust'+attr(0)+'''         - build app rust code
- '''+fg('yellow')+'build_containers'+attr(0)+'''   - build app Docker containers
- '''+fg('yellow')+'publish_containers'+attr(0)+''' - publish app Docker containers
- '''+fg('yellow')+'notify_completion'+attr(0)+'''  - notify remote HTTP endpoint of build completion
		''')

	#-------------
	# # MONGODB_HOST
	# arg_parser.add_argument('-mongodb_host',
	#     action =  "store",
	#     default = "all",
	#     help =    '''mongodb host to connect to (for testing, etc.)''')

	#----------------------------
	# RUN_WITH_SUDO - boolean flag
	# in the default Docker setup the daemon is run as root and so docker client commands have to be run with "sudo".
	# newer versions of Docker allow for non-root users to run Docker daemons. 
	# also CI systems might run this command in containers as root-level users in which case "sudo" must not be specified.
	arg_parser.add_argument("-test_services", action = "store_true",
		help = "if testing services (mongo, elasticsearch, etc.) should be started in containers when running tests")

	#----------------------------
	# RUN_WITH_SUDO - boolean flag
	# in the default Docker setup the daemon is run as root and so docker client commands have to be run with "sudo".
	# newer versions of Docker allow for non-root users to run Docker daemons. 
	# also CI systems might run this command in containers as root-level users in which case "sudo" must not be specified.
	arg_parser.add_argument("-docker_sudo", action = "store_true",
		help = "specify if certain Docker CLI commands are to run with 'sudo'")

	#-------------
	# ENV_VARS
	drone_commit_str             = os.environ.get("DRONE_COMMIT", None) # Drone defined ENV var
	gf_dockerhub_user_str        = os.environ.get("GF_DOCKERHUB_USER", None)
	gf_dockerhub_pass_str        = os.environ.get("GF_DOCKERHUB_P", None)
	gf_notify_completion_url_str = os.environ.get("GF_NOTIFY_COMPLETION_URL", None)

	print(gf_dockerhub_user_str)
	print(gf_dockerhub_pass_str)
	
	#-------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)
	return {
		"run":                      args_namespace.run,
		"drone_commit":             drone_commit_str,
		"gf_dockerhub_user":        gf_dockerhub_user_str,
		"gf_dockerhub_pass":        gf_dockerhub_pass_str,
		"gf_notify_completion_url": gf_notify_completion_url_str,
		"test_services":            args_namespace.test_services,
		"docker_sudo":              args_namespace.docker_sudo
		# "mongodb_host": args_namespace.mongodb_host,
	}

#--------------------------------------------------

if __name__ == "__main__":
	main()