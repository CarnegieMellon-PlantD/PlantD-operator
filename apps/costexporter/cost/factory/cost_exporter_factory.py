from cost.exporters.aws_cost_exporter import AWSCostExporter
from cost.exporters.azure_cost_exporter import AzureCostExporter


class CostExporterFactory:
    """
    A factory class to instantiate and return cost exporter objects based on the specified cloud service provider (CSP).

    The CostExporterFactory provides a mechanism to create cost exporter instances for supported cloud
    service providers, such as AWS and Azure. By using the factory pattern, the main application can
    easily obtain the appropriate cost exporter object without directly dealing with specific class
    initializations.

    Usage:
        aws_cost_exporter = CostExporterFactory.create('aws')  # Returns an instance of AWSCostExporter
        azure_cost_exporter = CostExporterFactory.create('azure)  # Returns an instance of AzureCostExporter
    """

    @staticmethod
    def create(exporter_type):
        """
        Create and return a cost exporter object based on the provided exporter type.

        This static method instantiates and returns the appropriate cost exporter object
        based on the given cloud service provider type.
        If an unsupported exporter type is provided, a ValueError is raised.

        :param exporter_type: A string representing the type of the cloud service provider (e.g., 'aws', 'azure).
        :return: An instance of the appropriate cost exporter class (e.g., AWSCostExporter, AzureCostExporter).
        :raises ValueError: If the provided exporter_type is not supported.
        """
        if exporter_type == "aws":
            return AWSCostExporter()
        elif exporter_type == "azure":
            return AzureCostExporter()
        else:
            raise ValueError("Unsupported CSP type")
