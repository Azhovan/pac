# PAC - Portable AWS Credentials

Transfer AWS SSO session credentials between machines. Useful when you need valid AWS credentials on a machine that doesn't have SSO configured (e.g., CI runners, dev VMs, or remote workstations).

## How It Works

1. **Export** on your SSO-authenticated machine — resolves temporary credentials and saves them to a JSON file
2. **Transfer** the file to the target machine (scp, USB, etc.)
3. **Import** on the target machine — writes the credentials into `~/.aws/credentials` and region into `~/.aws/config`

## Install

```bash
make build
```

Produces a `pac` binary. Copy it to both machines.

## Usage

### Export (source machine)

```bash
# Ensure you're logged in
aws sso login --profile my-profile

# Export credentials
pac export --profile my-profile --output aws-creds.json
```

### Import (target machine)

```bash
# Import under the same profile name
pac import --file aws-creds.json

# Or override the profile name
pac import --file aws-creds.json --profile custom-name
```

### Flags

| Command  | Flag               | Description                          |
|----------|--------------------|--------------------------------------|
| `export` | `-p, --profile`    | AWS SSO profile name (required)      |
| `export` | `-o, --output`     | Output file path (default `aws-creds.json`) |
| `import` | `-f, --file`       | Path to credentials JSON file (required) |
| `import` | `-p, --profile`    | Override target profile name         |

## Credential File Format

The exported JSON contains:

```json
{
  "access_key_id": "ASIA...",
  "secret_access_key": "...",
  "session_token": "...",
  "expiration": "2026-04-23T18:00:00Z",
  "region": "eu-central-1",
  "profile_name": "my-profile"
}
```

Treat this file as a secret. It is written with `0600` permissions.

## Cross-Compile

```bash
# Linux target
GOOS=linux GOARCH=amd64 go build -o pac-linux .

# macOS ARM
GOOS=darwin GOARCH=arm64 go build -o pac-darwin .
```
