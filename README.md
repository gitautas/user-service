# User Service

An example for a user microservice written in golang.

## Description

An in-depth paragraph about your project and overview of use.

## Getting Started
* Clone the repository
```
git clone https://github.com/gitautas/user-service.git
```
* Run tests

```
make test
```

`
* Build and run the project
```
make start
```
which runs

```
docker-compose up
```

## Choices taken, assumptions made
I decided to use MongoDB, as after implementing MySQL first, I realised that this would be a perfect opportunity to show off my ability to learn quickly, this service being my first experience with a NoSQL database.
I decided to use Redis as my pubsub, as it is simple and does the job.
I decided to use Gin as my HTTP framework, mainly because I wasn't knowledgeable enough about gRPC, Gin is a great framework for REST, but when I also have gRPC endpoints my controller logic ended up split between the two. The gRPC and REST APIs are different because of this.


## Improvements to be made
If I had unlimited time I would re-do the entire project structure, fix the naming conventions, fix the error model to one that's not tied to any specific technology, implement a proper logger, improve on the asynchronicity of the service, handle my protobuf dependecies sanely, use a gateway instead of a separate REST API (this is the biggest one).
Essentially, if I had more time I would use the lessons learned to rewrite this service, as this was my first foray into gRPC and MongoDB, learned quite a bit and can see all of the inconsistencies and bad choices made. However since I only had 4.5 days of free time to do this, this was my best attempt.

## Additional thoughts
This was, if anything, a learning experience. I had a week to do this in theory, but I only had 4 days free to spend on this, not to mention getting sick in the middle of them. I learned a lot about Mongo and gRPC, which are really cool technologies, as especially with gRPC, very different from anything I've used before, so as I was working on this, I kept learning more and more, realising that the choices I made before were not the best.
