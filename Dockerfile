FROM node:23-alpine

ENTRYPOINT ["/usr/local/lib/tongli-release-monitor/bin/tongli-release-monitor", "--config", "/usr/local/etc/tongli-release-monitor/config.yml", "--data", "/var/local/tongli-release-monitor"]
VOLUME /var/local/tongli-release-monitor

# ARG TONGLI_UID=1000
# ARG TONGLI_GID=1000

# RUN addgroup -g ${TONGLI_GID} tongli && \
#     adduser -u ${TONGLI_UID} -D -h /usr/local/lib/tongli-release-monitor -s /bin/sh -G tongli tongli

WORKDIR /usr/local/lib/tongli-release-monitor

COPY bin ./bin

COPY package.json package-lock.json ./

RUN npm install --production

COPY lib ./lib

# USER tongli
