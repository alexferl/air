FROM golang:1.17.6
MAINTAINER Alexandre Ferland <me@alexferl.com>

WORKDIR /build

RUN apt-get update \
            && apt-get install -y \
            libvips-dev

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build
RUN mv air /air

ENTRYPOINT ["/air"]

EXPOSE 1323
CMD []
