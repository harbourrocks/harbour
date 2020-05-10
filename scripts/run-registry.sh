#!/bin/sh

openssl req -nodes -newkey rsa:4096 \
  -keyout ./deployments/registry/registry-auth.key \
  -out ./deployments/registry/registry-auth.csr \
  -subj "/CN=http:\/\/localhost:5100"

openssl x509 \
  -in ./deployments/registry/registry-auth.csr \
  -out ./deployments/registry/registry-auth.crt \
  -req -signkey ./deployments/registry/registry-auth.key -days 3650

docker-compose \
    --file ./deployments/registry/docker-compose.yml \
    --project-name harbour \
    up \
    --force-recreate \
    --detach
