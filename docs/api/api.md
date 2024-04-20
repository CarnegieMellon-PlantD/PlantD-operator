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



#### Column



Column defines the column in Schema.

_Appears in:_
- [SchemaSpec](#schemaspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the column. |
| `type` _string_ | Data type of the random data to be generated in the column. Used together with the `params` field. It should be a valid function name in gofakeit, which can be parsed by gofakeit.GetFuncLookup(). `formula` field has precedence over this field. See https://plantd.org/docs/reference/types-and-params for available values. |
| `params` _object (keys:string, values:string)_ | Map of parameters for generating the data in the column. Used together with the `type` field. For any parameters not provided but required by the data type, the default value will be used, if available. Will ignore any parameters not used by the data type. See https://plantd.org/docs/reference/types-and-params for available values. |
| `formula` _[Formula](#formula)_ | Formula to be applied for populating the data in the column. This field has precedence over the `type` fields. |


#### ComponentStatus



ComponentStatus defines the status of a component.

_Appears in:_
- [PlantDCoreStatus](#plantdcorestatus)

| Field | Description |
| --- | --- |
| `text` _[ComponentStatusText](#componentstatustext)_ | Component status string. |
| `numReady` _integer_ | Number of ready replicas. |
| `numDesired` _integer_ | Number of desired replicas. |


#### ComponentStatusText

_Underlying type:_ _string_

ComponentStatusText defines the status of a component.

_Appears in:_
- [ComponentStatus](#componentstatus)



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


#### DataSetErrorType

_Underlying type:_ _string_

DataSetErrorType defines the type of error occurred.

_Appears in:_
- [DataSetStatus](#datasetstatus)



#### DataSetJobStatus

_Underlying type:_ _string_

DataSetJobStatus defines the status of the data generator job.

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
| `image` _string_ | Image of the data generator job. |
| `parallelism` _integer_ | Number of parallel jobs when generating the dataset. Default to 1. |
| `storageSize` _[Quantity](#quantity)_ | Size of the PVC for the data generator job. Default to 2Gi. |
| `fileFormat` _string_ | Format of the output file containing generated data. Available values are `csv` and `binary`. |
| `compressedFileFormat` _string_ | Format of the compressed file containing output files. Available value is `zip`. Leave empty to disable compression. |
| `compressPerSchema` _boolean_ | Flag for compression behavior. Takes effect only if `compressedFileFormat` is set. When set to `false` (default), files from all Schemas will be compressed into a single compressed file in each repetition. When set to `true`, files from each Schema will be compressed into a separate compressed file in each repetition. |
| `numFiles` _integer_ | Number of files to be generated. If `compressedFileFormat` is unset, this is the number of files for each Schema. If `compressedFileFormat` is set and `compressPerSchema` is `false`, this is the number of compressed files for each Schema. If `compressedFileFormat` is set and `compressPerSchema` is `true`, this is the total number of compressed files. |
| `schemas` _[SchemaSelector](#schemaselector) array_ | List of Schemas in the DataSet. |




#### DataSpec



DataSpec defines the data to be sent to an endpoint.

_Appears in:_
- [EndpointSpec](#endpointspec)

| Field | Description |
| --- | --- |
| `plainText` _string_ | PlainText data to be sent. `dataSetRef` field has precedence over this field. |
| `dataSetRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#localobjectreference-v1-core)_ | Reference to the DataSet to be sent. The DataSet must be in the same namespace as the Experiment. This field has precedence over the `plainText` field. |


#### DeploymentConfig



DeploymentConfig defines the desired state of a component deployed as Deployment.

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `replicas` _integer_ | Number of replicas. |
| `image` _string_ | Container image to use. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources requirements. |


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


#### DigitalTwinJobStatus

_Underlying type:_ _string_

DigitalTwinJobStatus defines the status of the Experiments created by DigitalTwin.

_Appears in:_
- [DigitalTwinStatus](#digitaltwinstatus)



#### DigitalTwinList



DigitalTwinList contains a list of DigitalTwin



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `DigitalTwinList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[DigitalTwin](#digitaltwin) array_ |  |


#### DigitalTwinSpec



DigitalTwinSpec defines the desired state of DigitalTwin.

_Appears in:_
- [DigitalTwin](#digitaltwin)

| Field | Description |
| --- | --- |
| `modelType` _string_ | Type of digital twin model. Available values are `simple`, `quickscaling`, and `autoscaling`. |
| `digitalTwinType` _string_ | Type of digital twin. Available values are `regular` and `schemaaware`. |
| `experiments` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core) array_ | Existing Experiments to retrieve metrics data from to train the DigitalTwin. Effective only when `digitalTwinType` is `regular`. |
| `dataSet` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#localobjectreference-v1-core)_ | DataSet to break down into Schemas to train the DigitalTwin. Effective only when `digitalTwinType` is `schemaaware`. |
| `pipeline` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#localobjectreference-v1-core)_ | Pipeline to use to train the DigitalTwin. Effective only when `digitalTwinType` is `schemaaware`. |
| `pipelineCapacity` _integer_ | Maximum RPS in the populated LoadPatterns. Effective only when `digitalTwinType` is `schemaaware`. |








#### EndpointSpec



EndpointSpec defines the test upon an endpoint.

_Appears in:_
- [ExperimentSpec](#experimentspec)

| Field | Description |
| --- | --- |
| `endpointName` _string_ | Name of endpoint. It should be an existing endpoint defined in the Pipeline used by the Experiment. |
| `dataSpec` _[DataSpec](#dataspec)_ | Data to be sent to the endpoint. |
| `loadPatternRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | LoadPattern to follow for the endpoint. |
| `storageSize` _[Quantity](#quantity)_ | Size of the PVC for the load generator job. Only effective when `dataSpec.dataSetRef` is set. Default to the PVC size of the DataSet. |


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


#### ExperimentJobStatus

_Underlying type:_ _string_

ExperimentJobStatus defines the status of the load generator job.

_Appears in:_
- [ExperimentStatus](#experimentstatus)



#### ExperimentList



ExperimentList contains a list of Experiments.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `ExperimentList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Experiment](#experiment) array_ |  |


#### ExperimentSpec



ExperimentSpec defines the desired state of Experiment.

_Appears in:_
- [Experiment](#experiment)

| Field | Description |
| --- | --- |
| `pipelineRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#localobjectreference-v1-core)_ | Reference to the Pipeline to use for the Experiment. |
| `endpointSpecs` _[EndpointSpec](#endpointspec) array_ | List of tests upon endpoints. |
| `scheduledTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#time-v1-meta)_ | Scheduled time to run the Experiment. |
| `drainingTime` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#duration-v1-meta)_ | Time to wait after the load generator job is completed before finishing the Experiment. It allows the pipeline-under-test to finish its processing. Default to no draining time. |




#### Formula



Formula defines the formula in column.

_Appears in:_
- [Column](#column)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the formula. Used together with the `args` field. See https://plantd.org/docs/reference/formulas for available values. |
| `args` _string array_ | Arguments to be passed to the formula. Used together with the `name` field. See https://plantd.org/docs/reference/formulas for available values. |


#### HTTP



HTTP defines the configurations of HTTP protocol in endpoint.

_Appears in:_
- [MetricsEndpoint](#metricsendpoint)
- [PipelineEndpoint](#pipelineendpoint)

| Field | Description |
| --- | --- |
| `url` _string_ | URL of the HTTP request. |
| `method` _string_ | Method of the HTTP request. |
| `headers` _object (keys:string, values:string)_ | Headers of the HTTP request. |


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



LoadPatternSpec defines the desired state of LoadPattern.

_Appears in:_
- [LoadPattern](#loadpattern)

| Field | Description |
| --- | --- |
| `stages` _[Stage](#stage) array_ | List of stages in the LoadPattern. Equivalent to the "ramping-arrival-rate" executor's `stages` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |
| `preAllocatedVUs` _integer_ | Number of VUs to pre-allocate before Experiment start. Equivalent to the "ramping-arrival-rate" executor's `preAllocatedVUs` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |
| `startRate` _integer_ | Number of requests per `timeUnit` period at Experiment start. Equivalent to the "ramping-arrival-rate" executor's `startRate` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |
| `timeUnit` _string_ | Period of time to apply to the `startRate` and `stages[].target` fields. Equivalent to the "ramping-arrival-rate" executor's `timeUnit` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |
| `maxVUs` _integer_ | Maximum number of VUs to allow for allocation during Experiment. Equivalent to the "ramping-arrival-rate" executor's `maxVUs` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |




#### MetricsEndpoint



MetricsEndpoint defines the endpoint for metrics scraping in Pipeline.

_Appears in:_
- [PipelineSpec](#pipelinespec)

| Field | Description |
| --- | --- |
| `http` _[HTTP](#http)_ | Configurations of the HTTP protocol. Only the `http.url` field will be used. Must be set if `inCluster` is set to `false` in the Pipeline. |
| `serviceRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#localobjectreference-v1-core)_ | Reference to the Service. The Service must be in the same namespace as the Pipeline. Must be set if `inCluster` is set to `true` in the Pipeline. |
| `port` _string_ | Name of the Service port to use. Effective only when `inCluster` is set to `true` in the Pipeline. Default to "metrics". |
| `path` _string_ | Path of the endpoint. Effective only when `inCluster` is set to `true` in the Pipeline. Default to "/metrics". |


#### NaturalIntRange



NaturalIntRange defines a range using two non-negative integers as boundaries.

_Appears in:_
- [ScenarioTask](#scenariotask)
- [SchemaSelector](#schemaselector)

| Field | Description |
| --- | --- |
| `min` _integer_ | Minimum value of the range. |
| `max` _integer_ | Maximum value of the range. |


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



NetCostSpec defines the desired state of NetCost.

_Appears in:_
- [NetCost](#netcost)

| Field | Description |
| --- | --- |
| `netCostPerMB` _string_ | The cost per MB of data transfer. The value should be a float number in string format. |
| `rawDataStoreCostPerMBMonth` _string_ | The cost per MB per month of raw data storage. The value should be a float number in string format. |
| `processedDataStoreCostPerMBMonth` _string_ | The cost per MB per month of processed data storage. The value should be a float number in string format. |
| `rawDataRetentionPolicyMonths` _integer_ | The number of months the raw data is retained. |
| `processedDataRetentionPolicyMonths` _integer_ | The number of months the processed data is retained. |




#### OpenCostConfig



OpenCostConfig defines the desired state of an OpenCost component.

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `replicas` _integer_ | Number of replicas. |
| `image` _string_ | Container image to use for OpenCost. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources requirements for OpenCost. |
| `uiImage` _string_ | Container image to use for OpenCost-UI. |
| `uiResources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources requirements for OpenCost-UI. |


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


#### PipelineAvailability

_Underlying type:_ _string_

PipelineAvailability defines the availability of the Pipeline.

_Appears in:_
- [PipelineStatus](#pipelinestatus)



#### PipelineEndpoint



PipelineEndpoint defines the endpoint for data ingestion in Pipeline.

_Appears in:_
- [PipelineSpec](#pipelinespec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the endpoint. |
| `http` _[HTTP](#http)_ | Configurations of the HTTP protocol. |


#### PipelineList



PipelineList contains a list of Pipeline



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `windtunnel.plantd.org/v1alpha1`
| `kind` _string_ | `PipelineList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Pipeline](#pipeline) array_ |  |


#### PipelineSpec



PipelineSpec defines the desired state of Pipeline.

_Appears in:_
- [Pipeline](#pipeline)

| Field | Description |
| --- | --- |
| `inCluster` _boolean_ | Whether the Pipeline is deployed within the cluster or not. When set to `false`, Services of type ExternalName will be created to access the Pipeline. When set to `true`, the Pipeline will be accessed by its Services. |
| `pipelineEndpoints` _[PipelineEndpoint](#pipelineendpoint) array_ | List of endpoints for data ingestion. |
| `metricsEndpoint` _[MetricsEndpoint](#metricsendpoint)_ | Endpoint for metrics scraping. |
| `healthCheckURLs` _string array_ | List of URLs for health check. An HTTP GET request will be made to each URL, and all of them should return 200 OK to pass the health check. If the list is empty, no health check will be performed. |
| `enableCostCalculation` _boolean_ | Whether to enable cost calculation for the Pipeline. |
| `cloudProvider` _string_ | Cloud provider of the Pipeline. Available values are `aws`, `azure`, and `gcp`. |
| `tags` _object (keys:string, values:string)_ | Map of tags to select cloud resources of the Pipeline. Equivalent to the tags in the cloud service provider. |




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



PlantDCoreSpec defines the desired state of PlantDCore.

_Appears in:_
- [PlantDCore](#plantdcore)

| Field | Description |
| --- | --- |
| `proxy` _[DeploymentConfig](#deploymentconfig)_ | PlantD-Proxy configuration. |
| `studio` _[DeploymentConfig](#deploymentconfig)_ | PlantD-Studio configuration. |
| `prometheus` _[PrometheusConfig](#prometheusconfig)_ | Prometheus configuration. |
| `thanos` _[ThanosConfig](#thanosconfig)_ | Thanos configuration. |
| `redis` _[StatefulSetConfig](#statefulsetconfig)_ | Redis configuration. |
| `opencost` _[OpenCostConfig](#opencostconfig)_ | OpenCost configuration. |




#### PrometheusConfig



PrometheusConfig defines the desired state of a Prometheus component.

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `replicas` _integer_ | Number of replicas. |
| `scrapeInterval` _[Duration](https://pkg.go.dev/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1#Duration)_ | Interval at which Prometheus scrapes metrics. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources requirements. |


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
| `tasks` _[ScenarioTask](#scenariotask) array_ | List of tasks in the Scenario. |




#### ScenarioTask



ScenarioTask defines the task in the Scenario.

_Appears in:_
- [ScenarioSpec](#scenariospec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the task. Should be a Schema name. |
| `size` _[Quantity](#quantity)_ | The size of a single upload in bytes. |
| `sendingDevices` _[NaturalIntRange](#naturalintrange)_ | Range of the number range of the devices to send the data. |
| `pushFrequencyPerMonth` _[NaturalIntRange](#naturalintrange)_ | Range of the frequency of data pushes per month. |
| `monthsRelevant` _integer array_ | List of months the task will apply to. For example, `[1, 12]` means the task will apply to January and December. |


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
| `numRecords` _[NaturalIntRange](#naturalintrange)_ | Range of number of rows to be generated in each output file. |
| `numFilesPerCompressedFile` _[NaturalIntRange](#naturalintrange)_ | Range of number of files to be generated in the compressed file. Take effect only if `compressedFileFormat` is set in the DataSet. |


#### SchemaSpec



SchemaSpec defines the desired state of Schema.

_Appears in:_
- [Schema](#schema)

| Field | Description |
| --- | --- |
| `columns` _[Column](#column) array_ | List of columns in the Schema. |




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
| `digitalTwinRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | DigitalTwin object for the Simulation. |
| `trafficModelRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | TrafficModel object for the Simulation. |
| `netCostRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | NetCost object for the Simulation. Optional. |
| `scenarioRef` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectreference-v1-core)_ | Scenario object for the Simulation. The task names in the Scenario must be the name of a Schema in the DataSet used by the DigitalTwin. Mandatory if the `digitalTwinType` field of the DigitalTwin is `schemaaware`. Always ignored otherwise. |




#### Stage



Stage defines how the load ramps up or down.

_Appears in:_
- [LoadPatternSpec](#loadpatternspec)

| Field | Description |
| --- | --- |
| `target` _integer_ | Target load to reach at the end of the stage. Equivalent to the "ramping-arrival-rate" executor's `stages[].target` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |
| `duration` _string_ | Duration of the stage, also the time to reach the target load. Equivalent to the "ramping-arrival-rate" executor's `stages[].duration` option in K6. See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details. |


#### StatefulSetConfig



StatefulSetConfig defines the desired state of a component deployed as StatefulSet.

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)
- [ThanosConfig](#thanosconfig)

| Field | Description |
| --- | --- |
| `replicas` _integer_ | Number of replicas. |
| `image` _string_ | Container image to use. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#resourcerequirements-v1-core)_ | Resources requirements. |
| `storageSize` _[Quantity](#quantity)_ | Storage size. |


#### ThanosConfig



ThanosConfig defines the desired state of a Thanos component.

_Appears in:_
- [PlantDCoreSpec](#plantdcorespec)

| Field | Description |
| --- | --- |
| `image` _string_ | Thanos image to use. Must be synced with the `version` field. |
| `version` _string_ | Thanos version to use. Must be synced with the `image` field. |
| `objectStoreConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#secretkeyselector-v1-core)_ | Object store configuration for Thanos. Set this field will enable upload in Thanos-Sidecar, and deploy Thanos-Store and Thanos-Compactor. |
| `sidecar` _[StatefulSetConfig](#statefulsetconfig)_ | Thanos-Sidecar configuration. The `sidecar.replicas`, `sidecar.image` and `sidecar.storageSize` fields are always ignored. |
| `store` _[StatefulSetConfig](#statefulsetconfig)_ | Thanos-Store configuration. The `store.image` field is always ignored. This field is ignored if `objectStoreConfig` is not set. |
| `compactor` _[StatefulSetConfig](#statefulsetconfig)_ | Thanos-Compactor configuration. The `compactor.image` field is always ignored. This field is ignored if `objectStoreConfig` is not set. |
| `querier` _[StatefulSetConfig](#statefulsetconfig)_ | Thanos-Querier configuration. The `querier.image` and `querier.storageSize` fields are always ignored. |


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



TrafficModelSpec defines the desired state of TrafficModel.

_Appears in:_
- [TrafficModel](#trafficmodel)

| Field | Description |
| --- | --- |
| `config` _string_ | TrafficModel configuration in JSON. |




