#!/bin/bash

INSTANCE="standalone-redis"

if [ "$1" != "" ]; then
  INSTANCE=$1
fi

echo "Instance: $INSTANCE"

docker run -it --network "harbour_$INSTANCE" --rm redis:5-alpine redis-cli -h "$INSTANCE"