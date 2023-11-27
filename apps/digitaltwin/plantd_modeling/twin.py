from dataclasses import dataclass  
import json
import os
import pandas as pd
from plantd_modeling import trafficmodel, metrics
import math

@dataclass
class PipelineModel:
    pass


# FUTURE WORK:
#    This just assumes the same cost per hour during the experiment as during real life.
#    but we should also gather cpu and ram utilization, and learn to infer those for future
#    traffic, then from that estimate cost.
    
@dataclass
class SimpleModel(PipelineModel):
    maxrate_rph: float    # max records per hour this pipe can process
    per_vm_hourcost: float   
    avg_latency_s: float  # average latency in seconds of whole pipeline, assuming no queueing
    policy: str
    numproc: int = 1

    def __post_init__(self):
        if not self.policy in ["fifo","lifo","random"]:
            raise Exception(f"Unknown scaling policy {self.policy}")
        
    # Serialize and deserialize as json
    def serialize(self):
        return json.dumps({
            "model_type": "simple",
            "maxrate_rph": self.maxrate_rph,
            "per_vm_hourcost": self.per_vm_hourcost,  #misnomer; this is per "nomimal pipeline", not per "vm".
            "avg_latency_s": self.avg_latency_s,
            "policy": self.policy
        })
    
    @classmethod
    def deserialize(cls, jsonstr):
        params = json.loads(jsonstr)
        if params["model_type"] != "simple":
            raise Exception(f"Unknown model type {params['model_type']}")
        return SimpleModel(params["maxrate_rph"], params["per_vm_hourcost"], params["avg_latency_s"], params["policy"])
        
    def reset(self):
        #self.upq = []
        #self.dnq = []
        #self.numproc = 1
        self.cumu_cost = 0.0
        self.queue = 0
        self.queue_worstcase_age_s = 0
        self.throughput_rph = 0
        self.latency_fifo_s = 0
        self.latency_lifo_s = 0
        self.hourcost = 0
        
    def input(self, recs_this_hour):  
        self.throughput_rph = min(recs_this_hour + self.queue, self.maxrate_rph)
        self.latency_fifo_s = self.avg_latency_s + self.queue * 1.0 / self.maxrate_rph
        self.queue = self.queue + recs_this_hour - self.throughput_rph
        self.queue_worstcase_age_s += 3600
        if self.queue < 0: 
            self.queue = 0
        if self.queue == 0:
            self.queue_worstcase_age_s = 0
        self.latency_lifo = self.avg_latency_s + self.queue_worstcase_age_s
        self.hourcost = self.per_vm_hourcost
        self.cumu_cost = self.cumu_cost + self.hourcost
        
        
"""        
===> TO DO: 
    v copy the code that deals with queueing stuff to here.  That should be in the pipe model, not the traffic model!
    - simplemodel should have a second-by-second simulation mode for testing, and a pandas mode for the real bulk work
    - validate pandas mode against s-by-s
    - output a characterization and trace.
    -v maybe the traffic model shoudl take this as an object? Or vice versa?
   """

@dataclass    
class QuickscalingModel(PipelineModel):
    fixed_hourcost: float
    basemodel: SimpleModel
    policy: str
        
    def __post_init__(self):
        if not self.policy in ["fifo","lifo","random"]:
            raise Exception(f"Unknown scaling policy {self.policy}")
        

    # Serialize and deserialize as json
    def serialize(self):
        return json.dumps({
            "model_type": "quickscaling",
            "fixed_hourcost": self.fixed_hourcost,
            "basemodel": self.basemodel.serialize(),
            "policy": self.policy
        })
    
    @classmethod
    def deserialize(cls, jsonstr):
        params = json.loads(jsonstr)
        if params["model_type"] != "quickscaling":
            raise Exception(f"Unknown model type {params['model_type']}")
        return AutoscalingModel(params["fixed_hourcost"], SimpleModel.deserialize(params["basemodel"]), params["policy"])

    def reset(self):
        self.numproc = 1
        self.cumu_cost = 0.0
        self.queue = 0
        self.queue_worstcase_age_s = 0
        self.throughput_rph = 0
        self.latency_fifo_s = 0
        self.latency_lifo_s = 0
        self.hourcost = 0
        
    def input(self, recs_this_hour):  
        self.numproc = math.ceil(recs_this_hour / self.basemodel.maxrate_rph)
        self.throughput_rph = min(recs_this_hour + self.queue, self.basemodel.maxrate_rph * self.numproc)
        self.latency_fifo_s = self.basemodel.avg_latency_s + self.queue * 1.0 / (self.basemodel.maxrate_rph * self.numproc)
        self.queue = self.queue + recs_this_hour - self.throughput_rph
        self.queue_worstcase_age_s += 3600
        if self.queue < 0: 
            self.queue = 0
        if self.queue == 0:
            self.queue_worstcase_age_s = 0
        self.latency_lifo_s = self.basemodel.avg_latency_s + self.queue_worstcase_age_s
        self.hourcost = self.fixed_hourcost + self.basemodel.per_vm_hourcost * self.numproc
        self.cumu_cost = self.cumu_cost + self.hourcost  

