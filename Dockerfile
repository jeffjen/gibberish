FROM golang:latest
MAINTAINER YI-HUNG JEN <yihungjen@gmail.com>

COPY main.go /go/src/github.com/jeffjen/gibberish/
COPY server /go/src/github.com/jeffjen/gibberish/server

COPY public /var/public

RUN go get github.com/Sirupsen/logrus
RUN go install github.com/jeffjen/gibberish

EXPOSE 9090

ENTRYPOINT ["gibberish"]
