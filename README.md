1. Introduction

NomadX is a Slackbot written in Golang to monitor cryptocurrency prices and alert when necessary.
It uses the coinmarketcap API.

1. Installation

* `go get https://github.com/snufkin/nomadx`
* Create an app-env file with your bot's credentials (`NOMADX_BOTTOKEN`, `NOMADX_CHANNELID`, `NOMADX_BOTID`)
* Run `sh build.sh`

This should produce a binary that you can use for local testing. Deployment and Docker config is pending.

1. Testing

* `go test`