# Heartbeat0 is now LogCrunch
*A WIP SIEM for CDC Competitions. Name subject to change again :)*

This SIEM (though currently only a log aggregator) works by reading and parsing a user-defined list of log files and sending them over the wire to a single stash. This project aims to undercut existing SIEMs in computational strength needed to run on even the weakest of endpoints. At this stage in development the prebuilt binaries will run on most modern linux with minimal resources. Support for other OSs are in development.

## Notes
This is a WIP project for my honors undergrad thesis. Issues, etc. are welcome but since this is tied to my graduation status I may not be able to accept collaboration/PRs at this point in time.

# How to run
There are a few ways to run this project. You can download the latest [release](https://github.com/TLop503/LogCrunch/releases) or clone the project to manually set up the server and agents. Or, you can use the setup [script](https://github.com/TLop503/LogCrunch/blob/main/scripts/server_setup.py) to automatically install the server, but will still need to manually deploy agents. Agent deployment scripts are coming soon!

## Setup
1. Download the tarball and extract it, or clone the repo
   1. If you are building from source you'll need golang installed. Compile w/ `go build ./server; go build ./agent`
1. If desired, edit the target file to add or remove log locations you'd like to watch.
1. Generate some TLS certs!
   1. `openssl req -x509 -newkey rsa:4096 -keyout path_to_server.key -out path_to_server.crt -days 365 -nodes -subj "/CN=localhost"`
1. Start the intake server as root, specifying:
   1. Host (0.0.0.0 to recieve from the network)
   1. Port
      1. Note! you will need to open this port on your firewall if you intend to recieve logs over the network
   1. The path to the .crt
   1. The path to the .key
1. Start the agent, specifying:
   1. Host (IP of the box with the server)
   1. Port
   1. The path to a target.yaml (see agent/target.yaml for example)
   1. Whether to verify the TLS certs being used (y/n)
1. Logs will appear in `/var/log/LogCrunch/firehose.log` on the box running the intake server. The webui will appear at :8080.
