from .params import PARAMS
from .create_etl import create_etl
from .log_create import log_create
from .run_etl import run_etl
from .summary_read import summary_read

import requests
import json
import sys
import time


def test_pocketchange():
    global PARAMS
    NUMBER_OF_PENNIES = 100
    existing_pennies = 0

    create_etl('starlark/sum-pocket-change.star')
    run_etl("sum-pocket-change")
    js = summary_read("total.pocketchange")

    if js and "count" in js:
        existing_pennies = js["count"]
        print(f"🟢 found {NUMBER_OF_PENNIES} pennies pre-existing")
    else:
        print(f"🟢 found no pennies")

    for i in range(NUMBER_OF_PENNIES):
        res: requests.Response = log_create(None, ["pocketchange"], json.dumps({
            "pennies": 1
        }))
        if res.status_code > 399:
            res.close()
            sys.exit(-1)
    res.close()

    # We have to wait until the buffer flushes to count the pennies.
    time.sleep(3)

    run_etl("sum-pocket-change")
    js = summary_read("total.pocketchange")

    if js["count"] == NUMBER_OF_PENNIES + existing_pennies:
        print(f"🟢 found {NUMBER_OF_PENNIES + existing_pennies} pennies")
    else:
        print("🔴 bad penny")
