# cf-redis-benchmark

Logs performance of the cf-redis-release. Experimental.

### Intentions:

#### Measure some Service Level Indicators (SLIs)
1. How long does it take to provision a Redis service instance?
2. How long does it take to write and read from a long running service instance?

Each indicator will be tested from time t time based on the SLI interval.
Each indicator will have a timeout after which it is considered to have failed.

#### Output the SLIs to an appropriate location
The benchmark will run as a Cloud Foundry App and emit the SLI measurements to an endpoint such as DataDog.

These SLI measurements should probably contain:

- The timestamp that the measurement was initiated
- If it succeeded, the time it took to succeed
- If it failed fast, the time it took to fail
- If it had not finished before the timeout, a `timeout` message should be emitted, this will usually be considered a failure

Therefore a single attempt to measure an SLI may emit up to two messages. One for success or failure, and sometimes one for `timeout`. These two messages will always have matching timestamps.

#### Provide advice on SLI interpretation
There are two distinct ways that an SLI can be calculated:
1. Indicator duration - how long does the measurement tend to take? Such as the mean time it takes to provision an 
service instance, or the median time taken to complete a write/read from Redis. This needs to take into account fast 
failures and timeouts.
2. Indicator success - what proportion of the time does a measurement succeed? i.e. The percentage of the time when 
Redis is able to read/write. This depends on the timeout and measurement interval granularity.
