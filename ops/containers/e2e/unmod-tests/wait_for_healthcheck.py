from .params import PARAMS
import time


def wait_for_healthcheck():
    global PARAMS
    url = f'{PARAMS.base()}/v1/health'
    alive = False
    while not alive:
        try:
            res = PARAMS.session().get(url, timeout=PARAMS.get("TIMEOUT"))
            if res.status_code > 399:
                print('Waiting for healthcheck...')
                time.sleep(3)
            res.close()
        except:
            print('Waiting for healthcheck...')
            time.sleep(2)
        else:
            alive = True
            res.close()
