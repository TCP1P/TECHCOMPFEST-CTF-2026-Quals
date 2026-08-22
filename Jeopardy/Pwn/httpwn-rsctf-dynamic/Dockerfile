FROM httpd:2.4

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      build-essential \
      vim \
      make \
      gcc \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -r ctf && useradd -r -g ctf ctf
RUN rm -rf /usr/local/apache2/htdocs/* /usr/local/apache2/cgi-bin/*
RUN sed -i 's/Listen 80/Listen 8000/' /usr/local/apache2/conf/httpd.conf

WORKDIR /home/ctf/chall

COPY src/ .
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
RUN make all install

COPY flag.txt /
RUN chown root:root /flag.txt && chmod 644 /flag.txt
RUN chown -R ctf:ctf /usr/local/apache2 /home/ctf/chall

USER root
CMD ["/entrypoint.sh"]

# RSCTF dynamic-flag compatibility layer. The build-time scan records
# exact marker-bearing files; runtime replacement is therefore bounded
# to those files instead of walking the image on every team start.
USER root
COPY rsctf-entrypoint.sh /usr/local/bin/rsctf-entrypoint
RUN set -eu; \
    chmod 0555 /usr/local/bin/rsctf-entrypoint; \
    : > /etc/rsctf-flag-targets; \
    for root in / /entrypoint.sh /home/ctf/chall /app /src /srv /home /opt /var/www /challenge /chall /ctf /root /usr/src /usr/local/bin /workspace /flag /flag.txt; do \
      if [ -e "$root" ]; then \
        find "$root" -xdev -type f ! -path /usr/local/bin/rsctf-entrypoint -exec grep -lF -e RSCTF_DYNAMIC_FLAG_6ae4ce5bccdba88be394ee1817f05a0fd655955 -e RSCTF_DYNAMIC_FLAG_10c7b989087c182b5646eb4d8ff1000776 {} + >> /etc/rsctf-flag-targets 2>/dev/null || true; \
      fi; \
    done; \
    sort -u /etc/rsctf-flag-targets > /etc/rsctf-flag-targets.sorted; \
    mv /etc/rsctf-flag-targets.sorted /etc/rsctf-flag-targets; \
    while IFS= read -r target; do [ -f "$target" ] && chmod a+w "$target"; done < /etc/rsctf-flag-targets
ENTRYPOINT ["/usr/local/bin/rsctf-entrypoint"]
