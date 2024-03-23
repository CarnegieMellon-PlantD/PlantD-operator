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
- [NetCost](#netcost)
- [NetCostList](#netcostlist)
- [Pipeline](#pipeline)
- [PipelineList](#pipelinelist)
- [PlantDCore](#plantdcore)
- [PlantDCoreList](#plantdcorelist)
- [Scenario](#scenario)
- [ScenarioList](#scenariolist)
- [Schema](#schema)
- [SchemaList](#schemalist)
- [Simulation](#simulation)
- [SimulationList](#simulationlist)
- [TrafficModel](#trafficmodel)
- [TrafficModelList](#trafficmodellist)



#### ColumnSpec



ColumnSpec defines the column in Schema.

_Appears in:_
- [SchemaSpec](#schemaspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the column. |
| `type` _string_ | Data type of the random data to be generated in the column. Used together with the `params` field. It should be a valid function name in gofakeit, which can be parsed by gofakeit.GetFuncLookup(). `formula` field has precedence over this field. See https://plantd.org/docs/reference/types-and-params for available values. |
| `params` _object (keys:string, values:string)_ | Map of parameters for generating the data in the column. Used together with the `type` field. For any parameters not provided but required by the data type, the default value will be used, if available. Will ignore any parameters not used by the data type. See https://plantd.org/docs/reference/types-and-params for available values. |
| `formula` _[FormulaSpec](#formulaspec)_ | Formula to be applied for populating the data in the column. This field has precedence over the `type` fields. |


#### CostExporter



CostExporter is the Schema for the costexporters API

_Appears in:_
- [CostExporterList](#costexporterlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `CostExporter`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[CostExporterSpec](#costexporterspec)_ |  |


#### CostExporterList



CostExporterList contains a list of CostExporter



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `CostExporterList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[CostExporter](#costexporter) array_ |  |


#### CostExporterSpec



CostExporterSpec defines the desired state of CostExporter

_Appears in:_
- [CostExporter](#costexporter)

| Field | Description |
| --- | --- |
| `s3Bucket` _string_ | S3Bucket defines the AWS S3 bucket name where stores the cost logs. |
| `cloudServiceProvider` _string_ | CloudServiceProvider defines the target cloud service provide for calculating cost. |
| `secretRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | SecretRef defines the reference to the Kubernetes Secret where stores the credentials of cloud service provider |




#### DataSet



DataSet is the Schema for the datasets API

_Appears in:_
- [DataSetList](#datasetlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DataSet`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DataSetSpec](#datasetspec)_ |  |


#### DataSetConfig



DataSetConfig defines the parameters to generate DataSet

_Appears in:_
- [ScenarioSpec](#scenariospec)

| Field | Description |
| --- | --- |
| `compressPerSchema` _boolean_ |  |
| `compressedFileFormat` _string_ |  |
| `fileFormat` _string_ |  |


#### DataSetErrorType

_Underlying type:_ _string_

DataSetErrorType defines the type of error occurred.

_Appears in:_
- [DataSetStatus](#datasetstatus)



#### DataSetJobStatus

_Underlying type:_ _string_

DataSetJobStatus defines the status of the data generating job.

_Appears in:_
- [DataSetStatus](#datasetstatus)



#### DataSetList



DataSetList contains a list of DataSet



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DataSetList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[DataSet](#dataset) array_ |  |


#### DataSetSpec



DataSetSpec defines the desired state of DataSet.

_Appears in:_
- [DataSet](#dataset)

| Field | Description |
| --- | --- |
| `fileFormat` _string_ | Format of the output file containing generated data. Available values are `csv` and `binary`. |
| `compressedFileFormat` _string_ | Format of the compressed file containing output files. Available value is `zip`. Leave empty to disable compression. |
| `compressPerSchema` _boolean_ | Flag for compression behavior. Takes effect only if `compressedFileFormat` is set. When set to `false` (default), files from all Schemas will be compressed into a single compressed file in each repetition. When set to `true`, files from each Schema will be compressed into a separate compressed file in each repetition. |
| `numFiles` _integer_ | Number of repetitions of the data generation process. If `compressedFileFormat` is unset, this is the number of files for each Schema. If `compressedFileFormat` is set and `compressPerSchema` is `false`, this is the number of compressed files for each Schema. If `compressedFileFormat` is set and `compressPerSchema` is `true`, this is the total number of compressed files. |
| `schemas` _[SchemaSelector](#schemaselector) array_ | List of Schemas in the DataSet. |
| `parallelJobs` _integer_ | Number of parallel jobs when generating the dataset. |




#### DataSpec



DataSpec defines the data to be sent to the endpoint.

_Appears in:_
- [EndpointSpec](#endpointspec)

| Field | Description |
| --- | --- |
| `plainText` _string_ | PlainText defines a plain text data. |
| `dataSetRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | DataSetRef defines the reference of the DataSet object. |


#### DeploymentConfig



DeploymentConfig defines the desired state of modules managed as Deployment

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `image` _string_ | Image defines the container image to use |
| `replicas` _integer_ | Replicas defines the desired number of replicas |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources defines the resource requirements per replica |


#### DigitalTwin



DigitalTwin is the Schema for the digitaltwins API

_Appears in:_
- [DigitalTwinList](#digitaltwinlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DigitalTwin`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DigitalTwinSpec](#digitaltwinspec)_ |  |


#### DigitalTwinList



DigitalTwinList contains a list of DigitalTwin



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DigitalTwinList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[DigitalTwin](#digitaltwin) array_ |  |


#### DigitalTwinSpec



DigitalTwinSpec defines the desired state of DigitalTwin

_Appears in:_
- [DigitalTwin](#digitaltwin)

| Field | Description |
| --- | --- |
| `modelType` _string_ | ModelType defines the type of the DigitalTwin model. |
| `experiments` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core) array_ | Experiments contains the list of Experiment object references for the DigitalTwin. |




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
| `serviceRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | ServiceRef defines the Kubernetes Service that exposes metrics. |
| `port` _string_ | Name of the Service port which this endpoint refers to. <br /><br /> It takes precedence over `targetPort`. |
| `targetPort` _[IntOrString](https://pkg.go.dev/k8s.io/apimachinery/pkg/util/intstr#IntOrString)_ | Name or number of the target port of the `Pod` object behind the Service, the port must be specified with container port property. <br /><br /> Deprecated: use `port` instead. |
| `path` _string_ | HTTP path from which to scrape for metrics. <br /><br /> If empty, Prometheus uses the default value (e.g. `/metrics`). |
| `scheme` _string_ | HTTP scheme to use for scraping. <br /><br /> `http` and `https` are the expected values unless you rewrite the `__scheme__` label via relabeling. <br /><br /> If empty, Prometheus uses the default value `http`. |
| `params` _object (keys:string, values:string array)_ | params define optional HTTP URL parameters. |
| `interval` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | Interval at which Prometheus scrapes the metrics from the target. <br /><br /> If empty, Prometheus uses the global scrape interval. |
| `scrapeTimeout` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | Timeout after which Prometheus considers the scrape to be failed. <br /><br /> If empty, Prometheus uses the global scrape timeout unless it is less than the target's scrape interval value in which the latter is used. |
| `tlsConfig` _[TLSConfig](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#TLSConfig)_ | TLS configuration to use when scraping the target. |
| `bearerTokenFile` _string_ | File to read bearer token for scraping the target. <br /><br /> Deprecated: use `authorization` instead. |
| `bearerTokenSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#secretkeyselector-v1-core)_ | `bearerTokenSecret` specifies a key of a Secret containing the bearer token for scraping targets. The secret needs to be in the same namespace as the ServiceMonitor object and readable by the Prometheus Operator. <br /><br /> Deprecated: use `authorization` instead. |
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


#### EndpointSpec



EndpointSpec defines the DataSet and LoadPattern to be used for an endpoint.

_Appears in:_
- [ExperimentSpec](#experimentspec)

| Field | Description |
| --- | --- |
| `endpointName` _string_ | EndpointName defines the name of endpoint. It should be the name of an existing endpoint defined in the Pipeline used in the Experiment. |
| `dataSpec` _[DataSpec](#dataspec)_ | DataSpec defines the data to be sent to the endpoint. |
| `loadPatternRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | LoadPatternRef defines the reference of the LoadPattern object. |


#### Experiment



Experiment is the Schema for the experiments API

_Appears in:_
- [ExperimentList](#experimentlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Experiment`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ExperimentSpec](#experimentspec)_ |  |


#### ExperimentList



ExperimentList contains a list of Experiments.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `ExperimentList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Experiment](#experiment) array_ |  |


#### ExperimentSpec



ExperimentSpec defines the desired state of Experiment

_Appears in:_
- [Experiment](#experiment)

| Field | Description |
| --- | --- |
| `pipelineRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | PipelineRef defines a reference of the Pipeline object. |
| `endpointSpecs` _[EndpointSpec](#endpointspec) array_ | EndpointSpecs defines a list of configurations for the endpoints. |
| `scheduledTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#time-v1-meta)_ | ScheduledTime defines the scheduled time for the Experiment. |




#### ExtraMetrics



ExtraMetrics defines the configurations of getting extra metrics.

_Appears in:_
- [PipelineSpec](#pipelinespec)

| Field | Description |
| --- | --- |
| `system` _[SystemMetrics](#systemmetrics)_ | System defines the configurfation of getting system metrics. |
| `messageQueue` _[MessageQueueMetrics](#messagequeuemetrics)_ | MessageQueue defines the configurfation of getting message queue related metrics. |


#### FormulaSpec



FormulaSpec defines the formula in column.

_Appears in:_
- [ColumnSpec](#columnspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the formula. Used together with the `args` field. See https://plantd.org/docs/reference/formulas for available values. |
| `args` _string array_ | Arguments to be passed to the formula. Used together with the `name` field. See https://plantd.org/docs/reference/formulas for available values. |


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


#### LoadPattern



LoadPattern is the Schema for the loadpatterns API

_Appears in:_
- [LoadPatternList](#loadpatternlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `LoadPattern`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[LoadPatternSpec](#loadpatternspec)_ |  |


#### LoadPatternList



LoadPatternList contains a list of LoadPattern



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `LoadPatternList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[LoadPattern](#loadpattern) array_ |  |


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



#### NetCost



NetCost is the Schema for the netcosts API

_Appears in:_
- [NetCostList](#netcostlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `NetCost`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[NetCostSpec](#netcostspec)_ |  |


#### NetCostList



NetCostList contains a list of NetCost



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `NetCostList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[NetCost](#netcost) array_ |  |


#### NetCostSpec



NetCostSpec defines the desired state of NetCost

_Appears in:_
- [NetCost](#netcost)

| Field | Description |
| --- | --- |
| `netCostPerMB` _[Quantity](#quantity)_ | NetCostPerMB defines the cost per MB of data transfer. |
| `rawDataStoreCostPerMBMonth` _[Quantity](#quantity)_ | RawDataStoreCostPerMBMonth defines the cost per MB per month of raw data storage. |
| `processedDataStoreCostPerMBMonth` _[Quantity](#quantity)_ | ProcessedDataStoreCostPerMBMonth defines the cost per MB per month of processed data storage. |
| `rawDataRetentionPolicyMonths` _integer_ | RawDataRetentionPolicyMonths defines the months raw data is retained. |
| `processedDataRetentionPolicyMonths` _integer_ | ProcessedDataRetentionPolicyMonths defines the months processed data is retained. |




#### Pipeline



Pipeline is the Schema for the pipelines API

_Appears in:_
- [PipelineList](#pipelinelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Pipeline`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[PipelineSpec](#pipelinespec)_ |  |


#### PipelineList



PipelineList contains a list of Pipeline



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PipelineList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Pipeline](#pipeline) array_ |  |


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
| `experimentRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | Internal usage. For experiment object to lock the pipeline object. |




#### PlantDCore



PlantDCore is the Schema for the plantdcores API

_Appears in:_
- [PlantDCoreList](#plantdcorelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PlantDCore`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[PlantDCoreSpec](#plantdcorespec)_ |  |


#### PlantDCoreList



PlantDCoreList contains a list of PlantDCore



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PlantDCoreList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[PlantDCore](#plantdcore) array_ |  |


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
| `thanosEnabled` _boolean_ | ThanosEnabled defines if Thanos is enabled (True / False) |




#### PrometheusConfig



PrometheusConfig defines the desired state of Prometheus

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `scrapeInterval` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | ScrapeInterval defines the desired time length between scrapings |
| `replicas` _integer_ | Replicas defines the desired number of replicas |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources defines the resource requirements per replica |


#### Scenario



Scenario is the Schema for the scenarios API

_Appears in:_
- [ScenarioList](#scenariolist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Scenario`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ScenarioSpec](#scenariospec)_ |  |


#### ScenarioList



ScenarioList contains a list of Scenario



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `ScenarioList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Scenario](#scenario) array_ |  |


#### ScenarioSpec



ScenarioSpec defines the desired state of Scenario

_Appears in:_
- [Scenario](#scenario)

| Field | Description |
| --- | --- |
| `dataSetConfig` _[DataSetConfig](#datasetconfig)_ | DataSetConfig defines the parameters to generate DataSet. |
| `pipelineRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | PipelineRef defines the reference to the Pipeline object. |
| `tasks` _[ScenarioTask](#scenariotask) array_ | Tasks defines the list of tasks to be executed in the Scenario. |




#### ScenarioTask



ScenarioTask defines the task to be executed in the Scenario

_Appears in:_
- [ScenarioSpec](#scenariospec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name defines the name of the task. |
| `size` _[Quantity](#quantity)_ | Size defines the size of a single upload in bytes. |
| `sendingDevices` _object (keys:string, values:integer)_ | SendingDevices defines the range of the devices to send the data. |
| `pushFrequencyPerMonth` _object (keys:string, values:integer)_ | PushFrequencyPerMonth defines the range of how many times the data is pushed per month. |
| `monthsRelevant` _integer array_ | MonthsRelevant defines the months the task is relevant. |


#### Schema



Schema is the Schema for the schemas API

_Appears in:_
- [SchemaList](#schemalist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Schema`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[SchemaSpec](#schemaspec)_ |  |


#### SchemaList



SchemaList contains a list of Schema



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `SchemaList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Schema](#schema) array_ |  |


#### SchemaSelector



SchemaSelector defines the reference to a Schema and its usage in the DataSet.

_Appears in:_
- [DataSetSpec](#datasetspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the Schema. Note that the Schema must be present in the same namespace as the DataSet. |
| `numRecords` _object (keys:string, values:integer)_ | Range of number of rows to be generated in each output file. Should be a map containing `min` and `max` keys. For each output file, a random number is picked from the specified range. |
| `numFilesPerCompressedFile` _object (keys:string, values:integer)_ | Range of number of files to be generated in the compressed file. Take effect only if `compressedFileFormat` is set in the DataSet. Should be a map containing `min` and `max` keys. A random number is picked from the specified range. |


#### SchemaSpec



SchemaSpec defines the desired state of Schema.

_Appears in:_
- [Schema](#schema)

| Field | Description |
| --- | --- |
| `columns` _[ColumnSpec](#columnspec) array_ | List of columns in the Schema. |




#### Simulation



Simulation is the Schema for the simulations API

_Appears in:_
- [SimulationList](#simulationlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `Simulation`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[SimulationSpec](#simulationspec)_ |  |


#### SimulationList



SimulationList contains a list of Simulation



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `SimulationList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Simulation](#simulation) array_ |  |


#### SimulationSpec



SimulationSpec defines the desired state of Simulation

_Appears in:_
- [Simulation](#simulation)

| Field | Description |
| --- | --- |
| `trafficModelRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | TrafficModelRef defines the TrafficModel object reference for the Simulation. |
| `digitalTwinRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | DigitalTwinRef defines the DigitalTwin object reference for the Simulation. |




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
| `secretRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | SecretRef defines the reference to the Kubernetes Secret object for authentication on the cloud service provider. |


#### TrafficModel



TrafficModel is the Schema for the trafficmodels API

_Appears in:_
- [TrafficModelList](#trafficmodellist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `TrafficModel`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[TrafficModelSpec](#trafficmodelspec)_ | Spec defines the specifications of the TrafficModel. |


#### TrafficModelList



TrafficModelList contains a list of TrafficModel



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `TrafficModelList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
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


