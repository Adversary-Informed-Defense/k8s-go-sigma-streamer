FROM alpine:latest AS base
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add --no-cache gcompat
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache librdkafka

FROM golang:1.16.0 AS build
WORKDIR /opt

# Install dependencies
# Note: this is a duplicate copy to optimize
# build speed by caching dependency layers
ENV GO111MODULE=on
COPY kafka/go.mod kafka/go.sum ./
RUN go mod download

# Build executable
COPY kafka/ .
RUN go build -o singe main.go

FROM base
WORKDIR /opt
COPY --from=build /opt/singe .
COPY rules rules
CMD ["./main"]
