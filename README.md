# gaze
Stare at a logfile (or stdin) and send each line to [riemann](https://riemann.io/) wrapped in an event.

### Features:
```
 -- Pros
     -- Fast. About 10x-50x Logstash performance.
     -- Light. ~40M with 100k queue size. Negligible CPU under normal event load.
     -- Never loses data because of a log rotation.
     -- Sheds load on full queue, never fills tcp buffer or exceeds socket limit.
     -- Supports multiple reimann endpoints for simultanous transmission.
     -- Supports per-endpoint worker pool with connection reuse.
     -- Easy to switch to TLS if needed.
     -- Simple. 400 lines of go.

 -- Cons
     -- No dynamic parsing of incoming data stream
         -- Just builds a riemann event with the full log line stored as a string in the event.
```

### Usage:
```bash
# Tail a logfile and send to riemann with :host myhostname, :service myservice, :ttl 10
./gaze --hostname myhostname --service myservice --ttl 10 --riemann 127.0.0.1:5555 /path/to/file.log

# As above, but read from stdin instead
./gaze --hostname myhostname --service myservice --ttl 10 --riemann 127.0.0.1:5555 -

# Send all events to three different riemann servers simultaneously. Determine hostname automatically.
./gaze --service myservice --riemann 127.0.0.1:5555 --riemann 127.0.0.2:5555 --riemann 127.0.0.3:5555 /path/to/file.log

```
### Configuration:
```go
// Configure some compile-time options (in common/common.go)
var RIEMANN_CONNECT_TIMEOUT_MILLIS uint = 50
var RIEMANN_RETRY_DELAY_MILLIS uint = 200
var RIEMANN_RETRY_ATTEMPTS uint = 100
var QUEUE_LENGTH uint = 100000
var POOL_SIZE uint = 1
```
