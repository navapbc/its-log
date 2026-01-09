{
  actions: [
    // messages take one param
    // everything can have a "message" option
    {
      action: "message",
      message: "HELO",
    },
    // WARNING: The first thing we do is truncate the entire summary table.
    // Why? To make sure each pipeline run is clean. Perhaps this is a bad idea.
    {
      action: "sql",
      filename: "truncate-summary.sql",
    },
    {
      action: "sql",
      filename: "count-total.sql"
    },
    {
      action: "sql",
      filename: "count-by-source.sql",
    },
    {
      action: "sql",
      filename: "count-by-event.sql",
    },
    {
      action: "assert",
      filename: "assert-source-and-total.sql",
    },
    {
      action: "assert",
      filename: "assert-source-counts.sql",
    },
    {
      action: "sql",
      filename: "distinct-values.sql",
    },
    {
      action: "sql",
      filename: "distinct-values-by-source.sql",
    },
    // fileCopy copies from one location to another
    {
      action: "fileCopy",
      source: "env.sourcePath",
      destination: "env.destPath",
    },
    // toS3
    // fromS3
    // DBtoDB
    // DBtoFile
  ],
}
