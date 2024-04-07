import http from 'k6/http';
import { check } from 'k6';

let endpoint = JSON.parse(open('endpoint.json'));
let plainText = open('plaintext.txt');
let loadPattern = JSON.parse(open('loadpattern.json'));

const url = endpoint.http.url;
const method = endpoint.http.method;
const headers = endpoint.http.headers || {};
const data = plainText;

export let options = {
  scenarios: {
    sendPlainText: {
      executor: 'ramping-arrival-rate',
      startRate: loadPattern.spec.startRate,
      timeUnit: loadPattern.spec.timeUnit,
      preAllocatedVUs: loadPattern.spec.preAllocatedVUs,
      maxVUs: loadPattern.spec.maxVUs,
      stages: loadPattern.spec.stages,
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