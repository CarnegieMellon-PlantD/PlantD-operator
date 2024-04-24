package proxy

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/datagen"
)

// Constants defining the possible kinds that can be used in the schema.GroupVersionKind struct.
const (
	SchemaKind       string = "Schema"
	DatasetKind      string = "DataSet"
	LoadPatternKind  string = "LoadPattern"
	PipelineKind     string = "Pipeline"
	ExperimentKind   string = "Experiment"
	CostExporterKind string = "CostExporter"
	DigitalTwinKind  string = "DigitalTwin"
	SimulationKind   string = "Simulation"
	TrafficModelKind string = "TrafficModel"
	NetCostKind      string = "NetCost"
	ScenarioKind     string = "Scenario"
	PlantDCoreKind   string = "PlantDCore"
)

// AllKinds is the list of all possible kinds for import/export.
var AllKinds = []string{
	SchemaKind,
	DatasetKind,
	LoadPatternKind,
	PipelineKind,
	ExperimentKind,
	CostExporterKind,
	DigitalTwinKind,
	SimulationKind,
	TrafficModelKind,
	NetCostKind,
	ScenarioKind,
}

// ForObject returns a client.Object instance based on the provided kind.
func ForObject(kind string) (client.Object, error) {
	switch kind {
	case SchemaKind:
		return &windtunnelv1alpha1.Schema{}, nil
	case DatasetKind:
		return &windtunnelv1alpha1.DataSet{}, nil
	case LoadPatternKind:
		return &windtunnelv1alpha1.LoadPattern{}, nil
	case PipelineKind:
		return &windtunnelv1alpha1.Pipeline{}, nil
	case ExperimentKind:
		return &windtunnelv1alpha1.Experiment{}, nil
	case CostExporterKind:
		return &windtunnelv1alpha1.CostExporter{}, nil
	case DigitalTwinKind:
		return &windtunnelv1alpha1.DigitalTwin{}, nil
	case SimulationKind:
		return &windtunnelv1alpha1.Simulation{}, nil
	case TrafficModelKind:
		return &windtunnelv1alpha1.TrafficModel{}, nil
	case NetCostKind:
		return &windtunnelv1alpha1.NetCost{}, nil
	case ScenarioKind:
		return &windtunnelv1alpha1.Scenario{}, nil
	case PlantDCoreKind:
		return &windtunnelv1alpha1.PlantDCore{}, nil
	}
	return nil, fmt.Errorf("failed to find resource of kind \"%s\"", kind)
}

// AddFileToZip adds a file with the given content to a zip.Writer.
func AddFileToZip(zw *zip.Writer, name string, content []byte) error {
	fw, err := zw.Create(name)
	if err != nil {
		return err
	}

	_, err = fw.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// GetSampleDataSet generates a sample DataSet and compresses it into a ZIP file stream.
func GetSampleDataSet(ctx context.Context, c client.Client, namespace, name string) (*bytes.Buffer, error) {
	dataSet := &windtunnelv1alpha1.DataSet{}
	if err := c.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, dataSet); err != nil {
		return nil, fmt.Errorf("while getting DataSet: %w", err)
	}

	schemaMap := map[string]*windtunnelv1alpha1.Schema{}
	for _, schemaSelector := range dataSet.Spec.Schemas {
		schema := &windtunnelv1alpha1.Schema{}
		if err := c.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      schemaSelector.Name,
		}, schema); err != nil {
			return nil, fmt.Errorf("while getting Schema: %w", err)
		}
		schemaMap[schema.Name] = schema
	}

	// Modify DataSet to have only 1 repeat
	// and 1 file per Schema per compressed file if compression is enabled
	dataSet.Spec.NumberOfFiles = 1
	if dataSet.Spec.CompressedFileFormat != "" {
		for schemaSelectorIdx, _ := range dataSet.Spec.Schemas {
			dataSet.Spec.Schemas[schemaSelectorIdx].NumRecords.Min = 1
			dataSet.Spec.Schemas[schemaSelectorIdx].NumRecords.Max = 1
			dataSet.Spec.Schemas[schemaSelectorIdx].NumFilesPerCompressedFile.Min = 1
			dataSet.Spec.Schemas[schemaSelectorIdx].NumFilesPerCompressedFile.Max = 1
		}
	}

	// Generate data in a temporary directory
	tmpPath, err := os.MkdirTemp("/tmp", "sample_dataset_*")
	if err != nil {
		return nil, fmt.Errorf("while creating temporary directory: %w", err)
	}
	job := datagen.NewBuilderBasedDataGeneratorJob(0, 1, dataSet, schemaMap)
	if err := job.GenerateData(tmpPath); err != nil {
		return nil, fmt.Errorf("while generating data: %w", err)
	}

	// Compress the generated data into a ZIP file stream
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	defer zw.Close()

	if err := filepath.Walk(tmpPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return AddFileToZip(zw, info.Name(), content)
	}); err != nil {
		return nil, fmt.Errorf("while compressing files: %w", err)
	}

	// Cleanup temporary directory
	if err := os.RemoveAll(tmpPath); err != nil {
		return nil, fmt.Errorf("while removing temporary directory: %w", err)
	}

	return buf, nil
}

