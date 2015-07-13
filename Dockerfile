FROM ubuntu:wily

# compile to linux
# GOOS=linux GOARCH=amd64 go build -v -o linux_news_crawlers

# deploy
# docker build -t jeanepaul/news_crawlers .
# docker push jeanepaul/news_crawlers

# run container
# docker run -d -e "COREOS_PRIVATE_IPV4=" jeanepaul/news_crawlers

MAINTAINER jeanepaul@gmail.com

# RUN apt-get update --fix-missing
RUN apt-get install -y ca-certificates

COPY linux_news_crawlers /usr/bin/

WORKDIR /usr/bin

ENTRYPOINT ./linux_news_crawlers
