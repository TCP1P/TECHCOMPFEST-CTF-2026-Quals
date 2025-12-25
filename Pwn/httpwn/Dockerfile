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
