#!/bin/sh

echo starting containers...

docker-compose pull
docker-compose -p 5letters up -d --remove-orphans
