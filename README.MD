# LogCrunch - Logs from Everywhere, Anywhere.

LogCrunch is a proto-SIEM that aggregates logs by parsing user-defined targets and schemas at runtime to ingest logs. These are then sent to a centralized stash where they can be queried and investigated. LogCrunch's primary goal is to undercut the expected computational requirements for a SIEM to minimalize resource usage, focusing on Collegiate Cyber Defense Competitions and homelabs. 

LogCrunch is built for UNIX systems, but support for Windows is planned for future releases.

## Notes
LogCrunch is a student project for my Honor's Undergrad thesis. Issues and PRs are welcome, but please ask before attempting any major contributions, as they may need to wait until after my thesis defense. Feel free to reach out via GitHub issues or email me directly at tlop \[at] tlop \[dot] xyz. RFC 1149-compliant addresses available upon request.

## Installation
### Server (Must be Done First)
1. Download and extract the latest [release](https://github.com/TLop503/LogCrunch/releases) for your architecture.
2. Generate or otehrwise configure TLS certificates. Note the path of the `key` and `crt`.
3. Choose if you want to run the server as `root` or under a dedicated user
   1. Users will need access to read + write to `/var/log/LogCrunch/*` and `/opt/LogCrunch/`. 
5. Start the intake server, specifying each argument. Running the intake w/o arguments will list the up-to-date required parameters. You'll need to pass in the paths to your cryptography, where to listen for logs, and where to serve the WebUI. 
   1. Note you may need to update your firewall to allow LogCrunch to receive data from the network. 
6. On the first startup, a (randomized) default password for the WebUI is printed to `stdout`. Upon first sign in to the WebUI, you will need to update this to something unique before you can access the query interface. There are no hardcoded password requirements, use your best judgement for your environment.
   7. If you are locked out or otherwise unable to log in, you'll need to delete the `/opt/LogCrunch/users` directory and restart the server.
### Agent(s)
1. Download and extract the latest [release](https://github.com/TLop503/LogCrunch/releases) for your architecture.
2. Write a `targets.yaml` configuration file, specifying logs to read and how to parse them. An example can be found in [Github](https://github.com/TLop503/LogCrunch/blob/main/agent/custom_cfg.yaml). Some parsing schemas are built-in as well, these can be found in the [MetaParser Registry](https://github.com/TLop503/LogCrunch/blob/ec750e335f0ab9f7d891e14750a5a5c1ea80b563/structs/meta_parser_regisry.go#L9). 
3. Start the agent, specifying arguments for the IP of the SIEM server, intake port, config path, and whether to verify the TLS certificates. 

### Automated Server Deployment, Dockerfiles, Etc.
Scripted installation methods are hosted in the [utility repo](https://github.com/TLop503/LogCrunch-Utils).

## Building from Source
1. Install GoLang
2. Clone the repository
3. `cd LogCrunch`
4. Compile with: `go build -o LogCrunch-Server ./server; go build -o LogCrunch-agent ./agent`