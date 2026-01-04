{
  actions: [
    // messages take one param
    // everything can have a "message" option
    {
      action: "message",
      message: "helo",
    },
    {
      action: "sql",
      filename: "truncate-summary.sql",
    },
    // SQL actions should take a connection
    // and they should be defined in advance.
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
