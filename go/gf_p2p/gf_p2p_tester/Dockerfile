FROM ubuntu:22.04

#------------
# GROUP/USER
RUN groupadd -r gf
RUN useradd --create-home --no-log-init -r -g gf gf

#------------
RUN apt-get update
RUN apt-get install -y \
    python3 \
    python3-pip

#------------

WORKDIR /home/gf
USER gf

#------------
# DIRS
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/config
RUN mkdir -p /home/gf/web/templates

#------------
# WEB
ADD web/templates web/templates

#------------
# MAIN
ADD ./gf_p2p_tester gf_p2p_tester

#------------
CMD ["/home/gf/gf_p2p_tester"]