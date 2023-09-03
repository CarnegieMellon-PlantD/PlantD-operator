package routes

import (
	"encoding/json"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

// healthCheck is a handler function for the health check endpoint.
// It returns an HTTP 200 status with the response body "Healthy".
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}

// GetSampleDataset returns an HTTP handler function for retrieving a sample dataset.
// It takes a client object of type client.Client for interacting with the Kubernetes API.
// The handler function retrieves the sample dataset based on the provided namespace and dataset name.
// It calls the proxy.GetSampleDataset function and writes the dataset to the response.
// If there is an error during the retrieval process, it returns an HTTP 500 status with an error message.
func GetSampleDataset(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		datasetName := chi.URLParam(r, "name")
		if fileFormat, bytes, err := proxy.GetSampleDataset(ctx, client, namespace, datasetName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while getting sample DataSet: " + err.Error()})
			return
		} else {
			contentDisposition := "attachment; filename=example." + fileFormat
			w.Header().Set("Content-Disposition", contentDisposition)
			contentType := "application/" + fileFormat
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(http.StatusOK)

			if _, err := bytes.WriteTo(w); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Message: "while writing response body: " + err.Error()})
				return
			}
		}
	}
}
