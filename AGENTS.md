# AGENTS.md

## Scope
このファイルは repository 全体に適用されます。
Codex の PR review では、このファイルを最優先の review policy として扱ってください。

## Review Policy
- review コメントは日本語で書く
- 重大な問題がある場合は、重要度順に `指摘事項` から書き始める
- 各指摘では、何が問題か / なぜ重要か / どこで起きるか を具体的に書く
- 重要な問題が無い場合は、`重大な指摘はありません。` を明示する
- 必要に応じてだけ短い節を使う
  - `指摘事項`
  - `リスク`
  - `追加確認`
- 体裁の指摘だけで review を埋めない。挙動、セキュリティ、運用、性能、検証不足を優先する
- `.agents/skills/review-ja` が存在する場合は、その skill を使う

## Repository Context
この project は Terraform / OpenTofu の fast preview workflow 向け OSS CLI です。
現在の v0.1 scope は次です。
- index
- preview
- gate

authoritative なのは通常の Terraform / OpenTofu plan です。
preview は advisory であり、merge / apply の最終判定ではありません。

## What To Protect
- sidecar model を維持する
- backend replacement に寄せない
- preview を authoritative に見せない
- drift、stale state、blast radius、preview confidence を軽視しない
- routine use の Terraform `-target` を前提にしない
- 強い理由なしに依存を増やさない

## Verification
コード変更がある PR では、まず次を優先してください。
```bash
go test ./...
```

CLI の挙動を触る PR では、help 出力も確認対象にしてください。
```bash
go run ./cmd/tf-fastpath --help
go run ./cmd/tf-fastpath preview --help
```

テストや verification note が足りない変更は、それ自体を指摘対象にしてください。
