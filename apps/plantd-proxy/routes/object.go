package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getObjectList retrieves a list of objects of a specified kind and namespace.
// It returns an HTTP handler function that handles GET requests to fetch the object list.
// The handler function reads the namespace parameter from the request URL query.
// It calls the proxy.GetObjectList function to fetch the object list using the provided client and kind.
// If successful, it encodes the list as JSON and writes it to the response.
// If an error occurs, it writes an error response with the corresponding status code and error message.
func getObjectList(client client.Client, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := r.URL.Query().Get("namespace")
		list, err := proxy.GetObjectList(ctx, client, kind, namespace)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(list)
		}
	}
}

// getObject retrieves a single object of a specified kind, namespace, and name.
// It returns an HTTP handler function that handles GET requests to fetch the object.
// The handler function reads the namespace and name parameters from the request URL.
// It calls the proxy.GetObject function to fetch the object using the provided client, kind, namespace, and name.
// If successful, it encodes the object as JSON and writes it to the response.
// If the object is not found, it writes an error response with the status code 404.
// If an error occurs, it writes an error response with the corresponding status code and error message.
func getObject(client client.Client, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		obj, err := proxy.GetObject(ctx, client, kind, namespace, name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(obj)
		}
	}
}

// createObject creates a new object of a specified kind, namespace, and name.
// It returns an HTTP handler function that handles POST requests to create the object.
// The handler function reads the namespace and name parameters from the request URL.
// It reads the request body and unmarshals it into a new object of the specified kind.
// It then sets the object's namespace and name based on the URL parameters.
// It calls the proxy.CreateObject function to create the object using the provided client and kind.
// If the object already exists, it writes an error response with the status code 409.
// If the creation fails, it writes an error response with the corresponding status code and error message.
// If successful, it writes a response with the status code 200.
func createObject(client client.Client, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		obj := proxy.ForObject(kind)
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		err = json.Unmarshal([]byte(body), obj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}
		obj.SetName(name)
		obj.SetNamespace(namespace)
		if objectExistsErr, creationFailedErr := proxy.CreateObject(ctx, client, obj, kind); objectExistsErr != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Message: objectExistsErr.Error()})
		} else if creationFailedErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: creationFailedErr.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// updateObject updates an existing object of a specified kind, namespace, and name.
// It returns an HTTP handler function that handles PUT requests to update the object.
// The handler function reads the namespace and name parameters from the request URL.
// It reads the request body and unmarshals it into an existing object of the specified kind.
// It then sets the object's namespace and name based on the URL parameters.
// It calls the proxy.UpdateObject function to update the object using the provided client and kind.
// If the object is not found, it writes an error response with the status code 404.
// If the updation fails, it writes an error response with the corresponding status code and error message.
// If successful, it writes a response with the status code 200.
func updateObject(client client.Client, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		obj := proxy.ForObject(kind)
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		err = json.Unmarshal([]byte(body), obj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}
		obj.SetName(name)
		obj.SetNamespace(namespace)
		if objectNotFoundErr, updationFailedErr := proxy.UpdateObject(ctx, client, obj, kind); objectNotFoundErr != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: objectNotFoundErr.Error()})
		} else if updationFailedErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: updationFailedErr.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// deleteObject deletes an existing object of a specified kind, namespace, and name.
// It returns an HTTP handler function that handles DELETE requests to delete the object.
// The handler function reads the namespace and name parameters from the request URL.
// It calls the proxy.DeleteObject function to delete the object using the provided client, kind, namespace, and name.
// If the object is not found, it writes an error response with the status code 404.
// If the deletion fails, it writes an error response with the corresponding status code and error message.
// If successful, it writes a response with the status code 200.
func deleteObject(client client.Client, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		if objectNotFoundErr, deletionFailedErr := proxy.DeleteObject(ctx, client, kind, namespace, name); objectNotFoundErr != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: objectNotFoundErr.Error()})
		} else if deletionFailedErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: deletionFailedErr.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