func ListKinds() []string {
	return AllKinds
}

func ListResources(ctx context.Context, c client.Client) ([]*ResourceLocator, error) {
	result := make([]*ResourceLocator, 0)
	for _, kind := range AllKinds {
		objList := &unstructured.UnstructuredList{}
		objList.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   windtunnelv1alpha1.GroupVersion.Group,
			Version: windtunnelv1alpha1.GroupVersion.Version,
			Kind:    kind,
		})

		if err := c.List(ctx, objList); err != nil {
			return nil, err
		}

		for _, item := range objList.Items {
			result = append(result, &ResourceLocator{
				Kind:      kind,
				Namespace: item.GetNamespace(),
				Name:      item.GetName(),
			})
		}
	}
	return result, nil
}

func ImportResources(ctx context.Context, c client.Client, buf *bytes.Buffer) (*ImportStatistics, error) {
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return nil, fmt.Errorf("while opening zip file: %s", err.Error())
	}

	result := &ImportStatistics{
		NumSucceeded:  0,
		NumFailed:     0,
		ErrorMessages: []string{},
	}
	for _, file := range zr.File {
		fr, err := file.Open()
		if err != nil {
			result.NumFailed++
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("while opening file \"%s\": %s", file.Name, err.Error()))
			continue
		}

		fileContent, err := io.ReadAll(fr)
		if err != nil {
			result.NumFailed++
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("while reading file \"%s\": %s", file.Name, err.Error()))
			continue
		}

		obj := &unstructured.Unstructured{}
		if err = yaml.Unmarshal(fileContent, obj); err != nil {
			result.NumFailed++
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("while unmarshalling file \"%s\": %s", file.Name, err.Error()))
			continue
		}

		if err = c.Create(ctx, obj); err != nil {
			result.NumFailed++
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("while creating object in file \"%s\": %s", file.Name, err.Error()))
			continue
		}

		result.NumSucceeded++
	}

	return result, nil
}

func ExportResources(ctx context.Context, c client.Client, resInfoList []*ResourceLocator) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	for idx, info := range resInfoList {
		obj := &unstructured.Unstructured{}
		obj.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   windtunnelv1alpha1.GroupVersion.Group,
			Version: windtunnelv1alpha1.GroupVersion.Version,
			Kind:    info.Kind,
		})

		if err := c.Get(ctx, types.NamespacedName{
			Namespace: info.Namespace,
			Name:      info.Name,
		}, obj); err != nil {
			return nil, fmt.Errorf("while getting object at pos %d: %s", idx, err.Error())
		}

		// For the output object, we manually set the GVK, which will result in the apiVersion and kind fields.
		// We also manually set the namespace and name, so that the metadata field will only contain these two fields.
		// Last, we copy the spec field from the fetched object.
		objToOutput := &unstructured.Unstructured{}
		objToOutput.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   windtunnelv1alpha1.GroupVersion.Group,
			Version: windtunnelv1alpha1.GroupVersion.Version,
			Kind:    info.Kind,
		})
		objToOutput.SetNamespace(info.Namespace)
		objToOutput.SetName(info.Name)
		objToOutput.Object["spec"] = obj.Object["spec"]

		fileContent, err := yaml.Marshal(objToOutput)
		if err != nil {
			return nil, fmt.Errorf("while marshalling object at pos %d: %s", idx, err.Error())
		}

		err = AddFileToZip(zw, fmt.Sprintf("%s_%s_%s.yaml", info.Kind, info.Namespace, info.Name), fileContent)
		if err != nil {
			return nil, fmt.Errorf("while writing object at pos %d to file: %s", idx, err.Error())
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}
