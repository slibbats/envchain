# envchain

> Manages layered environment variable overrides across dev/staging/prod without leaking secrets.

---

## Installation

```bash
go install github.com/yourname/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envchain.git && cd envchain && go build ./...
```

---

## Usage

Define layered environment files per environment:

```
.env.base
.env.dev
.env.staging
.env.prod
```

Run your command with the resolved environment chain:

```bash
envchain --env prod -- ./myapp serve
```

Variables are merged in order — base first, then the target environment layer on top. Secrets are never written to stdout or logs.

**Example config (`.envchain.yaml`):**

```yaml
layers:
  - .env.base
  - .env.${ENV}
secrets:
  - DATABASE_URL
  - API_SECRET_KEY
```

Marked secrets are masked in all output and never exported to child process environments unless explicitly allowed.

```bash
# Override a single variable at runtime
envchain --env staging --set PORT=9090 -- ./myapp serve
```

---

## Why envchain?

- **Layered overrides** — base + environment-specific values, no duplication
- **Secret isolation** — sensitive keys stay out of logs and shell history
- **Zero runtime deps** — single static binary, drop it anywhere

---

## License

MIT © 2024 yourname