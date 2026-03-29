# AGENTS.md

## Scope
This file applies to the entire repository.
Treat it as the primary review policy for native Codex GitHub review.
Keep the review grounded in the current diff and repository context.

## Review Goal
Review pull requests the way a strong human reviewer would.
Prioritize material engineering risk over style or polish.

Focus on:
- bugs and behavioral regressions
- security and permission risks
- performance and scalability regressions
- rollout and operational risks
- missing tests or weak verification

Do not spend review effort on style nits unless they hide a real maintenance problem.
Do not recommend broad rewrites unless the current approach is unsafe or clearly unworkable.

## Output Format
When there are material issues, start with findings.
Order findings by severity.
Keep each finding concrete:
- what is wrong
- why it matters
- where it appears

Use short sections only when needed:
- Findings
- Risks
- Follow-ups

If there are no material findings, say so explicitly.

## Repository Context
This project is an OSS CLI for Terraform and OpenTofu fast preview workflows.
The current product scope is v0.1:
- index
- preview
- gate

The authoritative result remains the normal Terraform or OpenTofu plan before merge or apply.
Any preview path is advisory, not authoritative.

## What To Protect
- Preserve the sidecar model. Do not turn this project into a state backend replacement.
- Keep the CLI surface explicit and stable.
- Treat drift, stale state, blast radius, and preview confidence as first-class concerns.
- Be skeptical of any change that makes preview results look authoritative.
- Be skeptical of routine Terraform `-target` usage.
- Avoid adding dependencies without strong justification.

## Verification
When code changes are present, prefer these checks:
```bash
go test ./...
```

If the change affects CLI behavior, also confirm help output still works:
```bash
go run ./cmd/tf-fastpath --help
go run ./cmd/tf-fastpath preview --help
```

If a PR changes behavior without updating tests or verification notes, call that out.