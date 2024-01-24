from plantd_modeling.configuration import Experiment
from plantd_modeling import metrics
from datetime import timedelta, datetime
import requests
import json
import pandas as pd
import redis
import os
import io

class PrometheusAgedOutException(Exception):
    pass

REALTIME_CALLS_QUERY = 'calls_total{{job="{params.name}", namespace="{params.namespace}"}}'
REALTIME_THROUGHPUT_QUERY = 'sum by(span_name)(irate(calls_total{{status_code="STATUS_CODE_UNSET", job="{params.name}", namespace="{params.namespace}"}}[{step}s]))'
REALTIME_LATENCY_QUERY = 'irate(duration_milliseconds_sum{{status_code="STATUS_CODE_UNSET", job="{params.name}", namespace="{params.namespace}"}}[30s]) / irate(duration_milliseconds_count{{status_code="STATUS_CODE_UNSET", job="{params.name}", namespace="{params.namespace}"}}[30s])'

class Redis:
    def __init__(self, redis_host, redis_password) -> None:
        self.redis_host = redis_host
        self.redis_password = redis_password
        if redis_host is None or redis_host == "":
            raise Exception("REDIS_HOST environment variable not set")
        self.r = redis.Redis(host=redis_host, port=6379, db=0, decode_responses=True, password=redis_password)

    def save_dict(self, type, name, data):
        self.r.set(f"plantd:{type}:{name}", json.dumps(data))

    def load_dict(self, type, name):
        return json.loads(self.r.get(f"plantd:{type}:{name}"))

    def save_str(self, type, name, data):
        self.r.set(f"plantd:{type}:{name}", data)

    def load_str(self, type, name):
        return self.r.get(f"plantd:{type}:{name}")

redis_password = os.environ.get("REDIS_PASSWORD", None)
redis_host = os.environ.get("REDIS_HOST", None)

redis = Redis(redis_host, redis_password)    

class Metrics:
    def __init__(self, prometheus_host) -> None:
        self.prometheus_host = prometheus_host
        

    # First, query from start time till now.  Figure out the right step interval to get < 11,000 samples
    # next, find where it goes quiet
    # then recalculate end time, and get new step interval
    # then requery from start time till that new end time

    def get_rough_end_time(self, experiment: Experiment):
        url = f"{self.prometheus_host}/api/v1/query_range"
        step_interval = 30
        start_ts = experiment.start_time.timestamp()
        end_ts = datetime.now().timestamp()
        step_interval = max(30, int((end_ts - start_ts) / 11000))
        #
        query = REALTIME_THROUGHPUT_QUERY.format(params=experiment.experiment_name, step=step_interval)
        print(query, datetime.utcfromtimestamp(start_ts), datetime.utcfromtimestamp(end_ts), step_interval)
        print(url)  
        print(query)
        response = requests.get(url, params={'query': query, 'start': start_ts, 'end': end_ts, 'step': step_interval}, 
            #auth=('prometheus', prometheus_password), 
            verify=False, stream=False)
        response.raise_for_status()

        
        #
        #import pdb; pdb.set_trace()

        if len(response.json()['data']['result']) == 0:
            raise PrometheusAgedOutException(f"No data found for {experiment.experiment_name}")
        
        dfs = []
        for result in response.json()['data']['result']:
            span = result['metric']['span_name']
            df = pd.DataFrame(result['values'], columns=['time', span])
            df['time'] = pd.to_datetime(df['time'], unit='s')
            df.set_index('time', inplace=True)
            
            for column, dtype in df.dtypes.iteritems():
                if dtype == 'object':
                    df[column] = pd.to_numeric(df[column], errors='coerce')
            dfs.append(df)

        if len(dfs) > 0:
            df = pd.concat(dfs, axis=1)
            last_non_zero_index = df.astype(bool).any(axis=1)[::-1].idxmax()
            return df[:last_non_zero_index]

        return None



    def get_metrics(self, experiment: Experiment, from_cached=False, also_latency=False):
        #import pdb; pdb.set_trace()
        if from_cached:
            return pd.read_csv(io.StringIO(metrics.redis.load_str("metrics", experiment.experiment_name)), index_col=0, parse_dates=True)
        
        try:
            roughtime = self.get_rough_end_time(experiment)
            if roughtime is None:
                return None
            
            end_ts = roughtime.index[-1].timestamp()
            url = f"{self.prometheus_host}/api/v1/query_range"

            start_ts = experiment.start_time.timestamp() #- step_interval*15
            step_interval = max(30, int((end_ts - start_ts) / 11000))

            query = REALTIME_THROUGHPUT_QUERY.format(params=experiment.experiment_name, step=step_interval)
            print(query, datetime.utcfromtimestamp(start_ts), datetime.utcfromtimestamp(end_ts), step_interval)
            print(url)  
            print(query)
       
            response = requests.get(url, params={'query': query, 'start': start_ts, 'end': end_ts, 'step': step_interval}, 
                #auth=('prometheus', prometheus_password), 
                verify=False, stream=False)
            response.raise_for_status()
            
            latency = {}

            if also_latency:
                response_latency = requests.get(url, params={'query': REALTIME_LATENCY_QUERY.format(params=experiment.experiment_name, step=step_interval), 'start': start_ts, 'end': end_ts, 'step': step_interval}, 
                    #auth=('prometheus', prometheus_password), 
                    verify=False, stream=False)
            
                response_latency.raise_for_status()

                latency = {r['metric']['span_name']: r for r in response_latency.json()['data']['result']}

            dfs = []
            for result in response.json()['data']['result']:
                span = result['metric']['span_name']
                df = pd.DataFrame(result['values'], columns=['time', span])
                df['time'] = pd.to_datetime(df['time'], unit='s')
                df.set_index('time', inplace=True)
                if also_latency and span in latency:
                    df_l = pd.DataFrame(latency[span]['values'], columns=['time', span])
                    df_l['time'] = pd.to_datetime(df_l['time'], unit='s')
                    df_l.set_index('time', inplace=True)
                    df[span + '_latency'] = df_l[span]
                dfs.append(df)

            if len(dfs) > 0:
                df = pd.concat(dfs, axis=1)
                for column, dtype in df.dtypes.iteritems():
                    if dtype == 'object':
                        df[column] = pd.to_numeric(df[column], errors='coerce')
                metrics.redis.save_str("metrics", experiment.experiment_name,df.to_json(orient='split', date_format='iso'))
            else:
                df = pd.DataFrame()
        except requests.exceptions.HTTPError as e:
            print(f"Error getting metrics for {experiment.experiment_name}: {e.response.text}")
            
            print(f"Trying redis")
            # read from redis
            df = pd.read_json(io.StringIO(metrics.redis.load_str("metrics", experiment.experiment_name)), orient='split', convert_dates=True)

            #df = pd.read_csv(io.StringIO(metrics.redis.load_str("metrics", experiment.experiment_name)), index_col=0, parse_dates=True)
        except Exception as e:
            print(f"Error getting metrics for {experiment.experiment_name}: {type(e)} {e}")
            print(f"Trying redis")
            # read from redis
            df = pd.read_json(io.StringIO(metrics.redis.load_str("metrics", experiment.experiment_name)), orient='split', convert_dates=True)

            #df = pd.read_csv(io.StringIO(metrics.redis.load_str("metrics", experiment.experiment_name)), index_col=0, parse_dates=True)

        return df