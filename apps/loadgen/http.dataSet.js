import http from 'k6/http';
import { check } from 'k6';
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

let endpoint = JSON.parse(open('endpoint.json'));
let dataset = JSON.parse(open('dataset.json'));
let loadpattern = JSON.parse(open('loadpattern.json'));

const url = endpoint.http.url;
const method = endpoint.http.method;
const headers = endpoint.http.headers || {};

const dataSetName = dataset.metadata.name;
const numFiles = dataset.spec.numFiles;
const numSchemas = dataset.spec.schemas.length;
const compressedFileFormat = dataset.spec.compressedFileFormat || "";
const fileFormat = dataset.spec.fileFormat;
const compressPerSchema = dataset.spec.compressPerSchema || false;

const fileExtensions = {
  csv: 'csv',
  binary: 'bin'
};
const ext = fileExtensions[fileFormat];

function filePerSchemaArray() {
  const n = numSchemas * numFiles;
  const arr = new Array(n);
  for (let i = 0; i < numSchemas; i++) {
    let k = i * numFiles;
    const schemaName = dataset.spec.schemas[i].name;
    for (let j = 0; j < numFiles; j++) {
      const fname = `${schemaName}/${dataSetName}_${schemaName}_${j}.${ext}`;
      arr[k + j] = {
        name: fname,
        content: open(fname, 'b')
      };
    }
  }
  return arr;
}

function filePerCompressedArray() {
  const arr = new Array(numFiles);
  for (let i = 0; i < numFiles; i++) {
    const fname = `${dataSetName}_${i}.${compressedFileFormat}`;
    arr[i] = {
      name: fname,
      content: open(fname, 'b')
    };
  }
  return arr;
}

function filePerCompressedPerSchemaArray() {
  const n = numSchemas * numFiles;
  const arr = new Array(n);
  for (let i = 0; i < numSchemas; i++) {
    let k = i * numFiles;
    let schemaName = dataset.spec.schemas[i].name;
    for (let j = 0; j < numFiles; j++) {
      const fname = `${dataSetName}_${schemaName}_${j}.${compressedFileFormat}`;
      arr[k + j] = {
        name: fname,
        content: open(fname, 'b')
      };
    }
  }
  return arr;
}

let maxIndex;
let dataCache;
if (compressedFileFormat === "") {
  maxIndex = numSchemas * numFiles - 1;
  dataCache = filePerSchemaArray()
} else if (compressPerSchema === true) {
  maxIndex = numSchemas * numFiles - 1;
  dataCache = filePerCompressedPerSchemaArray()
} else {
  maxIndex = numFiles - 1;
  dataCache = filePerCompressedArray()
}

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
  const i = randomIntBetween(0, maxIndex)
  let payload = {
    file: http.file(dataCache[i]['content'], dataCache[i]['name'], 'multipart/form-data'),
  };
  let res = http.request(method, url, payload, {
    headers: headers,
  });
  check(res, {
    'status was 200': (r) => r.status === 200,
  });
}
