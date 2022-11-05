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

#------------
# MAIN
ADD ./go/gf_p2p/gf_p2p_tester/gf_p2p_tester /home/gf/
ADD ./go/gf_p2p/gf_p2p_tester/web           /home/gf/

#------------
CMD ["/home/gf/gf_p2p_tester"]