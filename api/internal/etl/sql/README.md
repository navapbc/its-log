# SQL for ETL

SQL in this directory is compiled into the application and loaded into the database as default ETL actions. Therefore, they are designed to be generic against a reasonable set of use cases.

Every action is assumed to return a 0 or 1 value. Therefore:

- ETL actions that insert values into the summary table should explicitly end with a `RETURNING 0` to indicate success. If a failure happens prior to that, an SQL error will be raised and caught via other mechanisms.
- ETL actions that assert correctness in the pipeline should return 0 or 1. This is useful in a sequence, and a 0 would imply success. 

Failure to return a value will raise an error back to the end-user of the API.