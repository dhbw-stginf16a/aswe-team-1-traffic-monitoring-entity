#!/bin/bash

docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
docker build -t asweteam1/traffic-monitor:latest -t asweteam1/traffic-monitor:$TRAVIS_TAG --label version="$TRAVIS_TAG" --label go_version="$TRAVIS_GO_VERSION" .
docker push asweteam1/traffic-monitor:latest
docker push asweteam1/traffic-monitor:$TRAVIS_TAG
