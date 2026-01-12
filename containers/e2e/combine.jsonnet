{
  server: {
    url: "https://localhost:8443/v1/etl",
  },
  actions: [
    // messages take one param
    // everything can have a "message" option
    {
      action: "message",
      message: "HELO",
    },
  ] + std.map(function(npad) {
    action: "combine",
    name: "combine-summaries",
    table: "itslog_summary",
    source: "2026-01-" + npad,
    destination: "summary",
  }, std.map(function(n) if n < 10
  then "0" + std.toString(n)
  else std.toString(n), std.range(1, 29))),
}
