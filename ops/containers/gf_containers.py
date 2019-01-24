import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import subprocess
import delegator
from colored import fg,bg,attr

sys.path.append('%s/../../meta'%(cwd_str))
import gf_meta
import gf_web_meta
#-------------------------------------------------------------
def build(p_app_name_str,
	p_log_fun,
	p_user_name_str = 'local'):
	p_log_fun('FUN_ENTER','gf_containers.build()')
	assert isinstance(p_app_name_str, basestring)
	
	#------------------
	#META
	build_meta_map = gf_meta.get()['build_info_map']

	if not build_meta_map.has_key(p_app_name_str):
		p_log_fun("ERROR","supplied app (%s) does not exist in gf_meta"%(p_app_name_str))
		return
	app_meta_map = build_meta_map[p_app_name_str]

	service_name_str     = app_meta_map['service_name_str']
	service_base_dir_str = app_meta_map['service_base_dir_str']
	assert os.path.isdir(service_base_dir_str)

	service_version_str  = app_meta_map['version_str']
	assert len(service_version_str.split(".")) == 4 #format x.x.x.x
	#------------------

	prepare_context_dir(p_app_name_str, service_base_dir_str, p_log_fun)


	build_docker_container(service_name_str,
		service_base_dir_str,
		service_version_str,
		p_user_name_str,
		p_log_fun)
#-------------------------------------------------------------
def prepare_context_dir(p_app_name_str,
	p_service_base_dir_str,
	p_log_fun):
	p_log_fun('FUN_ENTER','gf_containers.prepare_context_dir()')

	apps_meta_map = gf_web_meta.get()
	assert apps_meta_map.has_key(p_app_name_str)

	app_meta_map = apps_meta_map[p_app_name_str]

	assert app_meta_map.has_key('pages_map')
	pages_map = app_meta_map['pages_map']

	for pg_name_str, pg_info_map in pages_map.items():

		assert pg_info_map.has_key('build_dir_str')
		assert os.path.isdir(pg_info_map['build_dir_str'])
		build_dir_str = pg_info_map['build_dir_str']

		#------------------
		#TARGET_DIR
		target_dir_str = '%s/static'%(p_service_base_dir_str)
		r = delegator.run('mkdir -p %s'%(target_dir_str))

		if not r.out == '': print r.out
		if not r.err == '': print '%sFAILED%s >>>>>>>\n%s'%(fg('red'), attr(0), r.err)
		#------------------
		#COPY_PAGE_WEB_CODE
		c = 'cp -p -r %s/* %s'%(build_dir_str, target_dir_str)
		print(c)
		
		r = delegator.run(c)
		if not r.out == '': print r.out
		if not r.err == '': print '%sFAILED%s >>>>>>>\n%s'%(fg('red'), attr(0), r.err)
		#------------------
#-------------------------------------------------------------
def build_docker_container(p_service_name_str,
	p_service_base_dir_str,
	p_service_version_str,
	p_user_name_str,
	p_log_fun):
	p_log_fun('FUN_ENTER','gf_containers.build_docker_container()')
	assert os.path.isdir(p_service_base_dir_str)

	image_name_str       = '%s/%s:%s'%(p_user_name_str, p_service_name_str, p_service_version_str)
	dockerfile_path_str  = '%s/Dockerfile'%(p_service_base_dir_str)
	context_dir_path_str = p_service_base_dir_str

	p_log_fun('INFO','====================+++++++++++++++=====================')
	p_log_fun('INFO','                 BUILDING PACKAGE/SERVICE IMAGE')
	p_log_fun('INFO','              %s'%(p_service_name_str))
	p_log_fun('INFO','Dockerfile     - %s'%(dockerfile_path_str))
	p_log_fun('INFO','image_name_str - %s'%(image_name_str))
	p_log_fun('INFO','====================+++++++++++++++=====================')

	cmd_lst = [
		'sudo docker build',
		'-f %s'%(dockerfile_path_str),
		'--tag=%s'%(image_name_str),
		context_dir_path_str
	]

	cmd_str = ' '.join(cmd_lst)
	p_log_fun('INFO',' - %s'%(cmd_str))

	#change to the dir where the Dockerfile is located, for the 'docker'
	#tool to have the proper context
	old_cwd = os.getcwd()
	os.chdir(context_dir_path_str)
	
	r = subprocess.Popen(cmd_str, shell = True, stdout = subprocess.PIPE, bufsize = 1)

	#---------------------------------------------------
	def get_image_id_from_line(p_stdout_line_str):
		p_lst = p_stdout_line_str.split(' ')

		assert len(p_lst) == 3
		image_id_str = p_lst[2]

		#IMPORTANT!! - check that this is a valid 12 char Docker ID
		assert len(image_id_str) == 12
		return image_id_str
	#---------------------------------------------------

	for line in r.stdout:
		line_str = line.strip() #strip() - to remove '\n' at the end of the line

		#------------------
		#display the line, to update terminal display
		print line_str
		#------------------

		if line_str.startswith('Successfully built'):
			image_id_str = get_image_id_from_line(line_str)
			return image_id_str

	#change back to old dir
	os.chdir(old_cwd)



