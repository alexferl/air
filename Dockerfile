ARG GOLANG_VERSION=1.17.7
FROM golang:${GOLANG_VERSION} AS builder
MAINTAINER Alexandre Ferland <me@alexferl.com>

RUN groupadd -g 1337 appuser && \
    useradd -r -d /app -u 1337 -g appuser appuser

WORKDIR /build

RUN apt-get update && apt-get install -y \
    libvips-dev

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build
RUN mv /build/air /air

USER appuser

EXPOSE 1323

ENTRYPOINT ["/air"]
