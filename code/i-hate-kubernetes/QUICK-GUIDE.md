## Defaults

```yml
project: here-goes-your-project-name #Everything within a project can talk to eachother easily. Everything outside the project can only talk over the internet and such
cache: false #If true, we deploy a cache for you. Can also be a string if a specific cache is wanted, or an image.
logging: true #We handle logging for you if this is turned on
analytics: false #We are going to make our own analytics, think posthog with good defaults for EU privacy laws
dashboard: true #A UI dashboard where you can deploy services, nodes, stop pods, add new projects, look at cpu, mem, etc, etc
registry: false #If true, then we will deploy a docker registry for you. No reason for you to pay for a docker registry license for internal images
build-system: false #If true, then we will deploy a build system. Something like jenkins, github actions, or something similar. TBD.

i-have-a-service:
    image: i-h-k/webpage
    ports:
        hostPort: 80
        containerPort: 80
        protocol: TCP
    domain: none
    path: none
    public: false #Should be exposed to the internet, even though no domain or path is listed (You must specify a host port)
    autoscale:
        initial: 1
        autoscale: true
    www: true #If true, redirect non-www traffic to www, if false we do no redirecting, if "non-www" then we redirect www to non-www.
    https: true #Redirect all http to https, and get me a certificate. (We will also add some options for certificates in the future)
    probes: #We need some probes
i-have-a-job: #In general the same defaults as any service, but has a few extra options listed below.
    image: some-docker-image
    job: true #This is a job
    chron: "45 1 * * *"
```

# Run program

```
go run main.go deploy
go run main.go deploy -f examples/hello-world-dns-routing.yml
gow run main.go deploy -f examples/hello-world-dns-routing.yml
go run main.go stop
```

# Run tests

```
go test -p 1 --parallel 1 ./test/e2e_test/
```

# Test webhook

### Run webhook short example for github pushes
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "ref": "refs/heads/master",
  "repository": {
    "html_url": "https://github.com/user/example-repo.git",
    "git_url": "git://github.com/user/example-repo.git",
    "ssh_url": "git@github.com:user/example-repo.git",
    "clone_url": "https://github.com/user/example-repo.git"
  }
}' http://localhost:6444/webhook/github
```

### Run webhook short example for github pushes to ihk
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "ref": "refs/heads/master",
  "repository": {
    "html_url": "https://github.com/ogdans3/i-hate-k8s.git",
    "git_url": "git://github.com/ogdans3/i-hate-k8s.git",
    "ssh_url": "git@github.com:ogdans3/i-hate-k8s.git",
    "clone_url": "https://github.com/ogdans3/i-hate-k8s.git"
  }
}' http://localhost:6444/webhook/github
```
