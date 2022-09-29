kubectl apply -f resource/svc-deployment-v1.yaml
kubectl get rs # see rs rotate
while true; sleep 1 && do curl http://35.201.21.74; done # see version change in curl