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

		r.Get("/namespaces", ListNamespaces(client))
		r.Post("/namespaces/{namespace}", CreateNamespace(client))
		r.Delete("/namespaces/{namespace}", DeleteNamespace(client))

		r.Get("/schemas", getObjectList(client, proxy.SchemaPlural))
		r.Get("/schemas/{namespace}/{name}", getObject(client, proxy.SchemaPlural))
		r.Post("/schemas/{namespace}/{name}", createObject(client, proxy.SchemaPlural))
		r.Put("/schemas/{namespace}/{name}", updateObject(client, proxy.SchemaPlural))
		r.Delete("/schemas/{namespace}/{name}", deleteObject(client, proxy.SchemaPlural))

		r.Get("/datasets", getObjectList(client, proxy.DatasetPlural))
		r.Get("/datasets/{namespace}/{name}", getObject(client, proxy.DatasetPlural))
		r.Post("/datasets/{namespace}/{name}", createObject(client, proxy.DatasetPlural))
		r.Put("/datasets/{namespace}/{name}", updateObject(client, proxy.DatasetPlural))
		r.Delete("/datasets/{namespace}/{name}", deleteObject(client, proxy.DatasetPlural))

		r.Get("/loadpatterns", getObjectList(client, proxy.LoadPatternPlural))
		r.Get("/loadpatterns/{namespace}/{name}", getObject(client, proxy.LoadPatternPlural))
		r.Post("/loadpatterns/{namespace}/{name}", createObject(client, proxy.LoadPatternPlural))
		r.Put("/loadpatterns/{namespace}/{name}", updateObject(client, proxy.LoadPatternPlural))
		r.Delete("/loadpatterns/{namespace}/{name}", deleteObject(client, proxy.LoadPatternPlural))

		r.Get("/pipelines", getObjectList(client, proxy.PipelinePlural))
		r.Get("/pipelines/{namespace}/{name}", getObject(client, proxy.PipelinePlural))
		r.Post("/pipelines/{namespace}/{name}", createObject(client, proxy.PipelinePlural))
		r.Put("/pipelines/{namespace}/{name}", updateObject(client, proxy.PipelinePlural))
		r.Delete("/pipelines/{namespace}/{name}", deleteObject(client, proxy.PipelinePlural))

		r.Get("/experiments", getObjectList(client, proxy.ExperimentPlural))
		r.Get("/experiments/{namespace}/{name}", getObject(client, proxy.ExperimentPlural))
		r.Post("/experiments/{namespace}/{name}", createObject(client, proxy.ExperimentPlural))
		r.Put("/experiments/{namespace}/{name}", updateObject(client, proxy.ExperimentPlural))
		r.Delete("/experiments/{namespace}/{name}", deleteObject(client, proxy.ExperimentPlural))

		r.Get("/costexporters", getObjectList(client, proxy.CostExporterPlural))
		r.Get("/costexporters/{namespace}/{name}", getObject(client, proxy.CostExporterPlural))
		r.Post("/costexporters/{namespace}/{name}", createObject(client, proxy.CostExporterPlural))
		r.Put("/costexporters/{namespace}/{name}", updateObject(client, proxy.CostExporterPlural))
		r.Delete("/costexporters/{namespace}/{name}", deleteObject(client, proxy.CostExporterPlural))

		r.Get("/plantdcores/{namespace}/{name}", getObject(client, proxy.PlantDCorePlural))
		r.Put("/plantdcores/{namespace}/{name}", updateObject(client, proxy.PlantDCorePlural))

		r.Get("/datasets/{namespace}/{name}/sample", GetSampleDataset(client))
		r.Get("/pipelines/healthCheck/{info}", CheckPipelineHealth(client))
		r.Get("/export/{info}", ExportCustomResources(client))
		r.Post("/import", ImportCustomResources(client))
	})

	router.Route("/data", func(r chi.Router) {
		r.Get("/bi-channel", GetQueryHandler(agent, BI_CHANN))
		r.Get("/tri-channel", GetQueryHandler(agent, TRI_CHANN))
	})
}
