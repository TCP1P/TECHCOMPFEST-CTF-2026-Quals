FROM python:3.13-alpine

RUN apk update && apk add socat gcc g++ make libffi-dev openssl-dev python3-dev exiftool

RUN adduser -D -u 1001 -s /bin/bash ctf

RUN mkdir /home/ctf/chall

COPY ./requirements.txt /home/ctf/chall
RUN pip3 install -r /home/ctf/chall/requirements.txt

RUN mkdir /home/ctf/chall/src

COPY ./chall.py /home/ctf/chall/src
COPY ./run.sh /home/ctf/chall/src
COPY ./flag.png /home/ctf/chall/src
COPY ./images.jpg /home/ctf/chall/src

RUN chown -R root:root /home/ctf/chall
RUN chmod -R 555 /home/ctf/chall
USER ctf
WORKDIR /home/ctf/chall/src

CMD ["./run.sh"]

# RSCTF dynamic-flag compatibility layer. The build-time scan records
# exact marker-bearing files; runtime replacement is therefore bounded
# to those files instead of walking the image on every team start.
USER root
COPY rsctf-entrypoint.sh /usr/local/bin/rsctf-entrypoint
COPY rsctf-prepare-flag.sh /usr/local/bin/rsctf-prepare-flag
RUN set -eu; \
    chmod 0555 /usr/local/bin/rsctf-entrypoint; \
    : > /etc/rsctf-flag-targets; \
    for root in /home/ctf/chall /home/ctf/chall/src /app /src /srv /home /opt /var/www /challenge /chall /ctf /root /usr/src /usr/local/bin /workspace /flag /flag.txt; do \
      if [ -e "$root" ]; then \
        find "$root" -xdev -type f ! -path /usr/local/bin/rsctf-entrypoint -exec grep -lF -e RSCTF_DYNAMIC_FLAG_8c22b404b39926293945fc9ed72cc0c1d0 -e RSCTF_DYNAMIC_FLAG_c5588a7ef51f112a28f2c96511d61 {} + >> /etc/rsctf-flag-targets 2>/dev/null || true; \
      fi; \
    done; \
    sort -u /etc/rsctf-flag-targets > /etc/rsctf-flag-targets.sorted; \
    mv /etc/rsctf-flag-targets.sorted /etc/rsctf-flag-targets; \
    while IFS= read -r target; do [ -f "$target" ] && chmod a+w "$target"; done < /etc/rsctf-flag-targets; \
    chmod 0555 /usr/local/bin/rsctf-prepare-flag; \
    chmod a+w /home/ctf/chall/src/flag.png
USER ctf
ENTRYPOINT ["/usr/local/bin/rsctf-entrypoint"]
