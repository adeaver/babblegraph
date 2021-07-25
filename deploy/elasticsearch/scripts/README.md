# ElasticSearch Scripts

## Backup

The first step of configuring elasticsearch snapshots is to setup the snapshot repository by running:
```
deploy/elasticsearch/scripts/setup-snapshot-repostiory
```

This will configure ElasticSearch to use DigitalOcean spaces as an S3-compatible snapshot repostiory. Next, you'll need to restart the node and run:
```
deploy/elasticsearch/scripts/setup-snapshot-configuration
```

to actually configure the snapshots.

Lastly, edit the crontab to call:
```
deploy/elasticsearch/scripts/remote-capture-snapshot
```
