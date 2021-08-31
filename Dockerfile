FROM golang:1.17
LABEL org.opencontainers.image.source https://github.com/skaji/debug-server
ARG VERSION=v0.0.1

WORKDIR /go/src/github.com/skaji/debug-server

COPY go.* *.go ./

RUN set -eux; \
  go build -ldflags "-X main.version=$VERSION"; \
  curl -fsSL -o /sbin/tini https://github.com/krallin/tini/releases/download/v0.19.0/tini; \
  chmod +x /sbin/tini; \
  export DEBIAN_FRONTEND=noninteractive; \
  apt-get update; \
  apt-get install -y net-tools dnsutils; \
  :

EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/go/src/github.com/skaji/debug-server/debug-server"]
