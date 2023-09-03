import http from 'k6/http';
import { check } from 'k6';

let loadpattern = JSON.parse(open('loadpattern.json'));
let pipeline = JSON.parse(open('pipeline.json'));
const url = pipeline.http.url;
const method = pipeline.http.method;
const data = JSON.stringify(pipeline.http.body.data || "");
const headers = pipeline.http.headers || {};
export let options = {
  scenarios: {
    ramping_arrival_rate: {
      executor: 'ramping-arrival-rate',
      startRate: loadpattern.startRate,
      timeUnit: loadpattern.timeUnit,
      preAllocatedVUs: loadpattern.preAllocatedVUs,
      maxVUs: loadpattern.maxVUs,
      stages: loadpattern.stages,
    },
  },
  discardResponseBodies: true,
  noVUConnectionReuse: true,
};

export default function () {
  let res = http.request(method, url, data, {
    headers: headers,
  });
  check(res, {
    'status was 200': (r) => r.status === 200,
  });
}