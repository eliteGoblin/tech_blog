#!/bin/bash
export SERVICE_ACCOUNT=api-explorer
export SECRET=$(kubectl get serviceaccount ${SERVICE_ACCOUNT} -o json -n demo| jq -Mr '.secrets[].name | select(contains("token"))')
export TOKEN=$(kubectl get secret ${SECRET} -o json -n demo| jq -Mr '.data.token' | base64 -d)
kubectl get secret ${SECRET} -o json -n demo| jq -Mr '.data["ca.crt"]' | base64 -d > /tmp/ca.crt
kubectl -n default get endpoints kubernetes --no-headers
export APISERVER=https://{One IP addr from above}:443


