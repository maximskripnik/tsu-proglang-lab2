FROM golang:1.13-buster

WORKDIR /opt/app
ENV GOPATH /opt/app

COPY src /opt/app/src

RUN go get all
RUN go build main

ENTRYPOINT [ "./main" ]