@dataclass    
class AutoscalingModel(PipelineModel):
    fixed_hourcost: float
    upPctTrigger: float
    upDelay_h: int            # wait this long before scaling up
    dnPctTrigger: float
    dnDelay_h: float          # wait this long before scaling down
    basemodel: SimpleModel
    policy: str
        
    def __post_init__(self):
        if not self.policy in ["fifo","lifo","random"]:
            raise Exception(f"Unknown scaling policy {self.policy}")
        

    # Serialize and deserialize as json
    def serialize(self):
        return json.dumps({
            "model_type": "autoscaling",
            "fixed_hourcost": self.fixed_hourcost,
            "upPctTrigger": self.upPctTrigger,
            "upDelay_h": self.upDelay_h,
            "dnPctTrigger": self.dnPctTrigger,
            "dnDelay_h": self.dnDelay_h,
            "basemodel": self.basemodel.serialize(),
            "policy": self.policy
        })
    
    @classmethod
    def deserialize(cls, jsonstr):
        params = json.loads(jsonstr)
        if params["model_type"] != "autoscaling":
            raise Exception(f"Unknown model type {params['model_type']}")
        return AutoscalingModel(params["fixed_hourcost"], params["upPctTrigger"], params["upDelay_h"], params["dnPctTrigger"], params["dnDelay_h"], SimpleModel.deserialize(params["basemodel"]), params["policy"])

    def reset(self):
        self.upq_rph = []      # recs per hour processed for last upDelay_h hours
        self.dnq_rph = []      # recs per hour processed for last dnDelay_h hours
        self.numproc = 1
        self.cumu_cost = 0.0
        self.queue = 0
        self.queue_worstcase_age_s = 0
        self.throughput_rph = 0
        self.latency_fifo_s = 0
        self.latency_lifo_s = 0
        self.hourcost = 0
        self.time_since_scale_h = 0   # hours since scaling last changed (scale up or down)
        
    def input(self, recs_this_hour):  
        self.throughput_rph = min(recs_this_hour + self.queue, self.basemodel.maxrate_rph * self.numproc)
        self.latency_fifo_s = self.basemodel.avg_latency_s + self.queue * 1.0 / (self.basemodel.maxrate_rph * self.numproc)
        self.queue = self.queue + recs_this_hour - self.throughput_rph
        self.queue_worstcase_age_s += 3600
        if self.queue < 0: 
            self.queue = 0
        if self.queue == 0:
            self.queue_worstcase_age_s = 0
        self.latency_lifo_s = self.basemodel.avg_latency_s + self.queue_worstcase_age_s
        self.hourcost = self.fixed_hourcost + self.basemodel.per_vm_hourcost * self.numproc
        self.cumu_cost = self.cumu_cost + self.hourcost  
        self.scale()
        
    def scale(self):
        def avg(x): return sum(x)/len(x)
        self.time_since_scale_h += 1
        self.upq_rph = (self.upq_rph + [self.throughput_rph])[-self.upDelay_h:]
        self.dnq_rph = (self.dnq_rph + [self.throughput_rph])[-self.dnDelay_h:]
        if self.time_since_scale_h >= self.upDelay_h \
          and avg(self.upq_rph) > self.basemodel.maxrate_rph * self.numproc * self.upPctTrigger/100.0:
            #print(f"Scale up from {self.numproc}: {avg(self.upq_rph)} > {self.basemodel.maxrate_rph} * {self.numproc} * {self.upPctTrigger/100.0}")
            #import pdb; pdb.set_trace()
            self.numproc += 1
            self.time_since_scale_h = 0
        elif self.time_since_scale_h >= self.dnDelay_h \
             and avg(self.dnq_rph) < self.basemodel.maxrate_rph * self.numproc * self.dnPctTrigger/100.0 \
             and self.numproc > 1:
            #print(f"Scale down from {self.numproc}: {avg(self.dnq_rph)} < {self.basemodel.maxrate_rph} * {self.numproc} * {self.dnPctTrigger/100.0}")
            self.numproc -= 1
            self.time_since_scale_h = 0
            

