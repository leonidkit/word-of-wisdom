Design and implement “Word of Wisdom” tcp server.
* TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
* The choice of the POW algorithm should be explained.
* After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
* Docker file should be provided both for the server and for the client that solves the POW challenge

## Requirements
* go 1.21 and upper
* [Task manager](https://taskfile.dev/installation/)
* Docker

## Why hashcash?
* Easy to understand.
* Uses efficient algorithms of the SHA family.
* Has [documentation](http://hashcash.org/).

## How to start?

```
$ docker build -t wow-server -f Dockerfile.server
$ docker build -t wow-client -f Dockerfile.client
$ docker run --network host -d wow-server
$ docker run --network host wow-client
```

## Want to do
* Review project structure.
* Review config approach.
* Review handlers approach.
* Server, handlers tests.