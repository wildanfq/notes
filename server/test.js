import http from "k6/http";
import { check, sleep } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

export const options = {
  stages: [
    { duration: "1m", target: 10 },
    { duration: "2m", target: 30 },
    { duration: "1m", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<1500"],
    http_req_failed: ["rate<0.05"],
  },
};

const BASE_URL = "http://34.29.132.244:8080";

export default function () {
  const url = `${BASE_URL}/notes`;

  const payload = JSON.stringify({
    title: `Note ${uuidv4()}`,
    content: "Testing performance",
  });

  const params = {
    headers: { "Content-Type": "application/json" },
  };

  const postResponse = http.post(url, payload, params);
  check(postResponse, {
    "post status 201": (r) => r.status === 201,
  });

  const getResponse = http.get(url);
  check(getResponse, {
    "get status 200": (r) => r.status === 200,
  });

  sleep(1);
}
