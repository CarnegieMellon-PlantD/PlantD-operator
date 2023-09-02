package routes

import (
	"encoding/json"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListNamespaces returns an HTTP handler function that handles GET requests to fetch a list of namespaces.
// The handler function calls the proxy.ListNamespaces function to retrieve the namespaces using the provided client.
// If successful, it encodes the namespace list as JSON and writes it to the response.
// If an error occurs, it writes an error response with the corresponding status code and error message.
func ListNamespaces(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		nsList, err := proxy.ListNamespaces(ctx, client)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nsList)
		}
	}
}

// CreateNamespace returns an HTTP handler function that handles POST requests to create a new namespace.
// The handler function reads the namespace parameter from the request URL.
// It calls the proxy.CreateNamespace function to create the namespace using the provided client and namespace name.
// If the namespace already exists, it writes an error response with the status code 409.
// If the creation fails, it writes an error response with the corresponding status code and error message.
// If successful, it writes a response with the status code 200.
func CreateNamespace(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		if namespaceExistsErr, creationFailedErr := proxy.CreateNamespace(ctx, client, namespace); namespaceExistsErr != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Message: namespaceExistsErr.Error()})
		} else if creationFailedErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: creationFailedErr.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// DeleteNamespace returns an HTTP handler function that handles DELETE requests to delete a namespace.
// The handler function reads the namespace parameter from the request URL.
// It calls the proxy.DeleteNamespace function to delete the namespace using the provided client and namespace name.
// If the namespace is not found, it writes an error response with the status code 404.
// If the deletion fails, it writes an error response with the corresponding status code and error message.
// If successful, it writes a response with the status code 200.
func DeleteNamespace(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		if namespaceNotFoundErr, deletionFailedErr := proxy.DeleteNamespace(ctx, client, namespace); namespaceNotFoundErr != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: namespaceNotFoundErr.Error()})
		} else if deletionFailedErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: deletionFailedErr.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
