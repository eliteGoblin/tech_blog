
## Account

*  elitegoblinrb@gmail.com
*  Projects: k8s-talk

## Init

re-init

```
# 每次重新login, 会让你选择projects
gcloud init --skip-diagnostics
gcloud config set compute/region australia-southeast1
gcloud config set compute/zone australia-southeast1-a
```

重建cluster
```
gcloud container clusters create k8s-talks --num-nodes=3 --enable-autoscaling --min-nodes=3 --max-nodes=3
# 设置kubectl, gcloud赋予其credential
gcloud container clusters get-credentials k8s-talks
```

配置
```
gcloud container clusters resize k8s-talks --num-nodes=3
# delete
gcloud container clusters delete k8s-talks
```

## Kube-ops-view

```
kubectl proxy &
docker run -it --net=host hjacobs/kube-ops-view
```

## 限制与k8s-talks namespace

```
kubectl create namespace k8s-talks
kubectl config set-context --current --namespace=k8s-talks
// 删除, 重建namespace最快
kubectl delete all --all -n k8s-talks
```

## Demo pod

```
kubectl run --restart=Never --image=gcr.io/kuar-demo/kuard-amd64:blue kuard
```

主页会显示 pod name

有3个版本 tag

*  blue
*  green
*  purple