# What was i working on again?

Sunday 17 november I was working on TLS certificates. I am creating a composite
job that adds more actions after a first pass to build image and deploy certbot
container. The second pass should wait for the container to be ready
(readinessprobe), then get certs via certbot, then something?

Currently this does not work because the readiness probe fails because of nil
pointer?
