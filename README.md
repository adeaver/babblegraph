# Babblegraph

The main repository for Babblegraph

## To Run System Locally

Start local databases:
`./scripts/run-local-dbs`

Start local services:
`./scripts/start-services`

## To Run Vim

Start local shell with configured GOPATH:
`./scripts/shell`

## Structure of this repository
- `backend`: all of the go code to run all babblegraph services.
    - `model`: all database models
    - `util`: all utility packages
    - `wordsmith`: special code for interacting with wordsmith database
    - `services`: all runnable go services.
- `deploy`: configurations for all deployable services. each service gets its own directory
- `ops`: configurations for local environment
- `scripts`: all useful bash scripts for working with babblegraph
