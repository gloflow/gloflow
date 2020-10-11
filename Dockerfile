FROM alpine:3.12.0

RUN apk update
RUN apk upgrade

RUN adduser --home /home/gf --disabled-password gf
RUN mkdir -p /home/gf/bin
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/config

RUN apk add --update bash
# RUN apk add --no-cache libc6-compat

#------------
# SUPERVISOR
RUN apk --update add supervisor

ADD ./config/supervisor.conf /home/gf/config/supervisor.conf

# "-c" - config file path
# "-n" - --nodaemon -- run in the foreground (same as 'nodaemon true' in config file).
#        this way the supervisor will run as the main process, and the container wont exit
CMD /usr/bin/supervisord -n -c /home/gf/config/supervisor.conf

# process supervisord events
ADD build/gf_supervisord_events.py /home/gf/bin/gf_supervisord_events.py

#------------
ADD config/gf_eth_monitor.yaml /home/gf/config/gf_eth_monitor.yaml
ADD build/gf_eth_monitor       /home/gf/bin/gf_eth_monitor
ADD static                     /home/gf/static

RUN chown -R gf /home/gf