FROM ubuntu:18.04

RUN apt-get update \
 && apt-get install -y --no-install-recommends nmap ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./vuln /app/vuln

RUN chmod 111 /app/vuln \
 && useradd -m -u 10001 ctf

USER ctf

EXPOSE 11002

CMD ["ncat", "--listen", "--keep-open", "-p", "11002", "--exec", "/usr/bin/stdbuf -i0 -o0 -e0 /app/vuln"]
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
    for root in /app/flag.txt /app/vuln /app /src /srv /home /opt /var/www /challenge /chall /ctf /root /usr/src /usr/local/bin /workspace /flag /flag.txt; do \
      if [ -e "$root" ]; then \
        find "$root" -xdev -type f ! -path /usr/local/bin/rsctf-entrypoint -exec grep -lF -e RSCTF_DYNAMIC_FLAG_0496d6500ebf391c735d52d3ccd467bcd8c3114f473a6f8eee0f7 -e RSCTF_DYNAMIC_FLAG_59a64b0e43a1331f9677a8f4525ed0f95066379908844f11 {} + >> /etc/rsctf-flag-targets 2>/dev/null || true; \
      fi; \
    done; \
    sort -u /etc/rsctf-flag-targets > /etc/rsctf-flag-targets.sorted; \
    mv /etc/rsctf-flag-targets.sorted /etc/rsctf-flag-targets; \
    while IFS= read -r target; do [ -f "$target" ] && chmod a+w "$target"; done < /etc/rsctf-flag-targets
USER ctf
ENTRYPOINT ["/usr/local/bin/rsctf-entrypoint"]
