FROM ubuntu:20.04

RUN apt-get update \
  && apt-get install --no-install-recommends -y docker

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
