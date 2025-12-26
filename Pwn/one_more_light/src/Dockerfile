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
