# Nomad Toast

[![Build Status](https://travis-ci.org/jrasell/nomad-toast.svg?branch=master)](https://travis-ci.org/jrasell/nomad-toast) [![Go Report Card](https://goreportcard.com/badge/github.com/jrasell/nomad-toast)](https://goreportcard.com/report/github.com/jrasell/nomad-toast)

nomad-toast in an open source tool for receiving notifications based on [HashiCorp Nomad](https://www.nomadproject.io/) events. It is designed to increase observability throughout an organisation and provide insights within chatops style environments. 

### Supported Nomad Endpoints

nomad-toast currently supports watching the following Nomad API endpoints:

* Allocations - https://www.nomadproject.io/api/allocations.html
* Deployments - https://www.nomadproject.io/api/deployments.html

### Supported Notification Endpoints

* Slack - https://slack.com/

## Download & Install

* The nomad-toast binary can be downloaded from the [GitHub releases page](https://github.com/jrasell/nomad-toast/releases) using `curl -L https://github.com/jrasell/nomad-toast/releases/download/0.0.1/nomad-toast_linux_amd64 -o nomad-toast`

* A docker image can be found on [Docker Hub](https://hub.docker.com/r/jrasell/nomad-toast/), the latest version can be downloaded using `docker pull jrasell/nomad-toast`.

* nomad-toast can be built from source by firstly cloning the repository `git clone github.com/jrasell/nomad-toast.git`. Once cloned the binary can be built using the `make` command.

## Commands and Flags

nomad-toast supports the following global flags:

* **--log-format** (string: "AUTO") Specify the format of nomad-toast logs. Valid values are AUTO, ZEROLOG or HUMAN.
* **--log-level** (string: "INFO") The level at which nomad-toast will log. Valid values are DEBUG, INFO, WARNING, ERROR and FATAL.
* **--log-use-color** (bool: true) Use ANSI colors in logging output.
* **--nomad-address** (string: "http://localhost:4646") The HTTP(S) API endpoint for Nomad where all calls will be made.
* **--nomad-allow-stale** (bool: true) Allow stale Nomad consistency when making API calls.
* **--slack-auth-token** (string: "") The Slack API auth token for connectivity to slack.
* **--slack-channel** (string: "") The Slack channel to send Nomad notifications to.

#### Command: `allocations`

Allocations triggers a watcher on the Nomad allocations endpoint and will notify you when allocations go through a state change. This can be helpful when keeping an eye on failed allocations, or just to have a general insight into allocation churn.

The allocations command supports the following flags:

* **--include-states** (comma separated list of strings: "") Comma-separated list of allocation client states that will be whitelisted for notifications. If specified, *only* these states will be included in notifications.
* **--exclude-states** (comma separated list of strings: "") Comma-separated list of allocation client states that will be excluded from notifications. This takes priority over include-states.

The set of client states is defined by the Nomad API - see the `AllocClientStatus` constants [in the API docs](https://godoc.org/github.com/hashicorp/nomad/api#pkg-constants).

#### Command: `deployments`

Deployments triggers a watcher on the Nomad deployments endpoint. This allows you to get notified of deployment activities on your cluster and allows stakeholders to gain insight.
