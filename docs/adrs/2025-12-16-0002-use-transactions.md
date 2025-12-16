# 2. Use buffers and bulk writes/transactions

Date: 2025-12-16

## Status

Accepted

## Context

We want event logging for 10-20 services that get hit approximately 1M times/month. This works out to 25 hits/minute, or 0.5 hits/second. With 10 services logging to the same service, that would work out to roughly 5 hits/second. It is possible that with 20 services and double the load, we might see 100 hits/second (e.g. 20 services, 2M hits/month).

We want one logger that runs in minimal memory and has performance capacity to spare. Therefore, we want to maximize throughput and be confident we have plenty of performance overhead.


The initial approach logged events as they came in (one-by-one) was very slow (no performance data). Adding a buffer improved things, but the writing of events was then a loop over the buffer, writing the events to the DB one-by-one. This provided a throughput of approximately 3600 events/second, but we would occasionally see pipeline stalls of 2-3 seconds while the events were written.

| state     | writes     | throughput |
| -- | -- | -- |
| buffered  | one-by-one | 3600 events/sec |
| buffered  | transaction | 60K events/sec |

By adding a transaction around the writing of the buffer, so that all events from the buffer are written into the database "at once," we see a 20x improvement in logging throughput. The p(90) for event handling is on the order of 500us (**micro**seconds), and the p(95) is approximately 700us. The median time to handle an event is 190us, and the maximum time taken to handle an event is now 260ms (milliseconds).

This information is all from local testing on a Mac M4. Performance in a cloud environment may vary, but the value of the approach will not change.

## Decision

We will buffer events and write the buffers as single transactions. This adds minimal complexity and maximizes throughput. 

## Consequences

* At 60K events/second we are operating at close to the performance limits of the language and hardware. This is a good thing, overall, and provides a hard upper bound on logging rates.
* Buffering events could lead to data loss in the event of a crash. Buffers are automatically flushed after &lt;N> seconds of inactivity for this reason.
* Using transactions could lose an entire buffer write (as opposed to writing events individually); however, this is probably no more or less likely than losing data when writing events one-by-one.