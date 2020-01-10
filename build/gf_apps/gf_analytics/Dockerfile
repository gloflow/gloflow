
#FIX!! - when using "gliderlabs/alpine:3.8" the container build process was failing at "apk upgrade".
#        figure out how to fix and use latest alpine 3.8
FROM gliderlabs/alpine:3.7

RUN apk update
RUN apk add --update bash
#------------
#DIRS
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/config
RUN mkdir -p /home/gf/data

#for crawler image data
RUN mkdir -p /home/gf/data/images
RUN mkdir -p /home/gf/data/images/thumbnails
RUN mkdir -p /home/gf/static

#templates
RUN mkdir -p /home/gf/templates
RUN mkdir -p /home/gf/templates/subtemplates
#------------
#IMPORTANT!! - without ca-certificates when the fetcher runs a GET on the image url, 
#              if its a https:// address, it will return an error
#              "x509: failed to load system roots and no roots provided"
RUN apk add ca-certificates
#------------
#PYTHON CODE

RUN apk --update add python3 \
    python3-dev

RUN apk --update add py-pip
RUN pip install --upgrade pip

#MATPLOTLIB DEPENDENCIES
#RUN apk add alpine-sdk #adds ~400mb to container size!
RUN apk add build-base
RUN apk add freetype-dev
RUN apk add libpng-dev

# cython needed for matplotlib
RUN python3 -mpip install cython
RUN python3 -mpip install --upgrade matplotlib
RUN python3 -mpip install pandas
RUN python3 -mpip install pymongo
RUN python3 -mpip install delegator-py
RUN python3 -mpip install colored
#------------
#SUPERVISOR
RUN apk --update add supervisor
#------------
#"-c" - config file path
#"-n" - --nodaemon -- run in the foreground (same as 'nodaemon true' in config file).
#       this way the supervisor will run as the main process, and the container wont exit
CMD /usr/bin/supervisord -n -c /home/gf/conf/supervisor.conf
#------------
ADD ./supervisor.conf /home/gf/conf/supervisor.conf
#------------
#PY SCRIPTS
ADD ./py/cli_stats.py /home/gf/py/cli_stats.py
ADD ./py/stats/       /home/gf/py/stats

#WEB
ADD ./templates /home/gf/templates
ADD ./static    /home/gf/static

#default crawler definitions
ADD ./config/crawl_config.yaml /home/gf/config
ADD ./gf_analytics_service     /home/gf/