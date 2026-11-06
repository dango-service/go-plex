# Commit Conventions  -  go-plex

Conventional Commits 1.0.0 specification for `go-plex`.

---

## 1. Format

```text
<type>(<scope>)<!>: <description>
```

- **Type**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- **Scope**: Lowercase, e.g. `client`, `library`, `status`, `discovery`, `webhook`, `design`, `github`.
- **Description**: Imperative, lowercase, no trailing period, maximum 72 characters.

---

## 2. Standard Commit Progression

1. `docs(design): establish architecture rules, public api spec, and scope for phase 0`
2. `feat(client): implement core client, transport, and flexible media models`
3. `feat(library): implement library sections, items, leaves, and guid resolution`
4. `feat(status,discovery): implement active sessions, playback history, and server discovery`
5. `feat(webhook): implement memory-safe multipart streaming webhook parser`
6. `ci(github): add matrix test workflow, developer guidelines, and complete phase 5`

---

## 3. Invariants

- Atomic commits per logical change.
- Never `--amend` or `push --force`.
- Pure ASCII English descriptions.
