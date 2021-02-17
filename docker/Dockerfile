# FROM karalabe/xgo-latest
FROM golang:latest
MAINTAINER Go Reverse Engineering Tool Kit
RUN \
  apt-get update && \
  apt-get install -y automake autogen build-essential ca-certificates clang \
    cmake gcc-mingw-w64 make wget git libxml2-dev libssl-dev libz-dev \
    --no-install-recommends

# Instructions taken from https://github.com/karalabe/xgo
# and https://wapl.es/rust/2019/02/17/rust-cross-compile-linux-to-macos.html
ENV OSX_SDK     MacOSX10.11.sdk
ENV OSX_NDK_X86 /usr/local/osx-ndk-x86

RUN \
  OSX_SDK_PATH=https://s3.dockerproject.org/darwin/v2/$OSX_SDK.tar.xz && \
  wget -nc $OSX_SDK_PATH && \
  \
  git clone https://github.com/tpoechtrager/osxcross.git && \
  mv `basename $OSX_SDK_PATH` osxcross/tarballs/        && \
  \
  PORTABLE=yes UNATTENDED=yes OSX_VERSION_MIN=10.6 osxcross/build.sh               && \
  mv osxcross/target $OSX_NDK_X86                                                   && \
  \
  rm -rf ./osxcross

ENV PATH $OSX_NDK_X86/bin:$PATH
ENV LD_LIBRARY_PATH $OSX_NDK_X86/lib:LD_LIBRARY_PATH

COPY build.sh /build.sh
ENTRYPOINT ["/build.sh"]
