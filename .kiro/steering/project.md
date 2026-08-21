# Project Overview

## Purpose

`terraform-provider-ranger` is a Terraform provider for
[Apache Ranger](https://ranger.apache.org/), the Hadoop ecosystem's
centralized security and authorization framework. It lets teams manage
Ranger access-control policies as declarative Terraform resources.

Published to the Terraform Registry under `gr8-toolkit/ranger`.
Protocol version: `5.0` (declared in `terraform-registry-manifest.json`).

## Module

```text
github.com/gr8-toolkit/terraform-provider-ranger
```

## Directory Structure

```text
.
├── main.go                          # Entry point — calls plugin.Serve
├── ranger/
│   ├── provider.go                  # Provider schema + resource registration
│   ├── client.go                    # HTTP client (resty) + providerConfigure
│   ├── resource_policy.go           # ranger_policy CRUD resource
│   ├── client_test.go               # Unit tests for providerConfigure
│   └── provider_test.go             # Unit tests for provider schema
├── tools/
│   └── tools.go                     # Build-tag tools pin (tfplugindocs)
├── examples/
│   ├── provider/provider.tf         # Provider usage example (for docs)
│   └── resources/ranger_policy/
│       ├── resource.tf              # Resource usage example (for docs)
│       └── import.sh                # Import example (for docs)
├── docs/                            # Generated — do not edit manually
│   ├── index.md
│   └── resources/policy.md
├── .kiro/steering/                  # Kiro steering files (this directory)
├── goreleaser.yml                   # Multi-platform release config
├── terraform-registry-manifest.json
├── Makefile
├── .pre-commit-config.yaml
└── .markdownlint.yaml
```

## Key Dependencies

| Package | Version | Role |
| --- | --- | --- |
| `hashicorp/terraform-plugin-sdk/v2` | v2.40.1 | Terraform provider SDK |
| `go-resty/resty/v2` | v2.17.2 | HTTP client for Ranger REST API |
| `hashicorp/terraform-plugin-docs` | v0.25.0 | Doc generation (`tfplugindocs`) |
| `stretchr/testify` | v1.12.0 | Test assertions |

## Resources

| Resource | Description |
| --- | --- |
| `ranger_policy` | Manages a Ranger access-control policy |

No data sources exist yet.

## Provider Configuration

```hcl
provider "ranger" {
  url             = "https://ranger.example.com"
  username        = "admin"
  password        = "secret"
  skip_tls_verify = false
}
```

## Makefile Targets

| Target | Command | Description |
| --- | --- | --- |
| `build` | `go build -o ./dist/` | Compile provider binary |
| `test` | `go test ./... -v` | Run all unit tests |
| `docs-gen` | `tfplugindocs generate && pre-commit run markdownlint --all-files` | Regenerate and lint docs |
