import json
import os
from pathlib import Path
import requests
import random
import time

KEYS = []


def get_key(kind):
    for k in KEYS:
        if k["kind"] == kind:
            return k
    return None


def load_env():
    global KEYS
    with open("/app/.env.local", "r") as ip:
        for line in ip:
            if "=" in line:
                line = line.replace("'", "")
                line = line.replace("\n", "")
                # if it isn't JSON...
                if line.count('"') == 2:
                    line = line.replace('"', "")
                pieces = line.split("=")
                # print(pieces[0], "=", pieces[1])
                os.environ[pieces[0]] = pieces[1]
    KEYS = json.loads(os.environ["ITSLOG_APIKEYS"])


def run_script(runscript, script, date, to_do):
    for a in script["actions"]:
        if a["action"] in to_do:
            match a["action"]:
                case "message":
                    # print(f"-- {a['message']} --")
                    pass
                case "load":
                    # print(f"-- loading: {a['filename']}")
                    base = Path(runscript).parent
                    contents = open(os.path.join(
                        base, "sql", a['filename'])).read()
                    url = "http://" + \
                        os.environ["ITSLOG_SERVE_HOST"] + ":" + \
                        os.environ["ITSLOG_SERVE_PORT"] + \
                        f"/v1/etl/{date}/{a['name']}"
                    key = get_key("log")["key"]
                    headers = {
                        "x-api-key": key
                    }
                    res = requests.post(url, headers=headers, json={
                        "sql": contents,
                    })
                    # print(res.json())
                case "run":
                    # print(f"-- running: {a['name']}")
                    base = Path(runscript).parent
                    url = "http://" + \
                        os.environ["ITSLOG_SERVE_HOST"] + ":" + \
                        os.environ["ITSLOG_SERVE_PORT"] + \
                        f"/v1/etl/{date}/{a['name']}"
                    key = get_key("log")["key"]
                    headers = {
                        "x-api-key": key
                    }
                    res = requests.put(url, headers=headers)
                    # print(res.json())

                case _:
                    print(f"-- skipping: {a['action']}")


def load_etl(date):
    runscript = "/app/e2e/load.json"
    script = json.load(open(runscript))
    run_script(runscript, script, date, ["message", "load"])


def run_etl(date):
    runscript = "/app/e2e/run.json"
    script = json.load(open(runscript))
    run_script(runscript, script, date, ["message", "run"])


applications = [
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

testclient = ["EOB", "Patient", "Coverage",
              "DigitalInsuranceCard", "Profile", "Metadata", "OIDC"]
fhir = ["Patient", "Coverage", "ExplanationOfBenefit"]

sources = {
    "testclient.v2": testclient,
    "testclient.v3": testclient,
    "fhir.v2": fhir,
    "fhir.v3": fhir,
}


def generate_events(n, date):
    for i in range(1, n):
        app = random.choice(applications)
        k = random.choice(list(sources.keys()))
        source = random.choice(sources[k])

        url = "http://" + \
            os.environ["ITSLOG_SERVE_HOST"] + ":" + \
            os.environ["ITSLOG_SERVE_PORT"] + \
            f"/v1/dse/{date}/{k}.{source}/{app}"
        key = get_key("test")["key"]
        headers = {
            "x-api-key": key
        }
        res = requests.put(url, headers=headers)


def main():
    load_env()
    for i in range(1, 30):
        date = f"2026-01-{i:02d}"
        generate_events(1000, date)
        load_etl(date)
        run_etl(date)
        time.sleep(0.5)


if __name__ in "__main__":
    main()
