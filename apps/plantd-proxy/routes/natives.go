package routes

import (
	"encoding/json"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// listNamespacesHandler returns an HTTP handler function that handles GET requests to fetch a list of namespaces.
// The handler function calls the proxy.ListNamespaces function to retrieve the namespaces using the provided client.
// If successful, it encodes the namespace list as JSON and writes it to the response.
// If an error occurs, it writes an error response with an HTTP 500 status code and error message.
func listNamespacesHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		nsList, err := proxy.ListNamespaces(ctx, client)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nsList)
		}
	}
}

// createNamespaceHandler returns an HTTP handler function that handles POST requests to create a new namespace.
// The handler function reads the namespace parameter from the request URL.
// It calls the proxy.CreateNamespace function to create the namespace using the provided client and namespace name.
// If the creation fails, it writes an error response with an HTTP 500 status code and error message.
// If successful, it writes a response with the status code 200.
func createNamespaceHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		if err := proxy.CreateNamespace(ctx, client, namespace); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	}
}

// deleteNamespaceHandler returns an HTTP handler function that handles DELETE requests to delete a namespace.
// The handler function reads the namespace parameter from the request URL.
// It calls the proxy.DeleteNamespace function to delete the namespace using the provided client and namespace name.
// If the deletion fails, it writes an error response with an HTTP 500 status code and error message.
// If successful, it writes a response with the status code 200.
func deleteNamespaceHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		if err := proxy.DeleteNamespace(ctx, client, namespace); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	}
}

// listServicesHandler returns an HTTP handler function that handles GET requests to fetch a list of Services.
// The handler function calls the proxy.ListServices function to retrieve the Services using the provided client.
// If successful, it encodes the Service list as JSON and writes it to the response.
// If an error occurs, it writes an error response with an HTTP 500 status code and error message.
func listServicesHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		svcList, err := proxy.ListServices(ctx, client)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(svcList)
		}
	}
}

// listSecretsHandler returns an HTTP handler function that handles GET requests to fetch a list of Secrets.
// The handler function calls the proxy.ListSecrets function to retrieve the Secrets using the provided client.
// If successful, it encodes the Secret list as JSON and writes it to the response.
// If an error occurs, it writes an error response with an HTTP 500 status code and error message.
func listSecretsHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		secretList, err := proxy.ListSecrets(ctx, client)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(secretList)
		}
	}
}

// getSecretHandler returns an HTTP handler function that handles GET requests to fetch a Secret.
// It returns an HTTP handler function that handles requests to fetch the Secret.
// The handler function reads the namespace and name parameters from the request URL.
// It calls the proxy.GetSecret function to fetch the Secret using the provided client.
// If successful, it responds an HTTP 200 status code with an Secret in JSON.
// If an error occurs, it responds an HTTP 500 status code with an ErrorResponse in JSON.
func getSecretHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")

		obj, err := proxy.GetSecret(ctx, client, namespace, name)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(obj)
		}
	}
}
