import http from "k6/http";
import { check, sleep } from "k6";

const ITSLOG_HOST = __ENV.ITSLOG_HOST;
const ITSLOG_PORT = __ENV.ITSLOG_PORT;

// https://k6.io/blog/how-to-generate-a-constant-request-rate-with-the-new-scenarios-api/
export const options = {
    thresholds: {
        // http errors should be less than 1%
        http_req_failed: ['rate<0.01'],
        // 95% of requests should be below 2ms
        http_req_duration: ['p(95)<2'],
    },
    insecureSkipTLSVerify: true,
};

// CONSTANTS
// names
const patient_names = []
for (const source_number of Array(30).keys()) {
    patient_names.push("Alice." + source_number)
}
for (const source_number of Array(30).keys()) {
    patient_names.push("Bob." + source_number)
}


const endpoints = ["EOB", "Patient", "Coverage", "DigitalInsuranceCard"];
const versions = ["v1", "v2", "v3"];
const applications = [
    "AaNeel - UHP",
    "AgentCubed",
    "Apple Research",
    "bwell",
    "Casedok",
    "CommonHealth",
    "Crescendo Health",
    "DrOwl",
    "FastenHealth",
    "Kidney Choices",
    "myFHR",
    "RubyWell",
    "WhatMeds"
]

const url = "http://" + ITSLOG_HOST + ":" + ITSLOG_PORT + "/v1/log";
// Simulated user behavior
export default function () {
    const N = 10;

    const params = {
        headers: { "x-api-key": __ENV.ITSLOG_APIKEY }
    };

    versions.forEach((v) => {
        endpoints.forEach((e) => {
            applications.forEach((a) => {
                for (let i = 0; i < N; ++i) {
                    var data = {
                        "cluster": i.toString(),
                        "tags": [v, e]
                    };
                    if ((i % 2) == 0) {
                        data["value"] = a;
                    }
                    var resp = http.post(
                        url,
                        JSON.stringify(data),
                        params);
                    if (resp.status > 299) {
                        console.log(url);
                        console.log(data);
                        console.log(params);
                        console.log(resp.json());
                        sleep(10);
                    }
                }
            });
        });
    });

    sleep(2);
};