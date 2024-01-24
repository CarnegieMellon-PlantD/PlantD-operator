"""
k8s-extperiment-cost.py
Calculates cost of one pipeline experiment and writes results to redis.

Utilizes OpenCost and Prometheus APIs to fetch required records and 
performs necessary calculations.
"""
import os
import datetime as dt
import json
import time
import requests
from plantd_modeling import metrics

def get_cost(source, experiment_name, pipeline_namespace, start_time, end_time, from_cached=False):
    if source != "opencost":
        return None
    
    if from_cached:
        experiment_cost = metrics.redis.load_dict("experiment_cost", experiment_name)
        return experiment_cost
    
    # Get endpoints from environment variables
    # load_dotenv(".env")   # read from .env file for local testing only
    prometheus_endpoint = os.environ.get("PROMETHEUS_HOST", "http://localhost:9090")
    opencost_endpoint = os.environ.get("OPENCOST_ENDPOINT", "http://localhost:9003")
    
    # Get experiment tag and start and end times from environment variables
    pipeline_label_key = os.environ.get("PIPELINE_LABEL_KEY", "app")
    pipepline_label_value = os.environ.get("PIPELINE_LABEL_VALUE", "unzipper")
   
    # Get experiment tag and start and end times from environment variables
    pipeline_label_key = os.environ.get("PIPELINE_LABEL_KEY", None)
    pipepline_label_value = os.environ.get("PIPELINE_LABEL_VALUE", None)
    pipeline_namespace = os.environ.get("PIPELINE_NAMESPACE", "ubi")
    #start_time = os.environ.get("START_TIME", "2023-11-16T22:50:00Z")
    #end_time = os.environ.get("END_TIME", "2023-11-17T00:40:00Z")

    
    print("Start time: ", start_time)
    print("OpenCost endpoint: ", opencost_endpoint)

    # convert times to datetime objects with UTC timezone
    #start_time = dt.datetime.strptime(start_time, '%Y-%m-%dT%H:%M:%SZ')
    #start_time = start_time.replace(tzinfo=dt.timezone.utc)
    #end_time = dt.datetime.strptime(end_time, '%Y-%m-%dT%H:%M:%SZ')
    #end_time = end_time.replace(tzinfo=dt.timezone.utc)

    # get cost data from OpenCost API
    cost_data = get_cost_data(opencost_endpoint, pipeline_label_key, 
                                pipepline_label_value, pipeline_namespace, 
                                start_time, end_time)


    # get additional data from prometheus
    prometheus_data = get_prometheus_data(prometheus_endpoint)

    # calculate cost of experiment
    experiment_cost = calculate_experiment_cost(cost_data, prometheus_data,
                        (end_time - start_time).total_seconds(),
                        opencost_endpoint)

    # write experiment cost to redis
    write_experiment_cost(experiment_name, experiment_cost)

    return(experiment_cost)

def get_cost_data(opencost_endpoint, pipeline_label_key, pipeline_label_value, 
        pipeline_namespace, start_time, end_time):
    """
    Fetches cost data from OpenCost API.

    Input:
        opencost_endpoint (str): URL of OpenCost API
        experiment_label (str)
        start_time (UTC datetime)
        end_time (UTC datetime)
    Returns:
        cost_data (dict): usage data from OpenCost API
    """
    # store parameters for API request
    params = {} 
    # convert start_time to string
    params["window"] = start_time.strftime('%Y-%m-%dT%H:%M:%SZ') + "," + \
        end_time.strftime('%Y-%m-%dT%H:%M:%SZ')
    # params["aggregate"] = "namespace"
    if pipeline_label_key is not None and len(pipeline_label_key) > 0 :
        params["aggregate"] = "namespace,label:" + pipeline_label_key + ",pod"
    else:
        params["aggregate"] = "namespace,pod"
    params["accumulate"] = "false"
    params["resolution"] = "1m"
    
    print("Calling: ", opencost_endpoint, " with params: ", params)

    # make API request
    response = requests.get(opencost_endpoint + "/allocation", params=params)
    if response.status_code >= 500:
        print("Error querying OpenCost API: ", response.status_code)
        print("Exiting...")
        # exit(1)
    elif response.status_code >= 400 and response.status.code < 500:
        print("Error querying OpenCost API: ", response.status_code)
        print("Ignoring...")

    # load required records into dict and return 
    try: 
        if pipeline_label_key is not None and len(pipeline_label_key) > 0:
            response_key = pipeline_namespace + "/" + pipeline_label_value
        else:
            response_key = pipeline_namespace
        print("Response key: ", response_key)        

        pod_records = response.json()['data'][0]
        #print(pod_records)
        # iterate through pod_records, find all that have key starting with 
        # response_key, and store selected records in opencost_records

        opencost_records = {}
        for key in pod_records.keys():
            if key.startswith(response_key):
                # create a dict of relevant records for this pod
                pod_dict = {}
                pod_dict["cpuCore"] = max(float(pod_records[key]["cpuCoreRequestAverage"]),
                                        float(pod_records[key]["cpuCoreUsageAverage"]))
                pod_dict["ramByteUsageAverage"] = float(pod_records[key]["ramByteUsageAverage"])
                pod_dict["loadBalancerCost"] = float(pod_records[key]["loadBalancerCost"])
                pod_dict["pvCost"] = float(pod_records[key]["pvCost"])
                #print(key)
                #print(pod_dict)
                #print("------")
                # add pod_dict to opencost_records
                opencost_records[key] = pod_dict
    except:
        print("Error parsing OpenCost response: ", response.json())
        print("Exiting...")
        # exit(1)

    return opencost_records

