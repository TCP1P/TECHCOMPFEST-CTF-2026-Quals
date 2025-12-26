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
