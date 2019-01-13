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

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

from colored import fg,bg,attr
import delegator

print ''
print '                              %sBUILD GLOFLOW%s'%(fg('green'),attr(0))
print ''

#---------------------------------
#META
gf_images_service__path_str        = '%s/../go/apps/gf_images/gf_images_service.go'%(cwd_str)
gf_images_service__output_path_str = '%s/../bin/gf_images_service'%(cwd_str)

gf_publisher_service__path_str        = '%s/../go/apps/gf_publisher/gf_publisher_service.go'%(cwd_str)
gf_publisher_service__output_path_str = '%s/../bin/gf_publisher_service'%(cwd_str)

gf_tagger_service__path_str        = '%s/../go/apps/gf_tagger/gf_tagger_service.go'%(cwd_str)
gf_tagger_service__output_path_str = '%s/../bin/gf_tagger_service'%(cwd_str)

gf_landing_page_service__path_str        = '%s/../go/apps/gf_landing_page/gf_landing_page_service.go'%(cwd_str)
gf_landing_page_service__output_path_str = '%s/../bin/gf_landing_page_service'%(cwd_str)

gf_analytics_service__path_str        = '%s/../go/apps/gf_analytics/gf_analytics_service.go'%(cwd_str)
gf_analytics_service__output_path_str = '%s/../bin/gf_analytics_service'%(cwd_str)
#---------------------------------
def build__go_bin(p_name_str,
    p_main_go_file_path_str,
    p_output_path_str):

    print ' -- build %s%s%s service'%(fg('green'), p_name_str, attr(0))
    
    cwd_str = os.getcwd()
    os.chdir(os.path.dirname(p_main_go_file_path_str)) #change into the target main package dir

    c = 'go build -o %s'%(p_output_path_str)
    print c
    r = delegator.run(c)
    if not r.out == '': print r.out
    if not r.err == '': print '%sFAILED%s >>>>>>>\n%s'%(fg('red'),attr(0),r.err)

    os.chdir(cwd_str) #return to initial dir
#---------------------------------
build__go_bin('gf_image_service',       gf_images_service__path_str,      gf_images_service__output_path_str)
build__go_bin('gf_publisher_service',   gf_publisher_service__path_str,   gf_publisher_service__output_path_str)
build__go_bin('gf_tagger_service',      gf_tagger_service__path_str,      gf_tagger_service__output_path_str)
build__go_bin('gf_landing_page_service',gf_landing_page_service__path_str,gf_landing_page_service__output_path_str)
build__go_bin('gf_analytics_service',   gf_landing_page_service__path_str,gf_landing_page_service__output_path_str)