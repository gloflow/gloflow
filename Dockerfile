FROM alpine:3.12.0

RUN mkdir -p /home/gf/bin
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/config

#------------
# SUPERVISOR
RUN apk --update add supervisor

ADD ./config/supervisor.conf /home/gf/config/supervisor.conf

# "-c" - config file path
# "-n" - --nodaemon -- run in the foreground (same as 'nodaemon true' in config file).
#        this way the supervisor will run as the main process, and the container wont exit
CMD /usr/bin/supervisord -n -c /home/gf/config/supervisor.conf

#------------

ADD build/gf_eth_monitor /home/gf/bin/gf_eth_monitor