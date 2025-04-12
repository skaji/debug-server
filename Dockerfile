FROM golang:latest AS builder

WORKDIR /app
RUN --mount=target=. go build -o /debug-server

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN <<EOF bash
set -eux

apt-get update
apt-get install -y \
  curl \
  dnsutils \
  net-tools \
  perl \
  procps \
  tzdata \
  ;
apt-get clean
rm -rf /var/lib/apt/lists/*

curl -fsSL -o /sbin/tini https://github.com/krallin/tini/releases/download/v0.19.0/tini
chmod +x /sbin/tini

echo Asia/Tokyo > /etc/timezone
rm -f /etc/localtime
dpkg-reconfigure -f noninteractive tzdata
EOF

COPY --from=builder /debug-server /debug-server

WORKDIR /root
EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/debug-server"]
