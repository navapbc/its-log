# effective logging

The `its-log` approach to logging (event-based, valueless) requires some thought about what and how you log.

At some level, it does not matter how the framework is used; it impacts summarization, at best.

## example

The simplest `its-log` model is a source and event. For example, in an API, the source might be `v1` and the event might be the endpoint being called.

| source | event |
| ------ | ----- |
| v1     | user  |

The problem then is that more information might be desirable. If the API follows a traditional RESTful design, it is possible that the HTTP verbs GET, POST, DELETE, and PUT all play a role on that endpoint. There are two ways to capture this. One is to put the verb in the *source*:

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

When we put it the source, we end up with four different sources, each with the same event. 


The first example would generate a summary that looked like the following:

| operation    | source    | count   |
| ------------ | --------- | ------- |
| count.by_day | v1.GET    | 394343  |
| count.by_day | v1.POST   | 1894834 |
| count.by_day | v1.DELETE | 1232    |
| count.by_day | v1.PUT    | 5       |

In the second, the same count of sources-by-day would look like the following:

| operation    | source | count   |
| ------------ | ------ | ------- |
| count.by_day | v1     | 2290414 |

If we add a count, by-day, by-event, we then see the breakout. In the first example, we would see:

| operation             | source    | event | count   |
| --------------------- | --------- | ----- | ------- |
| count.by_day          | v1.GET    | NULL  | 394343  |
| count.by_day          | v1.POST   | NULL  | 1894834 |
| count.by_day          | v1.DELETE | NULL  | 1232    |
| count.by_day          | v1.PUT    | NULL  | 5       |
| count.by_day.by_event | v1.GET    | user  | 394343  |
| count.by_day.by_event | v1.POST   | user  | 1894834 |
| count.by_day.by_event | v1.DELETE | user  | 1232    |
| count.by_day.by_event | v1.PUT    | user  | 5       |

The second approach would have less duplication:

| operation             | source | event       | count   |
| --------------------- | ------ | ----------- | ------- |
| count.by_day          | v1     | NULL        | 2290414 |
| count.by_day.by_event | v1     | GET.user    | 394343  |
| count.by_day.by_event | v1     | POST.user   | 1894834 |
| count.by_day.by_event | v1     | DELETE.user | 1232    |
| count.by_day.by_event | v1     | PUT.user    | 5       |

## what does this mean?

In the first example, we had more sources, which yielded duplication in our summaries:

| source    | event |
| --------- | ----- |
| v1.GET    | user  |
| v1.POST   | user  |
| v1.DELETE | user  |
| v1.PUT    | user  |

| operation             | source    | event | count   |
| --------------------- | --------- | ----- | ------- |
| count.by_day          | v1.GET    | NULL  | 394343  |
| count.by_day          | v1.POST   | NULL  | 1894834 |
| count.by_day          | v1.DELETE | NULL  | 1232    |
| count.by_day          | v1.PUT    | NULL  | 5       |
| count.by_day.by_event | v1.GET    | user  | 394343  |
| count.by_day.by_event | v1.POST   | user  | 1894834 |
| count.by_day.by_event | v1.DELETE | user  | 1232    |
| count.by_day.by_event | v1.PUT    | user  | 5       |

At first, this appears redundant---all the values are the same---but if we had more endpoints, putting the verb in the source starts to make more sense.

| source    | event |
| --------- | ----- |
| v1.GET    | song  |
| v1.POST   | song  |
| v1.DELETE | song  |
| v1.PUT    | song  |

Now, the summary tables might look like this:

| operation             | source    | event | count   |
| --------------------- | --------- | ----- | ------- |
| count.by_day          | v1.GET    | NULL  | 8272673 |
| count.by_day          | v1.POST   | NULL  | 2342342 |
| count.by_day          | v1.DELETE | NULL  | 12325   |
| count.by_day          | v1.PUT    | NULL  | 6       |
| count.by_day.by_event | v1.GET    | user  | 394343  |
| count.by_day.by_event | v1.POST   | user  | 1894834 |
| count.by_day.by_event | v1.DELETE | user  | 1232    |
| count.by_day.by_event | v1.PUT    | user  | 1       |
| count.by_day.by_event | v1.GET    | song  | 7878330 |
| count.by_day.by_event | v1.POST   | song  | 447508  |
| count.by_day.by_event | v1.DELETE | song  | 11093   |
| count.by_day.by_event | v1.PUT    | song  | 1       |

We have a concise summary, by-day, of how many verb hits we have on the v1 endpoint. We can see that, once we take the `song` endpoint into account, there are far more GET actions than POST, which was not obvious from only the `user` endpoint. The breakdown by *event* then shows us the counts per endpoint for each verb. Using this table, it is straight forward to answer questions like:

1. How many GET/POST/DELETE/PUT events were there on the v1 endpoint each day?
    `SELECT sum(count) FROM summary WHERE operation = 'count.by_day' AND source LIKE 'v1%`
1. How many events were on the `user` endpoint? 
    `SELECT sum(count) FROM summary where event='user'`

The second approach yields different results. We would still have a top-level summary that excluded the verbs.

| operation    | source | count   |
| ------------ | ------ | ------- |
| count.by_day | v1     | 2290414 |

The summaries would now look like this:

| operation             | source | event       | count    |
| --------------------- | ------ | ----------- | -------- |
| count.by_day          | v1     | NULL        | 10627346 |
| count.by_day.by_event | v1     | GET.user    | 394343   |
| count.by_day.by_event | v1     | POST.user   | 1894834  |
| count.by_day.by_event | v1     | DELETE.user | 1232     |
| count.by_day.by_event | v1     | PUT.user    | 5        |
| count.by_day.by_event | v1     | GET.song    | 7878330  |
| count.by_day.by_event | v1     | POST.song   | 447508   |
| count.by_day.by_event | v1     | DELETE.song | 11093    |
| count.by_day.by_event | v1     | PUT.song    | 1        |

The same questions now are answered slightly differently:

1. How many GET/POST/DELETE/PUT events were there on the v1 endpoint each day?
    `SELECT count FROM summary WHERE operation = 'count.by_day'`
2. How many events were on the `user` endpoint? 
    `SELECT sum(count) FROM summary where event LIKE '%user'`

## distinctness matters

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

## clustering

Finally, it is possible that a *group* of events are releated; it might be that we want to log more than one thing at a single moment in time. These would then be a *cluster* of events.

To do this, we would expand our event with one more value: a `cluster` identifier.

| cluster | source  | event | value    |
| ------- | ------- | ----- | -------- |
| 53      | v1.GET  | user  | 0a3c     |
| 53      | api     | RAM   | 3000000  |
| 174     | v1.POST | user  | 45b1     |
| 174     | api     | RAM   | 34999100 |
| NULL    | v1.POST | user  | 0a3c     |
| ...     | ...     | ...   | ...      |

Here, we see the first two events are *clustered*. Even if their timestamps are (slightly) different, they represent log state that is considered to be "at the same moment." We might capture a user at the same time as we ask what the RAM usage on the API was. Generic/general analysis of clusters is probably not possible; that is, every application will use clusters differently. However, dashboards might use this information to ask "what sources use the most RAM?" On one hand, these could be independent events (e.g. we could log it without clustering), but we would lose the connection through the cluster id to the fact that user `0a3c` made a query that used much less RAM than user `45b1`.

*I am very skeptical of the utility of clustering. I think I'll leave it out of v1.*

