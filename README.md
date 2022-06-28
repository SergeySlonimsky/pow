### TODOs
+ Refactor client architecture
+ Refactor protocol to decrease coupling between server impl and protocol impl
+ Add docker scratch image instead of alpine
+ Add conf struct instead of os.Getenv()

### Requirements

+ Go 1.18+
+ Docker, docker-compose

### Makefile commands

+ `make build-client` - builds client application
+ `make build-server` - builds server application \n
+ `make format` - formats source code with gofumpt
+ `make lint` - lints source code with golang-ci\n
+ `make test` - tests client and server code\n

### Run with Docker Compose

+ `docker-compose up`

### Proof of Work defenition

Proof of work (PoW) is a form of cryptographic proof in which one party (the prover) proves to others (the verifiers) 
that a certain amount of a specific computational effort has been expended. 
Verifiers can subsequently confirm this expenditure with minimal effort on their part.
A key feature of proof-of-work schemes is their asymmetry: the work – the computation – must be moderately hard (yet feasible)
on the prover or requester side but easy to check for the verifier or service provider.

### Definition of hashcah algorithm 

There are a few possible algorithms to implement PoW solution.
+ Hashcash
+ Merkle tree
+ Guided tour puzzle

I chose hashcash because this algorithm is easy to implement, it does not require much computation on the server side,
and it allows you to flexibly adjust the complexity for the client.
