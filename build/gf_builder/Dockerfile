


FROM ubuntu:18.04

RUN mkdir /home/gf


RUN apt-get update
RUN apt-get install -y python






ADD gf_builder.py /home/gf














ENTRYPOINT ["python", "/home/gf/gf_builder.py"]