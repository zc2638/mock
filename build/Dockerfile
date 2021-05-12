FROM golang:1.16-alpine as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o mock cmd/main.go

FROM gruebel/upx:latest as upx
WORKDIR /work
COPY --from=builder /work/mock /work/mock.pre
RUN upx --best --lzma -o /work/mock mock.pre

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine
LABEL maintainer="zc" \
    email="zc2638@qq.com" \
    version="v1.0.0"

ENV TS="Asia/Shanghai"
#ENV GIN_MODE=release
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /work/config.yaml /work/config.yaml
COPY --from=upx /work/mock /work/mock

WORKDIR /work
CMD ["./mock"]