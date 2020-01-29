
# FIX!! - when using "gliderlabs/alpine:3.8" the container build process was failing at "apk upgrade".
#         figure out how to fix and use latest alpine 3.8
FROM gliderlabs/alpine:3.7

RUN apk update
RUN apk upgrade
RUN apk add bash
#------------
# DIRS
RUN mkdir -p /home/gf/logs
RUN mkdir -p /home/gf/config
RUN mkdir -p /home/gf/static

# used for fetched images, and intermediate processing results
RUN mkdir -p /home/gf/images
RUN mkdir -p /home/gf/images/thumbnails
#------------
# SUPERVISOR
RUN apk --update add supervisor
#------------
# "-c" - config file path
# "-n" - --nodaemon -- run in the foreground (same as 'nodaemon true' in config file).
#        this way the supervisor will run as the main process, and the container wont exit
CMD /usr/bin/supervisord -n -c /home/gf/config/supervisor.conf
#------------
# IMPORTANT!! - without ca-certificates when the fetcher runs a GET on the image url, 
#               if its a https:// address, it will return an error
#               "x509: failed to load system roots and no roots provided"
RUN apk add ca-certificates
#------------
ADD ./config/supervisor.conf       /home/gf/config/supervisor.conf
ADD ./config/gf_images_config.yaml /home/gf/config/gf_images_config.yaml
ADD ./templates                    /home/gf/templates
ADD ./static                       /home/gf/static
ADD ./gf_images_service            /home/gf/gf_images_service