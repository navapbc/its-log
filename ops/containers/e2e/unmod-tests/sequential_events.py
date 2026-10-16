from .log_create import log_create
from .params import PARAMS
from .run_etl import run_etl
from .create_etl import create_etl
from .summary_read import summary_read
from http import HTTPStatus
import sys
import time


def create_N_events(N, cluster, tags, value):
    # Create one event
    created = 0
    for _ in range(N):
        res = log_create(cluster, tags, value)
        if res.status_code == HTTPStatus.OK:
            created += 1
        else:
            print("create_N_events:", res.text)
            res.close()
            sys.exit(-1)
    return created


def test_sequential_event_logging(id):
    global PARAMS
    print(f"== test_sequential_event_logging: {id}")

    time.sleep(PARAMS.get("TIMEOUT"))
    create_etl("sql/count-total.sql")
    run_etl("count-total")

    js = summary_read("count.total")
    initial_events = -1
    if js and "count" in js:
        initial_events = js["count"]
        print(f"🔵 {initial_events} found")
    else:
        print(js)
    create_N_events(1, "", ["v1"], "")
    print(f"🟢 created 1 event")

    time.sleep(PARAMS.get("TIMEOUT"))

    # Create 100 events, pass values
    N = PARAMS.get("LESS_THAN_BUFFER")
    inserted = create_N_events(N, "", ["v2", "centi"], "")
    if N == inserted:
        print(f"🟢 created {N} events")
    else:
        print("🔴 FAILURE in test_sequential_event_logging")
        sys.exit(-1)

    # Create 5000 events, which will require 3 buffer flushes
    # Pass values and clusters
    N = PARAMS.get("MORE_THAN_BUFFER")
    inserted = create_N_events(N, "", ["v2", "centi"], "")
    if N == inserted:
        print(f"🟢 created {N} events")
    else:
        print("🔴 FAILURE in test_sequential_event_logging")
        sys.exit(-1)

    time.sleep(PARAMS.get("TIMEOUT"))

    run_etl("count-total")
    js = summary_read("count.total")
    final_events = -1
    if js and "count" in js:
        final_events = js["count"]
        print(f"🔵 {final_events} found")
    final_events = final_events - initial_events
    inserted = PARAMS.get("LESS_THAN_BUFFER") + \
        PARAMS.get("MORE_THAN_BUFFER") + 1
    if final_events == inserted:
        print(f"🟢 all events accounted for")
    else:
        print(
            f"🔴 FAILURE event insertion error: expected {inserted}, found {final_events}")

    # There should now be 5101 events in the events table.
