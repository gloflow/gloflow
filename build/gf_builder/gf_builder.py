










print("BUILD...")



import delegator
print(delegator.run("ls -al").out)
print(delegator.run("pwd").out)





import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))





sys.path.append('%s/../../meta'%(cwd_str))
import gf_meta


sys.path.append('%s/../../ops/utils'%(cwd_str))
import gf_build_changes



print("DIFF")
apps_changes_deps_map = gf_meta.get()['apps_changes_deps_map']


changed_apps_map = gf_build_changes.list_changed_apps(apps_changes_deps_map)
gf_build_changes.view_changed_apps(changed_apps_map)