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
import signal
import subprocess

from colored import fg,bg,attr

#--------------------------------------------------
def run(p_app_name_str,
    p_test_name_str,
    p_app_meta_map,
    p_aws_s3_creds_map,
    p_exit_on_fail_bool     = False,
    p_test_mongodb_host_str = "127.0.0.1"):
    assert isinstance(p_test_name_str,    basestring)
    assert isinstance(p_app_meta_map,     dict)
    assert isinstance(p_aws_s3_creds_map, dict)

    print ''
    print ' -- test %s%s%s package'%(fg('green'), p_app_name_str, attr(0))

    if p_app_meta_map.has_key('test_data_to_serve_dir_str'): use_test_server_bool = True
    else:                                                    use_test_server_bool = False

    #GO_PACKAGE_DIR
    go_package_dir_path_str = p_app_meta_map['go_path_str']
    assert os.path.isdir(go_package_dir_path_str)
    print("go_package_dir_path_str - %s"%(go_package_dir_path_str))

    #-------------
    # TEST_SERVER - used to server assets/images that various Go functions
    #               that are tested that do fetching of remote resources.
    
    if use_test_server_bool:
        test_data_dir_str = p_app_meta_map['test_data_to_serve_dir_str']
        assert os.path.isdir(test_data_dir_str)

        print('')
        print('STARTING TEST DATA PYTHON HTTP SERVER ----------------------------')

        #run the python simple server in the dir where the test data is located, so that its served over http
        c_lst = ["cd %s && python -m SimpleHTTPServer 8000"%(test_data_dir_str)]
        print(' '.join(c_lst))

        # IMPORTANT!! - "cd" and py server are run by the shell, which is their parent process, so a session ID is attached
        #               so that its made a group leader. later when the os.killpg() termination signal is sent to that 
        #               group leader (the shell), its child processes will get shutdown as well (py server).
        #               otherwise the py server will be left running after the tests have finished
        server_p = subprocess.Popen(c_lst, stdout=subprocess.PIPE, preexec_fn=os.setsid, shell=True)
        
    #-------------
    cwd_str = os.getcwd()
    os.chdir(go_package_dir_path_str) #change into the target main package dir

    #-------------
    # CMD
    # ADD!! - per app timeout values, so that in gf_meta.py
    #         a test timeout can be specified per app/package.
    cmd_lst = [
        "go test",
        "-timeout 30s",
        "-mongodb_host=%s"%(p_test_mongodb_host_str)
    ]

    # specific test was selected for running, not all tests
    if not p_test_name_str == 'all':
        cmd_lst.append("-v") #verbose
        cmd_lst.append(p_test_name_str)


    c = " ".join(cmd_lst)
    print(c)
    #-------------

    e = os.environ.copy()
    e.update(p_aws_s3_creds_map)
    p = subprocess.Popen(c.split(' '), stderr=subprocess.PIPE, env=e)
    
    # IMPORTANT!! - stderr is used and read, because thats what contains the log lines from Go programs that has
    #               color codes preserved in log lines.
    for l in iter(p.stderr.readline, ""):
        print(l.rstrip())

    #if not p.stderr == None: print '%sTESTS FAILED%s >>>>>>>\n'%(fg('red'), attr(0))

    p.wait() # has to be called so that p.returncode is set
    os.chdir(cwd_str) # return to initial dir
    #-------------

    # kill HTTP test server used to serve assets that need to come over HTTP
    if use_test_server_bool:
        os.killpg(server_p.pid, signal.SIGTERM)


    # in certain scenarios (such as CI) we want this test run to fail 
    # completelly in case "go test" returns a non-zero return code (failed test).
    # this way CI pipeline will get stoped and marked as failed.
    if p_exit_on_fail_bool:
        print("test exited with code - %s"%(p.returncode))
        assert not p.returncode == None # makesure returncode is set by p.wait()
        
        if not p.returncode == 0:
            exit(p.returncode)
