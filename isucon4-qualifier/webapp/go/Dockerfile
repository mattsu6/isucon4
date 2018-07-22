FROM matsuu/isucon4-qualifier-base:latest

MAINTAINER matsuu@gmail.com

COPY webapp /home/isucon/webapp/go
WORKDIR /home/isucon/webapp/go

RUN \
  curl -L http://golang.org/dl/go1.10.3.linux-amd64.tar.gz | tar zxf - -C /usr/local && \
  /home/isucon/env.sh ./build.sh 

EXPOSE 8080
CMD ["/home/isucon/env.sh", "./golang-webapp"]