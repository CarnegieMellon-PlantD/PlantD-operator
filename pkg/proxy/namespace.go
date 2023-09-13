package proxy

import (
	"context"
	"time"

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
func CreateNamespace(ctx context.Context, client client.Client, namespaceName string) error {
	newNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	err := client.Create(ctx, newNs)
	if err != nil {
		return err
	}

	return nil
}

// DeleteNamespace deletes the namespace with the provided name.
func DeleteNamespace(ctx context.Context, client client.Client, namespaceName string) error {
	nsToDelete := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	if err := client.Delete(ctx, nsToDelete); err != nil {
		return err
	}

	// Wait until the namespace is completely deleted
	for {
		err := client.Get(ctx, types.NamespacedName{Name: namespaceName}, nsToDelete)
		if err != nil {
			return nil
		}
		time.Sleep(time.Second)
	}
}
