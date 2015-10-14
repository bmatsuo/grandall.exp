FROM ubuntu:14.04.1

MAINTAINER "Bryan Matsuo <bryan.matsuo@gmail.com>"

# for now precompiled binaries are used because of loose external dependencies.
RUN mkdir -p /opt/bin
RUN apt-get update -q && apt-get install -qy wget
RUN wget \
    -O /tmp/grandall.exp0.2-linux-amd64.tar.gz \
    https://github.com/bmatsuo/grandall.exp/releases/download/0.2/grandall.exp0.2-linux-amd64.tar.gz
RUN tar -C /tmp -xvzf /tmp/grandall.exp0.2-linux-amd64.tar.gz
RUN cp /tmp/grandall.exp0.2-linux-amd64/grandalld /opt/bin/grandalld

RUN mkdir -p /etc/grandall
COPY grandalld.conf /etc/grandall/grandalld.conf

RUN mkdir -p /etc/grandall/sites-enabled
RUN mkdir -p /etc/grandall/sites-available

ENTRYPOINT ["/opt/bin/grandalld",  "-config=/etc/grandall/grandalld.conf", "-sites=/etc/grandall/sites-enabled"]
