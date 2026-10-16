from pathlib import Path
from .params import PARAMS
import base64
import requests
import sys


def get_contents(filename):
    contents = open(filename, 'r').read()
    encoded = base64.b64encode(contents.encode("utf-8")).decode("ascii")
    if len(encoded) % 4 != 0:
        the_length = len(encoded)
        mod = the_length % 4
        encoded = encoded + "="*mod
    return encoded


def create_etl(filename):
    encoded = get_contents(filename)

    if ".sql" in filename:
        kind = "sql"
    elif ".star" in filename:
        kind = "starlark"
    else:
        kind = "BAD_ETL_TYPE"

    body = {
        "name": Path(filename).stem,
        "kind": kind,
        "body": encoded,
        "date": "today"
    }

    try:
        url = f'{PARAMS.base()}/v1/etl/create'
        res: requests.Response = PARAMS.session().post(url,
                                                       headers=PARAMS.headers(
                                                           "PUPPERADMIN"),
                                                       json=body,
                                                       timeout=PARAMS.get("TIMEOUT"))
        if res.status_code > 399:
            print("create_etl:", res.status_code)
            print(res.text)
            res.close()
            sys.exit(-1)
        if res.status_code == 201:
            return res.status_code
        res.close()
    except requests.exceptions.Timeout as e:
        print(e)


def test_create_etl(filename):
    global PARAMS
    print(f"== test_create_etl {filename}")
    code = create_etl(filename)
    print(f"🟢 test_create_etl {code}")
