FROM docker:19.03

# Pinned versions
ARG BASH_VERSION='5.0.17-r0'
ARG OPENSSL_VERSION='1.1.1g-r0'
ARG OPENSSH_VERSION='8.3_p1-r0'

RUN apk add --no-cache \
  bash=${BASH_VERSION} \
  openssl=${OPENSSL_VERSION} \
  openssh=${OPENSSH_VERSION}

COPY deploy.sh /deploy.sh

ENTRYPOINT ["/deploy.sh"]
