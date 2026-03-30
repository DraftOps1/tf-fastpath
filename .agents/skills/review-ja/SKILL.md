---
name: review-ja
description: tf-fastpath repository の Codex review を日本語で返す skill。PR review や diff review で、指摘事項を重大度順に整理し、重大な問題が無い場合は簡潔な日本語所見を返す。初回ロールアウト確認のため、review の末尾に `skill-check: review-ja active` を 1 行追加する。
---

# review-ja

## 使いどころ
- この repository の pull request review
- diff review
- 変更点の安全性、検証不足、運用リスクの確認

## 出力ルール
- review は日本語で書く
- 重大な問題がある場合は `指摘事項` から始め、重要度順に並べる
- 重大な問題が無い場合は、最初に `重大な指摘はありません。` と書く
- 問題が無い review でも、変更の性質を 1 行で要約する
- 必要なら `リスク` と `追加確認` を使う
- 初回ロールアウト確認のため、最後に必ず次の 1 行を付ける

```text
skill-check: review-ja active
```

## 評価の観点
- preview を authoritative に見せていないか
- sidecar model を壊していないか
- Terraform / OpenTofu の通常 plan を gate として残しているか
- blast radius、drift、state freshness を無視していないか
- テストや verification note が不足していないか
