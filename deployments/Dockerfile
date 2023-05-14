FROM ubuntu:22.10
LABEL maintainer="Sergey Turchenkov <ssturch@gmail.com>"

COPY ["Baltimore CyberTrust Root.crt", "/usr/local/share/ca-certificates"]
RUN apt-get update
RUN apt-get install ca-certificates -y
RUN update-ca-certificates
COPY xtapi /root
#EXPOSE 8080
#EXPOSE 5432
