package proxy

import (
	"context"
	"time"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListNamespaces retrieves a list of all namespaces.
func ListNamespaces(ctx context.Context, client client.Client) (*corev1.NamespaceList, error) {
	nsList := &corev1.NamespaceList{}

	if err := client.List(ctx, nsList); err != nil {
		return nil, err
	}

	return nsList, nil
}

// CreateNamespace creates a new namespace with the provided name.
func CreateNamespace(ctx context.Context, client client.Client, namespaceName string) (namespaceExists, creationFailed error) {

	namespace := &corev1.Namespace{}
	err := client.Get(ctx, types.NamespacedName{Name: namespaceName}, namespace)
	if err == nil {
		return errors.DuplicateIDError(namespaceName), nil
	}

	newNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	err = client.Create(ctx, newNamespace)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteNamespace deletes the namespace with the provided name.
func DeleteNamespace(ctx context.Context, client client.Client, namespaceName string) (namespaceNotFound, deletionFailed error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	err := client.Get(ctx, types.NamespacedName{Name: namespaceName}, namespace)
	if err != nil {
		return err, nil
	}

	if err := client.Delete(ctx, namespace); err != nil {
		return nil, err
	}

	for {
		err := client.Get(ctx, types.NamespacedName{Name: namespaceName}, namespace)
		if err != nil {
			return nil, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			time.Sleep(time.Second)
		}
	}
}
