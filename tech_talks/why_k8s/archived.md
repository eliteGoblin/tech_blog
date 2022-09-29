
## Compared to Linux as an OS

Key words: 

- Container
- Orchestration, managing
- Automation 
- Scaling, Auto heal

## Resource limit in docker container

*  Resources: CPU, memory, GPU
*  Requests means reserve, are what the container is guaranteed to get
*  Limits, make sure a container never goes above a certain value
// CPU limit, memory limit


```
docker run -it --cpus=".5" ubuntu /bin/bash
```

*  1m CPU = 1 / 1000 Core (mili-core)
*  1m Memory = 1 Medibyte(megabyte)

Generally, you should do resource planning before deploy: specify resource request and limit

// millicores. If your container needs two full cores to run, you would put the value “2000m”
//  If your app starts hitting your CPU limits, Kubernetes starts throttling your container. This means the CPU will be artificially restricted, giving your app potentially worse performance! However, it won’t be terminated or evicted

//  mebibyte 2^20
// , if a container goes past its memory limit it will be terminated.

//  it’s important that the containers have enough resources to actually run

// K8s pod schedule build based on it




// ## Proxy

// Create a proxy to remote Pods
// Kubectl proxy save your trouble to put request credential because kubectl already has it.
 //  ![20200705211742](https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20200705211742.png)


 ## Challenges of Microservice

Mainly come from: we break a system to many pieces: distributed system.

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201024121824.png" alt="20201024121824" style="width:200px"/> 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201024121958.png" alt="20201024121958" style="width:200px"/>



## Deploy new version

Now you have a app running in K8s: 

*  Backed by ReplicateSet, say 3 replicas(pods). 
*  Service object point to 3 pods.

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201021092340.png" alt="20201021092340" style="width:450px"/>

You have a new version, how you gonna release it?

// Beijing Udesk story: 

## Old ways

*  'Stone age': tar all binary and config, stop, replace, restart
*  Semi automated tools: Ansible

Imperative, a lot of human intervention
*  Tar related files
*  Copy to all hosts
*  Stop, replace, restart

// Things works in test not always in prod
// boring, risky, important event, make people nervous
// Can miss some, at some time, or some instance unhealthy.

// Release in Amap: provide map service in China, long list of hosts, shell script to release, silently fail, inconsistent version.

## K8s way

// very typical tech scenario, one of big reason I like k8s so much, my own story, last job
With K8s(Deployment object), we can achieve it by changing one field(image url), k8s will do the rest.

With k8s, to release:

*  Build a docker image of new version, push it to registry.
*  Change image url field of Deployment object

You get zero downtime release(also depend on your app)

## Rolling update

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201021102401.png" alt="20201021102401" style="width:1000px"/>

## Create a Deployment

.code resource/svc-deployment-v1.yaml /START OMIT/,/END OMIT/

.code resource/sample/deploy_v2.sh

TODO: may not need a demo for deployments at all