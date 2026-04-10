# starlark embedding in its-log

Date: 2026-04-10

## status

Accepted

## context

SQL can be cumbersome for some ETL actions, and having to compile in procedural ETLs as Golang can be a bottleneck. A scripting language would be appropriate for some classes of ETL.

## related documents

* [Starlark Language](https://starlark-lang.org/)
* [Starlark Language (Bazel)](https://bazel.build/rules/language)
* [Starlark GitHub](https://github.com/google/starlark-go)
* [Bazel](https://bazel.build/)
* [Goja](https://github.com/dop251/goja) (Javascript in Golang)
* [QuickJS in Golang](https://pkg.go.dev/modernc.org/quickjs)
* [Go-Lua](https://github.com/shopify/go-lua)

## embedded/scripting languages

It is possible, in the context of a binary application, to embed an *interpreter*. That is, even though Golang is compiled (and cannot be dynamically compiled/run), we can embed an interpreter for another programming language, and run *that* code as part of our binary.

This is not a radical idea; in a way, SQL is an interpreted language: we send SQL to a database server, and it is transformed/interpreted/executed against the data that is stored.

As initially designed, `its-log` allows ETLs to be implemented:

1. **As SQL**. This is SQL written and executed within the SQLite engine. It can only operate on data within a single database (a single file), and it cannot communicate with or operate on external systems (e.g. it cannot make network calls or write to the disk). This is "safe" in that it is sandboxed/constrained, and is appropriate for analyzing an `events` table and writing summaries to the `summary` table. ETL actions written in SQL can be loaded at runtime via API.
2. **In Golang**. Sometimes, we want to do things with data that are more complex than is easily expressed in SQL. Some things are even impossible. For example, we might want to export one or more tables as CSV, or copy an entire database file to S3. Golang ETL actions must be compiled into the binary.

The upside to SQL is that it is sandboxed; the downside is that it is limited to that sandbox. The upside to Golang is that it can do anything; the downside is that Golang ETLs must be compiled in.

This suggests a third route: embedding a general purpose programming language within `its-log`, and interpreting it at runtime.

## embedding considerations

There are multiple concerns when embedding a scripting language in a runtime environment.

1. **Sandboxing**. The embedded language is a way to do things at runtime that were not planned for. Ideally, the scripting language cannot write to the filesystem (thus altering data without going through our planned pathways), nor should it be able to make network calls (potentially pulling in malicious code or exfiltrating data). 
2. **Predictability**. While it is impossible to solve the [halting problem](https://en.wikipedia.org/wiki/Halting_problem), it is possible to design a programming language in such a way that all programs are guaranteed to halt. A scripting language with this property guarantees that we will not, in the middle of an ETL, accidentally fall into an infinite loop and never return.
3. **Native implementation**. It is possible to write an interpreter in C, and then call that C code from Go. This would mean that a C compiler would be needed to build `its-log`, which we have avoided so far. Therefore, implementations that are "pure Go" (meaning, the interpreter is written in Go only) is preferable for portability and other reasons.

There are others, but those are three big ones. 

## embedding non-choices

Generally, the Javascript and Lua choices were hard to integrate, or lacked good documentation. Further, their implementations varied in age and activity. 

By way of comparison, Starlark is used as the embedded language for [Bazel](https://bazel.build/), the build system used at Google for many of their products and systems. It is a stable language implementation, committed to remaining true to the language spec (and remaining unchanged) unless the core spec changes. It is actively maintained, and used in large production systems every day. Plus, given its Pythonic nature, it felt like a good fit for the teams who will likely use `its-log` and systems they work on.

## embedding choice: starlark

Three languages were considered for embedding: Javascript, Lua, and Starlark. Of these, we chose Starlark.

1. Starlark is deterministic. Every single Starlark program is guaranteed to terminate; it is not possible to express infinite loops or recursion in Starlark.
2. Starlark is hermetic. It is not possible to communicate with the filesystem, network, or any other systems outside of its execution sandbox.
3. Starlark supports parallelism. Not essential in our current usecase, but may become a benefit for performance reasons.
4. Starlark is immutable by default. Starlark makes all variables immutable *except* for dictionaries and lists. This property forces a slightly different style of programming, but is generally accepted as a safer way to develop code (and helps Starlark enforce some of its other properties).
5. Starlark is Python-like. The language looks and behaves like Python, with garbage collection and everything. It is *not* Python, but for developers familiar with Python, it will behave "Python-like" enough that it should be easy to quickly write data manipulation code in the language.

## starlark: example

Our embedding of Starlark passes one function in from the outside world: `query()`. This is a Golang function that can be called from within the Starlark interpreter. It consumes a string and returns an array of JSON objects. That string is executed as SQL against the current `its-log` SQLite database. In this way, a Starlark ETL author can choose what data they want to analyze.

That query comes back as a list of dictionaries; each object represents a row from the `itslog_events` table, with keys matching the columns.

In the example below, we 

1. iterate through each row, 
2. count how many vowels are in the `tags` column
3. count how many consonants are in the `tags` column
4. construct and return a list of dictionaries representing `itslog_summary` rows 

```python
def summarize():
    query_rows = query("SELECT * from itslog_events")
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
```

`its-log` then takes the array of dictionaries and converts those into `itslog_summary` rows, locks the database, and inserts (or replaces) all of the rows returned. 

We intentionally "breach" the Starlark sandbox by providing the `query()` function to the runtime environment. We avoid a second breach by expecting data to be returned in a specific format, which we then parse and insert (on the ETL author's behalf) into the database. 

## decision

Embed Starlark for ETL authoring.