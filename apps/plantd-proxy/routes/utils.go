package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/utils"

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

// getSampleDataSet returns an HTTP handler function for retrieving a sample dataset.
// It takes a client object of type client.Client for interacting with the Kubernetes API.
// The handler function retrieves the sample dataset based on the provided namespace and dataset name.
// It calls the proxy.GetSampleDataSet function and writes the dataset to the response.
// If there is an error during the retrieval process, it returns an HTTP 500 status with an error message.
func getSampleDataSet(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		if fileExt, bytes, err := proxy.GetSampleDataSet(ctx, client, namespace, name); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while getting sample dataset: " + err.Error()})
			return
		} else {
			contentDisposition := fmt.Sprintf("attachment; filename=sample-%s-%s.%s", namespace, name, fileExt)
			w.Header().Set("Content-Disposition", contentDisposition)
			contentType := "application/octet-stream"
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(http.StatusOK)
			bytes.WriteTo(w)
		}
	}
}

type CheckHTTPHealthRequest struct {
	URL string `json:"url,omitempty"`
}

// checkHTTPHealth returns an HTTP handler function for checking health status of a URL using HTTP protocol.
// The handler function retrieves the sample dataset based on the provided namespace and dataset name.
// It calls utils.CheckHTTPHealth to make a request to the designated URL. Upon receiving an HTTP non-200 response,
// it returns an HTTP 500 status with an error message.
func checkHTTPHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		data := CheckHTTPHealthRequest{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}
		_, err = utils.CheckHTTPHealth(data.URL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}
}
