# My portfolio

フロントエンドはSvelte、バックエンドはGoを使って実装しています。

デプロイ先はオンプレミスの Talos Linux 単一ノード Kubernetes クラスタです。クラスタ自体は別リポジトリ `thin-k8s` が管理しており、このリポジトリはアプリ本体（`deploy/`）と Argo CD の `Application` だけを持ちます。

外部公開は Cloudflare Tunnel → Traefik 経由で、ドメインは `cyokozai.net` です。
