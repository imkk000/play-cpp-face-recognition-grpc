FROM debian:bookworm-slim AS builder

RUN apt-get update && apt-get install -y \
  libopencv-dev libgrpc++-dev libprotobuf-dev \
  make pkgconf g++

WORKDIR /src
COPY . .

RUN make compile
