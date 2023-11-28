from cost.factory.cost_exporter_factory import CostExporterFactory
import os

CLOUD_SERVICE_PROVIDER = os.environ.get('CLOUD_SERVICE_PROVIDER', 'aws')


def collect_cost_logs():
  cost_service = CostExporterFactory().create(CLOUD_SERVICE_PROVIDER)
  cost_service.get_cost_logs()


if __name__ == "__main__":
  collect_cost_logs()
