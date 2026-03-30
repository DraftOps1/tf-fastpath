# AGENTS.md

## Scope
このファイルは repository 全体に適用されます。
Codex の PR review では、このファイルを最優先の review policy として扱ってください。
現在の diff と repository context に基づいて review してください。

## Review Goal
強い human reviewer のように pull request を review してください。
style や polish よりも、実質的な engineering risk を優先してください。

特に次を重視します。
- bug と behavioral regression
- security と permission risk
- performance と scalability regression
- rollout と operational risk
- test 不足や verification の弱さ

体裁だけの指摘で review を埋めないでください。
現在の approach が unsafe または clearly unworkable でない限り、大きな rewrite は勧めないでください。

## Output Format
- review コメントは日本語で書く
- 重大な問題がある場合は `指摘事項` から始め、重要度順に並べる
- 各指摘では、何が問題か / なぜ重要か / どこで起きるか を具体的に書く
- 必要に応じてだけ短い節を使う
  - `指摘事項`
  - `リスク`
  - `追加確認`
- 重大な問題が無い場合は、`重大な指摘はありません。` を明示する
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
- state backend replacement に寄せない
- CLI surface を explicit かつ stable に保つ
- drift、stale state、blast radius、preview confidence を first-class concern として扱う
- preview を authoritative に見せる変更には懐疑的であること
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

PR が挙動を変えるのに test や verification note を更新していない場合は、それ自体を指摘してください。
