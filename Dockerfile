FROM golang
ADD . /go/src/github.com/daraszkrzysztof/secure-voting
RUN go install github.com/daraszkrzysztof/secure-voting
ENTRYPOINT /go/bin/secure-voting
EXPOSE 8080