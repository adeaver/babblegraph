# ElasticSearch Scripts

## Server Setup

- Move `limits.conf` to `/etc/security/limits.conf` on the host machine.
- Move `sysctl.conf` to `/etc/sysctl.conf` on the host machine and reboot using `sudo shutdown -r now`.
- Start ElasticSearch

## Passwords
Set password using
```
deploy/elasticsearch/scripts/setup-passwords
```

## Backup
Used this guide: https://sanacl.wordpress.com/2020/03/30/elasticsearch-snapshots-in-digitalocean-spaces/

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
