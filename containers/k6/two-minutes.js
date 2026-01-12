import http from "k6/http";
import { check, sleep } from "k6";

// https://k6.io/blog/how-to-generate-a-constant-request-rate-with-the-new-scenarios-api/
export const options = {
  thresholds: {
    // http errors should be less than 1%
    http_req_failed: ['rate<0.01'], 
    // 95% of requests should be below 2ms
    http_req_duration: ['p(95)<2'], 
  },
  scenarios: {
    typical_usage: {
      executor: 'ramping-vus',
      startVUs: 1,
      gracefulRampDown: "5s",
      stages: [
        { duration: "30s", target: 10 },
        { duration: "60s", target: 30 },
        { duration: "30s", target: 10 },
      ],
    },
  },
};

function getRandE(list) {
  const randomIndex = Math.floor(Math.random() * list.length);
  return list[randomIndex];
}

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}

// Simulate 30 possible sources
var possible_sources = []
for (const source_number of Array(30).keys()) {
  possible_sources.push("app.source." + source_number) 
}

// Simulate 100 possible different events
var possible_events = []
for (const event_number of Array(100).keys()) {
  possible_events.push("event_" + event_number)
}
const params = {
            headers: {"x-api-key": "not-a-real-api-key-but-it-needs-to-be-long"}
        };

// Simulated user behavior
export default function () {
    http.put(
      "http://localhost:9999/v1/event/" + getRandE(possible_sources) + "/" + getRandE(possible_events), 
      null,
      params);
}