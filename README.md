# pad-operator
A kubernetes operator to manage [prometheus-anomaly-detector](https://github.com/AICoE/prometheus-anomaly-detector) instances. Based on [Operator-SDK](https://sdk.operatorframework.io/). 

## Prometheus Anomaly Detector
Prometheus-Anomaly-Detector (PAD) is a machine learning framework that enables us to use various models such as Fourier, LSTM to perform time-series forecasting on metric data collected from a given prometheus source. PAD collects the metrics of interest from the specified Prometheus data source, trains a model to forecast the future values of these metrics. These predicted future values (upper and lower bounds when considering a margin of error) and the actual values of the metrics are compared. 

![PAD Diagram](https://user-images.githubusercontent.com/7343099/64876301-d9062e00-d61c-11e9-80b6-35cb5c9e4540.jpg)

If the actual value of a particular metric at a point of time is very different from what was predicted (i.e it has either larger than the upper bound or smaller than the lower bound) then it is conidered to be an anomaly. 

![PAD Metric comparison](https://user-images.githubusercontent.com/7343099/64876403-081c9f80-d61d-11e9-84df-266c91a75dde.jpg)

## Using the PAD Operator

### Installing the operator
The operator can be introduced into a given kubernetes cluster in two ways. 

#### Deployment using the make file

1. Set the appropriate kubernetes context (in case you have multiple clusters).
2. Clone this repository.
3. Run `make deploy`. 

This should create a `pad-operator` Deployment in a new `pad-operator-system` namespace. 

`make undeploy` deletes the operator from the cluster. (It would be good to delete all the pad resources created via the operator first). 

#### Deployment using OLM

This repo also contains bundle files that enable the management of the operator via [OLM](https://sdk.operatorframework.io/docs/olm-integration/tutorial-bundle/). Please follow the steps mentioned in [enabling olm](https://sdk.operatorframework.io/docs/olm-integration/tutorial-bundle/#enabling-olm) if your cluster does not have olm enabled. 

The bundle image for this operator can be found in [docker hub](https://hub.docker.com/repository/docker/arjunshenoymec/pad-operator-bundle)

Run `operator-sdk run bundle docker.io/arjunshenoymec/pad-operator-bundle:v0.0.1` 

Run `operator-sdk cleanup --delete-all pad-operator` to delete the operator and all related resources. 

### PAD CRDs

After getting the pad-operator up and running in your cluster, you can manage PAD instances by using the pad CustomResource.

```
apiVersion: indicator.padoperator/v1alpha1
kind: Pad
metadata:
  name: pad-sample
spec:
   replicas: 1
   source: "http://demo.robustperception.io:9090/"
   metrics: "up"
   retraining_interval: "10"
```

The above snippet if applied will create a Deployment which will access [demo.robustperception.io](http://demo.robustperception.io:9090/), collect the `up` metric(s) and perform the forecasting, anomaly detection process. The following table specifies the currently available PadSpec parameters, what they mean and their default values. 

Parameter | Definition | Default 
----- | ----- | -----
replicas | The number of Replicas in the deployment | 1
source | The URL corresponding to the prometheus datasource. The port is also to be included. Corresponds to `FLT_PROM_URL` in the PAD source code. | "http://demo.robustperception.io:9090/"
metrics | The list of metrics to be worked on. This is a string where each metric is to be separated by a `;`. Corresponds to `FLT_METRICS_LIST` in the pad repo. | "up"
retraining_interval | Specifies how often the model will be retrained. Corresponds to `FLT_RETRAINING_INTERVAL_MINUTES` in the PAD repo. | "15"
training_window_size | Limits the size of the data considered for training. Also deletes older than the training_window during each trainign iteration. Corresponds to `FLT_ROLLING_TRAINING_WINDOW_SIZE` in the PAD Repo. | "24h" 
image | The PAD container image to be used. In case you want to modify and use your own private container image | "quay.io/aicoe/prometheus-anomaly-detector:latest"

## Work(s) In Progress

We plan to introduce other parameters such as ones corresponding to `FLT_PARALLELISM`, modifying the number of CPUs available to the container and enabling the choice of model being trained (currently an old version of [Prophet](https://facebook.github.io/prophet/) is used as default in the default image specified above). 
