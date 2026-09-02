#!/usr/bin/env bash
#
# GHCR のパッケージを private にしたまま、portfolio namespace から
# イメージを pull できるようにする。
#
# この Secret はリポジトリに置かない。cloudflared の認証情報と同じ扱いで、
# 帯域外（この script）でクラスタへ直接登録する。Argo CD の同期対象にも入れない
# （deploy/kustomization.yaml の resources に載せないこと）。
#
# 使い方:
#   export KUBECONFIG=/path/to/kubeconfig
#   export GHCR_TOKEN=<read:packages だけを持つ classic PAT>
#   ./scripts/ghcr-pull-secret.sh
#
# 再実行すると Secret を上書きする。PAT のローテーションはこれで足りる。
#
# PAT は argv に載せない（host の process 一覧に出るため）。
# dockerconfigjson を stdin 経由で kubectl へ渡している。
set -euo pipefail

GHCR_USER="${GHCR_USER:-cyokozai}"
NAMESPACE="${NAMESPACE:-portfolio}"
SECRET_NAME="${SECRET_NAME:-ghcr-pull}"
# クラスタのバージョンから大きく離れたものは使わない。
KUBECTL_IMAGE="${KUBECTL_IMAGE:-registry.k8s.io/kubectl:v1.34.1}"

if [ -z "${GHCR_TOKEN:-}" ]; then
  cat >&2 <<'MSG'
GHCR_TOKEN が未設定です。

read:packages だけを持つ classic PAT を発行し、環境変数で渡してください。
write:packages は不要です（push は GitHub Actions の GITHUB_TOKEN が行います）。

  export GHCR_TOKEN=<PAT>
MSG
  exit 1
fi

if [ -z "${KUBECONFIG:-}" ]; then
  cat >&2 <<'MSG'
KUBECONFIG が未設定です。クラスタの kubeconfig を指してください。

  export KUBECONFIG=~/myproject/thin-k8s/_out/kubeconfig-tailscale
MSG
  exit 1
fi

# KUBECONFIG は複数パスを : で連結できるが、bind mount は 1 ファイルしか渡せない。
case "$KUBECONFIG" in
  *:*)
    echo "KUBECONFIG に複数のパスが指定されています。この script では 1 ファイルだけを指してください: $KUBECONFIG" >&2
    exit 1
    ;;
esac

# 相対パスは bind mount で named volume と解釈されるため絶対パスへ直す。
# ~ はシェルが展開しないケース（クオート内など）があるので明示的に開く。
# 環境変数から来た値なので ~ は未展開のリテラル。展開させず照合するのが正しい。
# shellcheck disable=SC2088
case "$KUBECONFIG" in
  "~/"*) KUBECONFIG="${HOME}/${KUBECONFIG#\~/}" ;;
esac
if [ "${KUBECONFIG#/}" = "$KUBECONFIG" ]; then
  KUBECONFIG="$(pwd)/$KUBECONFIG"
fi

# ここを検証せずに docker run へ渡すと、**存在しないパスを Docker が
# 空ディレクトリとして自動作成**し、`/kubeconfig: is a directory` で落ちる。
# しかもホスト側にはその空ディレクトリが residue として残る。
if [ ! -e "$KUBECONFIG" ]; then
  echo "KUBECONFIG のパスが存在しません: $KUBECONFIG" >&2
  exit 1
fi
if [ -d "$KUBECONFIG" ]; then
  cat >&2 <<MSG
KUBECONFIG がディレクトリを指しています: $KUBECONFIG

以前に存在しないパスを指したまま docker run したことで、Docker が
空ディレクトリを作った可能性があります。空であれば削除してください。

  rmdir "$KUBECONFIG"
MSG
  exit 1
fi
if [ ! -f "$KUBECONFIG" ]; then
  echo "KUBECONFIG が通常ファイルではありません: $KUBECONFIG" >&2
  exit 1
fi
if [ ! -r "$KUBECONFIG" ]; then
  echo "KUBECONFIG を読めません: $KUBECONFIG" >&2
  exit 1
fi

echo "registry=ghcr.io user=${GHCR_USER} namespace=${NAMESPACE} secret=${SECRET_NAME}"

# 1 段目は --dry-run=client でサーバへ接触しないため kubeconfig を渡さない。
# 2 段目へは host のパイプで渡すので、PAT が中間ファイルに落ちない。
printf '{"auths":{"ghcr.io":{"auth":"%s"}}}' \
  "$(printf '%s:%s' "$GHCR_USER" "$GHCR_TOKEN" | base64 | tr -d '\n')" \
| docker run --rm -i "$KUBECTL_IMAGE" \
    create secret generic "$SECRET_NAME" \
      --namespace "$NAMESPACE" \
      --type=kubernetes.io/dockerconfigjson \
      --from-file=.dockerconfigjson=/dev/stdin \
      --dry-run=client -o yaml \
| if [ -n "${DRY_RUN:-}" ]; then
    # 生成した Secret の構造だけを見る。クラスタへは書き込まない。
    # 動作確認は必ずこちらで行うこと（誤って実在の PAT を上書きしないため）。
    sed 's/^\( *\.dockerconfigjson: \).*/\1<省略>/'
  else
    docker run --rm -i \
      -v "${KUBECONFIG}:/kubeconfig:ro" -e KUBECONFIG=/kubeconfig \
      "$KUBECTL_IMAGE" apply -f -
  fi

if [ -n "${DRY_RUN:-}" ]; then
  echo "DRY_RUN のため書き込んでいません。"
else
  echo "登録しました。deploy/deployment.yaml の imagePullSecrets がこの名前を参照しています。"
fi
