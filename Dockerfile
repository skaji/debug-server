FROM golang:1.15
LABEL org.opencontainers.image.source https://github.com/skaji/debug-server

WORKDIR /go/src/github.com/skaji/debug-server
ENV DEBIAN_FRONTEND noninteractive

COPY go.* *.go ./

RUN set -eux; \
  go build; \
  curl -fsSL -o /sbin/tini https://github.com/krallin/tini/releases/download/v0.19.0/tini; \
  chmod +x /sbin/tini; \
  apt-get update; \
  apt-get install -y net-tools dnsutils; \
  :

EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/go/src/github.com/skaji/debug-server/debug-server"]
