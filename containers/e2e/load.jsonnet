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
    {
      action: "load",
      name: "truncate-summary",
      filename: "truncate-summary.sql",
    },
    {
      action: "load",
      name: "count-total",
      filename: "count-total.sql",
    },
    {
      action: "load",
      name: "count-by-source",
      filename: "count-by-source.sql",
    },
    {
      action: "load",
      name: "count-by-event",
      filename: "count-by-event.sql",
    },
    {
      action: "load",
      name: "assert-source-and-total",
      filename: "assert-source-and-total.sql",
    },
    {
      action: "load",
      name: "assert-source-counts",
      filename: "assert-source-counts.sql",
    },
    {
      action: "load",
      name: "distinct-values",
      filename: "distinct-values.sql",
    },
    {
      action: "load",
      name: "distinct-values-by-source",
      filename: "distinct-values-by-source.sql",
    },
  ],
}
