import http from 'k6/http';
import { check } from 'k6';
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

const endpoint = JSON.parse(open('endpoint.json'));
const dataSet = JSON.parse(open('dataset.json'));
const loadPattern = JSON.parse(open('loadpattern.json'));

const url = endpoint.http.url;
const method = endpoint.http.method;
const headers = endpoint.http.headers || {};

const dataSetName = dataSet.metadata.name;
const numFiles = dataSet.spec.numFiles;
const numSchemas = dataSet.spec.schemas.length;
const compressedFileFormat = dataSet.spec.compressedFileFormat || "";
const fileFormat = dataSet.spec.fileFormat;
const compressPerSchema = dataSet.spec.compressPerSchema || false;

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
    const schemaName = dataSet.spec.schemas[i].name;
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

function filePerCompressedPerSchemaArray() {
  const n = numSchemas * numFiles;
  const arr = new Array(n);
  for (let i = 0; i < numSchemas; i++) {
    let k = i * numFiles;
    let schemaName = dataSet.spec.schemas[i].name;
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

export const options = {
  scenarios: {
    sendDataSet: {
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
