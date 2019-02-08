
#FIX!! - when using "gliderlabs/alpine:3.8" the container build process was failing at "apk upgrade".
#        figure out how to fix and use latest alpine 3.8
FROM gliderlabs/alpine:3.7

RUN apk add --update bash

#------------
#DIRS
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/conf
RUN mkdir -p /home/gf/data
#------------
ADD ./templates               /home/gf/templates
ADD ./gf_landing_page_service /home/gf/gf_landing_page_service