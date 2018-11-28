# Kubo Monitoring Using Prometheus & Grafana
## Deploy Prometheus
kubectl create namespace monitoring\
kubectl config set-context $(kubectl config current-context) --namespace=monitoring\
kubectl apply -f prometheus-rbac.yaml\
kubectl create -f prometheus-config-map.yaml\
kubectl create  -f prometheus-deployment.yaml\
kubectl get deployments\
kubectl create -f prometheus-service.yaml\
kubectl get svc |grep prometheus-service 

**Open web browser to** http://*prometheus-service EXTERNAL-IP*:8080\
**Select > Status > Targets**\
Verify Data Collection

## Deploy Grafana
Install helm client - https://docs.helm.sh/using_helm/#installing-helm or
curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh\
chmod 700 get_helm.sh\
./get_helm.sh\
kubectl apply -f tiller-rbac.yaml\
helm init --service-account tiller\
helm install --name grafana-app --namespace monitoring stable/grafana\
helm ls\
kubectl get secret --namespace monitoring grafana-app -o jsonpath="{.data.admin-password}" | base64 --decode 

Record Password\
kubectl create -f grafana-service.yaml\
kubectl get svc |grep grafana-service

**Open web browser to** http:*grafana-service EXTERNAL-IP*:3000

**User Name:** admin

**Pasword:** "Output from Above"

**Select > Add data source**

**Name:** Prometheus

**Type:** Prometheus

**URL:** http://*prometheus-service EXTERNAL-IP*:8080

**Select > Save and Test**

**Select > "+"** in left pane then **"Import"**

Import the .json file located in this repository.

In the Options table, select the Prometheus drop-down and select Prometheus 
**Select > Import** 

You will redirect to the new dashboard. 
