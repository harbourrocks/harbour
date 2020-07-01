#!/bin/sh

openssl req -nodes -newkey rsa:4096 \
  -keyout ./registry-auth.key \
  -out ./registry-auth.csr \
  -subj "/CN=https:\/\/harbour.rocks"

openssl x509 \
  -in ./registry-auth.csr \
  -out ./registry-auth.crt \
  -req -signkey ./registry-auth.key -days 3650
