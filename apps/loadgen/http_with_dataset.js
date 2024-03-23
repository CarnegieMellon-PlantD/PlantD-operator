import http from 'k6/http';
import { check } from 'k6';
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

let loadpattern = JSON.parse(open('loadpattern.json'));
let pipeline = JSON.parse(open('pipeline.json'));
let dataset = JSON.parse(open('dataset.json'));
const url = pipeline.http.url;
const method = pipeline.http.method;
const headers = pipeline.http.headers || {};
const numFiles = dataset.spec.numFiles;
const numSchemas = dataset.spec.schemas.length;
const compressedFileFormat = dataset.spec.compressedFileFormat || "";
const fileFormat = dataset.spec.fileFormat;
const compressPerSchema = dataset.spec.compressPerSchema || false;
const datasetName = dataset.metadata.name;
const fileExtention = {
  zip: 'zip',
  binary: 'bin'
};
let maxIndex;

const ext = fileExtention[fileFormat];

function filePerSchemaArray() {
  const n = numSchemas * numFiles;
  const arr = new Array(n);
  for (let i = 0; i < numSchemas; i++) {
    let k = i * numFiles;
    let schemaName = dataset.spec.schemas[i].name;
    for (let j = 0; j < numFiles; j++) {
      const fname = `${schemaName}/${datasetName}_${schemaName}_${j}.${ext}`;
      arr[k + j] = {
        name: fname,
        content: open(fname, 'b')
      };
    }
  }
  return arr;
}

function filepathPerCompressedArray() {
  const arr = new Array(numFiles);
  for (let i = 0; i < numFiles; i++) {
    const fname = `${datasetName}_${i}.${compressedFileFormat}`;
    arr[i] = {
      name: fname,
      content: open(fname, 'b')
    };
  }
  return arr;
}

function filepathPerCompressedPerSchemaArray() {
  const n = numSchemas * numFiles;
  const arr = new Array(n);
  for (let i = 0; i < numSchemas; i++) {
    let k = i * numFiles;
    let schemaName = dataset.spec.schemas[i].name;
    for (let j = 0; j < numFiles; j++) {
      const fname = `${schemaName}/${datasetName}_${schemaName}_${j}.${compressedFileFormat}`;
      arr[k + j] = {
        name: fname,
        content: open(fname, 'b')
      };
    }
  }
  return arr;
}

let dataCache;

if (compressedFileFormat === "") {
  maxIndex = numSchemas * numFiles - 1;
  dataCache = filePerSchemaArray()
} else if (compressPerSchema === true) {
  maxIndex = numSchemas * numFiles - 1;
  dataCache = filepathPerCompressedPerSchemaArray()
} else {
  maxIndex = numFiles - 1;
  dataCache = filepathPerCompressedArray()
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
  console.log(dataCache[i]['content'])
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
