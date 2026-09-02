# GHCR を private で運用する

`ghcr.io/cyokozai/portfolio` を private のまま、Talos クラスタの `portfolio`
namespace から pull できるようにする手順。[要件書](./ci-requirements.md) の
「判断が必要な項目」#1 に対する回答。

## 分担

| 経路 | 認証 | 備考 |
|---|---|---|
| push（GitHub Actions → GHCR） | `GITHUB_TOKEN` の `packages: write` | **追加の Secret は不要**。private でも変わらない |
| pull（kubelet → GHCR） | `portfolio/ghcr-pull`（`kubernetes.io/dockerconfigjson`） | 帯域外で登録する。リポジトリには置かない |

Argo CD はレジストリに接触しない（Git を読むだけ）ので、Argo CD 側の設定は要らない。

## 手順

### 1. パッケージを private にする

**GHCR のパッケージは初回 push 時に private で作られる。** リポジトリが public でも
パッケージの可視性は独立しているため、意図せず public になることはない。

初回の CD が走ったあと、`https://github.com/users/cyokozai/packages/container/portfolio/settings`
で `Private` になっていることを確認する。

### 2. pull 用の PAT を発行する

**classic PAT で `read:packages` だけ**を付ける。`write:packages` は付けない
（push は Actions の `GITHUB_TOKEN` が行うため不要）。

有効期限を付けた場合は、失効前に手順 3 を再実行する。

### 3. Secret をクラスタへ登録する

`portfolio` namespace が存在してから実行する（namespace は `deploy/namespace.yaml`
を Argo CD が同期して作る）。

```bash
export KUBECONFIG=/path/to/kubeconfig
export GHCR_TOKEN=<手順 2 の PAT>
./scripts/ghcr-pull-secret.sh
```

再実行は上書きになる。PAT のローテーションはこれで足りる。

### 4. 確認する

```bash
kubectl -n portfolio get secret ghcr-pull
kubectl -n portfolio rollout restart deploy/portfolio
kubectl -n portfolio get pod -w          # ImagePullBackOff が出ないこと
kubectl -n portfolio get pod -o jsonpath='{..image}'
```

## 落とし穴

- **namespace ごと消すと Secret も消える。** `deploy/namespace.yaml` は Argo CD の
  同期対象で `prune: true` が効いている。namespace が作り直されると、この Secret は
  Git に無いので復活しない。手順 3 を再実行する
- **初回 bootstrap の順序。** Application を登録した直後は Secret がまだ無いため
  `ImagePullBackOff` になる。手順 3 を実行すれば kubelet が自動で再試行して復帰する
- **PAT の失効は静かに効く。** タグは commit SHA なので、既に pull 済みのイメージで
  動いている Pod は落ちない。**次のデプロイで初めて失敗する**
- **Secret を `deploy/kustomization.yaml` の `resources` に追加しないこと。**
  追加すると認証情報が Git に載る

## public にする選択肢

このリポジトリは public であり、イメージの中身は公開済みのソースを静的リンクした
バイナリと SPA だけで、秘匿すべきものを含まない。**private にして得られるものは小さく、
代わりに長期の PAT を 1 本維持する必要が生じる**（失効すればデプロイが止まる）。

public にするなら、手順 1 で `Public` に変更し、`deploy/deployment.yaml` の
`imagePullSecrets` を外すだけでよい。この判断は後からどちらにも変えられる。
