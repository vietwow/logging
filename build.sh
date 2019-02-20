#!/bin/bash

docker login && docker build -t vietwow/logging . && docker push vietwow/logging
