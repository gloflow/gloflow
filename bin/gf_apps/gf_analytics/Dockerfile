
#FIX!! - when using "gliderlabs/alpine:3.8" the container build process was failing at "apk upgrade".
#        figure out how to fix and use latest alpine 3.8
FROM gliderlabs/alpine:3.7

RUN apk update
RUN apk add --update bash

#IMPORTANT!! - without ca-certificates when requests are made to external urls, 
#              if its a https:// address, it will return an error
#              "x509: failed to load system roots and no roots provided"
RUN apk add ca-certificates

#------------
#DIRS
RUN mkdir -p /home/gf/static/

#for crawler image data
RUN mkdir -p /home/gf/image_data 
#------------
#PYTHON CODE

RUN apk --update add python2

RUN apk --update add py-pip
RUN pip install --upgrade pip

#MATPLOTLIB DEPENDENCIES
#RUN apk add alpine-sdk #adds ~400mb to container size!
RUN apk add build-base
RUN apk add python-dev
RUN apk add freetype-dev
RUN apk add libpng-dev

RUN pip install matplotlib
RUN pip install pandas
RUN pip install pymongo
RUN pip install delegator-py
RUN pip install colored

#PY SCRIPTS
ADD ./py/cli_stats.py /home/gf/py/cli_stats.py
ADD ./py/stats/       /home/gf/py/stats
#------------
ADD ./templates/gf_domains_browser.html /home/gf/templates/gf_domains_browser.html
ADD ./static                            /home/gf/static
ADD ./gf_analytics_service              /home/gf/