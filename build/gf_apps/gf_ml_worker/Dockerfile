



FROM ubuntu:18.04


RUN apt-get update


RUN mkdir -p /home/gf/py




RUN apt-get install -y python3-pip


#------------
# PY SCRIPTS
ADD ./py/gf_simple_model.py /home/gf/py/
ADD ./py/requirements.txt   /home/gf/py/






RUN pip3 install -r /home/gf/py/requirements.txt




#------------
# RUN
# "-u" - unbuffered stdout/stderr
CMD python3 -u /home/gf/py/gf_simple_model.py