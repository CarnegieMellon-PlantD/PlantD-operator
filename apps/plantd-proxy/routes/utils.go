package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while getting sample dataset: " + err.Error()})
			return
		} else {
			contentDisposition := fmt.Sprintf("attachment; filename=sample_%s_%s.%s", namespace, name, fileExt)
			w.Header().Set("Content-Disposition", contentDisposition)
			w.Header().Set("Content-Type", "application/octet-stream")
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		data := CheckHTTPHealthRequest{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}

		_, err = utils.CheckHTTPHealth(data.URL)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type ExportResourcesRequest struct {
	Items []proxy.ResourceInfo `json:"items,omitempty"`
}

// exportResources return an HTTP handler function for exporting custom resource definitions to YAML files.
// The handler function receives an array of proxy.ResourceInfo objects as the request, calls proxy.ExportResources
// to export all the specified objects to YAML files, and return a ZIP file that contains them.
// When any error occurs, it returns an HTTP 500 status code and error message.
func exportResources(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		data := ExportResourcesRequest{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}

		bytes, err := proxy.ExportResources(ctx, client, &data.Items)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		contentDisposition := fmt.Sprintf("attachment; filename=%s.zip", time.Now().Format("2006-01-02-15-04-05"))
		w.Header().Set("Content-Disposition", contentDisposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		bytes.WriteTo(w)
	}
}
