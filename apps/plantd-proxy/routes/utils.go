package routes

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"io"
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

// CheckPipelineHealth returns an HTTP handler function for checking the health of a pipeline.
// It takes a client object of type client.Client for interacting with the Kubernetes API.
// The handler function decodes the base64-encoded URL parameter "info" and checks the pipeline health using the proxy.CheckPipelineHealth function.
// If the pipeline is healthy, it returns an HTTP 200 status.
// If there is an error during the decoding or health check process, it returns an HTTP 500 status with an error message.
func CheckPipelineHealth(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		b64Info := chi.URLParam(r, "info")
		info, err := b64.RawURLEncoding.DecodeString(b64Info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while decoding base64 URL param: " + err.Error()})
			return
		}
		healthCheckMeta := &proxy.HealthCheckMeta{}
		err = json.Unmarshal(info, healthCheckMeta)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling base64 URL param: " + err.Error()})
			return
		}

		if err := proxy.CheckPipelineHealth(ctx, client, healthCheckMeta.URL, healthCheckMeta.HealthCheckEndpoint); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// ExportCustomResources returns an HTTP handler function for exporting custom resources.
// It takes a client object of type client.Client for interacting with the Kubernetes API.
// The handler function decodes the base64-encoded URL parameter "info" and exports the custom resources using the proxy.ExportCustomResources function.
// If the export is successful, it writes the exported resources as a ZIP file to the response.
// If there is an error during the decoding or export process, it returns an HTTP 500 status with an error message.
func ExportCustomResources(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		b64Info := chi.URLParam(r, "info")
		info, err := b64.RawURLEncoding.DecodeString(b64Info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while decoding base64 URL param: " + err.Error()})
			return
		}
		exportResourcesInfo := proxy.ExportResourcesInfo{}
		err = json.Unmarshal(info, &exportResourcesInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling base64 URL param: " + err.Error()})
			return
		}
		zipBuffer, err := proxy.ExportCustomResources(ctx, client, &exportResourcesInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", "attachment; filename=crds.zip")
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write(zipBuffer); err != nil {
				return
			}
		}
	}
}

// ImportCustomResources returns an HTTP handler function for importing custom resources.
// It takes a client object of type client.Client for interacting with the Kubernetes API.
// The handler function reads the uploaded file from the request form data and imports the custom resources using the proxy.ImportCustomResources function.
// If the import is successful, it returns an HTTP 200 status.
// If there is an error during the file reading or import process, it returns an HTTP 500 status with an error message.
func ImportCustomResources(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		file, _, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while parsing request form data: " + err.Error()})
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading file: " + err.Error()})
		}

		if err := proxy.ImportCustomResources(ctx, client, buf.Bytes()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
