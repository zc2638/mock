FROM golang:1.15.1-alpine as builder
ENV GOPROXY https://goproxy.io
ENV GO111MODULE=on
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download
WORKDIR /app
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o mock cmd/main.go

FROM alpine
#ENV GIN_MODE=release
COPY --from=builder /app/config.yaml /config.yaml
COPY --from=builder /app/mock /mock
CMD ["/mock"]