#!/bin/sh

docker-compose \
    --file ./deployments/registry/docker-compose.yml \
    --project-name harbour \
    up \
    --force-recreate \
    --detach
