# Babblegraph Production Playbook

## Deploying Services

Currently, the only deployable service is the worker. This is deployed with:
`./scripts/deploy-worker`

## Adding new users

Users currently have to be added manually. This is, unfortunately, a combination of 3 commands.

Starting on `prod-db-1`, run:
- `docker exec -it babblegraph_main_db /home/postgres/scripts/insert-user --email {email_address}`
This will insert the user into the database. You next need to figure out what their assigned user ID is with
- `docker exec -it babblegraph_main_db psql -U bgmainuser -c "SELECT _id FROM users WHERE email_address='{email_address}'"`
This will give you the ID of the user. You will be prompted for the postgres password, which you can get from ./env. Lastly, you need to set their reading level or email sending will fail:
`docker exec -it babblegraph_main_db /home/postgres/scripts/add-user-reading-level --level {level} --lang {language_code} --user {_id}`

## Manually sending emails

If there's been an error sending emails, you can use the task runner to manually send the day's emails.
`./scripts/run-task-daily-email`

At present, there is no mechanism to send to a subset of users in the case that there is a panic on a single user.
