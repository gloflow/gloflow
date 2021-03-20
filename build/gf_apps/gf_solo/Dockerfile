FROM alpine:3.12.0

RUN apk update
RUN apk upgrade
RUN apk add bash

#------------
# DIRS
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/config
RUN mkdir -p /home/gf/data

#------------
ADD ./templates /home/gf/templates
ADD ./static    /home/gf/static
ADD ./gf_solo   /home/gf/gf_solo


ADD ./gf_analytics /home/gf/gf_analytics
ADD ./gf_ml_worker /home/gf/gf_ml_worker