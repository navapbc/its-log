import click
import json
import os
from pathlib import Path
import requests
import random
import time
from dotenv import dotenv_values
from datetime import datetime, timedelta
from math import floor

KEYS = []
CONFIG = {}


applications = [
    'AaNeel - CS',
    'AaNeel - CCA',
    'AaNeel - UHP',
    'Achievement',
    'AgentCubed',
    'Apple Research',
    'bwell',
    'CIG',
    'Casedok',
    'ClaimShare',
    'CommonHealth',
    'ConnectureDRX',
    'Crescendo Health',
    'DocSpera',
    'DrOwl',
    'FastenHealth',
    'HealthAgg',
    'HealthHive',
    'HealthLink Secure',
    'Kidney Choices',
    'MaxMD App',
    'myFHR',
    'PicnicHealth',
    'Project Baseline',
    'RubyWell',
    'Rush UMC',
    'Think Agent',
    'WhatMeds'
]

testclient = ['EOB', 'Patient', 'Coverage',
              'DigitalInsuranceCard', 'Profile', 'Metadata', 'OIDC']
fhir = ['Patient', 'Coverage', 'ExplanationOfBenefit']

sources = {
    'testclient.v2': testclient,
    'testclient.v3': testclient,
    'fhir.v2': fhir,
    'fhir.v3': fhir,
}


def get_key(kind):
    for k in KEYS:
        if k['kind'] == kind:
            return k
    return None


def load_env():
    global KEYS
    global CONFIG
    CONFIG = dotenv_values('/app/.env.local')
    KEYS = json.loads(CONFIG['ITSLOG_APIKEYS'])
    return CONFIG


def itslog_base_url():
    global CONFIG
    return ''.join(['http://', CONFIG['ITSLOG_SERVE_HOST'], ':', CONFIG['ITSLOG_SERVE_PORT']])


def itslog_base_headers(d={}, key='test'):
    headers = dict()
    headers['x-api-key'] = get_key(key)['key']
    for k, v in d.items():
        headers[k] = v
    return headers


def handle_response(res):
    if res.status_code >= 300:
        print(res.json())


def _load(uri, runscript, action={}):
    # print(f'-- loading: {a['filename']}')
    base = Path(runscript).parent
    contents = open(os.path.join(
        base, 'sql', action['filename'])).read()
    contents = f'-- LOADED {time.time()}\n' + contents
    url = itslog_base_url() + uri
    res = requests.post(url,
                        headers=itslog_base_headers(),
                        json={'sql': contents})
    handle_response(res)


def _run(uri, runscript, action={}):
    # print(f'-- running: {a['name']}')
    url = itslog_base_url() + uri
    res = requests.put(url, headers=itslog_base_headers())
    handle_response(res)


def _combine(uri, runscript, action={}):
    url = itslog_base_url() + uri
    handle_response(requests.put(url, headers=itslog_base_headers()))


def run_script(runscript, script, date, to_do):
    for action in script['actions']:
        if action['action'] in to_do:
            match action['action']:
                case 'message':
                    pass
                case 'load':
                    _load(
                        f'/v1/etl/{date}/{action["name"]}', runscript, action)
                case 'run':
                    _run(f'/v1/etl/{date}/{action["name"]}', runscript, action)
                case 'combine':
                    _combine(
                        f'/v1/combine/{action["source"]}/{action["destination"]}/{action["table"]}', runscript, action)
                case _:
                    print(f'-- skipping: {action["action"]}')


def _action_runner(file, date, todo):
    runscript = f'/app/e2e/{file}.json'
    script = json.load(open(runscript))
    run_script(runscript, script, date, todo)


def day_number_to_month_day(year, day_number):
    start_date = datetime(year, 1, 1).date()
    target_date = start_date + timedelta(days=day_number - 1)
    month = target_date.month
    day = target_date.day
    return month, day


def generate_events(events, date):
    for _ in range(events):
        # Every time we run, we get different distributions.
        # This way, different days/weeks look different in plotting.
        weights = map(lambda _: random.randint(1, 10), applications)
        app = random.choices(applications, weights=weights)[0]

        weights = map(lambda _: random.randint(1, 10), sources.keys())
        k = random.choices(list(sources.keys()), weights=weights)[0]

        weights = map(lambda _: random.randint(1, 10), list(sources[k]))
        source = random.choices(sources[k], weights=weights)[0]

        url = itslog_base_url() + f'/v1/dsev/{date}/{k}/{source}/{app}'
        headers = itslog_base_headers({'connection': 'close'})
        res = requests.put(url, headers=headers)
        handle_response(res)

# make jsonnet ; rm -f ../../its-log/data/itslog/*.sqlite ; ITSLOG_DAYS=365 ITSLOG_EVENTS=40000 docker compose up e2e
# ITSLOG_DAYS=1 ITSLOG_ACTIONS='combine' docker compose up e2e


@click.command()
@click.option('--actions', '-a', multiple=True)
@click.option('--events', '-e', default=5000)
@click.option('--days', '-d', default=4)
def main(actions, events, days):
    load_env()
    start = 1
    if os.getenv("ITSLOG_DAYS"):
        days = int(os.getenv("ITSLOG_DAYS"))
    if os.getenv("ITSLOG_EVENTS"):
        events = int(os.getenv("ITSLOG_EVENTS"))
    if os.getenv("ITSLOG_ACTIONS"):
        actions = os.getenv("ITSLOG_ACTIONS")
    end = days
    print(f"simulating {events} per day for {days} days", flush=True)
    time.sleep(5)
    # 'range' is exclusive of its end value.
    for day_number in range(start, end+1):
        (month, day) = day_number_to_month_day(1975, day_number)
        date = f"2026-{month:02d}-{day:02d}"
        t0 = time.time()
        if 'generate' in actions:
            generate_events(events, date)
            t1 = time.time()
            delta = t1 - t0
            print(f'{floor(delta)}s ({floor(events/delta)} EPS)', flush=True)
            # we must wait for the buffers to flush before trying to count
            time.sleep(3)
        if 'load' in actions:
            _action_runner('load', date, ['message', 'load'])
        if 'run' in actions:
            _action_runner('run', date, ['message', 'run'])
        if 'combine' in actions:
            _action_runner('combine', None, ['message', 'combine'])
        time.sleep(0.5)


if __name__ in '__main__':
    main()
