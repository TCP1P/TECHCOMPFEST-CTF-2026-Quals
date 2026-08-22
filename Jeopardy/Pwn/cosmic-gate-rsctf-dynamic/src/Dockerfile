FROM ubuntu@sha256:12aa2ae97b9b295d28962e83cb36cb478f3afbdda68196e5fd62e68bde3a67d5

RUN apt-get update && \
    apt-get install -y socat && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY flag.txt cosmic_gate /app/

RUN chmod +x /app/cosmic_gate
RUN chmod og+r /app/flag.txt
RUN chmod og-w /app/*

EXPOSE 2201

RUN useradd -m -s /usr/sbin/nologin ctf
USER ctf
# Original final-stage ENTRYPOINT is wrapped by RSCTF at runtime.

# RSCTF dynamic-flag compatibility layer. The build-time scan records
# exact marker-bearing files; runtime replacement is therefore bounded
# to those files instead of walking the image on every team start.
USER root
COPY rsctf-entrypoint.sh /usr/local/bin/rsctf-entrypoint
RUN set -eu; \
    chmod 0555 /usr/local/bin/rsctf-entrypoint; \
    : > /etc/rsctf-flag-targets; \
    for root in /app /src /srv /home /opt /var/www /challenge /chall /ctf /root /usr/src /usr/local/bin /workspace /flag /flag.txt; do \
      if [ -e "$root" ]; then \
        find "$root" -xdev -type f ! -path /usr/local/bin/rsctf-entrypoint -exec grep -lF -e RSCTF_DYNAMIC_FLAG_388af0ad05229a0ee367b436df935af1010c292e6d5ffec10c382f822a2273e8388af0ad05229a0ee367b436df935af1010c2 -e RSCTF_DYNAMIC_FLAG_bcb0ea2fc5c78db462ea6ded53af3e8770f8d5993d135f8611a20490f2ede7febcb0ea2fc5c78db462ea6ded53af3e87 {} + >> /etc/rsctf-flag-targets 2>/dev/null || true; \
      fi; \
    done; \
    sort -u /etc/rsctf-flag-targets > /etc/rsctf-flag-targets.sorted; \
    mv /etc/rsctf-flag-targets.sorted /etc/rsctf-flag-targets; \
    while IFS= read -r target; do [ -f "$target" ] && chmod a+w "$target"; done < /etc/rsctf-flag-targets
USER ctf
ENTRYPOINT ["/usr/local/bin/rsctf-entrypoint", "socat", "-d", "-d", "TCP-LISTEN:2201,reuseaddr,fork", "EXEC:'/usr/bin/timeout 600 stdbuf -i0 -o0 -e0 /app/cosmic_gate',pty,rawer,echo=0,ctty,setsid,stderr"]
