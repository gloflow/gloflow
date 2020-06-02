FROM ubuntu:18.04

RUN apt-get update


RUN mkdir -p /home/gf/py

#------------
# PIP
RUN apt-get install -y python3-pip

ADD ./py/requirements.txt /home/gf/py
RUN pip3 install -r /home/gf/py/requirements.txt

#------------
# PY CODE AND .so LIBS
ADD ./py /home/gf/py

# FIX!! - this should be compiled staticaly into gf_images_jobs_py.so
RUN apt-get update
# RUN apt-get install -y libcairo2-dev

#------------
# RUN
# "-u" - unbuffered stdout/stderr

WORKDIR /home/gf/py
CMD ls -al /home/gf/py
ENTRYPOINT LD_LIBRARY_PATH=. python3 -u /home/gf/py/gf_simple_model.py