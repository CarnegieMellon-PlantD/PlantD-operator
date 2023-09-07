package main

import (
	"context"
	"net/http"

	datav1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/apps/plantd-proxy/routes"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func MethodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get overriding method from header
		methodToOverride := r.Header.Get("X-Http-Method-Override")
		// When method from header is not empty and original method is POST
		if methodToOverride != "" && r.Method == http.MethodPost {
			// Change the method in the context
			newCtx := chi.NewRouteContext()
			newCtx.RouteMethod = methodToOverride
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, newCtx))
		}
		// Go to next handler
		next.ServeHTTP(w, r)
	})
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(datav1alpha1.AddToScheme(scheme))
}

func main() {
	r := chi.NewRouter()
	r.Use(MethodOverrideMiddleware)
	r.Use(middleware.Logger)

	cfg, err := ctrl.GetConfig()
	if err != nil {
		panic(err)
	}

	client, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		panic(err)
	}

	agent, err := proxy.NewQueryAgent(config.GetString("database.prometheus.url"))
	if err != nil {
		panic(err)
	}
	routes.Run(r, client, agent)
}
