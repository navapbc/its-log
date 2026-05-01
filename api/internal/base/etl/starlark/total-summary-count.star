## Testing this locally
# starlark internal/base/etl/starlark/total-summary-count.star

# def query(table, s):
#     print(s)
#     return [{"tags": "something.v3", "count": 3 }, {"tags": "another.v2", "count": 5 }]

def summarize():
    query_rows = query("summary", "SELECT * from itslog_summary")
    total = 0
    for row in query_rows:
        print("count:", row["count"])
        total += row["count"]

    result = [
        {"operation": "total.summary", "count": total},
    ]
    print(result)
    return result

# print(summarize())