# effective logging

`its-log` is both a simple and opinionated logger. In order to use it effectively, it is best to think about log events in terms of *categories* and *uniqueness*.


## sources and events are categories

When logging a source and event, you are logging two *categories*. For example, in an API, the source might be `v1` and the event might be the endpoint being called.

| source | event |
| ------ | ----- |
| v1     | user  |
| v1     | song  |
| v2     | song  |
| v1     | song  |
| ...    | ...   |

`its-log` has a lightweight ETL pipeline integrated for summarization. When working with categories, it is easy to have "generic" ETL transforms that generate summaries on a category-by-category basis. Using the above as an example a summary table might look like:

| operation                       | source | event | count |
| ------------------------------- | ------ | ----- | ----- |
| count.by_source.by_day          | v1     | NULL  | 4235  |
| count.by_source.by_event.by_day | v1     | user  | 1034  |
| count.by_source.by_event.by_day | v1     | song  | 3189  |
| count.by_source.by_event.by_day | v2     | song  | 12    |

Without knowledge of the application doing the logging, it is possible to generate categorical counts automatically. That is, we do not need to know the particular sources or events involved; we can treat them as categories and count events in the table grouped by these values.

Even if we later decide we want to include HTTP verbs in our logging of our events, it largely does not matter where the verb goes. We can attach the verb to the source *or* to the event. 

| source    | event |
| --------- | ----- |
| v1.GET    | user  |
| v1.POST   | user  |
| v1.DELETE | user  |
| v1.PUT    | user  |

Or, the verb can go in the event:

| source | event       |
| ------ | ----------- |
| v1     | GET.user    |
| v1     | POST.user   |
| v1     | DELETE.user |
| v1     | PUT.user    |

Either way, we are defining *categories*, and our automatically-generated summaries will "make sense."


## values are unique

As seen, this largely does not matter with tables where the source and event are largely *categorical*. That is, we have three categories:

* The version of the API (v1, v2)
* The HTTP verb (GET, POST, etc.)
* The endpoint (user, song, etc.)

None of these represent *unique* things. Therefore, shifting the categorical values between the columns has some effect, but in terms of ETL and summarization, nothing radical changes. We are not in danger of losing information, and it has a negligible impact on the kinds of queries we write in order to analyze the event data.

Where this will matter is when it comes to capturing *distinct* or *unique* events. Now, consider a unique user ID. Assume it is salted and hashed, but it comes thorugh to us as a unique, four-chracter code (e.g. `0a3c` or similar). Where we put this now matters *a lot*. 

If we put the unique value in the *source*, we have made it far more complex to collate our sources.

| source       | event |
| ------------ | ----- |
| v1.GET.0a3c  | user  |
| v1.POST.45b1 | user  |
| v1.POST.0a3c | user  |
| ...          | user  |

Now, every time user `0a3c` hits an endpoint, we have a new, unique source. This is *bad*. We could instead put the user ID in the *event* column.

| source  | event     |
| ------- | --------- |
| v1.GET  | user.0a3c |
| v1.POST | user.45b1 |
| v1.POST | user.0a3c |
| ...     | user      |

This at least lets us count sources in our summaries again, but we now have to do more work to count events. That is, we used to have a summary table that looked like this:


| operation             | source    | event     | count |
| --------------------- | --------- | --------- | ----- |
| count.by_day.by_event | v1.GET    | user.0a3c | 18    |
| count.by_day.by_event | v1.POST   | user.0a3c | 234   |
| count.by_day.by_event | v1.DELETE | user.0a3c | 7     |
| count.by_day.by_event | v1.GET    | user.45b1 | 3     |
| count.by_day.by_event | v1.PUT    | user.45b1 | 12    |
| ...                   | ...       | ...       | ...   |

By putting a unique value in the event column, we will end up with one `count.by_day.by_event` row for every single user. In order to count up the number of events that were on the `user` endpoint, we have to do something like:

`SELECT sum(count) FROM summary WHERE operation = 'count.by_day.by_event' and event LIKE 'user%'` However, as a summary table, it implies that we're more interested in the user than we are the categorical event. So, we can expand our logging model to include a *value*. Now, the event logging table looks like


| source  | event | value |
| ------- | ----- | ----- |
| v1.GET  | user  | 0a3c  |
| v1.POST | user  | 45b1  |
| v1.POST | user  | 0a3c  |
| ...     | ...   | ...   |

This lets us once again create categorical summaries. However, we can also now automatically generate a new class of summary data. The `value` column, if assumed to be unique, generate summary data rows that answer questions like:

1. How many unique values were there per day?
    `SELECT count(distinct value) FROM events`
1. How many unique values were there per source?
    `SELECT source, count(distinct value) FROM events GROUP BY source`
2. How many unique values were there per event type?
    `SELECT event, count(distinct value) FROM events GROUP BY event`

This yields new summary rows like:

| operation                      | source    | event | count |
| ------------------------------ | --------- | ----- | ----- |
| unique.values.by_day           | NULL      | NULL  | 38125 |
| unique.values.by_source.by_day | v1.GET    | NULL  | 10843 |
| unique.values.by_source.by_day | v1.POST   | NULL  | 4843  |
| unique.values.by_source.by_day | v1.DELETE | NULL  | 9303  |
| unique.values.by_source.by_day | v1.PUT    | NULL  | 2293  |
| unique.values.by_event.by_day  | NULL      | user  | 8333  |
| unique.values.by_event.by_day  | NULL      | song  | 29792 |
| ...                            | ...       | ...   | ...   |
