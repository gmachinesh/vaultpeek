# vaultpeek

A CLI tool to inspect and diff HashiCorp Vault secrets across environments.

---

## Installation

```bash
go install github.com/youruser/vaultpeek@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/vaultpeek.git && cd vaultpeek && go build -o vaultpeek .
```

---

## Usage

Set your Vault address and token, then start inspecting secrets:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your_token_here"

# Inspect a secret path
vaultpeek inspect secret/data/myapp/config

# Diff secrets between two environments
vaultpeek diff secret/data/myapp/staging secret/data/myapp/production
```

### Example Output

```
[~] DB_HOST       staging.db.internal  →  prod.db.internal
[+] CACHE_TTL     (missing)            →  3600
[-] DEBUG_MODE    true                 →  (missing)
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: `text`, `json`, `yaml` |
| `--mask` | Mask secret values in output |
| `--token` | Vault token (overrides `VAULT_TOKEN`) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault v1.x

---

## License

MIT © 2024 youruser