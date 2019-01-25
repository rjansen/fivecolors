# raizel [![Build Status](https://travis-ci.org/rjansen/fivecolors.svg?branch=master)](https://travis-ci.org/rjansen/fivecolors) [![Coverage Status](https://codecov.io/gh/rjansen/fivecolors/branch/master/graph/badge.svg)](https://codecov.io/gh/rjansen/fivecolors) [![Go Report Card](https://goreportcard.com/badge/github.com/rjansen/fivecolors)](https://goreportcard.com/report/github.com/rjansen/fivecolors)

A library to manipulate and manage card games collections. You can query cards through a GraphQL api.
The server can be deployd as a function or as a classic http server.
This project has deploy configurations for Cloud Functions, docker containers or kubernetes.
The persistence layer can be switched between PostgreSQL or Cloud Firestore.

##### Cassandra and MySQL support will be avaiable soon

# dependencies
### tools (you must provide the installation)
- [Docker](https://www.docker.com/)

### libraries
- [raizel](https://godoc.org/cloud.google.com/go/firestore)
- [l](https://github.com/rjansen/l)
- [migi](https://github.com/rjansen/migi)
- [yggdrasil](https://github.com/rjansen/yggdrasil)
- [graphql-go](https://github.com/rjansen/yggdrasil)

# tests and coverage
- run all ci tasks: `make ci`
- run unit tests: `make test`
- run integration tests: `make itest`
- run coverage: `make coverage.text`
- run html coverage: `make coverage.html`

# run
- `MODULE_NAME=server make run`

# on docker
- build container: `make docker.build`
- container bash: `make docker.bash`
- follow [run](#run) instrunctions
- or follow [test and coverage](#test_and_coverage) instructions
