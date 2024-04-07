package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/datagen"
)

func main() {
	// Environment variable provided by the Kubernetes if the Job is indexed
	// See more information at https://kubernetes.io/docs/concepts/workloads/controllers/job/#completion-mode
	jobIndex, err := strconv.Atoi(os.Getenv("JOB_COMPLETION_INDEX"))
	if err != nil {
		log.Panic(err)
	}
	jobSize, err := strconv.Atoi(os.Getenv("JOB_STEP_SIZE"))
	if err != nil {
		log.Panic(err)
	}
	totalRepeat, err := strconv.Atoi(os.Getenv("TOTAL_REPEAT"))
	if err != nil {
		log.Panic(err)
	}

	repeatStart := jobIndex * jobSize
	repeatEnd := min(repeatStart+jobSize, totalRepeat)

	dataSetString := os.Getenv("DATASET")
	var dataSet windtunnelv1alpha1.DataSet
	if err := json.Unmarshal([]byte(dataSetString), &dataSet); err != nil {
		log.Panic(err)
	}

	schemaMapString := os.Getenv("SCHEMA_MAP")
	var schemaMap map[string]*windtunnelv1alpha1.Schema
	if err := json.Unmarshal([]byte(schemaMapString), &schemaMap); err != nil {
		log.Panic(err)
	}

	path := os.Getenv("OUTPUT_PATH")

	job := datagen.NewBuilderBasedDataGeneratorJob(repeatStart, repeatEnd, &dataSet, schemaMap)
	if err := job.GenerateData(path); err != nil {
		log.Panic(err)
	}
}
