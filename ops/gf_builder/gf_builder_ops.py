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

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

sys.path.append('%s/../containers'%(cwd_str))
import gf_containers

#---------------------------------------------------
def cont__build(p_log_fun):
    
    image_name_str          = "gf_builder_org"
    image_tag_str           = "latest"
    docker_context_dir_str  = "%s/../../build/gf_builder"%(cwd_str)
    dockerhub_user_name_str = "gloflow"

    gf_containers.build_docker_image(image_name_str,
        image_tag_str,
		docker_context_dir_str,
		dockerhub_user_name_str,
		p_log_fun)