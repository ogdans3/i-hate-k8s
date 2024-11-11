# This shit does not work yet, dont even try it

# I hate kubernetes

I hate kubernetes so im trying to make my own, where simplicity for the user is
everything.


We pride ourself on:

- Easy config files
- The default is good enough for 99% of use cases
- As few mandatory steps as possible to achieve a working webpage
- Our documentation with examples
- Everything included: dns, ingress, logging, analytics, https, (database?),
  (cache?), dashboard
- Autoscaling by default
- No fucking helm

## This is the main program

It sucks atm

## Documentation by example

Looking for a hello world example? Here you go:

```yml
project: hello-world

webserver:
    image: strm/helloworld-http
```

This is equivalent to:

```yml
project: hello-world
engine: docker

logging: true
loadbalancer: true
dashboard: true

webserver:
    image: strm/helloworld-http
    https: true
    www: true
    ports:
        - "80"
    autoscale:
        initial: 1
        autoscale: true
```

### "Thats a dumb example"

Well lets make something harder. How about

- Webpage running svelte
- Webserver running node to serve wegbpage
- Microservice 1
- Microservice 2
- Microservice 3 (What are you even doing at this point?)
- Microservice 4 (for payments) (More services than you have users)
- A valkey (new redis) cache
- A job running every night at 01:45
- A domain
- Https
- Redirect requests to www (example.com -> www.example.com)

```yml
project: i-h-k
cache: valkey

webpage:
    image: i-h-k/webpage
    domain: example.com
webserver:
    image: i-h-k/webserver
    domain:
        domain: example.com
        path: /api
microservice1:
    image: i-h-k/microservice1
microservice2:
    image: i-h-k/another1
microservice3:
    image: i-h-k/why-though
payment-microservice:
    image: i-h-k/webserver
    autoscale: false # Who knows why, but this is not allowed to scale (Are you maybe running payment without transactions?)
nightly-job:
    image: i-h-k/delete-incriminating-database
    job: true
    chron: "45 1 * * *"
```

#### How is it this simple?

Why is kubernetes so fucking hard?

Ports default to 80 in the container and ihk determines the host port.
Everything is loadbalanced by default using the default loadbalancer. www
redirect and https is assumed by default and automatically handled by ihk.

By default everything is allowed to talk to everything (within your project).

Only the webpage and webserver services are exposed to the internet via the
loadbalancer. Microservice 1-4 and the nightly job is not exposed as it does not
specify a domain, or a path. (Maybe we should add another option on the service
to allow them to be exposed? The user could hardcode a port which we should
probably expose, is 'public' a good key?).

Everything is autoscaled by default (for now atleast).
