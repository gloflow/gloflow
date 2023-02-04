



install the latest docker version (>=23.0.0)
    - needed to be able to build the latest ubuntu:20.10 containers
    - containers used for the gf_builder instances for Go and Rust

easy way to install is via static docker binaries:
- https://download.docker.com/linux/static/stable/x86_64/
- download `docker-23.0.0.tgz` file

> cp ~/Download ./t/
> tar -xvzf docker-23.0.0.tgz
> cp docker/* /usr/bin/