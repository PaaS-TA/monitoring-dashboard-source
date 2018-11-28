# Deploy Kubo Logging Elasticseach-Fluentd-Kibana (EFK) along with Presistance Volume

## 1) Setup Helm 
    Install helm client - https://docs.helm.sh/using_helm/#installing-helm or

    curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh

    chmod 700 get_helm.sh

    ./get_helm.sh

    kubectl apply -f tiller-rbac.yaml

    helm init --service-account tiller

**Note**: Please Check triller deployment before using helm.

kubectl get deployment tiller-deploy -n kube-system
```
NAME            DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
tiller-deploy   1         1         1            1           46s
```
**Note**: Please Check if you have Storage Class According to your IaaS.

If you dont have Storage Class, [Create Storage Class According to your IaaS](https://kubernetes.io/docs/concepts/storage/storage-classes/) 

kubectl get storageclass
```
NAME      PROVISIONER            AGE
gold      kubernetes.io/cinder   10d
```

## 2) Install Elasticsearch store and analyze your application Log data

      kubectl create -f elasticsearch.yaml -n kube-system

**Note**: Please Check Elasticsearch StatefulSets, Deployment, Service and Presestance Volumes before deploying FluentD.

kubectl get statefulsets,pods,pvc,svc,deployments,configmaps -n kube-system
```
NAME                                        DESIRED   CURRENT   AGE
statefulsets/logging-elasticsearch-data     3         3         13m
statefulsets/logging-elasticsearch-master   3         3         13m

NAME                                               READY     STATUS    RESTARTS   AGE
po/logging-elasticsearch-client-6557d855df-gn9q5   1/1       Running   1          34m
po/logging-elasticsearch-client-6557d855df-r294g   1/1       Running   1          34m
po/logging-elasticsearch-data-0                    1/1       Running   0          34m
po/logging-elasticsearch-data-1                    1/1       Running   0          33m
po/logging-elasticsearch-data-2                    1/1       Running   0          32m
po/logging-elasticsearch-master-0                  1/1       Running   0          34m
po/logging-elasticsearch-master-1                  1/1       Running   0          33m
po/logging-elasticsearch-master-2                  1/1       Running   0          32m

NAME                                      STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
pvc/data-logging-elasticsearch-data-0     Bound     pvc-e28d127c-acf0-11e8-8c13-fa163ea1f9ba   30Gi       RWO            gold           13m
pvc/data-logging-elasticsearch-data-1     Bound     pvc-1c93feba-acf1-11e8-8c13-fa163ea1f9ba   30Gi       RWO            gold           11m
pvc/data-logging-elasticsearch-data-2     Bound     pvc-41dd2ad4-acf1-11e8-8c13-fa163ea1f9ba   30Gi       RWO            gold           10m
pvc/data-logging-elasticsearch-master-0   Bound     pvc-e2a717f5-acf0-11e8-8c13-fa163ea1f9ba   4Gi        RWO            gold           13m
pvc/data-logging-elasticsearch-master-1   Bound     pvc-18b3af95-acf1-11e8-8c13-fa163ea1f9ba   4Gi        RWO            gold           12m
pvc/data-logging-elasticsearch-master-2   Bound     pvc-43868c5a-acf1-11e8-8c13-fa163ea1f9ba   4Gi        RWO            gold           10m

NAME                                       TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)               AGE
svc/logging-elasticsearch-client           ClusterIP   10.100.200.3     <none>        9200/TCP              13m
svc/logging-elasticsearch-discovery        ClusterIP   None             <none>        9300/TCP              13m

NAME                                  DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/logging-elasticsearch-client   2         2         2            2           13m

NAME                                    DATA      AGE
cm/logging-elasticsearch                4         13m
```      
## Access the elasticserach service
From any worker-node on the Kubernetes cluster (that's running kube-proxy or similar), run:
```
curl http://logging-elasticsearch-client-CLUSTER-IP:9200
```
One should see something similar to the following:
```
{
  "name" : "logging-elasticsearch-client-6557d855df-r294g",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "w-psA2dXTQKI5OHzKOpbxA",
  "version" : {
    "number" : "6.4.0",
    "build_flavor" : "oss",
    "build_type" : "tar",
    "build_hash" : "595516e",
    "build_date" : "2018-08-17T23:18:47.308994Z",
    "build_snapshot" : false,
    "lucene_version" : "7.4.0",
    "minimum_wire_compatibility_version" : "5.6.0",
    "minimum_index_compatibility_version" : "5.0.0"
  },
  "tagline" : "You Know, for Search"
}
```
Or if one wants to see cluster information:
```
curl http://logging-elasticsearch-client-CLUSTER-IP:9200/_cluster/health?pretty
```
```
{
  "cluster_name" : "elasticsearch",
  "status" : "green",
  "timed_out" : false,
  "number_of_nodes" : 8,
  "number_of_data_nodes" : 3,
  "active_primary_shards" : 5,
  "active_shards" : 10,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 0,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 100.0
}
```
## 3) Install fluentd to forward your application Log data to Elasticsearch

    kubectl create -f fluentd.yaml -n kube-system

**Note**: Please Check Fluentd configmaps,pods,serviceaccounts,ClusterRole,ClusterRoleBinding,daemonsets before deploying Kibana.

kubectl get configmaps,pods,serviceaccounts,ClusterRole,ClusterRoleBinding,daemonsets -n kube-system
```
NAME                                    DATA      AGE
cm/fluentd-fluentd-elasticsearch        6         22m

NAME                                               READY     STATUS    RESTARTS   AGE
po/fluentd-fluentd-elasticsearch-c22rx             1/1       Running   0          22m
po/fluentd-fluentd-elasticsearch-wlm68             1/1       Running   0          22m
po/fluentd-fluentd-elasticsearch-z64p4             1/1       Running   0          22m

NAME                                    SECRETS   AGE
sa/fluentd-fluentd-elasticsearch        1         22m

NAME                                                                                AGE
clusterroles/fluentd-fluentd-elasticsearch                                          22m

NAME                                                                       AGE
clusterrolebindings/fluentd-fluentd-elasticsearch                          22m

NAME                               DESIRED   CURRENT   READY     UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
ds/fluentd-fluentd-elasticsearch   3         3         3         3            3           <none>          22m
```      
## Access the elasticserach service
From any worker-node on the Kubernetes cluster (that's running kube-proxy or similar), run:
```
curl http://logging-elasticsearch-client-CLUSTER-IP:9200/_cat/indices?v&pretty
```
One should see something similar to the following:
```
health status index               uuid                   pri rep docs.count docs.deleted store.size pri.store.size
green  open   logstash-2018.08.31 o-2uxO6UTluSop3NLhnR1g   5   1      60348            0     66.1mb         38.7mb
```
## 4) Install Kibana to visualize you log data 
    helm install stable/kibana --name kibana --set env.ELASTICSEARCH_URL=http://logging-elasticsearch-client:9200,env.SERVER_BASEPATH=/api/v1/namespaces/kube-system/services/kibana/proxy --namespace=kube-system --version 0.13.0

**Note**: Please Check kibana deployment,pods,service,configmaps, before deploying Prometheus.

kubectl get deployments,pods,services,configmaps -n kube-system
```
NAME                                  DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/kibana                         1         1         1            1           4m

NAME                                               READY     STATUS    RESTARTS   AGE
po/kibana-5b8c44f6f6-p2t8s                         1/1       Running   0          4m

NAME                                       TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)               AGE
svc/kibana                                 ClusterIP   10.100.200.139   <none>        443/TCP               4m

NAME                                    DATA      AGE
cm/kibana                               1         4m
cm/kibana.v1                            1         4m

```
## Run your Kubernetes Cluster Proxy on your Local Notebook (or) pc
```
kubectl proxy
```
**Open web browser to** http://localhost:8001/api/v1/namespaces/kube-system/services/kibana/proxy/app/kibana

**Select > Management > Index Patterns**

**Create index pattern**

**step1 > logstash-* > Next Step**

**step2 > @timestamp > Create index pattern**

**Select > Discover**

You will redirect to the Kibana new dashboard. 
