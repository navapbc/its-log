from .params import Params, PARAMS
from .sequential_events import create_N_events
from .run_etl import run_etl
from .summary_read import summary_read
from .create_etl import create_etl
import time
from multiprocessing import Pool, Process


def init_process():
    global PARAMS
    PARAMS = Params()


def test_parallel_inserts(N):
    global PARAMS
    print(f"test_parallel_inserts {N}")
    p1 = Pool(processes=N, initializer=init_process)
    p2 = Pool(processes=N, initializer=init_process)

    run_etl("count-total")
    js = summary_read("count.total")
    initial_events = -1
    if js and "count" in js:
        initial_events = js["count"]
        print(f"🔵 {initial_events} found")
    else:
        print(js)

    params = []
    for _ in range(N):
        params.append([N*1000, "", ["P", "parallel"], ""])
    r1 = p1.map_async(create_N_events, params)
    p1.close()
    r2 = p2.map_async(create_etl, ['sql/sentinel.sql'] * N * N)
    p2.close()
    p1.join()
    p2.join()

    time.sleep(PARAMS.get("TIMEOUT"))

    run_etl("count-total")
    js = summary_read("count.total")
    final_events = -1
    if js and "count" in js:
        final_events = js["count"]
        print(f"🔵 {final_events} found")
    total = final_events - initial_events
    print(f"parallel initial {initial_events} {final_events} {total}")
