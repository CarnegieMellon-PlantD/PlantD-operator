from plantd_modeling import build, trafficmodel, twin
import sys
from plantd_modeling import configuration, metrics
import json
import os

try:
	pmodel = build.build_twin(os.environ['MODEL_TYPE'], from_cached=False)
	tmodel = trafficmodel.forecast(2025)
	twin.simulate(pmodel, tmodel)
except Exception as e:
	print("SIMULATION_STATUS: Failure")
	print(f"ERROR_REASON: {type(e)}: {e}")
	raise e
print("SIMULATION_STATUS: Success")
