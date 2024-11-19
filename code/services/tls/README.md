# TLS agent for hive

This service is responsible for all TLS certificates for hive. It should be run
as a job in the cluster periodically (twice a day?)

Will by default run certbot to fetch and re-fetch certificates. It will expose
paths to `/ready`, `/live` and the acme challenges.

Only HTTP_01 is supported.

The certificates will be exported to /var/www/certbot? The cluster is
responsible for storing these certificates in volumes and exposing them to the
loadbalancer.

This service starts an nginx instance and waits for more command before it does
anything else. This is because the cluster needs time to enter the containers ip
into the loadbalancing.

When the cluster is ready we expect the cluster to send a command to the
container (e.g. via docker exec) with the following pattern. Not that all
certbot commands are available.

```bash
certbot certonly --webroot -w /var/www/certbot -d yourdomain.com --agree-tos --email your-email@example.com --non-interactive --dry-run
```
