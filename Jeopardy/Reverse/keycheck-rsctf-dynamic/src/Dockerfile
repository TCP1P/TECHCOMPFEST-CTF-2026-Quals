FROM ubuntu:18.04

RUN apt-get update \
 && apt-get install -y --no-install-recommends nmap ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY keycheck /app/keycheck

RUN chmod 555 /app/keycheck \
 && useradd -m -u 10001 ctf

USER ctf

EXPOSE 11001

CMD ["ncat", "--listen", "--keep-open", "-p", "11001", "--exec", "/usr/bin/stdbuf -i0 -o0 -e0 /app/keycheck"]
# Bake the legacy Compose flag mount into RSCTF's single-image runtime.
COPY ["flag.txt", "/app/flag.txt"]

# RSCTF dynamic-flag compatibility layer. The build-time scan records
# exact marker-bearing files; runtime replacement is therefore bounded
# to those files instead of walking the image on every team start.
USER root
COPY rsctf-entrypoint.sh /usr/local/bin/rsctf-entrypoint
RUN set -eu; \
    chmod 0555 /usr/local/bin/rsctf-entrypoint; \
    : > /etc/rsctf-flag-targets; \
    for root in /app/flag.txt /app/keycheck /app /src /srv /home /opt /var/www /challenge /chall /ctf /root /usr/src /usr/local/bin /workspace /flag /flag.txt; do \
      if [ -e "$root" ]; then \
        find "$root" -xdev -type f ! -path /usr/local/bin/rsctf-entrypoint -exec grep -lF -e RSCTF_DYNAMIC_FLAG_92f24d0c0840158155e4d04fcd96a44967de0 -e RSCTF_DYNAMIC_FLAG_acf4fcaa1bf77c30504da402eab6e7f7 {} + >> /etc/rsctf-flag-targets 2>/dev/null || true; \
      fi; \
    done; \
    sort -u /etc/rsctf-flag-targets > /etc/rsctf-flag-targets.sorted; \
    mv /etc/rsctf-flag-targets.sorted /etc/rsctf-flag-targets; \
    while IFS= read -r target; do [ -f "$target" ] && chmod a+w "$target"; done < /etc/rsctf-flag-targets
USER ctf
ENTRYPOINT ["/usr/local/bin/rsctf-entrypoint"]
