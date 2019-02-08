
#FIX!! - when using "gliderlabs/alpine:3.8" the container build process was failing at "apk upgrade".
#        figure out how to fix and use latest alpine 3.8
FROM gliderlabs/alpine:3.7

RUN apk add --update bash

ADD ./gf_tagger_service /home/gf/gf_tagger_service