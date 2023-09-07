package routes

import (
	"encoding/json"
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
		datasetName := chi.URLParam(r, "name")
		if fileFormat, bytes, err := proxy.GetSampleDataSet(ctx, client, namespace, datasetName); err != nil {
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
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		data := CheckHTTPHealthRequest{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}
		_, err = utils.CheckHTTPHealth(data.URL)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while checking health: " + err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	}
}
