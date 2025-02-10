# Heartbeat0
*A WIP SIEM for CDC Competitions. Name subject to change*

This SIEM (though currently only a log aggregator) works by `tail -f`'ing a (easily editable!) list of log files and sending them over the wire to a single stash. Because of Go's channels and routines, writing is both thread safe and elegant. This project aims to undercut existing SIEMs in computational strength needed to run on even the weakest of endpoints.

## What works
As of 2/9/25, logs are able to be read from anywhere and sent to the server. Additionally, a heartbeat log (which led to the repo name) is sent every few minutes to serve as a proof of life update. This is an extremely early version, so useful features like a UI, query interface, or event detection hasn't been implemented yet.

## What doesn't work (Or, planned features)
- Parsing syslogs into unified format (see eof)
- Querying logs
- At-a-glance ui of which boxes' agents are alive
- Event detection
- lolbin detection

## Other notes
This is a WIP project for my honors undergrad thesis. Issues, etc are welcome but since this is tied to my graduation status I may not be able to accept collaboration/PRs at this point in time.

---


Log format:
```json
{
   "host":"placeholder_host",
   "timestamp":1734649442,
   "type":"{HB : Seq_Err}",
   "payload":{
      "exp_seq":0,
      "recv_seq":9
   }
}
{
   "host":"placeholder_host",
   "timestamp":1734649447,
   "type":"{HB : Good}",
   "payload":{
      "type":"proof_of_life",
      "timestamp":1734649447,
      "seq":10
   }
}
```