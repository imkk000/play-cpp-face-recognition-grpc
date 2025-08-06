FROM archlinux:latest AS builder

RUN pacman -Syu --noconfirm && \
  pacman -S --noconfirm \
  base-devel \
  cmake \
  grpc \
  protobuf \
  opencv \
  pkg-config

WORKDIR /app

COPY . .

RUN make compile_minimal

FROM archlinux:latest

RUN pacman -Syu --noconfirm && \
  pacman -S --noconfirm \
  grpc \
  protobuf \
  opencv && \
  pacman -Scc --noconfirm

COPY --from=builder /app/server /usr/local/bin/
COPY --from=builder /app/models /models

EXPOSE 8080

USER nobody

ENTRYPOINT ["/usr/local/bin/server"]
