package proxy

import (
	"context"
	"fmt"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/errors"

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

// Helper functions to update the Spec field of different resource types.

// updateCostExporterSpec updates the Spec field of a fetched Schema object with the updated Schema object.
func updateCostExporterSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.CostExporter)
	if !ok {
		return fmt.Errorf("could not convert fetched object to schemas")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.CostExporter)
	if !ok {
		return fmt.Errorf("could not convert updated object to schemas")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateSchemaSpec updates the Spec field of a fetched Schema object with the updated Schema object.
func updateSchemaSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.Schema)
	if !ok {
		return fmt.Errorf("could not convert fetched object to schemas")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.Schema)
	if !ok {
		return fmt.Errorf("could not convert updated object to schemas")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateDatasetSpec updates the Spec field of a fetched DataSet object with the updated DataSet object.
func updateDatasetSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.DataSet)
	if !ok {
		return fmt.Errorf("could not convert fetched object to datasets")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.DataSet)
	if !ok {
		return fmt.Errorf("could not convert updated object to datasets")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateLoadPatternSpec updates the Spec field of a fetched LoadPattern object with the updated LoadPattern object.
func updateLoadPatternSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.LoadPattern)
	if !ok {
		return fmt.Errorf("could not convert fetched object to loadpatterns")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.LoadPattern)
	if !ok {
		return fmt.Errorf("could not convert updated object to loadpatterns")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updatePipelineSpec updates the Spec field of a fetched Pipeline object with the updated Pipeline object.
func updatePipelineSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.Pipeline)
	if !ok {
		return fmt.Errorf("could not convert fetched object to pipelines")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.Pipeline)
	if !ok {
		return fmt.Errorf("could not convert updated object to pipelines")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updateExperimentSpec updates the Spec field of a fetched Experiment object with the updated Experiment object.
func updateExperimentSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.Experiment)
	if !ok {
		return fmt.Errorf("could not convert fetched object to experiments")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.Experiment)
	if !ok {
		return fmt.Errorf("could not convert updated object to experiments")
	}

	fetchedTyped.Spec = updatedTyped.Spec
	return nil
}

// updatePlantDCoreSpec updates the Spec field of a fetched WindTunnelCluster object with the updated WindTunnelCluster object.
func updatePlantDCoreSpec(fetched client.Object, updated client.Object) error {
	fetchedTyped, ok := fetched.(*plantdv1alpha1.PlantDCore)
	if !ok {
		return fmt.Errorf("could not convert fetched object to windtunnelclusters")
	}
	updatedTyped, ok := updated.(*plantdv1alpha1.PlantDCore)
	if !ok {
		return fmt.Errorf("could not convert updated object to windtunnelclusters")
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

// GetObjectList retrieves a list of objects of the provided kind and namespace.
func GetObjectList(ctx context.Context, c client.Client, kind, namespace string) (client.ObjectList, error) {
	list := ForObjectList(kind)
	if list == nil {
		return nil, fmt.Errorf("does not exist kind: %s", kind)
	}
	if err := c.List(ctx, list, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	return list, nil
}

// GetObject retrieves an object of the provided kind, namespace, and name.
func GetObject(ctx context.Context, c client.Client, kind, namespace, name string) (client.Object, error) {
	obj := ForObject(kind)
	if obj == nil {
		return nil, fmt.Errorf("does not exist kind: %s", kind)
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
func CreateObject(ctx context.Context, c client.Client, newObj client.Object, kind string) (objectExists, creationFailed error) {

	key := types.NamespacedName{Name: newObj.GetName(), Namespace: newObj.GetNamespace()}

	obj := ForObject(kind)
	if obj == nil {
		return nil, fmt.Errorf("kind not found: %s", kind)
	}
	err := c.Get(ctx, key, obj)
	if err == nil {
		return errors.DuplicateIDError(newObj.GetName()), nil
	}

	if err := c.Create(ctx, newObj); err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateObject updates an existing object of the provided kind.
func UpdateObject(ctx context.Context, c client.Client, updated client.Object, kind string) (notFound, updationFailed error) {

	fetched := ForObject(kind)
	if fetched == nil {
		return nil, fmt.Errorf("does not exist kind: %s", kind)
	}
	key := types.NamespacedName{Name: updated.GetName(), Namespace: updated.GetNamespace()}
	if err := c.Get(ctx, key, fetched); err != nil {
		return err, nil
	}

	if err := updateSpec(fetched, updated, kind); err != nil {
		return nil, err
	}

	if err := c.Update(ctx, fetched); err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteObject deletes an existing object of the provided kind, namespace, and name.
func DeleteObject(ctx context.Context, c client.Client, kind, namespace, name string) (objectNotFound, deletionFailed error) {

	key := types.NamespacedName{Name: name, Namespace: namespace}

	obj := ForObject(kind)
	if obj == nil {
		return nil, fmt.Errorf("does not exist kind: %s", kind)
	}
	err := c.Get(ctx, key, obj)
	if err != nil {
		return err, nil
	}

	err = c.Delete(ctx, obj)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
