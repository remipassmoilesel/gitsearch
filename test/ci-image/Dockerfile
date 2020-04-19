FROM ubuntu:bionic

RUN apt-get update && apt-get install -y build-essential bash software-properties-common apt-transport-https \
    ca-certificates curl gnupg-agent software-properties-common

RUN add-apt-repository ppa:longsleep/golang-backports \
    && apt-get install -y golang-go

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash - \
    && apt-get install -y nodejs

ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOPATH/bin

RUN npm i -g yarn \
    && apt-get install -y git \
    && go get github.com/markbates/pkger/cmd/pkger
