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
      gracefulRampDown: "2s",
      stages: [
        { duration: "20s", target: 30 },
      ],
    },
  },
  insecureSkipTLSVerify: true,
};

function getRandE(list) {
  const randomIndex = Math.floor(Math.random() * list.length);
  return list[randomIndex];
}

function getRandD(d) {
  var keyAtIndex = getRandE(Object.keys(d))
  return { "source": keyAtIndex, "event": getRandE(d[keyAtIndex]) };
}

const patient_names = []
for (const source_number of Array(30).keys()) {
  patient_names.push("Alice." + source_number) 
}
for (const source_number of Array(30).keys()) {
  patient_names.push("Bob." + source_number) 
}

const params = {
            headers: {"x-api-key": "not-a-real-api-key-but-it-needs-to-be-long"}
        };


const generateHash = (string) => {
  let hash = 0;
  for (const char of string) {
    hash = (hash << 5) - hash + char.charCodeAt(0);
    hash |= 0; // Constrain to 32bit integer
  }
  return hash;
};

function addId(s) {
    return s + "-" + generateHash(getRandE(patient_names)).toString(16)
}

// Generate an event that is "authentic" to our particular context.
function generateEvents() {
    const root = "blue";
    var sources = {
        "testclient": ["EOB", addId("Patient"), "Coverage", "DigitalInsuranceCard"],
        "fhir.v2": [addId("Patient"), "Coverage", "ExplanationOfBenefit"],
        "fhir.v3": [addId("Patient"), "Coverage", "ExplanationOfBenefit"],
    }
    // From https://www.medicare.gov/providers-services/claims-appeals-complaints/claims/share-your-medicare-claims/connected-apps
    const applications = [
        "AaNeel - CS", 
        "AaNeel - CCA", 
        "AaNeel - UHP", 
        "Achievement", 
        "AgentCubed", 
        "Apple Research",
        "bwell",
        "CIG", 
        "Casedok",
        "ClaimShare",
        "CommonHealth",
        "ConnectureDRX",
        "Crescendo Health",
        "DocSpera",
        "DrOwl",
        "FastenHealth",
        "HealthAgg",
        "HealthHive",
        "HealthLink Secure",
        "Kidney Choices",
        "MaxMD App",
        "myFHR",
        "PicnicHealth",
        "Project Baseline",
        "RubyWell",
        "Rush UMC",
        "Think Agent",
        "WhatMeds"
    ]

    var e = getRandD(sources)

    return [
        { "source": root + "." + e.source, "event": e.event },
        { "source": root + ".application", "event": getRandE(applications)}
    ]

};

// Simulated user behavior
export default function () {
    generateEvents().forEach((e) => 
        http.put(
        "https://localhost:8443/v1/event/" + e.source + "/" + e.event, 
        null,
        params)  
    );
}