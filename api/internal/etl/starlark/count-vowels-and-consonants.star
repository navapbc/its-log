## Testing this locally
# starlark internal/base/etl/starlark/count-vowels-and-consonants.star
# Uncomment the query function for local testing (to produce mock data)
# and uncomment the call to summarize(). This is a rough simulation
# of what will happen in production.
# it is expected that an array of dictionaries matching a summary row
# will be returned (keys "operation" and "count" are required!)

# def query(s):
#     print(s)
#     return [{"tags": "something.v3"}, {"tags": "another.v2"}]

def summarize():
    query_rows = query("events", "SELECT * from itslog_events")
    summaries = []
    total_vowels = 0
    total_consonants = 0
    for row in query_rows:
        for l in "aeiou".elems():
            if l in row["tags"]:
                total_vowels += 1
        for l in "bcdfghjklmnpqrstvwxyz".elems():
            if l in row["tags"]:
                total_consonants += 1

    result = [
        {"operation": "total.vowels", "count": total_vowels},
        {"operation": "total.consonants", "count": total_consonants}
    ]
    return result

# summarize()