## Testing this locally
# starlark internal/base/etl/starlark/sum-pocket-change.star

# def query(_, s):
#     return [
#         {"tags": "pocketchange", "value": '{"pennies": 3, "nickels": 5}'}, 
#         ]

def summarize():
    query_rows = query("events", "SELECT * from itslog_events WHERE tags = 'pocketchange'")
    sums = {
        "pennies": 0,
        "nickels": 0,
        "dimes": 0,
        "quarters": 0,
    }
    for row in query_rows:
        # Convert the value field from JSON to a dictionary
        d = json.decode(row["value"])
        # Sum up what we find
        for coin in ["pennies", "nickels", "dimes", "quarters"]:
            if coin in d:
                sums[coin] = sums[coin] + d[coin]
    total = sums["pennies"] + (sums["nickels"] * 5) + (sums["dimes"] * 10) + (sums["quarters"] * 25)
    result = [
        {"operation": "total.pocketchange", "count": total},
        {"operation": "total.pocketchange.json", "tags": json.encode(sums), "count": 0}
    ]
    return result

# print(summarize())