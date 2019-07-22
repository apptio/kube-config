FROM golang:1.11.4-alpine3.8

WORKDIR /app

COPY . /app

ARG VERSION

RUN apk add git build-base \
    && go mod download \
    && go build -o kube-config -ldflags "-X main.version=${VERSION}"

ENTRYPOINT ["/app/kube-config"]
