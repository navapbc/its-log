## 2025-12-16 with dictionaries

A naive approach to writing the dictionary slows us down. Now at 36K events/second. Which... makes sense. I'm doing 2x the writes to the DB, so it slowed down by a factor of 2. 

I can add a trivial in-memory optimization to reduce the number of writes. Probably not worth it.

█ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=339.99µs min=25µs    med=231µs    max=55.82ms p(90)=710µs    p(95)=985µs 
      { expected_response:true }...: avg=339.99µs min=25µs    med=231µs    max=55.82ms p(90)=710µs    p(95)=985µs 
    http_req_failed................: 0.00%   0 out of 1903417
    http_reqs......................: 1903417 34607.420736/s

## 2025-12-16 comparison

Run on a Mac M4, OS 26.2, in two VSCode terminals. Assorted apps and streaming taking place in the background.

| method | buffer | flush intvl | reqs   | evts/sec | p(90)    | p(95)  |
| --     | --     | --          | --     | --       | --       | --     | 
| PUT    | 2K     | 10          | 3.14M  | 57K      | 480us    | 630us  |
| GET    | 2K     | 10          | 3.39M  | 61K      | 488us    | 637us  |
| PUT    | 10K    | 60          | 3.43M  | 62K      | 492us    | 636us  |
| GET    | 10K    | 60          | 3.30M  | 60K      | 495us    | 640us  |
| PUT&dagger;    | 50K    | 60          | 3.15M  | 57K      | 492us    | 636us  |

&dagger; Deleted the DB before this run; was ~700MB and 16M events or so. Did not seem to make a performance difference.

There is no meaningful difference between PUT and GET methods in terms of performance.

There is no meaningful difference in performance between buffer sizes of 2K, 10K, or 50K, especially once all of the events in the buffer can be inserted in a single transaction.

We lose fewer events from smaller buffers in the event of a crash. So, smaller buffers are probably preferable.

## 2025-12-16 PUT

This test added a way to write the buffered events as a single transaction. Notice that this took the logger from 3600 events/second to 58K/second. The max used to be 2.8s, which was probably beacuse I was writing 10K events one-at-a-time, and therefore doing a full database TX for each event.

Now, all 10K (or, however many are buffered) events get written as a single TX. That's 20K int64s, or... 160K bytes. That's easy for SQLite to eat in a single transaction, apparently.

Flushing more often (every 2K or 5K events) might be better, "just in case." Be interesting to see what the perf looks like.


         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: ./put.js
        output: -

     scenarios: (100.00%) 1 scenario, 30 max VUs, 1m25s max duration (incl. graceful stop):
              * typical_usage: Up to 30 looping VUs for 55s over 5 stages (gracefulRampDown: 5s, gracefulStop: 30s)



  █ THRESHOLDS 

    http_req_duration
    ✓ 'p(95)<2' p(95)=639µs

    http_req_failed
    ✓ 'rate<0.01' rate=0.00%


  █ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=291.14µs min=25µs    med=195µs    max=59.01ms p(90)=507µs    p(95)=639µs   
      { expected_response:true }...: avg=291.14µs min=25µs    med=195µs    max=59.01ms p(90)=507µs    p(95)=639µs   
    http_req_failed................: 0.00%   0 out of 3210526
    http_reqs......................: 3210526 58372.892216/s

    EXECUTION
    iteration_duration.............: avg=309.96µs min=34.83µs med=212.87µs max=59.13ms p(90)=528.66µs p(95)=662.16µs
    iterations.....................: 3210526 58372.892216/s
    vus............................: 29      min=1            max=29
    vus_max........................: 30      min=30           max=30

    NETWORK
    data_received..................: 443 MB  8.1 MB/s
    data_sent......................: 552 MB  10 MB/s

## 2025-12-15 GET

This is a cleaner/more authentic test. It runs up to 30 VUs (virtual users), and we see a throughput of somewhere around 3600 events/second. The median is around 200us. Probably becuase I gave it a 10K buffer for events.

        /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: ./get.js
        output: -

     scenarios: (100.00%) 1 scenario, 30 max VUs, 1m15s max duration (incl. graceful stop):
              * typical_usage: Up to 30 looping VUs for 45s over 6 stages (gracefulRampDown: 3s, gracefulStop: 30s)



  █ THRESHOLDS 

    http_req_duration
    ✓ 'p(95)<2' p(95)=837µs

    http_req_failed
    ✓ 'rate<0.01' rate=0.00%


  █ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=4.36ms min=23µs    med=192µs    max=2.8s p(90)=626µs   p(95)=837µs   
      { expected_response:true }...: avg=4.36ms min=23µs    med=192µs    max=2.8s p(90)=626µs   p(95)=837µs   
    http_req_failed................: 0.00%  0 out of 170007
    http_reqs......................: 170007 3614.096923/s

    EXECUTION
    iteration_duration.............: avg=4.38ms min=31.87µs med=209.08µs max=2.8s p(90)=649.2µs p(95)=867.25µs
    iterations.....................: 170007 3614.096923/s
    vus............................: 7      min=1           max=29
    vus_max........................: 30     min=30          max=30

    NETWORK
    data_received..................: 24 MB  499 kB/s
    data_sent......................: 26 MB  553 kB/s

running (0m47.0s), 00/30 VUs, 170007 complete and 0 interrupted iterations

## 2025-12-15 PUT

        /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: ./put.js
        output: -

     scenarios: (100.00%) 1 scenario, 30 max VUs, 1m15s max duration (incl. graceful stop):
              * typical_usage: Up to 30 looping VUs for 45s over 6 stages (gracefulRampDown: 3s, gracefulStop: 30s)



  █ THRESHOLDS 

    http_req_duration
    ✓ 'p(95)<2' p(95)=852µs

    http_req_failed
    ✓ 'rate<0.01' rate=0.00%


  █ TOTAL RESULTS 

    HTTP
    http_req_duration..............: avg=4.35ms min=27µs   med=197µs   max=2.74s p(90)=629µs    p(95)=852µs   
      { expected_response:true }...: avg=4.35ms min=27µs   med=197µs   max=2.74s p(90)=629µs    p(95)=852µs   
    http_req_failed................: 0.00%  0 out of 170008
    http_reqs......................: 170008 3636.893625/s

    EXECUTION
    iteration_duration.............: avg=4.37ms min=35.2µs med=214.7µs max=2.74s p(90)=654.16µs p(95)=883.54µs
    iterations.....................: 170008 3636.893625/s
    vus............................: 8      min=1           max=30
    vus_max........................: 30     min=30          max=30

    NETWORK
    data_received..................: 24 MB  502 kB/s
    data_sent......................: 29 MB  626 kB/s
