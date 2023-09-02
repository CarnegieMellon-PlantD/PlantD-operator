import io

import botocore.exceptions

from cost.exporters.base_cost_exporter import CostExporter
import os
import re
import datetime
import pandas as pd
from boto3 import client
import json

pd.set_option('display.max_columns', None)


class AWSCostExporter(CostExporter):
  def __init__(self, aws_access_key='', aws_secret_key='', *args, **kwargs):
    super(AWSCostExporter, self).__init__(*args, **kwargs)
    self.csp_credentials = os.environ.get('CSP_CREDENTIALS', '')
    # Parse the JSON
    data = json.loads(self.csp_credentials)

    # Extract AWS_ACCESS_KEY and AWS_SECRET_KEY
    self.aws_access_key = data.get("AWS_ACCESS_KEY")
    self.aws_secret_key = data.get("AWS_SECRET_KEY")
    experiment_tags = os.environ.get('EXPERIMENT_TAGS', [])
    self.cost_tags = json.loads(experiment_tags)
    self.s3_bucket_name = os.environ.get('S3_BUCKET_NAME', '')
    self.earliestTimestamp = os.environ.get('EARLIEST_EXPERIMENT', '')

  def _create_s3_connection(self):
    try:
      conn = client(
          's3',
          aws_access_key_id=self.aws_access_key,
          aws_secret_access_key=self.aws_secret_key
      )

      return conn
    except botocore.exceptions.PartialCredentialsError as e:
      print("Error: Check your AWS credentials. Ensure both AWS_ACCESS_KEY and AWS_SECRET_KEY are set.", e)
      quit()
    except botocore.exceptions.NoCredentialsError as e:
      print("Error: AWS credentials not found. Ensure both AWS_ACCESS_KEY and AWS_SECRET_KEY are set.", e)
      quit()
    except botocore.exceptions.EndpointConnectionError as e:
      print("Error: Unable to connect to the S3 endpoint.", e)
      quit()
    except Exception as e:
      print("Error: Unexpected error occurred while creating S3 connection.", e)
      quit()

  def _get_cost_files(self, conn):
    """
    Retrieve the last two cost files from AWS S3.

    This function retrieves a list of all cost logs from the specified S3 bucket, filters them
    based on a regular expression, sorts them by their key, and returns the last two.

    :param conn: The S3 client.
    :return: The last two cost log files.
    """
    try:
      # Calculate one hour before the earliestTimestamp
      self.earliestTimestamp = datetime.datetime.strptime(self.earliestTimestamp, '%Y-%m-%d %H:%M:%S')
      one_hour_before_earliest = self.earliestTimestamp - datetime.timedelta(hours=1)

      # Convert the timestamp to a string in the format 'YYYY-MM-DD HH:mm:ss'
      one_hour_before_str = one_hour_before_earliest.strftime('%Y-%m-%d %H:%M:%S')

      # Regex matching the name of a costlog on AWS
      cost_log_regex = re.compile(
        r"""costlog/cost-and-usage-report/\d{8}-\d{8}/cost-and-usage-report-\d+\.csv\.gz""")

      # Retrieve list of all cost logs, filter by regex and time
      cost_files = [
          cost_file for cost_file in conn.list_objects(
              Bucket=self.s3_bucket_name,
              Prefix='costlog/cost-and-usage-report/20'
          )['Contents'] if cost_log_regex.match(cost_file['Key'])
          and cost_file['LastModified'].strftime('%Y-%m-%d %H:%M:%S') >= one_hour_before_str
      ]

      return cost_files

    except botocore.exceptions.ClientError as e:
      print("Error: An error occurred while retrieving cost files from S3.", e)
      quit()
    except Exception as e:
      print("Error: Unexpected error occurred while retrieving cost files from S3.", e)
      quit()

  def _load_dataframe(self, conn, file_key):
    """
    Load a cost file into a pandas DataFrame.

    This function downloads a cost log from the specified S3 bucket using the given file key and
    loads it into a pandas DataFrame. If any error occurs during this process, an appropriate
    error message is printed and the program exits.

    :param conn: The S3 client.
    :param file_key: The file key of the cost log.
    :return: A pandas DataFrame of the cost log.
    """
    try:
      # Download cost log and load into pandas dataframe
      b1 = conn.get_object(Bucket=self.s3_bucket_name, Key=file_key)
      raw = b1["Body"].read()
      with open(file_key.split("/")[-1], "wb") as f:
        f.write(raw)

      df = pd.read_csv(io.BytesIO(raw), compression='gzip', low_memory=False)
      return df
    except botocore.exceptions.ClientError as e:
      print("Error: An error occurred while loading cost files from S3.", e)
      quit()
    except Exception as e:
      print("Error: Unexpected error occurred while loading cost files from S3.", e)
      quit()

  def put_logs_on_redis(self, cost_df, categories, exp_name):

    # Separate the tag columns from other categories
    tag_columns = categories[:-2]
    other_columns = [categories[-1]]

    # Iterate through the tag columns and values
    if all(tag in cost_df.columns for tag in tag_columns):
      for tag_column in tag_columns:
        for tag_value in cost_df[tag_column].dropna().unique():
          # Filter the DataFrame by the current tag column and value
          df_group = cost_df[cost_df[tag_column] == tag_value]
          # Group the filtered DataFrame by the other categories and iterate through the subgroups
          for gname, groupdf in df_group.groupby(other_columns, dropna=False):
            for hour, hgroupdf in groupdf.groupby(["starttime"]):
              values = {
                  "key": exp_name,
                  "tag": tag_value,
                  "resource": gname[-1],
                  "timestamp": int(hour[0].timestamp()) * 1000,
                  "cost": hgroupdf["lineItem/UnblendedCost"].sum().item(),
              }
              # print(values)
              self._write_to_db(values)

  def _filter_dataframe(self, df):
    """
    Filter a DataFrame of cost logs based on specified tags.

    This function filters a DataFrame of cost logs based on the tag keys and values set during
    initialization. The DataFrame is filtered such that it only contains rows where the tag key
    matches the tag value. Additional columns are also added for the start and end times, and
    any rows where the unblended cost is not greater than 0.0 are removed.

    :param df: Original cost logs
    :return df_filtered: df_filtered: Cost logs filtered by tags.
    :return categories: list of column names to be used when storing the logs.
    """
    # Initialize empty lists to store the extracted values
    # Initialize lists to store extracted values
    experiment_names = []
    tag_key_list_strings = []
    tag_value_list_strings = []

    for experiment_data in self.cost_tags:
      # Extract the experiment name
      experiment_name = experiment_data['Name']

      # Extract tag_key_list and tag_value_list
      tag_key_list = []
      tag_value_list = []

      for tag_item in experiment_data['Tags']:
        tag_key = tag_item['Key']
        tag_value = tag_item['Value']
        tag_key_list.append(tag_key)
        tag_value_list.append(tag_value)

      # Convert the tag_key_list and tag_value_list to strings
      tag_key_list_string = ', '.join(tag_key_list)
      tag_value_list_string = ', '.join(tag_value_list)

      # Append the extracted values to the respective lists
      experiment_names.append(experiment_name)
      tag_key_list_strings.append(tag_key_list_string)
      tag_value_list_strings.append(tag_value_list_string)

    # Process each experiment one by one
    for i in range(len(experiment_names)):
      key_list = tag_key_list_strings[i].split(",")
      value_list = tag_value_list_strings[i].split(",")

      condition = pd.Series([False] * len(df))

      for k, v in zip(key_list, value_list):
        column = f'resourceTags/user:{k}'
        if column not in df.columns:
          continue

        # Keep the column which is matched to the tag value
        condition |= (df[column] == v)

      df_filtered = df[condition]

      df_filtered.loc[:, "endtime"] = pd.to_datetime(df_filtered["lineItem/UsageEndDate"])
      df_filtered.loc[:, "starttime"] = pd.to_datetime(df_filtered["lineItem/UsageStartDate"])

      categories = ["resourceTags/user:" + key for key in key_list] + \
          ["lineItem/ProductCode", "lineItem/UsageType"]
      df_filtered = df_filtered[df_filtered["lineItem/UnblendedCost"] > 0.0]
      self.put_logs_on_redis(df_filtered, categories, experiment_names[i])
    return

  def get_cost_logs(self):
    s3_conn = self._create_s3_connection()

    cost_files = self._get_cost_files(s3_conn)
    for file in cost_files:
      df = self._load_dataframe(s3_conn, file["Key"])
      self._filter_dataframe(df)
    return
