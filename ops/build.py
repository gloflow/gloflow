import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))




from colored import fg,bg,attr
import delegator


print ''
print '                              %sBUILD GLOFLOW%s'%(fg('green'),attr(0))
print ''


#---------------------------------
def build__gf_images():
    print ' -- build %sgf_images%s service'%(fg('green'),attr(0))
    gf_images_service__path_str = '%s/../go/apps/gf_images/gf_images_service.go'%(cwd_str)
    output_path_str             = '%s/../bin/gf_images_service'%(cwd_str)
    c                           = 'go build -o %s %s'%(output_path_str,gf_images_service__path_str,)
    print c
    r = delegator.run(c)
    if not r.out == '': print r.out
    if not r.err == '': print '%sFAILED%s >>>>>>>\n%s'%(fg('red'),attr(0),r.err)
#---------------------------------


build__gf_images()