
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
RUN mkdir -p /home/gf/data/images
RUN mkdir -p /home/gf/data/images/thumbnails
RUN mkdir -p /home/gf/static

RUN mkdir -p /home/gf/templates
RUN mkdir -p /home/gf/templates/subtemplates
#------------
#IMPORTANT!! - without ca-certificates when the fetcher runs a GET on the image url, 
#              if its a https:// address, it will return an error
#              "x509: failed to load system roots and no roots provided"
RUN apk add ca-certificates
#------------
ADD ./templates/gf_images_flows_browser.html /home/gf/templates/gf_images_flows_browser.html
ADD ./static                                 /home/gf/static
ADD ./gf_images_service                      /home/gf/gf_images_service