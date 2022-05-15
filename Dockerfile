FROM ubuntu:18.04

RUN apt update && \
    apt install openssl -y && \
    apt install ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/app/

COPY premiumizearrd /opt/app/
COPY build/static /opt/app/static

ENTRYPOINT [ "/opt/app/premiumizearrd" ]