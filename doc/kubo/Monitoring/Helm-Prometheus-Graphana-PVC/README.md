# Deploy Kubo Monitoring Using Helm (Prometheus & Grafana) along with Presistance Volume

**Note**: Please Check below option file used that allow-privileged-containers.yml
```
bosh -d cfcr deploy /home/ubuntu/workspace/v0.21.0/kubo-deployment/manifests/cfcr.yml \
   -o /home/ubuntu/workspace/v0.21.0/kubo-deployment/manifests/ops-files/iaas/openstack/cloud-provider.yml \
   -o /home/ubuntu/workspace/v0.21.0/kubo-deployment/manifests/ops-files/allow-privileged-containers.yml \
   -v auth_url=http://182.252.135.131:5000/v3 \
   -v openstack_domain=default \
   -v openstack_username=crossent \
   -v openstack_password=crossent \
   -v region=RegionOne \
   -v openstack_project_id=7cd22e31952b478a8788f3f9a74c7a68 \
   -v ignore-volume-az=false
```

**Note**: Please Check if you have Storage Class According to your IaaS.

If you dont have Storage Class, [Create Storage Class According to your IaaS](https://kubernetes.io/docs/concepts/storage/storage-classes/) 

**Note**: Create Storage Class only on Openstack

```
For Example:
---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: gold  
provisioner: kubernetes.io/cinder
parameters:
  availability: nova
  type: KUBO
```


```
$ kubectl get storageclass

NAME      PROVISIONER            AGE
gold      kubernetes.io/cinder   10d
```

## 1) Setup Helm 
    Install helm client - https://docs.helm.sh/using_helm/#installing-helm or 

    curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh 

    chmod 700 get_helm.sh 

    ./get_helm.sh 

    kubectl apply -f tiller-rbac.yaml 

    helm init --service-account tiller 

**Note**: Please Check triller deployment before using helm.


```
$ kubectl get deployment tiller-deploy -n kube-system

NAME            DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
tiller-deploy   1         1         1            1           46s
```
## 2) Install Ingress with NGINX controller

      helm search nginx-ingress
          NAME                	CHART VERSION	APP VERSION	DESCRIPTION                                       
         stable/nginx-ingress	 0.25.1       	0.17.1     	An nginx Ingress controller that uses ConfigMap...

      helm install stable/nginx-ingress --name nginx-ingress --version 0.25.1 --set controller.publishService.enabled=true

**Note**: Please Check Nginx deployment and Service before deploying Prometheus.

```
$ kubectl get deployments -n default

NAME                            DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-ingress-controller        1         1         1            1           7m
nginx-ingress-default-backend   1         1         1            1           7m
```      
kubectl get services -n default
```
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP       PORT(S)                      AGE
kubernetes                      ClusterIP      10.100.200.1     <none>            443/TCP                      1d
nginx-ingress-controller        LoadBalancer   10.100.200.146   xxx.xxx.135.xxx   80:31763/TCP,443:32673/TCP   1h
nginx-ingress-default-backend   ClusterIP      10.100.200.189   <none>            80/TCP                       1h
```
## 3) Install Prometheus-Operator

    helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/
    helm repo update
    helm install --name prometheus-operator --namespace monitoring coreos/prometheus-operator

**Note**: Please Check Prometheus-Operator before deploying kube-prometheus.


```
$ kubectl get deployments -n monitoring

NAME                  DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
prometheus-operator   1         1         1            1           3m
```
**Note**: When you install the Prometheus operator you will get the new Custom Resource Definitions or CRDs


```
$ kubectl get CustomResourceDefinition

NAME                                    AGE
alertmanagers.monitoring.coreos.com     3m
prometheuses.monitoring.coreos.com      3m
prometheusrules.monitoring.coreos.com   3m
servicemonitors.monitoring.coreos.com   3m
```
## 4) Install kube-prometheus
    helm install --name mon --namespace monitoring -f custom-values.yaml coreos/kube-prometheus
**Note**: Please Check kube-prometheus deployment,Service, PresesitanceVolumeClaim, Ingress before deploying Prometheus.

```
$ kubectl get pods -n monitoring

NAME                                       READY     STATUS    RESTARTS   AGE
alertmanager-mon-0                         2/2       Running   0          3m
mon-exporter-kube-state-5c7c464686-2pkhw   2/2       Running   0          3m
mon-exporter-node-nrtmn                    1/1       Running   0          3m
mon-exporter-node-rzvk5                    1/1       Running   0          3m
mon-exporter-node-tn59k                    1/1       Running   0          3m
mon-grafana-74889c77fb-9r5tx               2/2       Running   0          3m
prometheus-mon-prometheus-0                3/3       Running   1          3m
prometheus-operator-858c485-v4l8j          1/1       Running   0          26m
```

```
$ kubectl get services -n monitoring

NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP       PORT(S)             AGE
alertmanager-operated     ClusterIP      None             <none>            9093/TCP,6783/TCP   5m
mon-alertmanager          ClusterIP      10.100.200.27    <none>            9093/TCP            5m
mon-exporter-kube-state   ClusterIP      10.100.200.124   <none>            80/TCP              5m
mon-exporter-node         ClusterIP      10.100.200.231   <none>            9100/TCP            5m
mon-grafana               ClusterIP      10.100.200.123   <none>            80/TCP              5m
mon-prometheus            LoadBalancer   10.100.200.105   xxx.xxx.135.xxx   9090:30921/TCP      5m
prometheus-operated       ClusterIP      None             <none>            9090/TCP            5m
```
**Open web browser to** http://*mon-prometheus-EXTERNAL-IP*:9090


```
$ kubectl get pvc -n monitoring

NAME                                                       STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
alertmanager-mon-db-alertmanager-mon-0                     Bound     pvc-7c2129a7-a4e6-11e8-9bf4-fa163ea5164b   10Gi       RWO            gold           9m
mon-grafana                                                Bound     pvc-7af0a9c2-a4e6-11e8-9bf4-fa163ea5164b   10Gi       RWO            gold           9m
prometheus-mon-prometheus-db-prometheus-mon-prometheus-0   Bound     pvc-7cdc57af-a4e6-11e8-9bf4-fa163ea5164b   10Gi       RWO            gold           9m
```

```
$ kubectl get ingress -n monitoring

NAME          HOSTS                       ADDRESS           PORTS     AGE
mon-grafana   grafana.test.akomljen.com   xxx.xxx.135.xxx   80, 443   18m
```
## 5) Configure Graphana Web Interface

Now add grafana.test.akomljen.com to the list of your /etc/hosts:

    $ {mon-grafana-EXTERNAL_IP} grafana.test.akomljen.com

Now go to grafana.test.akomljen.com here you are!

**Open web browser to** http://grafana.test.akomljen.com

**User Name:** admin

**Pasword:** "Obtain Password from above custom-values.yaml > grafana > adminPassword"

**Select > Configuration > Data Sources > Prometheus**
**Edit > HTTP > URL**

**URL:** http://*mon-prometheus-EXTERNAL-IP*:9090

**Select > Save and Test**

You will redirect to the new dashboard. 
