import click
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
                    contents = f"-- LOADED {time.time()}\n" + contents
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
                    if res.status_code >= 300:
                        print(res.json())
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
                    if res.status_code >= 300:
                        print(res.json())
                case "combine":
                    base = Path(runscript).parent
                    url = "http://" + \
                        os.environ["ITSLOG_SERVE_HOST"] + ":" + \
                        os.environ["ITSLOG_SERVE_PORT"] + \
                        f"/v1/combine/{a['source']}/{a['destination']}/{a['table']}"
                    key = get_key("log")["key"]
                    headers = {
                        "x-api-key": key
                    }
                    res = requests.put(url, headers=headers)
                    if res.status_code >= 300:
                        print(res.json())
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


def combine():
    runscript = "/app/e2e/combine.json"
    script = json.load(open(runscript))
    run_script(runscript, script, "no date", ["message", "combine"])


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
            f"/v1/dsev/{date}/{k}/{source}/{app}"
        key = get_key("test")["key"]
        headers = {
            "x-api-key": key,
            'Connection': 'close'
        }
        res = requests.put(url, headers=headers)
        if res.status_code >= 300:
            print(res.json())


@click.command()
@click.option('--actions', '-a', multiple=True)
@click.option('--events', '-e', default=40000)
@click.option("--days", '-d', default=4)
def main(actions, events, days):
    print(actions, events, days)
    events = events
    load_env()
    for i in range(1, 1+days):
        date = f"2026-01-{i:02d}"
        # 40K/day is authentic
        t0 = time.time()
        if "generate" in actions:
            generate_events(events, date)
        t1 = time.time()
        delta = t1 - t0
        print(f"{delta}s ({events/delta} events per second)", flush=True)
        # we must wait for the buffers to flush before trying to count
        time.sleep(3)
        if "load" in actions:
            load_etl(date)
        if "run" in actions:
            run_etl(date)
        if "combine" in actions:
            combine()
        time.sleep(0.5)


if __name__ in "__main__":
    main()
