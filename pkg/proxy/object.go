package proxy

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetObjectList retrieves a list of objects of the provided GVK.
func GetObjectList(ctx context.Context, c client.Client, group, version, kind string) (client.ObjectList, error) {
	objList := &unstructured.UnstructuredList{}
	objList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	if err := c.List(ctx, objList); err != nil {
		return nil, err
	}

	return objList, nil
}

// GetObject retrieves an object of the provided GVK, namespace, name.
func GetObject(ctx context.Context, c client.Client, group, version, kind, namespace, name string) (client.Object, error) {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	if err := c.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// CreateObject creates a new object of the provided GVK, namespace, name.
// Only the spec field of the obj parameter will be used.
func CreateObject(ctx context.Context, c client.Client, group, version, kind, namespace, name string, obj client.Object) error {
	objToCreate := &unstructured.Unstructured{}
	objToCreate.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	objToCreate.SetNamespace(namespace)
	objToCreate.SetName(name)

	// Only copy the spec field from the obj param
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}
	objToCreate.Object["spec"] = unstructuredObj["spec"]

	return c.Create(ctx, objToCreate)
}

// UpdateObject updates an existing object of the provided GVK, namespace, name.
// Only the spec field of the obj parameter will be used.
func UpdateObject(ctx context.Context, c client.Client, group, version, kind, namespace, name string, obj client.Object) error {
	objToUpdate := &unstructured.Unstructured{}
	objToUpdate.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	// Note: we need to get the full object based on the specified GVK, namespace and name, so that other fields such as
	// .metadata.resourceVersion will present, otherwise an error will occur when updating the object.
	if err := c.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, objToUpdate); err != nil {
		return err
	}

	// Only copy the spec field from the obj param
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}
	objToUpdate.Object["spec"] = unstructuredObj["spec"]

	return c.Update(ctx, objToUpdate)
}

// DeleteObject deletes an existing object of the provided GVK, namespace, name.
func DeleteObject(ctx context.Context, c client.Client, group, version, kind, namespace, name string) error {
	objToDelete := &unstructured.Unstructured{}
	objToDelete.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	if err := c.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, objToDelete); err != nil {
		return err
	}

	return c.Delete(ctx, objToDelete)
}
