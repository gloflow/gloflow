# GloFlow media management/publishing system
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
def run(p_name_str,
    p_app_meta_map,
    p_aws_s3_creds_map):
    assert isinstance(p_app_meta_map,    dict)
    assert isinstance(p_aws_s3_creds_map,dict)

    print ''
    print ' -- test %s%s%s package'%(fg('green'), p_name_str, attr(0))

    if p_app_meta_map.has_key('test_data_to_serve_dir_str'): use_test_server_bool = True
    else:                                                    use_test_server_bool = False

    go_package_dir_path_str = p_app_meta_map['go_path_str']
    assert os.path.isdir(go_package_dir_path_str)

    #-------------
    #TEST_SERVER
    if use_test_server_bool:
        test_data_dir_str = p_app_meta_map['test_data_to_serve_dir_str']
        assert os.path.isdir(test_data_dir_str)

        print('')
        print('STARTING TEST DATA PYTHON HTTP SERVER ----------------------------')
        c = '(cd %s && python -m SimpleHTTPServer 8000)'%(test_data_dir_str)
        print(c)

        #run the python simple server in the dir where the test data is located, so that its served over http
        c_lst = ["cd %s && python -m SimpleHTTPServer 8000"%(test_data_dir_str)]
        print(' '.join(c_lst))

        #IMPORTANT!! - "cd" and py server are run by the shell, which is their parent process, so a session ID is attached
        #              so that its made a group leader. later when the os.killpg() termination signal is sent to that 
        #              group leader (the shell), its child processes will get shutdown as well (py server).
        #              otherwise the py server will be left running after the tests have finished
        server_p = subprocess.Popen(c_lst, stdout=subprocess.PIPE, preexec_fn=os.setsid, shell=True)
        
    #-------------
    cwd_str = os.getcwd()
    os.chdir(go_package_dir_path_str) #change into the target main package dir

    c = "go test"
    print(c)

    e = os.environ.copy()
    e.update(p_aws_s3_creds_map)
    p = subprocess.Popen(c.split(' '), stdout=subprocess.PIPE, env=e)

    for l in iter(p.stdout.readline, ""):
        print(l.rstrip())

    if not p.stderr == None: print '%sFAILED%s >>>>>>>\n%s'%(fg('red'), attr(0), p.stderr)

    os.chdir(cwd_str) #return to initial dir
    #-------------

    #kill HTTP test server used to serve assets that need to come over HTTP
    if use_test_server_bool: os.killpg(server_p.pid, signal.SIGTERM)