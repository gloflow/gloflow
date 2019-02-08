
#FIX!! - when using "gliderlabs/alpine:3.8" the container build process was failing at "apk upgrade".
#        figure out how to fix and use latest alpine 3.8
FROM gliderlabs/alpine:3.7

RUN apk update
RUN apk upgrade
RUN apk add bash
#------------
#DIRS
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/conf
RUN mkdir -p /home/gf/data
#------------
ADD ./templates            /home/gf/templates
ADD ./static               /home/gf/static
ADD ./gf_publisher_service /home/gf/gf_publisher_service