package routes

import (
	"bytes"
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

// healthCheck is a handler function for the health check endpoint.
// It responds an HTTP 200 status with the response body "Healthy".
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}

// getSampleDataSetHandler returns an HTTP handler function for retrieving a sample dataset.
// It takes a client object of type client.Client for interacting with the Kubernetes API.
// The handler function retrieves the sample dataset based on the provided namespace and dataset name.
// It calls the proxy.GetSampleDataSet function and writes the dataset to the response.
// If there is an error during the retrieval process, it responds an HTTP 500 status with an error message.
func getSampleDataSetHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		if fileExt, bytes, err := proxy.GetSampleDataSet(ctx, client, namespace, name); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: Failed to get sample dataset: %s", err.Error())))
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

// checkHTTPHealthHandler returns an HTTP handler function for checking health status of a URL using HTTP protocol.
// The handler function retrieves the sample dataset based on the provided namespace and dataset name.
// It calls utils.CheckHTTPHealth to make a request to the designated URL. Upon receiving an HTTP non-200 response,
// it responds an HTTP 500 status with an ErrorResponse in JSON.
func checkHTTPHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		data := &proxy.CheckHTTPHealthRequest{}
		err = json.Unmarshal(body, data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}

		_, err = utils.CheckHTTPHealth(data.URL)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// listKindsHandler returns an HTTP handler function for listing all available kinds in custom resource definitions.
// If succeeded, it returns an array of string in JSON format.
func listKindsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(proxy.ListKinds())
	}
}

// listResourcesHandler returns an HTTP handler function for listing all custom resource definitions.
// If succeeded, it returns an array of proxy.ResourceLocator in JSON format.
func listResourcesHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resources, err := proxy.ListResources(ctx, client)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while listing resources: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resources)
	}
}

// importResourcesHandler returns an HTTP handler function for importing custom resource definitions from YAML files.
// The handler function reads the ZIP file content from the `file` field of the request body, which is a form.
// It calls proxy.ImportResources to extract the ZIP file and import each YAML file.
// If it completes successfully or with minor errors, it responds an HTTP 200 status code with a
// proxy.ImportStatistics in JSON format.
// If a fundamental error occurs, it responds a corresponding HTTP status code with a ErrorResponse in JSON.
func importResourcesHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while reading request form: " + err.Error()})
			return
		}
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, file); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while reading file content: " + err.Error()})
			return
		}
		if stat, err := proxy.ImportResources(ctx, client, buf); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stat)
		}
	}
}

// exportResourcesHandler return an HTTP handler function for exporting custom resource definitions to YAML files.
// The handler function gets an array of proxy.ResourceLocator objects from the `info` field of the request body, which is a form.
// It calls proxy.ExportResources to export all the specified objects to YAML files, and return a ZIP file that contains them.
// If successful, it responds an HTTP 200 status code.
// If a fundamental error occurs, it responds a corresponding HTTP status code with a ErrorResponse in plain text, as the
// client browser is expected to access this endpoint via a separate window instead of an AJAX request.
func exportResourcesHandler(client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		info := r.FormValue("info")
		if info == "" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Request form does not contain `info` field"))
			return
		}
		var data []*proxy.ResourceLocator
		if err := json.Unmarshal([]byte(info), &data); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: Failed to unmarshal `info` field: %s", err.Error())))
			return
		}

		bytes, err := proxy.ExportResources(ctx, client, data)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: Failed to export resources: %s", err.Error())))
			return
		}
		contentDisposition := fmt.Sprintf("attachment; filename=%s.zip", time.Now().Format("2006-01-02-15-04-05"))
		w.Header().Set("Content-Disposition", contentDisposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		bytes.WriteTo(w)
	}
}
