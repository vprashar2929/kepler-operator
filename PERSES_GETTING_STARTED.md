# Getting Started with Perses on Kubernetes

## Table of Contents

- [What is Perses?](#what-is-perses)
- [Prerequisites](#prerequisites)
- [Installing Perses Operator](#installing-perses-operator)
- [Deploying Perses](#deploying-perses)
- [Understanding Datasources](#understanding-datasources)
- [Creating Dashboards](#creating-dashboards)
- [Advanced Topics](#advanced-topics)
- [Troubleshooting](#troubleshooting)
- [Resources](#resources)

---

## What is Perses?

**Perses** is an open-source dashboard and visualization tool designed specifically for observability data. Think of it as a modern alternative to Grafana, built with a cloud-native approach.

### Key Features

- 📊 **Multi-datasource support**: Works with Prometheus, Tempo, Loki, and Pyroscope
- 🎨 **Modern UI**: Clean, intuitive interface for creating dashboards
- ☸️ **Kubernetes-native**: First-class Kubernetes support with CRDs (Custom Resource Definitions)
- 🔄 **GitOps-friendly**: Define dashboards and datasources as code
- 🚀 **Lightweight**: Minimal resource footprint

---

## Prerequisites

Before you begin, ensure you have:

1. **A running Kubernetes cluster** (Minikube, Kind, or any production cluster)
2. **kubectl** installed and configured
3. **Basic understanding of Kubernetes concepts** (Pods, Services, ConfigMaps)
4. **Optional but recommended**: A Prometheus instance running in your cluster

To check your setup:

```bash
# Verify kubectl is working
kubectl cluster-info

# Check your current context
kubectl config current-context
```

---

## Installing Perses Operator

The **Perses Operator** is a Kubernetes operator that manages Perses instances and resources (dashboards, datasources) using Custom Resource Definitions (CRDs).

### Step 1: Clone the Perses Operator Repository

```bash
git clone https://github.com/perses/perses-operator.git
cd perses-operator
```

### Step 2: Install Custom Resource Definitions (CRDs)

CRDs extend Kubernetes to understand Perses-specific resources like `PersesDashboard`, `GlobalDatasource`, etc.

```bash
# Install CRDs
make install
```

This command installs the following CRDs:

- `Perses`: The main Perses instance
- `PersesDashboard`: Dashboard definitions
- `GlobalDatasource`: Cluster-wide datasources
- `Datasource`: Project-scoped datasources
- And more...

You can verify the installation:

```bash
kubectl get crds | grep perses
```

### Step 3: Deploy the Operator

```bash
# Deploy the operator to your cluster
make deploy
```

This creates:

- A namespace (usually `perses-operator-system`)
- The operator deployment
- Necessary RBAC (Role-Based Access Control) resources

Verify the operator is running:

```bash
kubectl get pods -n perses-operator-system
```

You should see something like:

```text
NAME                                                   READY   STATUS    RESTARTS   AGE
perses-operator-controller-manager-xxxxxxxxxx-xxxxx   2/2     Running   0          1m
```

---

## Deploying Perses

Now that the operator is installed, let's deploy a Perses instance.

### Step 1: Create a Namespace for Perses

```bash
kubectl create namespace perses
```

### Step 2: Create a Perses Instance

Create a file named `perses-instance.yaml`:

```yaml
apiVersion: perses.dev/v1alpha1
kind: Perses
metadata:
  name: perses-demo
  namespace: perses
spec:
  # Number of Perses replicas
  replicas: 1

  # Resource limits (adjust based on your needs)
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi

  # Storage configuration
  database:
    # Use file-based storage (simple for getting started)
    file:
      folder: /var/lib/perses
      # Optionally use a PersistentVolumeClaim
      # persistence:
      #   enabled: true
      #   storageClass: standard
      #   size: 1Gi

  # Service configuration
  service:
    type: ClusterIP
    port: 8080
```

Apply the configuration:

```bash
kubectl apply -f perses-instance.yaml
```

### Step 3: Verify the Deployment

```bash
# Check if the Perses pod is running
kubectl get pods -n perses

# Check the service
kubectl get svc -n perses
```

### Step 4: Access the Perses UI

Forward the port to access Perses locally:

```bash
kubectl port-forward -n perses svc/perses-demo 8080:8080
```

Open your browser and navigate to: <http://localhost:8080>

🎉 Congratulations! Perses is now running!

---

## Understanding Datasources

Datasources in Perses are connections to external systems that provide data for your dashboards (e.g., Prometheus, Loki, Tempo).

### Datasource Scopes

Perses supports three scopes for datasources:

1. **Global Datasource**: Available to all projects and dashboards across the entire cluster
2. **Project Datasource**: Available to all dashboards within a specific project
3. **Dashboard Datasource**: Only available within a single dashboard

### 1. Global Datasource

Global datasources are defined cluster-wide and can be used by any dashboard.

#### Example: Global Prometheus Datasource

Create a file named `global-prometheus-datasource.yaml`:

```yaml
apiVersion: perses.dev/v1alpha1
kind: GlobalDatasource
metadata:
  name: prometheus-global
spec:
  # Display name shown in the UI
  display:
    name: "Global Prometheus"
    description: "Cluster-wide Prometheus datasource"

  # Mark as default datasource for this type
  default: true

  # Plugin configuration
  plugin:
    kind: PrometheusDatasource
    spec:
      # URL to your Prometheus server
      # If Prometheus is in the same cluster:
      directUrl: "http://prometheus-server.monitoring.svc.cluster.local:9090"

      # For external Prometheus:
      # directUrl: "http://prometheus.example.com:9090"

      # Optional: Scrape interval
      # scrapeInterval: "15s"
```

Apply it:

```bash
kubectl apply -f global-prometheus-datasource.yaml
```

#### Other Datasource Types

- **Tempo** (for traces):

  ```yaml
  plugin:
    kind: TempoDatasource
    spec:
      directUrl: "http://tempo.monitoring.svc.cluster.local:3200"
  ```

- **Loki** (for logs):

  ```yaml
  plugin:
    kind: LokiDatasource
    spec:
      directUrl: "http://loki.monitoring.svc.cluster.local:3100"
  ```

### 2. Project Datasource

Project datasources are scoped to a specific project (a namespace-like concept in Perses).

Create a file named `project-datasource.yaml`:

```yaml
apiVersion: perses.dev/v1alpha1
kind: Datasource
metadata:
  name: my-project-prometheus
  # This datasource belongs to the "my-project" project
  project: my-project
spec:
  display:
    name: "My Project Prometheus"
  plugin:
    kind: PrometheusDatasource
    spec:
      directUrl: "http://prometheus-server.monitoring.svc.cluster.local:9090"
```

### 3. Dashboard Datasource

Datasources can also be defined directly within a dashboard. See the [Creating Dashboards](#creating-dashboards) section below.

---

## Creating Dashboards

Dashboards in Perses consist of **panels** (individual visualizations) arranged in **layouts** (grid positions).

### Basic Dashboard Structure

```yaml
apiVersion: perses.dev/v1alpha1
kind: PersesDashboard
metadata:
  name: my-first-dashboard
  namespace: perses
spec:
  # Project this dashboard belongs to
  project: default

  # Display information
  display:
    name: "My First Dashboard"
    description: "A beginner's dashboard"

  # Time range
  duration: "1h"

  # Auto-refresh interval
  refreshInterval: "30s"

  # Datasources (optional - can also use global datasources)
  datasources:
    prometheus:
      kind: PrometheusDatasource
      spec:
        directUrl: "http://prometheus-server.monitoring.svc.cluster.local:9090"

  # Panels - the actual visualizations
  panels:
    panel-1:
      kind: Panel
      spec:
        display:
          name: "CPU Usage"
        plugin:
          kind: TimeSeriesChart
          spec: {}
        queries:
          - kind: TimeSeriesQuery
            spec:
              plugin:
                kind: PrometheusTimeSeriesQuery
                spec:
                  datasource:
                    kind: PrometheusDatasource
                  query: "rate(node_cpu_seconds_total{mode='user'}[5m])"

  # Layouts - how panels are arranged
  layouts:
    - kind: Grid
      spec:
        items:
          - x: 0      # X position (starts at 0)
            y: 0      # Y position (starts at 0)
            width: 12 # Width (grid is 24 units wide)
            height: 6 # Height (arbitrary units)
            content:
              $ref: "#/spec/panels/panel-1"
```

### Complete Example: Monitoring Dashboard

Here's a more complete example with multiple panels:

Create a file named `monitoring-dashboard.yaml`:

```yaml
apiVersion: perses.dev/v1alpha1
kind: PersesDashboard
metadata:
  name: kubernetes-monitoring
  namespace: perses
spec:
  project: default

  display:
    name: "Kubernetes Monitoring"
    description: "Overview of cluster metrics"

  duration: "6h"
  refreshInterval: "1m"

  # Variables for dynamic queries
  variables:
    - kind: ListVariable
      spec:
        name: namespace
        display:
          name: "Namespace"
        allowAllValue: true
        allowMultiple: false
        plugin:
          kind: PrometheusLabelValuesVariable
          spec:
            datasource:
              kind: PrometheusDatasource
            labelName: namespace
            matchers:
              - "up"

  # Panels
  panels:
    cpu-usage:
      kind: Panel
      spec:
        display:
          name: "CPU Usage by Namespace"
          description: "Rate of CPU usage across namespaces"
        plugin:
          kind: TimeSeriesChart
          spec:
            legend:
              position: bottom
        queries:
          - kind: TimeSeriesQuery
            spec:
              plugin:
                kind: PrometheusTimeSeriesQuery
                spec:
                  datasource:
                    kind: PrometheusDatasource
                  query: 'sum by (namespace) (rate(container_cpu_usage_seconds_total{namespace=~"$namespace"}[5m]))'
                  seriesNameFormat: "{{namespace}}"

    memory-usage:
      kind: Panel
      spec:
        display:
          name: "Memory Usage"
        plugin:
          kind: TimeSeriesChart
          spec: {}
        queries:
          - kind: TimeSeriesQuery
            spec:
              plugin:
                kind: PrometheusTimeSeriesQuery
                spec:
                  datasource:
                    kind: PrometheusDatasource
                  query: 'sum by (namespace) (container_memory_usage_bytes{namespace=~"$namespace"})'
                  seriesNameFormat: "{{namespace}}"

    pod-count:
      kind: Panel
      spec:
        display:
          name: "Running Pods"
        plugin:
          kind: StatChart
          spec:
            calculation: last
            unit:
              kind: decimal
        queries:
          - kind: TimeSeriesQuery
            spec:
              plugin:
                kind: PrometheusTimeSeriesQuery
                spec:
                  datasource:
                    kind: PrometheusDatasource
                  query: 'count(kube_pod_info{namespace=~"$namespace"})'

  # Layout - arrange panels in a grid
  layouts:
    - kind: Grid
      spec:
        items:
          # CPU panel - full width at top
          - x: 0
            y: 0
            width: 24
            height: 8
            content:
              $ref: "#/spec/panels/cpu-usage"

          # Memory panel - left side
          - x: 0
            y: 8
            width: 16
            height: 8
            content:
              $ref: "#/spec/panels/memory-usage"

          # Pod count - right side stat
          - x: 16
            y: 8
            width: 8
            height: 8
            content:
              $ref: "#/spec/panels/pod-count"
```

Apply the dashboard:

```bash
kubectl apply -f monitoring-dashboard.yaml
```

### Understanding Panel Types

Perses supports several panel types:

1. **TimeSeriesChart**: Line charts for time-series data
2. **StatChart**: Single value statistics
3. **GaugeChart**: Gauge visualization
4. **BarChart**: Bar charts
5. **TableChart**: Data tables
6. **MarkdownChart**: Markdown text panels

---

## Advanced Topics

### Managing Resources with ConfigMaps

For GitOps workflows, you can manage dashboards using ConfigMaps. This is useful when Perses uses filesystem storage.

Create a file named `dashboard-configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-dashboard-config
  namespace: perses
  labels:
    # This label tells Perses to load this as a dashboard
    perses.dev/resource: "true"
data:
  my-dashboard.json: |
    {
      "kind": "Dashboard",
      "metadata": {
        "name": "configmap-dashboard",
        "project": "default"
      },
      "spec": {
        "display": {
          "name": "Dashboard from ConfigMap",
          "description": "This dashboard is managed via ConfigMap"
        },
        "duration": "1h",
        "panels": {
          "simple-panel": {
            "kind": "Panel",
            "spec": {
              "display": {
                "name": "Simple Metric"
              },
              "plugin": {
                "kind": "TimeSeriesChart",
                "spec": {}
              },
              "queries": [
                {
                  "kind": "TimeSeriesQuery",
                  "spec": {
                    "plugin": {
                      "kind": "PrometheusTimeSeriesQuery",
                      "spec": {
                        "datasource": {
                          "kind": "PrometheusDatasource"
                        },
                        "query": "up"
                      }
                    }
                  }
                }
              ]
            }
          }
        },
        "layouts": [
          {
            "kind": "Grid",
            "spec": {
              "items": [
                {
                  "x": 0,
                  "y": 0,
                  "width": 24,
                  "height": 6,
                  "content": {
                    "$ref": "#/spec/panels/simple-panel"
                  }
                }
              ]
            }
          }
        ]
      }
    }
```

Apply it:

```bash
kubectl apply -f dashboard-configmap.yaml
```

### Datasource Discovery

Perses can automatically discover datasources in your cluster using Kubernetes service discovery.

Add this to your Perses instance configuration:

```yaml
apiVersion: perses.dev/v1alpha1
kind: Perses
metadata:
  name: perses-demo
  namespace: perses
spec:
  # ... other spec fields ...

  config:
    datasource:
      global:
        discovery:
          - name: "auto-discover-prometheus"
            kubernetes_sd:
              datasource_plugin_kind: "PrometheusDatasource"
              namespace: "monitoring"
              service_configuration:
                enable: true
                port_name: "http"
                serviceType: "ClusterIP"
              labels:
                app: prometheus
```

This will automatically discover Prometheus services in the `monitoring` namespace with the label `app: prometheus`.

---

## Troubleshooting

### Dashboard Not Showing Up

1. **Check if the dashboard was created:**

   ```bash
   kubectl get persesdashboards -n perses
   ```

2. **Check dashboard details:**

   ```bash
   kubectl describe persesdashboard <dashboard-name> -n perses
   ```

3. **Check Perses logs:**

   ```bash
   kubectl logs -n perses deployment/perses-demo
   ```

### Datasource Connection Issues

1. **Test connectivity from Perses pod:**

   ```bash
   kubectl exec -n perses deployment/perses-demo -- wget -O- http://prometheus-server.monitoring.svc.cluster.local:9090/-/healthy
   ```

2. **Verify datasource configuration:**

   ```bash
   kubectl get globaldatasources -o yaml
   ```

### No Data in Panels

1. **Verify your Prometheus query works:**
   - Access Prometheus UI directly
   - Run the query to ensure it returns data

2. **Check time range:**
   - Ensure the dashboard time range covers when data exists

3. **Check datasource reference:**
   - Verify the panel is referencing the correct datasource

### Operator Issues

1. **Check operator logs:**

   ```bash
   kubectl logs -n perses-operator-system deployment/perses-operator-controller-manager -c manager
   ```

2. **Verify CRDs are installed:**

   ```bash
   kubectl get crds | grep perses
   ```

---

## Resources

### Official Documentation

- **Perses Documentation**: <https://perses.dev/>
- **Perses GitHub**: <https://github.com/perses/perses>
- **Perses Operator GitHub**: <https://github.com/perses/perses-operator>
- **API Reference**: <https://perses.dev/perses/docs/api/>

### Community

- **Slack**: Join the CNCF Slack and find #perses channel
- **GitHub Discussions**: <https://github.com/perses/perses/discussions>

### Examples

- **Sample Dashboards**: <https://github.com/perses/perses/tree/main/docs/examples>
- **Perses Operator Samples**: <https://github.com/perses/perses-operator/tree/main/config/samples>

### Comparison with Other Tools

- **vs Grafana**: Perses is lighter weight, more Kubernetes-native, and uses a simpler architecture
- **vs Prometheus UI**: Perses provides a better dashboarding experience with support for multiple datasources

---

## Next Steps

Now that you have Perses running, here are some things to try:

1. ✅ **Create your first custom dashboard** based on metrics from your cluster
2. ✅ **Set up multiple datasources** (Prometheus, Loki, Tempo)
3. ✅ **Explore variables** to make dashboards more dynamic
4. ✅ **Organize dashboards into projects** for better management
5. ✅ **Integrate with your GitOps workflow** using ConfigMaps
6. ✅ **Set up datasource discovery** for automatic configuration

---

## Example: Quick Start Script

Here's a complete script to get Perses running quickly:

```bash
#!/bin/bash

# Clone and install operator
git clone https://github.com/perses/perses-operator.git
cd perses-operator
make install
make deploy

# Create namespace
kubectl create namespace perses

# Wait for operator to be ready
kubectl wait --for=condition=available --timeout=300s \
  deployment/perses-operator-controller-manager \
  -n perses-operator-system

# Deploy Perses instance
kubectl apply -f - <<EOF
apiVersion: perses.dev/v1alpha1
kind: Perses
metadata:
  name: perses-demo
  namespace: perses
spec:
  replicas: 1
  database:
    file:
      folder: /var/lib/perses
EOF

# Wait for Perses to be ready
kubectl wait --for=condition=available --timeout=300s \
  deployment/perses-demo \
  -n perses

# Port forward
echo "Perses is ready! Access it at http://localhost:8080"
kubectl port-forward -n perses svc/perses-demo 8080:8080
```

Save this as `install-perses.sh`, make it executable, and run:

```bash
chmod +x install-perses.sh
./install-perses.sh
```

---

## Happy Dashboarding! 📊

If you have questions or run into issues, check the [Troubleshooting](#troubleshooting) section or reach out to the Perses community.
