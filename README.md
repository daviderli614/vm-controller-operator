
<img width="428" alt="image" src="https://user-images.githubusercontent.com/64472425/168420021-a8dc337e-788b-4748-8c68-dd331dbb9d40.png">

## CR
```
apiVersion: vmcontroller.uk8s.ucloud.cn/v1
kind: VMClaim
metadata:
  name: vm-controller
  namespace: kube-system
spec:
  image: {{VM_CONTROLLER_IMAGE}}
  resources:
    limits:
      cpu: {{CPU_LIMIT}}
      memory: {{MEM_LIMIT}}
    requests:
      cpu: {{CPU_REQUEST}}
      memory: {{MEM_REQUEST}}
  envs:
    - name: UCLOUD_UK8S_CLUSTER_ID
      value: "{{CLUSTER_ID}}"
    - name: VERSION
      value: "{{VERSION}}"
  scaleUp:
    scaleUpUtilizationThreshold: "{{SCALE_UP_UTILIZATION-THRESHOLD}}"
  scaleDown:
    enabled: {{SCALE_DOWN_ENABLE}}
    delayAfterAdd: {{SCALE_DOWN_DELAY_AFTER_ADD}}
    unneededTime: {{SCALE_DOWN_UNNEEDED_TIME}}
    scaleDownUtilizationThreshold: "{{SCALE_DOWN_UTILIZATION-THRESHOLD}}"
  balanceSimilarNodeGroups: {{BALANCE_SIMILAR_NODE_GROUPS}}
  skipNodesWithLocalStorage: {{SKIP_NODES_WITH_LOCAL_STORAGE}}
  expander: "{{EXPANDER}}"
  scanInterval: 30s
```

## deploy
```
# kubectl apply -f deploy/crd/
# kubectl apply -f deploy/vm-controller/
# kubectl apply -f deploy/operator/ 
```
