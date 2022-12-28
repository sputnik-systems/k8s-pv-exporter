FROM golang:1.19.0-buster as build

RUN apt update \
    && apt install -y curl

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -o ./exporter ./cmd/exporter


FROM ubuntu

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/exporter /usr/local/bin/pv-exporter

ENTRYPOINT ["/usr/local/bin/pv-exporter"]
