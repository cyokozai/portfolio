# syntax=docker/dockerfile:1
FROM node:24-slim

ENV LANG=C.UTF-8 \
    TZ=Asia/Tokyo

WORKDIR /workspace

RUN apt-get update && apt-get install -y --no-install-recommends \
      git curl ca-certificates openssh-client \
    && update-ca-certificates \
    && curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
       | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
       > /etc/apt/sources.list.d/github-cli.list \
    && apt-get update && apt-get install -y --no-install-recommends gh \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /workspace/node_modules && chown -R node:node /workspace

USER node

CMD ["sleep", "infinity"]
