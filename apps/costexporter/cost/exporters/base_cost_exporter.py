from abc import ABC, abstractmethod
import os
import redis


class CostExporter(ABC):
  """
  An abstract base class representing a generic cost service. 

  The CostExporter class provides a foundational framework for integrating cost retrieval and management
  from various cloud service providers. Initially, this framework is designed to support AWS and Azure, 
  but it's architected in a way to easily accommodate additional cloud providers in the future. 

  To integrate a new cloud service provider, one should extend this abstract class and implement
  the required methods specific to that provider. 
  """

  def __init__(self):
    self.redis_host = os.environ.get('REDIS_HOST', 'localhost')
    self.redis_port = os.environ.get('REDIS_PORT', 6379)
    self.db = self._create_redis_time_series()
    self.cloud_service_provider = os.environ.get('CLOUD_SERVICE_PROVIDER', '')

  def _create_redis_time_series(self):
    """
    Establish a connection to the Redis time series database. 

    This method tries to create a connection to the Redis time series using the provided host and port. 
    :return: The Redis time series object if the connection is successfully established. 
    """
    try:
      redis_pool = redis.ConnectionPool(host=self.redis_host, port=self.redis_port, db=0)
      redis_conn = redis.Redis(connection_pool=redis_pool)
      return redis_conn
    except redis.ConnectionError as e:
      print("Error: Failed to connect to Redis.", e)
      quit()
    except Exception as e:
      print("Error: Unexpected error occurred while creating Redis time series.", e)
      quit()

  def _write_to_db(self, values):
    """
    Write cost log entries to the Redis time series database. 

    This method takes a dictionary of values and writes a time series record to Redis. 
    The values dictionary should provide key, timestamp, cost, tag, and resource for the entry. 
    If any error occrus during the write operation, an appropriate error message is printed. 

    :param values: Dictionary containing necessary data (key, timestamp, cost, tag, resource)
    :return: None
    """
    try:
      print("Writing to Redis")
      print(values)
      self.db.hmset(
          values.get("key"),
          {
              "timestamp": values.get("timestamp"),
              "cost": values.get("cost"),
              "tag": values.get("tag"),
              "resource": values.get("resource"),
          }
      )
    except redis.RedisError as e:
      print("Error: Failed to write to Redis.", e)
      raise e
    except Exception as e:
      print("Error: Unexpected error occurred while writing records to Redis.", e)
      raise e

  @abstractmethod
  def get_cost_logs(self):
    """
    Abstract method to retrieve cost logs. 

    This method should be implemented in each child class of CostExporter. 
    The main purpose of this method is to fetch and process cost logs specific to a cloud service provider. 
    The main application will call this method to obtain the cost logs. 
    """
    pass
