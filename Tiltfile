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

local_resource(
    "build-mcp-nats-image",
    cmd="docker build -f deploy/tilt/Dockerfile.tilt -t mcp-nats:tilt .",
    resource_deps=["nats"],
    allow_parallel=False,
)

local_resource(
    "mcp-nats",
    cmd=(
        HELM_ENV
        + " helm upgrade --install "
        + MCP_RELEASE
        + " ./deploy/charts/mcp-nats "
        + "--namespace "
        + NAMESPACE
        + " --create-namespace "
        + "-f deploy/tilt/mcp-nats-values.yaml"
    ),
    resource_deps=["build-mcp-nats-image"],
    allow_parallel=False,
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
        + " svc/"
        + MCP_RELEASE
        + "-mcp-nats 8000:8000'"
    ),
    resource_deps=["mcp-nats"],
    trigger_mode=TRIGGER_MODE_MANUAL,
)
