










print("BUILD...")



import delegator
print(delegator.run("ls -al").out)
print(delegator.run("pwd").out)
print(delegator.run("ls -al /home/gf").out)





import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))



sys.path.append('%s'%(cwd_str))