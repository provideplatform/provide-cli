FROM golang:1.15 AS builder
RUN mkdir -p /go/src/github.com/provideplatform
ADD . /go/src/github.com/provideplatform/provide-cli
WORKDIR /go/src/github.com/provideplatform/provide-cli
RUN make build

FROM alpine
RUN apk add --no-cache bash
WORKDIR /
COPY --from=builder /go/src/github.com/provideplatform/provide-cli/.bin/prvd /prvd
ENTRYPOINT ["./prvd"]
