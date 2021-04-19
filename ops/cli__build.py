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

from colored import fg, bg, attr
import delegator

sys.path.append("%s/../meta"%(modd_str))
import gf_meta
import gf_web_meta

sys.path.append("%s/utils"%(modd_str))
import gf_build_go
import gf_build_rust
import gf_build_changes
import gf_log

sys.path.append("%s/tests"%(modd_str))
import gf_tests

sys.path.append("%s/../py/gf_aws"%(modd_str))
import gf_aws_s3

sys.path.append("%s/web"%(modd_str))
import gf_web__build

sys.path.append("%s/containers"%(modd_str))
import gf_containers
import gf_local_cluster

sys.path.append("%s/gf_builder"%(modd_str))
import gf_builder_ops

#--------------------------------------------------
def main():
	
	print("")
	print("                              %sGLOFLOW BUILD TOOL%s"%(fg("green"), attr(0)))
	print("")
	
	build_meta_map        = gf_meta.get()["build_info_map"]
	apps_changes_deps_map = gf_meta.get()["apps_changes_deps_map"]
	args_map   = parse_args()
	run_str    = args_map["run"]

	app_name_str = args_map["app"]
	assert app_name_str in build_meta_map.keys()

	#--------------------------------------------------
	def go_build_app_in_cont(p_static_bool):


		app_meta_map = build_meta_map[app_name_str]
		if not "go_output_path_str" in app_meta_map.keys():
			print("not a main package")
			exit()

		gf_build_go.run_in_cont(app_name_str,
			app_meta_map["go_path_str"],
			app_meta_map["go_output_path_str"])

	#--------------------------------------------------
	def go_build_app(p_static_bool):
		
		app_meta_map = build_meta_map[app_name_str]
		if not "go_output_path_str" in app_meta_map.keys():
			print("not a main package")
			exit()
			
		fetch_deps_bool = args_map["fetch_deps"]

		gf_build_go.run(app_name_str,
			app_meta_map["go_path_str"],
			app_meta_map["go_output_path_str"],
			p_static_bool = p_static_bool,
			p_go_get_bool = fetch_deps_bool)

	#--------------------------------------------------
	def rust_build_apps_in_cont():
		print("RUST BUILD IN CONTAINER...")
		gf_build_rust.run_in_cont()

	#--------------------------------------------------
	def rust_build_apps():
		assert app_name_str == "gf_data_viz" or \
			app_name_str == "gf_images_jobs"

		# APP_META
		app_meta_map = build_meta_map[app_name_str]
		assert "type_str" in app_meta_map.keys()
		assert app_meta_map["type_str"] == "lib_rust"

		# CARGO_CRATE_DIR_PATHS
		assert "cargo_crate_specs_lst" in app_meta_map.keys()
		cargo_crate_specs_lst = app_meta_map["cargo_crate_specs_lst"]
		assert isinstance(cargo_crate_specs_lst, list)
		for s_map in cargo_crate_specs_lst:
			assert "dir_path_str" in s_map.keys()
			assert os.path.isdir(s_map["dir_path_str"])

		for s_map in cargo_crate_specs_lst:

			dir_path_str = s_map["dir_path_str"]

			print("")
			print("------------------------------------------------------------")
			print("       BUILD CARGO CRATE - %s"%(dir_path_str))
			print("")

			# STATIC
			static_bool = False
			if "static_bool" in s_map.keys():
				static_bool = s_map["static_bool"]

			# BUILD
			gf_build_rust.run(dir_path_str,
				p_static_bool = static_bool)

			# PREPARE_LIBS
			gf_build_rust.prepare_libs(app_name_str,
				dir_path_str,
				app_meta_map["type_str"])

	#--------------------------------------------------
	# AWS_CREDS
	def aws_creds_get():
		
		aws_creds_file_path_str     = args_map["aws_creds"]
		aws_creds_file_path_abs_str = os.path.abspath(aws_creds_file_path_str)
		print(aws_creds_file_path_abs_str)
		assert os.path.isfile(aws_creds_file_path_abs_str)

		aws_creds_map = gf_aws_s3.parse_creds(aws_creds_file_path_str)
		
		return aws_creds_map

	#-------------
	# BUILD_GO
	if run_str == "build" or run_str == "build_go":
		
		build_outof_cont_bool = args_map["build_outof_cont"]
		if build_outof_cont_bool:
			go_build_app(False)
		else:
			go_build_app_in_cont(False)

	#-------------
	# BUILD_RUST
	elif run_str == "build_rust":
		
		build_outof_cont_bool = args_map["build_outof_cont"]

		if build_outof_cont_bool:
			rust_build_apps()
		else:
			rust_build_apps_in_cont()

	#-------------
	# BUILD_WEB
	elif run_str == "build_web":
		apps_names_lst = [app_name_str]
		web_meta_map   = gf_web_meta.get() 

		gf_web__build.build(apps_names_lst, web_meta_map, gf_log.log_fun)

	#-------------
	# BUILD_CONTAINERS
	elif run_str == "build_containers":



		assert app_name_str in build_meta_map.keys()
		app_build_meta_map = build_meta_map[app_name_str]

		# if app_build_meta_map["type_str"] == "main_go":
		#
		# 	# STATIC_LINKING
		# 	# build using static linking, containers are based on Alpine linux, 
		# 	# which has a minimal stdlib and other libraries, so we want to compile 
		# 	# everything needed by this Go package into a single binary.
		# 	go_build_app_in_cont(True)
		
		dockerhub_user_str = args_map["dockerhub_user"]
		docker_sudo_bool   = args_map["docker_sudo"]



		# GF_BUILDER
		if app_name_str.startswith("gf_builder"):

			gf_builder_ops.cont__build(app_name_str,
				dockerhub_user_str,
				gf_log.log_fun,
				p_docker_sudo_bool = docker_sudo_bool)


		else:
			
			app_web_meta_map = None
			web_meta_map     = gf_web_meta.get()
			if app_name_str in web_meta_map.keys():
				app_web_meta_map = web_meta_map[app_name_str]

			gf_containers.build(app_name_str, 
				app_build_meta_map,
				gf_log.log_fun,
				p_app_web_meta_map = app_web_meta_map,
				p_user_name_str    = dockerhub_user_str,
				p_docker_sudo_bool = docker_sudo_bool)

	#-------------
	# PUBLISH_CONTAINERS
	elif run_str == "publish_containers":

		dockerhub_user_str = args_map["dockerhub_user"]
		dockerhub_pass_str = args_map["dockerhub_pass"]
		docker_sudo_bool   = args_map["docker_sudo"]

		assert app_name_str in build_meta_map.keys()
		app_build_meta_map = build_meta_map[app_name_str]

		# GIT_COMMIT_HASH
		git_commit_hash_str = None
		if "DRONE_COMMIT" in os.environ.keys():
			git_commit_hash_str = os.environ["DRONE_COMMIT"]

		# # GF_BUILDER
		# if app_name_str.startswith("gf_builder"):
		#
		# 	gf_builder_ops.cont__publish(app_name_str,
		# 		app_build_meta_map,
		# 		dockerhub_user_str,
		# 		dockerhub_pass_str,
		# 		gf_log.log_fun,
		# 		p_docker_sudo_bool = docker_sudo_bool)
		# else:
			
		gf_builder_ops.cont__publish(app_name_str,
			app_build_meta_map,
			dockerhub_user_str,
			dockerhub_pass_str,
			gf_log.log_fun,
			p_git_commit_hash_str = git_commit_hash_str,
			p_docker_sudo_bool    = docker_sudo_bool)
		
	#-------------
	# TEST
	elif run_str == "test":

		app_meta_map  = build_meta_map[app_name_str]
		test_name_str = args_map["test_name"]
		aws_creds_map = aws_creds_get()

		gf_tests.run(app_name_str,
			test_name_str,
			app_meta_map,
			aws_creds_map)

	#-------------
	# LIST_CHANGED_APPS
	elif run_str == "list_changed_apps":
		changed_apps_map = gf_build_changes.list_changed_apps(apps_changes_deps_map)
		gf_build_changes.view_changed_apps(changed_apps_map, "go")
		gf_build_changes.view_changed_apps(changed_apps_map, "web")
	
	#-------------
	# START_CLUSTER_LOCAL

	elif run_str == "start_cluster_local":
		
		docker_sudo_bool = args_map["docker_sudo"]
		aws_creds_map    = aws_creds_get()


		gf_local_cluster.start(aws_creds_map,
			p_docker_sudo_bool = docker_sudo_bool)

	#-------------
	# # GF_BUILDER__CONTAINER_BUILD
	# elif run_str == "gf_builder__cont_build":
	# 	dockerhub_user_str = args_map["dockerhub_user"]
	# 	docker_sudo_bool   = args_map["docker_sudo"]
	#
	# 	gf_builder_ops.cont__build(dockerhub_user_str,
	# 		gf_log.log_fun,
	# 		p_docker_sudo_bool = docker_sudo_bool)
	
	#-------------
	else:
		print("unknown run command - %s"%(run_str))
		exit()

