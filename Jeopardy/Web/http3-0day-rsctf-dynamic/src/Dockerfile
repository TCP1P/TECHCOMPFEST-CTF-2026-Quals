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
# Original final-stage ENTRYPOINT is wrapped by RSCTF at runtime.
CMD ["--listen=0.0.0.0:4433", "--dir=/app/h3/content/root"]
# Bake the legacy Compose flag mount into RSCTF's single-image runtime.
COPY ["flag.txt", "/app/h3/flag.txt"]

# RSCTF dynamic-flag compatibility layer. The build-time scan records
# exact marker-bearing files; runtime replacement is therefore bounded
# to those files instead of walking the image on every team start.
USER root
COPY rsctf-entrypoint.sh /usr/local/bin/rsctf-entrypoint
RUN set -eu; \
    chmod 0555 /usr/local/bin/rsctf-entrypoint; \
    : > /etc/rsctf-flag-targets; \
    for root in /app/h3/flag.txt /app/h3/examples/server.key /app/h3/examples/server.cert /app/h3 /app /src /srv /home /opt /var/www /challenge /chall /ctf /root /usr/src /usr/local/bin /workspace /flag /flag.txt; do \
      if [ -e "$root" ]; then \
        find "$root" -xdev -type f ! -path /usr/local/bin/rsctf-entrypoint -exec grep -lF -e RSCTF_DYNAMIC_FLAG_fe0b899d99bcc94d24906 -e RSCTF_DYNAMIC_FLAG_5f4a92e94112bb2d {} + >> /etc/rsctf-flag-targets 2>/dev/null || true; \
      fi; \
    done; \
    sort -u /etc/rsctf-flag-targets > /etc/rsctf-flag-targets.sorted; \
    mv /etc/rsctf-flag-targets.sorted /etc/rsctf-flag-targets; \
    while IFS= read -r target; do [ -f "$target" ] && chmod a+w "$target"; done < /etc/rsctf-flag-targets
USER appuser
ENTRYPOINT ["/usr/local/bin/rsctf-entrypoint", "/usr/local/bin/server"]
