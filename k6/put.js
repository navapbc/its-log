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
        { duration: "5s", target: 5 },
        { duration: "5s", target: 10 },
        { duration: "5s", target: 15 },
        { duration: "10s", target: 20 },
        { duration: "30s", target: 30 },
      ],
    },
  },
};

function getRandE(list) {
  const randomIndex = Math.floor(Math.random() * list.length);
  return list[randomIndex];
}

// Simulated user behavior
export default function () {
    const possible_events = ["event_a", "event_b", "event_c"];
    const params = {
                headers: {"x-api-key": "not-a-real-api-key-but-it-needs-to-be-long"}
            };

    http.put(
        "http://localhost:9999/v1/event/test.put.k6/" + getRandE(possible_events), 
        null,
        params)
}