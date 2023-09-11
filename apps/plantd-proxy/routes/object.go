package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getObjectList retrieves a list of objects of a specified GVK.
// It returns an HTTP handler function that handles requests to fetch the object list.
// It calls the proxy.GetObjectList function to fetch the object list using the provided client and GVK.
// If successful, it responds an HTTP 200 status code with an object list in JSON.
// If an error occurs, it responds an HTTP 500 status code with an ErrorResponse in JSON.
func getObjectList(client client.Client, group, version, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		objList, err := proxy.GetObjectList(ctx, client, group, version, kind)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(objList)
		}
	}
}

// getObject retrieves a single object of a specified GVK, namespace, and name.
// It returns an HTTP handler function that handles requests to fetch the object.
// The handler function reads the namespace and name parameters from the request URL.
// It calls the proxy.GetObject function to fetch the object using the provided client, GVK, namespace, and name.
// If successful, it responds an HTTP 200 status code with an object in JSON.
// If an error occurs, it responds an HTTP 500 status code with an ErrorResponse in JSON.
func getObject(client client.Client, group, version, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")

		obj, err := proxy.GetObject(ctx, client, proxy.PlantDGroup, proxy.V1Alpha1Version, kind, namespace, name)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(obj)
		}
	}
}

// createObject creates a new object of a specified GVK, namespace, and name.
// It returns an HTTP handler function that handles requests to create the object.
// The handler function reads the namespace and name parameters from the request URL.
// It reads the request body and unmarshalls it into an object of the specified GVK.
// It calls the proxy.CreateObject function to create the object using the provided client, GVK, namespace, and name.
// If successful, it responds an HTTP 201 status code.
// If an error occurs while creating the object of the specified GVK, reading or unmarshalling the request body,
// it responds an HTTP 400 status code with an ErrorResponse in JSON.
// If an error occurs in other stages, it responds an HTTP 500 status code with an ErrorResponse in JSON.
func createObject(client client.Client, group, version, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")

		// Note: though in proxy/object.go, we use unstructured objects to have more flexibility,
		// we still need to use typed objects here, because the TypeMeta and ObjectMeta can be determined by the URL
		// and the request body may only contain the spec field. Without other fields, it will cause an error when
		// unmarshalling to an unstructured object and can only be unmarshalled to a typed object.
		obj, err := proxy.ForObject(group, version, kind)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		err = json.Unmarshal(body, obj)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}

		if err := proxy.CreateObject(ctx, client, proxy.PlantDGroup, proxy.V1Alpha1Version, kind, namespace, name, obj); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

// updateObject updates an existing object of a specified GVK, namespace, and name.
// It returns an HTTP handler function that handles requests to update the object.
// The handler function reads the namespace and name parameters from the request URL.
// It reads the request body and unmarshalls it into an object of the specified kind.
// It calls the proxy.UpdateObject function to update the object using the provided client, GVK, namespace, name.
// If successful, it responds an HTTP 200 status code.
// If an error occurs while creating the object of the specified GVK, reading or unmarshalling the request body,
// it responds an HTTP 400 status code with an ErrorResponse in JSON.
// If an error occurs in other stages, it responds an HTTP 500 status code with an ErrorResponse in JSON.
func updateObject(client client.Client, group, version, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")

		// Note: though in proxy/object.go, we use unstructured objects to have more flexibility,
		// we still need to use typed objects here, because the TypeMeta and ObjectMeta can be determined by the URL
		// and the request body may only contain the spec field. Without other fields, it will cause an error when
		// unmarshalling to an unstructured object and can only be unmarshalled to a typed object.
		obj, err := proxy.ForObject(group, version, kind)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		err = json.Unmarshal(body, obj)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}

		if err := proxy.UpdateObject(ctx, client, group, version, kind, namespace, name, obj); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// deleteObject deletes an existing object of a specified GVK, namespace, and name.
// It returns an HTTP handler function that handles requests to delete the object.
// The handler function reads the namespace and name parameters from the request URL.
// It calls the proxy.DeleteObject function to delete the object using the provided client, GVK, namespace, and name.
// If successful, it responds an HTTP 200 status code.
// If an error occurs, it responds an HTTP 500 status code with an ErrorResponse in JSON.
func deleteObject(client client.Client, group, version, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")

		if err := proxy.DeleteObject(ctx, client, group, version, kind, namespace, name); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
