# syntax=docker/dockerfile:1.6

ARG RUST_VERSION=1.83
ARG EXAMPLE_NAME=server

FROM rust:${RUST_VERSION}-slim AS builder
ARG EXAMPLE_NAME
WORKDIR /app/h3
RUN apt-get update && \
    apt-get install -y --no-install-recommends pkg-config libssl-dev ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY h3 ./
RUN cargo build --release --example "${EXAMPLE_NAME}"

FROM debian:bookworm-slim AS runtime
ARG EXAMPLE_NAME
WORKDIR /app/h3
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -m appuser
RUN mkdir -p /app/h3/examples /app/h3/content/root && \
    chown -R appuser:appuser /app/h3/content
COPY --from=builder --chown=appuser:appuser /app/h3/target/release/examples/${EXAMPLE_NAME} /usr/local/bin/${EXAMPLE_NAME}
COPY --from=builder --chown=appuser:appuser /app/h3/examples/server.cert /app/h3/examples/server.cert
COPY --from=builder --chown=appuser:appuser /app/h3/examples/server.key /app/h3/examples/server.key
RUN mkdir -p /dist && \
    cp /usr/local/bin/${EXAMPLE_NAME} /dist/${EXAMPLE_NAME} && \
    chown appuser:appuser /dist/${EXAMPLE_NAME}
USER appuser
EXPOSE 4433
ENTRYPOINT ["/usr/local/bin/server"]
CMD ["--listen=0.0.0.0:4433", "--dir=/app/h3/content/root"]
