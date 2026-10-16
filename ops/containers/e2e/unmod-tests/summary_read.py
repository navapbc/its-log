from .params import Params, PARAMS

import requests
import sys


def summary_read(operation):
    global PARAMS
    url = f'{PARAMS.base()}/v1/summary/read'
    try:
        res: requests.Response = PARAMS.session().post(url,
                                                       headers=PARAMS.headers(
                                                           "PUPPERADMIN"),
                                                       json={
                                                           "operation": operation,
                                                           "date": "today",
                                                           "tags": "%",
                                                       },
                                                       timeout=PARAMS.get("TIMEOUT"))
        if res.status_code in [200, 201]:
            return res.json()
        if res.status_code > 399:
            print("🔴 summary_read:", res.status_code, res.text)
            print(res.text)
            res.close()
            sys.exit(-1)
        res.close()
    except requests.exceptions.Timeout as e:
        res.close()
        sys.exit(-1)
