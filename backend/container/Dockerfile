# syntax=docker/dockerfile:1

# ============================================================ #
#  開発用イメージ（devcontainer 用）                            #
#  ※本番用イメージはデプロイフェーズで別途マルチステージで作成  #
# ============================================================ #
FROM golang:1.27-bookworm

ENV LANG=C.UTF-8 \
    TZ=Asia/Tokyo \
    CGO_ENABLED=0 \
    GOFLAGS=-buildvcs=false

WORKDIR /workspace

# 基本ツール + GitHub CLI（PR ベース運用のため gh を同梱）
RUN apt-get update && apt-get install -y --no-install-recommends \
      git curl ca-certificates openssh-client make \
    && update-ca-certificates \
    && curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
       | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
       > /etc/apt/sources.list.d/github-cli.list \
    && apt-get update && apt-get install -y --no-install-recommends gh \
    && rm -rf /var/lib/apt/lists/*

# 非 root の開発ユーザ（golang イメージは既定 root のため作成する）
ARG USERNAME=dev
ARG USER_UID=1000
ARG USER_GID=1000
RUN groupadd --gid ${USER_GID} ${USERNAME} \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m -s /bin/bash ${USERNAME} \
    && mkdir -p /workspace /go/pkg/mod \
    && chown -R ${USER_UID}:${USER_GID} /workspace /go

# ログインシェルでも Go の PATH を保証する
# （Debian の /etc/profile が PATH を上書きし /usr/local/go/bin・/go/bin が落ちるため）
RUN printf 'export GOPATH=/go\nexport PATH=/usr/local/go/bin:/go/bin:$PATH\n' \
      > /etc/profile.d/go.sh \
    && chmod 0644 /etc/profile.d/go.sh

USER ${USERNAME}

# Go 開発ツール（言語サーバ / デバッガ / 静的解析）
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest

CMD ["sleep", "infinity"]
