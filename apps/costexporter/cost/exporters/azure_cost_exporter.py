from cost.exporters.base_cost_exporter import CostExporter
import os
from azure.identity import DefaultAzureCredential
from azure.mgmt.commerce import UsageManagementClient
from azure.core.exceptions import AzureError

import pandas as pd
import json
import datetime as dt
from datetime import datetime


class AzureCostExporter(CostExporter):
  def __init__(self, *args, **kwargs):
    super(AzureCostExporter, self).__init__(*args, **kwargs)
    self.azure_subscription_id = os.environ.get('AZURE_SUBSCRIPTION_ID', '')
    self.azure_tenant_id = os.environ.get('AZURE_TENANT_ID', '')
    self.tag_keys = os.environ.get('TAG_KEYS', '')
    self.tag_values = os.environ.get('TAG_VALUES', '')

  def _create_usage_client(self):
    """
    Create a client to interact with Azure's Usage Management services. 

    This method initializes an Azure Usage Management Client usng the DefaultAzureCredential.
    :return: An instance of UsageManagementClient if successful.
    """
    try:
      # Create client
      token_credential = DefaultAzureCredential()
      usage_client = UsageManagementClient(
          credential=token_credential,
          subscription_id=self.azure_subscription_id
      )

      return usage_client
    except Exception as e:
      print("Error: Check your Azure credentials.", e)
      quit()

  def _load_log_dict(self, item):
    """
    Given a single usage record from the Azure usage aggregates request,
    load selected fields into a Python dict and return the dict.
    :param item: nested json string
    :return: log_dict (Python dict)
    """
    log_dict = {}
    instance_data = json.loads(item.instance_data)["Microsoft.Resources"]
    log_dict["meter_category"] = item.meter_category
    log_dict["meter_name"] = item.meter_name
    log_dict["meter_id"] = item.meter_id
    log_dict["resourceUri"] = instance_data["resourceUri"]
    if "tags" in instance_data:
      log_dict["tags"] = instance_data["tags"]
    log_dict["usage_start_time"] = item.usage_start_time.strftime('%Y-%m-%d %H:%M:%S')
    log_dict["usage_end_time"] = item.usage_end_time.strftime('%Y-%m-%d %H:%M:%S')
    log_dict["quantity"] = item.quantity
    log_dict["unit"] = item.unit
    return log_dict

  def _get_usage_records(self, usage_report, tags):
    """
    Extract usage records from an Azure usage report based on resource tags. 

    This method iterates through each item in the provided Azure usage report,
    filters the items based on the provided tags, and returns the filtered records
    as a list of dictionaries.

    :param usage_report: Azure usage report
    :param tags: List of tuples representing key-value pairs or resource tags. 
    :return usage_records: A list of dictionaires representing filtered usage records. 
    """
    usage_records = []
    for item in usage_report:
      log_dict = self._load_log_dict(item)
      if "tags" not in log_dict:
        continue

      for tag in tags:
        if tag[0] in log_dict["tags"] and log_dict["tags"][tag[0]] == tag[1]:
          usage_records.append(log_dict)
          break

    return usage_records

  def _update_usage_records(self, usage_records, rate):
    """
    Update the provided usage records with rate and cost information.

    This method takes a list of usage records and a rate card from Azure. 
    It updates each record in the list with the rate for its meter and calculates its cost based
    on its quantity and rate. 

    :param usage_records: List of usage records.
    :param rate: Azure rate card.
    :return: List of usage records updated with rate and cost information.
    """
    test_meters = {}
    for record in usage_records:
      if record["meter_id"] not in test_meters:
        test_meters[record["meter_id"]] = 0
    for meter in rate.meters:
      if meter.meter_id in test_meters:
        test_meters[meter.meter_id] = meter.meter_rates["0"]
    for record in usage_records:
      record["rate"] = test_meters[record["meter_id"]]
      record["cost"] = record["rate"] * record["quantity"]

    return usage_records

  def _create_dataframe(self, usage_records):
    """
    Convert a list of usage records into a pandas DataFrame.

    This method takes a list of usage records, loads them into a pandas Dataframe,
    and groups them by meter information, usage start time, and tags. The resulting
    DataFrame aggregates costs for each group.

    :param usage_records: List of usage records.
    :return: A pandas DataFrame representing aggregated usage records. 
    """
    df = pd.DataFrame(usage_records)
    df["tags"] = df["tags"].astype(str)
    df["meter_info"] = df["meter_category"] + ":" + df["meter_name"]
    return df[["meter_info", "usage_start_time", "cost", "tags"]].groupby(
        ["meter_info", "usage_start_time", "tags"]).sum()

  def _convert_to_unix_timestamp(self, timestamp):
    """
    Convert an ISO formatted timestamp into a Unix timestamp in milliseconds.

    :param timestamp: ISO formatted timestamp.
    :return: Unix timestamp in milliseconds. 
    """
    datetime_obj = datetime.fromisoformat(timestamp)
    return int(datetime_obj.timestamp()) * 1000

  def get_cost_logs(self):
    usage_client = self._create_usage_client()

    # pre-fetch rate card
    rate = usage_client.rate_card.get(
        "OfferDurableId eq 'MS-AZR-0062P' and Currency eq 'USD' and Locale eq 'en-US' and RegionInfo eq 'US'"
    )

    tags = list(zip(self.tag_keys.split(","), self.tag_values.split(",")))

    start_time = str(dt.date.today() - dt.timedelta(days=14)) + 'T00:00:00Z'

    # set request endtime to nine hours in the past, rounded down to the nearest hour
    end_time = (dt.datetime.now(dt.timezone.utc) - dt.timedelta(hours=9)).strftime(
        '%Y-%m-%dT%H:%M:%SZ').split(":")[0] + ":00:00Z"

    usage_report = usage_client.usage_aggregates.list(
        start_time,
        end_time,
        aggregation_granularity='Hourly',
    )

    # collect usage records with at least one matching tag
    usage_records = self._get_usage_records(usage_report, tags)

    if len(usage_records) == 0:
      print("No usage records found\n")
      return

    usage_records = self._update_usage_records(usage_records, rate)

    df = self._create_dataframe(usage_records)

    # insert cost records into database
    for index, row in df.iterrows():
      values = {
          "key": index[2] + "/" + index[0],
          "resource": index[0],
          "timestamp": self._convert_to_unix_timestamp(index[1]),
          "tag": index[2],
          "cost": row["cost"]
      }
      self._write_to_db(values)

    return
