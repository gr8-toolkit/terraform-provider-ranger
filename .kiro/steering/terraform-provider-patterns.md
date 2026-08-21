# Terraform Provider Patterns

## SDK Version

This provider uses **Terraform Plugin SDK v2**
(`hashicorp/terraform-plugin-sdk/v2`). Do not mix in Plugin Framework
(`hashicorp/terraform-plugin-framework`) types or patterns.

## Provider Schema

Defined in `ranger/provider.go`. The four provider-level attributes are:

| Attribute | Type | Required | Sensitive |
| --- | --- | --- | --- |
| `url` | String | Yes | No |
| `username` | String | Yes | Yes |
| `password` | String | Yes | Yes |
| `skip_tls_verify` | Bool | No (default `false`) | No |

When adding new provider-level config (e.g. timeout, CA cert), follow the same
pattern: declare in `Provider()` schema, read in `providerConfigure`, pass into
`Client`.

## Client Initialization (`ranger/client.go`)

`providerConfigure` constructs the `Client` and validates connectivity by
calling `GET /service/public/v2/api/policy`. This connectivity check runs on
every `terraform init` / plan / apply. Keep it lightweight — it must not have
side effects.

`newClient()` sets:

- Base URL from `url`
- Basic auth from `username` / `password`
- `Content-Type: application/json`
- Optional TLS skip via `InsecureSkipVerify` when `skip_tls_verify = true`

New API methods belong on `Client` as methods, not as free functions.

## Resource Structure (`ranger/resource_policy.go`)

Every resource follows this layout:

```go
func resourcePolicy() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourcePolicyCreate,
        ReadContext:   resourcePolicyRead,
        UpdateContext: resourcePolicyUpdate,
        DeleteContext: resourcePolicyDelete,
        Importer: &schema.ResourceImporter{
            StateContext: resourcePolicyImport,
        },
        Schema: map[string]*schema.Schema{ ... },
    }
}
```

All CRUD functions have the signature:

```go
func resourcePolicyXxx(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics
```

Cast `meta` to `*ranger.Client` at the top of each function.

## State Management Rules

- Always call `d.SetId(id)` in Create after the API returns the new resource ID.
- In Delete, call `d.SetId("")` on success (including 404 — already deleted).
- Read must set **all** schema attributes from the API response or Terraform will
  show phantom diffs.
- Strip server-managed fields before storing JSON in state — see
  `parsePolicyResponse`. Never store fields that Ranger auto-populates (`id`,
  `guid`, `createdBy`, `updatedBy`, `createTime`, `updateTime`, `version`,
  `resourceSignature`).

## The `definition` Field Pattern

`ranger_policy.definition` is an opaque JSON string representing the full
Ranger policy body. The provider does not parse or validate its internal
structure — callers are responsible for valid JSON.

This pattern is intentional: it avoids tight coupling to the Ranger schema and
works across Ranger versions. Preserve it for new resources when the API body
is complex or version-sensitive.

## Import Support

All resources must implement `Importer`. Use `StateContext` (not the older
`State` field). The import ID convention for `ranger_policy` is the numeric
Ranger policy ID.

## Adding a New Resource

1. Create `ranger/resource_<noun>.go`.
2. Implement `resource<Noun>()` returning `*schema.Resource`.
3. Register it in `ranger/provider.go` `ResourcesMap`.
4. Add an example in `examples/resources/ranger_<noun>/resource.tf`.
5. Add an import example in `examples/resources/ranger_<noun>/import.sh`.
6. Run `make docs-gen` to regenerate `docs/`.

## Adding a New Data Source

1. Create `ranger/data_source_<noun>.go`.
2. Implement `dataSource<Noun>()` returning `*schema.Resource`.
3. Register it in `ranger/provider.go` `DataSourcesMap`.
4. Add an example in `examples/data-sources/ranger_<noun>/data-source.tf`.
5. Run `make docs-gen`.

## Quirks to Be Aware Of

- `resourcePolicyCreate` has a `time.Sleep(2 * time.Second)` before delegating
  to Read. This works around a Ranger API propagation delay. Keep it unless the
  underlying issue is fixed upstream.
- Read looks up policies by `service` + `name`, not by ID. The ID stored in
  state is only used for Update and Delete.
