FROM golang:1.15 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

ARG version="undefined"

WORKDIR /build/ip2location

ADD . /build/ip2location

RUN go build -o /ip2location -ldflags "-X main.version=${version} -s -w"  ./cmd/ip2location

# -----

FROM  debian:stretch-slim
COPY --from=build /ip2location /

RUN apt-get update \
     && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

ENTRYPOINT ["/ip2location"]
