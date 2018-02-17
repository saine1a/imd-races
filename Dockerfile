FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/saine1a/imd-races
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o imd-races .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/saine1a/imd-races/imd-races .
CMD ["./imd-races"]
