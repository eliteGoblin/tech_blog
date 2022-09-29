

Replacing pods with newer versions
 Updating managed pods
 Updating pods declaratively using Deployment
resources
 Performing rolling updates
 Automatically blocking rollouts of bad versions
 Controlling the rate of the rollout
 Reverting pods to a previous version


make your life easier by using a declarative
approach to deploying and updating applications in Kubernetes


Update: true zero-downtime update process

different version
```
luksa/kubia:v1
luksa/kubia:v2
```


```
while true; do curl http://130.211.109.222; done
```

A Deployment is a higher-level resource meant for deploying applications and
updating them declaratively,

When you create a Deployment, a ReplicaSet resource is created underneath
(eventually more of them)


the actual pods
are created and managed by the Deployment’s ReplicaSets, not by the Deployment
directly (the relationship is shown in figure 9.8)


<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201020104327.png" alt="20201020104327" style="width:500px"/>


create
```
kubectl create -f kubia-deployment-v1.yaml --record
```

Check rollout status

```
kubectl rollout status deployment kubia
```

update a Deployment. The only thing
you need to do is modify the pod template defined in the Deployment resource

Deployment strategy:

*  Rolling update(default strategy): removes old pods one by one,
while adding new ones at the same time
*  Recreate: deletes all the old pods at once and then creates new one; old pods to be deleted before the new ones are
created (when multiple version can't run same time). Not good and bad.

By changing the pod template in your Deployment resource, you’ve updated your app to a newer version—by changing a single
field!

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201020112547.png" alt="20201020112547" style="width:500px"/>

查看两个RS此消彼长的过程
```
kubectl get rs
```

he old ReplicaSet is still there, whereas the old ReplicationController was deleted at the end of the rolling-update process


What if deployments has problem

```
kubectl set image deployment kubia nodejs=luksa/kubia:v3
```

5次以后, 500

```
while true; do curl http://130.211.109.222; done
```

prod issue!

```
kubectl rollout undo deployment kubia
```

ROLLING BACK TO A SPECIFIC DEPLOYMENT REVISION

```
kubectl rollout undo deployment kubia --to-revision=1
```

 revision history is limited by the revisionHistoryLimit property; It defaults to two, so normally only the current and the previous revision
are shown in the history

不停的deploy不work的deployment, 对revision有何影响; revision只记录成功的?

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20201020162255.png" alt="20201020162255" style="width:500px"/>


can also pause and resume, e.g. : see how it behaves with only a fraction of all your
users
```
kubectl rollout pause deployment kubia
kubectl rollout resume deployment kubia
```

. Until the pod is available, the rollout process will not continue


no explicit readiness probe defined, the container and the
pod were always considered ready,


如果不设置minReadySeconds, If the readiness probe starts failing
shortly after, the bad version is rolled out across all pods.

新Version的pod status状态一直不ready, 会导致deployments超时 ProgressDeadlineExceeded; 超时时间配置: progressDeadlineSeconds, 