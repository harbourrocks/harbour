# Harbour.rocks
Harbour is a docker registry including build system focused on cloud native environments.

![Screenshot of Harbour](images/screenshot.png)

# Architecture

harbour.rocks is composed of several **Golang microservices**. It utilizes **Redis** as database and supports various **OpenId Connect** authentication provider.

![Architecture of harbour.rocks](images/architecture.png)

# Development

## How to develop harbour

* You have to install Go to your system for all the backend work: https://golang.org/doc/install
  
  *Right now there are no images build for the individual services but that is on the roadmap.*

* Install docker AND docker-compose to your system. Make sure both work by executing the following command line (note: versions do not have to match exactly)
```
$ docker -version
Docker version 19.03.8, build afacb8b

$ docker-compose -version
docker-compose version 1.25.4, build 8d51620a
```

* You need a redis database for storage. Simply setup the redis by executing this script: [scripts/run-redis.sh](scripts/run-redis.sh), or use your existing local one.

  *Script only works with docker installed before.*

* The actual registry that stores the images is the official docker registry. Simply setup the registry by executing this script: [scripts/run-registry.sh](scripts/run-registry.sh). This will also create all certificates that are required for authenticating with the registry.

  *Script only works with docker installed before.*
  
  *You will also have to adapt the REGISTRY_AUTH_TOKEN_REALM env variable which is the address of the IAM service on your PC. This has to be the ip address of your network adapter, because the registry (running inside a container) has to connect with it and localhost won't work from inside a container. Change the env variable in this file: [deployments/registry/docker-compose.yml](deployments/registry/docker-compose.yml)*

* We are using IntelliJ 2020 for development. You should open the project root as IntelliJ project. This will give you several run configurations which should work for the development environment.

  FOR THE TIME THIS PROJECT IS BEING DEVELOPED BY WEB SYSTEMS STUDENTS TRY TO WORK WITH THE PROVIDED CONFIGURATION.

* Now you should be able to run the several services of harbour from your IntelliJ by selecting the configuration and hitting the run button.

For further information see [CONTRIBUTING.md](CONTRIBUTING.md)

## <img style="float: left; margin-right: 10px;" src="https://redis.io/images/redis-small.png"> Redis

**Version:** 5+

Harbour stores its data in the key-value database [Redis](https://redis.io/).

To run a development redis database you can use the docker-compose file provided here [deployments/redis/docker-compose.yml](deployments/redis/docker-compose.yml).

There is also a script to quickly run a redis instance [scripts/run-redis.sh](scripts/run-redis.sh).

`./scripts/run-redis.sh`

You can access redis for debugging purposes by the redis-cli using the following script [scripts/redis-cli.sh](scripts/redis-cli.sh).

`./scripts/redis-cli.sh`

# Docker Registry API

* List all repositories
  * /v2/_catalog
* List all tags (images inside a repository)
  * /v2/\<name>/tags/list
  
https://docs.docker.com/registry/spec/api/#detail

# Licenses

## Font Awesome

The project, application, website including our harbour.rocks organizational logo uses Font Awesome licensed under https://fontawesome.com/license. There were no modifications made. Thank you Font Awesome for providing a wide variety of free icons!
