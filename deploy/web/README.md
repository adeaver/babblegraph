# Setting Up a Web Server

## Setup UFW

Install UFW (uncomplicated firewall)
`sudo apt install ufw` <br />

Setup allowable traffic
`sudo ufw allow ssh`<br />
`sudo ufw allow proto tcp from any to any port 80,443`<br />

Configure `/etc/default/ufw`<br />
- Set `DEFAULT_FORWARD_POLICY="ACCEPT"`
- Set `IPV6=no`

Enable the firewall: `sudo ufw enable`<br />
Ensure that it's running with `sudo ufw status verbose`

## Setup Docker

Allow docker to use the host's nameserver
`sudo ln -sf /run/systemd/resolve/resolv.conf /etc/resolv.conf`<br />

Optionally, volume mount /etc/resolv.conf on the host to /etc/resolv.conf on the container<br />

Set up docker daemon to respect the firewall by adding the following line to `/etc/docker/daemon.json`<br />
`{"iptables": false}`
