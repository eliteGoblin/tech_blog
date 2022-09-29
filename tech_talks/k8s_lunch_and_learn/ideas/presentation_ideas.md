

## Ideas

*  This draft provide list of possible topics we could share across teams, welcome team members to contribute this list, add more topic, ideas; and pick up topic to share.
*  Demo should be an important part of presentation.
*  It is better to keep the talk as code, for audience to review and try the demo, and also make demo self-contained as possible(eg. use docker-compose) to make it easier to try.
*  We can maintain a list of talks we already give, and links to presentation files.

###  Topic: Why k8s: from micro-service's view

*  Migration from monolith to microservice: complicated infra
*  Use a typical microservice as an example, show k8s's benefits
    -  Understanding the 8 Fallacies of Distributed Systems
*  Terminalogy explain briefly
    -  pods
    -  service
    -  deployment/replicateset
*  Challenges in a typical micro-services system
    -  Need to know your hosts/OS, capacity planning
    -  Difficult to scale horizonally
    -  Service discovery
    -  Coupling between services: OOM, CPU
    -  Complicated release: zero downtime
*  K8s Cool features(to solve previous problems)
    -  Pod scheduler 
        +  Visualize pods schedule with kops-view
        +  Take down a node, see pods scheduled
    -  Rolling update
        +  Show logs, visualize like https://octopus.com/blog/kubernetes-deployment-strategies-visualized
    -  Scaling easily: vertical, horizonal, node 
        +  demo of HPA, cluster autoscaler
    -  Service discovery: DNS

## Possible topic: Things you should know about Docker

*  What is docker/container, and why
    -  Realword examples: CircleCI, Bamboo run in docker...
*  How it works(client-server mode, docker-api, docker run demo)
*  Dockerfile walkthrough
*  Docker image and repo
*  Data sharing and volume
*  Docker's network basic
*  Demos
    -  How API team run integration test inside docker-compose:run docker in docker-compose
    -  Host a wordpress inside docker-compose: https://docs.docker.com/compose/wordpress/
*  Docker internal: Linux namespace, cgroup, process model
*  Other container runtime: rkt, CRI-O..

### Possible Topic: Create your own k8s cluster

*  play locally: Minikube
*  Kops
    -  design and pholosiphy
        +  state store
        +  philosophy
        +  boot sequence
    - Network model: VPC, subnet, multi-master in 3 AZs, bastion
    - Demo: Use kops and terraform to create gatling cluster
*  Create it manually on GCP: K8s the hardway: https://github.com/kelseyhightower/kubernetes-the-hard-way
    -  major steps walkthrough, to explain the key components of K8s

### Possible Topic: scaling in K8s

*  Resources in k8s/docker
*  HPA
*  Vertical
*  Cluster autoscaler
*  Some internel explain of HPA, cluster-autoscaler

### Possible Topic: K8s internels in depth

*  K8s objects introduction: 
    -  service
    -  deployments
    -  ingress
    -  custom resource
*  DNS of k8s
*  K8s network specification and CNI
*  How declare config works in K8s: controller loop
    -  how list-watch implemented(not polling): HTTP Chunked transfer encoding
*  Deployments and rolling update
    -  liveness, readiness
    -  maxSurge, maxUnavailable, minReadySeconds

Note: could be breaked into several sessions

### Possible Topic: Build CI/CD pipeline on K8s

*  Current CI approach:
*  Current CD: Kcd
*  Gitops
*  ArgoCD: declarative, GitOps continuous delivery tool for k8s

### Possible Topic: Etcd and distributed consensus

*  Why Etcd in k8s
*  Etcd internal, RAFT
*  Demo: use hasicorp/raft to implement distributed consensus;(distributed lock)
*  gRPC
*  Build a HA k8s cluster with ETCD cluster
*  Backup of ETCD

### Possible Topic: Realtime monitoring with Prometheus and Grafana

*  Prometheus architecture
*  PromQL
*  Demo: build monitoring dashboard for go-application

### Images processing and CV

### Istio topics

*  Distributed tracing
*  Visualize services
*  grpc load balancing