def get_prometheus_data(prometheus_endpoint):
    """
    Fetches additional data from Prometheus API.

    Input:
        prometheus_endpoint (str): URL of Prometheus API
    Return:
        prometheus_data (dict): resource costs stored in Prometheus
    """
    prometheus_records = {}
    params = {}
    for metric in ["kubecost_cluster_management_cost", "node_cpu_hourly_cost", 
                    "node_ram_hourly_cost"]:
        params["query"] = metric
        response = requests.get(prometheus_endpoint + "/api/v1/query", params=params)
        if response.status_code != 200:
            print("Error querying Prometheus API: ", response.status_code)
            print("Exiting...")
            # exit(1)
        try:
            prometheus_records[metric] = float(response.json()["data"]["result"][0]["value"][1])
        except:
            print("Error parsing Prometheus response: ", response.json())
            print("Exiting...")
            # exit(1)
    return prometheus_records

def calculate_experiment_cost(cost_data, prometheus_data, duration, opencost_endpoint):
    """
    Calculates cost of experiment from cost and prometheus data.

    Input:
        cost_data (dict): usage data from OpenCost API; {"pod_name": {"cpuCore": float,
                        "ramByteUsageAverage": float, "loadBalancerCost": float, pvCost: float}
        prometheus_data (dict): resource costs stored in Prometheus
        duration (int): duration of experiment in seconds
    Returns:
        experiment_cost (float): cost of experiment in USD
    """

    """
    Apply the following method to each pod in cost_data:

    # calculate cost of CPU and RAM
    cpu_cost = cost_data["cpuCore"] * prometheus_data["node_cpu_hourly_cost"] * duration / 3600
    ram_cost = cost_data["ramByteUsageAverage"] * 2**-30 * prometheus_data["node_ram_hourly_cost"] \
                * duration / 3600
    
    direct_cost = cpu_cost + ram_cost + cost_data["loadBalancerCost"] + cost_data["pvCost"]

    # distribute shared costs evenly across namespaces
    shared_cost = (prometheus_data["kubecost_cluster_management_cost"] * duration) / \
                    (count_namespaces(opencost_endpoint) * 3600) 
    """
    experiment_cost = {}
    num_namespaces = count_namespaces(opencost_endpoint)

    for pod_key in cost_data:
        # calculate cost of CPU and RAM
        cpu_cost = cost_data[pod_key]["cpuCore"] * prometheus_data["node_cpu_hourly_cost"] * duration / 3600
        ram_cost = cost_data[pod_key]["ramByteUsageAverage"] * 2**-30 * prometheus_data["node_ram_hourly_cost"] \
                    * duration / 3600
        
        pod_cost = {}
        pod_cost["direct_cost"] = cpu_cost + ram_cost + cost_data[pod_key]["loadBalancerCost"] + cost_data[pod_key]["pvCost"]
        # split overhead costs evenly among namespaces, then among pods in namespaces
        pod_cost["shared_cost"] = (prometheus_data["kubecost_cluster_management_cost"] * duration) / \
                        (num_namespaces * 3600 * len(cost_data))
        
        pod_cost["total_cost"] = pod_cost["direct_cost"] + pod_cost["shared_cost"]
        experiment_cost[pod_key] = pod_cost
    
    print("Experiment costs: ", experiment_cost)
    return experiment_cost


def count_namespaces(opencost_endpoint):
    """Determine the number of namespaces in cluster"""
    
    # get list of namespaces
    # store parameters for API request
    params = {} 
    # convert start_time to string
    params["window"] = '1d'
    params["aggregate"] = "namespace"
    params["accumulate"] = "false"
    params["resolution"] = "1m"
    
    print("Namespace Params: ", params)

    # make API request
    response = requests.get(opencost_endpoint, params=params)
    if response.status_code != 200:
        print("Error querying OpenCost API: ", response.status_code)
        print("Exiting...")
        # exit(1)

    try: 
        num_namespaces = len(response.json()['data'][0])
    except:
        print("Error parsing OpenCost response: ", response.json())
        print("Exiting...")
        # exit(1)
    print("There are ", num_namespaces, " namespaces in the cluster")
    return num_namespaces

def write_experiment_cost(experiment_name, experiment_cost):
    metrics.redis.save_dict("experiment_cost", experiment_name, experiment_cost)
 
    #with open(f"fakeredis/experiment_cost_{experiment_name}.json","w") as f:
    #    json.dump(experiment_cost, f)
    """
    Writes experiment cost to redis.

    Input:
        experiment_cost (dict): {"pod_name": {"direct_cost": float, "shared_cost": float, 
                            total_cost: float}}
    """
    


if __name__ == '__main__':
    main()