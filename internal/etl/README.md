# etl

When we're done with a day's worth of data, we want to transform it from it's raw form (many events) to a useful form (condensed, calculated values).

In databases, this is called an "ETL pipeline," which is an acronym for "extract, transform, and load." 

* We extract data from one source (our event tables)
* We transform it (by counting events)
* We load it (into a new table)

## setting up