@dataclass    
class AutoscalingModelFine(PipelineModel):
    fixed_hourcost: float
    upPctTrigger: float
    upDelay_s: int
    dnPctTrigger: float
    dnDelay_s: float
    basemodel: SimpleModel
    policy: str
        
    def __post_init__(self):
        if not self.policy in ["fifo","lifo","random"]:
            raise Exception(f"Unknown scaling policy {self.policy}")
        
    def reset(self):
        self.upq_s = []
        self.dnq_s = []
        self.numproc = 1
        self.cumu_cost = 0.0
        self.queue = 0
        self.queue_worstcase_age_s = 0
        self.throughput_rph = 0
        self.throughput_s = 0
        self.latency_fifo_s = 0
        self.latency_lifo_s = 0
        self.hourcost = 0
        self.time_since_scale_s = 0
        self.maxrate_s = self.basemodel.maxrate_rph/3600.0
        
    def input(self, recs_this_hour):
        
        slots = [0]*3600
        for i in range(int(recs_this_hour)):
            slots[int(i/3600.0)] += 1
        self.hourcost = 0
        self.throughput_rph = 0
        for recs_this_sec in slots:
            self.throughput_s = min(recs_this_sec + self.queue, self.maxrate_s * self.numproc)
            self.throughput_rph += self.throughput_s
            self.latency_fifo_s = self.basemodel.avg_latency_s + self.queue * 1.0 / (self.maxrate_s * self.numproc)
            self.queue = self.queue + recs_this_sec - self.throughput_s
            self.queue_worstcase_age_s += 1
            if self.queue < 0: 
                self.queue = 0
            if self.queue == 0:
                self.queue_worstcase_age_s = 0
            self.latency_lifo_s = self.basemodel.avg_latency_s + self.queue_worstcase_age_s
            self.hourcost += (self.fixed_hourcost/3600.0 + self.basemodel.per_vm_hourcost/3600.0 * self.numproc)
            self.scale()
        self.cumu_cost += self.hourcost
        
        
    def scale(self):
        def avg(x): return sum(x)/len(x)
        self.time_since_scale_s += 1
        self.upq_s = (self.upq_s + [self.throughput_s])[-self.upDelay_s:]
        self.dnq_s = (self.dnq_s + [self.throughput_s])[-self.dnDelay_s:]
        if self.time_since_scale_s >= self.upDelay_s and avg(self.upq_s) > self.maxrate_s * self.numproc * self.upPctTrigger/100.0:
            #print(f"Scale up from {self.numproc}: {avg(self.upq_s)} > {self.maxrate_s} * {self.numproc} * {self.upPctTrigger/100.0}")
            self.numproc += 1
            self.time_since_scale_s = 0
        elif self.time_since_scale_s >= self.dnDelay_s \
             and avg(self.dnq_s) < self.maxrate_s * self.numproc * self.dnPctTrigger/100.0 \
             and self.numproc > 1:
            #print(f"Scale down from {self.numproc}: {avg(self.dnq_s)} < {self.maxrate_s} * {self.numproc} * {self.dnPctTrigger/100.0}")
            self.numproc -= 1
            self.time_since_scale_s = 0



def simulate(twin, traffic):
    twin_name = os.environ['TWIN_NAME']
    sim_name = os.environ['SIM_NAME']
    traffic_model_name = os.environ['TRAFFIC_MODEL_NAME']
    
    #twin = SimpleModel.deserialize(open(f"fakeredis/twinmodel_{twin_name}.json").read())
    #traffic = trafficmodel.TrafficModel.deserialize_parameters(open(f"fakeredis/trafficmodel_{traffic_model_name}.json").read())
    #traffic.deserialize_forecast(f"fakeredis/trafficmodel_{traffic_model_name}.csv")
    thru = traffic.calculate_throughput(twin)
    queue = traffic.calculate_queue(twin)

    #traffic.traffic.to_csv(f"fakeredis/simulation_{sim_name}.csv")
    metrics.redis.save_str("simulation_traffic", sim_name, traffic.traffic.to_csv(index=True))
    
    sla_check = traffic.sla_check({"latency_sla_percent": 99.0, "latency_sla_limit": 70.0})


    metrics.redis.save_str("simulation_summary", sim_name, 
        json.dumps({"total_cost": float(traffic.traffic.cost.sum()), "avg_latency_s": float(traffic.traffic.latency_fifo.mean()),
               "max_latency_s": float(traffic.traffic.latency_fifo.max()), "avg_queue": float(traffic.traffic.queue_len.mean()),
               "max_queue": float(traffic.traffic.queue_len.max()), "avg_throughput_rph": float(traffic.traffic.throughput.mean()),
                "max_throughput_rph": float(traffic.traffic.throughput.max()),
                "sla_check": sla_check}))
    print(f"Wrote simulation results to redis/simulation_{sim_name}.csv and redis/simulation_{sim_name}.json")
    return twin




# QUick refactoring:
#   Traffic model should go in a redis file, not in the environment. So, for now, manually put it in a file, and refer by name in the code
#   Traffic model will have optional generated runs as separate files.  The main model file will have key-value dictionary of these, identifying date range -> filename
#   Here in simulate, I'll use the traffic model name to load the traffic model object, which in turn loads all of its runs.  The first one is the default