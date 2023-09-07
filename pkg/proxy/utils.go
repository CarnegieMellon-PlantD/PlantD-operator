package proxy

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"strconv"

	plantdv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/datagen"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/errors"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

	"github.com/brianvoe/gofakeit/v6"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetDataset retrieves a DataSet object by namespace and name.
func GetDataset(ctx context.Context, c client.Client, namespace string, name string) (*plantdv1alpha1.DataSet, error) {
	dataset := &plantdv1alpha1.DataSet{}
	if err := c.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, dataset); err != nil {
		return nil, err
	}
	return dataset, nil
}

// GetSchema retrieves a Schema object by namespace and name.
func GetSchema(ctx context.Context, c client.Client, namespace string, name string) (*plantdv1alpha1.Schema, error) {
	schema := &plantdv1alpha1.Schema{}
	if err := c.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, schema); err != nil {
		return nil, err
	}
	return schema, nil
}

// GetSampleDataSet generates a sample dataset based on the provided dataset name.
func GetSampleDataSet(ctx context.Context, c client.Client, namespace string, datasetName string) (string, *bytes.Buffer, error) {

	// Get the DataSet object
	dataset, err := GetDataset(ctx, c, namespace, datasetName)
	if err != nil {
		return "", nil, err
	}

	// Get the names of the associated schemas
	var schemaNames []string
	for _, schemaSelector := range dataset.Spec.Schemas {
		schemaName := schemaSelector.Name
		schemaObj, err := GetSchema(ctx, c, namespace, schemaName)
		if err != nil {
			return "", nil, err
		}
		schBldr, err := datagen.NewSchemaBuilder(schemaObj)
		if err != nil {
			return "", nil, err
		}
		datagen.PutSchemaBuilder(schemaName, schBldr)
		schemaNames = append(schemaNames, schemaName)
	}

	// Build the output data based on the dataset and associated schemas
	outputBuilder, err := datagen.NewOutputBuilder(dataset)
	if err != nil {
		return "", nil, err
	}

	// Generate random data based on the output builder and seed
	seed := gofakeit.New(0).Rand
	for _, schBldr := range outputBuilder.SchBuilders {
		err := schBldr.Build(seed)
		if err != nil {
			return "", nil, err
		}
	}

	// Generate the sample dataset based on the specified file format and compression
	if dataset.Spec.CompressedFileFormat != "" {
		// Compressed file format is specified
		if dataset.Spec.CompressedFileFormat != "zip" {
			return "", nil, errors.OperationUndefinedError(dataset.Spec.CompressedFileFormat)
		}
		if dataset.Spec.FileFormat == "csv" {
			// Generate CSV files and compress them into a ZIP archive
			numRecordsMap := make(map[string]int, len(dataset.Spec.Schemas))
			numberOfFilesPerCompressedFileMap := make(map[string]map[string]int, len(dataset.Spec.Schemas))

			// Generate CSV files for each schema and store the number of records
			for i, schemaName := range schemaNames {
				minRec := dataset.Spec.Schemas[i].NumRecords["min"]
				maxRec := dataset.Spec.Schemas[i].NumRecords["max"]

				numRecords := gofakeit.Number(minRec, maxRec)
				numRecordsMap[schemaName] = numRecords
				err := datagen.Raw2CSVAtCacheBySchema(numRecords, schemaName)
				if err != nil {
					return "", nil, err
				}

				numberOfFilesPerCompressedFile := dataset.Spec.Schemas[i].NumberOfFilesPerCompressedFile
				numberOfFilesPerCompressedFileMap[schemaName] = numberOfFilesPerCompressedFile
			}

			// Convert CSV files to a ZIP archive
			b, err := datagen.CSVAtCache2ZipAtBytes(schemaNames, numberOfFilesPerCompressedFileMap, 0, numRecordsMap)
			if err != nil {
				return "", nil, err
			}
			return "zip", bytes.NewBuffer(b), nil

		} else if dataset.Spec.FileFormat == "binary" {
			// Generate binary files and compress them into a ZIP archive
			numRecordsMap := make(map[string]int, len(dataset.Spec.Schemas))
			numberOfFilesPerCompressedFileMap := make(map[string]map[string]int, len(dataset.Spec.Schemas))
			for i, schemaName := range schemaNames {
				minRec := dataset.Spec.Schemas[i].NumRecords["min"]
				maxRec := dataset.Spec.Schemas[i].NumRecords["max"]

				numRecords := gofakeit.Number(minRec, maxRec)
				numRecordsMap[schemaName] = numRecords
				err := datagen.Raw2BinaryAtCacheBySchema(numRecords, schemaName)
				if err != nil {
					return "", nil, err
				}
				numberOfFilesPerCompressedFile := dataset.Spec.Schemas[i].NumberOfFilesPerCompressedFile
				numberOfFilesPerCompressedFileMap[schemaName] = numberOfFilesPerCompressedFile
			}

			// Convert binary files to a ZIP archive
			b, err := datagen.BinaryAtCache2ZipAtBytes(schemaNames, numberOfFilesPerCompressedFileMap, 0, numRecordsMap)
			if err != nil {
				return "", nil, err
			}
			return "zip", bytes.NewBuffer(b), nil
		} else {
			return "", nil, errors.OperationUndefinedError(dataset.Spec.FileFormat)
		}

	} else {
		// No compressed file format specified
		if dataset.Spec.FileFormat == "csv" {
			// Generate CSV files and store them in a TAR archive
			buf := new(bytes.Buffer)

			gw := gzip.NewWriter(buf)
			defer gw.Close()

			tw := tar.NewWriter(gw)
			defer tw.Close()

			for i, schemaName := range schemaNames {
				minRec := dataset.Spec.Schemas[i].NumRecords["min"]
				maxRec := dataset.Spec.Schemas[i].NumRecords["max"]

				numRecords := gofakeit.Number(minRec, maxRec)
				bytes, err := datagen.Raw2CSVAtBytesBySchema(numRecords, schemaName)
				if err != nil {
					return "", nil, err
				}
				AddFileToTar(tw, "file"+strconv.Itoa(i)+".csv", bytes)
			}
			return "csv", buf, nil

		} else if dataset.Spec.FileFormat == "binary" {
			// Generate binary files and store them in a TAR archive
			buf := new(bytes.Buffer)

			gw := gzip.NewWriter(buf)
			defer gw.Close()

			tw := tar.NewWriter(gw)
			defer tw.Close()
			for i, schemaName := range schemaNames {
				minRec := dataset.Spec.Schemas[i].NumRecords["min"]
				maxRec := dataset.Spec.Schemas[i].NumRecords["max"]

				numRecords := gofakeit.Number(minRec, maxRec)
				bytes, err := datagen.Raw2BinaryAtBytesBySchema(numRecords, schemaName)
				if err != nil {
					return "", nil, err
				}
				AddFileToTar(tw, "file"+strconv.Itoa(i)+".csv", bytes)
			}

			return "bin", buf, nil

		} else {
			return "", nil, errors.OperationUndefinedError(dataset.Spec.FileFormat)
		}
	}
}

// GetIndexByName retrieves the index of a schema by name from a list of schema selectors.
func GetIndexByName(schemas []plantdv1alpha1.SchemaSelector, schemaName string) (int, bool) {
	for i, schema := range schemas {
		if schema.Name == schemaName {
			return i, true
		}
	}
	return -1, false
}

// AddFileToTar adds a file with the given content to a tar.Writer.
func AddFileToTar(tw *tar.Writer, name string, content []byte) error {
	header := &tar.Header{
		Name: name,
		Mode: 0600,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err := tw.Write(content)
	return err
}

// CheckPipelineHealth checks the health of a pipeline.
func CheckPipelineHealth(ctx context.Context, c client.Client, URL string, HealthCheckEndpoint string) error {

	if HealthCheckEndpoint != "" {
		healthCheckURL, err := utils.GetHealthCheckURL(URL, HealthCheckEndpoint)
		if err != nil {
			return err
		}
		ok, err := utils.HealthCheck(healthCheckURL)
		if err != nil || !ok {
			return err
		}
	}

	return nil
}
