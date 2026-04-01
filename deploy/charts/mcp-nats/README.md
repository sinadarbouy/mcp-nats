# mcp-nats Helm chart

Helm chart for deploying the [mcp-nats](https://github.com/sinadarbouy/mcp-nats) MCP server on Kubernetes.

**Requirements:** Kubernetes **1.25+** (see `Chart.yaml`) and **Helm 3**. OCI registry install and `Chart.yaml` dependencies need **Helm 3.8+** (OCI support).

## Install / upgrade / uninstall

From the **repository root**:

```sh
# Install (release name and namespace are yours to choose)
helm install mcp-nats ./deploy/charts/mcp-nats --namespace mcp-nats --create-namespace

# Inspect rendered manifests
helm template mcp-nats ./deploy/charts/mcp-nats

# Upgrade after changing values
helm upgrade --install mcp-nats ./deploy/charts/mcp-nats --namespace mcp-nats -f my-values.yaml

# Remove
helm uninstall mcp-nats --namespace mcp-nats
```

From **this chart directory** (`deploy/charts/mcp-nats`):

```sh
helm install mcp-nats . --namespace mcp-nats --create-namespace
helm upgrade --install mcp-nats . --namespace mcp-nats -f my-values.yaml
```

## Install from OCI (GitHub Container Registry)

Tagged releases publish this chart to **GHCR** (see the repository [release workflow](https://github.com/sinadarbouy/mcp-nats/blob/main/.github/workflows/release.yml)). The OCI registry path is `oci://ghcr.io/sinadarbouy/charts`; the chart artifact name is **`mcp-nats`**.

Set **`version`** to the chart version you want (it must match a published tag’s chart, for example `0.1.4`).

**Direct install / upgrade**

```sh
# Authenticate when the registry requires it (private packages or rate limits)
helm registry login ghcr.io

helm install mcp-nats oci://ghcr.io/sinadarbouy/charts/mcp-nats \
  --version "0.1.4" \
  --namespace mcp-nats \
  --create-namespace

helm upgrade --install mcp-nats oci://ghcr.io/sinadarbouy/charts/mcp-nats \
  --version "0.1.4" \
  --namespace mcp-nats \
  -f my-values.yaml
```

Default [values.yaml](values.yaml) points at **`ghcr.io/sinadarbouy/mcp-nats`** with a tag matching chart **`appVersion`**. Override `image.registry`, `image.repository`, or `image.tag` if you use a different registry or pin another digest.

**As a dependency of another chart**

In the parent chart’s `Chart.yaml`:

```yaml
apiVersion: v2
name: my-platform
version: 0.1.0
dependencies:
  - name: mcp-nats
    version: "0.1.4"
    repository: "oci://ghcr.io/sinadarbouy/charts"
```

Then from the parent chart directory:

```sh
helm dependency update
helm install my-release . --namespace my-namespace
```

Helm resolves `name: mcp-nats` against that OCI repository (chart URL becomes `oci://ghcr.io/sinadarbouy/charts/mcp-nats`). Use a **GitHub personal access token** with `read:packages` as the Helm registry password if pulls are not anonymous.

## Values overview

| Area | Values keys | Notes |
|------|-------------|--------|
| Image | `image.repository`, `image.tag`, `image.registry`, `image.pullPolicy` | Default is **`ghcr.io/sinadarbouy/mcp-nats`**; `tag` defaults to chart `appVersion` when empty. |
| MCP HTTP | `server.transport`, `server.address`, `server.endpointPath`, `containerPort`, `service.*` | Default transport is `streamable-http` on port **8000**, path **`/mcp`**. |
| NATS | `nats.url`, `nats.noAuthentication` | Sets `NATS_URL` and `NATS_NO_AUTHENTICATION` on the pod. |
| User/password auth | `auth.existingSecret`, `auth.createSecret` / `auth.secretData` | When `nats.noAuthentication` is `false`, use an existing Secret or let the chart create one. |
| Credential env from Secrets | `auth.extraSecretEnv` | Adds env vars from Secret keys (for example `NATS_SYS_CREDS` / `NATS_A_CREDS` as **base64-encoded** `.creds` file contents—see the application [README](../../README.md)). |
| Pod extras | `podAnnotations`, `podLabels`, `command`, `args`, `extraArgs`, `env`, `envFrom`, `volumeMounts`, `volumes` | Use `podAnnotations` for Vault Agent Injector (see below). A writable **`/tmp`** `emptyDir` is included by default for NATS credential temp files when the root filesystem is read-only. |
| Service account | `serviceAccount.create`, `serviceAccount.name`, `serviceAccount.annotations` | Create a dedicated SA for Vault Kubernetes auth or cloud IAM where needed. |
| Ingress / Gateway API | `ingress.*`, `route.main.*` | Optional exposure outside the cluster. |

See [values.yaml](values.yaml) for the full list and defaults.

## Quick examples (`--set`)

Adjust `./deploy/charts/mcp-nats` if you run Helm from the chart directory (use `.`).

Streamable HTTP with explicit listen address and path (defaults already match this, shown for clarity):

```sh
helm install mcp-nats ./deploy/charts/mcp-nats \
  --set server.transport=streamable-http \
  --set server.address=0.0.0.0:8000 \
  --set server.endpointPath=/mcp
```

Anonymous NATS (no credentials):

```sh
helm install mcp-nats ./deploy/charts/mcp-nats \
  --set nats.url=nats://nats.default.svc:4222 \
  --set nats.noAuthentication=true
```

User/password NATS with an existing Kubernetes Secret:

```sh
kubectl create secret generic mcp-nats-auth \
  --from-literal=NATS_USER=myuser \
  --from-literal=NATS_PASSWORD=mypass

helm install mcp-nats ./deploy/charts/mcp-nats \
  --set nats.url=nats://nats.default.svc:4222 \
  --set nats.noAuthentication=false \
  --set auth.existingSecret.name=mcp-nats-auth
```

## Health probes

The chart wires HTTP probes to the MCP server: `startupProbe` and `livenessProbe` use **`/livez`**, `readinessProbe` uses **`/readyz`**. A `preStop` sleep helps drain endpoints during rollouts. Override any probe field via values (same keys as a container spec).

```sh
helm upgrade --install mcp-nats ./deploy/charts/mcp-nats \
  --set readinessProbe.httpGet.path=/readyz \
  --set livenessProbe.httpGet.path=/livez \
  --set startupProbe.failureThreshold=18
```

## Custom container command (Vault and other init patterns)

If you set `command`, the chart does **not** append the default `--transport` / `--address` arguments (so a `sh -c` wrapper can run a single composed command). Put the full `mcp-nats` invocation inside that script, or set `args` / `extraArgs` if you use a different entrypoint.

## Example: HashiCorp Vault Agent Injector (NATS credentials)

Use the Vault Agent sidecar to render a file (for example shell `export` lines) and source it before starting `mcp-nats`. Annotations go on the **pod** via `podAnnotations`. Align the Vault role, secret path, and KV version with your Vault deployment; the snippet below matches a common KV v2 layout.

**Encoding the `.creds` file for Vault**

`mcp-nats` expects `NATS_*_CREDS` to be a **single-line standard base64** string of the NATS user `.creds` file (for example from NSC). One way to produce that and paste into Vault KV (field `aCreds` in the example below):

```sh
# macOS (BSD base64): -i is the input file
base64 -i nats/nsc/keys/creds/OP/A/a.creds | tr -d '\n'

# Linux (GNU coreutils): no line wrapping
base64 -w 0 nats/nsc/keys/creds/OP/A/a.creds
```

Store that output in Vault; the Agent template then exports it as `NATS_A_CREDS`. Do not add quotes or newlines inside the stored value.

Create a values file (for example `values-vault.yaml`). The template assumes Vault key `aCreds` holds that base64 string.

```yaml
serviceAccount:
  create: true
  name: "nats-mcp-sa"

podAnnotations:
  vault.hashicorp.com/agent-inject: "true"
  vault.hashicorp.com/agent-inject-secret-creds: apps/kv/data/nats-mcp/creds
  vault.hashicorp.com/agent-inject-status: "update"
  vault.hashicorp.com/agent-inject-template-creds: |
    {{- with secret "apps/kv/data/nats-mcp/creds" -}}
    export NATS_A_CREDS="{{ .Data.data.aCreds }}"
    {{- end }}
  vault.hashicorp.com/agent-limits-cpu: 50m
  vault.hashicorp.com/agent-limits-mem: 64Mi
  vault.hashicorp.com/agent-pre-populate-only: "true"
  vault.hashicorp.com/agent-requests-cpu: 10m
  vault.hashicorp.com/agent-requests-mem: 16Mi
  vault.hashicorp.com/role: setup-role-nats-mcp
  vault.hashicorp.com/secret-volume-path: /usr/src/app/credentials

nats:
  url: nats://nats-svc.nats:4222
  noAuthentication: false

volumeMounts:
  - name: tmp
    mountPath: /tmp
volumes:
  - name: tmp
    emptyDir: {}

command:
  - /bin/sh
  - -c
  - . /usr/src/app/credentials/creds && exec /app/mcp-nats --transport streamable-http --address 0.0.0.0:8000 --endpoint-path /mcp --read-only
```

Install (from repository root):

```sh
helm upgrade --install mcp-nats ./deploy/charts/mcp-nats --namespace mcp-nats -f values-vault.yaml
```

Notes for this pattern:

- The injector writes a file under `secret-volume-path`; the name **`creds`** comes from the annotation suffix `agent-inject-secret-creds` (secret name `creds` in the volume). The `command` sources that file before `exec`.
- After sourcing, `NATS_A_CREDS` must stay a **single-line base64** `.creds` payload (see **Encoding the `.creds` file for Vault** above). The app decodes it and writes a temp file for the NATS CLI.
- Configure a Vault **role** bound to this Kubernetes service account and grant read access to the path used in `agent-inject-secret-*` / `secret "..."` in the template.
- If your cluster requires the projected service account token for Vault Kubernetes auth, set `serviceAccount.automountServiceAccountToken` appropriately for your policy (see [values.yaml](values.yaml)).

## Local integration test (Tilt)

End-to-end auth testing with NATS and this chart is described in the repository [README](../../README.md#tilt-integration-test-docker-desktop-kubernetes) (Tilt + Docker Desktop Kubernetes).
