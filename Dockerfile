FROM ubuntu:18.04

ENV DEBIAN_FRONTEND="noninteractive"

RUN apt-get update && apt-get -y install software-properties-common && add-apt-repository ppa:longsleep/golang-backports
RUN apt-get -y install r-base build-essential golang-go libgit2-dev libssl-dev libxml2-dev

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN pwd
RUN ls -l
RUN make

CMD [ "/usr/bin/bash" ]
