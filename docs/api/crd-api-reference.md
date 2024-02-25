# API Reference

## Packages
- [windtunnel.plantd.org/v1alpha1](#windtunnelplantdorgv1alpha1)


## windtunnel.plantd.org/v1alpha1

Package v1alpha1 contains API Schema definitions for the windtunnel v1alpha1 API group

### Resource Types
- [CostExporter](#costexporter)
- [CostExporterList](#costexporterlist)
- [DataSet](#dataset)
- [DataSetList](#datasetlist)
- [DigitalTwin](#digitaltwin)
- [DigitalTwinList](#digitaltwinlist)
- [Experiment](#experiment)
- [ExperimentList](#experimentlist)
- [LoadPattern](#loadpattern)
- [LoadPatternList](#loadpatternlist)
- [Pipeline](#pipeline)
- [PipelineList](#pipelinelist)
- [PlantDCore](#plantdcore)
- [PlantDCoreList](#plantdcorelist)
- [Schema](#schema)
- [SchemaList](#schemalist)
- [Simulation](#simulation)
- [SimulationList](#simulationlist)
- [TrafficModel](#trafficmodel)
- [TrafficModelList](#trafficmodellist)



#### Column



Column defines the metadata of the column data.

_Appears in:_
- [SchemaSpec](#schemaspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name defines the name of the column. |
| `type` _string_ | Type defines the data type of the column. Should match the type with one of the provided types. |
| `params` _object (keys:string, values:string)_ | Params defines the parameters for constructing the data give certain data type. |
| `formula` _[FormulaSpec](#formulaspec)_ | Formula defines the formula applies to the column data. |


#### CostExporter



CostExporter is the Schema for the costexporters API

_Appears in:_
- [CostExporterList](#costexporterlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `CostExporter`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[CostExporterSpec](#costexporterspec)_ | Spec defines the specifications of the CostExporter. |


#### CostExporterList



CostExporterList contains a list of CostExporter



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `CostExporterList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[CostExporter](#costexporter) array_ | Items defines a list of CostExporters. |


#### CostExporterSpec



CostExporterSpec defines the desired state of CostExporter

_Appears in:_
- [CostExporter](#costexporter)

| Field | Description |
| --- | --- |
| `s3Bucket` _string_ | S3Bucket defines the AWS S3 bucket name where stores the cost logs. |
| `cloudServiceProvider` _string_ | CloudServiceProvider defines the target cloud service provide for calculating cost. |
| `secretRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | SecretRef defines the reference to the Kubernetes Secret where stores the credentials of cloud service provider |




#### DataSet



DataSet is the Schema for the datasets API

_Appears in:_
- [DataSetList](#datasetlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DataSet`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DataSetSpec](#datasetspec)_ | Spec defines the specifications of the DataSet. |


#### DataSetList



DataSetList contains a list of DataSet



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DataSetList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[DataSet](#dataset) array_ | Items defines a list of DataSets. |


#### DataSetSpec



DataSetSpec defines the desired state of DataSet

_Appears in:_
- [DataSet](#dataset)

| Field | Description |
| --- | --- |
| `fileFormat` _string_ | FileFormat defines the file format of the each file containing the generated data. This may or may not be the output file format based on whether you want to compress these files. |
| `compressedFileFormat` _string_ | CompressedFileFormat defines the file format for the compressed files. Each file inside the compressed file is of "fileFormat" format specified above. This is the output format if specified for the files. |
| `compressPerSchema` _boolean_ | CompressPerSchema defines the flag of compression. If you wish files from all the different schemas to compressed into one compressed file leave this field as false. If you wish to have a different compressed file for every schema, mark this field as true. |
| `numFiles` _integer_ | NumberOfFiles defines the total number of output files irrespective of compression. Unless "compressPerSchema" is false, this field is applicable per schema. |
| `schemas` _[SchemaSelector](#schemaselector) array_ | Schemas defines a list of Schemas. |
| `parallelJobs` _integer_ | ParallelJobs defines the number of parallel jobs when generating the dataset. |




#### DeploymentConfig



DeploymentConfig defines the desired state of modules managed as Deployment

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `image` _string_ | Image defines the container image to use |
| `replicas` _integer_ | Replicas defines the desired number of replicas |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#resourcerequirements-v1-core)_ | Resources defines the resource requirements per replica |


#### DigitalTwin



DigitalTwin is the Schema for the digitaltwins API

_Appears in:_
- [DigitalTwinList](#digitaltwinlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DigitalTwin`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DigitalTwinSpec](#digitaltwinspec)_ | Spec defines the specifications of the DigitalTwin. |


#### DigitalTwinList



DigitalTwinList contains a list of DigitalTwin



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DigitalTwinList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[DigitalTwin](#digitaltwin) array_ | Items defines a list of DigitalTwins. |


#### DigitalTwinSpec



DigitalTwinSpec defines the desired state of DigitalTwin

_Appears in:_
- [DigitalTwin](#digitaltwin)

| Field | Description |
| --- | --- |
| `modelType` _string_ | ModelType defines the type of the DigitalTwin model. |
| `experiments` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core) array_ | Experiments contains the list of Experiment object references for the DigitalTwin. |




#### Endpoint



Endpoint defines the configuration of the endpoint.

_Appears in:_
- [PipelineSpec](#pipelinespec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name defines the name of the endpoint. It's required when it's for pipeline endpoint. |
| `http` _[HTTP](#http)_ | HTTP defines the configuration of the HTTP request. It's mutually exclusive with WebSocket and GRPC. |
| `websocket` _[WebSocket](#websocket)_ | WebSocket defines the configuration of the WebSocket connection. It's mutually exclusive with HTTP and GRPC. |
| `grpc` _[GRPC](#grpc)_ | GRPC defines the configuration of the gRPC request. It's mutually exclusive with HTTP and WebSocket. |
| `serviceRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | ServiceRef defines the Kubernetes Service that exposes metrics. |
| `port` _string_ | Name of the Service port which this endpoint refers to. <br /><br /> It takes precedence over `targetPort`. |
| `targetPort` _[IntOrString](https://pkg.go.dev/k8s.io/apimachinery/pkg/util/intstr#IntOrString)_ | Name or number of the target port of the `Pod` object behind the Service, the port must be specified with container port property. <br /><br /> Deprecated: use `port` instead. |
| `path` _string_ | HTTP path from which to scrape for metrics. <br /><br /> If empty, Prometheus uses the default value (e.g. `/metrics`). |
| `scheme` _string_ | HTTP scheme to use for scraping. <br /><br /> `http` and `https` are the expected values unless you rewrite the `__scheme__` label via relabeling. <br /><br /> If empty, Prometheus uses the default value `http`. |
| `params` _object (keys:string, values:string array)_ | params define optional HTTP URL parameters. |
| `interval` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | Interval at which Prometheus scrapes the metrics from the target. <br /><br /> If empty, Prometheus uses the global scrape interval. |
| `scrapeTimeout` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | Timeout after which Prometheus considers the scrape to be failed. <br /><br /> If empty, Prometheus uses the global scrape timeout unless it is less than the target's scrape interval value in which the latter is used. |
| `tlsConfig` _[TLSConfig](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#TLSConfig)_ | TLS configuration to use when scraping the target. |
| `bearerTokenFile` _string_ | File to read bearer token for scraping the target. <br /><br /> Deprecated: use `authorization` instead. |
| `bearerTokenSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#secretkeyselector-v1-core)_ | `bearerTokenSecret` specifies a key of a Secret containing the bearer token for scraping targets. The secret needs to be in the same namespace as the ServiceMonitor object and readable by the Prometheus Operator. <br /><br /> Deprecated: use `authorization` instead. |
| `authorization` _[SafeAuthorization](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#SafeAuthorization)_ | `authorization` configures the Authorization header credentials to use when scraping the target. <br /><br /> Cannot be set at the same time as `basicAuth`, or `oauth2`. |
| `honorLabels` _boolean_ | When true, `honorLabels` preserves the metric's labels when they collide with the target's labels. |
| `honorTimestamps` _boolean_ | `honorTimestamps` controls whether Prometheus preserves the timestamps when exposed by the target. |
| `trackTimestampsStaleness` _boolean_ | `trackTimestampsStaleness` defines whether Prometheus tracks staleness of the metrics that have an explicit timestamp present in scraped data. Has no effect if `honorTimestamps` is false. <br /><br /> It requires Prometheus >= v2.48.0. |
| `basicAuth` _[BasicAuth](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#BasicAuth)_ | `basicAuth` configures the Basic Authentication credentials to use when scraping the target. <br /><br /> Cannot be set at the same time as `authorization`, or `oauth2`. |
| `oauth2` _[OAuth2](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#OAuth2)_ | `oauth2` configures the OAuth2 settings to use when scraping the target. <br /><br /> It requires Prometheus >= 2.27.0. <br /><br /> Cannot be set at the same time as `authorization`, or `basicAuth`. |
| `metricRelabelings` _[RelabelConfig](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#RelabelConfig) array_ | `metricRelabelings` configures the relabeling rules to apply to the samples before ingestion. |
| `relabelings` _[RelabelConfig](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#RelabelConfig) array_ | `relabelings` configures the relabeling rules to apply the target's metadata labels. <br /><br /> The Operator automatically adds relabelings for a few standard Kubernetes fields. <br /><br /> The original scrape job's name is available via the `__tmp_prometheus_job_name` label. <br /><br /> More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config |
| `proxyUrl` _string_ | `proxyURL` configures the HTTP Proxy URL (e.g. "http://proxyserver:2195") to go through when scraping the target. |
| `followRedirects` _boolean_ | `followRedirects` defines whether the scrape requests should follow HTTP 3xx redirects. |
| `enableHttp2` _boolean_ | `enableHttp2` can be used to disable HTTP2 when scraping the target. |
| `filterRunning` _boolean_ | When true, the pods which are not running (e.g. either in Failed or Succeeded state) are dropped during the target discovery. <br /><br /> If unset, the filtering is enabled. <br /><br /> More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase |


#### Experiment



Experiment is the Schema for the experiments API

_Appears in:_
- [ExperimentList](#experimentlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Experiment`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ExperimentSpec](#experimentspec)_ | Spec defines the specifications of the Experiment. |


#### ExperimentList



ExperimentList contains a list of Experiments.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `ExperimentList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Experiment](#experiment) array_ | Items defines a list of Experiments. |


#### ExperimentSpec



ExperimentSpec defines the desired state of Experiment

_Appears in:_
- [Experiment](#experiment)

| Field | Description |
| --- | --- |
| `pipelineRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | PipelineRef defines s reference of the Pipeline object. |
| `loadPatterns` _[LoadPatternConfig](#loadpatternconfig) array_ | LoadPatterns defines a list of configuration of name of endpoints and LoadPatterns. |
| `scheduledTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | ScheduledTime defines the scheduled time for the Experiment. |




#### ExtraMetrics



ExtraMetrics defines the configurations of getting extra metrics.

_Appears in:_
- [PipelineSpec](#pipelinespec)

| Field | Description |
| --- | --- |
| `system` _[SystemMetrics](#systemmetrics)_ | System defines the configurfation of getting system metrics. |
| `messageQueue` _[MessageQueueMetrics](#messagequeuemetrics)_ | MessageQueue defines the configurfation of getting message queue related metrics. |


#### FormulaSpec



FormulaSpec defines the specification of the formula.

_Appears in:_
- [Column](#column)

| Field | Description |
| --- | --- |
| `name` _string_ | Name defines the name of the formula. Should match the name with one of the provided formulas. |
| `args` _string array_ | Args defines the arugments for calling the formula. |


#### GRPC



GRPC defines the configurations of gRPC protocol.

_Appears in:_
- [Endpoint](#endpoint)

| Field | Description |
| --- | --- |
| `address` _string_ | Placeholder. |
| `protoFiles` _string array_ | Placeholder. |
| `url` _string_ | Placeholder. |
| `params` _object (keys:string, values:string)_ | Placeholder. |
| `request` _object (keys:string, values:string)_ | Placeholder. |


#### HTTP



HTTP defines the configurations of HTTP protocol.

_Appears in:_
- [Endpoint](#endpoint)

| Field | Description |
| --- | --- |
| `url` _string_ | URL defines the absolute path for an entry point of the Pipeline. |
| `method` _string_ | Method defines the HTTP method used for the endpoint. |
| `headers` _object (keys:string, values:string)_ | Headers defines a map of HTTP headers. |
| `body` _[HTTPBodySpec](#httpbodyspec)_ | Body defines the configurations of the HTTP request body. |


#### HTTPBodySpec



HTTPBodySpec defines the configurations of the HTTP request body. User can specify either Data or DataSetRef, but not both fields.

_Appears in:_
- [HTTP](#http)

| Field | Description |
| --- | --- |
| `data` _string_ |  |
| `dataSetRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ |  |


#### LoadPattern



LoadPattern is the Schema for the loadpatterns API

_Appears in:_
- [LoadPatternList](#loadpatternlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `LoadPattern`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[LoadPatternSpec](#loadpatternspec)_ | Spec defines the specification of the LoadPattern. |


#### LoadPatternConfig



LoadPatternConfig defines the configuration of the load pattern in the experiment.

_Appears in:_
- [ExperimentSpec](#experimentspec)

| Field | Description |
| --- | --- |
| `endpointName` _string_ | EndpointName defines the name of endpoint where to send the requests. It should match the name of endpoint declared in the specification of the pipeline. |
| `loadPatternRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | LoadPatternRef defines s reference of the LoadPattern object. |


#### LoadPatternList



LoadPatternList contains a list of LoadPattern



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `LoadPatternList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[LoadPattern](#loadpattern) array_ | Items defines a list of LoadPatterns. |


#### LoadPatternSpec



LoadPatternSpec defines the desired state of LoadPattern

_Appears in:_
- [LoadPattern](#loadpattern)

| Field | Description |
| --- | --- |
| `stages` _[Stage](#stage) array_ | Stages defines a list of stages for the LoadPattern. |
| `preAllocatedVUs` _integer_ | PreAllocatedVUs defines pre-allocated virtual users for the K6 load generator. |
| `startRate` _integer_ | StartRate defines the initial requests per second when the K6 load generator starts. |
| `maxVUs` _integer_ | MaxVUs defines the maximum virtual users for the K6 load generator. |
| `timeUnit` _string_ | TimeUnit defines the unit of the time for K6 load generator. |




#### MessageQueueMetrics



MessageQueueMetrics defines the configurations of getting message queue related metrics.

_Appears in:_
- [ExtraMetrics](#extrametrics)



#### Pipeline



Pipeline is the Schema for the pipelines API

_Appears in:_
- [PipelineList](#pipelinelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Pipeline`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[PipelineSpec](#pipelinespec)_ | Spec defines the specifications of the Pipeline. |


#### PipelineList



PipelineList contains a list of Pipeline



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PipelineList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Pipeline](#pipeline) array_ | Items defines a list of Pipelines. |


#### PipelineSpec



PipelineSpec defines the desired state of Pipeline

_Appears in:_
- [Pipeline](#pipeline)

| Field | Description |
| --- | --- |
| `pipelineEndpoints` _[Endpoint](#endpoint) array_ | Endpoints for pipeline-under-test. |
| `healthCheckEndpoints` _string array_ | Endpoints for health check. |
| `metricsEndpoint` _[Endpoint](#endpoint)_ | Endpoints for metrics. |
| `extraMetrics` _[ExtraMetrics](#extrametrics)_ | Extra metrics, such as CPU utilzation, I/O and etc. |
| `inCluster` _boolean_ | In cluster flag. True indecates the pipeline-under-test is deployed in the same cluster as the plantD. Otherwise it should be False. |
| `cloudVendor` _string_ | State which cloud service provider the pipeline is deployed. |
| `enableCostCalculation` _boolean_ | Cost calculation flag. |
| `experimentRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | Internal usage. For experiment object to lock the pipeline object. |




#### PlantDCore



PlantDCore is the Schema for the plantdcores API

_Appears in:_
- [PlantDCoreList](#plantdcorelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PlantDCore`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[PlantDCoreSpec](#plantdcorespec)_ | Spec defines the specifications of the PlantDCore. |


#### PlantDCoreList



PlantDCoreList contains a list of PlantDCore



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PlantDCoreList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[PlantDCore](#plantdcore) array_ | Items defines a list of PlantDCores. |


#### PlantDCoreSpec



PlantDCoreSpec defines the desired state of PlantDCore

_Appears in:_
- [PlantDCore](#plantdcore)

| Field | Description |
| --- | --- |
| `kubeProxy` _[DeploymentConfig](#deploymentconfig)_ | KubeProxyConfig defines the desire state of PlantD Kube Proxy |
| `studio` _[DeploymentConfig](#deploymentconfig)_ | StudioConfig defines the desire state of PlantD Studio |
| `prometheus` _[PrometheusConfig](#prometheusconfig)_ | PrometheusConfig defines the desire state of Prometheus |
| `redis` _[DeploymentConfig](#deploymentconfig)_ | RedisConfig defines the desire state of Redis |




#### PrometheusConfig



PrometheusConfig defines the desired state of Prometheus

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `scrapeInterval` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | ScrapeInterval defines the desired time length between scrapings |
| `replicas` _integer_ | Replicas defines the desired number of replicas |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#resourcerequirements-v1-core)_ | Resources defines the resource requirements per replica |


#### Schema



Schema is the Schema for the schemas API

_Appears in:_
- [SchemaList](#schemalist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Schema`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[SchemaSpec](#schemaspec)_ | Spec defines the specifications of the Schema. |


#### SchemaList



SchemaList contains a list of Schema



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `SchemaList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Schema](#schema) array_ | Items defines a list of Schemas. |


#### SchemaSelector



SchemaSelector defines a list of Schemas and the required numbers and format.

_Appears in:_
- [DataSetSpec](#datasetspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name defines the name of the Schame. Should match the name of existing Schema in the same namespace as the DataSet. |
| `numRecords` _object (keys:string, values:integer)_ | NumRecords defines the number of records to be generated in each output file. A random number is picked from the specified range. |
| `numFilesPerCompressedFile` _object (keys:string, values:integer)_ | NumberOfFilesPerCompressedFile defines the number of intermediate files to be compressed into a single compressed file. A random number is picked from the specified range. |


#### SchemaSpec



SchemaSpec defines the desired state of Schema

_Appears in:_
- [Schema](#schema)

| Field | Description |
| --- | --- |
| `columns` _[Column](#column) array_ | Columns defines a list of column specifications. |




#### Simulation



Simulation is the Schema for the simulations API

_Appears in:_
- [SimulationList](#simulationlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Simulation`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[SimulationSpec](#simulationspec)_ | Spec defines the specifications of the Simulation. |


#### SimulationList



SimulationList contains a list of Simulation



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `SimulationList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Simulation](#simulation) array_ | Items defines a list of Simulations. |


#### SimulationSpec



SimulationSpec defines the desired state of Simulation

_Appears in:_
- [Simulation](#simulation)

| Field | Description |
| --- | --- |
| `trafficModelRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | TrafficModelRef defines the TrafficModel object reference for the Simulation. |
| `digitalTwinRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | DigitalTwinRef defines the DigitalTwin object reference for the Simulation. |




#### Stage



Stage defines the stage configuration of the load.

_Appears in:_
- [LoadPatternSpec](#loadpatternspec)

| Field | Description |
| --- | --- |
| `target` _integer_ | Target defines the target requests per second. |
| `duration` _string_ | Duration defines the duration of the current stage. |


#### SystemMetrics



SystemMetrics defines the configurations of getting system metrics.

_Appears in:_
- [ExtraMetrics](#extrametrics)

| Field | Description |
| --- | --- |
| `tags` _object (keys:string, values:string)_ | Tags defines the tags for the resources of the pipeline-under-test in the cloud service provider. |
| `secretRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core)_ | SecretRef defines the reference to the Kubernetes Secret object for authentication on the cloud service provider. |


#### TrafficModel



TrafficModel is the Schema for the trafficmodels API

_Appears in:_
- [TrafficModelList](#trafficmodellist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `TrafficModel`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[TrafficModelSpec](#trafficmodelspec)_ | Spec defines the specifications of the TrafficModel. |


#### TrafficModelList



TrafficModelList contains a list of TrafficModel



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `TrafficModelList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[TrafficModel](#trafficmodel) array_ | Items defines a list of TrafficModels. |


#### TrafficModelSpec



TrafficModelSpec defines the desired state of TrafficModel

_Appears in:_
- [TrafficModel](#trafficmodel)

| Field | Description |
| --- | --- |
| `config` _string_ | Config defines the configuration of the TrafficModel. |




#### WebSocket



WebSocket defines the configurations of websocket protocol.

_Appears in:_
- [Endpoint](#endpoint)

| Field | Description |
| --- | --- |
| `url` _string_ | Placeholder. |
| `params` _object (keys:string, values:string)_ | Placeholder. |
| `callback` _string_ | Placeholder. |


