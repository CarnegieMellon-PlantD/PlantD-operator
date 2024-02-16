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
func Run(router *chi.Mux, client client.Client, queryAgent *proxy.QueryAgent) {
	getRoutes(router, client, queryAgent)
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		panic(err)
	}
}

// getRoutes defines the API routes and their corresponding handlers.
// It takes a router object of type *chi.Mux and a client object of type client.Client.
// It registers the routes with the router, associating each route with its respective handler function.
func getRoutes(router *chi.Mux, client client.Client, queryAgent *proxy.QueryAgent) {
	router.Route("/api", func(r chi.Router) {
		r.Post("/health", healthCheck)

		r.Get("/namespaces", listNamespacesHandler(client))
		r.Post("/namespaces/{namespace}", createNamespaceHandler(client))
		r.Delete("/namespaces/{namespace}", deleteNamespaceHandler(client))

		r.Get("/schemas", getObjectListHandler(client, proxy.SchemaKind))
		r.Get("/schemas/{namespace}/{name}", getObjectHandler(client, proxy.SchemaKind))
		r.Post("/schemas/{namespace}/{name}", createObjectHandler(client, proxy.SchemaKind))
		r.Put("/schemas/{namespace}/{name}", updateObjectHandler(client, proxy.SchemaKind))
		r.Delete("/schemas/{namespace}/{name}", deleteObjectHandler(client, proxy.SchemaKind))

		r.Get("/datasets", getObjectListHandler(client, proxy.DatasetKind))
		r.Get("/datasets/{namespace}/{name}", getObjectHandler(client, proxy.DatasetKind))
		r.Post("/datasets/{namespace}/{name}", createObjectHandler(client, proxy.DatasetKind))
		r.Put("/datasets/{namespace}/{name}", updateObjectHandler(client, proxy.DatasetKind))
		r.Delete("/datasets/{namespace}/{name}", deleteObjectHandler(client, proxy.DatasetKind))

		r.Get("/loadpatterns", getObjectListHandler(client, proxy.LoadPatternKind))
		r.Get("/loadpatterns/{namespace}/{name}", getObjectHandler(client, proxy.LoadPatternKind))
		r.Post("/loadpatterns/{namespace}/{name}", createObjectHandler(client, proxy.LoadPatternKind))
		r.Put("/loadpatterns/{namespace}/{name}", updateObjectHandler(client, proxy.LoadPatternKind))
		r.Delete("/loadpatterns/{namespace}/{name}", deleteObjectHandler(client, proxy.LoadPatternKind))

		r.Get("/pipelines", getObjectListHandler(client, proxy.PipelineKind))
		r.Get("/pipelines/{namespace}/{name}", getObjectHandler(client, proxy.PipelineKind))
		r.Post("/pipelines/{namespace}/{name}", createObjectHandler(client, proxy.PipelineKind))
		r.Put("/pipelines/{namespace}/{name}", updateObjectHandler(client, proxy.PipelineKind))
		r.Delete("/pipelines/{namespace}/{name}", deleteObjectHandler(client, proxy.PipelineKind))

		r.Get("/experiments", getObjectListHandler(client, proxy.ExperimentKind))
		r.Get("/experiments/{namespace}/{name}", getObjectHandler(client, proxy.ExperimentKind))
		r.Post("/experiments/{namespace}/{name}", createObjectHandler(client, proxy.ExperimentKind))
		r.Put("/experiments/{namespace}/{name}", updateObjectHandler(client, proxy.ExperimentKind))
		r.Delete("/experiments/{namespace}/{name}", deleteObjectHandler(client, proxy.ExperimentKind))

		r.Get("/costexporters", getObjectListHandler(client, proxy.CostExporterKind))
		r.Get("/costexporters/{namespace}/{name}", getObjectHandler(client, proxy.CostExporterKind))
		r.Post("/costexporters/{namespace}/{name}", createObjectHandler(client, proxy.CostExporterKind))
		r.Put("/costexporters/{namespace}/{name}", updateObjectHandler(client, proxy.CostExporterKind))
		r.Delete("/costexporters/{namespace}/{name}", deleteObjectHandler(client, proxy.CostExporterKind))

		r.Get("/trafficmodels", getObjectListHandler(client, proxy.TrafficModelKind))
		r.Get("/trafficmodels/{namespace}/{name}", getObjectHandler(client, proxy.TrafficModelKind))
		r.Post("/trafficmodels/{namespace}/{name}", createObjectHandler(client, proxy.TrafficModelKind))
		r.Put("/trafficmodels/{namespace}/{name}", updateObjectHandler(client, proxy.TrafficModelKind))
		r.Delete("/trafficmodels/{namespace}/{name}", deleteObjectHandler(client, proxy.TrafficModelKind))

		r.Get("/digitaltwins", getObjectListHandler(client, proxy.DigitalTwinKind))
		r.Get("/digitaltwins/{namespace}/{name}", getObjectHandler(client, proxy.DigitalTwinKind))
		r.Post("/digitaltwins/{namespace}/{name}", createObjectHandler(client, proxy.DigitalTwinKind))
		r.Put("/digitaltwins/{namespace}/{name}", updateObjectHandler(client, proxy.DigitalTwinKind))
		r.Delete("/digitaltwins/{namespace}/{name}", deleteObjectHandler(client, proxy.DigitalTwinKind))

		r.Get("/simulations", getObjectListHandler(client, proxy.SimulationKind))
		r.Get("/simulations/{namespace}/{name}", getObjectHandler(client, proxy.SimulationKind))
		r.Post("/simulations/{namespace}/{name}", createObjectHandler(client, proxy.SimulationKind))
		r.Put("/simulations/{namespace}/{name}", updateObjectHandler(client, proxy.SimulationKind))
		r.Delete("/simulations/{namespace}/{name}", deleteObjectHandler(client, proxy.SimulationKind))

		r.Get("/plantdcores/{namespace}/{name}", getObjectHandler(client, proxy.PlantDCoreKind))
		r.Put("/plantdcores/{namespace}/{name}", updateObjectHandler(client, proxy.PlantDCoreKind))

		r.Get("/datasets/sample/{namespace}/{name}", getSampleDataSetHandler(client))
		r.Get("/health/http", checkHTTPHealthHandler())
		r.Get("/kinds", listKindsHandler())
		r.Get("/resources", listResourcesHandler(client))
		r.Post("/resources/import", importResourcesHandler(client))
		// We are violating RESTful API design principles and using POST instead of GET here,
		// because we want to accept a request body.
		r.Post("/resources/export", exportResourcesHandler(client))
	})

	router.Route("/data", func(r chi.Router) {
		r.Get("/raw/redis", queryRawHandler(queryAgent, proxy.Redis))

		r.Get("/bi-channel/prometheus", queryHandler(queryAgent, proxy.BiChan, proxy.Prometheus))
		r.Get("/bi-channel/redis-ts", queryHandler(queryAgent, proxy.BiChan, proxy.RedisTimeSeries))

		r.Get("/tri-channel/prometheus", queryHandler(queryAgent, proxy.TriChan, proxy.Prometheus))
		r.Get("/tri-channel/redis-ts", queryHandler(queryAgent, proxy.TriChan, proxy.RedisTimeSeries))
	})
}
