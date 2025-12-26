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
ENTRYPOINT ["socat","-d","-d","TCP-LISTEN:2201,reuseaddr,fork","EXEC:'/usr/bin/timeout 600 stdbuf -i0 -o0 -e0 /app/cosmic_gate',pty,rawer,echo=0,ctty,setsid,stderr"]
