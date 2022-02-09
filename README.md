# pad-operator
A kubernetes operator to manage [prometheus-anomaly-detector](https://github.com/AICoE/prometheus-anomaly-detector) instances. Based on [Operator-SDK](https://sdk.operatorframework.io/). 

## Prometheus Anomaly Detector
Prometheus-Anomaly-Detector (PAD) is a machine learning framework that enables us to use various models such as Fourier, LSTM to perform time-series forecasting on metric data collected from a given prometheus source. PAD collects the metrics of interest from the specified Prometheus data source, trains a model to forecast the future values of these metrics. These predicted future values (upper and lower bounds when considering a margin of error) and the actual values of the metrics are compared. 

If the actual value of a particular metric at a point of time is very different from what was predicted (i.e it has either larger than the upper bound or smaller than the lower bound) then it is conidered to be an anomaly. 

![PAD Diagram](https://user-images.githubusercontent.com/7343099/64876403-081c9f80-d61d-11e9-84df-266c91a75dde.jpg)
