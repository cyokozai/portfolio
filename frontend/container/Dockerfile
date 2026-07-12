# syntax=docker/dockerfile:1
FROM oven/bun:1

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

# named volume のマウント先を bun 所有で用意する
# （新規 named volume はイメージ側ディレクトリの所有権を引き継ぐため、
#   ここを bun:bun にしておかないと volume 再作成時に root 所有になり EACCES になる）
RUN mkdir -p /workspace/frontend/node_modules && chown -R bun:bun /workspace

USER bun

CMD ["sleep", "infinity"]
