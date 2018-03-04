FROM golang:1.7

RUN go get github.com/saine1a/imd-races

RUN go install github.com/saine1a/imd-races

ENTRYPOINT /go/bin/imd-races

EXPOSE 8080
