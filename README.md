# effective-invention

## JS 

Sends batches of 10 requests at a time with a brief wait to stay under a 1000 request per minut rate limit.

Result: fetched 69726 tasks in 268.69 seconds. RMP: 156.31

## Go

Sends a new request whenever a token is availible, with the rate limited by 1 token added every 60 milliseconds to stay under the 1000 request per minute rate limit. Because each request is in its own go routine it doesn't wait for other requests to finish. Responses are sent through a buffered channel into a processor that handles combining the tasks in a separate go routine. 

Result: fetched 69726 tasks in 44.543s seconds. RPM: 939.23

