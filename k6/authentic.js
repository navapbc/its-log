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
      startVUs: 10,
      gracefulRampDown: "2s",
      stages: [
        { duration: "2s",  target: 30 },
        { duration: "5s", target: 60 },
        { duration: "5s", target: 60 },
        { duration: "2s", target: 30 },
      ],
    },
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
// API PARAMS
const params = {
            headers: {"x-api-key": "not-a-real-api-key-but-it-needs-to-be-long"}
        };


// HELPER FUNS
function getRandE(list) {
  const randomIndex = Math.floor(Math.random() * list.length);
  return list[randomIndex];
}

function getRandD(d) {
  var keyAtIndex = getRandE(Object.keys(d))
  return { "source": keyAtIndex, "event": getRandE(d[keyAtIndex]) };
}

const generateHash = (string) => {
  let hash = 0;
  for (const char of string) {
    hash = (hash << 5) - hash + char.charCodeAt(0);
    hash |= 0; // Constrain to 32bit integer
  }
  return hash;
};


// Generate an event that is "authentic" to our particular context.
function generateEvents() {
    const root = "blue";
    var testclient = ["EOB", "Patient", "Coverage", "DigitalInsuranceCard", "Profile", "Metadata", "OIDC"];
    var fhir = ["Patient", "Coverage", "ExplanationOfBenefit"];

    var sources = {
        "testclient.v2": testclient,
        "testclient.v3": testclient,
        "fhir.v2": fhir,
        "fhir.v3": fhir,
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

    var e = getRandD(sources);
    var app = getRandE(applications);
    var bene = generateHash(getRandE(patient_names)).toString(16);

    return [
        // Now, blue.endpoint_app.{source}.{event}/{app}
        { 
          "endpoint": "sev",
          "source": [ root, "endpoint", "by_app", e.source ].join("."), 
          "event": e.event,
          "value": app,
        },
        { 
          "endpoint": "sev",
          "source": [ root, "endpoint", "by_user", e.source ].join("."),
          "event": app,
          "value": bene,
        },
    ]

};

const endpoints = {
  "cse": (e) => [ e.endpoint, e.cluster, e.source, e.event ].join("/"),
  "csev": (e) => [ e.endpoint, e.cluster, e.source, e.event, e.value ].join("/"),
  "sev": (e) => [ e.endpoint, e.source, e.event, e.value ].join("/"),
  "se": (e) => [ e.endpoint, e.source, e.event ].join("/"),
};

// Simulated user behavior
export default function () {
    generateEvents().forEach(function (e) {       
        http.put(
        // "https://localhost:8443/v1/event/" + e.source + "/" + e.event, 
        "https://localhost:8443/v1/" + endpoints[e.endpoint](e),
        null,
        params)  
});
}