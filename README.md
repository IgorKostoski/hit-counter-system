# Hit Counter System - DevOps Portfolio Project

This project demonstrates a range of DevOps practices and technologies by building a simple "Hit Counter" microservice
and wrapping it with a comprehensive infrastructure and CI/CD setup. The goal is to showcase skills in Go development,
containerization, Kubernetes, Infrastructure as Code (IaC), monitoring, artifact management, and automation.

**Live Demo (Conceptual - Not Deployed Publicly for this Portfolio Version)**

* API Endpoint: `http://hit-counter.local/api/v1/...` (when running locally with Ingress and `minikube tunnel`)
* Grafana Dashboard: `http://localhost:3000` (when port-forwarded)

## Features

* **Go API Service:**
    * `/api/v1/hit` (POST): Increments a counter for a given key in PostgreSQL.
    * `/api/v1/count/{key}` (GET): Retrieves the current count for a key.
    * `/metrics`: Exposes Prometheus metrics.
* **PostgreSQL Database:** Stores hit counts.
* **Containerization:** Dockerized Go application and database.
* **Orchestration:** Deployed on a local Kubernetes cluster (Minikube).
* **Package Management (K8s):** Helm chart for deploying the application and its dependencies.
* **Infrastructure as Code (IaC):**
    * Terraform for managing specific Kubernetes resources (e.g., ServiceMonitor).
    * Ansible for system configuration (demonstrated with a separate Nginx example) and K8s resource management (
      ConfigMap example).
* **Monitoring:**
    * Prometheus for metrics collection.
    * Grafana for metrics visualization.
* **CI/CD:**
    * GitHub Actions for basic Go app CI (build, test).
    * Azure DevOps Pipeline (`azure-pipelines.yml`) defined for a more complete CI/CD flow (build, Docker image, Helm
      package, simulated deployment).
* **Artifact Management (Simulated with Local Nexus):**
    * Local Nexus OSS instance for hosting Docker images and Helm charts.
* **Ingress:** Nginx Ingress Controller for exposing the API service via a hostname.
* **Automation Scripts:** Helper scripts for common tasks.

## Technologies Used

* **Programming Language:** Go
* **Database:** PostgreSQL
* **Containerization:** Docker
* **Orchestration:** Kubernetes (Minikube for local cluster)
* **K8s Package Management:** Helm
* **Infrastructure as Code (IaC):**
    * Terraform (Kubernetes provider)
    * Ansible
* **Monitoring:**
    * Prometheus (`kube-prometheus-stack` Helm chart)
    * Grafana
* **CI/CD:**
    * GitHub Actions
    * Azure DevOps (YAML Pipeline definition)
* **Artifact Repository:** Sonatype Nexus Repository OSS (run locally via Docker)
* **Web Server/Reverse Proxy:** Nginx (for Ingress Controller and Ansible demo)
* **Automation/Scripting:** Bash
* **Version Control:** Git, GitHub

## Project Structure

```bash
hit-counter-system/
    ├── app/ # Go application source
    │ ├── api/ # Main API service (main.go, handlers.go, db.go, Dockerfile)
    │ └── ...
    ├── ansible/ # Ansible playbooks, roles, inventory
    │ ├── roles/webserver/ # Role to configure Nginx on a target
    │ ├── Dockerfile.ansible-target # Dockerfile for Ansible target node
    │ ├── inventory.ini
    │ ├── k8s_configmap_playbook.yml
    │ └── playbook.yml # Main playbook for webserver role
    ├── helm-charts/ # Helm charts
    │ └── hit-counter-system/ # Umbrella chart for the application & PostgreSQL
    ├── helm-packages/ # Stores packaged .tgz Helm charts (local use)
    ├── kubernetes/ # Raw Kubernetes YAML manifests (initial setup, now managed by Helm/Terraform)
    ├── monitoring-stack/ # Configuration for Prometheus/Grafana stack
    │ ├── kube-prometheus-stack-values.yaml
    │ └── api-servicemonitor.yaml
    ├── terraform/ # Terraform configurations
    │ └── ... # (providers.tf, main.tf for ServiceMonitor, etc.)
    ├── .github/workflows/ # GitHub Actions workflows
    │ └── go-ci.yml
    ├── docker-compose.yml # For local development stack (Go API, PG, Nexus)
    ├── azure-pipelines.yml # Azure DevOps CI/CD pipeline definition
    ├── README.md # This file
    └── .gitignore
```


## Prerequisites

*   Go (version 1.21+ recommended)
*   Docker & Docker Compose
*   Minikube (or another local Kubernetes cluster like Kind, k3d)
*   `kubectl`
*   Helm (version 3+)
*   Terraform (version 1.x recommended)
*   Ansible (version 2.10+ recommended, with `kubernetes.core` collection)
*   `curl` or a similar HTTP client (for testing API)
*   (Optional) An account on Docker Hub / Azure Container Registry if you wish to push images externally.
*   (Optional) An Azure DevOps organization/project if you wish to run the pipeline.

## Setup and Running the Project

### 1. Clone the Repository
```bash
git clone https://github.com/IgorKostoski/hit-counter-system.git
cd hit-counter-system
```

### 2.Start Local Kubernetes Cluster(Minikube)
```bash
minikube start --cpus=4 --memory=8192mb # Adjust resources as needed
# For local image usage with application deployment:
eval $(minikube docker-env) # Run in terminals where you build Docker images for the app
```

### 3.Local Dvelopment Stack (Go API, PostgreSQL, Nexus) using Docker Compose
```bash
# Build the Go API Docker image (if not already built by Minikube's Docker daemon for K8s deployment)
# cd app/api && docker build -t hit-counter-api:dev . && cd ../..

# Start services
docker-compose up -d --build

# Initial Nexus Setup (first time only):
# 1. Wait a few minutes for Nexus to start (docker-compose logs -f nexus)
# 2. Access Nexus UI: http://localhost:8081
# 3. Get admin password: docker exec hit-counter-nexus cat /nexus-data/admin.password
# 4. Login as 'admin', change password, complete setup wizard.
# 5. Create repositories in Nexus UI (see "Artifact Management (Nexus)" section below).
# 6. Configure Docker insecure registry for localhost:8082 (see Nexus section).

```

### 4.Build Application Docker Image (for Kubernetes)
```bash
cd app/api
docker build -t hit-counter-api:k8s-v1 . # Or any tag you prefer
cd ../..
```

### 5.Deploy Monitoring Stack (Prometheus & Grafana)
```bash
# Create namespace
kubectl create namespace monitoring

# Install kube-prometheus-stack
cd monitoring-stack # (Assuming you are in hit-counter-system/ root)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  -f ./kube-prometheus-stack-values.yaml \
  
cd ..
```