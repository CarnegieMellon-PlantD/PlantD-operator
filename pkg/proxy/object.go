package proxy

import (
	"context"
	"fmt"

	plantdv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Constants defining the plural names of various resources.
const (
	SchemaPlural       string = "schemas"
	DatasetPlural      string = "datasets"
	LoadPatternPlural  string = "loadpatterns"
	PipelinePlural     string = "pipelines"
	ExperimentPlural   string = "experiments"
	PlantDCorePlural   string = "plantdcores"
	CostExporterPlural string = "costexporters"
)

// ForObject returns a client.Object instance based on the provided kind.
func ForObject(kind string) client.Object {
	switch kind {
	case SchemaPlural:
		return &plantdv1alpha1.Schema{}
	case DatasetPlural:
		return &plantdv1alpha1.DataSet{}
	case LoadPatternPlural:
		return &plantdv1alpha1.LoadPattern{}
	case PipelinePlural:
		return &plantdv1alpha1.Pipeline{}
	case ExperimentPlural:
		return &plantdv1alpha1.Experiment{}
	case PlantDCorePlural:
		return &plantdv1alpha1.PlantDCore{}
	case CostExporterPlural:
		return &plantdv1alpha1.CostExporter{}
	}
	return nil
}

// ForObjectList returns a client.ObjectList instance based on the provided kind.
func ForObjectList(kind string) client.ObjectList {
	switch kind {
	case SchemaPlural:
		return &plantdv1alpha1.SchemaList{}
	case DatasetPlural:
		return &plantdv1alpha1.DataSetList{}
	case LoadPatternPlural:
		return &plantdv1alpha1.LoadPatternList{}
	case PipelinePlural:
		return &plantdv1alpha1.PipelineList{}
	case ExperimentPlural:
		return &plantdv1alpha1.ExperimentList{}
	case PlantDCorePlural:
		return &plantdv1alpha1.PlantDCoreList{}
	case CostExporterPlural:
		return &plantdv1alpha1.CostExporterList{}
	}
	return nil
}

// updateSchemaSpec updates the Spec field of a fetched Schema object with the updated Schema object.
func updateSchemaSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.Schema)
	if !ok {
		return fmt.Errorf("could not convert fetched object to schema")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.Schema)
	if !ok {
		return fmt.Errorf("could not convert updated object to schema")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateDatasetSpec updates the Spec field of a fetched DataSet object with the updated DataSet object.
func updateDatasetSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.DataSet)
	if !ok {
		return fmt.Errorf("could not convert fetched object to dataset")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.DataSet)
	if !ok {
		return fmt.Errorf("could not convert updated object to dataset")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateLoadPatternSpec updates the Spec field of a fetched LoadPattern object with the updated LoadPattern object.
func updateLoadPatternSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.LoadPattern)
	if !ok {
		return fmt.Errorf("could not convert fetched object to loadpattern")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.LoadPattern)
	if !ok {
		return fmt.Errorf("could not convert updated object to loadpattern")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updatePipelineSpec updates the Spec field of a fetched Pipeline object with the updated Pipeline object.
func updatePipelineSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.Pipeline)
	if !ok {
		return fmt.Errorf("could not convert fetched object to pipeline")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.Pipeline)
	if !ok {
		return fmt.Errorf("could not convert updated object to pipeline")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateExperimentSpec updates the Spec field of a fetched Experiment object with the updated Experiment object.
func updateExperimentSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.Experiment)
	if !ok {
		return fmt.Errorf("could not convert fetched object to experiment")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.Experiment)
	if !ok {
		return fmt.Errorf("could not convert updated object to experiment")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updatePlantDCoreSpec updates the Spec field of a fetched PlantDCore object with the updated PlantDCore object.
func updatePlantDCoreSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.PlantDCore)
	if !ok {
		return fmt.Errorf("could not convert fetched object to plantdcore")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.PlantDCore)
	if !ok {
		return fmt.Errorf("could not convert updated object to plantdcore")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateCostExporterSpec updates the Spec field of a fetched Schema object with the updated Schema object.
func updateCostExporterSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.CostExporter)
	if !ok {
		return fmt.Errorf("could not convert fetched object to costexporter")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.CostExporter)
	if !ok {
		return fmt.Errorf("could not convert updated object to costexporter")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateSpec updates the Spec field of the fetched object with the updated object based on the provided kind.
func updateSpec(fetched client.Object, updated client.Object, kind string) error {
	switch kind {
	case SchemaPlural:
		return updateSchemaSpec(fetched, updated)
	case DatasetPlural:
		return updateDatasetSpec(fetched, updated)
	case LoadPatternPlural:
		return updateLoadPatternSpec(fetched, updated)
	case PipelinePlural:
		return updatePipelineSpec(fetched, updated)
	case ExperimentPlural:
		return updateExperimentSpec(fetched, updated)
	case PlantDCorePlural:
		return updatePlantDCoreSpec(fetched, updated)
	case CostExporterPlural:
		return updateCostExporterSpec(fetched, updated)
	}
	return nil
}

// GetObjectList retrieves a list of objects of the provided kind.
func GetObjectList(ctx context.Context, c client.Client, kind string) (client.ObjectList, error) {
	objList := ForObjectList(kind)
	if objList == nil {
		return nil, fmt.Errorf("kind \"%s\" not found", kind)
	}

	if err := c.List(ctx, objList); err != nil {
		return nil, err
	}

	return objList, nil
}

// GetObject retrieves an object of the provided kind, namespace, and name.
func GetObject(ctx context.Context, c client.Client, kind, namespace, name string) (client.Object, error) {
	obj := ForObject(kind)
	if obj == nil {
		return nil, fmt.Errorf("kind \"%s\" not found", kind)
	}

	if err := c.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// CreateObject creates a new object of the provided kind.
func CreateObject(ctx context.Context, c client.Client, newObj client.Object) error {
	if err := c.Create(ctx, newObj); err != nil {
		return err
	}

	return nil
}

// UpdateObject updates an existing object of the provided kind.
func UpdateObject(ctx context.Context, c client.Client, updated client.Object, kind string) error {
	fetched := ForObject(kind)
	if fetched == nil {
		return fmt.Errorf("kind \"%s\" not found", kind)
	}
	key := types.NamespacedName{Name: updated.GetName(), Namespace: updated.GetNamespace()}
	if err := c.Get(ctx, key, fetched); err != nil {
		return err
	}

	if err := updateSpec(fetched, updated, kind); err != nil {
		return err
	}

	if err := c.Update(ctx, fetched); err != nil {
		return err
	}

	return nil
}

// DeleteObject deletes an existing object of the provided kind, namespace, and name.
func DeleteObject(ctx context.Context, c client.Client, kind, namespace, name string) error {
	key := types.NamespacedName{Name: name, Namespace: namespace}
	obj := ForObject(kind)
	if obj == nil {
		return fmt.Errorf("kind \"%s\" not found", kind)
	}
	err := c.Get(ctx, key, obj)
	if err != nil {
		return err
	}

	err = c.Delete(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}
