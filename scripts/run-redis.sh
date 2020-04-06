#!/bin/sh

docker-compose \
    --file ./deployments/redis/docker-compose.yml \
    --project-name harbour \
    up \
    --force-recreate \
    --detach
