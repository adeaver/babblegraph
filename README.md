# Babblegraph

The main repository for Babblegraph

## To Run System Locally

Get AWS SES access key and token, set them in your local environment as `$BABBLEGRAPH_SES_ACCESS_KEY` and `BABBLEGRAPH_SES_ACCESS_SECRET` respectively

Start local databases:
`./scripts/run-local-dbs`

Start local services:
`./scripts/start-services`

## To Run Vim

Start local shell with configured GOPATH:
`./scripts/shell`

## Structure of this repository
- `backend`: all of the go code to run all babblegraph services.
    - `actions`: actions are any functionality based on models. This helps reduce cyclical dependencies.
    - `model`: all database models
    - `util`: all utility packages
    - `wordsmith`: special code for interacting with wordsmith database
    - `services`: all runnable go services.
    - `jobs`: larger extraneous tasks that are runnable from multiple services
- `deploy`: configurations for all deployable services. each service gets its own directory
- `ops`: configurations for local environment
- `scripts`: all useful bash scripts for working with babblegraph

## Other helpful commands

- Restarting the worker process locally
`docker restart babblegraph_worker_1`

- Stopping local databases
`docker-compose -f ops/local-dbs.compose.yaml down`

- Apply migrations to local databases
`docker exec -it ops_db_1 /home/postgres/scripts/apply-migrations --db-user dev --db-name babblegraph`

- Insert user into local database
`docker exec -it ops_db_1 /home/postgres/scripts/insert-user --db-name babblegraph --db-user dev --email`<br />
`docker exec -it ops_db_1 psql -U dev -d babblegraph -c "SELECT _id FROM users WHERE email_address='{email}'"`<br />
`docker exec -it ops_db_1 /home/postgres/scripts/add-user-reading-level --db-name babblegraph --db-user dev --lang es --level 4 --user {userID}`<br />
