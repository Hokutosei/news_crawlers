#!/bin/bash

GOOS=linux GOARCH=amd64 go build -v -o linux_news_crawlers

# deploy
echo "--->> building container...."
docker build -t jeanepaul/news_crawlers .

echo "--->> pushing container"
docker push jeanepaul/news_crawlers:latest

echo "--->> stoping newscrawlers pod"
kubectl stop pod newscrawlers

echo "--->> creating newscrawlers pod"
kubectl create -f "$(pwd)"/kubernets_pod.yaml

echo "done!"
# run container
# docker run -d -e "COREOS_PRIVATE_IPV4=" jeanepaul/news_crawlers
