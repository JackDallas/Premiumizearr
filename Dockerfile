FROM ubuntu:latest

RUN apt update && \
    apt install openssl -y && \
    apt install ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /data
RUN mkdir /unzip
RUN mkdir /downloads
RUN mkdir /transfers
RUN mkdir /blackhole

ENV PREMIUMIZEARR_CONFIG_DIR_PATH=/data
ENV PREMIUMIZEARR_LOGGING_DIR_PATH=/data

EXPOSE 8182

WORKDIR /opt/app/

COPY premiumizearrd /opt/app/
COPY build/static /opt/app/static

ENTRYPOINT [ "/opt/app/premiumizearrd" ]