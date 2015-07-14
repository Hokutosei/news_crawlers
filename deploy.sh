#!/bin/bash

GOOS=linux GOARCH=amd64 go build -v -o linux_news_crawlers

# deploy
echo "--->> building container...."
docker build -t jeanepaul/news_crawlers .

# echo "--->> pushing container"
# docker push jeanepaul/news_crawlers:latest

echo "--->> re-tag container..."
docker tag -f jeanepaul/news_crawlers gcr.io/chat-app-proto01/news_crawlers

echo "--->> pushing container"
gcloud docker push gcr.io/chat-app-proto01/news_crawlers

echo "--->> stoping newscrawlers pod"
kubectl stop pod newscrawlers

echo "--->> creating newscrawlers pod"
kubectl create -f "$(pwd)"/kubernets_pod.yaml

echo "done! ctrl+c to stop status!"
kubectl logs -f newscrawlers
while true; do kubectl get pods; sleep 5; done
# run container
# docker run -d -e "COREOS_PRIVATE_IPV4=" jeanepaul/news_crawlers
