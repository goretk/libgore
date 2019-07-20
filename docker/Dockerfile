# FROM karalabe/xgo-latest
FROM golang
MAINTAINER Go Reverse Engineering Tool Kit
RUN \
  apt-get update && \
  apt-get install -y automake autogen build-essential ca-certificates gcc-mingw-w64 make \
    --no-install-recommends

COPY build.sh /build.sh
ENTRYPOINT ["/build.sh"]