#--------------------------------------------------
def parse_args():

	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#-------------
	# RUN
	arg_parser.add_argument("-run", action = "store", default = "build",
		help = '''
- '''+fg('yellow')+'build | build_go'+attr(0)+'''    - build app golang code
- '''+fg('yellow')+'build_rust'+attr(0)+'''          - build app golang code
- '''+fg('yellow')+'build_web'+attr(0)+'''           - build app web code (ts/js/css/html)
- '''+fg('yellow')+'build_containers'+attr(0)+'''    - build app Docker containers
- '''+fg('yellow')+'publish_containers'+attr(0)+'''  - publish app Docker containers
- '''+fg('yellow')+'test'+attr(0)+'''                - run app code tests
- '''+fg('yellow')+'list_changed_apps'+attr(0)+'''   - list all apps (and files) that have changed from last to the last-1 commit (for monorepo CI)
- '''+fg('yellow')+'start_cluster_local'+attr(0)+''' - start a local GF cluster using docker-compose

		''')
		
	#-------------
	# APP
	arg_parser.add_argument('-app', action = "store", default = 'gf_images',
		help = '''
- '''+fg('yellow')+'gf_solo'+attr(0)+'''
- '''+fg('yellow')+'gf_images'+attr(0)+'''
- '''+fg('yellow')+'gf_images_lib'+attr(0)+'''
- '''+fg('yellow')+'gf_publisher'+attr(0)+'''
- '''+fg('yellow')+'gf_tagger'+attr(0)+'''
- '''+fg('yellow')+'gf_landing_page'+attr(0)+'''
- '''+fg('yellow')+'gf_analytics'+attr(0)+'''
- '''+fg('yellow')+'gf_crawl_lib'+attr(0)+'''
- '''+fg('yellow')+'gf_crawl_core'+attr(0)+'''

- '''+fg('yellow')+'gf_images_jobs'+attr(0)+'''
- '''+fg('yellow')+'gf_ml_worker'+attr(0)+'''

- '''+fg('yellow')+'gf_builder'+attr(0)+'''
- '''+fg('yellow')+'gf_builder_go_ubuntu'+attr(0)+'''
- '''+fg('yellow')+'gf_builder_rust_ubuntu'+attr(0)+'''

		''')

	#-------------
	# TEST_NAME
	arg_parser.add_argument('-test_name',
		action =  "store",
		default = "all",
		help =    '''if only a particular test needs to be run''')

	#-------------
	# AWS_S3_CREDS
	arg_parser.add_argument('-aws_creds',
		action =  "store",
		default = "%s/../../creds/aws/s3.txt"%(modd_str),
		help =    '''path to the file containing AWS S3 credentials to be used''')

	#-------------
	# DOCKERHUB_USER/PASS
	arg_parser.add_argument('-dockerhub_user',
		action =  "store",
		default = "glofloworg",
		help =    '''name of the dockerhub user to target''')

	#-------------
	# BUILD_OUTOF_CONT
	arg_parser.add_argument("-build_outof_cont", action = "store_true", default=False,
		help = "build outside of a gf_builder container")

	#-------------
	# FETCH_DEPENDENCIES
	arg_parser.add_argument("-fetch_deps", action = "store_true", default=False,
		help = "explicit fetch of dependencies for Go/Py/Rust/JS if its configurable")

	#----------------------------
	# RUN_WITH_SUDO - boolean flag
	# in the default Docker setup the daemon is run as root and so docker client commands have to be run with "sudo".
	# newer versions of Docker allow for non-root users to run Docker daemons. 
	# also CI systems might run this command in containers as root-level users in which case "sudo" must not be specified.
	arg_parser.add_argument('-docker_sudo', action = "store_true", default=False,
		help = "specify if certain Docker CLI commands are to run with 'sudo'")

	#-------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)


	if not args_namespace.dockerhub_user == None:
		gf_dockerhub_user_str = args_namespace.dockerhub_user
	else:
		gf_dockerhub_user_str = os.environ.get("GF_DOCKERHUB_USER", None)

	gf_dockerhub_pass_str = os.environ.get("GF_DOCKERHUB_P", None)


	args_map = {
		"run":              args_namespace.run,
		"app":              args_namespace.app,
		"test_name":        args_namespace.test_name,
		"aws_creds":        args_namespace.aws_creds,
		"dockerhub_user":   gf_dockerhub_user_str, # ,
		"dockerhub_pass":   gf_dockerhub_pass_str,
		"docker_sudo":      args_namespace.docker_sudo,
		"build_outof_cont": args_namespace.build_outof_cont,
		"fetch_deps":       args_namespace.fetch_deps
	}
	return args_map

#--------------------------------------------------
main()