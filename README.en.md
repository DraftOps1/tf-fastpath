# tf-fastpath

- 日本語: [README.md](README.md)

Graph-aware fast preview and authoritative gate workflow for Terraform and OpenTofu plans.

## Why

Large Terraform and OpenTofu environments often spend most of their time in state refresh and dependency resolution before a normal `plan` can show a useful diff.

`tf-fastpath` does not try to replace the authoritative full plan. It tries to shorten the developer inner loop so contributors can inspect likely impact and a provisional diff in seconds instead of waiting for a full refresh on every edit.

## What tf-fastpath is

`tf-fastpath` is a plan frontend / sidecar, not a backend replacement.

- Terraform or OpenTofu state remains authoritative in the existing backend or HCP Terraform
- `terraform plan` / `terraform apply` remain the execution path for authoritative changes
- `tf-fastpath` stores derived data only
- Derived data includes dependency graphs, file-to-resource address indexes, previous full-plan results, and state freshness metadata
- The first implementation stores those derived artifacts in SQLite

## What tf-fastpath optimizes

`tf-fastpath` does not magically make the authoritative full plan cheap. It optimizes the inner loop around that full plan.

- infer likely affected resource addresses from changed files
- compute dependency closure and blast radius
- run `terraform plan -json -refresh=false` for a fast provisional preview
- move heavyweight full plans to pre-merge, pre-apply, and scheduled verification

## Command model

### `tf-fastpath index`

Build the derived data set.

- read state with `terraform show -json`
- read provider schema with `terraform providers schema -json`
- parse HCL and map files, modules, and variables to resource addresses
- optionally use `terraform graph` as supporting input
- persist results to SQLite

### `tf-fastpath preview`

Return a fast provisional preview.

- collect changed files from `git diff`
- infer likely affected resource addresses
- compute blast radius
- run `terraform plan -json -refresh=false` and summarize the preview

### `tf-fastpath verify`

Run validation that supports preview confidence.

- execute `terraform plan -refresh-only` on a schedule and record drift
- lower preview confidence when a workspace drifts frequently
- optionally support refresh-only synchronization workflows

### `tf-fastpath gate`

Run the authoritative check before merge or apply.

- always execute the normal `terraform plan`
- compare preview output against gate output
- treat gate as the final authority for merge and apply decisions

## Authoritative boundary

The project is explicit about its boundary.

- Terraform or OpenTofu state plus a normal full plan remain authoritative
- `preview` is fast but provisional
- `gate` is the only command that should drive merge or apply decisions

## Non-goals

`tf-fastpath` does not aim to:

- replace the Terraform backend
- move the source-of-truth state into a new storage model
- present preview as equivalent to an authoritative plan
- rely on routine `-target` usage
- guarantee high preview accuracy in high-drift environments

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
- README and product framing

### v0.2

- GitHub Actions PR comments

### v0.3

- HCP Terraform integration
- better plan-only run integration

### v0.4

- MCP server for AI and IDE impact questions

## Current status

The repository is in bootstrap mode. The next steps are the CLI skeleton, indexing pipeline, and the first preview / gate implementation.
