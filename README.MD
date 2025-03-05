# Heartbeat0 is now LogCrunch
*A WIP SIEM for CDC Competitions. Name subject to change again :)*

This SIEM (though currently only a log aggregator) works by `tail -f`'ing an (easily editable!) list of log files and sending them over the wire to a single stash. This project aims to undercut existing SIEMs in computational strength needed to run on even the weakest of endpoints. At this stage in development the prebuilt binaries will run on most modern linux with minimal resources. Support for other OSs are in development.

## What works
As of 2/9/25, logs are able to be read from anywhere (config-based) and sent to the server. Additionally, a heartbeat log (which led to the repo name) is sent every few minutes to serve as a proof of life update. This is an extremely early version, so useful features like a UI, query interface, or event detection hasn't been implemented yet.

## What doesn't work (Or, planned features)
- Parsing syslogs into unified format (see eof)
- Querying logs
- At-a-glance ui of which boxes' agents are alive
- Event detection
- lolbin detection
- proxys/tunnels

## Other notes
This is a WIP project for my honors undergrad thesis. Issues, etc are welcome but since this is tied to my graduation status I may not be able to accept collaboration/PRs at this point in time.

# How to run
There are 2 ways to run this project. You can either download the latest [release](https://github.com/TLop503/LogCrunch/releases), or clone the project.

## Setup
1. Download the tarball and extract it, or clone the repo
   1. If you are building from source you'll need golang installed
1. If desired, edit the config file to add or remove log locations you'd like to watch.
1. Generate some TLS certs!
   1. `openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"`
1. Start the intake server, specifying:
   1. Host (0.0.0.0 to recieve from the network)
   1. Port
      1. Note! you will need to open this port on your firewall if applicable
   1. The path to the .crt
   1. The path to the .key
1. Start the agent, specifying:
   1. Host (IP of the box with the server)
   1. Port
   1. Whether or not to verify the TLS certs being used (y/n)
1. Do whatever you like. Logs will appear in ./logs on the box running the intake server


---


Log format (parsing pending, this is an early example):
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