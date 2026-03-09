import os
from random import choice
import requests
import time

ITSLOG_HOST = os.getenv('ITSLOG_HOST')
ITSLOG_PORT = os.getenv('ITSLOG_PORT')
ITSLOG_APIKEY = os.getenv('ITSLOG_APIKEY')

# 4x endpoints, 3x versions, 5x apps, and N=10 should yield 7200 events.
endpoints = ['EOB', 'Patient', 'Coverage', 'DigitalInsuranceCard']
versions = ['v1', 'v2', 'v3']
applications = [
    'Apple Research',
    'bwell',
    'CommonHealth',
    'DrOwl',
    'FastenHealth',
]


def log_url(path):
    return 'http://' + ITSLOG_HOST + ':' + ITSLOG_PORT + path


def headers(key):
    return {
        'x-api-key': key,
    }


def post(tags, value=None):
    json = {
        "tags": tags
    }
    if value is not None:
        json = json | {"value": value}
    requests.post(log_url('/v1/log'),
                  json=json,
                  headers=headers(ITSLOG_APIKEY)
                  )

# With N=20, at full speed, we get "cannot assign requested address" errors
# https://blog.paessler.com/cannot-assign-requested-address-10049-errors-in-webserver-stress-tool-under-vista
# https://stackoverflow.com/a/64330227
# This feels like a load-testing issue within/between containers in Docker.
# I put a 1-second sleep in to slow things down, and this seems to make a difference.
# Another solution is to keep N=10.
# This could be an issue in production, but it is likely we have to be under insanely high
# load for the port exhaustion to be an issue. TIME_WAIT may have to be twiddled... which
# I'm not confident we can do in a Fargate container? :shrug:


def simulate_events():
    N = 10
    for version in versions:
        # Generate N each of v1, v2, v3
        for _ in range(N):
            for ep in endpoints:
                # Generate N of each endpoint
                # without a value
                for _ in range(N):
                    post([version, ep])
                # Generate N of each endpoint
                # for each app
                for app in applications:
                    for _ in range(N):
                        post([version, ep], app)


def run_etl():
    pass


def main():
    simulate_events()
    run_etl()


if __name__ in "__main__":
    main()
