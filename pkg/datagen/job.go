package datagen

import (
	"os"
	"path/filepath"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
)

// DataGeneratorJob is an interface for generating data.
type DataGeneratorJob interface {
	GenerateData(path string) error
}

// BuilderBasedDataGeneratorJob is a data generator job based on the build strategy.
type BuilderBasedDataGeneratorJob struct {
	RepeatStart int
	RepeatEnd   int
	DataSet     *windtunnelv1alpha1.DataSet
	SchemaMap   map[string]*windtunnelv1alpha1.Schema
}

// NewBuilderBasedDataGeneratorJob creates a new BuilderBasedDataGeneratorJob instance.
func NewBuilderBasedDataGeneratorJob(start, end int, dataSet *windtunnelv1alpha1.DataSet, schemaMap map[string]*windtunnelv1alpha1.Schema) DataGeneratorJob {
	return &BuilderBasedDataGeneratorJob{
		RepeatStart: start,
		RepeatEnd:   end,
		DataSet:     dataSet,
		SchemaMap:   schemaMap,
	}
}

// MakeOutputDir creates the output directory for a Schema in the DataSet.
func MakeOutputDir(dataSet *windtunnelv1alpha1.DataSet, schemaIdx int, path string) error {
	schPath := filepath.Join(path, dataSet.Spec.Schemas[schemaIdx].Name)

	// Remove any existing directory
	err := os.RemoveAll(schPath)
	if err != nil {
		return err
	}

	// Create new directory
	err = os.MkdirAll(schPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// ApplyOperations applies the operations defined in the OutputBuilder to generate the final output.
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

// GenerateData generates the data using the build strategy.
func (dg *BuilderBasedDataGeneratorJob) GenerateData(path string) error {
	var err error

	// Initiate faker for gofakeit
	faker := gofakeit.NewFaker(source.NewCrypto(), true)

	// Create SchemaBuilders and put them to cache
	for _, schemaSelector := range dg.DataSet.Spec.Schemas {
		schemaName := schemaSelector.Name
		schemaObj := dg.SchemaMap[schemaName]
		schBldr, err := NewSchemaBuilder(schemaObj)
		if err != nil {
			return err
		}
		PutSchemaBuilder(schemaName, schBldr)
	}

	// Create the OutputBuilder
	outputBuilder, err := NewOutputBuilder(dg.DataSet, path)
	if err != nil {
		return err
	}

	// Create output directories for each Schema if compression is disabled
	if dg.DataSet.Spec.CompressedFileFormat == "" {
		numSchema := len(outputBuilder.SchBuilders)
		for i := 0; i < numSchema; i++ {
			err := MakeOutputDir(dg.DataSet, i, path)
			if err != nil {
				return err
			}
		}
	}

	// Generate data for each repeat
	for i := dg.RepeatStart; i < dg.RepeatEnd; i++ {
		// Initialize the randomness and cache for each SchemaBuilder
		outputBuilder.SetRandomnessAndCache(faker, dg.DataSet)
		// Build data for each Schema
		for _, schBldr := range outputBuilder.SchBuilders {
			err := schBldr.Build(faker)
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
