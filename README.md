# tf-fastpath

- English: [README.en.md](README.en.md)

Terraform / OpenTofu の inner loop を高速化するための graph-aware preview and gate workflow。

Graph-aware fast preview and authoritative gate workflow for Terraform and OpenTofu plans.

## Why

大規模な Terraform / OpenTofu 環境では、通常の `plan` が state refresh と依存関係の解決に時間を使い、ローカル変更のたびに数分待つことがあります。

`tf-fastpath` が狙うのは、authoritative な最終判定そのものを置き換えることではありません。開発者が毎回 full plan を待たなくても、変更範囲の見積もりと暫定差分を数秒で確認できるようにすることです。

## What tf-fastpath is

`tf-fastpath` は backend replacement ではなく、plan frontend / sidecar です。

- state の正本は、引き続き Terraform / OpenTofu backend または HCP Terraform 側にあります
- `apply` の実行経路は、引き続き通常の `terraform plan` / `terraform apply` を使います
- `tf-fastpath` が保持するのは派生データだけです
- 派生データの例は、依存グラフ、file-to-resource address の索引、前回 full plan の結果、state freshness metadata です
- 最初の実装では、これらの派生データを SQLite に保存する想定です

## What tf-fastpath optimizes

`tf-fastpath` は full plan を魔法のように短縮するツールではありません。最適化するのは開発者の inner loop です。

- 変更ファイルから影響を受けそうな resource address を推定する
- 依存 closure を使って blast radius を見積もる
- `terraform plan -json -refresh=false` を使って高速な preview を返す
- preview を通常の merge 判定に使い、重い full plan は不確実なケースに寄せる

## Command model

### `tf-fastpath index`

派生データを構築します。

- `terraform show -json` から state 情報を読む
- `terraform providers schema -json` から provider schema を読む
- HCL をパースして file / module / variable と resource address の対応を作る
- 必要に応じて `terraform graph` を補助的に使う
- 結果を SQLite に保存する

### `tf-fastpath preview`

高速な authoritative preview を返します。

- `git diff` から変更ファイルを取る
- 影響を受けそうな resource address 群を推定する
- blast radius を計算する
- `terraform plan -json -refresh=false` を実行して merge 判定に使う preview を返す

### `tf-fastpath verify`

preview の信頼度を支える検証を行います。

- `terraform plan -refresh-only` を定期実行して drift を記録する
- drift が多い workspace の preview confidence を下げる
- 必要なら opt-in で refresh-only apply 相当の同期運用を支援する

### `tf-fastpath gate`

必要な時だけ追加で実行する full plan です。

- 必要な時だけ通常の `terraform plan` を実行する
- preview と gate の差分を比較して精度を評価する
- preview confidence が十分なら、通常は preview を優先する

## Authoritative boundary

このプロジェクトの前提は明確です。

- authoritative なのは Terraform / OpenTofu の state です
- `preview` は通常の merge 判断を支える primary signal です
- `gate` は confidence が不足した時だけ使う escalation path です

## Non-goals

`tf-fastpath` は次のことを目的にしません。

- Terraform backend を置き換えること
- state の正本を別の保存形式へ移すこと
- どんな変更でも必ず full plan を待つこと
- routine use の `-target` に依存すること
- drift が激しい環境でも常に高精度だと保証すること

## Expected output

```text
Changed files: 4
Likely affected addresses: 17
Blast radius: 29 resources
State freshness: medium
Preview confidence: 0.82

Fast preview:
- 3 to change
- 1 to add
- 0 to destroy

Gate required before merge: yes
Reason: data source + IAM policy touched
```

## Roadmap

### v0.1

- `index`
- `preview`
- `gate`
- README と product framing

### v0.2

- GitHub Actions で PR comment を返す

### v0.3

- HCP Terraform 連携
- plan-only run との連携改善

### v0.4

- MCP server を追加し、AI / IDE から影響範囲を問い合わせられるようにする

## Current status

この repository は立ち上げ段階です。まずは problem statement、CLI skeleton、索引構築、preview / gate の最小実装から進めます。
