# syntax=docker/dockerfile:1
#
# 本番用イメージ。SPA を Go バイナリへ埋め込み、distroless 上で単体実行する。
# 開発用イメージ（ツールチェーン入り）は backend/container・frontend/container 側。
#
# ビルドコンテキストはリポジトリルート（frontend と backend の両方が要るため）:
#   docker build -f container/Dockerfile -t portfolio .

# ---- stage 1: SPA をビルドする ----
FROM oven/bun:1 AS web
WORKDIR /src/frontend

# 依存だけ先に入れてレイヤキャッシュを効かせる
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile

COPY frontend/ ./
RUN bun run build

# ---- stage 2: Go バイナリをビルドする ----
FROM golang:1.26-bookworm AS build
WORKDIR /src/backend

COPY backend/go.mod ./
RUN go mod download

COPY backend/ ./

# go:embed の対象へ SPA を配置してからビルドする。
# この COPY が build より前にある必要がある。
COPY --from=web /src/frontend/dist/ ./internal/assets/static/

# CGO を無効にして静的リンクする（distroless/static で動かすため）。
# -trimpath はビルドパスを、-s -w はデバッグ情報を落としてサイズを削る。
RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/portfolio ./cmd/server

# ---- stage 3: 実行イメージ ----
# シェルもパッケージマネージャも持たない。攻撃面が最小で、
# 侵入されても exec できるものが無い。
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/portfolio /portfolio

EXPOSE 8080
USER nonroot:nonroot

# exec 形式にすることでバイナリが PID 1 になり、
# SIGTERM が中間シェルを挟まず直接届く（graceful shutdown の前提）。
ENTRYPOINT ["/portfolio"]
