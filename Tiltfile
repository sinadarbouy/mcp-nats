allow_k8s_contexts("docker-desktop")

NAMESPACE = "mcp-nats-tilt"
NATS_RELEASE = "nats"
MCP_RELEASE = "mcp-nats"

HELM_ENV = "HELM_CACHE_HOME=$PWD/.helm/cache HELM_CONFIG_HOME=$PWD/.helm/config HELM_DATA_HOME=$PWD/.helm/data"

local_resource(
    "helm-repo-nats",
    cmd=(
        "mkdir -p .helm/cache .helm/config .helm/data && "
        + HELM_ENV
        + " helm repo add nats https://nats-io.github.io/k8s/helm/charts/ --force-update && "
        + HELM_ENV
        + " helm repo update nats"
    ),
    allow_parallel=False,
)

local_resource(
    "nats",
    cmd=(
        HELM_ENV
        + " helm upgrade --install "
        + NATS_RELEASE
        + " nats/nats "
        + "--namespace "
        + NAMESPACE
        + " --create-namespace "
        + "-f deploy/tilt/nats-values.yaml"
    ),
    resource_deps=["helm-repo-nats"],
    allow_parallel=False,
)

# Managed image build + cluster wiring: Tilt rewrites the Deployment image to a
# unique tag so nodes never keep serving an old blob when reuse_tag mcp-nats:tilt
# is rebuilt locally (fixes "flag provided but not defined: -address" from stale binaries).
docker_build(
    "mcp-nats:tilt",
    ".",
    dockerfile="deploy/tilt/Dockerfile.tilt",
)

k8s_yaml(helm(
    "./deploy/charts/mcp-nats",
    name=MCP_RELEASE,
    namespace=NAMESPACE,
    values=["deploy/tilt/mcp-nats-values.yaml"],
))

k8s_resource(
    "mcp-nats",
    resource_deps=["nats"],
    labels=["mcp"],
)

local_resource(
    "status",
    cmd=(
        "kubectl get pods,svc -n "
        + NAMESPACE
        + " && "
        + "echo && "
        + "echo 'Use: kubectl port-forward -n "
        + NAMESPACE
        + " svc/mcp-nats 8000:8000'"
    ),
    resource_deps=["mcp-nats"],
    trigger_mode=TRIGGER_MODE_MANUAL,
)
