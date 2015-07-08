#!/bin/bash

GOOS=linux GOARCH=amd64 go build -v -o linux_news_crawlers

# deploy
docker build -t jeanepaul/news_crawlers .
docker push jeanepaul/news_crawlers

# run container
# docker run -d -e "COREOS_PRIVATE_IPV4=" jeanepaul/news_crawlers
