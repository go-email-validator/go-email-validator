FROM node:current-buster-slim
RUN apt-get update && apt-get install -y build-essential git curl ca-certificates openssl

ARG cert_location=/usr/local/share/ca-certificates

# Get certificate from "github.com"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates

RUN curl https://download.docker.com/linux/static/stable/x86_64/docker-20.10.1.tgz --output docker.tar.gz
RUN tar xzvf docker.tar.gz
RUN cp docker/* /usr/bin/
RUN dockerd &