FROM ubuntu:18.04

WORKDIR /opt/app/

COPY premiumizearrd /opt/app/
COPY build/static /opt/app/static

ENTRYPOINT [ "/opt/app/premiumizearrd" ]