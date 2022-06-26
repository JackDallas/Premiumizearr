LABEL build_version="Premiumizearr version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="JackDallas"

COPY docker/root/ /

EXPOSE 8182

RUN mkdir /data
RUN mkdir /unzip
RUN mkdir /downloads
RUN mkdir /transfers
RUN mkdir /blackhole
RUN mkdir -p /opt/app/

WORKDIR /opt/app/

ENV PREMIUMIZEARR_CONFIG_DIR_PATH=/data
ENV PREMIUMIZEARR_LOGGING_DIR_PATH=/data

COPY premiumizearrd /opt/app/
COPY build/static /opt/app/static

ENTRYPOINT ["/init"]