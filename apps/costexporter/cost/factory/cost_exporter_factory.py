from cost.exporters.aws_cost_exporter import AWSCostExporter
from cost.exporters.azure_cost_exporter import AzureCostExporter


class CostServiceFactory:
    """
    A factory class to instantiate and return cost service objects based on the specified cloud service provider (CSP).

    The CostServiceFactory provides a mechanism to create cost service instances for supported cloud
    service providers, such as AWS and Azure. By using the factory pattern, the main application can
    easily obtain the appropriate cost service object without directly dealing with specific class
    initializations.

    Usage:
        aws_cost_service = CostServiceFactory.create('aws')  # Returns an instance of AWSCostService
        azure_cost_service = CostServiceFactory.create('azure)  # Returns an instance of AzureCostService
    """

    @staticmethod
    def create(service_type):
        """
        Create and return a cost service object based on the provided service type.

        This static method instantiates and returns the appropriate cost service object
        based on the given cloud service provider type.
        If an unsupported service type is provided, a ValueError is raised.

        :param service_type: A string representing the type of the cloud service provider (e.g., 'aws', 'azure).
        :return: An instance of the appropriate cost service class (e.g., AWSCostService, AzureCostService).
        :raises ValueError: If the provided service_type is not supported.
        """
        if service_type == "aws":
            return AWSCostService()
        elif service_type == "azure":
            return AzureCostService()
        else:
            raise ValueError("Unsupported CSP type")
