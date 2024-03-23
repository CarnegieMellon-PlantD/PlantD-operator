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
	jobIndex, _ := strconv.Atoi(os.Getenv("JOB_COMPLETION_INDEX"))
	stepSize, _ := strconv.Atoi(os.Getenv("JOB_STEP_SIZE"))
	maxRepeat, _ := strconv.Atoi(os.Getenv("MAX_REPEAT"))

	repeatStart := jobIndex * stepSize
	repeatEnd := min(repeatStart+stepSize, maxRepeat)

	dataGeneratorNamespace := os.Getenv("DG_NAMESPACE")
	dataGeneratorName := os.Getenv("DG_NAME")

	datasetString := os.Getenv("DATASET")

	var dataset windtunnelv1alpha1.DataSet
	err := json.Unmarshal([]byte(datasetString), &dataset)
	if err != nil {
		log.Panic(err)
	}

	schemaMapString := os.Getenv("SCHEMA_MAP")

	var schemaMap map[string]*windtunnelv1alpha1.Schema
	err = json.Unmarshal([]byte(schemaMapString), &schemaMap)
	if err != nil {
		log.Panic(err)
	}

	job := datagen.NewBuildBasedDataGeneratorJob(repeatStart, repeatEnd, dataGeneratorNamespace, dataGeneratorName, &dataset, schemaMap)
	err = job.GenerateData()
	if err != nil {
		log.Panic(err)
	}
}
