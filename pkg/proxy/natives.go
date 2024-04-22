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
	resultCh := make(chan struct{})
	go func() {
		for {
			err := client.Get(ctx, types.NamespacedName{Name: namespaceName}, nsToDelete)
			if err != nil {
				resultCh <- struct{}{}
				break
			}
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-resultCh:
		return nil
	}
}

// ListServices retrieves a list of all Services.
func ListServices(ctx context.Context, client client.Client) (*corev1.ServiceList, error) {
	svcList := &corev1.ServiceList{}

	if err := client.List(ctx, svcList); err != nil {
		return nil, err
	}

	return svcList, nil
}

// ListSecrets retrieves a list of all Secrets.
func ListSecrets(ctx context.Context, client client.Client) (*corev1.SecretList, error) {
	secretList := &corev1.SecretList{}

	if err := client.List(ctx, secretList); err != nil {
		return nil, err
	}

	return secretList, nil
}
