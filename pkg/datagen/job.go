package datagen

import (
	"os"
	"path/filepath"

	"github.com/brianvoe/gofakeit/v6"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
)

// DataGeneratorJob is an interface for generating data.
type DataGeneratorJob interface {
	GenerateData() error
}

// JobConfig holds the configuration for a data generator job.
type JobConfig struct {
	RepeatStart int
	RepeatEnd   int
}

// BuildBasedDataGeneratorJob is a data generator job based on the Build strategy.
type BuildBasedDataGeneratorJob struct {
	RepeatStart int
	RepeatEnd   int
	Namespace   string
	DGName      string
	Dataset     *windtunnelv1alpha1.DataSet
	SchemaMap   map[string]*windtunnelv1alpha1.Schema
}

// NewBuildBasedDataGeneratorJob creates a new BuildBasedDataGeneratorJob instance.
func NewBuildBasedDataGeneratorJob(start int, end int, dgNamespace string, dgName string, dataset *windtunnelv1alpha1.DataSet, schemaMap map[string]*windtunnelv1alpha1.Schema) DataGeneratorJob {
	return &BuildBasedDataGeneratorJob{
		RepeatStart: start,
		RepeatEnd:   end,
		Namespace:   dgNamespace,
		DGName:      dgName,
		Dataset:     dataset,
		SchemaMap:   schemaMap,
	}
}

// MakeOutputDir creates the output directory for a schema in the dataset.
func MakeOutputDir(dataGeneratorConfig *windtunnelv1alpha1.DataSet, seqNum int) error {
	schPath := filepath.Join(path, dataGeneratorConfig.Spec.Schemas[seqNum].Name)
	err := os.RemoveAll(schPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(schPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// ApplyOperations applies the operations defined in the output builder to generate the final output.
func ApplyOperations(outputBuilder *OutputBuilder, seqNum int) error {
	var err error
	for _, op := range outputBuilder.Operations {
		err = op(outputBuilder, seqNum)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateData generates the data using the Build strategy.
func (dg *BuildBasedDataGeneratorJob) GenerateData() error {
	var err error

	// Create schema builders and populate the schema builder cache
	for _, schemaSelector := range dg.Dataset.Spec.Schemas {
		schemaName := schemaSelector.Name
		schemaObj := dg.SchemaMap[schemaName]
		schBldr, err := NewSchemaBuilder(schemaObj)
		if err != nil {
			return err
		}
		PutSchemaBuilder(schemaName, schBldr)
	}

	// Create the output builder
	outputBuilder, err := NewOutputBuilder(dg.Dataset)
	if err != nil {
		return err
	}

	// Set the seed for random number generation
	seed := gofakeit.New(0).Rand

	// Create output directories for each schema
	scheNum := len(outputBuilder.SchBuilders)
	for i := 0; i < scheNum; i++ {
		err := MakeOutputDir(dg.Dataset, i)
		if err != nil {
			return err
		}
	}

	// Generate data for each repeat
	for i := dg.RepeatStart; i < dg.RepeatEnd; i++ {
		// Build data for each schema
		for _, schBldr := range outputBuilder.SchBuilders {
			err := schBldr.Build(seed)
			if err != nil {
				return err
			}
		}
		// Apply operations to generate the final output
		err = ApplyOperations(outputBuilder, i)
		if err != nil {
			return err
		}
	}

	return nil
}
