# Queue Package

Implementation of a message queue using Postgres

## TODO
- [ ] Support for maximum retries on a per queue basis
- [ ] Ability to put failed queued messages in the back of the line
- [ ] Fix assumption that postgres table will always be `queue_messages`
