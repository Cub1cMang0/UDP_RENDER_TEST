FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates jq \
    && rm -rf /var/lib/apt/lists/*

# Version confirmed working in the local Docker Compose persistence test.
ARG QURL_VERSION=1.5.0
RUN curl -fsSL "https://github.com/layervai/qurl-integrations/releases/download/v${QURL_VERSION}/qurl_${QURL_VERSION}_linux_amd64.tar.gz" -o /tmp/qurl.tar.gz \
    && tar -xzf /tmp/qurl.tar.gz -C /usr/local/bin qurl \
    && rm /tmp/qurl.tar.gz \
    && chmod +x /usr/local/bin/qurl

WORKDIR /app
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
