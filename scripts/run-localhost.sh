#!/bin/sh

openssl req -nodes -newkey rsa:4096 \
  -keyout ./deployments/localhost/registry-auth.key \
  -out ./deployments/localhost/registry-auth.csr \
  -subj "/CN=http:\/\/localhost:5100"

openssl x509 \
  -in ./deployments/localhost/registry-auth.csr \
  -out ./deployments/localhost/registry-auth.crt \
  -req -signkey ./deployments/localhost/registry-auth.key -days 3650

docker-compose \
    --file ./deployments/localhost/docker-compose.yml \
    --project-name harbour \
    up \
    --force-recreate \
    --detach
