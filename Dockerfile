FROM ubuntu:latest

RUN apt update && \
    apt install openssl -y && \
    apt install ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/app/

COPY premiumizearrd /opt/app/
COPY build/static /opt/app/static

EXPOSE 8182

ENTRYPOINT [ "/opt/app/premiumizearrd" ]