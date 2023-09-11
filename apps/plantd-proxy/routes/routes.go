package routes

import (
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Run starts the HTTP server with the provided router and client.
// It listens on port 5000 and handles incoming requests.
// If an error occurs while starting the server, it panics.
func Run(router *chi.Mux, client client.Client, agent proxy.QueryClient) {
	getRoutes(router, client, agent)
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		panic(err)
	}
}

// getRoutes defines the API routes and their corresponding handlers.
// It takes a router object of type *chi.Mux and a client object of type client.Client.
// It registers the routes with the router, associating each route with its respective handler function.
func getRoutes(router *chi.Mux, client client.Client, agent proxy.QueryClient) {
	router.Route("/api", func(r chi.Router) {
		r.Post("/health", healthCheck)

		r.Get("/namespaces", listNamespaces(client))
		r.Post("/namespaces/{namespace}", createNamespace(client))
		r.Delete("/namespaces/{namespace}", deleteNamespace(client))

		r.Get("/schemas", getObjectList(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.SchemaKind))
		r.Get("/schemas/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.SchemaKind))
		r.Post("/schemas/{namespace}/{name}", createObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.SchemaKind))
		r.Put("/schemas/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.SchemaKind))
		r.Delete("/schemas/{namespace}/{name}", deleteObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.SchemaKind))

		r.Get("/datasets", getObjectList(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.DatasetKind))
		r.Get("/datasets/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.DatasetKind))
		r.Post("/datasets/{namespace}/{name}", createObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.DatasetKind))
		r.Put("/datasets/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.DatasetKind))
		r.Delete("/datasets/{namespace}/{name}", deleteObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.DatasetKind))

		r.Get("/loadpatterns", getObjectList(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.LoadPatternKind))
		r.Get("/loadpatterns/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.LoadPatternKind))
		r.Post("/loadpatterns/{namespace}/{name}", createObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.LoadPatternKind))
		r.Put("/loadpatterns/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.LoadPatternKind))
		r.Delete("/loadpatterns/{namespace}/{name}", deleteObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.LoadPatternKind))

		r.Get("/pipelines", getObjectList(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PipelineKind))
		r.Get("/pipelines/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PipelineKind))
		r.Post("/pipelines/{namespace}/{name}", createObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PipelineKind))
		r.Put("/pipelines/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PipelineKind))
		r.Delete("/pipelines/{namespace}/{name}", deleteObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PipelineKind))

		r.Get("/experiments", getObjectList(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.ExperimentKind))
		r.Get("/experiments/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.ExperimentKind))
		r.Post("/experiments/{namespace}/{name}", createObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.ExperimentKind))
		r.Put("/experiments/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.ExperimentKind))
		r.Delete("/experiments/{namespace}/{name}", deleteObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.ExperimentKind))

		r.Get("/costexporters", getObjectList(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.CostExporterKind))
		r.Get("/costexporters/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.CostExporterKind))
		r.Post("/costexporters/{namespace}/{name}", createObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.CostExporterKind))
		r.Put("/costexporters/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.CostExporterKind))
		r.Delete("/costexporters/{namespace}/{name}", deleteObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.CostExporterKind))

		r.Get("/plantdcores/{namespace}/{name}", getObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PlantDCoreKind))
		r.Put("/plantdcores/{namespace}/{name}", updateObject(client, proxy.PlantDGroup, proxy.V1Alpha1Version, proxy.PlantDCoreKind))

		r.Get("/datasets/{namespace}/{name}/sample", getSampleDataSet(client))
		r.Get("/healthcheck/http", checkHTTPHealth())
		r.Get("/export", exportResources(client))
	})

	router.Route("/data", func(r chi.Router) {
		r.Get("/bi-channel", GetQueryHandler(agent, BI_CHANN))
		r.Get("/tri-channel", GetQueryHandler(agent, TRI_CHANN))
	})
}
