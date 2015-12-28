FROM golang:1.5-alpine

ADD . /go/src/github.com/daraszkrzysztof/secure-voting
RUN apk --update add git
RUN go get github.com/emicklei/go-restful
RUN go install github.com/daraszkrzysztof/secure-voting
ENTRYPOINT /go/bin/secure-voting

EXPOSE 8080