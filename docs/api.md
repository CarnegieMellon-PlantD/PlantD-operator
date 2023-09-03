---
title: "API reference"
description: "Online Programming Exercise Platform operator generated API reference docs"
draft: false
images: []
menu: "operator"
weight: 211
toc: true
---
> This page is automatically generated with `gen-crd-api-reference-docs`.
<p>Packages:</p>
<ul>
<li>
<a href="#windtunnel.plantd.org%2fv1alpha1">windtunnel.plantd.org/v1alpha1</a>
</li>
</ul>
<h2 id="windtunnel.plantd.org/v1alpha1">windtunnel.plantd.org/v1alpha1</h2>
<div>
<p>Package v1alpha1 contains API Schema definitions for the windtunnel v1alpha1 API group</p>
</div>
Resource Types:
<ul></ul>
<h3 id="windtunnel.plantd.org/v1alpha1.Column">Column
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.SchemaSpec">SchemaSpec</a>)
</p>
<div>
<p>Column defines the metadata of the column data.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines the name of the column.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>Type defines the data type of the column. Should match the type with one of the provided types.</p>
</td>
</tr>
<tr>
<td>
<code>params</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Params defines the parameters for constructing the data give certain data type.</p>
</td>
</tr>
<tr>
<td>
<code>formula</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.FormulaSpec">
FormulaSpec
</a>
</em>
</td>
<td>
<p>Formula defines the formula applies to the column data.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.CostExporter">CostExporter
</h3>
<div>
<p>CostExporter is the Schema for the costexporters API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.CostExporterSpec">
CostExporterSpec
</a>
</em>
</td>
<td>
<p>Spec defines the specifictions of the CostExporter.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
string
</em>
</td>
<td>
<p>S3Bucket defines the AWS S3 bucket name where stores the cost logs.</p>
</td>
</tr>
<tr>
<td>
<code>cloudServiceProvider</code><br/>
<em>
string
</em>
</td>
<td>
<p>CloudServiceProvider defines the target cloud service provide for calculating cost.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>SecretRef defines the reference to the Kubernetes Secret where stores the credentials of cloud service provider</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.CostExporterStatus">
CostExporterStatus
</a>
</em>
</td>
<td>
<p>Status defines the status of the CostExporter.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.CostExporterSpec">CostExporterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.CostExporter">CostExporter</a>)
</p>
<div>
<p>CostExporterSpec defines the desired state of CostExporter</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
string
</em>
</td>
<td>
<p>S3Bucket defines the AWS S3 bucket name where stores the cost logs.</p>
</td>
</tr>
<tr>
<td>
<code>cloudServiceProvider</code><br/>
<em>
string
</em>
</td>
<td>
<p>CloudServiceProvider defines the target cloud service provide for calculating cost.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>SecretRef defines the reference to the Kubernetes Secret where stores the credentials of cloud service provider</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.CostExporterStatus">CostExporterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.CostExporter">CostExporter</a>)
</p>
<div>
<p>CostExporterStatus defines the observed state of CostExporter</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>jobCompletionTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>JobCompletionTime defines the completion time of the cost calculation job.</p>
</td>
</tr>
<tr>
<td>
<code>podName</code><br/>
<em>
string
</em>
</td>
<td>
<p>PodName defines the name of the cost calculation pod.</p>
</td>
</tr>
<tr>
<td>
<code>jobStatus</code><br/>
<em>
string
</em>
</td>
<td>
<p>JobStatus defines the status of the cost calculation job.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
string
</em>
</td>
<td>
<p>Tags defines the json string of using tags.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.DataSet">DataSet
</h3>
<div>
<p>DataSet is the Schema for the datasets API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.DataSetSpec">
DataSetSpec
</a>
</em>
</td>
<td>
<p>Spec defines the specifications of the DataSet.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>fileFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>FileFormat defines the file format of the each file containing the generated data.
This may or may not be the output file format based on whether you want to compress these files.</p>
</td>
</tr>
<tr>
<td>
<code>compressedFileFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>CompressedFileFormat defines the file format for the compressed files.
Each file inside the compressed file is of &ldquo;fileFormat&rdquo; format specified above.
This is the output format if specified for the files.</p>
</td>
</tr>
<tr>
<td>
<code>compressPerSchema</code><br/>
<em>
bool
</em>
</td>
<td>
<p>CompressPerSchema defines the flag of compression.
If you wish files from all the different schemas to compressed into one compressed file leave this field as false.
If you wish to have a different compressed file for every schema, mark this field as true.</p>
</td>
</tr>
<tr>
<td>
<code>numFiles</code><br/>
<em>
int32
</em>
</td>
<td>
<p>NumberOfFiles defines the total number of output files irrespective of compression.
Unless &ldquo;compressPerSchema&rdquo; is false, this field is applicable per schema.</p>
</td>
</tr>
<tr>
<td>
<code>schemas</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.SchemaSelector">
[]SchemaSelector
</a>
</em>
</td>
<td>
<p>Schemas defines a list of Schemas.</p>
</td>
</tr>
<tr>
<td>
<code>parallelJobs</code><br/>
<em>
int32
</em>
</td>
<td>
<p>ParallelJobs defines the number of parallel jobs when generating the dataset.
TODO: Infer the optimal number of parallel jobs automatically.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.DataSetStatus">
DataSetStatus
</a>
</em>
</td>
<td>
<p>Status defines the status of the DataSet.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.DataSetSpec">DataSetSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.DataSet">DataSet</a>)
</p>
<div>
<p>DataSetSpec defines the desired state of DataSet</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>fileFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>FileFormat defines the file format of the each file containing the generated data.
This may or may not be the output file format based on whether you want to compress these files.</p>
</td>
</tr>
<tr>
<td>
<code>compressedFileFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>CompressedFileFormat defines the file format for the compressed files.
Each file inside the compressed file is of &ldquo;fileFormat&rdquo; format specified above.
This is the output format if specified for the files.</p>
</td>
</tr>
<tr>
<td>
<code>compressPerSchema</code><br/>
<em>
bool
</em>
</td>
<td>
<p>CompressPerSchema defines the flag of compression.
If you wish files from all the different schemas to compressed into one compressed file leave this field as false.
If you wish to have a different compressed file for every schema, mark this field as true.</p>
</td>
</tr>
<tr>
<td>
<code>numFiles</code><br/>
<em>
int32
</em>
</td>
<td>
<p>NumberOfFiles defines the total number of output files irrespective of compression.
Unless &ldquo;compressPerSchema&rdquo; is false, this field is applicable per schema.</p>
</td>
</tr>
<tr>
<td>
<code>schemas</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.SchemaSelector">
[]SchemaSelector
</a>
</em>
</td>
<td>
<p>Schemas defines a list of Schemas.</p>
</td>
</tr>
<tr>
<td>
<code>parallelJobs</code><br/>
<em>
int32
</em>
</td>
<td>
<p>ParallelJobs defines the number of parallel jobs when generating the dataset.
TODO: Infer the optimal number of parallel jobs automatically.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.DataSetStatus">DataSetStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.DataSet">DataSet</a>)
</p>
<div>
<p>DataSetStatus defines the observed state of DataSet</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>jobStatus</code><br/>
<em>
string
</em>
</td>
<td>
<p>JobStatus defines the status of the data generating job.</p>
</td>
</tr>
<tr>
<td>
<code>pvcStatus</code><br/>
<em>
string
</em>
</td>
<td>
<p>PVCStatus defines the status of the PVC mount to the data generating pod.</p>
</td>
</tr>
<tr>
<td>
<code>startTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>StartTime defines the start time of the data generating job.</p>
</td>
</tr>
<tr>
<td>
<code>completionTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>CompletionTime defines the duration of the data generating job.</p>
</td>
</tr>
<tr>
<td>
<code>lastGeneration</code><br/>
<em>
int64
</em>
</td>
<td>
<p>LastGeneration defines the last generation of the DataSet object.</p>
</td>
</tr>
<tr>
<td>
<code>errorCount</code><br/>
<em>
int
</em>
</td>
<td>
<p>ErrorCount defines the number of errors raised by the controller or data generating job.</p>
</td>
</tr>
<tr>
<td>
<code>errors</code><br/>
<em>
map[string][]string
</em>
</td>
<td>
<p>Errors defines the map of error messages.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.Endpoint">Endpoint
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.PipelineSpec">PipelineSpec</a>)
</p>
<div>
<p>Endpoint defines the configuration of the endpoint.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines the name of the endpoint. It&rsquo;s required when it&rsquo;s for pipeline endpoint.</p>
</td>
</tr>
<tr>
<td>
<code>http</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.HTTP">
HTTP
</a>
</em>
</td>
<td>
<p>HTTP defines the configuration of the HTTP request. It&rsquo;s mutually exclusive with WebSocket and GRPC.</p>
</td>
</tr>
<tr>
<td>
<code>websocket</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.WebSocket">
WebSocket
</a>
</em>
</td>
<td>
<p>WebSocket defines the configuration of the WebSocket connection. It&rsquo;s mutually exclusive with HTTP and GRPC.</p>
</td>
</tr>
<tr>
<td>
<code>grpc</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.GRPC">
GRPC
</a>
</em>
</td>
<td>
<p>GRPC defines the configuration of the gRPC request. It&rsquo;s mutually exclusive with HTTP and WebSocket.</p>
</td>
</tr>
<tr>
<td>
<code>serviceRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>ServiceRef defines the Kubernetes Service that exposes metrics.</p>
</td>
</tr>
<tr>
<td>
<code>port</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the service port this endpoint refers to. Mutually exclusive with targetPort.</p>
</td>
</tr>
<tr>
<td>
<code>targetPort</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/util/intstr#IntOrString">
k8s.io/apimachinery/pkg/util/intstr.IntOrString
</a>
</em>
</td>
<td>
<p>Name or number of the target port of the Pod behind the Service, the port must be specified with container port property. Mutually exclusive with port.</p>
</td>
</tr>
<tr>
<td>
<code>path</code><br/>
<em>
string
</em>
</td>
<td>
<p>HTTP path to scrape for metrics.
If empty, Prometheus uses the default value (e.g. <code>/metrics</code>).</p>
</td>
</tr>
<tr>
<td>
<code>scheme</code><br/>
<em>
string
</em>
</td>
<td>
<p>HTTP scheme to use for scraping.
<code>http</code> and <code>https</code> are the expected values unless you rewrite the <code>__scheme__</code> label via relabeling.
If empty, Prometheus uses the default value <code>http</code>.</p>
</td>
</tr>
<tr>
<td>
<code>params</code><br/>
<em>
map[string][]string
</em>
</td>
<td>
<p>Optional HTTP URL parameters</p>
</td>
</tr>
<tr>
<td>
<code>interval</code><br/>
<em>
github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.Duration
</em>
</td>
<td>
<p>Interval at which metrics should be scraped
If not specified Prometheus&rsquo; global scrape interval is used.</p>
</td>
</tr>
<tr>
<td>
<code>scrapeTimeout</code><br/>
<em>
github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.Duration
</em>
</td>
<td>
<p>Timeout after which the scrape is ended
If not specified, the Prometheus global scrape timeout is used unless it is less than <code>Interval</code> in which the latter is used.</p>
</td>
</tr>
<tr>
<td>
<code>tlsConfig</code><br/>
<em>
github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.TLSConfig
</em>
</td>
<td>
<p>TLS configuration to use when scraping the endpoint</p>
</td>
</tr>
<tr>
<td>
<code>bearerTokenFile</code><br/>
<em>
string
</em>
</td>
<td>
<p>File to read bearer token for scraping targets.</p>
</td>
</tr>
<tr>
<td>
<code>bearerTokenSecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Secret to mount to read bearer token for scraping targets. The secret
needs to be in the same namespace as the service monitor and accessible by
the Prometheus Operator.</p>
</td>
</tr>
<tr>
<td>
<code>authorization</code><br/>
<em>
github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.SafeAuthorization
</em>
</td>
<td>
<p>Authorization section for this endpoint</p>
</td>
</tr>
<tr>
<td>
<code>honorLabels</code><br/>
<em>
bool
</em>
</td>
<td>
<p>HonorLabels chooses the metric&rsquo;s labels on collisions with target labels.</p>
</td>
</tr>
<tr>
<td>
<code>honorTimestamps</code><br/>
<em>
bool
</em>
</td>
<td>
<p>HonorTimestamps controls whether Prometheus respects the timestamps present in scraped data.</p>
</td>
</tr>
<tr>
<td>
<code>basicAuth</code><br/>
<em>
github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.BasicAuth
</em>
</td>
<td>
<p>BasicAuth allow an endpoint to authenticate over basic authentication
More info: <a href="https://prometheus.io/docs/operating/configuration/#endpoints">https://prometheus.io/docs/operating/configuration/#endpoints</a></p>
</td>
</tr>
<tr>
<td>
<code>oauth2</code><br/>
<em>
github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.OAuth2
</em>
</td>
<td>
<p>OAuth2 for the URL. Only valid in Prometheus versions 2.27.0 and newer.</p>
</td>
</tr>
<tr>
<td>
<code>metricRelabelings</code><br/>
<em>
[]github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.RelabelConfig
</em>
</td>
<td>
<p>MetricRelabelConfigs to apply to samples before ingestion.</p>
</td>
</tr>
<tr>
<td>
<code>relabelings</code><br/>
<em>
[]github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1.RelabelConfig
</em>
</td>
<td>
<p>RelabelConfigs to apply to samples before scraping.
Prometheus Operator automatically adds relabelings for a few standard Kubernetes fields.
The original scrape job&rsquo;s name is available via the <code>__tmp_prometheus_job_name</code> label.
More info: <a href="https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config">https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config</a></p>
</td>
</tr>
<tr>
<td>
<code>proxyUrl</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProxyURL eg <a href="http://proxyserver:2195">http://proxyserver:2195</a> Directs scrapes to proxy through this endpoint.</p>
</td>
</tr>
<tr>
<td>
<code>followRedirects</code><br/>
<em>
bool
</em>
</td>
<td>
<p>FollowRedirects configures whether scrape requests follow HTTP 3xx redirects.</p>
</td>
</tr>
<tr>
<td>
<code>enableHttp2</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Whether to enable HTTP2.</p>
</td>
</tr>
<tr>
<td>
<code>filterRunning</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Drop pods that are not running. (Failed, Succeeded). Enabled by default.
More info: <a href="https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase">https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.Experiment">Experiment
</h3>
<div>
<p>Experiment is the Schema for the experiments API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.ExperimentSpec">
ExperimentSpec
</a>
</em>
</td>
<td>
<p>Spec defines the specifications of the Experiment.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>pipelineRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>PipelineRef defines s reference of the Pipeline object.</p>
</td>
</tr>
<tr>
<td>
<code>loadPatterns</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.LoadPatternConfig">
[]LoadPatternConfig
</a>
</em>
</td>
<td>
<p>LoadPatterns defines a list of configuration of name of endpoints and LoadPatterns.</p>
</td>
</tr>
<tr>
<td>
<code>scheduledTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>ScheduledTime defines the scheduled time for the Experiment.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.ExperimentStatus">
ExperimentStatus
</a>
</em>
</td>
<td>
<p>Status defines the status of the Experiment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.ExperimentSpec">ExperimentSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Experiment">Experiment</a>)
</p>
<div>
<p>ExperimentSpec defines the desired state of Experiment</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>pipelineRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>PipelineRef defines s reference of the Pipeline object.</p>
</td>
</tr>
<tr>
<td>
<code>loadPatterns</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.LoadPatternConfig">
[]LoadPatternConfig
</a>
</em>
</td>
<td>
<p>LoadPatterns defines a list of configuration of name of endpoints and LoadPatterns.</p>
</td>
</tr>
<tr>
<td>
<code>scheduledTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>ScheduledTime defines the scheduled time for the Experiment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.ExperimentStatus">ExperimentStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Experiment">Experiment</a>)
</p>
<div>
<p>ExperimentStatus defines the observed state of Experiment</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>experimentState</code><br/>
<em>
string
</em>
</td>
<td>
<p>ExperimentState defines the state of the Experiment.</p>
</td>
</tr>
<tr>
<td>
<code>protocols</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Protocols defines a map of name of endpoint (key) to request protocol (value).</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Tags defines the a map of key-value pair that use for tagging cloud resources.</p>
</td>
</tr>
<tr>
<td>
<code>duration</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
map[string]k8s.io/apimachinery/pkg/apis/meta/v1.Duration
</a>
</em>
</td>
<td>
<p>Duration defines the duration of the K6 load generator.</p>
</td>
</tr>
<tr>
<td>
<code>startTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>StartTime defines the start of the K6 load generator.</p>
</td>
</tr>
<tr>
<td>
<code>endTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>EndTime defines the end of the Experiment.
TODO: Add microservice to calculate the end time of the experiment.</p>
</td>
</tr>
<tr>
<td>
<code>cloudVendor</code><br/>
<em>
string
</em>
</td>
<td>
<p>CloudVendor defines the cloud service provider which the pipeline-under-test is deployed.</p>
</td>
</tr>
<tr>
<td>
<code>enableCostCalculation</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableCostCalculation defines teh flag of cost calculation.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.ExtraMetrics">ExtraMetrics
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.PipelineSpec">PipelineSpec</a>)
</p>
<div>
<p>ExtraMetrics defines the configurations of getting extra metrics.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>system</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.SystemMetrics">
SystemMetrics
</a>
</em>
</td>
<td>
<p>System defines the configurfation of getting system metrics.</p>
</td>
</tr>
<tr>
<td>
<code>messageQueue</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.MessageQueueMetrics">
MessageQueueMetrics
</a>
</em>
</td>
<td>
<p>MessageQueue defines the configurfation of getting message queue related metrics.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.FormulaSpec">FormulaSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Column">Column</a>)
</p>
<div>
<p>FormulaSpec defines the specification of the formula.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines the name of the formula. Should match the name with one of the provided formulas.</p>
</td>
</tr>
<tr>
<td>
<code>args</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Args defines the arugments for calling the formula.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.GRPC">GRPC
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Endpoint">Endpoint</a>)
</p>
<div>
<p>GRPC defines the configurations of gRPC protocol.
TODO: Validate the gRPC library in K6 and update the API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>address</code><br/>
<em>
string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
<tr>
<td>
<code>protoFiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
<tr>
<td>
<code>url</code><br/>
<em>
string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
<tr>
<td>
<code>params</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
<tr>
<td>
<code>request</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.HTTP">HTTP
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Endpoint">Endpoint</a>)
</p>
<div>
<p>HTTP defines the configurations of HTTP protocol.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>url</code><br/>
<em>
string
</em>
</td>
<td>
<p>URL defines the absolute path for an entry point of the Pipeline.</p>
</td>
</tr>
<tr>
<td>
<code>method</code><br/>
<em>
string
</em>
</td>
<td>
<p>Method defines the HTTP method used for the endpoint.</p>
</td>
</tr>
<tr>
<td>
<code>headers</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Headers defines a map of HTTP headers.</p>
</td>
</tr>
<tr>
<td>
<code>body</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.HTTPBodySpec">
HTTPBodySpec
</a>
</em>
</td>
<td>
<p>Body defines the configurations of the HTTP request body.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.HTTPBodySpec">HTTPBodySpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.HTTP">HTTP</a>)
</p>
<div>
<p>HTTPBodySpec defines the configurations of the HTTP request body.
User can specify either Data or DataSetRef, but not both fields.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>data</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>dataSetRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.LoadPattern">LoadPattern
</h3>
<div>
<p>LoadPattern is the Schema for the loadpatterns API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.LoadPatternSpec">
LoadPatternSpec
</a>
</em>
</td>
<td>
<p>Spec defines the specification of the LoadPattern.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>stages</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Stage">
[]Stage
</a>
</em>
</td>
<td>
<p>Stages defines a list of stages for the LoadPattern.</p>
</td>
</tr>
<tr>
<td>
<code>preAllocatedVUs</code><br/>
<em>
int
</em>
</td>
<td>
<p>PreAllocatedVUs defines pre-allocated virtual users for the K6 load generator.</p>
</td>
</tr>
<tr>
<td>
<code>startRate</code><br/>
<em>
int
</em>
</td>
<td>
<p>StartRate defines the initial requests per second when the K6 load generator starts.</p>
</td>
</tr>
<tr>
<td>
<code>maxVUs</code><br/>
<em>
int
</em>
</td>
<td>
<p>MaxVUs defines the maximum virtual users for the K6 load generator.</p>
</td>
</tr>
<tr>
<td>
<code>timeUnit</code><br/>
<em>
string
</em>
</td>
<td>
<p>TimeUnit defines the unit of the time for K6 load generator.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.LoadPatternStatus">
LoadPatternStatus
</a>
</em>
</td>
<td>
<p>Status defines the status of the LoadPattern.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.LoadPatternConfig">LoadPatternConfig
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.ExperimentSpec">ExperimentSpec</a>)
</p>
<div>
<p>LoadPatternConfig defines the configuration of the load pattern in the experiment.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>endpointName</code><br/>
<em>
string
</em>
</td>
<td>
<p>EndpointName defines the name of endpoint where to send the requests.
It should match the name of endpoint declared in the specification of the pipeline.</p>
</td>
</tr>
<tr>
<td>
<code>loadPatternRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>LoadPatternRef defines s reference of the LoadPattern object.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.LoadPatternSpec">LoadPatternSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.LoadPattern">LoadPattern</a>)
</p>
<div>
<p>LoadPatternSpec defines the desired state of LoadPattern</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>stages</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Stage">
[]Stage
</a>
</em>
</td>
<td>
<p>Stages defines a list of stages for the LoadPattern.</p>
</td>
</tr>
<tr>
<td>
<code>preAllocatedVUs</code><br/>
<em>
int
</em>
</td>
<td>
<p>PreAllocatedVUs defines pre-allocated virtual users for the K6 load generator.</p>
</td>
</tr>
<tr>
<td>
<code>startRate</code><br/>
<em>
int
</em>
</td>
<td>
<p>StartRate defines the initial requests per second when the K6 load generator starts.</p>
</td>
</tr>
<tr>
<td>
<code>maxVUs</code><br/>
<em>
int
</em>
</td>
<td>
<p>MaxVUs defines the maximum virtual users for the K6 load generator.</p>
</td>
</tr>
<tr>
<td>
<code>timeUnit</code><br/>
<em>
string
</em>
</td>
<td>
<p>TimeUnit defines the unit of the time for K6 load generator.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.LoadPatternStatus">LoadPatternStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.LoadPattern">LoadPattern</a>)
</p>
<div>
<p>LoadPatternStatus defines the observed state of LoadPattern</p>
</div>
<h3 id="windtunnel.plantd.org/v1alpha1.MessageQueueMetrics">MessageQueueMetrics
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.ExtraMetrics">ExtraMetrics</a>)
</p>
<div>
<p>MessageQueueMetrics defines the configurations of getting message queue related metrics.</p>
</div>
<h3 id="windtunnel.plantd.org/v1alpha1.Pipeline">Pipeline
</h3>
<div>
<p>Pipeline is the Schema for the pipelines API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.PipelineSpec">
PipelineSpec
</a>
</em>
</td>
<td>
<p>Spec defines the specifications of the Pipeline.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>pipelineEndpoints</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Endpoint">
[]Endpoint
</a>
</em>
</td>
<td>
<p>Endpoints for pipeline-under-test.</p>
</td>
</tr>
<tr>
<td>
<code>healthCheckEndpoints</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Endpoints for health check.</p>
</td>
</tr>
<tr>
<td>
<code>metricsEndpoint</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Endpoint">
Endpoint
</a>
</em>
</td>
<td>
<p>Endpoints for metrics.</p>
</td>
</tr>
<tr>
<td>
<code>extraMetrics</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.ExtraMetrics">
ExtraMetrics
</a>
</em>
</td>
<td>
<p>Extra metrics, such as CPU utilzation, I/O and etc.</p>
</td>
</tr>
<tr>
<td>
<code>inCluster</code><br/>
<em>
bool
</em>
</td>
<td>
<p>In cluster flag. True indecates the pipeline-under-test is deployed in the same cluster as the plantD. Otherwise it should be False.</p>
</td>
</tr>
<tr>
<td>
<code>cloudVendor</code><br/>
<em>
string
</em>
</td>
<td>
<p>State which cloud service provider the pipeline is deployed.</p>
</td>
</tr>
<tr>
<td>
<code>enableCostCalculation</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Cost calculation flag.</p>
</td>
</tr>
<tr>
<td>
<code>experimentRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>Internal usage. For experiment object to lock the pipeline object.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.PipelineStatus">
PipelineStatus
</a>
</em>
</td>
<td>
<p>Status defines the status of the Pipeline.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.PipelineSpec">PipelineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Pipeline">Pipeline</a>)
</p>
<div>
<p>PipelineSpec defines the desired state of Pipeline</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>pipelineEndpoints</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Endpoint">
[]Endpoint
</a>
</em>
</td>
<td>
<p>Endpoints for pipeline-under-test.</p>
</td>
</tr>
<tr>
<td>
<code>healthCheckEndpoints</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Endpoints for health check.</p>
</td>
</tr>
<tr>
<td>
<code>metricsEndpoint</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Endpoint">
Endpoint
</a>
</em>
</td>
<td>
<p>Endpoints for metrics.</p>
</td>
</tr>
<tr>
<td>
<code>extraMetrics</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.ExtraMetrics">
ExtraMetrics
</a>
</em>
</td>
<td>
<p>Extra metrics, such as CPU utilzation, I/O and etc.</p>
</td>
</tr>
<tr>
<td>
<code>inCluster</code><br/>
<em>
bool
</em>
</td>
<td>
<p>In cluster flag. True indecates the pipeline-under-test is deployed in the same cluster as the plantD. Otherwise it should be False.</p>
</td>
</tr>
<tr>
<td>
<code>cloudVendor</code><br/>
<em>
string
</em>
</td>
<td>
<p>State which cloud service provider the pipeline is deployed.</p>
</td>
</tr>
<tr>
<td>
<code>enableCostCalculation</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Cost calculation flag.</p>
</td>
</tr>
<tr>
<td>
<code>experimentRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>Internal usage. For experiment object to lock the pipeline object.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.PipelineStatus">PipelineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Pipeline">Pipeline</a>)
</p>
<div>
<p>PipelineStatus defines the observed state of Pipeline</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>pipelineState</code><br/>
<em>
string
</em>
</td>
<td>
<p>PipelineState defines the state of the Pipeline.</p>
</td>
</tr>
<tr>
<td>
<code>statusCheck</code><br/>
<em>
string
</em>
</td>
<td>
<p>StatusCheck defines the health status of the Pipeline.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.PlantDCore">PlantDCore
</h3>
<div>
<p>PlantDCore is the Schema for the plantdcores API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.PlantDCoreSpec">
PlantDCoreSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of PlantDCore. Edit plantdcore_types.go to remove/update</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.PlantDCoreStatus">
PlantDCoreStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.PlantDCoreSpec">PlantDCoreSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.PlantDCore">PlantDCore</a>)
</p>
<div>
<p>PlantDCoreSpec defines the desired state of PlantDCore</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of PlantDCore. Edit plantdcore_types.go to remove/update</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.PlantDCoreStatus">PlantDCoreStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.PlantDCore">PlantDCore</a>)
</p>
<div>
<p>PlantDCoreStatus defines the observed state of PlantDCore</p>
</div>
<h3 id="windtunnel.plantd.org/v1alpha1.Schema">Schema
</h3>
<div>
<p>Schema is the Schema for the schemas API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.SchemaSpec">
SchemaSpec
</a>
</em>
</td>
<td>
<p>Spec defines the specifications of the Schema.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>columns</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Column">
[]Column
</a>
</em>
</td>
<td>
<p>Columns defines a list of column specifications.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.SchemaStatus">
SchemaStatus
</a>
</em>
</td>
<td>
<p>Status defines the status of the Schema.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.SchemaSelector">SchemaSelector
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.DataSetSpec">DataSetSpec</a>)
</p>
<div>
<p>SchemaSelector defines a list of Schemas and the required numbers and format.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines the name of the Schame. Should match the name of existing Schema in the same namespace as the DataSet.</p>
</td>
</tr>
<tr>
<td>
<code>numRecords</code><br/>
<em>
map[string]int
</em>
</td>
<td>
<p>NumRecords defines the number of records to be generated in each output file. A random number is picked from the specified range.</p>
</td>
</tr>
<tr>
<td>
<code>numFilesPerCompressedFile</code><br/>
<em>
map[string]int
</em>
</td>
<td>
<p>NumberOfFilesPerCompressedFile defines the number of intermediate files to be compressed into a single compressed file.
A random number is picked from the specified range.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.SchemaSpec">SchemaSpec
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Schema">Schema</a>)
</p>
<div>
<p>SchemaSpec defines the desired state of Schema</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>columns</code><br/>
<em>
<a href="#windtunnel.plantd.org/v1alpha1.Column">
[]Column
</a>
</em>
</td>
<td>
<p>Columns defines a list of column specifications.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.SchemaStatus">SchemaStatus
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Schema">Schema</a>)
</p>
<div>
<p>SchemaStatus defines the observed state of Schema</p>
</div>
<h3 id="windtunnel.plantd.org/v1alpha1.Stage">Stage
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.LoadPatternSpec">LoadPatternSpec</a>)
</p>
<div>
<p>Stage defines the stage configuration of the load.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>target</code><br/>
<em>
int
</em>
</td>
<td>
<p>Target defines the target requests per second.</p>
</td>
</tr>
<tr>
<td>
<code>duration</code><br/>
<em>
string
</em>
</td>
<td>
<p>Duration defines the duration of the current stage.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.SystemMetrics">SystemMetrics
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.ExtraMetrics">ExtraMetrics</a>)
</p>
<div>
<p>SystemMetrics defines the configurations of getting system metrics.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Tags defines the tags for the resources of the pipeline-under-test in the cloud service provider.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectreference-v1-core">
Kubernetes core/v1.ObjectReference
</a>
</em>
</td>
<td>
<p>SecretRef defines the reference to the Kubernetes Secret object for authentication on the cloud service provider.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="windtunnel.plantd.org/v1alpha1.WebSocket">WebSocket
</h3>
<p>
(<em>Appears on:</em><a href="#windtunnel.plantd.org/v1alpha1.Endpoint">Endpoint</a>)
</p>
<div>
<p>WebSocket defines the configurations of websocket protocol.
TODO: Validate the websocket library in K6 and update the API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>url</code><br/>
<em>
string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
<tr>
<td>
<code>params</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
<tr>
<td>
<code>callback</code><br/>
<em>
string
</em>
</td>
<td>
<p>Placeholder.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
