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
      action: "run",
      name: "truncate-summary",
    },
    {
      action: "run",
      name: "count-total",
    },
    {
      action: "run",
      name: "count-by-source",
    },
    {
      action: "run",
      name: "count-by-event",
    },
    {
      action: "run",
      name: "distinct-values",
    },
    {
      action: "run",
      name: "distinct-values-by-source",
    },
  ],
}
